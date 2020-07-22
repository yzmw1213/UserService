resource "aws_db_subnet_group" "mysql-subnet-group" {
  name = "mysql-subnet"
  description = "RDS subnet group"
  subnet_ids = [aws_subnet.go-micro-main-private-1.id, aws_subnet.go-micro-main-private-2.id]
}

resource "aws_db_instance" "rds_mysql_instance" {
  allocated_storage    = 20
  availability_zone = aws_subnet.go-micro-main-private-1.availability_zone
  db_subnet_group_name = aws_db_subnet_group.mysql-subnet-group.name
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  name                 = var.RDS_NAME
  username             = var.RDS_USERNAME
  password             = var.RDS_PASSWORD
  parameter_group_name = "default.mysql5.7"
  vpc_security_group_ids  = [aws_security_group.allow_mysql.id]
  multi_az = "false"
  final_snapshot_identifier = "rds-mysql-instance-cluster-backup"
  skip_final_snapshot = "true"
}
