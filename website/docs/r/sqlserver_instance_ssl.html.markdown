---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_instance_ssl"
sidebar_current: "docs-tencentcloud-resource-sqlserver_instance_ssl"
description: |-
  Provides a resource to create a sqlserver instance_ssl
---

# tencentcloud_sqlserver_instance_ssl

Provides a resource to create a sqlserver instance_ssl

## Example Usage

```hcl
resource "tencentcloud_sqlserver_instance_ssl" "example" {
  instance_id = "mssql-qelbzgwf"
  type        = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `type` - (Required, String) Operation type. enable: turn on SSL; disable: turn off SSL; renew: update the certificate validity period.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeDBInstancesAttribute](https://www.tencentcloud.com/document/api/238/54811)
  - [ModifyDBInstanceSSL](https://www.tencentcloud.com/document/api/238/72764)
* Read: [DescribeDBInstancesAttribute](https://www.tencentcloud.com/document/api/238/54811)
* Update:
  - [DescribeDBInstancesAttribute](https://www.tencentcloud.com/document/api/238/54811)
  - [ModifyDBInstanceSSL](https://www.tencentcloud.com/document/api/238/72764)

