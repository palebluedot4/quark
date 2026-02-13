variable "environment" {
  description = "Deployment environment identifier (dev, stg, uat, prod)."
  type        = string
  default     = "dev"

  validation {
    condition     = contains(["dev", "stg", "uat", "prod"], var.environment)
    error_message = "Invalid environment. Must be one of: dev, stg, uat, prod."
  }
}

variable "common_tags" {
  description = "Common metadata tags for all resources."
  type        = map(string)
  default = {
    Project   = "quark"
    Owner     = "palebluedot4"
    ManagedBy = "terraform"
  }
}

variable "aws_region" {
  description = "Target AWS region for resource deployment."
  type        = string
  default     = "ap-northeast-1"
}

variable "bucket_name" {
  description = "Name of the S3 bucket for storing Terraform state (must be globally unique)."
  type        = string

  validation {
    condition     = can(regex("^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$", var.bucket_name))
    error_message = "Bucket name must be 3-63 characters, lowercase, alphanumeric, dots, or hyphens."
  }

  validation {
    condition     = !contains(split("", var.bucket_name), "..")
    error_message = "Bucket name cannot contain consecutive periods (..)."
  }

  validation {
    condition     = !can(regex("^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$", var.bucket_name))
    error_message = "Bucket name cannot be formatted as an IP address."
  }
}

variable "dynamodb_table_name" {
  description = "Name of the DynamoDB table for state locking."
  type        = string
  default     = "terraform-state-locks"
}
