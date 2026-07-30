## Context

The `tencentcloud_instance` resource manages CVM instances through multiple CVM APIs. Currently, the resource supports KMS encryption for data disks via `data_disks.*.kms_key_id` (a nested field within the `data_disks` block). System disk encryption through a custom KMS key is not yet supported at the Terraform level.

The vendor SDK (`tencentcloud-sdk-go/tencentcloud/cvm/v20170312`) has been updated: the `SystemDisk` struct now includes `KmsKeyId` and `Encrypt` fields. According to the SDK documentation, these fields are currently used for the `RunInstances` API and are in a grayscale/gated rollout phase.

The existing `tencentcloud_instance` resource already handles system disk parameters (`system_disk_type`, `system_disk_size`, `system_disk_id`, `system_disk_name`) at the top level of the schema. The new `kms_key_id` parameter follows the same pattern.

## Goals / Non-Goals

**Goals:**
- Add a top-level `kms_key_id` (TypeString, Optional, ForceNew) parameter to `tencentcloud_instance` for system disk encryption
- Map the parameter to `SystemDisk.KmsKeyId` in the `RunInstances` call during Create
- Map the parameter to `SystemDisk.KmsKeyId` in the `ResetInstance` call during Update (when image_id changes or reinstall is triggered)

**Non-Goals:**
- Not adding a separate `system_disk_encrypt` bool parameter (the `Encrypt` field on `SystemDisk` is implied by providing `KmsKeyId`)
- Not modifying the existing `data_disks.*.kms_key_id` behavior
- Not adding `kms_key_id` to the Read path return (the `DescribeInstances` API does not return `SystemDisk.KmsKeyId`)

## Decisions

**Decision 1: Top-level `kms_key_id` rather than nested under a `system_disk` block**

The existing resource uses top-level fields for system disk configuration (`system_disk_type`, `system_disk_size`, `system_disk_id`, `system_disk_name`). Adding `kms_key_id` as a top-level field follows the established pattern and avoids a breaking schema change.

**Decision 2: ForceNew for `kms_key_id`**

System disk encryption is configured at instance creation time via `RunInstances`. Changing the encryption key after creation requires recreating the instance. Therefore, `ForceNew: true` is the correct behavior.

**Decision 3: Pass `kms_key_id` to `ResetInstance` as well**

When the instance is reinstalled (via `ResetInstance` API), the system disk is recreated. Passing `kms_key_id` to `ResetInstance` ensures the encryption key is preserved during reinstall operations.

**Decision 4: No Read-back of `kms_key_id`**

The `DescribeInstances` API response does not include `SystemDisk.KmsKeyId`. Since `ForceNew: true` is set, the parameter is only used during creation and reinstall. The value in Terraform state will be the user-provided value, which is consistent with how other ForceNew-only parameters work.

## Risks / Trade-offs

- **Risk**: The `KmsKeyId` field in `SystemDisk` is documented as "灰度中" (in grayscale rollout). If the API does not yet support this field for all users, calling it may return an error.
  - **Mitigation**: The field is Optional. Users who don't set it are unaffected. Users who set it must have access to the feature.

- **Risk**: The `kms_key_id` value cannot be read back from the API, so drift detection is limited.
  - **Mitigation**: This is inherent to the API design. The `ForceNew: true` setting means any change to `kms_key_id` will trigger recreation, which is the correct behavior for this parameter.