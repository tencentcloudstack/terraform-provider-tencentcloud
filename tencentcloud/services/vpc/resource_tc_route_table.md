Provides a resource to create a VPC routing table.

Example Usage

```hcl
resource "tencentcloud_vpc" "foo" {
  name       = "ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "foo" {
  vpc_id = tencentcloud_vpc.foo.id
  name   = "ci-temp-test-rt"
}
```

Import

Vpc routetable instance can be imported, e.g.

```
$ terraform import tencentcloud_route_table.test route_table_id
```