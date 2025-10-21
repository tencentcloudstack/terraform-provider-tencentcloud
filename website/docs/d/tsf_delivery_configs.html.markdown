---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_delivery_configs"
sidebar_current: "docs-tencentcloud-datasource-tsf_delivery_configs"
description: |-
  Use this data source to query detailed information of tsf delivery_configs
---

# tencentcloud_tsf_delivery_configs

Use this data source to query detailed information of tsf delivery_configs

## Example Usage

```hcl
data "tencentcloud_tsf_delivery_configs" "delivery_configs" {
  search_word = "test"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) search word.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - deploy group information about the deployment group associated with a delivery item.Note: This field may return null, which means that no valid value was obtained.
  * `content` - content. Note: This field may return null, which means that no valid value was obtained.
    * `collect_path` - harvest log path. Note: This field may return null, which means that no valid value was obtained.
    * `config_id` - config id.
    * `config_name` - config name.
    * `create_time` - Creation time.Note: This field may return null, indicating that no valid values can be obtained.
    * `custom_rule` - CustomRule specifies a custom line separator rule.Note: This field may return null, which means that no valid value was obtained.
    * `enable_auth` - whether use auth for kafka. Note: This field may return null, which means that no valid value was obtained.
    * `enable_global_line_rule` - Indicates whether a single row rule should be applied.Note: This field may return null, which means that no valid value was obtained.
    * `groups` - Associated deployment group information.Note: This field may return null, indicating that no valid values can be obtained.
      * `associate_time` - Associate Time. Note: This field may return null, indicating that no valid values can be obtained.
      * `cluster_id` - Cluster ID. Note: This field may return null, indicating that no valid values can be obtained.
      * `cluster_name` - Cluster Name. Note: This field may return null, indicating that no valid values can be obtained.
      * `cluster_type` - Cluster type.
      * `group_id` - Group Id.
      * `group_name` - Group Name.
      * `namespace_name` - Namespace Name. Note: This field may return null, indicating that no valid values can be obtained.
    * `kafka_address` - KafkaAddress refers to the address of a Kafka server.Note: This field may return null, which means that no valid value was obtained.
    * `kafka_infos` - Kafka Infos. Note: This field may return null, which means that no valid value was obtained.
      * `custom_rule` - Custom Line Rule.
      * `line_rule` - Line rule specifies the type of line separator used in a file. It can have one of the following values: default: The default line separator is used to separate lines in the file. time: The lines in the file are separated based on time. custom: A custom line separator is used. In this case, the CustomRule field should be filled with the specific custom value. Note: This field may return null, which means that no valid value was obtained.
      * `path` - harvest log path. Note: This field may return null, which means that no valid value was obtained.
      * `topic` - Kafka topic. Note: This field may return null, which means that no valid value was obtained.
    * `kafka_v_ip` - Kafka VIP. Note: This field may return null, which means that no valid value was obtained.
    * `kafka_v_port` - Kafka VPort. Note: This field may return null, which means that no valid value was obtained.
    * `line_rule` - Line Rule for log. Note: This field may return null, which means that no valid value was obtained.
    * `password` - Password. Note: This field may return null, which means that no valid value was obtained.
    * `topic` - Topic. Note: This field may return null, which means that no valid value was obtained.
    * `username` - user Name. Note: This field may return null, which means that no valid value was obtained.
  * `total_count` - total count. Note: This field may return null, which means that no valid value was obtained.


