data "aws_secretsmanager_secret" "credentials" {
  name = "go_chat_secrets"
}

data "aws_secretsmanager_secret_version" "credentials_retriever" {
  secret_id = data.aws_secretsmanager_secret.credentials.id
}

locals {
  credentials = jsondecode(data.aws_secretsmanager_secret_version.credentials_retriever.secret_string)
}