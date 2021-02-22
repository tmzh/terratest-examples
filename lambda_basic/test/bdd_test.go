package test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
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
	stepValues       map[string]string
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
	o.stepValues = make(map[string]string)

	ctx.Step(`^Terraform code is deployed with these variables:$`, o.terraformIsDeployedWithVariables)
	ctx.Step(`^For given inputs Lambda function output is as expected:$`, o.givenInputsLambdaReturnsValuesAsExpected)
	ctx.Step(`^Cloudwatch log stream is generated$`, o.cloudwatchLogIsGenerated)
	ctx.AfterScenario(o.destroyTerraform)
}

func (o *godogFeaturesScenario) terraformIsDeployedWithVariables(tbl *godog.Table) error {
	tfVars := make(map[string]interface{})
	for _, row := range tbl.Rows {
		tfVars[row.Cells[0].Value] = row.Cells[1].Value
	}
	o.stepValues["awsRegion"] = "us-east-1"

	terraformOptions := terraform.WithDefaultRetryableErrors(o.testing, &terraform.Options{
		TerraformDir: "..",
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": o.stepValues["awsRegion"],
		},
	})

	o.terraformOptions = terraformOptions
	terraform.InitAndApply(o.testing, terraformOptions)
	return nil
}

func (o *godogFeaturesScenario) givenInputsLambdaReturnsValuesAsExpected(tbl *godog.Table) error {
	o.stepValues["functionName"] = terraform.Output(o.testing, o.terraformOptions, "lambda_function")
	for _, row := range tbl.Rows {
		input := row.Cells[0].Value
		expected := row.Cells[1].Value
		response := aws.InvokeFunction(o.testing, o.stepValues["awsRegion"], o.stepValues["functionName"], Payload{Name: input})
		actual := string(response)
		if expected != actual {
			return fmt.Errorf("Not equal: \n"+
				"expected: %s\n"+
				"actual  : %s", expected, actual)
		}
	}
	return nil
}

func (o *godogFeaturesScenario) cloudwatchLogIsGenerated() error {
	logGroupName := fmt.Sprintf("/aws/lambda/%s", o.stepValues["functionName"])
	client := aws.NewCloudWatchLogsClient(o.testing, o.stepValues["awsRegion"])
	output, _ := client.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &logGroupName,
	})
	if len(output.LogGroups) < 1 {
		return fmt.Errorf("Expected at least one log group. Found %d log groups", len(output.LogGroups))
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
}

func (o *godogFeaturesScenario) destroyTerraform(sc *godog.Scenario, err error) {
	terraform.Destroy(o.testing, o.terraformOptions)
}
