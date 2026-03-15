data "aws_vpc" "default" {
  default = true
}

data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-arm64"]
  }

  filter {
    name   = "architecture"
    values = ["arm64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  lifecycle {
    postcondition {
      condition     = self.architecture == "arm64"
      error_message = "The selected AMI must be ARM64 architecture."
    }

    postcondition {
      condition     = self.root_device_type == "ebs"
      error_message = "The selected AMI must be EBS-backed."
    }
  }
}

resource "aws_security_group" "this" {
  name_prefix = "simple-web-server-sg-"
  description = "Allow HTTP inbound traffic and all outbound traffic."
  vpc_id      = data.aws_vpc.default.id

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "simple-web-server-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "http" {
  security_group_id = aws_security_group.this.id
  description       = "Allow HTTP inbound traffic."

  cidr_ipv4   = "0.0.0.0/0"
  ip_protocol = "tcp"
  from_port   = var.server_port
  to_port     = var.server_port
}

resource "aws_vpc_security_group_egress_rule" "all" {
  security_group_id = aws_security_group.this.id
  description       = "Allow all outbound traffic."

  cidr_ipv4   = "0.0.0.0/0"
  ip_protocol = "-1"
}

resource "aws_instance" "this" {
  ami           = data.aws_ami.amazon_linux.id
  instance_type = var.instance_type

  vpc_security_group_ids      = [aws_security_group.this.id]
  associate_public_ip_address = true

  root_block_device {
    volume_size           = 8
    volume_type           = "gp3"
    encrypted             = true
    delete_on_termination = true
  }

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 1
  }

  user_data_replace_on_change = true

  user_data = <<-EOF
              #!/usr/bin/env bash
              set -euo pipefail
              WORKDIR="/var/www/quark"
              mkdir -p "$WORKDIR"
              cd "$WORKDIR"
              echo "<h1>Hello, World!</h1><p>Server is running on port ${var.server_port}</p>" > index.html
              nohup python3 -m http.server ${var.server_port} --bind 0.0.0.0 >/dev/null 2>&1 &
              EOF

  tags = {
    Name = "simple-web-server-instance"
  }
}
