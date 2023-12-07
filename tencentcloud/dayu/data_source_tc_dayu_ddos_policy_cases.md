Use this data source to query dayu DDoS policy cases

Example Usage

```hcl
data "tencentcloud_dayu_ddos_policy_cases" "id_test" {
  resource_type = tencentcloud_dayu_ddos_policy_case.test_policy_case.resource_type
  scene_id      = tencentcloud_dayu_ddos_policy_case.test_policy_case.scene_id
}
```