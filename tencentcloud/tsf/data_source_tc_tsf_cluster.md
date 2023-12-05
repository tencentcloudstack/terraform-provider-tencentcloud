Use this data source to query detailed information of tsf cluster

Example Usage

```hcl
data "tencentcloud_tsf_cluster" "cluster" {
  cluster_id_list = ["cluster-vwgj5e6y"]
  cluster_type = "V"
  # search_word = ""
  disable_program_auth_check = true
}
```