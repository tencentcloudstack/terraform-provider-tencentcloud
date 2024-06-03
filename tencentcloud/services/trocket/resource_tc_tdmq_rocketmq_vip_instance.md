Provides a resource to create a tdmq rocketmq_vip_instance

~> **NOTE:** The instance cannot be downgraded, Include parameters `node_count`, `spec`, `storage_size`.
~> **NOTE:** If `spec` is `rocket-vip-basic-2`, configuration changes are not supported.

Example Usage

```hcl
# query availability zones
data "tencentcloud_availability_zones" "zones" {}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.1.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

# create rocketmq vip instance
resource "tencentcloud_tdmq_rocketmq_vip_instance" "example" {
  name         = "tx-example"
  spec         = "rocket-vip-basic-2"
  node_count   = 2
  storage_size = 200
  zone_ids     = [
    data.tencentcloud_availability_zones.zones.zones.0.id,
    data.tencentcloud_availability_zones.zones.zones.1.id
  ]

  vpc_info {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }

  time_span = 1
  ip_rules {
    ip_rule = "0.0.0.0/0"
    allow   = true
    remark  = "remark."
  }
}
```