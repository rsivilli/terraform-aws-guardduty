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

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

	println(cloudwatchEventRule)
	println(detectorID)

	//start new aws session
	s := session.Must(session.NewSession())
	c := guardduty.New(s, aws.NewConfig().WithRegion((region)))

	_, error := c.GetDetector(&guardduty.GetDetectorInput{DetectorId: &detectorID})

	require.NoError(t, error)

}
