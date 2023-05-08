---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_immutable_tag_rule"
sidebar_current: "docs-tencentcloud-resource-tcr_immutable_tag_rule"
description: |-
  Provides a resource to create a tcr immutable_tag_rule
---

# tencentcloud_tcr_immutable_tag_rule

Provides a resource to create a tcr immutable_tag_rule

## Example Usage

```hcl
resource "tencentcloud_tcr_immutable_tag_rule" "my_rule" {
  registry_id    = "%s"
  namespace_name = "%s"
  rule {
    repository_pattern    = "**"
    tag_pattern           = "**"
    repository_decoration = "repoMatches"
    tag_decoration        = "matches"
    disabled              = false
  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `namespace_name` - (Required, String) namespace name.
* `registry_id` - (Required, String) instance id.
* `rule` - (Required, List) rule.
* `tags` - (Optional, Map) Tag description list.

The `rule` object supports the following:

* `repository_decoration` - (Required, String) repository decoration type:repoMatches or repoExcludes.
* `repository_pattern` - (Required, String) repository matching rules.
* `tag_decoration` - (Required, String) tag decoration type: matches or excludes.
* `tag_pattern` - (Required, String) tag matching rules.
* `disabled` - (Optional, Bool) disable rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tcr immutable_tag_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_immutable_tag_rule.immutable_tag_rule immutable_tag_rule_id
```

