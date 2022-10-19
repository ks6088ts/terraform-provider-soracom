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

resource "soracom_eventhandler_handlers" "handlers" {
  description        = "An example for terraform-provider-soracom"
  name               = "terraform-provider-soracom"
  status             = var.status
  target_operator_id = var.operator_id
  rule_config {
    type = "SubscriberMonthlyTrafficRule"
    properties {
      inactive_timeout_date_const   = "BEGINNING_OF_NEXT_MONTH"
      limit_total_traffic_mega_byte = "1024"
      run_once_among_target         = true
    }
  }
  action_config_list {
    type = "ChangeSpeedClassAction"
    properties {
      speed_class               = "s1.minimum"
      execution_date_time_const = "IMMEDIATELY"
    }
  }

  action_config_list {
    type = "ChangeSpeedClassAction"
    properties {
      speed_class               = "s1.standard"
      execution_date_time_const = "BEGINNING_OF_NEXT_MONTH"
    }
  }
  action_config_list {
    type = "SendMailToOperatorAction"
    properties {
      execution_date_time_const = "IMMEDIATELY"
      execution_offset_minutes  = "1"
      title                     = "速度制限のお知らせ"
      message                   = "対象 IoT SIM: imsi\n\nIoT SIMの月次データ通信量が 1024 MiBに到達したため、通信速度が \"s1.minimum\" に制限されました。"
    }
  }
  action_config_list {
    type = "SendMailToOperatorAction"
    properties {
      execution_date_time_const = "BEGINNING_OF_NEXT_MONTH"
      execution_offset_minutes  = "1"
      title                     = "速度制限解除のお知らせ"
      message                   = "対象 IoT SIM: imsi\n\n速度制限期間が終了したため、通信速度が \"si.standard\" に設定されました。"
    }
  }
  action_config_list {
    type = "ExecuteWebRequestAction"
    properties {
      content_type              = "application/json"
      execution_date_time_const = "IMMEDIATELY"
      execution_offset_minutes  = "1"
      # fixme: not supported for now since currently not working yet. When the headers element is specified, 400 Bad Request is returned.
      # headers = {
      #   hello = 1
      #   hello2 = 2
      # }
      http_method = "GET"
      url         = "https://yoursite.com"
    }
  }
}
