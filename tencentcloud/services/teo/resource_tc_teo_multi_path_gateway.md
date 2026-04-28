Provides a resource to create a TEO multi-path gateway.

Example Usage

Cloud type gateway

```hcl
resource "tencentcloud_teo_multi_path_gateway" "example" {
  zone_id      = "zone-3fkff38fyw8s"
  gateway_type = "cloud"
  gateway_name = "tf-example-cloud"
  gateway_port = 8080
  region_id    = "ap-guangzhou"
}
```

Private type gateway

```hcl
resource "tencentcloud_teo_multi_path_gateway" "example" {
  zone_id      = "zone-3fkff38fyw8s"
  gateway_type = "private"
  gateway_name = "tf-example-private"
  gateway_port = 8080
  gateway_ip   = "10.0.0.1"
}
```

Import

TEO multi-path gateway can be imported using the zoneId#gatewayId, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway.example zone-3fkff38fyw8s#gw-2qrk328yw8s
```
