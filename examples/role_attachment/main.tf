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

resource "soracom_user" "user" {
  operator_id = var.operator_id
  user_name   = var.user_name
  description = "An example user managed by terraform-provider-soracom"
}

resource "soracom_role" "role_default" {
  operator_id = var.operator_id
  role_id     = "default_role"
  description = "A default role managed by terraform-provider-soracom"
  permission  = "{ \"statements\": [ { \"api\": \"*\", \"effect\": \"allow\" } ] }"
}

resource "soracom_role" "role_admin" {
  operator_id = var.operator_id
  role_id     = "default_admin"
  description = "An admin role managed by terraform-provider-soracom"
  permission  = "{ \"statements\": [ { \"api\": \"*\", \"effect\": \"allow\" } ] }"
}

locals {
  role_ids = [
    soracom_role.role_default.role_id,
    soracom_role.role_admin.role_id,
  ]
}
resource "soracom_role_attachment" "role_attachment" {
  for_each    = toset(local.role_ids)
  operator_id = var.operator_id
  role_id     = each.value
  user_name   = soracom_user.user.user_name
}
