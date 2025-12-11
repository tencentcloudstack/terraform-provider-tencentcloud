Provides a resource to create a APM prometheus rule

Example Usage

```hcl
resource "tencentcloud_apm_prometheus_rule" "example" {
  instance_id       = "apm-lhqHyRBuA"
  name              = "tf-example"
  service_name      = "java-market-service"
  metric_match_type = 0
  metric_name_rule  = "task.duration"
  status            = 1
}
```

Import

APM prometheus rule can be imported using the instanceId#ruleId, e.g.

```
terraform import tencentcloud_apm_prometheus_rule.example apm-lhqHyRBuA#140
```
