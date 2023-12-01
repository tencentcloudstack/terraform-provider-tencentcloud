resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_monitor_grafana_instance" "foo" {
  instance_name         = "test-grafana"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_ids            = [tencentcloud_subnet.subnet.id]
  grafana_init_password = "1234567890"
  enable_internet = false

  tags = {
    "createdBy" = "test"
  }
}

resource "tencentcloud_monitor_grafana_integration" "grafanaIntegration" {
  instance_id = tencentcloud_monitor_grafana_instance.foo.id
  kind        = "tencentcloud-monitor-app"
  content     = "{\"kind\":\"tencentcloud-monitor-app\",\"spec\":{\"dataSourceSpec\":{\"authProvider\":{\"__anyOf\":\"使用密钥\",\"useRole\":true,\"secretId\":\"arunma@tencent.com\",\"secretKey\":\"12345678\"},\"name\":\"uint-test\"},\"grafanaSpec\":{\"organizationIds\":[]}}}"
}

resource "tencentcloud_monitor_alarm_notice" "foo" {
  name                  = "tf_alarm_notice"
  notice_type           = "ALL"
  notice_language       = "zh-CN"

  user_notices    {
      receiver_type              = "USER"
      start_time                 = 0
      end_time                   = 1
      notice_way                 = ["SMS","EMAIL"]
      user_ids                   = [10001]
      group_ids                  = []
      phone_order                = [10001]
      phone_circle_times         = 2
      phone_circle_interval      = 50
      phone_inner_interval       = 60
      need_phone_arrive_notice   = 1
      phone_call_type            = "CIRCLE"
      weekday                    =[1,2,3,4,5,6,7]
  }

  url_notices {
      url    = "https://www.mytest.com/validate"
      end_time =  0
      start_time = 1
      weekday = [1,2,3,4,5,6,7]
  }
}

resource "tencentcloud_monitor_grafana_notification_channel" "grafanaNotificationChannel" {
  instance_id   = tencentcloud_monitor_grafana_instance.foo.id
  channel_name  = "tf-channel"
  org_id        = 1
  receivers     = [tencentcloud_monitor_alarm_notice.foo.amp_consumer_id]
  extra_org_ids = ["1"]
}

resource "tencentcloud_monitor_grafana_plugin" "grafanaPlugin" {
  instance_id = tencentcloud_monitor_grafana_instance.foo.id
  plugin_id   = "grafana-piechart-panel"
  version     = "1.6.2"
}

resource "tencentcloud_monitor_grafana_sso_account" "ssoAccount" {
  instance_id = tencentcloud_monitor_grafana_instance.foo.id
  user_id     = var.user_id
  notes       = "desc12222"
  role {
    organization  = "Main Org."
    role          = "Admin"
  }
}