# Add tencentcloud_config_rule Resource

## What

Add a new Terraform resource `tencentcloud_config_rule` for managing Tencent Cloud Config rules. This resource supports the full CRUD lifecycle including enabling/disabling rules via OpenConfigRule/CloseConfigRule.

## Why

Config rules define compliance requirements for cloud resources. Currently no Terraform resource exists for managing them independently (they can be added via compliance packs but not as standalone rules).

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `AddConfigRule` | Returns `RuleId` (String) |
| Read | `DescribeConfigRule` | Query by `RuleId` |
| Update (content) | `UpdateConfigRule` | Update name, triggers, risk level, params, etc. |
| Update (status) | `OpenConfigRule` / `CloseConfigRule` | Toggle rule status |
| Delete | `DeleteConfigRule` | Pass `RuleId` |

## Resource ID

`RuleId` (String, e.g. `cr-3xhsd76j603v0a8ma0i73`).
