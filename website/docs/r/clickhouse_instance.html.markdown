---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_instance"
sidebar_current: "docs-tencentcloud-resource-clickhouse_instance"
description: |-
  Provides a resource to create a clickhouse instance.
---

# tencentcloud_clickhouse_instance

Provides a resource to create a clickhouse instance.

## Example Usage

```hcl
resource "tencentcloud_clickhouse_instance" "cdwch_instance" {
  zone            = "ap-guangzhou-6"
  ha_flag         = true
  vpc_id          = "vpc-xxxxxx"
  subnet_id       = "subnet-xxxxxx"
  product_version = "21.8.12.29"
  data_spec {
    spec_name = "SCH6"
    count     = 2
    disk_size = 300
  }
  common_spec {
    spec_name = "SCH6"
    count     = 3
    disk_size = 300
  }
  charge_type   = "POSTPAID_BY_HOUR"
  instance_name = "tf-test-clickhouse"
}
```

### PREPAID instance

```hcl
resource "tencentcloud_clickhouse_instance" "cdwch_instance_prepaid" {
  zone            = "ap-guangzhou-6"
  ha_flag         = true
  vpc_id          = "vpc-xxxxxx"
  subnet_id       = "subnet-xxxxxx"
  product_version = "21.8.12.29"
  data_spec {
    spec_name = "SCH6"
    count     = 2
    disk_size = 300
  }
  common_spec {
    spec_name = "SCH6"
    count     = 3
    disk_size = 300
  }
  charge_type   = "PREPAID"
  renew_flag    = 1
  time_span     = 1
  instance_name = "tf-test-clickhouse-prepaid"
}
```

## Argument Reference

The following arguments are supported:

* `charge_type` - (Required, String) Billing type: `PREPAID` prepaid, `POSTPAID_BY_HOUR` postpaid.
* `data_spec` - (Required, List) Data spec.
* `ha_flag` - (Required, Bool) Whether it is highly available.
* `instance_name` - (Required, String) Instance name.
* `product_version` - (Required, String) Product version.
* `subnet_id` - (Required, String) Subnet.
* `vpc_id` - (Required, String) Private network.
* `zone` - (Required, String) Availability zone.
* `cls_log_set_id` - (Optional, String) CLS log set id.
* `common_spec` - (Optional, List) ZK node.
* `cos_bucket_name` - (Optional, String) COS bucket name.
* `ha_zk` - (Optional, Bool) Whether ZK is highly available.
* `mount_disk_type` - (Optional, Int) Whether it is mounted on a bare disk.
* `renew_flag` - (Optional, Int) PREPAID needs to be passed. Whether to renew automatically. 1 means auto renewal is enabled.
* `tags` - (Optional, Map) Tag description list.
* `time_span` - (Optional, Int) Prepaid needs to be delivered, billing time length, how many months.

The `common_spec` object supports the following:

* `count` - (Required, Int) Node count. NOTE: Only support value 3.
* `disk_size` - (Required, Int) Disk size.
* `spec_name` - (Required, String) Spec name.

The `data_spec` object supports the following:

* `count` - (Required, Int) Data spec count.
* `disk_size` - (Required, Int) Disk size.
* `spec_name` - (Required, String) Spec name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `expire_time` - Expire time.


## Import

Clickhouse instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clickhouse_instance.foo cdwch-xxxxxx
```

