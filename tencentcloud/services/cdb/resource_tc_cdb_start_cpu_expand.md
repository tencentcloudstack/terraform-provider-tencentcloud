Provides a resource to create a CDB CPU elastic expand attachment

Example Usage

Auto expansion type

```hcl
resource "tencentcloud_cdb_start_cpu_expand" "example" {
  instance_id = "cdb-test1234"
  type        = "auto"

  auto_strategy {
    expand_threshold     = 80
    shrink_threshold     = 20
    expand_second_period = 300
    shrink_second_period = 600
  }
}
```

Manual expansion type

```hcl
resource "tencentcloud_cdb_start_cpu_expand" "example" {
  instance_id = "cdb-test1234"
  type        = "manual"
  expand_cpu  = 4
}
```

TimeInterval expansion type

```hcl
resource "tencentcloud_cdb_start_cpu_expand" "example" {
  instance_id = "cdb-test1234"
  type        = "timeInterval"
  expand_cpu  = 4

  time_interval_strategy {
    start_time = 1709251200
    end_time   = 1709337600
  }
}
```

Period expansion type

```hcl
resource "tencentcloud_cdb_start_cpu_expand" "example" {
  instance_id = "cdb-test1234"
  type        = "period"
  expand_cpu  = 4

  period_strategy {
    time_cycle {
      monday    = true
      tuesday   = true
      wednesday = true
      thursday  = true
      friday    = true
      saturday  = false
      sunday    = false
    }

    time_interval {
      start_time = "09:00"
      end_time   = "18:00"
    }
  }
}
```

Import

CDB start cpu expand can be imported using the instance_id, e.g.

```
terraform import tencentcloud_cdb_start_cpu_expand.example cdb-test1234
```
