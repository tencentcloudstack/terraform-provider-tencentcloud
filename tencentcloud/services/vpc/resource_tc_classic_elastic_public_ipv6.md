Provides a resource to create a vpc classic_elastic_public_ipv6

Example Usage

```hcl
resource "tencentcloud_classic_elastic_public_ipv6" "classic_elastic_public_ipv6" {
  ip6_address                = "xxxxxx"
  internet_max_bandwidth_out = 2
  tags = {
    "testkey" = "testvalue"
  }
}
```

Import

vpc classic_elastic_public_ipv6 can be imported using the id, e.g.

```
terraform import tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6 classic_elastic_public_ipv6_id
```
