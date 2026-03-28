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

## Updating fields

The following fields can be updated after instance creation:
- `node_spec`: Node specification (e.g., upgrade from rabbit-vip-basic-1 to rabbit-vip-profession-4c16g)
- `node_num`: Number of nodes (e.g., scale from 1 to 3 nodes)
- `storage_size`: Storage size in GB (e.g., expand from 200 to 500)
- `auto_renew_flag`: Auto-renewal flag (true/false)
- `band_width`: Public network bandwidth in Mbps
- `enable_public_access`: Toggle public network access on/off

Example - Update node specification and storage size:

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [100006]
  vpc_id                                = "vpc-i5yyodl9"
  subnet_id                             = "subnet-hhi88a58"
  cluster_name                          = "tf-example"
  node_spec                             = "rabbit-vip-profession-4c16g"  # Updated from rabbit-vip-basic-1
  node_num                              = 3                              # Updated from 1
  storage_size                          = 500                            # Updated from 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}
```

Example - Toggle public network access:

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
  enable_public_access                  = true   # Updated from false
  band_width                            = 200    # Updated from 100
}
```

Note: The following fields cannot be changed after creation:
- `zone_ids`: Availability zones
- `vpc_id`: VPC ID
- `subnet_id`: Subnet ID
- `cluster_version`: Cluster version (requires special upgrade API)
- `enable_create_default_ha_mirror_queue`: HA mirror queue setting
- `time_span`: Purchase duration
- `pay_mode`: Payment method

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

