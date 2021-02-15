package test

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"

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

	godog.TestSuite{
		Name:                 "LambdaTest",
		TestSuiteInitializer: ts.InitializeTestSuite,
		ScenarioInitializer:  ts.InitializeScenario,
		Options:              &opts,
	}.Run()

}

func (ts testingSuite) InitializeTestSuite(ctx *godog.TestSuiteContext) {
}

func (ts testingSuite) InitializeScenario(ctx *godog.ScenarioContext) {
	o := &godogFeaturesScenario{}
	o.testing = ts.testing

	ctx.Step(`^Terraform code is deployed with these variables:$`, o.terraformIsDeployedWithVariables)
	ctx.Step(`^For given inputs Lambda function output is as expected:$`, o.givenInputsLambdaReturnsValuesAsExpected)
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
	return nil
}

func (o *godogFeaturesScenario) givenInputsLambdaReturnsValuesAsExpected(tbl *godog.Table) error {
	functionName := terraform.Output(o.testing, o.terraformOptions, "lambda_function")
	awsRegion := "us-east-1"

	for _, row := range tbl.Rows {
		input := row.Cells[0].Value
		expected := row.Cells[1].Value
		response := aws.InvokeFunction(o.testing, awsRegion, functionName, Payload{Name: input})
		actual := string(response)
		if expected != actual {
			return fmt.Errorf("Not equal: \n"+
				"expected: %s\n"+
				"actual  : %s", expected, actual)
		}
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
}

func (o *godogFeaturesScenario) destroyTerraform(sc *godog.Scenario, err error) {
	terraform.Destroy(o.testing, o.terraformOptions)
}
