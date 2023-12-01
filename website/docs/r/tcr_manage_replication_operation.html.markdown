---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_manage_replication_operation"
sidebar_current: "docs-tencentcloud-resource-tcr_manage_replication_operation"
description: |-
  Provides a resource to start a tcr instance replication operation
---

# tencentcloud_tcr_manage_replication_operation

Provides a resource to start a tcr instance replication operation

## Example Usage

### Sync source tcr instance to target instance

Synchronize an existing tcr instance to the destination instance. This operation is often used in the cross-multiple region scenario.
Assume you have had two TCR instances before this operation. This example shows how to sync a tcr instance from ap-guangzhou(gz) to ap-shanghai(sh).

```hcl
# tcr instance on ap-guangzhou
resource "tencentcloud_tcr_instance" "example_gz" {
  name          = "tf-example-tcr-gz"
  instance_type = "premium"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example_gz" {
  instance_id    = tencentcloud_tcr_instance.example_gz.id
  name           = "tf_example_ns_gz"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

# tcr instance on ap-shanghai
resource "tencentcloud_tcr_instance" "example_sh" {
  name          = "tf-example-tcr-sh"
  instance_type = "premium"
  delete_bucket = true
}

resource "tencentcloud_tcr_namespace" "example_sh" {
  instance_id    = tencentcloud_tcr_instance.example_sh.id
  name           = "tf_example_ns_sh"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}
```



```hcl
# Run this on region ap-guangzhou
locals {
  src_id  = tencentcloud_tcr_instance.example_gz.id
  dest_id = tencentcloud_tcr_instance.example_sh.id
  src_ns  = tencentcloud_tcr_namespace.example_gz.name
  dest_ns = tencentcloud_tcr_instance.example_sh.id
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

resource "tencentcloud_tcr_manage_replication_operation" "example_sync" {
  source_registry_id      = local.src_id
  destination_registry_id = local.dest_id
  rule {
    name           = "tf_example_sync_gz_to_sh"
    dest_namespace = local.dest_ns
    override       = true
    filters {
      type  = "name"
      value = join("/", [local.src_ns, "**"])
    }
    filters {
      type  = "tag"
      value = ""
    }
    filters {
      type  = "resource"
      value = ""
    }
  }
  description           = "example for tcr sync operation"
  destination_region_id = var.tcr_region_map["ap-shanghai"] # 4
  peer_replication_option {
    peer_registry_uin       = ""
    peer_registry_token     = ""
    enable_peer_replication = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `destination_registry_id` - (Required, String, ForceNew) copy destination instance Id.
* `rule` - (Required, List, ForceNew) synchronization rules.
* `source_registry_id` - (Required, String, ForceNew) copy source instance Id.
* `description` - (Optional, String, ForceNew) rule description.
* `destination_region_id` - (Optional, Int, ForceNew) the region ID of the target instance, such as Guangzhou is 1.
* `peer_replication_option` - (Optional, List, ForceNew) enable synchronization of configuration items across master account instances.

The `filters` object supports the following:

* `type` - (Required, String) type (name, tag, and resource).
* `value` - (Optional, String) empty by default.

The `peer_replication_option` object supports the following:

* `enable_peer_replication` - (Required, Bool) whether to enable cross-master account instance synchronization.
* `peer_registry_token` - (Required, String) access permanent token of the instance to be synchronized.
* `peer_registry_uin` - (Required, String) uin of the instance to be synchronized.

The `rule` object supports the following:

* `dest_namespace` - (Required, String) target namespace.
* `filters` - (Required, List) sync filters.
* `name` - (Required, String) synchronization rule names.
* `override` - (Required, Bool) whether to cover.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



