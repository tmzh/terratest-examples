package test

import (
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"testing"
)

type Payload struct {
	Name string
}

func TestTerraformAwsLambdaFunction(t *testing.T) {
	t.Parallel()

	awsRegion := "us-east-1"
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "..",
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
	functionName := terraform.Output(t, terraformOptions, "lambda_function")

	response := aws.InvokeFunction(t, awsRegion, functionName, Payload{Name: "World"})

	assert.Equal(t, `"Hello World!"`, string(response))

}
