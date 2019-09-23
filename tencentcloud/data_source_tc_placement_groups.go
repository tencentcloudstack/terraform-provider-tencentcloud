/*
Use this data source to query placement groups.

Example Usage

```hcl
data "tencentcloud_placement_groups" "foo" {
  placement_group_id = "ps-21q9ibvr"
  name               = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func dataSourceTencentCloudPlacementGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPlacementGroupsRead,

		Schema: map[string]*schema.Schema{
			"placement_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the placement group to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the placement group to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"placement_group_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of placement group. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"placement_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the placement group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the placement group.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the placement group.",
						},
						"cvm_quota_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of hosts in the placement group.",
						},
						"current_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of hosts in the placement group.",
						},
						"instance_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Host IDs in the placement group.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the placement group.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudPlacementGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_placement_groups.read")
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var placementGroupId string
	var name string
	if v, ok := d.GetOk("placement_group_id"); ok {
		placementGroupId = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	var placementGroups []*cvm.DisasterRecoverGroup
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		placementGroups, errRet = cvmService.DescribePlacementGroupByFilter(ctx, placementGroupId, name)
		if errRet != nil {
			return retryError(errRet, "InternalError")
		}
		return nil
	})
	if err != nil {
		return err
	}

	placementGroupList := make([]map[string]interface{}, 0, len(placementGroups))
	ids := make([]string, 0, len(placementGroups))
	for _, placement := range placementGroups {
		mapping := map[string]interface{}{
			"placement_group_id": placement.DisasterRecoverGroupId,
			"name":               placement.Name,
			"type":               placement.Type,
			"cvm_quota_total":    placement.CvmQuotaTotal,
			"current_num":        placement.CurrentNum,
			"instance_ids":       flattenStringList(placement.InstanceIds),
			"create_time":        placement.CreateTime,
		}
		placementGroupList = append(placementGroupList, mapping)
		ids = append(ids, *placement.DisasterRecoverGroupId)
	}

	d.SetId(dataResourceIdsHash(ids))
	err = d.Set("placement_group_list", placementGroupList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set placement group list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), placementGroupList); err != nil {
			return err
		}
	}
	return nil
}
