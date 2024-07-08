#!/bin/bash

ARCH=$(uname -m)

# Install HashiCorp tools (Terraform and Vault)
sudo apt-get update
sudo apt-get install -y gnupg software-properties-common ca-certificates

apt-get install -y wget unzip jq

wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --yes --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list

sudo apt update
sudo apt-get install -y terraform=1.5.3-1 vault=1.14.0-1

# Install AWS CLI
curl --silent https://awscli.amazonaws.com/awscli-exe-linux-${ARCH}.zip \
    --output awscliv2.zip

unzip awscliv2.zip
sudo ./aws/install --update

rm -rf awscliv2.zip aws

EXTENSION_ARCH=$(uname -m | sed s/x86_64/amd64/)
rm -rf extensions/
# Install Vault Lambda Extension
curl --silent https://releases.hashicorp.com/vault-lambda-extension/0.10.1/vault-lambda-extension_0.10.1_linux_${EXTENSION_ARCH}.zip \
    --output vault-lambda-extension.zip

unzip vault-lambda-extension.zip -d ./
rm -rf vault-lambda-extension.zip

# Build the Docker image
docker build --build-arg="ARCH=${ARCH}" -t gamify:latest .
