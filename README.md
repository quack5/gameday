# gamify-aws-hcp-participant
This is a template repository designed to be used by participants of an AWS HCP Gamify challenge.

## Getting Started 

Click the "Use this template" button to create a new team repository under the appropriate org. Once this is done, launch a GitHub Codespace using the <> Code button to start the challenge with the new repository.

```shell
./setup/tool-setup.sh
```

This installs Terraform, Vault and the AWS CLI. It also builds the needed Docker container needed for the challenge. Verify you have a Docker image tagged with `gamify:latest` using:

```
docker image list
```

## AWS Account Access

Refer to your facilitator's AWS Workshop Studio event link.

## HCP Vault Access

Setup the Vault CLI with your provided HCP Vault endpoint and namespace:
```
export VAULT_ADDR=<hcp-vault-endpoint-here>
export VAULT_NAMESPACE=admin/<your-namespace-here>
```

Login to your HCP Vault namespace using a [GitHub personal access token](https://github.com/settings/tokens) (verify it has the `read:org` scope):
```shell
vault login -method github
```

You should now be able to list your kv secrets using this command:
```shell
vault kv get -format=json kv/terraform
```

## Terraform Cloud Access

Next, connect to your codespace with a Terraform Cloud workspace. In order to authenticate with Terraform Cloud, you will need an API team token. Your facilitator has supplied this token in your respective HCP Vault namespace. Please verify that you have access to:

- A Terraform Cloud organization and project (UI)
- An HCP Vault namespace (UI)

Once ready, open the file `terraform.tf` and replace the organization and workspace name values accordingly.

Now that we have the remote backend defined, copy & paste the Terraform team token value and use it in the login prompt.

```shell
vault kv get -format=json kv/terraform | jq -r .data.data.team_token
```

Copy the above token and use it in the following prompt:
```shell
terraform login
```

Assuming you're successfully authenticated, initialize the repository with:

```shell
terraform init
```

## The Architecture

![gameday_participant](https://github.com/acornies/gamify-aws-hcp-participant/assets/2882297/fb6646d7-9042-4031-8474-a9d5a4b580c6)

## The App (Lambda function)

The app in this repository is already written, so it just needs to be deployed as a container image. The Vault Lambda Extension needs these environment variables to run:

| Environment variable      | Value |
| ----------- | ----------- |
| VAULT_ADDR      | https://your-event-hcp-vault-cluster.z1.hashicorp.cloud:8200       |
| VAULT_NAMESPACE   | admin/your-namespace        |
| VAULT_AUTH_ROLE   | your-aws-auth-role        |
| VAULT_AUTH_PROVIDER   | aws        |
| VAULT_SECRET_PATH_DB   | database/creds/your-db-role        |
| VAULT_SECRET_FILE_DB   | /tmp/vault_secret.json        |

The app needs these environment variables to run:

| Environment variable      | Value |
| ----------- | ----------- |
| VAULT_SECRET_FILE_DB   | /tmp/vault_secret.json        |
| DATABASE_ADDR   | your-rds-address.us-east-2.rds.amazonaws.com:5432/dbname        |

## The Challenge

Build and deploy the lambda function that receives messages from a specified SQS queue. The function is built to use the Vault Lambda Extension to secure the Postgres database connection using a dynamic database credential. 

### Scoring

Once your function starts receiving and processing messages from the queue, your team will receive points for very message processed. Your team will also receive points for the types of AWS resources that are living in your AWS account.

Additional points may be awarded for Terraform style and the use of Terraform Cloud.

The team with the highest amount of points wins the challenge. 🏆

## Suggested Steps

1. 📝 Launch a Github Codespace from this repository, run `./setup/tool-setup.sh`
   1. Hook up your Terraform code/workspace to Terraform Cloud (instructions above)
   2. Add AWS and Vault credentials to your Terraform Cloud workspace variables
2. 🐳 Create an ECR repository for the Docker image in AWS
   1. Push the `gamify:latest` Docker container to the ECR repository
3. 🐘 Create a RDS Postgres database instance
   1. Think about a security group for the RDS instance
4. 🚀 Create a Lambda function with package type "image"
   1. Think about an IAM policy and role needed for the Vault integration
5. 📄 Configure the Lambda with the image url from the ECR repository
   1. Provide the [config](#the-app-lambda-function) needed for the app to run
6. 📬 Map the SQS event source to your Lambda (Queue ARN provided)
7.  🔒 Configure your Vault namespace for the Lambda to fetch a dynamic database credential
    1. The AWS auth method is needed
    2. An AWS auth role is also needed for the Lambda
    3. Think about a suitable Vault policy to assign to the role
    4. A database secrets engine of type Postgres is needed
    5. A database engine role is also needed to vend Postgres accounts
8.  🎉 Test your Lambda function!

Don't forget to commit and push your code the repository!

## Debugging

In order debug your function code, add the `AWSLambdaBasicExecutionRole` managed policy to your Lambda execution role to view AWS Cloud Watch logs.

## References to help

- [The AWS Terraform provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [The Vault Terraform provider](https://registry.terraform.io/providers/hashicorp/vault/latest/docs)
- [Vault AWS auth](https://developer.hashicorp.com/vault/tutorials/cloud-ops/vault-auth-method-aws)
- [Introduction to the Vault AWS Lambda extension](https://developer.hashicorp.com/vault/tutorials/app-integration/intro-vault-aws-lambda-extension)
- [Learn Vault Lambda Extension Github](https://github.com/hashicorp-education/learn-vault-lambda-extension)
