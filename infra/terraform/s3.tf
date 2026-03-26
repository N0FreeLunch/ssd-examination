resource "aws_s3_bucket" "backup_bucket" {
  bucket_prefix = "${var.project_name}-backup-"
  force_destroy = true
}

resource "aws_s3_bucket_versioning" "backup_bucket_ver" {
  bucket = aws_s3_bucket.backup_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# Deployment Artifacts & Metadata Bucket
# ---------------------------------------------------------------------------------------------------------------------
resource "aws_s3_bucket" "deploy_bucket" {
  bucket        = "sdd-exam-deploy-artifacts-dev-992481306230"
  force_destroy = true
}

resource "aws_s3_bucket_versioning" "deploy_bucket_ver" {
  bucket = aws_s3_bucket.deploy_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

# CI/CD Metadata Registry (GitHub Actions consumption)
resource "aws_s3_object" "deploy_config" {
  bucket = aws_s3_bucket.deploy_bucket.id
  key    = "deploy-config.json"
  content = jsonencode({
    LightsailStaticIP          = aws_lightsail_static_ip.static_ip.ip_address
    S3DeployBucket             = aws_s3_bucket.deploy_bucket.bucket
    CodeDeployApplication      = "examination-app"
    CodeDeployDeploymentGroup  = "examination-dev-group"
  })
}


# IAM User for Application (Litestream)
resource "aws_iam_user" "app_user" {
  name = "${var.project_name}-app-user"
}

resource "aws_iam_access_key" "app_user_key" {
  user = aws_iam_user.app_user.name
}

resource "aws_iam_user_policy" "app_user_policy" {
  name = "${var.project_name}-app-policy"
  user = aws_iam_user.app_user.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:ListBucket",
          "s3:DeleteObject"
        ]
        Effect = "Allow"
        Resource = [
          aws_s3_bucket.backup_bucket.arn,
          "${aws_s3_bucket.backup_bucket.arn}/*",
          aws_s3_bucket.deploy_bucket.arn,
          "${aws_s3_bucket.deploy_bucket.arn}/*"
        ]
      },
      {
        Action = [
          "ssm:GetParametersByPath",
          "ssm:GetParameter",
          "ssm:GetParameters"
        ]
        Effect   = "Allow"
        Resource = "arn:aws:ssm:*:*:parameter/sdd-exam/*"
      },
      {
        Action = [
          "codedeploy:*"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "kms:Decrypt"
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
}
