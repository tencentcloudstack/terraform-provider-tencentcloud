Use this data source to query detailed information of CLB

Example Usage

```hcl
data "tencentcloud_clb_instances" "example" {
  clb_id             = "lb-k2zjp9lv"
  network_type       = "OPEN"
  clb_name           = "tf-example"
  project_id         = 0
  result_output_file = "myOutputPath"
}

# Parse JSON fields
output "exclusive_cluster_info" {
  value = jsondecode(data.tencentcloud_clb_instances.example.clb_list[0].exclusive_cluster)
}
``` 