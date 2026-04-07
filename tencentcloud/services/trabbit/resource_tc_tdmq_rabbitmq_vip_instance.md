Provides a resource to create a TDMQ rabbitmq vip instance

## Update Policy

The following fields can be updated after instance creation:
- `cluster_name` - Instance name
- `resource_tags` - Instance tags

The following fields currently cannot be updated due to Tencent Cloud API limitations:
- `node_spec` - Node specification (contact Tencent Cloud support for spec changes)
- `node_num` - Node count (contact Tencent Cloud support for scaling)
- `storage_size` - Storage size (contact Tencent Cloud support for storage expansion)
- `band_width` - Public network bandwidth (contact Tencent Cloud support for bandwidth adjustments)
- `enable_public_access` - Public network access (contact Tencent Cloud support for access changes)

The following fields are immutable and cannot be changed after instance creation:
- `zone_ids` - Availability zone(s)
- `vpc_id` - VPC ID
- `subnet_id` - Subnet ID
- `enable_create_default_ha_mirror_queue` - Mirror queue configuration
- `auto_renew_flag` - Auto-renewal flag
- `time_span` - Subscription period
- `pay_mode` - Payment mode
- `cluster_version` - Cluster version

Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {
  name = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rabbitmq instance
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = tencentcloud_vpc.vpc.id
  subnet_id                             = tencentcloud_subnet.subnet.id
  cluster_name                          = "tf-example-rabbitmq-vip-instance"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}

# create postpaid rabbitmq instance
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example2" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = tencentcloud_vpc.vpc.id
  subnet_id                             = tencentcloud_subnet.subnet.id
  cluster_name                          = "tf-example-rabbitmq-vip-instance"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
  pay_mode                              = 0
  cluster_version                       = "3.11.8"
  resource_tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

Enable public network access

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [100006]
  vpc_id                                = "vpc-i5yyodl9"
  subnet_id                             = "subnet-hhi88a58"
  cluster_name                          = "tf-example"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  pay_mode                              = 0
  cluster_version                       = "3.11.8"
  enable_public_access                  = true
  band_width                            = 100

  resource_tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

Import

TDMQ rabbitmq vip instance can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rabbitmq_vip_instance.example amqp-mok52gmn
```

