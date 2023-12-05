package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCamGroupUserAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamGroupUserAccountRead,
		Schema: map[string]*schema.Schema{
			"uid": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Sub-user uid.",
			},

			"rp": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number per page. The default is 20.",
			},

			"sub_uin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Sub-user uin.",
			},

			"total_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of user groups the sub-user has joined.",
			},

			"group_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "User group information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User group ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User group name.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark.",
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

func dataSourceTencentCloudCamGroupUserAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_group_user_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("uid"); ok {
		paramMap["Uid"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("rp"); ok {
		paramMap["Rp"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("sub_uin"); ok {
		paramMap["SubUin"] = helper.IntUint64(v.(int))
	}

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	var groupInfoList []*cam.GroupInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCamGroupUserAccountByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		groupInfoList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(groupInfoList))
	tmpList := make([]map[string]interface{}, 0, len(groupInfoList))

	if len(groupInfoList) > 0 {
		_ = d.Set("total_num", len(groupInfoList))
	}

	if groupInfoList != nil {
		for _, groupInfo := range groupInfoList {
			groupInfoMap := map[string]interface{}{}

			if groupInfo.GroupId != nil {
				groupInfoMap["group_id"] = groupInfo.GroupId
			}

			if groupInfo.GroupName != nil {
				groupInfoMap["group_name"] = groupInfo.GroupName
			}

			if groupInfo.CreateTime != nil {
				groupInfoMap["create_time"] = groupInfo.CreateTime
			}

			if groupInfo.Remark != nil {
				groupInfoMap["remark"] = groupInfo.Remark
			}

			ids = append(ids, helper.UInt64ToStr(*groupInfo.GroupId))
			tmpList = append(tmpList, groupInfoMap)
		}

		_ = d.Set("group_info", tmpList)
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
