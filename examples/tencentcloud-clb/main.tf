resource "tencentcloud_security_group" "foo" {
  name = "example"
}

resource "tencentcloud_vpc" "foo" {
  name       = "example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  name              = "example"
  availability_zone = "${var.availability_zone}"
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  cidr_block        = "10.0.0.0/24"
  is_multicast      = false
}

resource "tencentcloud_instance" "foo" {
  instance_name              = "example"
  availability_zone          = "${var.availability_zone}"
  image_id                   = "img-9qabwvbn"
  instance_type              = "S2.SMALL1"
  system_disk_type           = "CLOUD_PREMIUM"
  internet_max_bandwidth_out = 0
  vpc_id                     = "${tencentcloud_vpc.foo.id}"
  subnet_id                  = "${tencentcloud_subnet.foo.id}"
}

resource "tencentcloud_clb_instance" "example" {
  clb_name                  = "example"
  network_type              = "${var.network_type}"
  project_id                = 0
  vpc_id                    = "${tencentcloud_vpc.foo.id}"
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = "${tencentcloud_vpc.foo.id}"
  security_groups           = ["${tencentcloud_security_group.foo.id}"]
}

resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = "${tencentcloud_clb_instance.example.id}"
  listener_name              = "listener_tcp"
  port                       = 22
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
}

resource "tencentcloud_clb_attachment" "attachment_tcp" {
  clb_id      = "${tencentcloud_clb_instance.example.id}"
  listener_id = "${tencentcloud_clb_listener.listener_tcp.id}"

  targets {
    instance_id = "${tencentcloud_instance.foo.id}"
    port        = 22
    weight      = 10
  }
}

resource "tencentcloud_clb_listener" "listener_https" {
  clb_id               = "${tencentcloud_clb_instance.example.id}"
  listener_name        = "listener_https"
  port                 = 443
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "VfqO4zkB"
}

resource "tencentcloud_clb_listener_rule" "rule_https" {
  clb_id              = "${tencentcloud_clb_instance.example.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_https.id}"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_attachment" "attachment_https" {
  clb_id      = "${tencentcloud_clb_instance.example.id}"
  listener_id = "${tencentcloud_clb_listener.listener_https.id}"
  rule_id     = "${tencentcloud_clb_listener_rule.rule_https.id}"

  targets {
    instance_id = "${tencentcloud_instance.foo.id}"
    port        = 443
    weight      = 10
  }
}

resource "tencentcloud_clb_listener" "listener_http_src" {
  clb_id        = "${tencentcloud_clb_instance.example.id}"
  port          = 8080
  protocol      = "HTTP"
  listener_name = "listener_http_src"
}

resource "tencentcloud_clb_listener_rule" "rule_http_src" {
  clb_id              = "${tencentcloud_clb_instance.example.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_http_src.id}"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_listener" "listener_http_dst" {
  clb_id        = "${tencentcloud_clb_instance.example.id}"
  port          = 80
  protocol      = "HTTP"
  listener_name = "listener_http_dst"
}

resource "tencentcloud_clb_listener_rule" "rule_http_dst" {
  clb_id              = "${tencentcloud_clb_instance.example.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_http_dst.id}"
  domain              = "abcd.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_redirection" "redirection_http" {
  clb_id             = "${tencentcloud_clb_instance.example.id}"
  source_listener_id = "${tencentcloud_clb_listener.listener_http_src.id}"
  target_listener_id = "${tencentcloud_clb_listener.listener_http_dst.id}"
  source_rule_id     = "${tencentcloud_clb_listener_rule.rule_http_src.id}"
  target_rule_id     = "${tencentcloud_clb_listener_rule.rule_http_dst.id}"
}

data "tencentcloud_clb_instances" "instances" {
  clb_id = "${tencentcloud_clb_instance.example.id}"
}

data "tencentcloud_clb_listeners" "listeners" {
  clb_id      = "${tencentcloud_clb_instance.example.id}"
  listener_id = "${tencentcloud_clb_listener.listener_tcp.id}"
}

data "tencentcloud_clb_listener_rules" "rules" {
  clb_id      = "${tencentcloud_clb_instance.example.id}"
  listener_id = "${tencentcloud_clb_listener.listener_https.id}"
  domain      = "${tencentcloud_clb_listener_rule.rule_https.domain}"
  url         = "${tencentcloud_clb_listener_rule.rule_https.url}"
}

data "tencentcloud_clb_attachments" "attachments" {
  clb_id      = "${tencentcloud_clb_instance.example.id}"
  listener_id = "${tencentcloud_clb_listener.listener_https.id}"
  rule_id     = "${tencentcloud_clb_attachment.attachment_https.id}"
}

data "tencentcloud_clb_redirections" "redirections" {
  clb_id             = "${tencentcloud_clb_instance.example.id}"
  source_listener_id = "${tencentcloud_clb_redirection.redirection_http.source_listener_id}"
  source_rule_id     = "${tencentcloud_clb_redirection.redirection_http.source_rule_id}"
}
