Use this data source to query detailed information of HA VIPs.

Example Usage

```hcl
data "tencentcloud_ha_vips" "havips" {
  id         = "havip-kjqwe4ba"
  name       = "test"
  vpc_id     = "vpc-gzea3dd7"
  subnet_id  = "subnet-4d4m4cd4"
  address_ip = "10.0.4.16"
}
```