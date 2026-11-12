Provides a resource to create a monitor dispense external rule (dispense rule 3.0) which forwards cloud monitor external metrics to Kafka.

-> **Note:** `ext_metrics` is not refreshed from the describe response. The describe API reports metric details through `rule.ext_metric`, but the mapping between the requested metric names and the reported details is not documented, so drift on `ext_metrics` cannot be detected. Keep the configured value authoritative.

-> **Note:** Removing `dispense_regions`, `ext_metrics`, `period`, `dispense_conditions`, `producer.merge`, `producer.global_tags`, `producer.default_tags`, `producer.username` or `producer.password` from the configuration is rejected, because the modify API does not document whether an omitted value means "keep" or "clear". Keep the attribute in the configuration or recreate the resource instead.

-> **Note:** Destroying this resource calls `DeleteDispenseExternalRule` with the single `rule_id` of the rule. The effect on the Kafka side and on related resources is not documented by the API.

Example Usage

```hcl
resource "tencentcloud_monitor_dispense_external_rule" "example" {
  name          = "tf-example"
  ext_namespace = "ns-tf-example"

  producer = {
    protocol_type = 0
    type          = "kafka"
    brokers       = "broker-1:9092"
    topic         = "topic-1"
    merge         = 1
    global_tags = [
      {
        key   = "env"
        value = "prod"
      },
    ]
    default_tags = ["dim-1"]
    username     = "user"
    password     = "pass"
  }

  dispense_regions = ["ap-guangzhou"]
  ext_metrics      = ["metric-1"]
  period           = [60]

  dispense_conditions = {
    ext_metric = "metric-1"
    dispense_filters = [
      {
        key        = "dim"
        values     = ["v1", "v2"]
        expression = "eq"
      },
    ]
  }
}
```

Import

Monitor dispense external rule can be imported using the rule id, e.g.

```
terraform import tencentcloud_monitor_dispense_external_rule.example 12345
```
