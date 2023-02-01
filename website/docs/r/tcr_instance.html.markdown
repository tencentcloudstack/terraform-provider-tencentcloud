---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_instance"
sidebar_current: "docs-tencentcloud-resource-tcr_instance"
description: |-
  Use this resource to create tcr instance.
---

# tencentcloud_tcr_instance

Use this resource to create tcr instance.

## Example Usage

```hcl
resource "tencentcloud_tcr_instance" "foo" {
  name          = "example"
  instance_type = "basic"

  tags = {
    test = "tf"
  }
}
```

Using public network access whitelist

```hcl
resource "tencentcloud_tcr_instance" "foo" {
  name                  = "example"
  instance_type         = "basic"
  open_public_operation = true
  security_policy {
    cidr_block = "10.0.0.1/24"
  }
  security_policy {
    cidr_block = "192.168.1.1"
  }
}
```

Create with Replications

```hcl
resource "tencentcloud_tcr_instance" "foo" {
  name          = "example"
  instance_type = "premium"
  replications {
    region_id = var.tcr_region_map["ap-guangzhou"] # 1
  }
  replications {
    region_id = var.tcr_region_map["ap-singapore"] # 9
  }
}

variable "tcr_region_map" {
  default = {
    "ap-guangzhou"     = 1
    "ap-shanghai"      = 4
    "ap-hongkong"      = 5
    "ap-beijing"       = 8
    "ap-singapore"     = 9
    "na-siliconvalley" = 15
    "ap-chengdu"       = 16
    "eu-frankfurt"     = 17
    "ap-seoul"         = 18
    "ap-chongqing"     = 19
    "ap-mumbai"        = 21
    "na-ashburn"       = 22
    "ap-bangkok"       = 23
    "eu-moscow"        = 24
    "ap-tokyo"         = 25
    "ap-nanjing"       = 33
    "ap-taipei"        = 39
    "ap-jakarta"       = 72
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, String) TCR types. Valid values are: `standard`, `basic`, `premium`.
* `name` - (Required, String, ForceNew) Name of the TCR instance.
* `delete_bucket` - (Optional, Bool) Indicate to delete the COS bucket which is auto-created with the instance or not.
* `instance_charge_type_prepaid_period` - (Optional, Int) Length of time to purchase an instance (in month). Must set when registry_charge_type is prepaid.
* `instance_charge_type_prepaid_renew_flag` - (Optional, Int) Auto renewal flag. 1: manual renewal, 2: automatic renewal, 3: no renewal and no notification. Must set when registry_charge_type is prepaid.
* `open_public_operation` - (Optional, Bool) Control public network access.
* `registry_charge_type` - (Optional, Int) Charge type of instance. 1: postpaid; 2: prepaid. Default is postpaid.
* `replications` - (Optional, List) Specify List of instance Replications, premium only. The available [source region list](https://www.tencentcloud.com/document/api/1051/41101) is here.
* `security_policy` - (Optional, Set) Public network access allowlist policies of the TCR instance. Only available when `open_public_operation` is `true`.
* `tags` - (Optional, Map) The available tags within this TCR instance.

The `replications` object supports the following:

* `region_id` - (Optional, Int) Replication region ID, check the example at the top of page to find out id of region.
* `syn_tag` - (Optional, Bool) Specify whether to sync TCR cloud tags to COS Bucket. NOTE: You have to specify when adding, modifying will be ignored for now.

The `security_policy` object supports the following:

* `cidr_block` - (Optional, String) The public network IP address of the access source.
* `description` - (Optional, String) Remarks of policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `expired_at` - Instance expiration time (prepaid).
* `internal_end_point` - Internal address for access of the TCR instance.
* `public_domain` - Public address for access of the TCR instance.
* `public_status` - Status of the TCR instance public network access.
* `status` - Status of the TCR instance.


## Import

tcr instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_instance.foo cls-cda1iex1
```

