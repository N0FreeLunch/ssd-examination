terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "sdd-exam-tfstate-ap-northeast-2-992481306230"
    key            = "infra/terraform.tfstate"
    region         = "ap-northeast-2"
    dynamodb_table = "sdd-exam-tfstate-lock"
    encrypt        = true
  }
}


provider "aws" {
  region = "ap-northeast-2"
  # profile = "terraform-admin" # Uncomment if using a specific profile
}
