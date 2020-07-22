variable "AWS_ACCESS_KEY" {}
variable "AWS_SECRET_KEY" {}
variable "AWS_REGION" {
  default = "ap-northeast-1"
}

variable "PUBLIC_KEY_PATH" {
  default = "keys/TerraformKey.pub"
}

variable "ECS_INSTANCE_TYPE" {
  default = "t2.micro"
}

variable "ECS_AMIS" {
  type = map
  default = {
    us-east-1 = "ami-13be557e"
    ap-northeast-1 = "ami-08d175f1b493f205f"
  }
}

variable "RDS_NAME" {}
variable "RDS_PASSWORD" {}
variable "RDS_USERNAME" {}

variable "aws_resource_prefix" {
  description = "Prefix to be used in the naming of some of the created AWS resources"
}

variable "TASK_DEFINITION_NAME" {
  description = "Name of ECS Cluster Task definition"
}

variable "SERVICE_NAME" {
}

variable "KEY_NAME" {
  description = "Name of EC2 Instance Connection Key"
}

variable "LOAD_BALANCER_NAME" {
  description = "Name of ECS Instance Load Balancer"
}