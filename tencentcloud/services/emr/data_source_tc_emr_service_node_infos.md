Use this data source to query detailed information of emr emr_service_node_infos

Example Usage

```hcl
data "tencentcloud_emr_service_node_infos" "emr_service_node_infos" {
  instance_id = "emr-rzrochgp"
  offset = 1
  limit = 10
  search_text = ""
  conf_status = 2
  maintain_state_id = 2
  operator_state_id = 1
  health_state_id = "2"
  service_name = "YARN"
  node_type_name = "master"
  data_node_maintenance_id = 0
}
```
