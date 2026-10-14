---
subcategory: "TencentCloud Elastic Microservice(TEM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tem_app_config"
sidebar_current: "docs-tencentcloud-resource-tem_app_config"
description: |-
  Provides a resource to create a tem appConfig
---

# tencentcloud_tem_app_config

Provides a resource to create a tem appConfig

## Example Usage

```hcl
resource "tencentcloud_tem_app_config" "appConfig" {
  environment_id = "en-o5edaepv"
  name           = "demo"
  config_data {
    key   = "key"
    value = "value"
  }
  config_data {
    key   = "key1"
    value = "value1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `config_data` - (Required, List) payload.
* `environment_id` - (Required, String, ForceNew) environment ID.
* `name` - (Required, String, ForceNew) appConfig name.

The `config_data` object supports the following:

* `key` - (Required, String) key.
* `value` - (Required, String) value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateConfigData](https://cloud.tencent.com/document/api/1371/78337) - This API is used to create a configuration.
* Read: [DescribeConfigData](https://cloud.tencent.com/document/api/1371/78336) - This API is used to query details of a configuration.
* Update: [ModifyConfigData](https://cloud.tencent.com/document/api/1371/78333) - This API is used to modify a configuration.
* Delete: [DestroyConfigData](https://cloud.tencent.com/document/api/1371/78334) - This API is used to terminate a configuration.

## Import

tem appConfig can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_app_config.appConfig environmentId#name
```

