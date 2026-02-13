output "s3_bucket_name" {
  description = "The name of the S3 bucket to be used for the Terraform state."
  value       = aws_s3_bucket.this.id
}

output "dynamodb_table_name" {
  description = "The name of the DynamoDB table to be used for state locking."
  value       = aws_dynamodb_table.this.name
}

output "aws_region" {
  description = "The AWS region where the backend resources are located."
  value       = var.aws_region
}

output "backend_config_example" {
  description = "A helper snippet for configuring the S3 backend in other projects."
  value       = <<EOT
terraform {
  backend "s3" {
    bucket         = "${aws_s3_bucket.this.id}"
    key            = "${var.environment}/infrastructure/my-project/terraform.tfstate"
    region         = "${var.aws_region}"
    dynamodb_table = "${aws_dynamodb_table.this.name}"
    encrypt        = true
  }
}
EOT
}
