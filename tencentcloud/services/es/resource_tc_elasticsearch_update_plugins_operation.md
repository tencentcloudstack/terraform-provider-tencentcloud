Provides a resource to update elasticsearch plugins

Example Usage

```hcl
resource "tencentcloud_elasticsearch_update_plugins_operation" "update_plugins_operation" {
  instance_id = "es-xxxxxx"
  install_plugin_list = ["analysis-pinyin"]
  force_restart = false
  force_update = true
}
```