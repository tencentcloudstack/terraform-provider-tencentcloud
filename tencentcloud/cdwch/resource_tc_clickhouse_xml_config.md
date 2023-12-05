Provides a resource to create a clickhouse xml_config

Example Usage

```hcl
resource "tencentcloud_clickhouse_xml_config" "xml_config" {
  instance_id = "cdwch-datuhk3z"
  modify_conf_context {
    file_name      = "metrika.xml"
    old_conf_value = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPHlhbmRleD4KPC95YW5kZXg+Cg=="
    new_conf_value = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPHlhbmRleD4KICAgIDx6b29rZWVwZXItc2VydmVycz4KICAgIDwvem9va2VlcGVyLXNlcnZlcnM+CjwveWFuZGV4Pgo="
    file_path      = "/etc/clickhouse-server"
  }
}
```