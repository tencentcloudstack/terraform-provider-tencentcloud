Use this data source to query detailed information of CWP machines

Example Usage

```hcl
data "tencentcloud_cwp_machines" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
}
```

Query by Keyword filter

```hcl
data "tencentcloud_cwp_machines" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"

  filters {
    name        = "Keywords"
    values      = ["tf_example"]
    exact_match = true
  }
}
```