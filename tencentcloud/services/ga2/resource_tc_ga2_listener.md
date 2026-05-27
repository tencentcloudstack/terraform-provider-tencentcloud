Provides a resource to create a GA2 (Global Accelerator 2) listener.

Example Usage

```hcl
resource "tencentcloud_ga2_listener" "example" {
  global_accelerator_id = "ga2-xxxxxxxx"
  name                  = "tf-example"
  protocol              = "TCP"
  listener_type         = "INTELLIGENT"

  port_ranges {
    from_port = 80
    to_port   = 80
  }

  description    = "tf example listener"
  idle_timeout   = 900
  client_affinity = "SOURCE_IP"
  get_real_ip_type = "TOA"
}
```

Import

GA2 listener can be imported using the composite ID `<global_accelerator_id>#<listener_id>`, e.g.

```
terraform import tencentcloud_ga2_listener.example ga2-xxxxxxxx#lbl-xxxxxxxx
```
