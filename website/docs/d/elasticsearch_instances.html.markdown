---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_instances"
sidebar_current: "docs-tencentcloud-datasource-elasticsearch_instances"
description: |-
  Use this data source to query elasticsearch instances.
---

# tencentcloud_elasticsearch_instances

Use this data source to query elasticsearch instances.

## Example Usage

```hcl
data "tencentcloud_elasticsearch_instances" "foo" {
  instance_id = "es-17634f05"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional) ID of the instance to be queried.
* `instance_name` - (Optional) Name of the instance to be queried.
* `result_output_file` - (Optional) Used to save results.
* `tags` - (Optional) Tag of the instance to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - An information list of elasticsearch instance. Each element contains the following attributes:
  * `availability_zone` - Availability zone.
  * `basic_security_type` - Whether to enable X-Pack security authentication in Basic Edition 6.8 and above.
  * `charge_type` - The charge type of instance.
  * `create_time` - Instance creation time.
  * `deploy_mode` - Cluster deployment mode.
  * `elasticsearch_domain` - Elasticsearch domain name.
  * `elasticsearch_port` - Elasticsearch port.
  * `elasticsearch_vip` - Elasticsearch VIP.
  * `instance_id` - ID of the instance.
  * `instance_name` - Name of the instance.
  * `kibana_url` - Kibana access URL.
  * `license_type` - License type.
  * `multi_zone_infos` - Details of AZs in multi-AZ deployment mode.
    * `availability_zone` - Availability zone.
    * `subnet_id` - The id of a VPC subnetwork.
  * `node_info_list` - Node information list, which describe the specification information of various types of nodes in the cluster.
    * `disk_size` - Node disk size.
    * `disk_type` - Node disk type.
    * `node_num` - Number of nodes.
    * `node_type` - Node specification.
    * `type` - Node type.
  * `subnet_id` - The id of a VPC subnetwork.
  * `tags` - A mapping of tags to assign to the instance.
  * `version` - Version of the instance.
  * `vpc_id` - The id of a VPC network.


