output "alb_dns_name" {
  description = "The domain name of the Application Load Balancer."
  value       = aws_lb.this.dns_name
}

output "alb_url" {
  description = "The fully qualified URL to access the scalable web service."
  value       = "http://${aws_lb.this.dns_name}/"
}

output "asg_name" {
  description = "The name of the Auto Scaling Group."
  value       = aws_autoscaling_group.this.name
}
