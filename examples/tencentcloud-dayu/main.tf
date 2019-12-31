resource "tencentcloud_dayu_ddos_policy" "example" {
  resource_type = "bgpip"
  name          = "example"

  drop_options {
    drop_tcp           = true
    drop_udp           = true
    drop_icmp          = true
    drop_other         = true
    drop_abroad        = true
    check_sync_conn    = true
    source_new_limit   = 100
    dst_new_limit      = 100
    source_conn_limit  = 100
    dst_conn_limit     = 100
    tcp_mbps_limit     = 100
    udp_mbps_limit     = 100
    icmp_mbps_limit    = 100
    other_mbps_limit   = 100
    bad_conn_threshold = 100
    null_conn_enable   = true
    conn_timeout       = 500
    syn_rate           = 50
    syn_limit          = 100
  }

  black_white_ips {
    ip   = "1.1.1.1"
    type = "black"
  }

  port_limits {
    start_port = "2000"
    end_port   = "2500"
    protocol   = "all"
    action     = "drop"
    kind       = 1
  }

  packet_filters {
    protocol       = "tcp"
    action         = "drop"
    d_start_port   = 1000
    d_end_port     = 1500
    s_start_port   = 2000
    s_end_port     = 2500
    pkt_length_max = 1400
    pkt_length_min = 1000
    is_include     = true
    match_begin    = "begin_l5"
    match_type     = "pcre"
    depth          = 1000
    offset         = 500
  }

  water_prints {
    tcp_port_list = ["2000-3000", "3500-4000"]
    udp_port_list = ["5000-6000"]
    offset        = 50
    auto_remove   = true
    open_switch   = true
  }
}

resource "tencentcloud_dayu_ddos_policy_case" "example" {
  resource_type       = "bgp-multip"
  name                = "example"
  platform_types      = ["PC", "MOBILE"]
  app_type            = "WEB"
  app_protocols       = ["tcp", "udp"]
  tcp_start_port      = "1000"
  tcp_end_port        = "2000"
  udp_start_port      = "3000"
  udp_end_port        = "4000"
  has_abroad          = "yes"
  has_initiate_tcp    = "yes"
  has_initiate_udp    = "yes"
  peer_tcp_port       = "1111"
  peer_udp_port       = "3333"
  tcp_foot_print      = "511"
  udp_foot_print      = "500"
  web_api_urls        = ["abc.com", "test.cn/aaa.png"]
  min_tcp_package_len = "1000"
  max_tcp_package_len = "1200"
  min_udp_package_len = "1000"
  max_udp_package_len = "1200"
  has_vpn             = "yes"
}

resource "tencentcloud_dayu_ddos_policy_attachment" "example" {
  resource_type = tencentcloud_dayu_ddos_policy.example.resource_type
  resource_id   = var.resource_bgpip
  policy_id     = tencentcloud_dayu_ddos_policy.example.policy_id
}

resource "tencentcloud_dayu_l4_rule" "example" {
  resource_type             = "net"
  resource_id               = var.resource_net
  name                      = "example"
  protocol                  = "TCP"
  source_port               = 80
  dest_port                 = 66
  source_type               = 2
  health_check_switch       = true
  health_check_timeout      = 30
  health_check_interval     = 35
  health_check_health_num   = 5
  health_check_unhealth_num = 10
  session_switch            = false
  session_time              = 300

  source_list {
    source = "1.1.1.1"
    weight = 100
  }
  source_list {
    source = "2.2.2.2"
    weight = 50
  }
}

resource "tencentcloud_dayu_l7_rule" "example" {
  resource_type             = "bgpip"
  resource_id               = var.resource_bgpip
  name                      = "example"
  domain                    = var.default_domain
  protocol                  = "https"
  switch                    = true
  source_type               = 2
  source_list               = ["1.1.1.1:80", "2.2.2.2"]
  ssl_id                    = var.resource_ssl
  health_check_switch       = true
  health_check_code         = 31
  health_check_interval     = 30
  health_check_method       = "GET"
  health_check_path         = "/"
  health_check_health_num   = 5
  health_check_unhealth_num = 10
}

resource "tencentcloud_dayu_cc_https_policy" "example" {
  resource_type = tencentcloud_dayu_l7_rule.example.resource_type
  resource_id   = tencentcloud_dayu_l7_rule.example.resource_id
  rule_id       = tencentcloud_dayu_l7_rule.example.rule_id
  domain        = tencentcloud_dayu_l7_rule.example.domain
  name          = "example"
  exe_mode      = "drop"
  switch        = true

  rule_list {
    skey     = "cgi"
    operator = "include"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "example_bgpip" {
  resource_type = "bgpip"
  resource_id   = var.resource_bgpip
  name          = "example_bgpip"
  smode         = "matching"
  exe_mode      = "drop"
  switch        = true
  rule_list {
    skey     = "host"
    operator = "include"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "example_net" {
  resource_type = "net"
  resource_id   = var.resource_net
  name          = "example_net"
  smode         = "matching"
  exe_mode      = "drop"
  switch        = true
  rule_list {
    skey     = "cgi"
    operator = "equal"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "example_bgpmultip" {
  resource_type = "bgp-multip"
  resource_id   = var.resource_bgpmultip
  name          = "example_bgpmultip"
  smode         = "matching"
  exe_mode      = "alg"
  switch        = true
  ip            = var.bgpmultip_ip

  rule_list {
    skey     = "referer"
    operator = "not_include"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "example_bgp" {
  resource_type = "bgp"
  resource_id   = var.resource_bgp
  name          = "example_bgp"
  smode         = "matching"
  exe_mode      = "alg"
  switch        = true

  rule_list {
    skey     = "ua"
    operator = "not_include"
    value    = "123"
  }
}

data "tencentcloud_dayu_cc_http_policies" "id_test" {
  resource_type = tencentcloud_dayu_cc_http_policy.example_bgpip.resource_type
  resource_id   = tencentcloud_dayu_cc_http_policy.example_bgpip.resource_id
  policy_id     = tencentcloud_dayu_cc_http_policy.example_bgpip.policy_id
}

data "tencentcloud_dayu_cc_http_policies" "name_test" {
  resource_type = tencentcloud_dayu_cc_http_policy.example_bgpip.resource_type
  resource_id   = tencentcloud_dayu_cc_http_policy.example_bgpip.resource_id
  name          = tencentcloud_dayu_cc_http_policy.example_bgpip.name
}

data "tencentcloud_dayu_cc_https_policies" "name_test" {
  resource_type = tencentcloud_dayu_cc_https_policy.example.resource_type
  resource_id   = tencentcloud_dayu_cc_https_policy.example.resource_id
  name          = tencentcloud_dayu_cc_https_policy.example.name
}

data "tencentcloud_dayu_cc_https_policies" "id_test" {
  resource_type = tencentcloud_dayu_cc_https_policy.example.resource_type
  resource_id   = tencentcloud_dayu_cc_https_policy.example.resource_id
  policy_id     = tencentcloud_dayu_cc_https_policy.example.policy_id
}

data "tencentcloud_dayu_ddos_policies" "id_test" {
  resource_type = tencentcloud_dayu_ddos_policy.example.resource_type
  policy_id     = tencentcloud_dayu_ddos_policy.example.policy_id
}

data "tencentcloud_dayu_ddos_policy_attachments" "foo_type" {
  resource_type = tencentcloud_dayu_ddos_policy_attachment.example.resource_type
}

data "tencentcloud_dayu_ddos_policy_attachments" "foo_resource" {
  resource_id   = tencentcloud_dayu_ddos_policy_attachment.example.resource_id
  resource_type = tencentcloud_dayu_ddos_policy_attachment.example.resource_type
}

data "tencentcloud_dayu_ddos_policy_attachments" "foo_policy" {
  resource_type = tencentcloud_dayu_ddos_policy_attachment.example.resource_type
  policy_id     = tencentcloud_dayu_ddos_policy_attachment.example.policy_id
}

data "tencentcloud_dayu_ddos_policy_cases" "id_test" {
  resource_type = tencentcloud_dayu_ddos_policy_case.example.resource_type
  scene_id      = tencentcloud_dayu_ddos_policy_case.example.scene_id
}

data "tencentcloud_dayu_l4_rules" "name_test" {
  resource_type = tencentcloud_dayu_l4_rule.example.resource_type
  resource_id   = tencentcloud_dayu_l4_rule.example.resource_id
  name          = tencentcloud_dayu_l4_rule.example.name
}

data "tencentcloud_dayu_l4_rules" "id_test" {
  resource_type = tencentcloud_dayu_l4_rule.example.resource_type
  resource_id   = tencentcloud_dayu_l4_rule.example.resource_id
  rule_id       = tencentcloud_dayu_l4_rule.example.rule_id
}

data "tencentcloud_dayu_l7_rules" "domain_test" {
  resource_type = tencentcloud_dayu_l7_rule.example.resource_type
  resource_id   = tencentcloud_dayu_l7_rule.example.resource_id
  domain        = tencentcloud_dayu_l7_rule.example.domain
}

data "tencentcloud_dayu_l7_rules" "id_test" {
  resource_type = tencentcloud_dayu_l7_rule.example.resource_type
  resource_id   = tencentcloud_dayu_l7_rule.example.resource_id
  rule_id       = tencentcloud_dayu_l7_rule.example.rule_id
}
