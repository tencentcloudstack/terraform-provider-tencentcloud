Use this data source to query detailed information of cam list_entities_for_policy

Example Usage

```hcl
data "tencentcloud_cam_list_entities_for_policy" "list_entities_for_policy" {
  policy_id = 1
  entity_filter = "All"
    }
```