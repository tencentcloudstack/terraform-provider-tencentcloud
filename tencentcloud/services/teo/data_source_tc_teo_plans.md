Use this data source to query detailed information of TEO plans

Example Usage

Query all plans

```hcl
data "tencentcloud_teo_plans" "example" {}
```

Query plans by filters

```hcl
data "tencentcloud_teo_plans" "example" {
  order     = "expire-time"
  direction = "desc"
  filters {
    name = "plan-id"
    values = [
      "edgeone-2o1xvpmq7nn",
      "edgeone-2mezmk9s2xdx"
    ]
  }

  filters {
    name = "plan-type"
    values = [
      "plan-trial",
      "plan-personal",
      "plan-basic",
      "plan-standard",
      "plan-enterprise"
    ]
  }

  filters {
    name = "area"
    values = [
      "overseas",
      "mainland",
      "global"
    ]
  }

  filters {
    name = "status"
    values = [
      "normal",
      "expiring-soon",
      "expired",
      "isolated"
    ]
  }
}
```
