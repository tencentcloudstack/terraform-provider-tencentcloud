---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_config_attachment"
sidebar_current: "docs-tencentcloud-resource-cls_config_attachment"
description: |-
  Provides a resource to create a cls config attachment
---

# tencentcloud_cls_config_attachment

Provides a resource to create a cls config attachment

## Example Usage



## Argument Reference

The following arguments are supported:

* `config_id` - (Required, String, ForceNew) Collection configuration id.
* `group_id` - (Required, String, ForceNew) Machine group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ApplyConfigToMachineGroup](https://www.tencentcloud.com/document/product/614/42814) - This API is used to apply the collection configuration to a specified machine group.
* Read: [DescribeConfigMachineGroups](https://www.tencentcloud.com/document/product/614/42810) - This API is used to get the machine group bound to collection rule configuration.
* Delete: [DeleteConfigFromMachineGroup](https://www.tencentcloud.com/document/product/614/42811) - This API is used to delete the collection configuration applied to a machine group.

## Import

cls config_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cls_config_attachment.attach config_id#group_id
```

```

