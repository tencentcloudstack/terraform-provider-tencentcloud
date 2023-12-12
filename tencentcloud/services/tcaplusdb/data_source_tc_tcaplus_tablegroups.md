Use this data source to query table groups of the TcaplusDB cluster.

Example Usage

```hcl
data "tencentcloud_tcaplus_tablegroups" "null" {
  cluster_id = "19162256624"
}
data "tencentcloud_tcaplus_tablegroups" "id" {
  cluster_id    = "19162256624"
  tablegroup_id = "19162256624:1"
}
data "tencentcloud_tcaplus_tablegroups" "name" {
  cluster_id      = "19162256624"
  tablegroup_name = "test"
}
data "tencentcloud_tcaplus_tablegroups" "all" {
  cluster_id      = "19162256624"
  tablegroup_id   = "19162256624:1"
  tablegroup_name = "test"
}
```