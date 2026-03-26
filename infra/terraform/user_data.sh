#!/bin/bash
# Install Docker and Docker Compose
apt-get update
apt-get install -y ca-certificates curl gnupg lsb-release git

mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin awscli jq

# Create App Directory
mkdir -p /app
chown ubuntu:ubuntu /app

# Add ubuntu user to docker group
usermod -aG docker ubuntu

# Configure AWS Credentials for Root (Used by CodeDeploy Agent)
mkdir -p /root/.aws
cat <<EOF > /root/.aws/credentials
[default]
aws_access_key_id = ${aws_access_key_id}
aws_secret_access_key = ${aws_secret_access_key}
EOF

cat <<EOF > /root/.aws/config
[default]
region = ${aws_region}
EOF

# Install CodeDeploy Agent
apt-get update
apt-get install -y ruby-full wget
cd /home/ubuntu
wget https://aws-codedeploy-${aws_region}.s3.${aws_region}.amazonaws.com/latest/install
chmod +x ./install
./install auto

service codedeploy-agent start

# Register instance with CodeDeploy (On-premises mode for Lightsail)
# This allows CodeDeploy to identify this instance by name
aws deploy register-on-premises-instance \
  --instance-name ${project_name}-server \
  --iam-user-arn ${iam_user_arn} \
  --region ${aws_region}

aws deploy add-tags-to-on-premises-instances \
  --instance-names ${project_name}-server \
  --tags Key=Name,Value=${project_name}-server \
  --region ${aws_region}

echo "Infrastructure setup complete."
