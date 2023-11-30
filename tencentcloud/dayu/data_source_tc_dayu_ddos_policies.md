Use this data source to query dayu DDoS policies

Example Usage

```hcl
data "tencentcloud_dayu_ddos_policies" "id_test" {
  resource_type = tencentcloud_dayu_ddos_policy.test_policy.resource_type
  policy_id     = tencentcloud_dayu_ddos_policy.test_policy.policy_id
}
```