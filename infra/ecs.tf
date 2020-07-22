resource "aws_ecs_cluster" "gomicro-cluster" {
  name = local.aws_ecs_cluster_name
}

resource "aws_launch_configuration" "aws_go_micro_main_launchconfig" {
  name_prefix = "ecs-launchconfig"
  image_id = var.ECS_AMIS[var.AWS_REGION]
  instance_type = var.ECS_INSTANCE_TYPE
  key_name = aws_key_pair.mykeypair.key_name
  iam_instance_profile = aws_iam_instance_profile.go-micro-main-ecs-ec2-role.id
  security_groups = [aws_security_group.ecs-securitygroup.id]
  user_data = "#!/bin/bash\necho 'ECS_CLUSTER=${local.aws_ecs_cluster_name}'>/etc/ecs/ecs.config\nstart ecs"
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "ecs-gomicro-autoscaling" {
  name = "ecs-gomicro-autoscaling"
  vpc_zone_identifier = [aws_subnet.go-micro-main-public-1.id, aws_subnet.go-micro-main-public-2.id]
  launch_configuration = aws_launch_configuration.aws_go_micro_main_launchconfig.name
  min_size = 1
  max_size = 2
  tag {
    key = "Name"
    value = "ecs-ec2-container"
    propagate_at_launch = true
  }
}
