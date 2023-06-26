---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_tag_retention_execution_config"
sidebar_current: "docs-tencentcloud-resource-tcr_tag_retention_execution_config"
description: |-
  Provides a resource to create a tcr tag_retention_execution_config
---

# tencentcloud_tcr_tag_retention_execution_config

Provides a resource to create a tcr tag_retention_execution_config

## Example Usage

```hcl
resource "tencentcloud_tcr_namespace" "my_ns" {
  instance_id    = tencentcloud_tcr_instance.mytcr_retention.id
  name           = "tf_test_ns_retention"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_tag_retention_rule" "my_rule" {
  registry_id    = tencentcloud_tcr_instance.mytcr_retention.id
  namespace_name = tencentcloud_tcr_namespace.my_ns.name
  retention_rule {
    key   = "nDaysSinceLastPush"
    value = 2
  }
  cron_setting = "manual"
  disabled     = true
}

resource "tencentcloud_tcr_tag_retention_execution_config" "tag_retention_execution_config" {
  registry_id  = tencentcloud_tcr_tag_retention_rule.my_rule.registry_id
  retention_id = tencentcloud_tcr_tag_retention_rule.my_rule.retention_id
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


