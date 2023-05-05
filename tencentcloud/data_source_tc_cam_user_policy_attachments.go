/*
Use this data source to query detailed information of CAM user policy attachments

Example Usage

```hcl
# query by user_id
data "tencentcloud_cam_user_policy_attachments" "foo" {
  user_id = tencentcloud_cam_user.foo.id
}

# query by user_id and policy_id
data "tencentcloud_cam_user_policy_attachments" "bar" {
  user_id   = tencentcloud_cam_user.foo.id
  policy_id = tencentcloud_cam_policy.foo.id
}
```
*/
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

func dataSourceTencentCloudCamUserPolicyAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamUserPolicyAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"user_id", "user_name"},
				Deprecated:   "It has been deprecated from version 1.59.6. Use `user_name` instead.",
				Description:  "ID of the attached CAM user to be queried.",
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"user_id", "user_name"},
				Description:  "Name of the attached CAM user as unique key to be queried.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of CAM policy to be queried.",
			},
			"create_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{1, 2}),
				Description:  "Mode of Creation of the CAM user policy attachment. `1` means the CAM policy attachment is created by production, and the others indicate syntax strategy ways.",
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
			"user_policy_attachment_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CAM user policy attachments. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "It has been deprecated from version 1.59.6. Use `user_name` instead.",
							Description: "ID of CAM user.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CAM user as unique key.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CAM user.",
						},
						"create_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Mode of Creation of the CAM user policy attachment. `1` means the cam policy attachment is created by production, and the others indicate syntax strategy ways.",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time of the CAM user policy attachment.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the policy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamUserPolicyAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_user_policy_attachments.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	params := make(map[string]interface{})
	userId, _, err := getUserId(d)
	if err != nil {
		return err
	}
	params["user_id"] = userId
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
	var policyOfUsers []*cam.AttachPolicyInfo
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeUserPolicyAttachmentsByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		policyOfUsers = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM user policy attachments failed, reason:%s\n", logId, err.Error())
		return err
	}
	policyOfUserList := make([]map[string]interface{}, 0, len(policyOfUsers))
	ids := make([]string, 0, len(policyOfUsers))
	for _, policy := range policyOfUsers {
		mapping := map[string]interface{}{
			"user_id":     userId,
			"user_name":   userId,
			"policy_id":   strconv.Itoa(int(*policy.PolicyId)),
			"create_time": *policy.AddTime,
			"create_mode": *policy.CreateMode,
			"policy_type": *policy.PolicyType,
			"policy_name": *policy.PolicyName,
		}
		policyOfUserList = append(policyOfUserList, mapping)
		ids = append(ids, userId+"#"+strconv.Itoa(int(*policy.PolicyId)))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("user_policy_attachment_list", policyOfUserList); e != nil {
		log.Printf("[CRITAL]%s provider set CAM user polilcy attachment list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), policyOfUserList); e != nil {
			return e
		}
	}

	return nil
}
