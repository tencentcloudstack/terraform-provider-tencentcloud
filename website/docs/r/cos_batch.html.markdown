---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_batch"
sidebar_current: "docs-tencentcloud-resource-cos_batch"
description: |-
  Provides a resource to create a cos bucket batch.
---

# tencentcloud_cos_batch

Provides a resource to create a cos bucket batch.

## Example Usage

```hcl
resource "tencentcloud_cos_batch" "cos_batch" {
  uin                   = "100022975249"
  appid                 = "1308919341"
  confirmation_required = true
  description           = "cos_batch"
  priority              = 1
  status                = "Cancelled"
  role_arn              = "qcs::cam::uin/100022975249:roleName/COSBatch_QCSRole"
  manifest {
    location {
      etag       = "64357de8fd75a3abae2200135a2c9627"
      object_arn = "qcs::cos:ap-guangzhou:uid/1308919341:keep-test-1308919341/cos_bucket_inventory/1308919341/keep-test/test/20230621/manifest.json"
    }
    spec {
      format = "COSInventoryReport_CSV_V1"
    }
  }
  operation {
    cos_put_object_copy {
      access_control_directive = "Copy"
      metadata_directive       = "Copy"
      prefix_replace           = false
      storage_class            = "STANDARD"
      tagging_directive        = "Copy"
      target_resource          = "qcs::cos:ap-guangzhou:uid/1308919341:cos-lock-1308919341"
    }
  }
  report {
    bucket       = "qcs::cos:ap-guangzhou:uid/1308919341:keep-test-1308919341"
    enabled      = "true"
    format       = "Report_CSV_V1"
    report_scope = "AllTasks"
  }
}
```

## Argument Reference

The following arguments are supported:

* `appid` - (Required, Int, ForceNew) Appid.
* `manifest` - (Required, List, ForceNew) List of objects to be processed.
* `operation` - (Required, List, ForceNew) Select the action to be performed on the objects in the manifest file.
* `priority` - (Required, Int) Mission priority. The higher the value, the higher the priority of the task. Priority values range from 0 to 2147483647.
* `report` - (Required, List, ForceNew) Task completion report.
* `role_arn` - (Required, String, ForceNew) COS resource identifier, which is used to identify the role you created. You need this resource identifier to verify your identity.
* `uin` - (Required, String, ForceNew) Uin.
* `confirmation_required` - (Optional, Bool, ForceNew) Whether to confirm before performing the task. The default is false.
* `description` - (Optional, String, ForceNew) Mission description. If you configured this information when you created the task, the content is returned. The description length ranges from 0 to 256 bytes.
* `status` - (Optional, String) Current status of the task.
Legal parameter values include Active, Cancelled, Cancelling, Complete, Completing, Failed, Failing, New, Paused, Pausing, Preparing, Ready, Suspended.
For Update status, when you move a task to the Ready state, COS will assume that you have confirmed the task and will perform it. When you move a task to the Cancelled state, COS cancels the task. Optional parameters include: Ready, Cancelled.

The `access_control_grants` object supports the following:

* `identifier` - (Required, String) User ID (UIN) in qcs format. For example: qcs::cam::uin/100000000001:uin/100000000001.
* `permission` - (Required, String) Specify a permission to be granted. Enumerated value: READ,WRITE,FULL_CONTROL.
* `type_identifier` - (Required, String) Specifies the type of Identifier. Currently, only user ID is supported. Enumerated value: ID.
* `display_name` - (Optional, String) User name.

The `cos_initiate_restore_object` object supports the following:

* `expiration_in_days` - (Required, Int) Sets the number of days after which the copy will be automatically expired and deleted, an integer in the range of 1-365.
* `job_tier` - (Required, String) Select archive recovery model. Available values: Bulk, Standard.

The `cos_put_object_copy` object supports the following:

* `target_resource` - (Required, String) Sets the target bucket for the Copy. Use qcs to specify, for example, qcs::cos:ap-chengdu:uid/1250000000:examplebucket-1250000000.
* `access_control_directive` - (Optional, String) This element specifies how ACL is copied. Valid values:
- Copy: inherits the source object ACL
- Replaced: replace source ACL
- Add: add a new ACL based on the source ACL.
* `access_control_grants` - (Optional, List) Controls the specific access to the object.
* `canned_access_control_list` - (Optional, String) Defines the ACL property of the object. Valid values: private, public-read.
* `metadata_directive` - (Optional, String) This element specifies whether to copy object metadata from the source object or replace it with metadata in the < NewObjectMetadata > element. Valid values are: Copy, Replaced, Add. Copy: inherit source object metadata; Replaced: replace source metadata; Add: add new metadata based on source metadata.
* `modified_since_constraint` - (Optional, Int) When the object is modified after the specified time, the operation is performed, otherwise 412 is returned.
* `new_object_metadata` - (Optional, List) Configure the metadata for the object.
* `new_object_tagging` - (Optional, List) The label of the configuration object, which must be specified when the < TaggingDirective > value is Replace or Add.
* `prefix_replace` - (Optional, Bool) Specifies whether the prefix of the source object needs to be replaced. A value of true indicates the replacement object prefix, which needs to be used with <ResourcesPrefix> and <TargetKeyPrefix>. Default value: false.
* `resources_prefix` - (Optional, String) This field is valid only when the < PrefixReplace > value is true. Specify the source object prefix to be replaced, and the replacement directory should end with `/`. Can be empty with a maximum length of 1024 bytes.
* `storage_class` - (Optional, String) Sets the storage level of the object. Enumerated value: STANDARD,STANDARD_IA. Default value: STANDARD.
* `tagging_directive` - (Optional, String) This element specifies whether to copy the object tag from the source object or replace it with the tag in the < NewObjectTagging > element. Valid values are: Copy, Replaced, Add. Copy: inherits the source object tag; Replaced: replaces the source tag; Add: adds a new tag based on the source tag.
* `target_key_prefix` - (Optional, String) This field is valid only when the <PrefixReplace> value is true. This value represents the replaced prefix, and the replacement directory should end with /. Can be empty with a maximum length of 1024 bytes.
* `unmodified_since_constraint` - (Optional, Int) When the object has not been modified after the specified time, the operation is performed, otherwise 412 is returned.

The `location` object supports the following:

* `etag` - (Required, String) Specifies the etag of the object list. Length 1-1024 bytes.
* `object_arn` - (Required, String) Specifies the unique resource identifier of the object manifest, which is 1-1024 bytes long.
* `object_version_id` - (Optional, String) Specifies the version of the object manifest ID, which is 1-1024 bytes long.

The `manifest` object supports the following:

* `location` - (Required, List) The location information of the list of objects.
* `spec` - (Required, List) Format information that describes the list of objects. If it is a CSV file, this element describes the fields contained in the manifest.

The `new_object_metadata` object supports the following:

* `cache_control` - (Optional, String) The caching instructions defined in RFC 2616 are saved as object metadata.
* `content_disposition` - (Optional, String) The file name defined in RFC 2616 is saved as object metadata.
* `content_encoding` - (Optional, String) The encoding format defined in RFC 2616 is saved as object metadata.
* `content_type` - (Optional, String) The content types defined in RFC 2616 are saved as object metadata.
* `http_expires_date` - (Optional, String) The cache expiration time defined in RFC 2616 is saved as object metadata.
* `sse_algorithm` - (Optional, String) Server encryption algorithm. Currently, only AES256 is supported.
* `user_metadata` - (Optional, List) Includes user-defined metadata.

The `new_object_tagging` object supports the following:

* `key` - (Required, String) key.
* `value` - (Required, String) value.

The `operation` object supports the following:

* `cos_initiate_restore_object` - (Optional, List) Specifies the specific parameters for the batch restore operation for archive storage type objects in the inventory.
* `cos_put_object_copy` - (Optional, List) Specifies the specific parameters for the batch copy operation on the objects in the list.

The `report` object supports the following:

* `bucket` - (Required, String) Delivery bucket for task completion reports.
* `enabled` - (Required, String) Whether to output the task completion report.
* `format` - (Required, String) Task completion report format information. Legal value: Report_CSV_V1.
* `report_scope` - (Required, String) Task completion report the task information that needs to be recorded to determine whether to record the execution information of all operations or the information of failed operations. Legal values: AllTasks, FailedTasksOnly.
* `prefix` - (Optional, String) Prefix information for the task completion report. Length 0-256 bytes.

The `spec` object supports the following:

* `format` - (Required, String) Specifies the format information for the list of objects. Legal fields are: COSBatchOperations_CSV_V1, COSInventoryReport_CSV_V1.
* `fields` - (Optional, List) Describes the fields contained in the listing, which you need to use to specify CSV file fields when Format is COSBatchOperations_CSV_V1. Legal fields are: Ignore, Bucket, Key, VersionId.

The `user_metadata` object supports the following:

* `key` - (Required, String) key.
* `value` - (Required, String) value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `job_id` - Job id.


## Import

cos bucket batch can be imported using the id, e.g.

```
$ terraform import tencentcloud_cos_batch.cos_batch ${uin}#${appid}#{job_id}
```

