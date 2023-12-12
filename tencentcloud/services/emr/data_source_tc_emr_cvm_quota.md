Use this data source to query detailed information of emr cvm_quota

Example Usage

```hcl
data "tencentcloud_emr_cvm_quota" "cvm_quota" {
  cluster_id = "emr-0ze36vnp"
  zone_id    = 100003
}
```