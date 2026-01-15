Provides a resource to create a WeData authorize data source

Example Usage

Authorize by project ids

```hcl
resource "tencentcloud_wedata_authorize_data_source" "example" {
  data_source_id = "116203"
  auth_project_ids = [
    "1857740139240632320",
    "1857740139240632318",
  ]
}
```

Authorize by users

```hcl
resource "tencentcloud_wedata_authorize_data_source" "example" {
  data_source_id = "116203"
  auth_users = [
    "1857740139240632320_100028448903",
    "1857740139240632320_100028578751",
  ]
}

```

Import

WeData authorize data source can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_authorize_data_source.example 116203
```
