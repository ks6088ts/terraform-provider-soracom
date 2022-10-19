terraform {
  required_providers {
    soracom = {
      source = "ks6088ts/soracom"
    }
  }
}

provider "soracom" {
  profile = var.soracom_profile
}

# AWS
resource "soracom_credentials" "aws_credentials" {
  credentials_id = "aws-credentials"
  description    = "AWS credentials by terraform-provider-soracom"
  type           = "aws-credentials"
  credentials = {
    accessKeyId     = var.aws_access_key_id
    secretAccessKey = var.aws_secret_access_key
  }
}

# AWS IAM role
resource "soracom_credentials" "aws_iam_role_credentials" {
  credentials_id = "aws-iam-role-credentials"
  description    = "AWS IAM role credentials by terraform-provider-soracom"
  type           = "aws-iam-role-credentials"
  credentials = {
    roleArn    = var.aws_iam_role_arn
    externalId = var.aws_iam_role_external_id
  }
}

# Azure credentials
resource "soracom_credentials" "azure_credentials" {
  credentials_id = "azure-credentials"
  description    = "Azure credentials by terraform-provider-soracom"
  type           = "azure-credentials"
  credentials = {
    key        = var.azure_credentials_key
    policyName = var.azure_credentials_policy_name
  }
}

# Azure IoT Hub credentials
resource "soracom_credentials" "azure_iot_hub_credentials" {
  credentials_id = "azure-iot-credentials"
  description    = "Azure IoT credentials by terraform-provider-soracom"
  type           = "azureIoT-credentials"
  credentials = {
    sharedAccessKey     = var.azure_credentials_iot_shared_access_key
    sharedAccessKeyName = var.azure_credentials_iot_shared_access_key_name
  }
}

# X.509 certificate
resource "soracom_credentials" "x_509_certificate" {
  credentials_id = "x_509_certificate"
  description    = "X.509 certificate by terraform-provider-soracom"
  type           = "x509"
  credentials = {
    ca   = var.x_509_certificate_ca
    cert = var.x_509_certificate_cert
    key  = var.x_509_certificate_key
  }
}
