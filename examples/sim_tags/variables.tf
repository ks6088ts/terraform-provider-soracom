variable "soracom_profile" {
  description = "SORACOM Profile"
  type        = string
  default     = "default"
}

variable "sim_id" {
  description = "Target SIM ID"
  type        = string
}

variable "sandbox" {
  description = "Flag for Sandbox"
  type        = bool
  default     = true
}

variable "coverage_type" {
  description = "Coverage type"
  type        = string
  default     = "g"

  validation {
    condition     = contains(["g", "jp"], var.coverage_type)
    error_message = "Invalid variable for coverage_type, expected g or jp."
  }
}

variable "name" {
  description = "Target SIM's name"
  type        = string
}
