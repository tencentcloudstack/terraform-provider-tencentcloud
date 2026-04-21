# Add EMR Boot Script Config Resource

## What

Add a new Terraform config resource `tencentcloud_emr_boot_script_config` that manages the boot script configuration for a specific EMR instance and boot type.

## Why

Users managing EMR clusters via Terraform need to configure and maintain boot scripts (pre-execution scripts) that run at different lifecycle stages of the cluster. Without this resource, they cannot manage boot scripts declaratively.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create / Update | `ModifyBootScript` | Sync, no async polling required |
| Read | `DescribeBootScript` | Returns `Detail` with per-type script lists |
| Delete | no-op | Config resource, delete does nothing |

## Resource ID

`{instance_id}#{boot_type}` — composite of the two required fields that uniquely identify a boot script configuration slot.

## Config Resource Pattern

This resource follows the **config resource** pattern (like `tencentcloud_config_deliver_config`):
- Create = `d.SetId(...)` then call Update
- Update = call `ModifyBootScript` then call Read  
- Delete = no-op
