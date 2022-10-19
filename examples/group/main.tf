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
    name = var.name
    team = "soracom"
  }
}
