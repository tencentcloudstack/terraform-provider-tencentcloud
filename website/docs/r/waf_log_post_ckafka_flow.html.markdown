---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_log_post_ckafka_flow"
sidebar_current: "docs-tencentcloud-resource-waf_log_post_ckafka_flow"
description: |-
  Provides a resource to create a WAF log post ckafka flow
---

# tencentcloud_waf_log_post_ckafka_flow

Provides a resource to create a WAF log post ckafka flow

## Example Usage

### If vip_type is 1

```hcl
resource "tencentcloud_waf_log_post_ckafka_flow" "example" {
  ckafka_region = "ap-guangzhou"
  ckafka_id     = "ckafka-k9m5vwar"
  brokers       = "ckafka-k9m5vwar.ap-guangzhou.ckafka.tencentcloudmq.com:50000"
  compression   = "snappy"
  vip_type      = 1
  log_type      = 2
  topic         = "tf-test"
  kafka_version = "2.8.1"
}
```

### If vip_type is 2

```hcl
resource "tencentcloud_waf_log_post_ckafka_flow" "example" {
  ckafka_region = "ap-guangzhou"
  ckafka_id     = "ckafka-k9m5vwar"
  brokers       = "11.135.14.110:18737"
  compression   = "snappy"
  vip_type      = 2
  log_type      = 1
  topic         = "tf-test"
  kafka_version = "2.8.1"
}
```

## Argument Reference

The following arguments are supported:

* `brokers` - (Required, String) The supporting environment is IP:PORT, The external network environment is domain:PORT.
* `ckafka_id` - (Required, String) CKafka ID.
* `ckafka_region` - (Required, String) The region where CKafka is located for delivery.
* `compression` - (Required, String) Default to none, supports snappy, gzip, and lz4 compression, recommended snappy.
* `kafka_version` - (Required, String) Version number of Kafka cluster.
* `log_type` - (Required, Int) 1- Access log, 2- Attack log, the default is access log.
* `topic` - (Required, String) Theme name, default not to pass or pass empty string, default value is waf_post_access_log.
* `vip_type` - (Required, Int) 1. External network TGW, 2. Supporting environment, default is supporting environment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `flow_id` - Unique ID for post cls flow.
* `status` - Status 0- Off 1- On.


## Import

WAF log post ckafka flow can be imported using the id, e.g.

```
# If log_type is 1
terraform import tencentcloud_waf_log_post_ckafka_flow.example 100536#1

# If log_type is 2
terraform import tencentcloud_waf_log_post_ckafka_flow.example 100541#2
```

