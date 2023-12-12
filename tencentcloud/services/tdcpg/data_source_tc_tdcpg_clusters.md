Use this data source to query detailed information of tdcpg clusters.

~> **NOTE:** This data source is still in internal testing. To experience its functions, you need to apply for a whitelist from Tencent Cloud.

Example Usage

```hcl
data "tencentcloud_tdcpg_clusters" "clusters" {
  cluster_id = ""
  cluster_name = ""
  status = ""
  pay_mode = ""
  project_id = ""
  }
```