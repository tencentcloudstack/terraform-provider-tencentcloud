Use this data source to query the detail information of CFS access group.

Example Usage

```hcl
data "tencentcloud_cfs_access_groups" "access_groups" {
  access_group_id = "pgroup-7nx89k7l"
  name            = "test"
}
```