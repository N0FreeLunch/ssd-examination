output "static_ip" {
  value = aws_lightsail_static_ip.static_ip.ip_address
}

output "backup_bucket_name" {
  value = aws_s3_bucket.backup_bucket.bucket
}

output "litestream_access_key_id" {
  value = aws_iam_access_key.app_user_key.id
}

output "litestream_secret_access_key" {
  value     = aws_iam_access_key.app_user_key.secret
  sensitive = true
}
