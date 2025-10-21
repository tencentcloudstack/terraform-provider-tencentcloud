---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_manage_grafana_attachment"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_manage_grafana_attachment"
description: |-
  Provides a resource to create a monitor tmp_manage_grafana_attachment
---

# tencentcloud_monitor_tmp_manage_grafana_attachment

Provides a resource to create a monitor tmp_manage_grafana_attachment

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

resource "tencentcloud_monitor_grafana_instance" "foo" {
  instance_name         = "tf-grafana"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_ids            = [tencentcloud_subnet.subnet.id]
  grafana_init_password = "1234567890"
  enable_internet       = false
  is_destroy            = true

  tags = {
    "createdBy" = "test"
  }
}

resource "tencentcloud_monitor_tmp_manage_grafana_attachment" "foo" {
  grafana_id  = tencentcloud_monitor_grafana_instance.foo.id
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `grafana_id` - (Required, String, ForceNew) Grafana instance ID.
* `instance_id` - (Required, String, ForceNew) Prometheus instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor tmp_manage_grafana_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_manage_grafana_attachment.manage_grafana_attachment prom-xxxxxxxx
```

