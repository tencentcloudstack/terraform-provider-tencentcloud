Provides a resource to create a antiddos ddos speed limit config

Example Usage

```hcl
resource "tencentcloud_antiddos_ddos_speed_limit_config" "ddos_speed_limit_config" {
  instance_id = "bgp-xxxxxx"
  ddos_speed_limit_config {
		mode = 1
		speed_values {
			type = 1
			value = 1
		}
        speed_values {
			type = 2
			value = 2
		}
		protocol_list = "ALL"
		dst_port_list = "8000"
  }
}
```

Import

antiddos ddos_speed_limit_config can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config ${instanceId}#${configId}s
```