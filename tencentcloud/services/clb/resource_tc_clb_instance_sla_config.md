Provides a resource to create a clb instance_sla_config

Example Usage

```hcl
resource "tencentcloud_clb_instance_sla_config" "instance_sla_config" {
  load_balancer_id = "lb-5dnrkgry"
  sla_type         = "SLA"
}
```

Import

clb instance_sla_config can be imported using the id, e.g.

```
terraform import tencentcloud_clb_instance_sla_config.instance_sla_config instance_id
```