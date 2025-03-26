Provides a resource to create a postgres postgresql_time_window

Example Usage

```hcl
resource "tencentcloud_postgresql_time_window" "postgresql_time_window" {
    db_instance_id      = "postgres-45b0vlmr"
    maintain_duration   = 2
    maintain_start_time = "04:00"
    maintain_week_days  = [
        "friday",
        "monday",
        "saturday",
        "sunday",
        "thursday",
        "tuesday",
        "wednesday",
    ]
}
```

Import

postgres postgresql_time_window can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_time_window.postgresql_time_window instance_id
```
