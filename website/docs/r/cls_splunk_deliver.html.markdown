---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_splunk_deliver"
sidebar_current: "docs-tencentcloud-resource-cls_splunk_deliver"
description: |-
  Provides a resource to create a CLS splunk deliver.
---

# tencentcloud_cls_splunk_deliver

Provides a resource to create a CLS splunk deliver.

## Example Usage

```hcl
resource "tencentcloud_cls_splunk_deliver" "example" {
  topic_id  = "eb417b1d-fc00-46f2-85f7-47ca0848a6b3"
  name      = "tf-example"
  index_ack = 1

  net_info {
    host     = "110.11.22.106"
    port     = 8088
    token    = "e27274fb-****-****-****-****00206282"
    net_type = 2
    is_ssl   = true
  }

  metadata_info {
    format = "json"
    meta_fields = [
      "__HOSTNAME__",
      "__FILENAME__",
      "__TIMESTAMP__",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `metadata_info` - (Required, List) Splunk deliver task metadata information.
* `name` - (Required, String) Splunk deliver task name.
* `net_info` - (Required, List) Splunk deliver task target network information.
* `topic_id` - (Required, String, ForceNew) Log topic ID.
* `channel` - (Optional, String) Advanced configuration - channel. Required if indexer is enabled.
* `dsl_filter` - (Optional, String) Log pre-filtering DSL statement for raw data written to Splunk.
* `enable` - (Optional, Int) Delivery task enable status. 0: disable, 1: enable.
* `external_role` - (Optional, List) Advanced configuration - cross-account delivery role authorization information.
* `has_service_log` - (Optional, Int) Whether to enable service log. 1: disable, 2: enable. Default: enable.
* `index_ack` - (Optional, Int) Whether to enable indexer. 1: disable, 2: enable. Default: 1.
* `index` - (Optional, String) Advanced configuration - Splunk index. No more than 64 characters.
* `source_type` - (Optional, String) Advanced configuration - data source type. No more than 64 characters.
* `source` - (Optional, String) Advanced configuration - data source. No more than 64 characters.

The `external_role` object supports the following:

* `external_id` - (Required, String) Cross-account delivery role name.
* `role_arn` - (Required, String) Cross-account delivery role RoleArn.

The `metadata_info` object supports the following:

* `format` - (Required, String) Data format. Valid values: rawlog, json.
* `enable_tag` - (Optional, Bool) Whether to deliver __TAG__ field.
* `meta_fields` - (Optional, Set) Delivery fields, including __SOURCE__, __FILENAME__, __TIMESTAMP__, __HOSTNAME__, __PKG_ID__.
* `tag_json_tiled` - (Optional, Bool) Whether to flatten JSON. Required when enable_tag is true.

The `net_info` object supports the following:

* `host` - (Required, String) Network address.
* `net_type` - (Required, Int) Network type. 1: internal network, 2: external network.
* `port` - (Required, Int) Port.
* `token` - (Required, String) Authentication token.
* `is_ssl` - (Optional, Bool) Whether to use SSL. Default is false.
* `virtual_gateway_type` - (Optional, Int) Network service type. Required when net_type is internal network. 0: CVM, 3: Direct Connect Gateway, 11: CCN, 1025: CLB.
* `vpc_id` - (Optional, String) VPC ID. Required when net_type is internal network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `task_id` - Splunk deliver task ID.


## Import

cls splunk deliver can be imported using the topicId#taskId, e.g.

```
terraform import tencentcloud_cls_splunk_deliver.example eb417b1d-fc00-46f2-85f7-47ca0848a6b3#716dae34-7eba-4b09-a3d9-6fe8e510e236
```

