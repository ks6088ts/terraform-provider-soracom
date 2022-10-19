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

resource "soracom_role" "role" {
  operator_id = var.operator_id
  role_id     = var.role_id
  description = "An example role managed by terraform-provider-soracom"
  permission  = "{ \"statements\": [ { \"api\": \"*\", \"effect\": \"allow\" } ] }"
}
