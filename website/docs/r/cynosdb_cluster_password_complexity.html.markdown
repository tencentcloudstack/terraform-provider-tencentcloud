---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_password_complexity"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cluster_password_complexity"
description: |-
  Provides a resource to create a cynosdb cluster_password_complexity
---

# tencentcloud_cynosdb_cluster_password_complexity

Provides a resource to create a cynosdb cluster_password_complexity

## Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_password_complexity" "cluster_password_complexity" {
  cluster_id                           = "cynosdbmysql-cgd2gpwr"
  validate_password_length             = 8
  validate_password_mixed_case_count   = 1
  validate_password_special_char_count = 1
  validate_password_number_count       = 1
  validate_password_policy             = "STRONG"
  validate_password_dictionary = [
    "cccc",
    "xxxx",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `validate_password_length` - (Required, Int) Password length.
* `validate_password_mixed_case_count` - (Required, Int) Number of uppercase and lowercase characters.
* `validate_password_number_count` - (Required, Int) Number of digits.
* `validate_password_policy` - (Required, String) Password strength (MEDIUM, STRONG).
* `validate_password_special_char_count` - (Required, Int) Number of special characters.
* `validate_password_dictionary` - (Optional, Set: [`String`]) Data dictionary.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb cluster_password_complexity can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity cluster_password_complexity_id
```

