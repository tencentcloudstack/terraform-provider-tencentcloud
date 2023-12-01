---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_immutable_tag_rule"
sidebar_current: "docs-tencentcloud-resource-tcr_immutable_tag_rule"
description: |-
  Provides a resource to create a tcr immutable tag rule.
---

# tencentcloud_tcr_immutable_tag_rule

Provides a resource to create a tcr immutable tag rule.

## Example Usage

### Create a immutable tag rule with specified tags and exclude specified repositories

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "premium"
  delete_bucket = true
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_example_ns"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_immutable_tag_rule" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  rule {
    repository_pattern    = "deprecated_repo" # specify exclude repo
    tag_pattern           = "**"              # all tags
    repository_decoration = "repoExcludes"
    tag_decoration        = "matches"
    disabled              = false
  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

### With specified repositories and exclude specified version tag

```hcl
resource "tencentcloud_tcr_immutable_tag_rule" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  rule {
    repository_pattern    = "**" # all repo
    tag_pattern           = "v1" # exlude v1 tags
    repository_decoration = "repoMatches"
    tag_decoration        = "excludes"
    disabled              = false
  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

### Disabled the specified rule

```hcl
resource "tencentcloud_tcr_immutable_tag_rule" "example_rule_A" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  rule {
    repository_pattern    = "deprecated_repo" # specify exclude repo
    tag_pattern           = "**"              # all tags
    repository_decoration = "repoExcludes"
    tag_decoration        = "matches"
    disabled              = false
  }
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_immutable_tag_rule" "example_rule_B" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  rule {
    repository_pattern    = "**" # all repo
    tag_pattern           = "v1" # exlude v1 tags
    repository_decoration = "repoMatches"
    tag_decoration        = "excludes"
    disabled              = true # disable it
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

