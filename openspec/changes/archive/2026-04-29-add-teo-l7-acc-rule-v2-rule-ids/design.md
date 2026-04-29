## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource manages individual TEO L7 acceleration rules. It previously had a `rule_ids` computed attribute that exposed the full list of rule IDs from the `CreateL7AccRules` API response. However, since Terraform manages one rule per resource instance, only `rule_id` (singular) is needed. The SDK's `RuleEngineAction` struct also supports `VaryParameters` (a simple `Switch` on/off field) which was not exposed.

## Goals / Non-Goals

**Goals:**
- Remove redundant `rule_ids` computed attribute
- Add `Vary` action support via `vary_parameters` block
- Add `OriginAuthentication` action support via `origin_authentication_parameters` block
- Maintain backward compatibility for existing configurations

**Non-Goals:**
- Changing `rule_id` behavior

## Decisions

1. **Remove `rule_ids`**: Since TF handles single resources, the array `rule_ids` is redundant with `rule_id`. Removing simplifies the schema.

2. **Add `vary_parameters` following `ContentCompression` pattern**: `VaryParameters` has the same structure as `ContentCompressionParameters` (single `Switch *string` field). It's placed after `content_compression_parameters` in all 3 locations (schema, Get, Set).

3. **Add `origin_authentication_parameters` with nested list**: `OriginAuthenticationParameters` contains `RequestProperties []*OriginAuthenticationRequestProperties`, each with `Type`, `Name`, `Value` string fields. Implemented as `origin_authentication_parameters.request_properties` list of objects.

4. **Update `name` field description**: Fix `SetContentIdentifierParameters` → `SetContentIdentifier`, add missing action names (`Vary`, `ContentCompression`, `OriginAuthentication`) to align with SDK's `RuleEngineAction.Name` definition. `Shield` is excluded as it is not exposed in this resource.

5. **Remove unit tests for `rule_ids`**: The 5 gomonkey mock tests were specifically for `rule_ids`. Acceptance tests remain.

## Risks / Trade-offs

- **Breaking change for `rule_ids` users**: Users who referenced `rule_ids` will get an error. Since the resource manages a single rule, `rule_id` provides the same information.
- **Adding `Vary`**: No risk, purely additive Optional action.
