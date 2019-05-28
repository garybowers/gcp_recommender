package terraform

import (
	"github.com/garybowers/gcp_recommender/inputs/recommender"
	"log"
)

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
