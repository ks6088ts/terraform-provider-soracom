terraform {
  required_providers {
    soracom = {
      source = "ks6088ts.github.io/ks6088ts/soracom"
    }
  }
}

provider "soracom" {
  profile = var.soracom_profile
}

resource "soracom_group" "group" {
  tags = {
    name = var.soracom_group_name
  }
}

resource "soracom_credentials" "aws_credentials" {
  credentials_id = var.credentials_id
  description    = "An example AWS credential via terraform-provider-soracom"
  type           = "aws-credentials"
  credentials = {
    accessKeyId     = var.aws_access_key_id
    secretAccessKey = var.aws_secret_access_key
  }
}

resource "soracom_group_configuration_funk" "group_configuration_funk" {
  group_id       = soracom_group.group.id
  enabled        = true
  credentials_id = soracom_credentials.aws_credentials.credentials_id
  content_type   = "json"

  # AWS Lambda
  destination {
    provider     = "aws"
    service      = "lambda"
    resource_url = "arn:aws:lambda:region:xxxxxxxxxxxx:function:funcname"
  }

  # Azure Functions
  # destination {
  #   provider     = "azure"
  #   service      = "function-app"
  #   resource_url = "https://appname.azurewebsites.net/api/funcname"
  # }

  # Google Cloud Functions
  # destination {
  #   provider     = "google"
  #   service      = "cloud-functions"
  #   resource_url = "https://region-project.cloudfunctions.net/funcname"
  # }
}
