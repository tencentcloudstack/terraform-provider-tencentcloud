/*
Use this data source to query detailed information of chdfs access_groups

Example Usage

```hcl
data "tencentcloud_chdfs_access_groups" "access_groups" {
  vpc_id = "vpc-pewdpc0d"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudChdfsAccessGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudChdfsAccessGroupsRead,
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "get groups belongs to the vpc id, must set but only can use one of VpcId and OwnerUin to get the groups.",
			},

			"owner_uin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "get groups belongs to the owner uin, must set but only can use one of VpcId and OwnerUin to get the groups.",
			},

			"access_groups": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "access group list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "access group id.",
						},
						"access_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "access group name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "access group description.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"vpc_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "vpc network type(1:CVM, 2:BM 1.0).",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudChdfsAccessGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_chdfs_access_groups.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("vpc_id"); ok {
		paramMap["vpc_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("owner_uin"); ok {
		paramMap["owner_uin"] = helper.IntUint64(v.(int))
	}

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var accessGroups []*chdfs.AccessGroup

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeChdfsAccessGroupsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		accessGroups = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(accessGroups))
	tmpList := make([]map[string]interface{}, 0, len(accessGroups))

	if accessGroups != nil {
		for _, accessGroup := range accessGroups {
			accessGroupMap := map[string]interface{}{}

			if accessGroup.AccessGroupId != nil {
				accessGroupMap["access_group_id"] = accessGroup.AccessGroupId
			}

			if accessGroup.AccessGroupName != nil {
				accessGroupMap["access_group_name"] = accessGroup.AccessGroupName
			}

			if accessGroup.Description != nil {
				accessGroupMap["description"] = accessGroup.Description
			}

			if accessGroup.CreateTime != nil {
				accessGroupMap["create_time"] = accessGroup.CreateTime
			}

			if accessGroup.VpcType != nil {
				accessGroupMap["vpc_type"] = accessGroup.VpcType
			}

			if accessGroup.VpcId != nil {
				accessGroupMap["vpc_id"] = accessGroup.VpcId
			}

			ids = append(ids, *accessGroup.AccessGroupId)
			tmpList = append(tmpList, accessGroupMap)
		}

		_ = d.Set("access_groups", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
