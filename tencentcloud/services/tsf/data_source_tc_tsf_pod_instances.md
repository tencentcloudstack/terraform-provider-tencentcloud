Use this data source to query detailed information of tsf pod_instances

Example Usage

```hcl
data "tencentcloud_tsf_pod_instances" "pod_instances" {
  group_id = "group-ynd95rea"
  pod_name_list = ["keep-terraform-6f8f977688-zvphm"]
}
```