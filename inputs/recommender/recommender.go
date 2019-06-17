// Connect to the recommendation engine api (alpha) and get the recommendations to modify the terraform code.

// TODO: At the moment this only supports Instances and InstanceTemplates
//		 This needs to be made more generic to support future recommendations
//		 i.e. IAM Recommendations.
//		 GetFindings func needs rework

package recommender

import (
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"encoding/json"
	"path/filepath"
	"strings"
)

type recommendation struct {
	Recommendations []struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		PrimaryImpact struct {
			Category       string `json:"category"`
			CostProjection struct {
				Cost struct {
					CurrencyCode string `json:"currencyCode"`
					Units        string `json:"units"`
				}
				Duration string `json:"duration"`
			}
		}
		Content struct {
			Justification struct {
				Summary string `json:"summary"`
			}
			OperationGroups []struct {
				Operations []struct {
					Action       string `json:"action"`
					ResourceType string `json:"resourceType"`
					Resource     string `json:"resource"`
				} `json:"operations"`
			} `json:"operationGroups"`
		}
		StateInfo struct {
			State string `json:"state"`
		}
		Etag            string `json:"etag"`
		LastRefreshTime string `json:"lastRefreshTime"`
	} `json:"recommendations"`
}

type Findings struct {
	InstanceType     string
	Instance         string
	InstanceSizeFrom string
	InstanceSizeTo   string
	Currency         string
	CostDifference   string
	Status           string
	Justification    string
}

func getRecommendations(credentialsFile string, projectId string, location string, recommenderId string) (body []byte) {

	url := "https://recommender.googleapis.com/v1alpha1/projects/" + projectId + "/locations/" + location + "/recommenders/" + recommenderId + "/recommendations"

	log.Println("Getting recommendations from ", url)

	ctx := context.Background()

	cred, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(cred, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Fatalf("Unable to parse credentials file: %v", err)
	}

	client := config.Client(ctx)

	request, err := http.NewRequest("GET", url, nil)
	request = request.WithContext(ctx)

	resp, err := client.Do(request)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)
  //log.Println(string(b))

	return b
}

func unmarshal(f []byte) (*recommendation, error) {
	recommendation := &recommendation{}

	if err := json.Unmarshal(f, recommendation); err != nil {
		return nil, err
	}
	return recommendation, nil
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func GetFindings(credentialsFile string, projectId string, location string, recommenderId string) []Findings {
	f := getRecommendations(credentialsFile, projectId, location, recommenderId)

	finding, err := unmarshal(f)
	if err != nil {
		log.Fatalf("Unable to parse recommender api: %v", err)
	}

	found := make([]Findings, 0, 2)

	for i := range finding.Recommendations {
		for og := range finding.Recommendations[i].Content.OperationGroups {
			for op := range finding.Recommendations[i].Content.OperationGroups[og].Operations {
				switch finding.Recommendations[i].Content.OperationGroups[og].Operations[op].ResourceType {
				case "compute.googleapis.com/InstanceTemplate":
					if finding.Recommendations[i].Content.OperationGroups[og].Operations[op].Action == "test" {
						instanceSizeFrom := strings.TrimSpace(between(finding.Recommendations[i].Description, "from", "to"))
						instanceSizeTo := strings.TrimSpace(between(finding.Recommendations[i].Description, "to", "."))

						found = append(found, Findings{
							//finding.Recommendations[i].Content.OperationGroups[og].Operations[op].ResourceType,
							"google_compute_instance_template",
							filepath.Base(finding.Recommendations[i].Content.OperationGroups[og].Operations[op].Resource),
							instanceSizeFrom,
							instanceSizeTo,
							finding.Recommendations[i].PrimaryImpact.CostProjection.Cost.Units,
							finding.Recommendations[i].PrimaryImpact.CostProjection.Cost.CurrencyCode,
							finding.Recommendations[i].StateInfo.State,
							finding.Recommendations[i].Content.Justification.Summary})
						break
					}
				case "compute.googleapis.com/Instance":
					instanceSizeFrom := strings.TrimSpace(between(finding.Recommendations[i].Description, "from", "to"))
					instanceSizeTo := strings.TrimSpace(between(finding.Recommendations[i].Description, "to", "."))
					found = append(found, Findings{
						//finding.Recommendations[i].Content.OperationGroups[og].Operations[op].ResourceType,
						"google_compute_instance",
						filepath.Base(finding.Recommendations[i].Content.OperationGroups[og].Operations[op].Resource),
						instanceSizeFrom,
						instanceSizeTo,
						finding.Recommendations[i].PrimaryImpact.CostProjection.Cost.Units,
						finding.Recommendations[i].PrimaryImpact.CostProjection.Cost.CurrencyCode,
						finding.Recommendations[i].StateInfo.State,
						finding.Recommendations[i].Content.Justification.Summary})
					break
				}
			}
		}
	}
	return found
}
