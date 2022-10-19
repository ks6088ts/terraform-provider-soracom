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

resource "soracom_group" "group" {
  tags = {
    name = var.soracom_group_name
  }
}

resource "soracom_group_configuration_air" "group_configuration_air" {
  group_id       = soracom_group.group.id
  use_custom_dns = true
  dns_servers = [
    "8.8.8.8",
    "8.8.4.4",
  ]
  meta_data {
    enabled      = true
    read_only    = true
    allow_origin = "http://some.example.com"
  }
  user_data = "foobar"
}
