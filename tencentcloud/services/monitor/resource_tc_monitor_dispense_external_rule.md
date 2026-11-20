Provides a resource to create a monitor dispense external rule (dispense rule 3.0) which forwards cloud monitor external metrics to Kafka.

Example Usage

```hcl
resource "tencentcloud_monitor_dispense_external_rule" "example" {
  name             = "tf-example"
  ext_namespace    = "QCE/CVM"
  dispense_regions = ["gz"]
  status           = 1
  ext_metrics = [
    "AccOuttraffic",
    "Acs",
    "ActiveMtu",
    "AtsEnabled",
    "BaseCpuUsage",
    "BondingMode",
    "Cnp802pPrio",
    "CbsVolumeFsUsage",
  ]
  period = [60]
  producer = {
    protocol_type = 1
    type          = "Kafka"
    brokers       = "13.132.191.101:8080"
    topic         = "topic-demo"
  }

  dispense_conditions = {
    ext_metric = "AtsEnabled"
    dispense_filters = [
      {
        key        = "InstanceId"
        expression = "="
        values     = ["ins-p97oyxt7"]
      },
      {
        key        = "InstanceId"
        expression = "in"
        values     = ["ins-p97oyxt7,ins-ly4ibb5w"]
      },
    ]
  }
}
```

Import

Monitor dispense external rule can be imported using the rule id, e.g.

```
terraform import tencentcloud_monitor_dispense_external_rule.example 2918
```
