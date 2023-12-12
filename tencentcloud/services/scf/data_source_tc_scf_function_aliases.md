Use this data source to query detailed information of scf function_aliases

Example Usage

```hcl
data "tencentcloud_scf_function_aliases" "function_aliases" {
  function_name = "keep-1676351130"
  namespace     = "default"
}
```