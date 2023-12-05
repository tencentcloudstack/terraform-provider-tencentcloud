package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCamPolicyGrantingServiceAccess() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamPolicyGrantingServiceAccessRead,
		Schema: map[string]*schema.Schema{
			"target_uin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Sub-account uin, one of the three (TargetUin, RoleId, GroupId) must be passed.",
			},

			"role_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Role Id, one of the three (TargetUin, RoleId, GroupId) must be passed.",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Group Id, one of the three (TargetUin, RoleId, GroupId) must be passed.",
			},

			"service_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service type, this field needs to be passed when viewing the details of the service authorization interface.",
			},

			"list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Service info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service type.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service name.",
									},
								},
							},
						},
						"action": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Action list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action name.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action description.",
									},
								},
							},
						},
						"policy": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Policy list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy Id.",
									},
									"policy_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy name.",
									},
									"policy_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Polic type.",
									},
									"policy_description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy description.",
									},
								},
							},
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

func dataSourceTencentCloudCamPolicyGrantingServiceAccessRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_policy_granting_service_access.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		targetUin string
		roleId    string
		groupId   string
	)
	paramMap := make(map[string]interface{})
	if v, _ := d.GetOkExists("target_uin"); v != nil {
		paramMap["TargetUin"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOkExists("role_id"); v != nil {
		paramMap["RoleId"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOkExists("group_id"); v != nil {
		paramMap["GroupId"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("service_type"); ok {
		paramMap["ServiceType"] = helper.String(v.(string))
	}

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	var list []*cam.ListGrantServiceAccessNode

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCamPolicyGrantingServiceAccessByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		list = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(list))

	if list != nil {
		for _, listGrantServiceAccessNode := range list {
			listGrantServiceAccessNodeMap := map[string]interface{}{}

			if listGrantServiceAccessNode.Service != nil {
				serviceMap := map[string]interface{}{}

				if listGrantServiceAccessNode.Service.ServiceType != nil {
					serviceMap["service_type"] = listGrantServiceAccessNode.Service.ServiceType
				}

				if listGrantServiceAccessNode.Service.ServiceName != nil {
					serviceMap["service_name"] = listGrantServiceAccessNode.Service.ServiceName
				}

				listGrantServiceAccessNodeMap["service"] = []interface{}{serviceMap}
			}

			if listGrantServiceAccessNode.Action != nil {
				var actionList []interface{}
				for _, action := range listGrantServiceAccessNode.Action {
					actionMap := map[string]interface{}{}

					if action.Name != nil {
						actionMap["name"] = action.Name
					}

					if action.Description != nil {
						actionMap["description"] = action.Description
					}

					actionList = append(actionList, actionMap)
				}

				listGrantServiceAccessNodeMap["action"] = actionList
			}

			if listGrantServiceAccessNode.Policy != nil {
				var policyList []interface{}
				for _, policy := range listGrantServiceAccessNode.Policy {
					policyMap := map[string]interface{}{}

					if policy.PolicyId != nil {
						policyMap["policy_id"] = policy.PolicyId
					}

					if policy.PolicyName != nil {
						policyMap["policy_name"] = policy.PolicyName
					}

					if policy.PolicyType != nil {
						policyMap["policy_type"] = policy.PolicyType
					}

					if policy.PolicyDescription != nil {
						policyMap["policy_description"] = policy.PolicyDescription
					}

					policyList = append(policyList, policyMap)
				}

				listGrantServiceAccessNodeMap["policy"] = policyList
			}

			tmpList = append(tmpList, listGrantServiceAccessNodeMap)
		}

		_ = d.Set("list", tmpList)
	}

	d.SetId(targetUin + FILED_SP + roleId + groupId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
