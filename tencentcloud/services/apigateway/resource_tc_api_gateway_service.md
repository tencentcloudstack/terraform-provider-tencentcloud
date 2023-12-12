Use this resource to create API gateway service.

~> **NOTE:** After setting `uniq_vpc_id`, it cannot be modified.

Example Usage

Shared Service

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "example-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf-example"
  protocol     = "http&https"
  service_desc = "desc."
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
  uniq_vpc_id  = tencentcloud_vpc.vpc.id

  tags = {
    createdBy = "terraform"
  }

  release_limit = 500
  pre_limit     = 500
  test_limit    = 500
}
```

Exclusive Service

```hcl
resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf-example"
  protocol     = "http&https"
  service_desc = "desc."
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
  uniq_vpc_id  = tencentcloud_vpc.vpc.id
  instance_id  = "instance-rc6fcv4e"

  tags = {
    createdBy = "terraform"
  }

  release_limit = 500
  pre_limit     = 500
  test_limit    = 500
}
```

Import

API gateway service can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_service.service service-pg6ud8pa
```