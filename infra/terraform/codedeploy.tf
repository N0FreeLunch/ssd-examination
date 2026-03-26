# CodeDeploy Application
resource "aws_codedeploy_app" "app" {
  name = "examination-app"
}

# CodeDeploy Service Role
resource "aws_iam_role" "codedeploy_role" {
  name = "${var.project_name}-codedeploy-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "codedeploy.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "codedeploy_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSCodeDeployRole"
  role       = aws_iam_role.codedeploy_role.name
}

# CodeDeploy Deployment Group
resource "aws_codedeploy_deployment_group" "dev_group" {
  app_name              = aws_codedeploy_app.app.name
  deployment_group_name = "examination-dev-group"
  service_role_arn      = aws_iam_role.codedeploy_role.arn

  deployment_config_name = "CodeDeployDefault.OneAtATime"

  # For Lightsail (on-premises tag based)
  on_premises_instance_tag_filter {
    key   = "Name"
    type  = "KEY_AND_VALUE"
    value = aws_lightsail_instance.server.name
  }

  # NOTE: Lightsail instances are treated as on-premises instances by CodeDeploy 
  # because they aren't managed via standard EC2 integration.
  # We will tag the instance in the agent configuration on the server.
}
