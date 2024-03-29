package cam

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCamRolePolicyAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamRolePolicyAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the attached CAM role to be queried.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of CAM policy to be queried.",
			},
			"create_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2}),
				Description:  "Mode of Creation of the CAM user policy attachment. `1` means the cam policy attachment is created by production, and the others indicate syntax strategy ways.",
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CAM_POLICY_CREATE_STRATEGY),
				Description:  "Type of the policy strategy. Valid values are 'User', 'QCS'. 'User' means customer strategy and 'QCS' means preset strategy.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"role_policy_attachment_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CAM role policy attachments. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of CAM role.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CAM role.",
						},
						"create_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Mode of Creation of the CAM role policy attachment. `1` means the cam policy attachment is created by production, and the others indicate syntax strategy ways.",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the CAM role policy attachment.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the policy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamRolePolicyAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cam_role_policy_attachments.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]interface{})
	roleId := d.Get("role_id").(string)
	params["role_id"] = roleId
	if v, ok := d.GetOk("policy_id"); ok {
		policyId, err := strconv.Atoi(v.(string))
		if err != nil {
			return err
		}
		params["policy_id"] = uint64(policyId)
	}
	if v, ok := d.GetOk("policy_type"); ok {
		params["policy_type"] = v.(string)
	}
	if v, ok := d.GetOk("create_mode"); ok {
		params["create_mode"] = v.(int)
	}

	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var policyOfRoles []*cam.AttachedPolicyOfRole
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeRolePolicyAttachmentsByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		policyOfRoles = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM role policy attachments failed, reason:%s\n", logId, err.Error())
		return err
	}
	policyOfRoleList := make([]map[string]interface{}, 0, len(policyOfRoles))
	ids := make([]string, 0, len(policyOfRoles))
	for _, policy := range policyOfRoles {
		mapping := map[string]interface{}{
			"role_id":     roleId,
			"policy_id":   strconv.Itoa(int(*policy.PolicyId)),
			"create_time": *policy.AddTime,
			"create_mode": *policy.CreateMode,
			"policy_type": *policy.PolicyType,
			"policy_name": *policy.PolicyName,
		}
		policyOfRoleList = append(policyOfRoleList, mapping)
		ids = append(ids, roleId+"#"+strconv.Itoa(int(*policy.PolicyId)))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("role_policy_attachment_list", policyOfRoleList); e != nil {
		log.Printf("[CRITAL]%s provider set role polilcy attachment list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), policyOfRoleList); e != nil {
			return e
		}
	}

	return nil
}
