Use this data source to query detailed information of cynosdb proxy_version

Example Usage

```hcl
data "tencentcloud_cynosdb_proxy_version" "proxy_version" {
  cluster_id     = "cynosdbmysql-bws8h88b"
  proxy_group_id = "cynosdbmysql-proxy-l6zf9t30"
}
```