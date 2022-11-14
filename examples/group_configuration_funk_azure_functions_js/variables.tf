# Azure
variable "azure_prefix" {
  description = "Prefix"
  type        = string
  default     = "azsrcm"
}

variable "azure_resource_location" {
  description = "Location"
  type        = string
  default     = "japaneast"
}

# SORACOM
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
  default     = "terraform-provider-soracom-azfunc"
}
