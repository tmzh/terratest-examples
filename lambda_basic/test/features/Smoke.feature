Feature: Simple test to confirm lambda function behavior
	Confirms that given a valid terraform variable
	Lambda resources are deployed
	The Lambda function executes as intended
	Scenario: Deploy a Lambda function
		Given Terraform code is deployed with these variables:
			|function_name | random_name|
		Then For given inputs Lambda function output is as expected:
			|world | "Hello world!"|
		Then Cloudwatch log stream is generated
