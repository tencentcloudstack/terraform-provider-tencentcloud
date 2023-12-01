/*
Use this data source to query detailed information of cam list_attached_user_policy

Example Usage

```hcl
data "tencentcloud_cam_list_attached_user_policy" "list_attached_user_policy" {
  target_uin = 100032767426
  attach_type = 0
    }
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCamListAttachedUserPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamListAttachedUserPolicyRead,
		Schema: map[string]*schema.Schema{
			"target_uin": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Target User ID.",
			},

			"attach_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "0: Return direct association and group association policies, 1: Only return direct association policies, 2: Only return group association policies.",
			},

			"strategy_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Policy type.",
			},

			"keyword": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search Keywords.",
			},

			"policy_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Policy List Data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy ID.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy Name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy Description.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"strategy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy type (1 represents custom policy, 2 represents preset policy).",
						},
						"create_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation mode (1 represents policies created by product or project permissions, others represent policies created by policy syntax).",
						},
						"groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated information with groupNote: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Group ID.",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group Name.",
									},
								},
							},
						},
						"deactived": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Has it been taken offline (0: No 1: Yes)Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"deactived_detail": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of offline productsNote: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudCamListAttachedUserPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_list_attached_user_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOkExists("target_uin"); v != nil {
		paramMap["TargetUin"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOkExists("attach_type"); v != nil {
		paramMap["AttachType"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOkExists("strategy_type"); v != nil {
		paramMap["StrategyType"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	var policyList []*cam.AttachedUserPolicy

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCamListAttachedUserPolicyByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		policyList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(policyList))
	tmpList := make([]map[string]interface{}, 0, len(policyList))

	if policyList != nil {
		for _, attachedUserPolicy := range policyList {
			attachedUserPolicyMap := map[string]interface{}{}

			if attachedUserPolicy.PolicyId != nil {
				attachedUserPolicyMap["policy_id"] = attachedUserPolicy.PolicyId
			}

			if attachedUserPolicy.PolicyName != nil {
				attachedUserPolicyMap["policy_name"] = attachedUserPolicy.PolicyName
			}

			if attachedUserPolicy.Description != nil {
				attachedUserPolicyMap["description"] = attachedUserPolicy.Description
			}

			if attachedUserPolicy.AddTime != nil {
				attachedUserPolicyMap["add_time"] = attachedUserPolicy.AddTime
			}

			if attachedUserPolicy.StrategyType != nil {
				attachedUserPolicyMap["strategy_type"] = attachedUserPolicy.StrategyType
			}

			if attachedUserPolicy.CreateMode != nil {
				attachedUserPolicyMap["create_mode"] = attachedUserPolicy.CreateMode
			}

			if attachedUserPolicy.Groups != nil {
				groupsList := []interface{}{}
				for _, groups := range attachedUserPolicy.Groups {
					groupsMap := map[string]interface{}{}

					if groups.GroupId != nil {
						groupsMap["group_id"] = groups.GroupId
					}

					if groups.GroupName != nil {
						groupsMap["group_name"] = groups.GroupName
					}

					groupsList = append(groupsList, groupsMap)
				}

				attachedUserPolicyMap["groups"] = groupsList
			}

			if attachedUserPolicy.Deactived != nil {
				attachedUserPolicyMap["deactived"] = attachedUserPolicy.Deactived
			}

			if attachedUserPolicy.DeactivedDetail != nil {
				attachedUserPolicyMap["deactived_detail"] = attachedUserPolicy.DeactivedDetail
			}

			ids = append(ids, *attachedUserPolicy.PolicyId)
			tmpList = append(tmpList, attachedUserPolicyMap)
		}

		_ = d.Set("policy_list", tmpList)
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
