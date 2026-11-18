Use this data source to query the list of Makers deployment records (build history) of a site project.

Example Usage

Query all deployment records of a Makers project

```hcl
data "tencentcloud_makers_deployments" "example" {
  zone_id    = "zone-abcdefg"
  project_id = "prj-123456"
}
```

Query deployment records with filters and pagination

```hcl
data "tencentcloud_makers_deployments" "example" {
  zone_id    = "zone-abcdefg"
  project_id = "prj-123456"
  filters {
    name   = "Status"
    values = ["Success"]
  }
  order_by   = "CreatedOn"
  order_type = "Desc"
  limit      = 100
  offset     = 0
}
```
