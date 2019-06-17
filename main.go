package main

import (
	"fmt"
	"github.com/garybowers/gcp_recommender/inputs/recommender"
	"github.com/garybowers/gcp_recommender/util/git"
	"github.com/garybowers/gcp_recommender/util/storage"
	"github.com/garybowers/gcp_recommender/util/terraform"
	//"io/ioutil"
	"log"
	"os"
)

func main() {
	// 1. Get the recommendations from the recommendation input.
	//		See inputs/ folder for supported inputs
	//		TODO: Make things other than GCP Recommendation engine work.

	gitUrl := os.Args[1]
	gitUser := os.Args[2]
	gitToken := os.Args[3]
	gcpProject := os.Args[4]
	gcpZone := os.Args[5]
	gcpCreds := os.Args[6]
	tfStateBucket := os.Args[7]
	tfStatePath := os.Args[8]

	GCPFindings := recommender.GetFindings(gcpCreds,
		gcpProject,
		gcpZone,
		"google.compute.instanceGroupManager.MachineTypeRecommender")

	// 2. Parase the terraform state file, below we use GCS bucket in a central project to store all terraform state files.  You can also use local storage.
	//	TODO: Implement S3 storage -- i hate aws & implement local storage

	tfState, err := storage.ReadObject(tfStateBucket, tfStatePath)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	state, err := terraform.ParseTerraformState(tfState)
	if err != nil {
		log.Fatalf("Unable to parse terraform state: %v", err)
	}

	// 3. Find the recommended change in the state, this will give us a clue as to what resource block to change and where in the hierarchy it lives.
  var findings []recommender.Findings
  findings = GCPFindings
  for _, finding := range findings {
    log.Println(finding.Instance)
   }

	resources := terraform.FindResources(state, GCPFindings)
  fmt.Println(resources)

	// 4. Get the terraform code repo from our git repository, token, username and repo can come from various places.
	// At KPMG as part of the project creation we put a breadcrumb in the project metadata, we have a base64 hash of the git repo url.
	fs, repo := git.Clone(gitUrl,
		gitUser,
		gitToken)
	fmt.Println(repo)

	//var empty []string
	//paths := terraform.Walk(fs, ".", empty)

	terraform.DiscoverResource(fs, ".", resources)

}
