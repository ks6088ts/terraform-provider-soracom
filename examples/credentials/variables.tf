variable "soracom_profile" {
  description = "SORACOM Profile"
  type        = string
  default     = "default"
}

variable "aws_access_key_id" {
  description = "AWS AccessKeyId"
  type        = string
}

variable "aws_secret_access_key" {
  description = "AWS SecretAccessKey"
  type        = string
}

variable "aws_iam_role_arn" {
  description = "AWS IAM role ARN"
  type        = string
}

variable "aws_iam_role_external_id" {
  description = "AWS IAM role External ID"
  type        = string
}

variable "azure_credentials_key" {
  description = "Azure credentials key"
  type        = string
}

variable "azure_credentials_policy_name" {
  description = "Azure credentials policy name"
  type        = string
}

variable "azure_credentials_iot_shared_access_key" {
  description = "Azure credentials IoT shared access key"
  type        = string
}

variable "azure_credentials_iot_shared_access_key_name" {
  description = "Azure credentials shared access key name"
  type        = string
}

variable "x_509_certificate_ca" {
  description = "Certificate Authority for X.509 certificate"
  type        = string
}

variable "x_509_certificate_cert" {
  description = "A certificate for X.509 certificate"
  type        = string
}

variable "x_509_certificate_key" {
  description = "A secret key for X.509 certificate"
  type        = string
}
