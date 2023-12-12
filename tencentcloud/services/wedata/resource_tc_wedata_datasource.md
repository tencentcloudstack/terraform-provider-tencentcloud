Provides a resource to create a wedata datasource

Example Usage

```hcl
resource "tencentcloud_wedata_datasource" "example" {
  name                = "tf_example"
  category            = "DB"
  type                = "MYSQL"
  owner_project_id    = "1612982498218618880"
  owner_project_name  = "project_demo"
  owner_project_ident = "体验项目"
  description         = "description."
  display             = "tf_example_demo"
  status              = 1
  cos_bucket          = "wedata-agent-sh-1257305158"
  cos_region          = "ap-shanghai"
  params              = jsonencode({
    "connectType" : "public",
    "authorityType" : "true",
    "deployType" : "CONNSTR_PUBLICDB",
    "url" : "jdbc:mysql://1.1.1.1:8080/database",
    "username" : "root",
    "password" : "password",
    "type" : "MYSQL"
  })
}
```