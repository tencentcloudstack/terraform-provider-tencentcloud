Use this data source to query the list of Makers deployment records (build history) of a site project.

Example Usage

Query all deployment records of a Makers project

```hcl
data "tencentcloud_makers_deployments" "example" {
  zone_id    = "zone-3hzk0ndvnnyp"
  project_id = "makers-jw7jha6dktob"
}
```

Query deployment records with filters and sorting

```hcl
data "tencentcloud_makers_deployments" "example" {
  zone_id    = "zone-3hzk0ndvnnyp"
  project_id = "makers-jw7jha6dktob"
  filters = [{
    name   = "Status"
    values = ["Success"]
  }]
  order_by   = "CreatedOn"
  order_type = "Desc"
}
```
