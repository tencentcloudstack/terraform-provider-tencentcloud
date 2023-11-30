Use this data source to query detailed information of scf function_address

Example Usage

```hcl
data "tencentcloud_scf_function_address" "function_address" {
  function_name = "keep-1676351130"
  namespace     = "default"
  qualifier     = "$LATEST"
}
```