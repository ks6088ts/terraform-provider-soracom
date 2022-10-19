variable "soracom_profile" {
  description = "SORACOM Profile"
  type        = string
  default     = "default"
}

variable "soracom_group_name" {
  description = "SIM group name"
  type        = string
  default     = "terraform-provider-soracom"
}

variable "credentials_id" {
  description = "Credentials ID"
  type        = string
  default     = "terraform-provider-soracom-aws"
}

variable "aws_access_key_id" {
  description = "AWS AccessKeyId"
  type        = string
}

variable "aws_secret_access_key" {
  description = "AWS SecretAccessKey"
  type        = string
}
