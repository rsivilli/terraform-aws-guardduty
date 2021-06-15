output "aws_cloudwatch_event_rule_name" {
  description = "Name of the CloudWatch rule that watches for AWS GuardDuty findings."
  value       = module.guardduty.aws_cloudwatch_event_rule_name
}
output "aws_guardduty_detector_id" {
  description = "The ID of the GuardDuty detector."
  value       = module.guardduty.aws_guardduty_detector_id
}
