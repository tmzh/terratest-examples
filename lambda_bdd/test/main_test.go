package test

import (
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	//	"github.com/stretchr/testify/require"
	"testing"
)

type Payload struct {
	Name string
}

func TestTerraformAwsLambdaFunction(t *testing.T) {
	t.Parallel()

	awsRegion := "us-east-1"
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "..",
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
	functionName := terraform.Output(t, terraformOptions, "lambda_function")

	// Invoke the function, so we can test its output
	response := aws.InvokeFunction(t, awsRegion, functionName, Payload{Name: "World"})

	assert.Equal(t, `"Hello World!"`, string(response))

}
