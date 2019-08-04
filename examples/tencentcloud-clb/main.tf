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
  health_check_switch        = true
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

resource "tencentcloud_clb_server_attachment" "server_attachment_http" {
  clb_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_https.id}"
  location_id = "${tencentcloud_clb_listener_rule.rule.id}"
  targets {
    instance_id = "ins-1flbqyp8"
    port        = 23
    weight      = 10
  }
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = "lb-p7olt9e5"
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}


resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_basic.id}"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}
resource "tencentcloud_clb_listener" "listener_target" {
  clb_id        = "${tencentcloud_clb_instance.clb_basic.id}"
  port          = 44
  protocol      = "HTTP"
  listener_name = "listener_basic1"
}
resource "tencentcloud_clb_listener_rule" "rule_target" {
  clb_id              = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_target.id}"
  domain              = "abcd.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}
resource "tencentcloud_clb_rewrite" "rewrite_basic" {
  clb_id                = "${tencentcloud_clb_instance.clb_basic.id}"
  source_listener_id    = "${tencentcloud_clb_listener.listener_basic.id}"
  target_listener_id    = "${tencentcloud_clb_listener.listener_target.id}"
  rewrite_source_loc_id = "${tencentcloud_clb_listener_rule.rule_basic.id}"
  rewrite_target_loc_id = "${tencentcloud_clb_listener_rule.rule_target.id}"
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
data "tencentcloud_clb_server_attachments" "attachments" {
  clb_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_https.id}"
  location_id = "${tencentcloud_clb_server_attachment.server_attachment_http.id}"
}
data "tencentcloud_clb_rewrites" "rewrites" {
  clb_id                = "${tencentcloud_clb_instance.clb_basic.id}"
  source_listener_id    = "${tencentcloud_clb_rewrite.rewrite_basic.source_listener_id}"
  rewrite_source_loc_id = "${tencentcloud_clb_rewrite.rewrite_basic.rewrite_source_loc_id}"
}