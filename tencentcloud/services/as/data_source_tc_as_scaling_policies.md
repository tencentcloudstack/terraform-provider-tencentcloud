Use this data source to query detailed information of scaling policy.

Example Usage

```hcl
data "tencentcloud_as_scaling_policies" "as_scaling_policies" {
  scaling_policy_id  = "asg-mvyghxu7"
  result_output_file = "mytestpath"
}
```