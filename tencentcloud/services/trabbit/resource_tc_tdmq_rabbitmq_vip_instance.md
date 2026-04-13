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
  remark                                = "test-remark"
  enable_deletion_protection            = true
  enable_risk_warning                    = true
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

Update instance attributes

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  # ... other fields ...

  # Update instance remark
  remark = "updated-remark"

  # Toggle deletion protection
  enable_deletion_protection = false

  # Toggle risk warning
  enable_risk_warning = false
}
```

Arguments Reference

The following arguments are supported:

* `cluster_name` - (Required, ForceNew) Cluster name.
* `cluster_version` - (Optional) Cluster version, the default is `3.8.30`, valid values: `3.8.30`, `3.11.8` and `3.13.7`.
* `enable_create_default_ha_mirror_queue` - (Optional) Mirrored queue, the default is false.
* `enable_public_access` - (Optional) Whether to enable public network access. Default is false.
* `band_width` - (Optional) Public network bandwidth in Mbps.
* `node_num` - (Optional) The number of nodes, a minimum of 3 nodes for a multi-availability zone. If not passed, the default single availability zone is 1, and the multi-availability zone is 3.
* `node_spec` - (Optional, ForceNew) Node specifications. Valid values: rabbit-vip-basic-5 (for 2C4G), rabbit-vip-profession-2c8g (for 2C8G), rabbit-vip-basic-1 (for 4C8G), rabbit-vip-profession-4c16g (for 4C16G), rabbit-vip-basic-2 (for 8C16G), rabbit-vip-profession-8c32g (for 8C32G), rabbit-vip-basic-4 (for 16C32G), rabbit-vip-profession-16c64g (for 16C64G). The default is rabbit-vip-basic-1. NOTE: The above specifications may be sold out or removed from the shelves.
* `pay_mode` - (Optional) Payment method: 0 indicates postpaid; 1 indicates prepaid. Default: prepaid.
* `subnet_id` - (Required, ForceNew) Private network SubnetId.
* `time_span` - (Optional) Purchase duration, the default is 1 (month).
* `vpc_id` - (Required, ForceNew) Private network VpcId.
* `zone_ids` - (Required, ForceNew) availability zone.
* `auto_renew_flag` - (Optional) Automatic renewal, the default is true.
* `storage_size` - (Optional) Single node storage specification, the default is 200G.
* `resource_tags` - (Optional) Instance resource tags. Each tag is a key-value pair for resource identification and management.
* `remark` - (Optional) Instance remark.
* `enable_deletion_protection` - (Optional) Whether to enable deletion protection. Default is false.
* `enable_risk_warning` - (Optional) Whether to enable cluster risk warning. Default is false.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the instance.
* `public_access_endpoint` - Public Network Access Point.
* `vpcs` - List of VPC Access Points.

**Fields that cannot be updated:**

* `zone_ids`
* `vpc_id`
* `subnet_id`
* `node_spec`
* `node_num`
* `storage_size`
* `enable_create_default_ha_mirror_queue`
* `auto_renew_flag`
* `time_span`
* `pay_mode`
* `cluster_version`
* `band_width`
* `enable_public_access`

Import

TDMQ rabbitmq vip instance can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rabbitmq_vip_instance.example amqp-mok52gmn
```

