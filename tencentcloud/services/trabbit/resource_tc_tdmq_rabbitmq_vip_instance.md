Provides a resource to create a TDMQ rabbitmq vip instance

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
  enable_deletion_protection            = true
  remark                                = "Example RabbitMQ instance"
  enable_risk_warning                   = false
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

Update instance configuration

The following parameters can be updated after instance creation:
- `cluster_name` - Cluster name
- `resource_tags` - Resource tags
- `auto_renew_flag` - Auto renew flag (prepaid instances only)
- `enable_public_access` - Enable public network access
- `band_width` - Public network bandwidth (when public access is enabled)

Example of updating instance configuration:

```hcl
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

  # Initially disable public access
  enable_public_access                  = false
  band_width                            = 100

  resource_tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}

# Later, update to enable public access with higher bandwidth
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = tencentcloud_vpc.vpc.id
  subnet_id                             = tencentcloud_subnet.subnet.id
  cluster_name                          = "tf-example-rabbitmq-vip-instance-updated"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = false
  time_span                             = 1

  # Enable public access with updated bandwidth
  enable_public_access                  = true
  band_width                            = 200

  resource_tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

**Note**: The following parameters are immutable and cannot be changed after instance creation. If you need to modify these parameters, you must recreate the instance:
- `zone_ids` - Availability zones
- `vpc_id` - VPC ID
- `subnet_id` - Subnet ID
- `node_spec` - Node specification
- `node_num` - Number of nodes
- `storage_size` - Storage size
- `enable_create_default_ha_mirror_queue` - HA mirror queue flag
- `time_span` - Purchase duration (prepaid instances)
- `pay_mode` - Payment mode
- `cluster_version` - Cluster version

Import

TDMQ rabbitmq vip instance can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rabbitmq_vip_instance.example amqp-mok52gmn
```

Update Behavior

The following parameters can be updated after instance creation:
- `cluster_name` - Instance cluster name
- `enable_deletion_protection` - Whether to enable deletion protection
- `remark` - Instance remark/description
- `enable_risk_warning` - Whether to enable cluster risk warning
- `resource_tags` - Instance resource tags

The following parameters are immutable and cannot be changed after creation:
- `zone_ids` - Availability zones
- `vpc_id` - VPC ID
- `subnet_id` - Subnet ID
- `node_spec` - Node specifications
- `node_num` - Number of nodes
- `storage_size` - Storage size
- `enable_create_default_ha_mirror_queue` - Mirror queue setting
- `auto_renew_flag` - Auto renew flag
- `time_span` - Purchase duration
- `pay_mode` - Payment method
- `cluster_version` - Cluster version
- `band_width` - Public network bandwidth
- `enable_public_access` - Whether to enable public network access

