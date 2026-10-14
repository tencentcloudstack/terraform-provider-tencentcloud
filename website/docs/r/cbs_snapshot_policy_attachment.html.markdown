---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_snapshot_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-cbs_snapshot_policy_attachment"
description: |-
  Provides a CBS snapshot policy attachment resource.
---

# tencentcloud_cbs_snapshot_policy_attachment

Provides a CBS snapshot policy attachment resource.

~> **NOTE:** To distinguish between `storage_id` and `storage_id`, use `storage_id` when there is only one diskId, otherwise use `storage_ids`.

## Example Usage

### Attachment CBS snapshot policy by storage_id

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 60
  availability_zone = "ap-guangzhou-6"
  project_id        = 0
  encrypt           = false

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cbs_snapshot_policy" "example" {
  snapshot_policy_name = "tf-example"
  repeat_weekdays      = [1, 4]
  repeat_hours         = [1]
  retention_days       = 7
}

resource "tencentcloud_cbs_snapshot_policy_attachment" "example" {
  storage_id         = tencentcloud_cbs_storage.example.id
  snapshot_policy_id = tencentcloud_cbs_snapshot_policy.example.id
}
```

### Attachment CBS snapshot policy by storage_ids

```hcl
resource "tencentcloud_cbs_storage" "example1" {
  storage_name      = "tf-example1"
  storage_type      = "CLOUD_SSD"
  storage_size      = 60
  availability_zone = "ap-guangzhou-6"
  project_id        = 0
  encrypt           = false

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cbs_storage" "example2" {
  storage_name      = "tf-example2"
  storage_type      = "CLOUD_SSD"
  storage_size      = 60
  availability_zone = "ap-guangzhou-6"
  project_id        = 0
  encrypt           = false

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cbs_snapshot_policy" "example" {
  snapshot_policy_name = "tf-example"
  repeat_weekdays      = [1, 4]
  repeat_hours         = [1]
  retention_days       = 7
}

resource "tencentcloud_cbs_snapshot_policy_attachment" "example" {
  storage_ids = [
    tencentcloud_cbs_storage.example1.id,
    tencentcloud_cbs_storage.example2.id,
  ]
  snapshot_policy_id = tencentcloud_cbs_snapshot_policy.example.id
}
```

## Argument Reference

The following arguments are supported:

* `snapshot_policy_id` - (Required, String, ForceNew) ID of CBS snapshot policy.
* `storage_id` - (Optional, String, ForceNew) ID of CBS.
* `storage_ids` - (Optional, Set: [`String`], ForceNew) IDs of CBS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [BindAutoSnapshotPolicy](https://www.tencentcloud.com/document/product/362/32168) - This API (BindAutoSnapshotPolicy) is used to bind the cloud disk to the specified scheduled snapshot policy.
* Read: [DescribeAutoSnapshotPolicies](https://www.tencentcloud.com/document/product/362/32165) - This API is used to query scheduled snapshot policies.
* Read: [DescribeDiskAssociatedAutoSnapshotPolicy](https://www.tencentcloud.com/document/product/362/32164) - This API (DescribeDiskAssociatedAutoSnapshotPolicy) is used to query the scheduled snapshot policy bound to a cloud disk.
* Delete: [UnbindAutoSnapshotPolicy](https://www.tencentcloud.com/document/product/362/32161) - This API (UnbindAutoSnapshotPolicy) is used to unbind the cloud disk from the specified scheduled snapshot policy.

## Import

CBS snapshot policy attachment can be imported using the id, e.g.

If use storage_id

```
$ terraform import tencentcloud_cbs_snapshot_policy_attachment.example disk-fesgc43m#asp-8abupspr
```

If use storage_ids

```
$ terraform import tencentcloud_cbs_snapshot_policy_attachment.example disk-ghylus9y,disk-0tm61hla#asp-ng87uf4t
```

