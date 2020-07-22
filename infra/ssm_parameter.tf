resource "aws_ssm_parameter" "DB_NAME" {
  name  = "DB_NAME"
  type  = "String"
  value = aws_db_instance.rds_mysql_instance.name
}

resource "aws_ssm_parameter" "DB_PASSWORD" {
  name  = "DB_PASSWORD"
  type  = "String"
  value = aws_db_instance.rds_mysql_instance.password
}

resource "aws_ssm_parameter" "DB_USER" {
  name  = "DB_USER"
  type  = "String"
  value = aws_db_instance.rds_mysql_instance.username
}

resource "aws_ssm_parameter" "DB_ADRESS" {
  name  = "DB_ADRESS"
  type  = "String"
  value = aws_db_instance.rds_mysql_instance.endpoint
}
