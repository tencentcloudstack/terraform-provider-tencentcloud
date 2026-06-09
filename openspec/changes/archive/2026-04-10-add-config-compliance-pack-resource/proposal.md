# Add tencentcloud_config_compliance_pack Resource

## What

Add a new Terraform resource `tencentcloud_config_compliance_pack` for managing Tencent Cloud Config compliance packs. This resource allows users to create, read, update, and delete compliance packs that batch-manage and evaluate resource configuration compliance.

## Why

The Tencent Cloud Config service (配置审计) compliance pack feature is not yet supported in the terraform-provider-tencentcloud. Users need to manage compliance packs to organize and evaluate multiple resource configuration rules at scale. This resource closes the gap by supporting the full CRUD lifecycle.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `AddCompliancePack` | Returns `CompliancePackId` as the resource unique ID |
| Read | `DescribeCompliancePack` | Query by `CompliancePackId` |
| Update (content) | `UpdateCompliancePack` | Update name, description, risk level, rules |
| Update (status) | `UpdateCompliancePackStatus` | Enable/disable the compliance pack |
| Delete | `DeleteCompliancePack` | Requires pack to be in UN_ACTIVE status first |

## Resource ID

The unique resource ID is `CompliancePackId` (e.g., `cp-33mA27YUlOJWG4sJ53Sx`).
