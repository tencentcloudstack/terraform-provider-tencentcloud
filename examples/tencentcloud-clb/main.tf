resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}

resource "tencentcloud_vpc" "foo" {
  name       = "vpc-temp-test"
  cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_clb_instance" "my_clb" {
  network_type              = "${var.network_type}"
  clb_name                  = "tf-test-clb"
  project_id                = 0
  vpc_id                    = "${tencentcloud_vpc.foo.id}"
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = "${tencentcloud_vpc.foo.id}"


  security_groups = ["${tencentcloud_security_group.foo.id}"]
}

resource "tencentcloud_clb_listener" "my_listener" {
  clb_id                     = "${tencentcloud_clb_instance.my_clb.id}"
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = 1
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
}

resource "tencentcloud_clb_listener" "listener_https" {
  clb_id               = "${tencentcloud_clb_instance.my_clb.id}"
  listener_name        = "listener_https"
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "VfqcL1ME"

}
resource "tencentcloud_clb_listener_rule" "rule" {
  clb_id              = "${tencentcloud_clb_instance.my_clb.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_https.id}"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}
data "tencentcloud_clb_instances" "clbs" {
  clb_id = "${tencentcloud_clb_instance.my_clb.id}"
}

data "tencentcloud_clb_listeners" "listeners" {
  clb_id      = "${tencentcloud_clb_instance.my_clb.id}"
  listener_id = "${tencentcloud_clb_listener.my_listener.id}"
}
data "tencentcloud_clb_listener_rules" "rules" {
  listener_id = "${tencentcloud_clb_listener.listener_https.id}"
  domain      = "${tencentcloud_clb_listener_rule.rule.domain}"
  url         = "${tencentcloud_clb_listener_rule.rule.url}"
}