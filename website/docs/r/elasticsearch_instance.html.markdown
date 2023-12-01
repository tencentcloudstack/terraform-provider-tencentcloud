---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_instance"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_instance"
description: |-
  Provides an elasticsearch instance resource.
---

# tencentcloud_elasticsearch_instance

Provides an elasticsearch instance resource.

## Example Usage

### Create a basic version of elasticsearch instance paid by the hour

```hcl
data "tencentcloud_availability_zones_by_product" "availability_zone" {
  product = "es"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_es_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_availability_zones_by_product.availability_zone.zones.0.name
  name              = "tf_es_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_elasticsearch_instance" "example" {
  instance_name       = "tf_example_es"
  availability_zone   = data.tencentcloud_availability_zones_by_product.availability_zone.zones.0.name
  version             = "7.10.1"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  password            = "Test12345"
  license_type        = "basic"
  basic_security_type = 2

  web_node_type_info {
    node_num  = 1
    node_type = "ES.S1.MEDIUM4"
  }

  node_info_list {
    node_num  = 2
    node_type = "ES.S1.MEDIUM8"
    encrypt   = false
  }

  es_acl {
    # black_list = [
    #   "9.9.9.9",
    #   "8.8.8.8",
    # ]
    white_list = [
      "127.0.0.1",
    ]
  }

  tags = {
    test = "test"
  }
}
```

### Create a basic version of elasticsearch instance for multi-availability zone deployment

```hcl
data "tencentcloud_availability_zones_by_product" "availability_zone" {
  product = "es"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_es_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_availability_zones_by_product.availability_zone.zones.0.name
  name              = "tf_es_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_subnet" "subnet_multi_zone" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_availability_zones_by_product.availability_zone.zones.1.name
  name              = "tf_es_subnet"
  cidr_block        = "10.0.2.0/24"
}

resource "tencentcloud_elasticsearch_instance" "example_multi_zone" {
  instance_name       = "tf_example_es"
  availability_zone   = "-"
  version             = "7.10.1"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = "-"
  password            = "Test12345"
  license_type        = "basic"
  basic_security_type = 2
  deploy_mode         = 1

  multi_zone_infos {
    availability_zone = data.tencentcloud_availability_zones_by_product.availability_zone.zones.0.name
    subnet_id         = tencentcloud_subnet.subnet.id
  }

  multi_zone_infos {
    availability_zone = data.tencentcloud_availability_zones_by_product.availability_zone.zones.1.name
    subnet_id         = tencentcloud_subnet.subnet_multi_zone.id
  }

  web_node_type_info {
    node_num  = 1
    node_type = "ES.S1.MEDIUM4"
  }

  node_info_list {
    type      = "dedicatedMaster"
    node_num  = 3
    node_type = "ES.S1.MEDIUM8"
    encrypt   = false
  }

  node_info_list {
    type      = "hotData"
    node_num  = 2
    node_type = "ES.S1.MEDIUM8"
    encrypt   = false
  }

  es_acl {
    # black_list = [
    #   "9.9.9.9",
    #   "8.8.8.8",
    # ]
    white_list = [
      "127.0.0.1",
    ]
  }

  tags = {
    test = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `node_info_list` - (Required, List) Node information list, which is used to describe the specification information of various types of nodes in the cluster, such as node type, node quantity, node specification, disk type, and disk size.
* `password` - (Required, String) Password to an instance, the password needs to be 8 to 16 characters, including at least two items ([a-z,A-Z], [0-9] and [-!@#$%&^*+=_:;,.?] special symbols.
* `version` - (Required, String) Version of the instance. Valid values are `5.6.4`, `6.4.3`, `6.8.2`, `7.5.1` and `7.10.1`.
* `vpc_id` - (Required, String, ForceNew) The ID of a VPC network.
* `availability_zone` - (Optional, String, ForceNew) Availability zone. When create multi-az es, this parameter must be omitted or `-`.
* `basic_security_type` - (Optional, Int) Whether to enable X-Pack security authentication in Basic Edition 6.8 and above. Valid values are `1` and `2`. `1` is disabled, `2` is enabled, and default value is `1`. Notice: this parameter is only take effect on `basic` license.
* `charge_period` - (Optional, Int, ForceNew) The tenancy of the prepaid instance, and uint is month. NOTE: it only works when charge_type is set to `PREPAID`.
* `charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`.
* `deploy_mode` - (Optional, Int, ForceNew) Cluster deployment mode. Valid values are `0` and `1`. `0` is single-AZ deployment, and `1` is multi-AZ deployment. Default value is `0`.
* `es_acl` - (Optional, List) Kibana Access Control Configuration.
* `instance_name` - (Optional, String) Name of the instance, which can contain 1 to 50 English letters, Chinese characters, digits, dashes(-), or underscores(_).
* `license_type` - (Optional, String) License type. Valid values are `oss`, `basic` and `platinum`. The default value is `platinum`.
* `multi_zone_infos` - (Optional, List, ForceNew) Details of AZs in multi-AZ deployment mode (which is required when deploy_mode is `1`).
* `renew_flag` - (Optional, String, ForceNew) When enabled, the instance will be renew automatically when it reach the end of the prepaid tenancy. Valid values are `RENEW_FLAG_AUTO` and `RENEW_FLAG_MANUAL`. NOTE: it only works when charge_type is set to `PREPAID`.
* `subnet_id` - (Optional, String, ForceNew) The ID of a VPC subnetwork. When create multi-az es, this parameter must be omitted or `-`.
* `tags` - (Optional, Map) A mapping of tags to assign to the instance. For tag limits, please refer to [Use Limits](https://intl.cloud.tencent.com/document/product/651/13354).
* `web_node_type_info` - (Optional, List) Visual node configuration.

The `es_acl` object supports the following:

* `black_list` - (Optional, Set) Blacklist of kibana access.
* `white_list` - (Optional, Set) Whitelist of kibana access.

The `multi_zone_infos` object supports the following:

* `availability_zone` - (Required, String) Availability zone.
* `subnet_id` - (Required, String) The ID of a VPC subnetwork.

The `node_info_list` object supports the following:

* `node_num` - (Required, Int) Number of nodes.
* `node_type` - (Required, String) Node specification, and valid values refer to [document of tencentcloud](https://intl.cloud.tencent.com/document/product/845/18376).
* `disk_size` - (Optional, Int) Node disk size. Unit is GB, and default value is `100`.
* `disk_type` - (Optional, String) Node disk type. Valid values are `CLOUD_SSD` and `CLOUD_PREMIUM`, `CLOUD_HSSD`. The default value is `CLOUD_SSD`.
* `encrypt` - (Optional, Bool) Decides to encrypt this disk or not.
* `type` - (Optional, String) Node type. Valid values are `hotData`, `warmData` and `dedicatedMaster`. The default value is 'hotData`.

The `web_node_type_info` object supports the following:

* `node_num` - (Required, Int) Visual node number.
* `node_type` - (Required, String) Visual node specifications.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Instance creation time.
* `elasticsearch_domain` - Elasticsearch domain name.
* `elasticsearch_port` - Elasticsearch port.
* `elasticsearch_vip` - Elasticsearch VIP.
* `kibana_url` - Kibana access URL.


## Import

Elasticsearch instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_elasticsearch_instance.foo es-17634f05
```

