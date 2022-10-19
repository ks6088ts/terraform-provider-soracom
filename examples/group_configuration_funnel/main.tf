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

resource "soracom_credentials" "aws_credentials" {
  credentials_id = var.credentials_id
  description    = "An example AWS credential via terraform-provider-soracom"
  type           = "aws-credentials"
  credentials = {
    accessKeyId     = var.aws_access_key_id
    secretAccessKey = var.aws_secret_access_key
  }
}

resource "soracom_group_configuration_funnel" "group_configuration_funnel" {
  group_id       = soracom_group.group.id
  enabled        = true
  credentials_id = soracom_credentials.aws_credentials.credentials_id
  content_type   = "json"
  add_sim_id     = true

  # # Amazon Kinesis Streams
  # destination {
  #   provider     = "aws"
  #   service      = "kinesis"
  #   resource_url = "https://kinesis.<region>.amazonaws.com/<delivery stream name>"
  #   # fixme: `randomizePartitionKey` is disabled by default and currently not supported yet
  # }
  # # Amazon Kinesis Firehose
  # destination {
  #   provider     = "aws"
  #   service      = "firehose"
  #   resource_url = "https://firehose.<region>.amazonaws.com/<delivery stream name>"
  # }
  # AWS IoT
  destination {
    provider     = "aws"
    service      = "aws-iot"
    resource_url = "https://<id>.iot.<region>.amazonaws.com/<topic>"
  }
  # # Microsoft Azure Event Hubs
  # destination {
  #   provider     = "azure"
  #   service      = "eventhubs"
  #   resource_url = "https://<namespace>.servicebus.windows.net/<event hubs name>/messages"
  # }
  # # Google Cloud Pub/Sub
  # destination {
  #   provider     = "google"
  #   service      = "pubsub"
  #   resource_url = "<topic>"
  # }
}
