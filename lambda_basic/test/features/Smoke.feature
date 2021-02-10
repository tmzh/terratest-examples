Feature: Simple test to confirm lambda function behavior
	Confirms that given a valid terraform variable
	Lambda resources are deployed
	The Lambda function can be executes as intended
	Scenario: Deploy a Lambda function
		Given Terraform code is deployed with these variables:
			|function_name | random_name|
		When Lambda function is invoked the following input:
			|function_name | random_name|
		Then the response is successful
