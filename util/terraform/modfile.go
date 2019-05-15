package terraform

import (
	"bufio"
	"fmt"
	"gopkg.in/src-d/go-billy.v4"
	"log"
	"strconv"
	"strings"
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

/*
func StringInFile(file *File, searchString string) {
}
*/
func getSearchString() {

}

func FindStringInFiles(fs billy.Filesystem, files []string, searchString string) string {
	for _, file := range files {
		f, err := fs.Open(file)
		if err != nil {
			log.Fatalf("Unable to open file: %v", err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		line := 1
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), searchString) {
				return file
			}
			line++
		}
	}
	return ""
}

/*
func FindResourceAttribute(files *Files, resourceType string, resourceName string, attribute) {
	searchString := resourceType + resourceName
	fileWithResource := FindStringInFiles(files, searchString)
	if len(GetOccurencesInFile(file, attributeString)) == 1 {
	    ReplaceStringInFile(file, oldValue newValue)
	}
	else {
		print("Multiple strings not yet supported.")
	}
}
*/

func Unpackfile(fs billy.Filesystem, file string, inputs []StateResources) {
	//fmt.Println(file)
	f, err := fs.Open(file)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	line := 1
	for scanner.Scan() {
		for _, ip := range inputs {
			resourceType := strings.Split(ip.TFResource, ".")[0]
			resourceName := strings.Join(strings.Split(ip.TFResource, ".")[1:], ".")
			searchString := "resource " + strconv.Quote(resourceType) + " " + strconv.Quote(resourceName)

			if strings.Contains(scanner.Text(), searchString) {
				log.Printf("Found matching resource in file %v", file, " at line ", line, " matching resource", searchString)
				fmt.Println(file)
				fmt.Println(ip.OldValue)
			}
			line++
		}
	}

	//b, err := ioutil.ReadAll(f)
	//if err != nil {
	//	log.Fatalf("Unable to open file: %v", err)
	//}
	//fmt.Println(file)
	//fmt.Println(string(b))

}
