/*
Package barf contains code relevant to barfing.

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

const (
	gpr = `-- Generated Gnat file
-- Example use:
--   gprbuild -P {{.ProjectName}} -Xmode=debug -p
project {{.ProjectName}} is

   -- Standard configurations
   for Main        use ("main.adb");
   for Source_Dirs use ("src/**");
   for Exec_Dir    use "bin/";

   -- Ignore git scm stuff
   for Ignore_Source_Sub_Dirs use (".git/");

   for Object_Dir use "obj/" & external ("mode", "debug");
   for Object_Dir use "obj/" & external ("mode", "release");

   package Builder is
      for Executable ("main.adb") use "{{.ProjectName}}";

   end Builder;

   -- To invoke either case, you need to set the -X flag at gnatmake in command
   -- line. You will also notice the Mode_Type type. This constrains the values
   -- of possible valid flags.
   type Mode_Type is ("debug", "release");
   Mode : Mode_Type := external ("mode", "debug");
   package Compiler is
      -- Either debug or release mode
      case Mode is
         when "debug" =>
            for Switches ("Ada") use ("-g");
         when "release" =>
            for Switches ("Ada") use ("-O2");
      end case;
   end Compiler;

   package Binder is end Binder;

   package Linker is end Linker;

end {{.ProjectName}};`

	mainProgram = `with Ada.Text_IO;
procedure Main is begin
   Ada.Text_IO.Put_Line ("hello world");
end Main;
`
)

func ada(args common.RunParams) common.RunReturn {
	if len(args) == 0 {
		return ErrNeedProjectName
	}

	type project struct {
		ProjectName string
	}

	t := template.Must(
		template.New("ada-project").Parse(gpr),
	)
	p := project{args[0]}

	err := os.MkdirAll(p.ProjectName, 0755)
	if err != nil {
		return err
	}

	srcDir := filepath.Join(p.ProjectName, "src")
	err = os.MkdirAll(srcDir, 0755)
	if err != nil {
		return err
	}

	gprBuff := bytes.NewBuffer([]byte{})
	err = t.Execute(gprBuff, p)
	if err != nil {
		return err
	}

	gprFilePath := filepath.Join(p.ProjectName, p.ProjectName+".gpr")
	file, err := os.Create(gprFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(gprBuff.Bytes()); err != nil {
		return err
	}

	mainFilePath := filepath.Join(srcDir, "main.adb")
	file, err = os.Create(mainFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(mainProgram); err != nil {
		return err
	}

	fmt.Println("ada project barfed successfully")

	return nil
}
