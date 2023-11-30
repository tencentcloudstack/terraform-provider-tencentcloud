Provides a resource to create a antiddos ip. Only support for bgp-multip.

Example Usage

```hcl
resource "tencentcloud_dayu_ddos_ip_attachment_v2" "boundip" {
  id = "bgp-xxxxxx"
  bound_ip_list {
	ip = "1.1.1.1"
	biz_type = "public"
	instance_id = "ins-xxx"
	device_type = "cvm"
  }
}
```