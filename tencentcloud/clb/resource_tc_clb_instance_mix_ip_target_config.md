Provides a resource to create a clb instance_mix_ip_target_config

Example Usage

```hcl
resource "tencentcloud_clb_instance_mix_ip_target_config" "instance_mix_ip_target_config" {
  load_balancer_id = "lb-5dnrkgry"
  mix_ip_target = false
}
```

Import

clb instance_mix_ip_target_config can be imported using the id, e.g.

```
terraform import tencentcloud_clb_instance_mix_ip_target_config.instance_mix_ip_target_config instance_id
```