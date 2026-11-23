---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_resource_package"
sidebar_current: "docs-tencentcloud-resource-cynosdb_resource_package"
description: |-
  Provides a resource to create a cynosdb resource_package
---

# tencentcloud_cynosdb_resource_package

Provides a resource to create a cynosdb resource_package

## Example Usage

```hcl
resource "tencentcloud_cynosdb_resource_package" "resource_package" {
  instance_type   = "cdb"
  package_region  = "china"
  package_type    = "CCU"
  package_version = "base"
  package_spec    = 1
  expire_day      = 180
  package_count   = 1
  package_name    = "PackageName"
}
```

## Argument Reference

The following arguments are supported:

* `expire_day` - (Required, Int) Validity period of resource package, in days.
* `instance_type` - (Required, String) Instance Type.
* `package_count` - (Required, Int) Number of purchased resource packs.
* `package_region` - (Required, String) Resource package usage region China - common in mainland China, overseas - common in Hong Kong, Macao, Taiwan, and overseas.
* `package_spec` - (Required, Float64) Resource package size, calculated in 10000 units; Storage resources: GB.
* `package_type` - (Required, String) Resource package type: CCU computing resource package, DISK storage resource package.
* `package_version` - (Required, String) Resource package version base basic version, common general version, enterprise enterprise version.
* `package_name` - (Optional, String) Resource Package Name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateResourcePackage](https://www.tencentcloud.com/document/api/1098/55899)
* Read: [DescribeResourcePackageDetail](https://www.tencentcloud.com/document/api/1098/55898)
* Update: [ModifyResourcePackageName](https://www.tencentcloud.com/document/api/1098/55894)
* Delete: [RefundResourcePackage](https://www.tencentcloud.com/document/api/1098/55893)

## Import

cynosdb resource_package can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_resource_package.resource_package resource_package_id
```

