Use this data source to query detailed information of WeData projects

Example Usage

Query all projects

```hcl
data "tencentcloud_wedata_projects" "example" {}
```

Query projects by filter

```hcl
data "tencentcloud_wedata_projects" "example" {
  project_ids = [
    "2982667120655491072",
    "2853989879663501312"
  ]

  project_name  = "tf_example"
  status        = 1
  project_model = "SIMPLE"
}
```
