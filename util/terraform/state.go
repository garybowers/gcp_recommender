package terraform

import (
	"encoding/json"
	//"fmt"
	"log"
	"sync"
	//"github.com/hashicorp/terraform/terraform"
)

type TerraformBackend struct {
	Type   string
	Config map[string]interface{}
}

type TerraformState struct {
	Version int
	Serial  int
	Backend *TerraformBackend
	Modules []TerraformModuleState
}

type TerraformModuleState struct {
	Path         []string                           `json:"path"`
	Locals       map[string]interface{}             `json:"-"`
	Outputs      map[string]*TerraformOutputState   `json:"outputs"`
	Resources    map[string]*TerraformResourceState `json:"resources"`
	Dependencies []string                           `json:"depends_on"`
}

type TerraformResourceState struct {
	Type         string                    `json:"type"`
	Dependencies []string                  `json:"depends_on"`
	Primary      *TerraformInstanceState   `json:"primary"`
	Deposed      []*TerraformInstanceState `json:"deposed"`
	Provider     string                    `json:"provider"`
}

type TerraformInstanceState struct {
	ID         string                 `json:"id"`
	Attributes map[string]string      `json:"attributes"`
	Ephemeral  EphemeralState         `json:"-"`
	Meta       map[string]interface{} `json:"meta"`
}

type TerraformOutputState struct {
	Sensitive bool        `json:"sensitive"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`

	mu sync.Mutex
}

type EphemeralState struct {
	ConnInfo map[string]string `json:"-"`
	Type     string            `json:"-"`
}

/*
func (state *TerraformState) Print() {
	fmt.Println(state.Version)
}

func (state *TerraformState) String() string {
	//fmt.Println(state.Modules)
	for _, module := range state.Modules {
		//fmt.Println(strings.Join(module.Path, "/"))
		fmt.Println(module.Path)
		for tfResource, resource := range module.Resources {
			fmt.Println(tfResource)
			fmt.Println(resource.Primary.ID)
			//		fmt.Println(resource.Primary.Attributes)
		}
	}
	return fmt.Sprintf("%d", state.Version)
}
*/
func ParseTerraformState(terraformStateData []byte) (*TerraformState, error) {
	log.Printf("Parsing terraform state")
	terraformState := &TerraformState{}

	if err := json.Unmarshal(terraformStateData, terraformState); err != nil {
		return nil, err
	}
	return terraformState, nil
}
