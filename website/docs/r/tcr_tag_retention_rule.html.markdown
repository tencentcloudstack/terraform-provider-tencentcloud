---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_tag_retention_rule"
sidebar_current: "docs-tencentcloud-resource-tcr_tag_retention_rule"
description: |-
  Provides a resource to create a TCR tag retention rule.
---

# tencentcloud_tcr_tag_retention_rule

Provides a resource to create a TCR tag retention rule.

## Example Usage

### Create and enable a tcr tag retention rule instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example"
  instance_type = "standard"
  delete_bucket = true
  tags = {
    "createdBy" = "Terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example"
  severity    = "medium"
}

resource "tencentcloud_tcr_tag_retention_rule" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  advanced_rule_items {
    repository_filter {
      decoration = "repoMatches"
      pattern    = "**"
    }

    retention_policy {
      key   = "nDaysSinceLastPush"
      value = 2
    }

    tag_filter {
      decoration = "matches"
      pattern    = "**"
    }
  }

  cron_setting = "daily"
}
```

## Argument Reference

The following arguments are supported:

* `cron_setting` - (Required, String, ForceNew) Execution cycle, currently only available selections are: manual; daily; weekly; monthly.
* `namespace_name` - (Required, String) The Name of the namespace.
* `registry_id` - (Required, String) The main instance ID.
* `advanced_rule_items` - (Optional, List) The advanced retention policy takes precedence; when both the basic and advanced retention policies are configured, the advanced retention policy will be used.
* `disabled` - (Optional, Bool) Whether to disable the rule, with the default value of false.
* `retention_rule` - (Optional, List) Retention Policy.

The `advanced_rule_items` object supports the following:

* `repository_filter` - (Optional, List) Warehouse filter.
* `retention_policy` - (Optional, List) Version retention rules.
* `tag_filter` - (Optional, List) Tag filter.

The `repository_filter` object of `advanced_rule_items` supports the following:

* `decoration` - (Optional, String) Filter rule types: In tag filtering, the available options are matches (match) and excludes (exclude). In repository filtering, the available options are repoMatches (repository match) and repoExcludes (repository exclude).
* `pattern` - (Optional, String) Filter expression.

The `retention_policy` object of `advanced_rule_items` supports the following:

* `key` - (Required, String) Supported strategies, with possible values: latestPushedK (retain the latest K pushed versions), nDaysSinceLastPush (retain versions pushed within the last n days).
* `value` - (Required, Int) Corresponding values under the rule settings.

The `retention_rule` object supports the following:

* `key` - (Required, String) The supported policies are latestPushedK (retain the latest `k` pushed versions) and nDaysSinceLastPush (retain pushed versions within the last `n` days).
* `value` - (Required, Int) corresponding values for rule settings.

The `tag_filter` object of `advanced_rule_items` supports the following:

* `decoration` - (Optional, String) Filter rule types: In tag filtering, the available options are matches (match) and excludes (exclude). In repository filtering, the available options are repoMatches (repository match) and repoExcludes (repository exclude).
* `pattern` - (Optional, String) Filter expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `retention_id` - The ID of the retention task.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateTagRetentionRule](https://cloud.tencent.com/document/api/1141/54803) - This API is used to create a tag retention rule.
* Create: [DescribeNamespaces](https://cloud.tencent.com/document/api/1141/42765) - This API is used to query the namespace list or the information of the specified namespace.
* Create: [DescribeTagRetentionRules](https://cloud.tencent.com/document/api/1141/54799) - This API is used to query tag retention rules.
* Read: [DescribeTagRetentionRules](https://cloud.tencent.com/document/api/1141/54799) - This API is used to query tag retention rules.
* Update: [DescribeNamespaces](https://cloud.tencent.com/document/api/1141/42765) - This API is used to query the namespace list or the information of the specified namespace.
* Update: [ModifyTagRetentionRule](https://cloud.tencent.com/document/api/1141/54798) - This API is used to update a tag retention rule.
* Delete: [DeleteTagRetentionRule](https://cloud.tencent.com/document/api/1141/54802) - This API is used to delete a tag retention rule.

## Import

TCR tag retention rule can be imported using the registryId#namespaceName#retentionId, e.g.

```
terraform import tencentcloud_tcr_tag_retention_rule.example tcr-s1jud21h#tf_example#3
```

