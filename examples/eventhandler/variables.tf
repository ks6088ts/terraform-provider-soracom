variable "soracom_profile" {
  description = "SORACOM Profile"
  type        = string
  default     = "default"
}

variable "status" {
  description = "EventHandler status"
  type        = string
  default     = "inactive"

  validation {
    condition     = contains(["inactive", "active"], var.status)
    error_message = "Invalid variable for status, expected inactive or active."
  }
}

variable "operator_id" {
  description = "Target operator ID"
  type        = string
}
