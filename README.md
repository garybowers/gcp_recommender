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



#### To Do List:

1. inputs/ folder needs to move to a plugin architecture, allowing inputs from outside recommendation engine to feed in changes, e.g. from prometheus metrics or forseti.

2. terraform hcl modification is awful at the moment, it's scanning whole files for string matches.  This is problematic as if the recommendation engine comes back with a change from n1-standard-4 we may have multiple resources of this type.

	Answer is to build a abstract syntax tree to find actual resource blocks to modify the resource in question.
	
3.	I need to write tests ..oops

4.	Some terraform states are in s3, i hate aws and really don't want to look into their API's


