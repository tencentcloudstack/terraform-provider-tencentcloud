---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_sso_account"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_sso_account"
description: |-
  Provides a resource to create a monitor grafana ssoAccount
---

# tencentcloud_monitor_grafana_sso_account

Provides a resource to create a monitor grafana ssoAccount

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
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

resource "tencentcloud_monitor_grafana_instance" "foo" {
  instance_name         = "test-grafana"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_ids            = [tencentcloud_subnet.subnet.id]
  grafana_init_password = "1234567890"
  enable_internet       = false

  tags = {
    "createdBy" = "test"
  }
}

resource "tencentcloud_monitor_grafana_sso_account" "ssoAccount" {
  instance_id = tencentcloud_monitor_grafana_instance.foo.id
  user_id     = "111"
  notes       = "desc12222"
  role {
    organization = "Main Org."
    role         = "Admin"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) grafana instance id.
* `user_id` - (Required, String) sub account uin of specific user.
* `notes` - (Optional, String) account related description.
* `role` - (Optional, List) grafana role.

The `role` object supports the following:

* `organization` - (Required, String) Grafana organization id string.
* `role` - (Required, String) Grafana role, one of {Admin,Editor,Viewer}.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor grafana ssoAccount can be imported using the instance_id#user_id, e.g.
```
$ terraform import tencentcloud_monitor_grafana_sso_account.ssoAccount grafana-50nj6v00#111
```

