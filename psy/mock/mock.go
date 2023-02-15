/*
Package mock will mock tcp/udp endpoints

Copyright 2020-2022 Simon Symeonidis (psyomn)

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
package mock

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"sync"

	"github.com/psyomn/ecophagy/psy/common"

	"github.com/go-yaml/yaml"
)

type record struct {
	Type   string      `yaml:"type"`
	Port   int         `yaml:"port"`
	Return interface{} `yaml:"return"`

	// only relevant for http endpoints
	Root string `yaml:"root"`
}

type config map[string]record

func usage(fs *flag.FlagSet) common.RunReturn {
	fs.Usage()
	return ErrWrongCmdUsage
}

// Run net mocker
func Run(args common.RunParams) common.RunReturn {
	t := &config{}

	type session struct {
		generate string
		config   string
	}
	sess := session{}

	mockCmd := flag.NewFlagSet("mock", flag.ExitOnError)
	mockCmd.StringVar(&sess.generate, "generate", sess.generate, "generate a sample config file")
	mockCmd.StringVar(&sess.config, "config", sess.config, "use the config to run the server")
	if err := mockCmd.Parse(args); err != nil {
		return err
	}

	if sess.generate != "" {
		return generateYamlConfig(sess.generate)
	}

	if sess.config == "" {
		return usage(mockCmd)
	}

	configContents, err := readFile(sess.config)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configContents, t)
	if err != nil {
		return err
	}

	return processEntries(t)
}

func processEntries(conf *config) error {
	var wg sync.WaitGroup

	for _, v := range *conf {
		switch v.Type {
		case "udp":
			wg.Add(1)
			go createUDP(v.Port, v.Return, &wg)
		case "tcp":
			wg.Add(1)
			go createTCP(v.Port, v.Return, &wg)

		case "http":
			wg.Add(1)
			go createHTTP(v.Port, v.Return, v.Root, &wg)
		default:
			return fmt.Errorf("%w: %v", ErrUnknownService, v.Type)
		}
	}

	wg.Wait()

	return nil
}

func createUDP(port int, ret interface{}, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	buf := make([]byte, 1024)
	portStr := fmt.Sprintf(":%d", port)
	pc, err := net.ListenPacket("udp", portStr)
	log.Println(err)

	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Println("could not read udp packet: ", err)
		}

		log.Println(n, addr, string(buf[:n]))

		if ret != nil {
			val := processReturn(ret)

			_, err := pc.WriteTo(val, addr)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func createTCP(port int, ret interface{}, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	portStr := fmt.Sprintf(":%d", port)

	l, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Println("error:", err)
		return
	}
	defer l.Close()

	var buff [1024]byte
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("error:", err)
			return
		}

		n, err := conn.Read(buff[:])
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("read num bytes:", n)
		if ret != nil {
			val := processReturn(ret)
			_, err := conn.Write(val)
			if err != nil {
				log.Println(err)
			}
		}
		conn.Close()
	}
}

func createHTTP(port int, ret interface{}, root string, wg *sync.WaitGroup) {
	val := processReturn(ret)
	server := http.NewServeMux()

	server.HandleFunc(root, func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write(val); err != nil {
			log.Println(err)
		}
	})

	go func(wait *sync.WaitGroup) {
		if wait != nil {
			defer wg.Done()
		}

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
	}(wg)
}

func processReturn(value interface{}) []byte {
	var ret []byte
	switch v := value.(type) {
	case string:
		ret = []byte(v)
	case []interface{}:
		ret = make([]byte, len(v))
		for i := range v {
			if val, ok := v[i].(int); ok {
				ret[i] = byte(val)
			} else {
				log.Fatal("not a byte array")
			}
		}
	case []byte:
		ret = v
	case []int:
		log.Fatal("array in return should be byte magnitude")
	default:
		log.Fatal("that type is not supported:", reflect.TypeOf(v).String())
	}
	return ret
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	return data, err
}

func generateYamlConfig(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(example)
	return err
}
