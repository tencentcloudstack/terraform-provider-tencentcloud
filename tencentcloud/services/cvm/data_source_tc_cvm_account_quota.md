Use this data source to query CVM account quota details.

Example Usage

Basic query without filters

```hcl
data "tencentcloud_cvm_account_quota" "quota" {}

output "app_id" {
  value = data.tencentcloud_cvm_account_quota.quota.app_id
}
```

Query by availability zone

```hcl
data "tencentcloud_cvm_account_quota" "quota_zone" {
  zone = ["ap-guangzhou-3", "ap-guangzhou-4"]
}
```

Query by quota type

```hcl
data "tencentcloud_cvm_account_quota" "quota_type" {
  quota_type = "PostPaidQuotaSet"
}
```

Query with multiple filters

```hcl
data "tencentcloud_cvm_account_quota" "quota_filtered" {
  zone       = ["ap-guangzhou-3"]
  quota_type = "PostPaidQuotaSet"
}
```

Query with result output file

```hcl
data "tencentcloud_cvm_account_quota" "quota_output" {
  zone               = ["ap-guangzhou-3"]
  result_output_file = "./quota.json"
}
```
