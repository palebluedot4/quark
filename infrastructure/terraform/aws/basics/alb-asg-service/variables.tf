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

variable "instance_type" {
  description = "EC2 instance type (must be a Graviton/ARM64 family)."
  type        = string
  default     = "t4g.micro"

  validation {
    condition     = can(regex("^[a-z]+[0-9]+g[a-z]*\\.", var.instance_type))
    error_message = "The instance_type must belong to a Graviton (ARM64) family (e.g., t4g.micro, c8g.large)."
  }
}

variable "server_port" {
  description = "The TCP port the HTTP server will listen on."
  type        = number
  default     = 8080

  validation {
    condition     = var.server_port >= 1 && var.server_port <= 65535
    error_message = "The server_port must be a valid port number (1-65535)."
  }
}

variable "min_size" {
  description = "The minimum number of EC2 instances the ASG can scale down to."
  type        = number
  default     = 2

  validation {
    condition     = var.min_size >= 0
    error_message = "Minimum ASG size must be a non-negative integer."
  }
}

variable "max_size" {
  description = "The maximum number of EC2 instances the ASG can scale up to."
  type        = number
  default     = 4

  validation {
    condition     = var.max_size > 0
    error_message = "Maximum ASG size must be greater than 0."
  }
}

variable "desired_capacity" {
  description = "The initial number of EC2 instances to launch in the ASG."
  type        = number
  default     = 2

  validation {
    condition     = var.desired_capacity >= 0
    error_message = "Desired capacity must be a non-negative integer."
  }
}
