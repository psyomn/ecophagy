/*
Package barf contains code relevant to barfing.

A lot of this can be refactored, but this is not the point of this
sub-project for now.

Copyright 2019 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package barf

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/psyomn/ecophagy/psy/common"
)

func cmake(args common.RunParams) common.RunReturn {
	if len(args) == 0 {
		return ErrNeedProjectName
	}

	type project struct {
		ProjectName string
	}

	p := project{args[0]}

	const (
		cmakeFile = `cmake_minimum_required(VERSION 3.9)
project({{.ProjectName}})

# Took some many of these parts for cmake off
#   https://github.com/RAttab/optics

enable_testing()

add_definitions("-Wall")
add_definitions("-Wextra")
add_definitions("-Wundef")
add_definitions("-Wformat=2")
add_definitions("-Winit-self")
add_definitions("-Wcast-align")
add_definitions("-Wswitch-enum")
add_definitions("-Wwrite-strings")
add_definitions("-Wswitch-default")
add_definitions("-Wunreachable-code")
add_definitions("-Wno-strict-aliasing")
add_definitions("-Wno-format-nonliteral")
add_definitions("-Wno-missing-field-initializers")
add_definitions("-pipe -g -O3 -Werror -march=native")

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -g -std=gnu11")

set({{.ProjectName}}_SOURCES
  src/helper.c)

add_library({{.ProjectName}} STATIC ${{"{"}}{{.ProjectName}}_SOURCES{{"}"}})

include_directories(include)

function({{.ProjectName}}_add_test name)
  add_executable(${name} test/${name}.c)
  target_link_libraries(${name} {{.ProjectName}})
  add_test(
    NAME    ${name}_valgrind
    COMMAND valgrind --leak-check=full
                     --error-exitcode=1 $<TARGET_FILE:${name}>)
  add_test(${name} ${name})
endfunction({{.ProjectName}}_add_test)

{{.ProjectName}}_add_test({{.ProjectName}}_example)
`
		mainProgram = `#include "helper.h"
int main(int argc, char* argv[]) {
  return {{.ProjectName}}_add(0,0);
}`

		helperProgramHeader = `#pragma once
int {{.ProjectName}}_add(int a, int b);`

		helperProgram = `#include <{{.ProjectName}}/helper.h>
int {{.ProjectName}}_add(int a, int b) {
  return a + b;
}`
		helperTestExample = `#include <{{.ProjectName}}/test.h>
#include <{{.ProjectName}}/helper.h>

#include <stdio.h>
#include <assert.h>

int some_test(void **data)
{
  (void) data;
  return 0;
}

int main(void)
{
  {{.ProjectName}}_test("some test", some_test, NULL);
  return 0;
}`

		testConfig = `#pragma once

#include <stdio.h>
#include <time.h>
#include <stdio.h>

#define {{.ProjectName}}_test(l, fn, dt) internal_{{.ProjectName}}_test(__FILE__ ": " l, fn, dt)

void internal_{{.ProjectName}}_test(const char* label, int (*func)(void **data), void **data)
{
  printf("%s:", label);

  const clock_t start = clock();
  const time_t time_start = time(NULL);
  const int ret = func(data);
  const clock_t end = clock();
  const time_t time_end = time(NULL);
  const double elapsed = (end - start) / (double) CLOCKS_PER_SEC;
  const size_t elapsed_time = (time_end - time_start);

  fprintf(stdout, " %s [wc:%f][tm:%zu]\n", !ret ? "ok" : "fail", elapsed, elapsed_time);
}`
	)

	err := os.MkdirAll(p.ProjectName, 0755)
	if err != nil {
		return err
	}

	srcDir := filepath.Join(p.ProjectName, "src")
	err = os.MkdirAll(srcDir, 0755)
	if err != nil {
		return err
	}

	incDir := filepath.Join(p.ProjectName, "include", p.ProjectName)
	err = os.MkdirAll(incDir, 0755)
	if err != nil {
		return err
	}

	{ // helper.h, helper.c
		helperHeaderFile, err := os.Create(filepath.Join(incDir, "helper.h"))
		if err != nil {
			return err
		}
		defer helperHeaderFile.Close()
		t := template.Must(
			template.New("helper-header").Parse(helperProgramHeader),
		)
		helperHeaderBuff := bytes.NewBuffer([]byte{})

		if err := t.Execute(helperHeaderBuff, p); err != nil {
			return err
		}

		if _, err := helperHeaderFile.Write(helperHeaderBuff.Bytes()); err != nil {
			return err
		}

		helperSourceFile, err := os.Create(filepath.Join(srcDir, "helper.c"))
		if err != nil {
			return err
		}
		defer helperSourceFile.Close()
		t2 := template.Must(
			template.New("helper-source").Parse(helperProgram),
		)
		helperSourceBuff := bytes.NewBuffer([]byte{})
		err = t2.Execute(helperSourceBuff, p)
		if err != nil {
			return err
		}

		if _, err := helperSourceFile.Write(helperSourceBuff.Bytes()); err != nil {
			return err
		}
	}

	{ // test/examples
		if err := os.MkdirAll(incDir, 0755); err != nil {
			return err
		}

		testHelper, err := os.Create(filepath.Join(incDir, "test.h"))
		if err != nil {
			return err
		}
		defer testHelper.Close()

		t := template.Must(
			template.New("test-config").Parse(testConfig),
		)

		testBuff := bytes.NewBuffer([]byte{})
		err = t.Execute(testBuff, p)
		if err != nil {
			return err
		}

		if _, err := testHelper.Write(testBuff.Bytes()); err != nil {
			return err
		}

		testSrcDir := filepath.Join(p.ProjectName, "test")
		err = os.MkdirAll(testSrcDir, 0755)
		if err != nil {
			return err
		}

		t2 := template.Must(
			template.New("test-sample-config").Parse(helperTestExample),
		)

		sampleBuff := bytes.NewBuffer([]byte{})
		err = t2.Execute(sampleBuff, p)
		if err != nil {
			return err
		}

		testSampleFile, err := os.Create(
			filepath.Join(testSrcDir, p.ProjectName+"_example.c"),
		)
		if err != nil {
			return err
		}
		defer testSampleFile.Close()

		if _, err := testSampleFile.Write(sampleBuff.Bytes()); err != nil {
			return err
		}
	}

	{ // main
		mainFile, err := os.Create(filepath.Join(srcDir, "main.c"))
		if err != nil {
			return err
		}
		defer mainFile.Close()

		t := template.Must(
			template.New("main-program").Parse(mainProgram),
		)
		mainBuff := bytes.NewBuffer([]byte{})
		err = t.Execute(mainBuff, p)
		if err != nil {
			return err
		}

		if _, err := mainFile.Write(mainBuff.Bytes()); err != nil {
			return err
		}
	}

	cmakeFileHandle, err := os.Create(filepath.Join(p.ProjectName, "CMakeLists.txt"))
	if err != nil {
		return err
	}
	defer cmakeFileHandle.Close()

	t := template.Must(
		template.New("cmake-project").Parse(cmakeFile),
	)

	cmakeBuff := bytes.NewBuffer([]byte{})
	err = t.Execute(cmakeBuff, p)
	if err != nil {
		return err
	}

	if _, err := cmakeFileHandle.Write(cmakeBuff.Bytes()); err != nil {
		return err
	}

	fmt.Println("cmake project barfed successfully")

	return nil
}
