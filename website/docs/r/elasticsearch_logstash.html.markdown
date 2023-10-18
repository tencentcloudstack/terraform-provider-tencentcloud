---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_logstash"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_logstash"
description: |-
  Provides a resource to create a elasticsearch logstash
---

# tencentcloud_elasticsearch_logstash

Provides a resource to create a elasticsearch logstash

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_logstash" "logstash" {
  instance_name    = "logstash-test"
  zone             = "ap-guangzhou-6"
  logstash_version = "7.14.2"
  vpc_id           = "vpc-4owdpnwr"
  subnet_id        = "subnet-4o0zd840"
  node_num         = 1
  charge_type      = "POSTPAID_BY_HOUR"
  node_type        = "LOGSTASH.SA2.MEDIUM4"
  disk_type        = "CLOUD_SSD"
  disk_size        = 20
  license_type     = "xpack"
  operation_duration {
    periods    = [1, 2, 3, 4, 5, 6, 0]
    time_start = "02:00"
    time_end   = "06:00"
    time_zone  = "UTC+8"
  }
  tags = {
    tagKey = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, String) Instance name (compose of 1-50 letter, number, - or _).
* `logstash_version` - (Required, String) Instance version(6.8.13, 7.10.1).
* `subnet_id` - (Required, String) Subnet id.
* `vpc_id` - (Required, String) VPC id.
* `zone` - (Required, String) Available zone.
* `auto_voucher` - (Optional, Int) whether to use voucher auto, 1 when use, else 0.
* `charge_period` - (Optional, Int) Period when charged by months or years(unit depends on TimeUnit).
* `charge_type` - (Optional, String) Charge type. PREPAID: charged by months or years; POSTPAID_BY_HOUR: charged by hours; default vaule: POSTPAID_BY_HOUR.
* `disk_size` - (Optional, Int) node disk size (unit GB).
* `disk_type` - (Optional, String) Disk type. CLOUD_SSD: SSD cloud disk; CLOUD_PREMIUM: high hard energy cloud disk; default: CLOUD_SSD.
* `license_type` - (Optional, String) License type. oss: open source version; xpack:xpack version; default: xpack.
* `node_num` - (Optional, Int) Node num(range 2-50).
* `node_type` - (Optional, String) Node type. Valid values:
- LOGSTASH.S1.SMALL2: 1 core 2G;
- LOGSTASH.S1.MEDIUM4:2 core 4G;
- LOGSTASH.S1.MEDIUM8:2 core 8G;
- LOGSTASH.S1.LARGE16:4 core 16G;
- LOGSTASH.S1.2XLARGE32:8 core 32G;
- LOGSTASH.S1.4XLARGE32:16 core 32G;
- LOGSTASH.S1.4XLARGE64:16 core 64G.
* `operation_duration` - (Optional, List) operation time by tencent clound.
* `renew_flag` - (Optional, String) Automatic renewal flag. RENEW_FLAG_AUTO: auto renewal; RENEW_FLAG_MANUAL: do not renew automatically, users renew manually. It needs to be set when ChargeType is PREPAID. If this parameter is not passed, ordinary users will not renew automatically by default, and SVIP users will renew automatically.
* `tags` - (Optional, Map) Tag description list.
* `time_unit` - (Optional, String) charge time unit(set when ChargeType is PREPAID, default value: ms).
* `voucher_ids` - (Optional, Set: [`String`]) Voucher list(only can use one voucher by now).

The `operation_duration` object supports the following:

* `periods` - (Required, Set) day of week, from Monday to Sunday, value range: [0, 6]notes: may return null when missing.
* `time_end` - (Required, String) operation end time.
* `time_start` - (Required, String) operation start time.
* `time_zone` - (Required, String) time zone, for example: UTC+8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

elasticsearch logstash can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_logstash.logstash logstash_id
```

