Provides a resource to create a WeData project

Example Usage

```hcl
resource "tencentcloud_wedata_project" "example" {
  project {
    project_name  = "tf_example"
    display_name  = "display_name"
    project_model = "SIMPLE"
  }

  dlc_info {
    compute_resources     = ["svmgao_stability"]
    region                = "ap-guangzhou"
    default_database      = "db_name"
    standard_mode_env_tag = "Dev"
    access_account        = "OWNER"
  }

  resource_ids = [
    "20250909193110713075",
    "20250820215449817917"
  ]

  status = 1
}
```
