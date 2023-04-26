---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_manage_replication_operation"
sidebar_current: "docs-tencentcloud-resource-tcr_manage_replication_operation"
description: |-
  Provides a resource to create a tcr manage_replication_operation
---

# tencentcloud_tcr_manage_replication_operation

Provides a resource to create a tcr manage_replication_operation

## Example Usage

```hcl
resource "tencentcloud_tcr_instance" "mytcr_dest" {
  name          = "tf-test-tcr-%s"
  instance_type = "premium"
  delete_bucket = true
}

resource "tencentcloud_tcr_namespace" "myns_dest" {
  instance_id    = tencentcloud_tcr_instance.mytcr_dest.id
  name           = "tf_test_ns_dest"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_manage_replication_operation" "my_replica" {
  source_registry_id      = local.tcr_id
  destination_registry_id = tencentcloud_tcr_instance.mytcr_dest.id
  rule {
    name           = "test_sync_%d"
    dest_namespace = tencentcloud_tcr_namespace.myns_dest.name
    override       = true
    filters {
      type  = "name"
      value = join("/", [var.tcr_namespace, "**"])
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
  description           = "this is the tcr sync operation"
  destination_region_id = 1 // "ap-guangzhou"
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



