---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_release_api_group"
sidebar_current: "docs-tencentcloud-resource-tsf_release_api_group"
description: |-
  Provides a resource to create a tsf release_api_group
---

# tencentcloud_tsf_release_api_group

Provides a resource to create a tsf release_api_group

## Example Usage

```hcl
resource "tencentcloud_tsf_release_api_group" "release_api_group" {
  group_id = "grp-qp0rj3zi"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) api group Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ReleaseApiGroup](https://www.tencentcloud.com/search/?wd=ReleaseApiGroup&type=document&lang=en)
* Read: [DescribeApiGroup](https://www.tencentcloud.com/search/?wd=DescribeApiGroup&type=document&lang=en)

