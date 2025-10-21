---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket_inventory"
sidebar_current: "docs-tencentcloud-resource-cos_bucket_inventory"
description: |-
  Provides a resource to create a cos bucket inventory
---

# tencentcloud_cos_bucket_inventory

Provides a resource to create a cos bucket inventory

~> **NOTE:** The current resource does not support cdc.

## Example Usage

```hcl
# get user info
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

# create cos
resource "tencentcloud_cos_bucket" "example" {
  bucket = "private-bucket-${local.app_id}"
  acl    = "private"
}

# create cos bucket inventory
resource "tencentcloud_cos_bucket_inventory" "example" {
  name                     = "tf-example"
  bucket                   = tencentcloud_cos_bucket.example.id
  is_enabled               = "true"
  included_object_versions = "Current"

  optional_fields {
    fields = ["Size", "ETag"]
  }

  filter {
    period {
      start_time = "1687276800"
    }
  }

  schedule {
    frequency = "Daily"
  }

  destination {
    bucket = "qcs::cos:ap-guangzhou::private-bucket-1309118522"
    format = "CSV"
    prefix = "frontends"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Bucket name.
* `destination` - (Required, List) Information about the inventory result destination.
* `included_object_versions` - (Required, String) Whether to include object versions in the inventory. All or No.
* `is_enabled` - (Required, String) Whether to enable the inventory. true or false.
* `name` - (Required, String, ForceNew) Inventory Name.
* `schedule` - (Required, List) Inventory job cycle.
* `filter` - (Optional, List) Filters objects prefixed with the specified value to analyze.
* `optional_fields` - (Optional, List) Analysis items to include in the inventory result	.

The `destination` object supports the following:

* `bucket` - (Required, String) Bucket name.
* `format` - (Required, String) Format of the inventory result. Valid value: CSV.
* `account_id` - (Optional, String) ID of the bucket owner.
* `encryption` - (Optional, List) Server-side encryption for the inventory result.
* `prefix` - (Optional, String) Prefix of the inventory result.

The `encryption` object of `destination` supports the following:

* `sse_cos` - (Optional, String) Encryption with COS-managed key. This field can be left empty.

The `filter` object supports the following:

* `period` - (Optional, List) Creation time range of the objects to analyze.
* `prefix` - (Optional, String) Prefix of the objects to analyze.

The `optional_fields` object supports the following:

* `fields` - (Optional, Set) Optional analysis items to include in the inventory result. The optional fields include Size, LastModifiedDate, StorageClass, ETag, IsMultipartUploaded, ReplicationStatus, Tag, Crc64, and x-cos-meta-*.

The `period` object of `filter` supports the following:

* `end_time` - (Optional, String) Creation end time of the objects to analyze. The parameter is a timestamp in seconds, for example, 1568688762.
* `start_time` - (Optional, String) Creation start time of the objects to analyze. The parameter is a timestamp in seconds, for example, 1568688761.

The `schedule` object supports the following:

* `frequency` - (Required, String) Frequency of the inventory job. Enumerated values: Daily, Weekly.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cos bucket inventory can be imported using the id, e.g.

```
terraform import tencentcloud_cos_bucket_inventory.example private-bucket-1309118522#tf-example
```

