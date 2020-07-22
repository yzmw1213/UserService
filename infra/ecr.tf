resource "aws_ecr_repository" "myRepository" {
  name = local.aws_ecr_repository_name
}
