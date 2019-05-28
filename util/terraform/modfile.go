package terraform

import (
	"fmt"
	"github.com/hashicorp/hcl"
	//"github.com/hashicorp/hcl/hcl/printer"
	//jsonParser "github.com/hashicorp/hcl/json/parser"
	"encoding/json"
	//"github.com/hashicorp/terraform/config"
	"gopkg.in/src-d/go-billy.v4"
	"io/ioutil"
	//  "os"
	"log"
	"path/filepath"
)

func walk(fs billy.Filesystem, fullPath string, paths []string) []string {
	children, _ := fs.ReadDir(fullPath)
	for _, fi := range children {
		var name string
		if fullPath == "." {
			name = fi.Name()
		} else {
			name = fullPath + "/" + fi.Name()
		}

		if fi.IsDir() {
			paths = walk(fs, name, paths)
		} else {
			paths = append(paths, name)
		}
	}
	return paths
}

func DiscoverResource(fs billy.Filesystem, fullPath string, resource []StateResources) {

	var files []string
	var empty []string

	files = walk(fs, ".", empty)

	for _, file := range files {
		if filepath.Ext(file) == ".tf" {
			fmt.Println(file)
			f, err := fs.Open(file)
			if err != nil {
				log.Fatalf("Unable to open file: %v", err)
			}

			fi, err := ioutil.ReadAll(f)
			if err != nil {
				log.Fatalf("Unable to read file: %v", err)
			}

			js, err := hclToJson(fi)
			if err != nil {
				log.Fatalf("Unable to convert HCL to JSON: %v", err)
			}
			//fmt.Println(string(js))
			var result map[string]interface{}
			json.Unmarshal([]byte(js), &result)

			for t1 := range result {
				switch t1 {
				case "resource":
					birds := result[t1].([]interface{})
					for key, value := range birds {
						fmt.Println(key)
            //fmt.Println(value)
            for k, v := range value.(map[string]interface{}) {
              fmt.Println(k)
              fmt.Println(v)
            }
            fmt.Println("=========================")
					}
				}
			}

		}
	}
}

/*
func DiscoverResource(fs billy.Filesystem, files []string, resource []StateResources) {
	for _, file := range files {
		f, err := fs.Open(file)
		if err != nil {
			log.Fatalf("Unable to open file: %v", err)
		}

		fi, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		if filepath.Ext(file) == ".tf" {

				hcl, err := hcl.Parse(string(fi))
				if err != nil {
					panic(err)
				}
				fmt.Println(hcl)

			var t interface{}
			hcl.Decode(&t, string(fi))
			if err != nil {
				panic(err)
			}
			fmt.Println(t.module)
		}

			if filepath.Ext(file) == ".tf" {
				tfJson, err := hclToJson(fi)
				if err != nil {
					panic(err)
				}

				var result map[string]interface{}
				json.Unmarshal(tfJson, &result)

				for key, _ := range result {
					switch key {
					case "resource":
						res1 := result["resource"].([]interface{})
						for key1, _ := range res1 {
							fmt.Println(key1)
							//fmt.Println(val1)
						}
					}

				}

			}

	}
}
*/

func hclToJson(tfSource []byte) ([]byte, error) {
	var v interface{}
	err := hcl.Unmarshal(tfSource, &v)
	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json, nil
}
