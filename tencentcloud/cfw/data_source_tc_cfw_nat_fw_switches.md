Use this data source to query detailed information of cfw nat_fw_switches

Example Usage

Query Nat instance'switch by instance id

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
}
```

Or filter by switch status

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
  status     = 1
}
```