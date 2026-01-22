Use this data source to query detailed information of CFW nat fw switches

Example Usage

Query Nat instance'switch by instance ID

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
}
```

Or filter by switch enable status

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
  enable     = 1
}
```
