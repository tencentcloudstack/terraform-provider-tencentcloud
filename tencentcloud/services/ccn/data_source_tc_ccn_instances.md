Use this data source to query detailed information of CCN instances.

Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

data "tencentcloud_ccn_instances" "id_instances" {
  ccn_id = tencentcloud_ccn.main.id
}

data "tencentcloud_ccn_instances" "name_instances" {
  name = tencentcloud_ccn.main.name
}
```