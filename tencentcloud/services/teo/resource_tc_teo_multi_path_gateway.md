Provides a resource to create a TEO multi-path gateway.

Example Usage

Cloud gateway type

```hcl
resource "tencentcloud_teo_multi_path_gateway" "cloud_example" {
  zone_id      = "zone-3fkff38fyw8s"
  gateway_type = "cloud"
  gateway_name = "tf-cloud-gateway"
  gateway_port = 8080
  region_id    = "ap-guangzhou"
}
```

Private gateway type

```hcl
resource "tencentcloud_teo_multi_path_gateway" "private_example" {
  zone_id      = "zone-3fkff38fyw8s"
  gateway_type = "private"
  gateway_name = "tf-private-gateway"
  gateway_port = 9090
  gateway_ip   = "10.0.0.1"
}
```

Import

TEO multi-path gateway can be imported using the zoneId#gatewayId, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway.example zone-3fkff38fyw8s#gw-abc123
```
