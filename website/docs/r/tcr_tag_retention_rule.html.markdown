---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_tag_retention_rule"
sidebar_current: "docs-tencentcloud-resource-tcr_tag_retention_rule"
description: |-
  Provides a resource to create a tcr tag retention rule.
---

# tencentcloud_tcr_tag_retention_rule

Provides a resource to create a tcr tag retention rule.

## Example Usage

### Create and enable a tcr tag retention rule instance

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

resource "tencentcloud_tcr_tag_retention_rule" "my_rule" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  retention_rule {
    key   = "nDaysSinceLastPush"
    value = 2
  }
  cron_setting = "daily"
}
```

## Argument Reference

The following arguments are supported:

* `cron_setting` - (Required, String, ForceNew) Execution cycle, currently only available selections are: manual; daily; weekly; monthly.
* `namespace_name` - (Required, String) The Name of the namespace.
* `registry_id` - (Required, String) The main instance ID.
* `retention_rule` - (Required, List) Retention Policy.
* `disabled` - (Optional, Bool) Whether to disable the rule, with the default value of false.

The `retention_rule` object supports the following:

* `key` - (Required, String) The supported policies are latestPushedK (retain the latest `k` pushed versions) and nDaysSinceLastPush (retain pushed versions within the last `n` days).
* `value` - (Required, Int) corresponding values for rule settings.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `retention_id` - The ID of the retention task.


