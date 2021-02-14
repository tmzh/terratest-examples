package test

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

type testingSuite struct {
	testing *testing.T
}

type godogFeaturesScenario struct {
	testing          *testing.T
	terraformOptions *terraform.Options
}

func TestLambdaFunctionBDD(t *testing.T) {
	t.Parallel()

	opts := godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(),
	}

	ts := &testingSuite{}
	ts.testing = t

	status := godog.TestSuite{
		Name:                 "LambdaTest",
		TestSuiteInitializer: ts.InitializeTestSuite,
		ScenarioInitializer:  ts.InitializeScenario,
		Options:              &opts,
	}.Run()

	fmt.Println(status)
}

func (ts testingSuite) InitializeTestSuite(ctx *godog.TestSuiteContext) {
}

func (ts testingSuite) InitializeScenario(ctx *godog.ScenarioContext) {
	o := &godogFeaturesScenario{}
	o.testing = ts.testing

	ctx.Step(`^Terraform code is deployed with these variables:$`, o.terraformIsDeployedWithVariables)
	ctx.AfterScenario(o.destroyTerraform)
}

func (o *godogFeaturesScenario) terraformIsDeployedWithVariables(tbl *godog.Table) error {
	tfVars := make(map[string]interface{})
	for _, row := range tbl.Rows {
		tfVars[row.Cells[0].Value] = row.Cells[1].Value
	}

	awsRegion := "us-east-1"

	terraformOptions := terraform.WithDefaultRetryableErrors(o.testing, &terraform.Options{
		TerraformDir: "..",
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	})

	o.terraformOptions = terraformOptions

	terraform.InitAndApply(o.testing, terraformOptions)
	functionName := terraform.Output(o.testing, terraformOptions, "lambda_function")

	// Invoke the function, so we can test its output
	response := aws.InvokeFunction(o.testing, awsRegion, functionName, Payload{Name: "World"})

	assert.Equal(o.testing, `"Hello World!"`, string(response))
	return nil
}

func (o *godogFeaturesScenario) destroyTerraform(sc *godog.Scenario, err error) {
	terraform.Destroy(o.testing, o.terraformOptions)
}
