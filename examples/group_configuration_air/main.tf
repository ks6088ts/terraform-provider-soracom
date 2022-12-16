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
  binary_parser_enabled = true
  binary_parser_format = "frameType:0:uint:1:big-endian:7 batLow:0:bool:6 boot:0:bool:5 coSensor:0:bool:4 temp:0:int:12:big-endian:3 hygro:2:uint:8 co:3:uint:8"
}
