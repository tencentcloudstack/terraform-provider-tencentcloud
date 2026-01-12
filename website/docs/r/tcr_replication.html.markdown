---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_replication"
sidebar_current: "docs-tencentcloud-resource-tcr_replication"
description: |-
  Provides a resource to create a TCR replication
---

# tencentcloud_tcr_replication

Provides a resource to create a TCR replication

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

resource "tencentcloud_tcr_replication" "example" {
  source_registry_id      = tencentcloud_tcr_instance.example_gz.id
  destination_registry_id = tencentcloud_tcr_instance.example_sh.id
  rule {
    name           = "tf-example"
    dest_namespace = tencentcloud_tcr_namespace.example_sh.name
    override       = true
    deletion       = true
    filters {
      type  = "name"
      value = join("/", [tencentcloud_tcr_namespace.example_gz.name, "**"])
    }
  }

  destination_region_id = 1
  description           = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `destination_registry_id` - (Required, String, ForceNew) Destination instance ID.
* `rule` - (Required, List, ForceNew) Synchronization rule.
* `source_registry_id` - (Required, String, ForceNew) Source instance ID.
* `description` - (Optional, String, ForceNew) Rule description.
* `destination_region_id` - (Optional, Int, ForceNew) Region ID of the destination instance. For example, `1` represents Guangzhou.
* `peer_replication_option` - (Optional, List, ForceNew) Configuration of the synchronization rule.

The `filters` object of `rule` supports the following:

* `type` - (Required, String, ForceNew) Type (`name`, `tag` and `resource`).
* `value` - (Optional, String, ForceNew) It is left blank by default. If the type is `resource` it supports `image`, `chart`, and an empty string. If the type is `name` it supports Namespace name/**, Namespace name/Repository name.

The `peer_replication_option` object supports the following:

* `enable_peer_replication` - (Required, Bool, ForceNew) Whether to enable cross-account synchronization.
* `peer_registry_token` - (Required, String, ForceNew) Permanent access Token for the destination instance.
* `peer_registry_uin` - (Required, String, ForceNew) UIN of the destination instance.

The `rule` object supports the following:

* `dest_namespace` - (Required, String, ForceNew) Destination namespace.
* `filters` - (Required, List, ForceNew) Synchronization filters.
* `name` - (Required, String, ForceNew) Name of synchronization rule.
* `override` - (Required, Bool, ForceNew) Whether to override.
* `deletion` - (Optional, Bool, ForceNew) Whether synchronous deletion event.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



