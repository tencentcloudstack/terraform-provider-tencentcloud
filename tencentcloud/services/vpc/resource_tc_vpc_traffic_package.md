Provides a resource to create a vpc traffic_package

Example Usage

```hcl
resource "tencentcloud_vpc_traffic_package" "example" {
  traffic_amount = 10
}
```

Import

vpc traffic_package can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_traffic_package.traffic_package traffic_package_id
```