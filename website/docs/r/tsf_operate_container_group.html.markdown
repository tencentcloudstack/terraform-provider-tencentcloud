---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_operate_container_group"
sidebar_current: "docs-tencentcloud-resource-tsf_operate_container_group"
description: |-
  Provides a resource to create a tsf operate_container_group
---

# tencentcloud_tsf_operate_container_group

Provides a resource to create a tsf operate_container_group

## Example Usage

```hcl
resource "tencentcloud_tsf_operate_container_group" "operate_container_group" {
  group_id = "group-ynd95rea"
  operate  = "stop"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) group Id.
* `operate` - (Required, String) Operation, `start`- start the container, `stop`- stop the container.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeContainerGroupAttribute](https://www.tencentcloud.com/search/?wd=DescribeContainerGroupAttribute&type=document&lang=en)
  - [StartContainerGroup](https://www.tencentcloud.com/search/?wd=StartContainerGroup&type=document&lang=en)
  - [StopContainerGroup](https://www.tencentcloud.com/search/?wd=StopContainerGroup&type=document&lang=en)
* Read: [DescribeContainerGroupAttribute](https://www.tencentcloud.com/search/?wd=DescribeContainerGroupAttribute&type=document&lang=en)
* Update:
  - [DescribeContainerGroupAttribute](https://www.tencentcloud.com/search/?wd=DescribeContainerGroupAttribute&type=document&lang=en)
  - [StartContainerGroup](https://www.tencentcloud.com/search/?wd=StartContainerGroup&type=document&lang=en)
  - [StopContainerGroup](https://www.tencentcloud.com/search/?wd=StopContainerGroup&type=document&lang=en)

