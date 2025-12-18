Provides a resource to create a antiddos ddos black white ip

Example Usage

```hcl
resource "tencentcloud_antiddos_ddos_black_white_ip" "ddos_black_white_ip" {
  instance_id = "bgp-xxxxxx"
  ip = "1.2.3.5"
  mask = 0
  type = "black"
}
```

Import

antiddos ddos_black_white_ip can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip ${instanceId}#${ip}
```