package terraform

import (
	"fmt"
	"github.com/hashicorp/hcl"
	//"github.com/hashicorp/hcl/hcl/printer"
	//jsonParser "github.com/hashicorp/hcl/json/parser"
	"encoding/json"
	"gopkg.in/src-d/go-billy.v4"
	"io/ioutil"
	"log"
	"path/filepath"
)

func Walk(fs billy.Filesystem, fullPath string, paths []string) []string {
	children, _ := fs.ReadDir(fullPath)
	for _, fi := range children {
		var name string
		if fullPath == "." {
			name = fi.Name()
		} else {
			name = fullPath + "/" + fi.Name()
		}

		if fi.IsDir() {
			paths = Walk(fs, name, paths)
		} else {
			paths = append(paths, name)
		}
	}
	return paths
}

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
			json := hclToJson(fi)
			fmt.Println(string(json))
		}
	}
}

func hclToJson(tfSource []byte) string {
	//fmt.Println(tfSource)
	var v interface{}
	err := hcl.Unmarshal(tfSource, &v)
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func jsonToHcl() {

}
