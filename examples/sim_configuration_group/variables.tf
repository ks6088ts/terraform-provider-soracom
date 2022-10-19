variable "soracom_profile" {
  description = "SORACOM Profile"
  type        = string
  default     = "default"
}

variable "soracom_sim_id_list" {
  default = [
    "1234567890123456789",
    "2345678901234567890",
  ]
}
