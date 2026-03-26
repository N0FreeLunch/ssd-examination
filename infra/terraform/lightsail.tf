resource "tls_private_key" "pk" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_lightsail_key_pair" "kp" {
  name       = "${var.project_name}-key"
  public_key = tls_private_key.pk.public_key_openssh
}

resource "local_file" "ssh_key" {
  filename        = "${path.module}/../../${var.project_name}-key.pem"
  content         = tls_private_key.pk.private_key_pem
  file_permission = "0400"
}

resource "aws_lightsail_instance" "server" {
  name              = "${var.project_name}-server"
  availability_zone = "${var.aws_region}a"
  blueprint_id      = "ubuntu_22_04"
  bundle_id         = "nano_3_0"
  key_pair_name     = aws_lightsail_key_pair.kp.name

  user_data = templatefile("${path.module}/user_data.sh", {
    project_name          = var.project_name
    aws_region            = var.aws_region
    aws_access_key_id     = aws_iam_access_key.app_user_key.id
    aws_secret_access_key = aws_iam_access_key.app_user_key.secret
    iam_user_arn          = aws_iam_user.app_user.arn
  })
}

resource "aws_lightsail_static_ip" "static_ip" {
  name = "${var.project_name}-static-ip"
}

resource "aws_lightsail_static_ip_attachment" "static_ip_attach" {
  static_ip_name = aws_lightsail_static_ip.static_ip.name
  instance_name  = aws_lightsail_instance.server.name
}

resource "aws_lightsail_instance_public_ports" "server_ports" {
  instance_name = aws_lightsail_instance.server.name

  port_info {
    protocol  = "tcp"
    from_port = 22
    to_port   = 22
  }

  port_info {
    protocol  = "tcp"
    from_port = 80
    to_port   = 80
  }

  port_info {
    protocol  = "tcp"
    from_port = 443
    to_port   = 443
  }
}
