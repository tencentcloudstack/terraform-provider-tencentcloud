Use this data source to get the available product configs of the postgresql instance.

Example Usage

```hcl
data "tencentcloud_postgresql_specinfos" "foo" {
  availability_zone = "ap-shanghai-2"
}
```