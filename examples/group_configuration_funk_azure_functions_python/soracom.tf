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

resource "soracom_group" "group" {
  tags = {
    name = var.soracom_group_name
  }
}

resource "soracom_credentials" "credentials" {
  credentials_id = var.credentials_id
  description    = "API Token credentials via terraform-provider-soracom"
  type           = "api-token-credentials"
  credentials = {
    token = data.azurerm_function_app_host_keys.example.default_function_key
  }
}

resource "soracom_group_configuration_funk" "group_configuration_funk" {
  group_id       = soracom_group.group.id
  enabled        = true
  credentials_id = soracom_credentials.credentials.credentials_id
  content_type   = "json"

  # Azure Functions
  destination {
    provider     = "azure"
    service      = "function-app"
    resource_url = azurerm_function_app_function.example.invocation_url
  }
}
