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
    name = "terraform-provider-soracom-example"
  }
}

resource "soracom_sim_configuration_group" "sim_configuration_group" {
  for_each     = toset(var.soracom_sim_id_list)
  sim_id       = each.value
  sim_group_id = soracom_group.group.id
}
