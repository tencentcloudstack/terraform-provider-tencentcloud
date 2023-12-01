data "tencentcloud_antiddos_basic_device_status" "basic_device_status" {
  ip_list = [
    "127.0.0.1"
  ]
  filter_region = 1
}

data "tencentcloud_antiddos_bgp_biz_trend" "bgp_biz_trend" {
  business    = "bgp-multip"
  start_time  = "2023-11-22 09:25:00"
  end_time    = "2023-11-22 10:25:00"
  metric_name = "intraffic"
  instance_id = "bgp-00000ry7"
  flag        = 0
}

data "tencentcloud_antiddos_list_listener" "list_listener" {
}

data "tencentcloud_antiddos_overview_attack_trend" "overview_attack_trend" {
  type       = "ddos"
  dimension  = "attackcount"
  period     = 86400
  start_time = "2023-11-21 10:28:31"
  end_time   = "2023-11-22 10:28:31"
}