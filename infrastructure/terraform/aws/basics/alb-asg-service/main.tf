data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }

  lifecycle {
    postcondition {
      condition     = length(self.ids) >= 2
      error_message = "Must have at least two subnets in different Availability Zones for high availability."
    }
  }
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

data "aws_launch_template" "current" {
  name = aws_launch_template.this.name

  depends_on = [aws_launch_template.this]
}

data "aws_default_tags" "this" {}

resource "aws_security_group" "lb" {
  name_prefix = "alb-asg-service-lb-sg-"
  description = "Allow HTTP inbound to Load Balancer."
  vpc_id      = data.aws_vpc.default.id

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "alb-asg-service-lb-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "lb_http" {
  security_group_id = aws_security_group.lb.id
  description       = "Allow inbound HTTP traffic from internet."

  cidr_ipv4   = "0.0.0.0/0"
  ip_protocol = "tcp"
  from_port   = 80
  to_port     = 80
}

resource "aws_vpc_security_group_egress_rule" "lb_all" {
  security_group_id = aws_security_group.lb.id
  description       = "Allow all outbound traffic from Load Balancer."

  cidr_ipv4   = "0.0.0.0/0"
  ip_protocol = "-1"
}

resource "aws_security_group" "instance" {
  name_prefix = "alb-asg-service-instance-sg-"
  description = "Allow HTTP inbound from Load Balancer."
  vpc_id      = data.aws_vpc.default.id

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "alb-asg-service-instance-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "instance_from_lb" {
  security_group_id = aws_security_group.instance.id
  description       = "Allow inbound HTTP traffic from Load Balancer only."

  referenced_security_group_id = aws_security_group.lb.id

  ip_protocol = "tcp"
  from_port   = var.server_port
  to_port     = var.server_port
}

resource "aws_vpc_security_group_egress_rule" "instance_all" {
  security_group_id = aws_security_group.instance.id
  description       = "Allow all outbound traffic from instances."

  cidr_ipv4   = "0.0.0.0/0"
  ip_protocol = "-1"
}

resource "aws_lb_target_group" "this" {
  name_prefix = "lbasg-"
  port        = var.server_port
  protocol    = "HTTP"
  vpc_id      = data.aws_vpc.default.id

  health_check {
    path                = "/"
    protocol            = "HTTP"
    matcher             = "200"
    interval            = 15
    timeout             = 3
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }

  tags = {
    Name = "alb-asg-service-tg"
  }
}

resource "aws_lb" "this" {
  name               = "alb-asg-service-lb"
  load_balancer_type = "application"
  subnets            = data.aws_subnets.default.ids
  security_groups    = [aws_security_group.lb.id]

  tags = {
    Name = "alb-asg-service-lb"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.this.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}

resource "aws_launch_template" "this" {
  name_prefix            = "alb-asg-service-lt-"
  image_id               = data.aws_ami.amazon_linux.id
  instance_type          = var.instance_type
  update_default_version = true

  network_interfaces {
    associate_public_ip_address = true
    delete_on_termination       = true
    security_groups             = [aws_security_group.instance.id]
  }

  block_device_mappings {
    device_name = data.aws_ami.amazon_linux.root_device_name

    ebs {
      volume_size           = 8
      volume_type           = "gp3"
      encrypted             = true
      delete_on_termination = true
    }
  }

  metadata_options {
    http_tokens   = "required"
    http_endpoint = "enabled"
  }

  user_data = base64encode(<<-EOF
              #!/usr/bin/env bash
              set -euo pipefail
              TOKEN=$(curl -s -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
              INSTANCE_ID=$(curl -s -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/instance-id)
              WORKDIR="/var/www/quark"
              mkdir -p "$WORKDIR"
              cd "$WORKDIR"
              echo "<h1>Hello, World!</h1><p>Instance ID: $INSTANCE_ID</p>" > index.html
              nohup python3 -m http.server ${var.server_port} --bind 0.0.0.0 >/dev/null 2>&1 &
              EOF
  )

  tags = {
    Name = "alb-asg-service-lt"
  }

  tag_specifications {
    resource_type = "instance"
    tags = merge(data.aws_default_tags.this.tags, {
      Name = "alb-asg-service-instance"
    })
  }

  tag_specifications {
    resource_type = "volume"
    tags = merge(data.aws_default_tags.this.tags, {
      Name = "alb-asg-service-volume"
    })
  }

  tag_specifications {
    resource_type = "network-interface"
    tags = merge(data.aws_default_tags.this.tags, {
      Name = "alb-asg-service-eni"
    })
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "this" {
  name_prefix = "alb-asg-service-"

  vpc_zone_identifier = data.aws_subnets.default.ids
  target_group_arns   = [aws_lb_target_group.this.arn]

  min_size         = var.min_size
  max_size         = var.max_size
  desired_capacity = var.desired_capacity

  health_check_type         = "ELB"
  health_check_grace_period = 120

  launch_template {
    id      = aws_launch_template.this.id
    version = data.aws_launch_template.current.latest_version
  }

  instance_refresh {
    strategy = "Rolling"

    preferences {
      min_healthy_percentage = 50
      instance_warmup        = 120
    }
  }

  lifecycle {
    create_before_destroy = true

    precondition {
      condition     = var.max_size >= var.min_size
      error_message = "The max_size (${var.max_size}) must be greater than or equal to min_size (${var.min_size})."
    }

    precondition {
      condition     = var.desired_capacity >= var.min_size && var.desired_capacity <= var.max_size
      error_message = "The desired_capacity (${var.desired_capacity}) must be between min_size and max_size."
    }
  }

  tag {
    key                 = "Name"
    value               = "alb-asg-service"
    propagate_at_launch = false
  }

  dynamic "tag" {
    for_each = data.aws_default_tags.this.tags
    content {
      key                 = tag.key
      value               = tag.value
      propagate_at_launch = false
    }
  }
}
