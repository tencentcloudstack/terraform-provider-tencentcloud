Use this data source to query detailed information of HA VIP EIP attachments

Example Usage

```hcl
data "tencentcloud_ha_vip_eip_attachments" "foo" {
  havip_id   = "havip-kjqwe4ba"
  address_ip = "1.1.1.1"
}
```