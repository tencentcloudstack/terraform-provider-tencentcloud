Provides a resource to create a CDC dedicated cluster

Example Usage

```hcl
# create cdc site
resource "tencentcloud_cdc_site" "example" {
  name         = "tf-example"
  country      = "China"
  province     = "Guangdong Province"
  city         = "Guangzhou"
  address_line = "Tencent Building"
  description  = "desc."
}

# create cdc dedicated cluster
resource "tencentcloud_cdc_dedicated_cluster" "example" {
  site_id     = tencentcloud_cdc_site.example.id
  name        = "tf-example"
  zone        = "ap-guangzhou-6"
  description = "desc."
}
```

Import

CDC dedicated cluster can be imported using the id, e.g.

```
terraform import tencentcloud_cdc_dedicated_cluster.example cluster-d574omhk
```
