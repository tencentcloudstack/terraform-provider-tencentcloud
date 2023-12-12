Use this data source to query detailed information of dedicated tunnels instances.

Example Usage

```hcl
data "tencentcloud_dcx_instances" "name_select" {
  name = "main"
}

data "tencentcloud_dcx_instances" "id" {
  dcx_id = "dcx-3ikuw30k"
}
```