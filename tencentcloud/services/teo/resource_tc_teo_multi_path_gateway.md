Provides a resource to create a teo multi path gateway for EdgeOne(TEO).

Example Usage

Cloud type gateway

```hcl
resource "tencentcloud_teo_multi_path_gateway" "cloud" {
  zone_id      = "zone-359h792djt7h"
  gateway_type = "cloud"
  gateway_name = "test-cloud-gw"
  region_id    = "ap-guangzhou"
  status       = "online"
}
```

Private type gateway

```hcl
resource "tencentcloud_teo_multi_path_gateway" "private" {
  zone_id      = "zone-359h792djt7h"
  gateway_type = "private"
  gateway_name = "test-private-gw"
  gateway_ip   = "1.2.3.4"
  gateway_port = 8080
  status       = "online"
}
```

The `status` attribute is optional and accepts the following values:

- `online`: enable the gateway.
- `offline`: disable the gateway.

If `status` is not set, the value is populated by the server automatically.

Import

teo multi path gateway can be imported using the id, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway.example zone-279qso5a4cw9#mpgw-g3176ppeye
```
