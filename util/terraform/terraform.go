package terraform

import (
	"github.com/garybowers/recommender/inputs/recommender"
	"log"
)

//func GetModulePath(pathComponents []string, currentPath string) string {
//	if len(pathComponents) == 0 {
//		return currentPath
//	} else {
//		targetModule := pathComponents[0]
//		localTerraformFiles := GetLocalTerraformFiles(currentPath)
//		relativeModulePath := GetModuleSourceDirectory(targetModule, localTerraformFiles)
//		fullModulePath := path.Join(currentPath, relativeModulePath)
//		reducedPathComponents := pathComponents[:len(pathComponents)-1]
//		return GetModulePath(reducedPathComponents, fullModulePath)
//	}
//}
/*
func GetLocalTerraformFiles(path string) []string {
	// TODO
	localTerraformFiles := make([]string, 0)
	return localTerraformFiles
}

func GetModuleSourceDirectory(localTerraformFilePaths []string) string {
	// TODO
	return ""
}
*/
//GetModulePath(["root", "foo", "bar", "baz"], "/gary/terraform")
//GetModulePath(["foo", "bar", "baz"], "/gary/terraform/")
//GetModulePath(["bar", "baz"], "/gary/terraform/foo")
//GetModulePath(["baz"], "/gary/terraform/foo/bar2")

type StateResources struct {
	ResourceType     string
	ResourceName     string
	ResourceLocation []string
	OldValue         string
	NewValue         string
	TFResource       string
}

func FindResources(State *TerraformState, Findings []recommender.Findings) []StateResources {
	log.Printf("Searching terraform state for matching resource")
	resources := make([]StateResources, 0, 2)

	for _, findings := range Findings {
		for _, module := range State.Modules {
			for tfResource, resource := range module.Resources {
				if resource.Type == findings.InstanceType {
					log.Printf("Found recommendation for %v", findings.Instance)
					if resource.Primary.ID == findings.Instance {
						resources = append(resources, StateResources{
							resource.Type,
							resource.Primary.ID,
							module.Path,
							findings.InstanceSizeFrom,
							findings.InstanceSizeTo,
							tfResource})
						break
					}
				}
			}
		}
	}

	return resources
}
