Provides a resource to create a antiddos ddos_geo_ip_block_config

Example Usage

```hcl
resource "tencentcloud_antiddos_ddos_geo_ip_block_config" "ddos_geo_ip_block_config" {
	instance_id = "bgp-xxxxxx"
	ddos_geo_ip_block_config {
	  region_type = "customized"
	  action = "drop"
	  area_list = [100002]
	}
}
```

Import

antiddos ddos_geo_ip_block_config can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config ${instanceId}#${configId}
```