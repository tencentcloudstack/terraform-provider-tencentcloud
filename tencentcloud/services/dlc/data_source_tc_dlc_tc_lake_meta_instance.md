Use this data source to query the activation status of DLC (Data Lake Compute) TCLake meta instance.

Example Usage

```hcl
data "tencentcloud_dlc_tc_lake_meta_instance" "example" {
}

output "status" {
  value = data.tencentcloud_dlc_tc_lake_meta_instance.example.status
}
```
