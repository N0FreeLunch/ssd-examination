#!/bin/bash
set -e

echo "⏳ Fetching DEPLOY_HOST from SSM..."
HOST=$(aws ssm get-parameter --name "/sdd-exam/dev/DEPLOY_HOST" --with-decryption --query "Parameter.Value" --output text)

echo "📡 Checking SSH connection to $HOST..."
for i in {1..30}; do
    if ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no -i sdd-exam-key.pem "$HOST" "echo 'SSH is ready!'"; then
        echo "✅ SSH connection established."
        break
    fi
    echo "  Wait $i/30... Retrying in 5 seconds."
    sleep 5
done

echo "⚙️ Setting up 1GB Swap memory and verifying Docker installation..."
ssh -o StrictHostKeyChecking=no -i sdd-exam-key.pem "$HOST" << 'EOF'
set -e
if [ ! -f /swapfile ]; then
    echo "Creating 1GB swap file..."
    sudo fallocate -l 1G /swapfile
    sudo chmod 600 /swapfile
    sudo mkswap /swapfile
    sudo swapon /swapfile
    echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
    echo "✅ Swap memory configured."
else
    echo "✅ Swap memory already exists."
fi

free -h

echo "Waiting for Docker to be ready..."
while ! command -v docker &> /dev/null; do
    echo "Waiting 5s for docker installation..."
    sleep 5
done

while ! sudo docker ps &> /dev/null; do
    echo "Waiting 5s for docker daemon to start..."
    sleep 5
done

echo "✅ Docker is up and running!"
EOF

echo "🎉 Server setup completed successfully!"
