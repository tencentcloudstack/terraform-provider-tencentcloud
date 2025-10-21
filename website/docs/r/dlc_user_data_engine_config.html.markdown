---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_user_data_engine_config"
sidebar_current: "docs-tencentcloud-resource-dlc_user_data_engine_config"
description: |-
  Provides a resource to create a dlc user_data_engine_config
---

# tencentcloud_dlc_user_data_engine_config

Provides a resource to create a dlc user_data_engine_config

## Example Usage

```hcl
resource "tencentcloud_dlc_user_data_engine_config" "user_data_engine_config" {
  data_engine_id = "DataEngine-cgkvbas6"
  data_engine_config_pairs {
    config_item  = "qq"
    config_value = "ff"
  }
  session_resource_template {
    driver_size          = "small"
    executor_size        = "small"
    executor_nums        = 1
    executor_max_numbers = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_id` - (Required, String) Engine unique id.
* `data_engine_config_pairs` - (Optional, List) Engine configuration items.
* `session_resource_template` - (Optional, List) Job engine resource configuration template.

The `data_engine_config_pairs` object supports the following:

* `config_item` - (Required, String) Config key.
* `config_value` - (Required, String) Config value.

The `session_resource_template` object supports the following:

* `driver_size` - (Optional, String) Engine driver size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.
* `executor_max_numbers` - (Optional, Int) Specify the executor max number (in a dynamic configuration scenario), the minimum value is 1, and the maximum value is less than the cluster specification (when ExecutorMaxNumbers is less than ExecutorNums, the value is set to ExecutorNums).
* `executor_nums` - (Optional, Int) Specify the number of executors. The minimum value is 1 and the maximum value is less than the cluster specification.
* `executor_size` - (Optional, String) Engine executor size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dlc user_data_engine_config can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_user_data_engine_config.user_data_engine_config user_data_engine_config_id
```

