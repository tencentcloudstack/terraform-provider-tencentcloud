Use this data source to query TcaplusDB tables.

Example Usage

```hcl
data "tencentcloud_tcaplus_tables" "null" {
  cluster_id = "19162256624"
}

data "tencentcloud_tcaplus_tables" "tablegroup" {
  cluster_id    = "19162256624"
  tablegroup_id = "19162256624:3"
}

data "tencentcloud_tcaplus_tables" "name" {
  cluster_id    = "19162256624"
  tablegroup_id = "19162256624:3"
  table_name    = "guagua"
}

data "tencentcloud_tcaplus_tables" "id" {
  cluster_id = "19162256624"
  table_id   = "tcaplus-faa65eb7"
}
data "tencentcloud_tcaplus_tables" "all" {
  cluster_id    = "19162256624"
  tablegroup_id = "19162256624:3"
  table_id      = "tcaplus-faa65eb7"
  table_name    = "guagua"
}
```