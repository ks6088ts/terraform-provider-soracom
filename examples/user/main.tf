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

resource "soracom_user" "user" {
  operator_id = var.operator_id
  user_name   = var.user_name
  description = "An example user managed by terraform-provider-soracom"
}
