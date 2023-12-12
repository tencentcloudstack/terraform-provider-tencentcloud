Use this data source to query detailed information of direct connect gateway instances.

Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test"
  network_instance_id = tencentcloud_ccn.main.id
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}

#You need to sleep for a few seconds because there is a cache on the server
data "tencentcloud_dc_gateway_instances" "name_select" {
  name = tencentcloud_dc_gateway.ccn_main.name
}

data "tencentcloud_dc_gateway_instances" "id_select" {
  dcg_id = tencentcloud_dc_gateway.ccn_main.id
}
```