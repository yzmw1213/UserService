# app

data "template_file" "myapp-task-definition-template" {
  template = file("templates/app.json.tpl")
  vars = {
    CONTAINER_NAME = var.SERVICE_NAME
    DB_ADRESS = aws_ssm_parameter.DB_ADRESS.value
    DB_NAME = aws_ssm_parameter.DB_NAME.value
    DB_PASSWORD = aws_ssm_parameter.DB_PASSWORD.value
    DB_USER = aws_ssm_parameter.DB_USER.value
    REPOSITORY_URL = replace(aws_ecr_repository.myRepository.repository_url, "https://", "")
  }
}

resource "aws_ecs_task_definition" "myapp-task-definition" {
  family                = var.TASK_DEFINITION_NAME
  container_definitions = data.template_file.myapp-task-definition-template.rendered
}

resource "aws_elb" "myapp-elb" {
  name = var.LOAD_BALANCER_NAME

  listener {
    instance_port     = 8888
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  health_check {
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 30
    target              = "HTTP:8888/"
    interval            = 60
  }

  cross_zone_load_balancing   = true
  idle_timeout                = 400
  connection_draining         = true
  connection_draining_timeout = 400

  subnets         = [aws_subnet.go-micro-main-public-1.id, aws_subnet.go-micro-main-public-2.id]
  security_groups = [aws_security_group.myapp-elb-securitygroup.id]

  tags = {
    Name = var.LOAD_BALANCER_NAME
  }
}

resource "aws_ecs_service" "myapp-service" {
  name            = var.SERVICE_NAME
  cluster         = aws_ecs_cluster.gomicro-cluster.id
  task_definition = aws_ecs_task_definition.myapp-task-definition.arn
  desired_count   = 2
  # deployment_maximum_percent = 200
  # deployment_minimum_healthy_percent = 50
  iam_role        = aws_iam_role.ecs-service-role.arn
  depends_on      = [aws_iam_policy_attachment.ecs-service-attach1]

  load_balancer {
    elb_name       = aws_elb.myapp-elb.name
    container_name = var.SERVICE_NAME
    container_port = 8888
  }
  lifecycle {
    ignore_changes = [task_definition]
  }
}
