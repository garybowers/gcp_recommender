# Recommendation actioner

## Purpose

The application takes recommendations from various 'sensors' or inputs and modifies infrastructure as code via terraform.

The basic flow of the application is as follows:

1. Take the input recommendation
2. Clone the source code repository into memory
3. Walk the terraform state file and find the resource block of the resource in question
4. Modify the resource parameters based on the recommendation
5. Commit the code and raise a Merge Request

Usage:
```
  ./recommender \
  		<url_to_terraform_repo> \
		<git_username> \
		<git_token> \
		<gcp_project> \
		<gcp_zone> \
		<gcp_credentials_json_file> \
		<tf_state_gcs_bucket_name> \
		<tf_state_path_in_bucket>
```

#### Note:
Only git over https:// url schemes are supported at the moment

There's lots to do with this to make it modular / plugin oriented
