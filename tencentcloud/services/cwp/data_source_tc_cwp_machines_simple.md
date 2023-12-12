Use this data source to query detailed information of cwp machines_simple

Example Usage

```hcl
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
  project_ids    = [1210293, 1157652]
}
```

Query by Keyword filter

```hcl
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
  project_ids    = [0]

  filters {
    name        = "Keywords"
    values      = ["tf_example"]
    exact_match = true
  }
}
```

Query by Version filter

```hcl
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
  project_ids    = [0]

  filters {
    name        = "Version"
    values      = ["BASIC_VERSION"]
    exact_match = true
  }
}
```

Query by TagId filter

```hcl
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "ALL"
  machine_region = "all-regions"

  filters {
    name        = "TagId"
    values      = ["13771"]
    exact_match = true
  }
}
```