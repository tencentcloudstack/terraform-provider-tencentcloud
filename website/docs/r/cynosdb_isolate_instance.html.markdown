---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_isolate_instance"
sidebar_current: "docs-tencentcloud-resource-cynosdb_isolate_instance"
description: |-
  Provides a resource to create a cynosdb isolate_instance
---

# tencentcloud_cynosdb_isolate_instance

Provides a resource to create a cynosdb isolate_instance

## Example Usage

```hcl
resource "tencentcloud_cynosdb_account" "account" {
  cluster_id           = "cynosdbmysql-bws8h88b"
  account_name         = "terraform_test"
  account_password     = "Password@1234"
  host                 = "%"
  description          = "testx"
  max_user_connections = 2
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `operate` - (Required, String) isolate, activate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ActivateInstance](https://www.tencentcloud.com/document/api/1098/45568)
* Create: [DescribeFlow](https://www.tencentcloud.com/document/api/1098/53463)
* Create: [IsolateInstance](https://www.tencentcloud.com/document/api/1098/40744)
* Update: [ActivateInstance](https://www.tencentcloud.com/document/api/1098/45568)
* Update: [DescribeFlow](https://www.tencentcloud.com/document/api/1098/53463)
* Update: [IsolateInstance](https://www.tencentcloud.com/document/api/1098/40744)

