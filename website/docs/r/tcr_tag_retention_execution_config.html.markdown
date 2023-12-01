---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_tag_retention_execution_config"
sidebar_current: "docs-tencentcloud-resource-tcr_tag_retention_execution_config"
description: |-
  Provides a resource to configure a tcr tag retention execution.
---

# tencentcloud_tcr_tag_retention_execution_config

Provides a resource to configure a tcr tag retention execution.

## Example Usage

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_example_ns_retention"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_tag_retention_rule" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  retention_rule {
    key   = "nDaysSinceLastPush"
    value = 2
  }
  cron_setting = "manual"
  disabled     = true
}

resource "tencentcloud_tcr_tag_retention_execution_config" "example" {
  registry_id  = tencentcloud_tcr_tag_retention_rule.example.registry_id
  retention_id = tencentcloud_tcr_tag_retention_rule.example.retention_id
  dry_run      = false
}
```

## Argument Reference

The following arguments are supported:

* `registry_id` - (Required, String) instance id.
* `retention_id` - (Required, Int) retention id.
* `dry_run` - (Optional, Bool) Whether to simulate execution, the default value is false, that is, non-simulation execution.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `execution_id` - execution id.


