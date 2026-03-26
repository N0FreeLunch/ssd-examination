terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-2"
}

# Terraform State 파일용 S3 버킷
resource "aws_s3_bucket" "terraform_state" {
  bucket = "sdd-exam-tfstate-ap-northeast-2-992481306230" # Replace with your account ID or a unique string

  # 실수로 삭제하는 것을 방지
  lifecycle {
    prevent_destroy = true
  }
}

# S3 버킷 버전 관리 활성화 (파일 손상 대비)
resource "aws_s3_bucket_versioning" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id
  versioning_configuration {
    status = "Enabled"
  }
}

# 기본 암호화 설정 (필요 시)
resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# S3 버킷 퍼블릭 액세스 차단
resource "aws_s3_bucket_public_access_block" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Terraform State Locking용 DynamoDB 테이블
resource "aws_dynamodb_table" "terraform_state_lock" {
  name           = "sdd-exam-tfstate-lock"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}
