/*
Use this data source to query detailed information of CAM groups

Example Usage

```hcl
data "tencentcloud_cam_groups" "foo" {
  group_id = "12515263"
  name     = "cam-role-test"
  remark   = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
)

func dataSourceTencentCloudCamGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamGroupsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CAM group to be queried.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of CAM group to be queried.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the cam group to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"group_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CAM groups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CAM group.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of CAM group.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the CAM group.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_groups.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		params["group_id"], _ = strconv.Atoi(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		params["name"] = v.(string)
	}
	if v, ok := d.GetOk("remark"); ok {
		params["remark"] = v.(string)
	}

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var groups []*cam.GroupInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeGroupsByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		groups = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM groups failed, reason:%s\n", logId, err.Error())
		return err
	}
	groupList := make([]map[string]interface{}, 0, len(groups))
	ids := make([]string, 0, len(groups))
	for _, group := range groups {
		mapping := map[string]interface{}{
			"name":        *group.GroupName,
			"create_time": *group.CreateTime,
		}
		if group.Remark != nil {
			mapping["remark"] = *group.Remark
		}
		groupList = append(groupList, mapping)
		ids = append(ids, strconv.Itoa(int(*group.GroupId)))
	}

	d.SetId(dataResourceIdsHash(ids))
	if e := d.Set("group_list", groupList); e != nil {
		log.Printf("[CRITAL]%s provider set group list fail, reason:%s\n", logId, e.Error())
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
