---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_instance"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_instance"
description: |-
  Provides a resource to create a monitor tmpInstance
---

# tencentcloud_monitor_tmp_instance

Provides a resource to create a monitor tmpInstance

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

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

resource "tencentcloud_monitor_tmp_instance" "foo" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 30
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `data_retention_time` - (Required, Int) Data retention time.
* `instance_name` - (Required, String) Instance name.
* `subnet_id` - (Required, String) Subnet Id.
* `vpc_id` - (Required, String) Vpc Id.
* `zone` - (Required, String) Available zone.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_root_path` - Prometheus HTTP API root address.
* `ipv4_address` - Instance IPv4 address.
* `proxy_address` - Proxy address.
* `remote_write` - Prometheus remote write address.


## Import

monitor tmpInstance can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_instance.tmpInstance tmpInstance_id
```

