---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_repository"
sidebar_current: "docs-tencentcloud-resource-tcr_repository"
description: |-
  Use this resource to create tcr repository.
---

# tencentcloud_tcr_repository

Use this resource to create tcr repository.

## Example Usage

### Create a tcr repository instance

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

resource "tencentcloud_tcr_repository" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  name           = "test"
  brief_desc     = "111"
  description    = "111111111111111111111111111111111111"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the TCR instance.
* `name` - (Required, String, ForceNew) Name of the TCR repository. Valid length is [2~200]. It can only contain lowercase letters, numbers and separators (`.`, `_`, `-`, `/`), and cannot start, end or continue with separators. Support the use of multi-level address formats, such as `sub1/sub2/repo`.
* `namespace_name` - (Required, String, ForceNew) Name of the TCR namespace.
* `brief_desc` - (Optional, String) Brief description of the repository. Valid length is [1~100].
* `description` - (Optional, String) Description of the repository. Valid length is [1~1000].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `is_public` - Indicate the repository is public or not.
* `update_time` - Last updated time.
* `url` - URL of the repository.


## Import

tcr repository can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_repository.foo instance_id#namespace_name#repository_name
```

