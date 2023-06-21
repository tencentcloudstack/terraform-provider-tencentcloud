---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_account"
sidebar_current: "docs-tencentcloud-resource-cynosdb_account"
description: |-
  Provides a resource to create a cynosdb account
---

# tencentcloud_cynosdb_account

Provides a resource to create a cynosdb account

## Example Usage

```hcl
resource "tencentcloud_cynosdb_account" "account" {
  cluster_id = "xxx"
  accounts {
    account_name         = ""
    account_password     = ""
    host                 = ""
    description          = ""
    max_user_connections =
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, String) Account name, including alphanumeric _, Start with a letter, end with a letter or number, length 1-16.
* `account_password` - (Required, String) Password, with a length range of 8 to 64 characters.
* `cluster_id` - (Required, String) Cluster ID.
* `host` - (Required, String) main engine.
* `description` - (Optional, String) describe.
* `max_user_connections` - (Optional, Int) The maximum number of user connections cannot be greater than 10240.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb account can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_account.account account_id
```

