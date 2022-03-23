---
subcategory: "Elasticsearch"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_instance"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_instance"
description: |-
  Provides an elasticsearch instance resource.
---

# tencentcloud_elasticsearch_instance

Provides an elasticsearch instance resource.

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_instance" "foo" {
  instance_name     = "tf-test"
  availability_zone = "ap-guangzhou-3"
  version           = "7.5.1"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  password          = "Test12345"
  license_type      = "oss"

  web_node_type_info {
    node_num  = 1
    node_type = "ES.S1.MEDIUM4"
  }

  node_info_list {
    node_num  = 2
    node_type = "ES.S1.MEDIUM4"
    encrypt   = false
  }

  tags = {
    test = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `node_info_list` - (Required) Node information list, which is used to describe the specification information of various types of nodes in the cluster, such as node type, node quantity, node specification, disk type, and disk size.
* `password` - (Required) Password to an instance.
* `version` - (Required) Version of the instance. Valid values are `5.6.4`, `6.4.3`, `6.8.2` and `7.5.1`.
* `vpc_id` - (Required, ForceNew) The ID of a VPC network.
* `availability_zone` - (Optional, ForceNew) Availability zone. When create multi-az es, this parameter must be omitted.
* `basic_security_type` - (Optional) Whether to enable X-Pack security authentication in Basic Edition 6.8 and above. Valid values are `1` and `2`. `1` is disabled, `2` is enabled, and default value is `1`.
* `charge_period` - (Optional, ForceNew) The tenancy of the prepaid instance, and uint is month. NOTE: it only works when charge_type is set to `PREPAID`.
* `charge_type` - (Optional, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`.
* `deploy_mode` - (Optional, ForceNew) Cluster deployment mode. Valid values are `0` and `1`. `0` is single-AZ deployment, and `1` is multi-AZ deployment. Default value is `0`.
* `instance_name` - (Optional) Name of the instance, which can contain 1 to 50 English letters, Chinese characters, digits, dashes(-), or underscores(_).
* `license_type` - (Optional) License type. Valid values are `oss`, `basic` and `platinum`. The default value is `platinum`.
* `multi_zone_infos` - (Optional, ForceNew) Details of AZs in multi-AZ deployment mode (which is required when deploy_mode is `1`).
* `renew_flag` - (Optional, ForceNew) When enabled, the instance will be renew automatically when it reach the end of the prepaid tenancy. Valid values are `RENEW_FLAG_AUTO` and `RENEW_FLAG_MANUAL`. NOTE: it only works when charge_type is set to `PREPAID`.
* `subnet_id` - (Optional, ForceNew) The ID of a VPC subnetwork. When create multi-az es, this parameter must be omitted.
* `tags` - (Optional) A mapping of tags to assign to the instance. For tag limits, please refer to [Use Limits](https://intl.cloud.tencent.com/document/product/651/13354).
* `web_node_type_info` - (Optional) Visual node configuration.

The `multi_zone_infos` object supports the following:

* `availability_zone` - (Required) Availability zone.
* `subnet_id` - (Required) The ID of a VPC subnetwork.

The `node_info_list` object supports the following:

* `node_num` - (Required) Number of nodes.
* `node_type` - (Required) Node specification, and valid values refer to [document of tencentcloud](https://intl.cloud.tencent.com/document/product/845/18376).
* `disk_size` - (Optional) Node disk size. Unit is GB, and default value is `100`.
* `disk_type` - (Optional) Node disk type. Valid values are `CLOUD_SSD` and `CLOUD_PREMIUM`. The default value is `CLOUD_SSD`.
* `encrypt` - (Optional) Decides to encrypt this disk or not.
* `type` - (Optional) Node type. Valid values are `hotData`, `warmData` and `dedicatedMaster`. The default value is 'hotData`.

The `web_node_type_info` object supports the following:

* `node_num` - (Required) Visual node number.
* `node_type` - (Required) Visual node specifications.

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

