---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_key_pair_attachment"
sidebar_current: "docs-tencentcloud-resource-lighthouse_key_pair_attachment"
description: |-
  Provides a resource to create a lighthouse key_pair_attachment
---

# tencentcloud_lighthouse_key_pair_attachment

Provides a resource to create a lighthouse key_pair_attachment

## Example Usage

```hcl
resource "tencentcloud_lighthouse_key_pair_attachment" "key_pair_attachment" {
  key_id      = "lhkp-xxxxxx"
  instance_id = "lhins-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `key_id` - (Required, String, ForceNew) Key pair ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [AssociateInstancesKeyPairs](https://www.tencentcloud.com/document/api/1103/42484)
  - [DescribeInstances](https://www.tencentcloud.com/document/api/1103/42498)
* Read: [DescribeKeyPairs](https://www.tencentcloud.com/document/api/1103/42480)
* Delete:
  - [DescribeInstances](https://www.tencentcloud.com/document/api/1103/42498)
  - [DisassociateInstancesKeyPairs](https://www.tencentcloud.com/document/api/1103/42479)

## Import

lighthouse key_pair_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_key_pair_attachment.key_pair_attachment key_pair_attachment_id
```

