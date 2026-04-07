Use this data source to get the available product configs of the PostgreSQL specifications.

Example Usage

```hcl
data "tencentcloud_postgresql_specinfos" "example" {
  availability_zone = "ap-guangzhou-7"
  storage_type      = "CLOUD_HSSD"
}
```
