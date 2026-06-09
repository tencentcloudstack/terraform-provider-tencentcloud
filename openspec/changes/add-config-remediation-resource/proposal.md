# Add tencentcloud_config_remediation Resource

## What

Add a new Terraform resource `tencentcloud_config_remediation` for managing Tencent Cloud Config rule remediation settings. This resource allows users to create, read, update, and delete remediation configurations bound to a config rule.

## Why

The Config compliance evaluation workflow requires remediation actions to fix non-compliant resources. Currently, the provider has no resource for managing remediation settings. Users need to:

- Bind a remediation (SCF function) to a config rule to automatically or manually fix non-compliant resources
- Manage the invoke type (manual/auto/none) and source template

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `CreateRemediation` | Returns `RemediationId` as the resource unique ID |
| Read | `ListRemediations` | Query by `RuleIds` filter; find matching `RemediationId` in response |
| Update | `UpdateRemediation` | Update `RemediationType`, `RemediationTemplateId`, `InvokeType`, `SourceType` |
| Delete | `DeleteRemediations` | Pass `[RemediationId]` array |

## Resource ID

The unique resource ID is `RemediationId` (e.g., `crr-lKj43O4nbSD78XYlvGS9`).
