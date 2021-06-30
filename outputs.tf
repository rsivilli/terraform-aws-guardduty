output "aws_cloudwatch_event_rule_name" {
  description = "Name of the CloudWatch rule that watches for AWS GuardDuty findings."
  value       = aws_cloudwatch_event_rule.guardduty_findings.name
}

output "aws_guardduty_detector_id" {
  description = "The ID of the GuardDuty detector."
  value       = aws_guardduty_detector.main.id
}

output "kms_key_id" {
  description = "The ID of the KMS key that is generated"
  value       = aws_kms_key.guardduty.key_id
}
