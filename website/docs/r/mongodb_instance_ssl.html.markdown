---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_ssl"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance_ssl"
description: |-
  Provides a resource to manage MongoDB instance SSL configuration.
---

# tencentcloud_mongodb_instance_ssl

Provides a resource to manage MongoDB instance SSL configuration.

~> **NOTE:** This resource is used to enable or disable SSL for MongoDB instances. When the resource is destroyed, SSL will be disabled automatically.

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance_ssl" "example" {
  instance_id = "cmgo-xxxxxxxx"
  enable      = true
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Bool) Whether to enable SSL. Valid values: `true` - enable SSL, `false` - disable SSL.
* `instance_id` - (Required, String, ForceNew) Instance ID, for example: cmgo-p8vnipr5.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cert_url` - Certificate download link. This field is only available when SSL is enabled.
* `expired_time` - Certificate expiration time, format: 2023-05-01 12:00:00. This field is only available when SSL is enabled.
* `status` - SSL status. Valid values: `0` - SSL is disabled, `1` - SSL is enabled.


## Import

MongoDB instance SSL configuration can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_instance_ssl.example cmgo-xxxxxxxx
```

