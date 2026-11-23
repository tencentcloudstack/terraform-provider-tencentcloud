---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_encrypt_attributes"
sidebar_current: "docs-tencentcloud-resource-mariadb_encrypt_attributes"
description: |-
  Provides a resource to create a mariadb encrypt_attributes
---

# tencentcloud_mariadb_encrypt_attributes

Provides a resource to create a mariadb encrypt_attributes

## Example Usage

```hcl
resource "tencentcloud_mariadb_encrypt_attributes" "encrypt_attributes" {
  instance_id     = "tdsql-ow728lmc"
  encrypt_enabled = 1
}
```

## Argument Reference

The following arguments are supported:

* `encrypt_enabled` - (Required, Int) whether to enable data encryption, it is not supported to turn it off after it is turned on. The optional values: 0-disable, 1-enable.
* `instance_id` - (Required, String) instance id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyDBEncryptAttributes](https://www.tencentcloud.com/document/api/237/53947)
* Read: [DescribeDBEncryptAttributes](https://www.tencentcloud.com/document/api/237/53462)
* Update: [ModifyDBEncryptAttributes](https://www.tencentcloud.com/document/api/237/53947)

