Provides a resource to create a Makers project.

The project id is returned by `CreateProject` and used as the resource `id`.
`ModifyProject` only supports updating `name`, `root_dir`, `output_dir`,
`build_cmd`, `install_cmd`, `framework` and `node_js_version`; every other
configurable field (zone_id, channel, provider, source, area, template_id,
repo_*, namespace, coding_project_id) forces resource replacement when changed.

Note: the cloud API echoes `node_js_version` back as `nodejs_version`.
Removing an updatable attribute from the configuration is rejected because the
cloud API does not document whether `ModifyProject` keeps or clears omitted
fields.

Example Usage

```hcl
resource "tencentcloud_makers_project" "example" {
  zone_id        = "zone-12345"
  name           = "tf-example-project"
  channel        = "Custom"
  provider       = "Github"
  area           = "mainland"
  node_js_version = "20.18.0"
  framework      = "react"
  root_dir       = "./"
  output_dir     = "dist"
  build_cmd      = "npm run build"
  install_cmd    = "npm install"
  repo_url       = "https://github.com/example/tf-example"
  repo_owner     = "example"
  repo_name      = "tf-example"
  repo_branch    = "main"
}
```

Import

Makers project can be imported using the project id, e.g.

```
terraform import tencentcloud_makers_project.example proj-67890
```
