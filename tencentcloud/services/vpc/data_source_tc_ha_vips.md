Use this data source to query detailed information of HA VIPs.

Example Usage

Query all HA vips

```hcl
data "tencentcloud_ha_vips" "ha_vips" {}
```

Query HA vips by filters

```hcl
data "tencentcloud_ha_vips" "ha_vips" {
  name = "tf-example"
}

data "tencentcloud_ha_vips" "ha_vips" {
  id = "havip-rg9y1k2c"
}

data "tencentcloud_ha_vips" "ha_vips" {
  vpc_id = "vpc-q23dnivj"
}

data "tencentcloud_ha_vips" "ha_vips" {
  subnet_id = "subnet-g6c7yi7o"
}
```
