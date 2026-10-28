## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource manages individual TEO L7 acceleration rules. The Read function previously only checked `if respData == nil` to determine if the resource was deleted. However, the DescribeL7AccRules API can return a valid response with an empty `Rules` array when the rule no longer exists, which was not handled.

## Goals / Non-Goals

**Goals:**
- Fix Read function to handle empty Rules array from API response
- Ensure resource is properly marked as deleted when API returns empty Rules

**Non-Goals:**
- Changing any other resource behavior

## Decisions

1. **Fix Read function nil check**: Change `if respData == nil {` to `if respData == nil || len(respData.Rules) == 0 {` to properly handle the case where the API returns a response with an empty Rules array (resource deleted).

## Risks / Trade-offs

- No risks. This is a purely defensive fix that handles an edge case in the API response.
