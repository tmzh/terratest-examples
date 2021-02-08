package test

import (
	"fmt"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformAwsLambdaFunction(t *testing.T) {
	t.Parallel()

	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "..",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Invoke the function, so we can test its output
	response := aws.InvokeFunction(t, awsRegion, functionName, "World")

	assert.Equal(t, `"hi!"`, string(response))

}
