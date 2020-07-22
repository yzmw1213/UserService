resource "aws_key_pair" "mykeypair" {
  key_name = var.KEY_NAME
  public_key = file(var.PUBLIC_KEY_PATH)
}
