## 1. Schema Definition

- [x] 1.1 Add computed `rules` attribute to the `tencentcloud_teo_l7_acc_rule_v2` resource schema in `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`. The attribute is a `TypeList` of `schema.Resource` with fields: `status`, `rule_id`, `rule_name`, `description`, `branches`, `rule_priority`, all set to `Computed: true`.

## 2. Read Function Update

- [x] 2.1 Update the `ResourceTencentCloudTeoL7AccRuleV2Read` function to populate the `rules` attribute from `respData.Rules` by flattening the `[]*RuleEngineItem` slice into the format expected by Terraform's `d.Set()`.

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` to document the new computed `rules` attribute.

## 4. Unit Tests

- [x] 4.1 Add unit tests in `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` to verify the `rules` attribute is correctly populated using gomonkey mock approach.
