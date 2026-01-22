Provides a resource to create a vpc ha_vip_instance_attachment

Example Usage

```hcl
resource "tencentcloud_ha_vip_instance_attachment" "ha_vip_instance_attachment" {
  instance_id = "eni-xxxxxx"
  ha_vip_id = "havip-xxxxxx"
  instance_type = "ENI"
}
```

Import

vpc ha_vip_instance_attachment can be imported using the id(${haVipId}#${instanceType}#${instanceId}), e.g.

```
terraform import tencentcloud_ha_vip_instance_attachment.ha_vip_instance_attachment ha_vip_instance_attachment_id
```
