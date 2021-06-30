// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/guardduty"
)

func TestTerraformSimpleExample(t *testing.T) {
	t.Parallel()

	//pull default region from environment. If it doesn't exist, terminate current test
	region := os.Getenv("AWS_DEFAULT_REGION")
	require.NotEmpty(t, region, "missing environment variable AWS_DEFAULT_REGION")

	//naming of test
	testName := fmt.Sprintf("terratest-guardduty-simple-%s", strings.ToLower(random.UniqueId()))

	//configuring root Terraform directory
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/simple",
		//variables that terraform uses
		Vars: map[string]interface{}{
			"test_name": testName,
			"tags": map[string]interface{}{
				"Automation": "Terraform",
				"Terratest":  "yes",
				"Test":       "TestTerraformSimpleExample",
			},
		},
		//environment variables that everything that this test runs might use
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": region,
		},
	})

	//enable test to not destroy things at end of test
	if os.Getenv("TT_SKIP_DESTROY") != "1" {
		defer terraform.Destroy(t, terraformOptions)
	}

	//Runs terraform init and terraform apply with the options specified
	terraform.InitAndApply(t, terraformOptions)

	//Get outputs of
	cloudwatchEventRule := terraform.Output(t, terraformOptions, "aws_cloudwatch_event_rule_name")
	detectorID := terraform.Output(t, terraformOptions, "aws_guardduty_detector_id")
	//kmsKeyID := terraform.Output(t, terraformOptions, "kms_key_id")

	//start new aws session
	s := session.Must(session.NewSession())
	gdClient := guardduty.New(s, aws.NewConfig().WithRegion((region)))

	//verify that the detector was created and configured correctly
	detectorresponse, detectorerror := gdClient.GetDetector(&guardduty.GetDetectorInput{DetectorId: &detectorID})
	require.NoError(t, detectorerror, "We should be able to access the detector based on the output id")
	require.Equal(t, "ENABLED", *detectorresponse.DataSources.CloudTrail.Status, "Cloudtrail logs should be enabled")
	require.Equal(t, "ENABLED", *detectorresponse.DataSources.DNSLogs.Status, "DNS logs should be enabled")
	require.Equal(t, "ENABLED", *detectorresponse.DataSources.FlowLogs.Status, "Flowlogs should be enabled")
	require.Equal(t, "ENABLED", *detectorresponse.DataSources.S3Logs.Status, "S3 Logs should be enabled")
	require.Equal(t, "ENABLED", *detectorresponse.Status, "The detector should be enabled")

	//create sample findings (NOTE this should also create alerts in s3/cloudwatch)
	_, createFindingsError := gdClient.CreateSampleFindings(&guardduty.CreateSampleFindingsInput{DetectorId: &detectorID})

	require.NoError(t, createFindingsError, "We should not receive an error when creating sample findings")

	//verify that findings can be seen in guardduty
	findingsResponse, findingsResponseError := gdClient.ListFindings(&guardduty.ListFindingsInput{DetectorId: &detectorID})

	require.NoError(t, findingsResponseError)
	require.GreaterOrEqual(t, len(findingsResponse.FindingIds), 5, "Findings should be more than 5")
	require.NotNil(t, findingsResponse.NextToken, "There should be a token indicating response was larger than one payload")

	fmt.Print(*findingsResponse)
	//var cwrule = cloudwatch.MetricDataQuery{MetricStat: &cloudwatch.MetricStat{&cloudwatch.Metric{}}}
	cwclient := cloudwatch.New(s, aws.NewConfig().WithRegion(region))
	cwclient.GetMetricData(&cloudwatch.GetMetricDataInput{EndTime: aws.Time(time.Now()), MetricDataQueries: {}})

	// kmsClient := kms.New(s, aws.NewConfig().WithRegion((region)))
	// kmsKeyPolicyResult, kmserror := kmsClient.GetKeyPolicy(&kms.GetKeyPolicyInput{KeyId: &kmsKeyID})
	// require.NoError(t, kmserror, "We should be able to access the key policy on the output id")
	// fmt.Println((*kmsKeyPolicyResult))

	//verify that a policy for guardduty was created and configured correctly

}
