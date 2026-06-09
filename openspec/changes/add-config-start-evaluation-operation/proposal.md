# Add tencentcloud_config_start_config_rule_evaluation_operation

## What

Add a new one-shot Terraform operation resource `tencentcloud_config_start_config_rule_evaluation_operation` that manually triggers a Config rule evaluation via the `StartConfigRuleEvaluation` API.

## Why

In the Tencent Cloud Config compliance workflow, evaluation results are generated either periodically or on configuration change. Users sometimes need to manually trigger an evaluation (e.g., after fixing a non-compliant resource) to get an updated compliance result immediately. This operation resource exposes that capability in Terraform.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create (trigger) | `StartConfigRuleEvaluation` | Triggers evaluation; returns only `RequestId` |

## Resource Lifecycle

This is a one-shot operation resource (no real cloud object is created):
- **Create**: Call `StartConfigRuleEvaluation`; set ID to `helper.BuildToken()`
- **Read**: No-op (always returns nil)
- **Delete**: No-op (nothing to tear down)
- No `Update`, no `Importer`
