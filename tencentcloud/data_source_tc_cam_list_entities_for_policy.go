/*
Use this data source to query detailed information of cam list_entities_for_policy

Example Usage

```hcl
data "tencentcloud_cam_list_entities_for_policy" "list_entities_for_policy" {
  policy_id = 1
  entity_filter = "All"
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

func dataSourceTencentCloudCamListEntitiesForPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamListEntitiesForPolicyRead,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Policy Id.",
			},

			"rp": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Per page size, default value is 20.",
			},

			"entity_filter": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Can take values of &amp;amp;#39;All&amp;amp;#39;, &amp;amp;#39;User&amp;amp;#39;, &amp;amp;#39;Group&amp;amp;#39;, and &amp;amp;#39;Role&amp;amp;#39;. &amp;amp;#39;All&amp;amp;#39; represents obtaining all entity types, &amp;amp;#39;User&amp;amp;#39; represents only obtaining sub accounts, &amp;amp;#39;Group&amp;amp;#39; represents only obtaining user groups, and &amp;amp;#39;Role&amp;amp;#39; represents only obtaining roles. The default value is&amp;amp;#39; All &amp;amp;#39;.",
			},

			"list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Entity ListNote: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entity ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entity NameNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Entity UinNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"related_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Association type. 1. User association; 2 User Group Association.",
						},
						"attachment_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy association timeNote: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudCamListEntitiesForPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_list_entities_for_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOkExists("policy_id"); v != nil {
		paramMap["PolicyId"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOkExists("rp"); v != nil {
		paramMap["Rp"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("entity_filter"); ok {
		paramMap["EntityFilter"] = helper.String(v.(string))
	}

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	var listEntitiesForPolicy []*cam.AttachEntityOfPolicy
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCamListEntitiesForPolicyByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		listEntitiesForPolicy = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(listEntitiesForPolicy))
	tmpList := make([]map[string]interface{}, 0)

	if listEntitiesForPolicy != nil {
		for _, attachEntityOfPolicy := range listEntitiesForPolicy {
			attachEntityOfPolicyMap := map[string]interface{}{}

			if attachEntityOfPolicy.Id != nil {
				attachEntityOfPolicyMap["id"] = attachEntityOfPolicy.Id
			}

			if attachEntityOfPolicy.Name != nil {
				attachEntityOfPolicyMap["name"] = attachEntityOfPolicy.Name
			}

			if attachEntityOfPolicy.Uin != nil {
				attachEntityOfPolicyMap["uin"] = attachEntityOfPolicy.Uin
			}

			if attachEntityOfPolicy.RelatedType != nil {
				attachEntityOfPolicyMap["related_type"] = attachEntityOfPolicy.RelatedType
			}

			if attachEntityOfPolicy.AttachmentTime != nil {
				attachEntityOfPolicyMap["attachment_time"] = attachEntityOfPolicy.AttachmentTime
			}

			ids = append(ids, helper.UInt64ToStr(*attachEntityOfPolicy.Uin))
			tmpList = append(tmpList, attachEntityOfPolicyMap)
		}

		_ = d.Set("list", tmpList)
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
