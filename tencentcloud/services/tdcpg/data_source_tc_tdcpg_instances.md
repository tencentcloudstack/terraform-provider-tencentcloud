Use this data source to query detailed information of tdcpg instances.

~> **NOTE:** This data source is still in internal testing. To experience its functions, you need to apply for a whitelist from Tencent Cloud.

Example Usage

```hcl
data "tencentcloud_tdcpg_instances" "instances" {
  cluster_id = ""
  instance_id = ""
  instance_name = ""
  status = ""
  instance_type = ""
  }
```