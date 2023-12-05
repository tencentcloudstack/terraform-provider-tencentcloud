Provides a resource to create a antiddos cc black white ip

Example Usage

```hcl
resource "tencentcloud_antiddos_cc_black_white_ip" "cc_black_white_ip" {
  instance_id = "bgpip-xxxxxx"
  black_white_ip {
    ip   = "1.2.3.5"
    mask = 0

  }
  type     = "black"
  ip       = "xxx.xxx.xxx.xxx"
  domain   = "t.baidu.com"
  protocol = "http"
}
```

Import

antiddos cc_black_white_ip can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_cc_black_white_ip.cc_black_white_ip ${instanceId}#${policyId}#${instanceIp}#${domain}#${protocol}
```