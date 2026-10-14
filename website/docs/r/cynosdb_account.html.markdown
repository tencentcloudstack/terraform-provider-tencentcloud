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

### If host is %

```hcl
resource "tencentcloud_cynosdb_account" "example" {
  cluster_id           = "cynosdbmysql-ddciqx2l"
  account_name         = "tf_example"
  account_password     = "Password@123"
  host                 = "%"
  description          = "remark."
  max_user_connections = 10
}
```

### If host is ip

```hcl
resource "tencentcloud_cynosdb_account" "example" {
  cluster_id           = "cynosdbmysql-ddciqx2l"
  account_name         = "tf_example"
  account_password     = "Password@123"
  host                 = "1.1.1.1"
  description          = "remark."
  max_user_connections = 0
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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAccounts](https://cloud.tencent.com/document/api/1003/71660) - This API is used to create user accounts.
* Create: [DescribeTasks](https://cloud.tencent.com/document/api/1003/103866) - This API is used to query task lists.
* Read: [DescribeAccounts](https://cloud.tencent.com/document/api/1003/48075) - This API is used to query the database account list.
* Update: [ModifyAccountDescription](https://cloud.tencent.com/document/api/1003/92716) - This API is used to modify the descriptions of a database account.
* Update: [ModifyAccountHost](https://cloud.tencent.com/document/api/1003/92715) - This API is used to modify account hosts.
* Update: [ModifyAccountParams](https://cloud.tencent.com/document/api/1003/70112)
* Update: [ResetAccountPassword](https://cloud.tencent.com/document/api/1003/83592) - This API is used to modify the database account password.
* Delete: [DeleteAccounts](https://cloud.tencent.com/document/api/1003/92718) - This API is used to delete user accounts.

## Import

cynosdb account can be imported using the clusterId#accountName#host, e.g.

```
terraform import tencentcloud_cynosdb_account.example cynosdbmysql-ddciqx2l#tf_example#%

or

terraform import tencentcloud_cynosdb_account.example cynosdbmysql-ddciqx2l#tf_example#1.1.1.1
```

