/*
Use this data source to query detailed information of CAM groups

Example Usage

```hcl
data "tencentcloud_cam_group_memberships" "foo" {
  group_id = "12515263"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudCamGroupMemberships() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamGroupMembershipsRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of CAM group to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"membership_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CAM group membership. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of CAM group.",
						},
						"user_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Id set of the cam group members.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamGroupMembershipsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ca,_groups.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	groupId := d.Get("group_id").(string)
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var memberships []*string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeGroupMembershipById(ctx, groupId)
		if e != nil {
			return retryError(e)
		}
		memberships = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM group memberships failed, reason:%s\n ", logId, err.Error())
		return err
	}
	groupList := make([]map[string]interface{}, 0, 1)
	ids := make([]string, 0, 1)
	mapping := map[string]interface{}{
		"group_id": groupId,
		"user_ids": memberships,
	}
	groupList = append(groupList, mapping)
	ids = append(ids, groupId)

	d.SetId(dataResourceIdsHash(ids))
	if e := d.Set("membership_list", groupList); e != nil {
		log.Printf("[CRITAL]%s provider set membershiplist fail, reason:%s\n ", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), groupList); e != nil {
			return e
		}
	}

	return nil
}
