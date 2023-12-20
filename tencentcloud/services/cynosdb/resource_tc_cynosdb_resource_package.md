Provides a resource to create a cynosdb resource_package

Example Usage

```hcl
resource "tencentcloud_cynosdb_resource_package" "resource_package" {
  instance_type = "cdb"
  package_region = "china"
  package_type = "CCU"
  package_version = "base"
  package_spec =
  expire_day = 180
  package_count = 1
  package_name = "PackageName"
}
```

Import

cynosdb resource_package can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_resource_package.resource_package resource_package_id
```
