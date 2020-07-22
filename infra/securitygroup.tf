resource "aws_security_group" "ecs-securitygroup" {
  vpc_id = aws_vpc.go-micro-main.id
  name = "${var.SERVICE_NAME}_allow_ssh_sg"
  description = "security group that allows ssh and egress traffic"
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port       = 8888
    to_port         = 8888
    protocol        = "tcp"
    security_groups = [aws_security_group.myapp-elb-securitygroup.id]
  }

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "allow-ssh"
  }
}

resource "aws_security_group" "myapp-elb-securitygroup" {
  vpc_id      = aws_vpc.go-micro-main.id
  name        = "${var.SERVICE_NAME}_elb_sg"
  description = "security group for ecs"
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = {
    Name = "myapp-elb"
  }
}

resource "aws_security_group" "allow_mysql" {
  vpc_id = aws_vpc.go-micro-main.id
  name = "allow_mysql_sg"
  description = "security group that allows mysql connection"

  ingress {
    from_port = 3306
    to_port = 3306
    protocol = "tcp"
    security_groups = [aws_security_group.ecs-securitygroup.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    self        = true
  }

  tags = {
    Name = "allow-mysql"
  }
}
