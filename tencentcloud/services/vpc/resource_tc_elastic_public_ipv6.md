Provides a resource to create a vpc elastic_public_ipv6

Example Usage

```hcl
resource "tencentcloud_elastic_public_ipv6" "elastic_public_ipv6" {
    address_name = "test"
    internet_max_bandwidth_out = 1
    tags = {
        "test1key" = "test1value"
    }
}
```

Import

vpc elastic_public_ipv6 can be imported using the id, e.g.

```
terraform import tencentcloud_elastic_public_ipv6.elastic_public_ipv6 elastic_public_ipv6_id
```
