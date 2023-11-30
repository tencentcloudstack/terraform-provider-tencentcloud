Provides a resource to create a dnspod download_snapshot
Example Usage
```hcl
resource "tencentcloud_dnspod_download_snapshot_operation" "download_snapshot" {
  domain = "dnspod.cn"
  snapshot_id = "456"
}
```