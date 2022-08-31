# Provider
terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
    }
  }
}

provider "tencentcloud" {
  region = "ap-guangzhou"
}

# 站点
resource "tencentcloud_teo_zone" "sfurnace_work" {
  name           = "sfurnace.work"
  plan_type      = "ent_cm_with_bot"
  type           = "full"
  paused         = false
  cname_speed_up = "enabled"

  vanity_name_servers {
    switch  = "on"
    servers = ["ns1.sfurnace.work", "ns2.sfurnace.work"]
  }
}

# 站点全局配置
resource "tencentcloud_teo_zone_setting" "sfurnace_work" {
  zone_id = tencentcloud_teo_zone.sfurnace_work.id

  cache {
    follow_origin {
      switch = "off"
    }

    no_cache {
      switch = "off"
    }
  }

  cache_key {
    full_url_cache = "off"
    ignore_case    = "on"

    query_string {
      action = "excludeCustom"
      switch = "on"
      value  = ["test", "apple"]
    }
  }

  cache_prefresh {
    percent = 90
    switch  = "off"
  }

  client_ip_header {
    switch = "off"
  }

  compression {
    switch = "off"
  }

  force_redirect {
    redirect_status_code = 302
    switch               = "on"
  }

  https {
    http2         = "on"
    ocsp_stapling = "off"
    tls_version   = [
      "TLSv1.2",
      "TLSv1.3",
    ]

    hsts {
      include_sub_domains = "off"
      max_age             = 0
      preload             = "off"
      switch              = "off"
    }
  }

  max_age {
    follow_origin = "off"
    max_age_time  = 600
  }

  offline_cache {
    switch = "off"
  }

  origin {
    origin_pull_protocol = "follow"
  }

  post_max_size {
    max_size = 524288000
    switch   = "on"
  }

  quic {
    switch = "on"
  }

  smart_routing {
    switch = "on"
  }

  upstream_http2 {
    switch = "off"
  }

  web_socket {
    switch  = "off"
    timeout = 30
  }
}

# DNS 记录
resource "tencentcloud_teo_dns_record" "sfurnace_work" {
  zone_id     = tencentcloud_teo_zone.sfurnace_work.id
  record_type = "A"
  name        = "sfurnace.work"
  mode        = "proxied"
  content     = "2.2.2.2"
  ttl         = 80
}

resource "tencentcloud_teo_dns_record" "www_sfurnace_work" {
  zone_id     = tencentcloud_teo_zone.sfurnace_work.id
  record_type = "A"
  name        = "www.sfurnace.work"
  mode        = "proxied"
  content     = "1.1.1.1"
  ttl         = 120
}

resource "tencentcloud_teo_dns_record" "vstest_sfurnace_work" {
  zone_id     = tencentcloud_teo_zone.sfurnace_work.id
  record_type = "A"
  name        = "vstest.sfurnace.work"
  mode        = "proxied"
  content     = "3.3.3.3"
  ttl         = 120
}

## DNS SEC
##resource "tencentcloud_teo_dns_sec" "sfurnace_work" {
##  zone_id = tencentcloud_teo_zone.sfurnace_work.id
##  status  = "disabled"
##}
#
## 源站组
#resource "tencentcloud_teo_origin_group" "group0" {
#  zone_id     = tencentcloud_teo_zone.sfurnace_work.id
#  origin_name = "group0"
#  origin_type = "self"
#  type        = "weight"
#
#  dynamic "record" {
#    for_each = local.group0
#    content {
#      record = record.value["record"]
#      port   = record.value["port"]
#      weight = record.value["weight"]
#      area   = []
#    }
#  }
#}
#
#locals {
#  group0 = [
#    {
#      "record" = "1.1.1.1"
#      "port"   = 80
#      "weight" = 30
#    }, {
#      "record" = "2.2.2.2"
#      "port"   = 443
#      "weight" = 70
#    }
#  ]
#}
#
#resource "tencentcloud_teo_origin_group" "group1" {
#  zone_id     = tencentcloud_teo_zone.sfurnace_work.id
#  origin_name = "group1"
#  origin_type = "self"
#  type        = "area"
#
#  dynamic "record" {
#    for_each = local.group1
#    content {
#      record = record.value["record"]
#      port   = record.value["port"]
#      area   = record.value["area"]
#      weight = 0
#    }
#  }
#}
#
#locals {
#  group1 = [
#    {
#      "record" = "1.1.1.1"
#      "port"   = 80
#      "area"   = []
#    }, {
#      "record" = "2.2.2.2"
#      "port"   = 80
#      "area"   = ["ER"]
#    }
#  ]
#}
#
## 负载均衡
#resource "tencentcloud_teo_load_balancing" "lb0" {
#  zone_id = tencentcloud_teo_zone.sfurnace_work.id
#
#  host      = "sfurnace.work"
#  origin_id = [
#    split("#", tencentcloud_teo_origin_group.group0.id)[1]
#  ]
#  ttl  = 600
#  type = "proxied"
#}
#
## 四层代理
#resource "tencentcloud_teo_application_proxy" "app0" {
#  zone_id   = tencentcloud_teo_zone.sfurnace_work.id
#  zone_name = "sfurnace.work"
#
#  accelerate_type      = 1
#  security_type        = 1
#  plat_type            = "domain"
#  proxy_name           = "www.sfurnace.work"
#  proxy_type           = "hostname"
#  session_persist_time = 2400
#}
#
#resource "tencentcloud_teo_application_proxy_rule" "app0_rule0" {
#  zone_id  = tencentcloud_teo_zone.sfurnace_work.id
#  proxy_id = tencentcloud_teo_application_proxy.app0.proxy_id
#
#  forward_client_ip = "TOA"
#  origin_type       = "custom"
#  origin_value      = [
#    "1.1.1.1:80",
#  ]
#  port = [
#    "80",
#  ]
#  proto           = "TCP"
#  session_persist = false
#}
#
##resource "tencentcloud_teo_application_proxy_rule" "app0_rule1" {
##  zone_id  = tencentcloud_teo_zone.sfurnace_work.id
##  proxy_id = tencentcloud_teo_application_proxy.app0.proxy_id
##
##  # forward_client_ip = "PPV2"
##  origin_type  = "origins"
##  origin_value = [
##    "origin-cd3fefb0-28e2-11ed-8d28-5254005a52aa",
##  ]
##  port = [
##    "80",
##  ]
##  proto           = "TCP"
##  session_persist = true
##}
#
## 规则引擎
#resource "tencentcloud_teo_rule_engine" "sfurnace_work" {
#  zone_id   = tencentcloud_teo_zone.sfurnace_work.id
#  rule_name = "规则0"
#  status    = "enable"
#
#  rules {
#    conditions {
#      conditions {
#        operator = "equal"
#        target   = "host"
#        values   = [
#          "www.sfurnace.work",
#        ]
#      }
#    }
#
#    actions {
#      normal_action {
#        action = "MaxAge"
#
#        parameters {
#          name   = "FollowOrigin"
#          values = [
#            "on",
#          ]
#        }
#        parameters {
#          name   = "MaxAgeTime"
#          values = [
#            "0",
#          ]
#        }
#      }
#    }
#  }
#}
#
#resource "tencentcloud_teo_rule_engine" "sfurnace_work_1" {
#  zone_id   = tencentcloud_teo_zone.sfurnace_work.id
#  rule_name = "规则1"
#  status    = "disable"
#
#  rules {
#    conditions {
#      conditions {
#        operator = "equal"
#        target   = "extension"
#        values   = [
#          "mp4",
#        ]
#      }
#      conditions {
#        operator = "equal"
#        target   = "host"
#        values   = [
#          "sfurnace.work",
#        ]
#      }
#    }
#
#    actions {
#      normal_action {
#        action = "CachePrefresh"
#
#        parameters {
#          name   = "Switch"
#          values = [
#            "on",
#          ]
#        }
#        parameters {
#          name   = "Percent"
#          values = [
#            "80",
#          ]
#        }
#      }
#    }
#
#    actions {
#      normal_action {
#        action = "CacheKey"
#
#        parameters {
#          name   = "Type"
#          values = [
#            "Header",
#          ]
#        }
#        parameters {
#          name   = "Switch"
#          values = [
#            "on",
#          ]
#        }
#        parameters {
#          name   = "Value"
#          values = [
#            "Duck",
#          ]
#        }
#      }
#    }
#  }
#}
#
#resource "tencentcloud_teo_rule_engine" "sfurnace_work_2" {
#  zone_id   = tencentcloud_teo_zone.sfurnace_work.id
#  rule_name = "规则3"
#  status    = "enable"
#
#  rules {
#    actions {
#      rewrite_action {
#        action = "ResponseHeader"
#
#        parameters {
#          action = "add"
#          name   = "A"
#          values = [
#            "A",
#          ]
#        }
#      }
#    }
#
#    conditions {
#      conditions {
#        operator = "equal"
#        target   = "host"
#        values   = [
#          "www.sfurnace.work",
#        ]
#      }
#    }
#  }
#}

## 默认证书
#resource "tencentcloud_teo_default_certificate" "sfurnace_work" {
#  zone_id = tencentcloud_teo_zone.sfurnace_work.id
#
#  cert_info {
#    cert_id = "teo-28i46c1gtmkl"
#    status  = "deployed"
#  }
#}
#
## 域名证书
#resource "tencentcloud_teo_host_certificate" "vstest_sfurnace_work" {
#  zone_id = tencentcloud_teo_zone.sfurnace_work.id
#  host    = tencentcloud_teo_dns_record.vstest_sfurnace_work.name
#
#  cert_info {
#    cert_id = "yqWPPbs7"
#    status  = "deployed"
#  }
#}

# DDoS 安全策略
#resource "tencentcloud_teo_ddos_policy" "sfurnace_work" {
#  zone_id   = tencentcloud_teo_zone.sfurnace_work.id
#  policy_id = 706
#
#  ddos_rule {
#    switch = "on"
#
#    ddos_geo_ip {
#      region_id = [226074, 1820814, 3190538]
#    }
#  }
#}

# Web/Bot 安全策略
#resource "tencentcloud_teo_security_policy" "sfurnace_work" {
#  zone_id = tencentcloud_teo_zone.sfurnace_work.id
#  entity  = "sfurnace.work"
#
#  config {
#    switch_config {
#      web_switch = "on"
#    }
#  }
#}