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

resource "soracom_sim_tags" "tags" {
  sim_id = var.sim_id
  tags = {
    name             = var.name
    firmware_version = "0.0.0"
    serial_number    = "XXXXXX"
    device_type      = "sensor"
  }
}
