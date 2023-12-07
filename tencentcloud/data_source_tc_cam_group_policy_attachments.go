package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCamGroupPolicyAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamGroupPolicyAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the attached CAM group to be queried.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of CAM policy to be queried.",
			},
			"create_mode": {
				Type:     schema.TypeInt,
				Optional: true,

				Description: "Mode of creation of the CAM user policy attachment. 1 means the cam policy attachment is created by production, and the others indicate syntax strategy ways.",
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CAM_POLICY_CREATE_STRATEGY),
				Description:  "Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"group_policy_attachment_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CAM group policy attachments. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of CAM group.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CAM group.",
						},
						"create_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Mode of Creation of the CAM group policy attachment. 1 means the cam policy attachment is created by production, and the others indicate syntax strategy ways.",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the CAM group policy attachment.",
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

func dataSourceTencentCloudCamGroupPolicyAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_group_policy_attachments.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	params := make(map[string]interface{})
	groupId := d.Get("group_id").(string)
	params["group_id"] = groupId
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
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var policyOfGroups []*cam.AttachPolicyInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeGroupPolicyAttachmentsByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		policyOfGroups = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM group policy attachments failed, reason:%s\n", logId, err.Error())
		return err
	}
	policyOfGroupList := make([]map[string]interface{}, 0, len(policyOfGroups))
	ids := make([]string, 0, len(policyOfGroups))
	for _, policy := range policyOfGroups {
		mapping := map[string]interface{}{
			"group_id":    groupId,
			"policy_id":   strconv.Itoa(int(*policy.PolicyId)),
			"create_time": *policy.AddTime,
			"create_mode": *policy.CreateMode,
			"policy_type": *policy.PolicyType,
			"policy_name": *policy.PolicyName,
		}
		policyOfGroupList = append(policyOfGroupList, mapping)
		ids = append(ids, groupId+"#"+strconv.Itoa(int(*policy.PolicyId)))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("group_policy_attachment_list", policyOfGroupList); e != nil {
		log.Printf("[CRITAL]%s provider set group polilcy attachment list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), policyOfGroupList); e != nil {
			return e
		}
	}

	return nil
}
