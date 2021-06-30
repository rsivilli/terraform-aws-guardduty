<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Usage

Creates an AWS GuardDuty Detector, KMS Key for encrypting exports to S3, and CloudWatch rule to watch for findings.

```hcl
module "guardduty" {
  source = "dod-iac/guardduty/aws"

  tags = {
    Application = var.application
    Environment = var.environment
    Automation  = "Terraform"
  }
}
```

You can customize the finding publishing frequency.

```hcl
module "guardduty" {
  source = "dod-iac/guardduty/aws"

  enable = true
  finding_publishing_frequency = "SIX_HOURS"
  tags = {
    Application = var.application
    Environment = var.environment
    Automation  = "Terraform"
  }
}
```

You can exports GuardDuty findings to a S3 bucket using the s3\_bucket\_name variable.

```hcl
module "guardduty" {
  source = "dod-iac/guardduty/aws"

  enable = true
  s3_bucket_name = module.logs.aws_logs_bucket
  tags = {
    Application = var.application
    Environment = var.environment
    Automation  = "Terraform"
  }
}
```

## Terraform Version

Terraform 0.12. Pin module version to ~> 1.0.0 . Submit pull-requests to master branch.

Terraform 0.11 is not supported.

## License

This project constitutes a work of the United States Government and is not subject to domestic copyright protection under 17 USC § 105.  However, because the project utilizes code licensed from contributors and other third parties, it therefore is licensed under the MIT License.  See LICENSE file for more information.

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.13 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 3.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | ~> 3.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_cloudwatch_event_rule.guardduty_findings](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_event_rule) | resource |
| [aws_guardduty_detector.main](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/guardduty_detector) | resource |
| [aws_guardduty_publishing_destination.main](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/guardduty_publishing_destination) | resource |
| [aws_kms_alias.guardduty](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_alias) | resource |
| [aws_kms_key.guardduty](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_key) | resource |
| [aws_s3_bucket_object.guardduty](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_object) | resource |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity) | data source |
| [aws_iam_policy_document.key_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_partition.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/partition) | data source |
| [aws_s3_bucket.main](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/s3_bucket) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_enable"></a> [enable](#input\_enable) | Enable monitoring and feedback reporting.  Setting to false is equivalent to "suspending" GuardDuty. | `bool` | `true` | no |
| <a name="input_finding_publishing_frequency"></a> [finding\_publishing\_frequency](#input\_finding\_publishing\_frequency) | Specifies the frequency of notifications sent for subsequent finding occurrences.  If the detector is a GuardDuty member account, the value is determined by the GuardDuty master account and cannot be modified, otherwise defaults to SIX\_HOURS.  For standalone and GuardDuty master accounts, it must be configured in Terraform to enable drift detection.  Valid values for standalone and master accounts: FIFTEEN\_MINUTES, ONE\_HOUR, SIX\_HOURS. | `string` | `"FIFTEEN_MINUTES"` | no |
| <a name="input_kms_alias_name"></a> [kms\_alias\_name](#input\_kms\_alias\_name) | The display name of the alias of the KMS key used to encrypt exports to S3. The name must start with the word "alias" followed by a forward slash (alias/). | `string` | `"alias/guardduty"` | no |
| <a name="input_kms_key_tags"></a> [kms\_key\_tags](#input\_kms\_key\_tags) | Tags to apply to the AWS KMS Key used to encrypt exports to S3. | `map(string)` | `{}` | no |
| <a name="input_s3_bucket_name"></a> [s3\_bucket\_name](#input\_s3\_bucket\_name) | The name of the S3 bucket that receives findings from GuardDuty.  If blank, then GuardDuty does not export findings to S3. | `string` | `""` | no |
| <a name="input_s3_bucket_prefix"></a> [s3\_bucket\_prefix](#input\_s3\_bucket\_prefix) | The prefix for where findings from GuardDuty are stored in the S3 bucket.  Should start with "/" if defined.  GuardDuty will build the full destination ARN using this format: <s3\_bucket\_arn><s3\_bucket\_prefix>/AWSLogs/<account\_id>/GuardDuty/<region>. | `string` | `"/guardduty"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_aws_cloudwatch_event_rule_name"></a> [aws\_cloudwatch\_event\_rule\_name](#output\_aws\_cloudwatch\_event\_rule\_name) | Name of the CloudWatch rule that watches for AWS GuardDuty findings. |
| <a name="output_aws_guardduty_detector_id"></a> [aws\_guardduty\_detector\_id](#output\_aws\_guardduty\_detector\_id) | The ID of the GuardDuty detector. |
| <a name="output_kms_key_id"></a> [kms\_key\_id](#output\_kms\_key\_id) | The ID of the KMS key that is generated |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
