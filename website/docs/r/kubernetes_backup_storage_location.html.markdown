---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_backup_storage_location"
sidebar_current: "docs-tencentcloud-resource-kubernetes_backup_storage_location"
description: |-
  Provide a resource to create tke backup storage location.
---

# tencentcloud_kubernetes_backup_storage_location

Provide a resource to create tke backup storage location.

~> **NOTE:** To create this resource, you need to create a cos bucket with prefix "tke-backup" in advance.

## Example Usage

```hcl
resource "tencentcloud_kubernetes_backup_storage_location" "example_backup" {
  name           = "example-backup-1"
  storage_region = "ap-guangzhou"         # region of you pre-created COS bucket
  bucket         = "tke-backup-example-1" # bucket name of your pre-created COS bucket
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Name of the bucket.
* `name` - (Required, String, ForceNew) Name of the backup storage location.
* `storage_region` - (Required, String, ForceNew) Region of the storage.
* `path` - (Optional, String, ForceNew) Prefix of the bucket.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `message` - Message of the backup storage location.
* `state` - State of the backup storage location.


## Import

tke backup storage location can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_backup_storage_location.test xxx
```

