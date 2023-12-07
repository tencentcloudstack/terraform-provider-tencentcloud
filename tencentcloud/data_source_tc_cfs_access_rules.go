package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCfsAccessRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfsAccessRulesRead,

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A specified access group ID used to query.",
			},
			"access_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A specified access rule ID used to query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"access_rule_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of CFS access rule. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the access rule.",
						},
						"auth_client_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Allowed IP of the access rule.",
						},
						"rw_permission": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Read and write permissions.",
						},
						"user_permission": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The permissions of accessing users.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The priority level of access rule.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCfsAccessRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cfs_access_rules.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var accessRuleId string
	accessGroupId := d.Get("access_group_id").(string)
	if v, ok := d.GetOk("access_rule_id"); ok {
		accessRuleId = v.(string)
	}

	var accessRules []*cfs.PGroupRuleInfo
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		accessRules, errRet = cfsService.DescribeAccessRule(ctx, accessGroupId, accessRuleId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	accessRuleList := make([]map[string]interface{}, 0, len(accessRules))
	ids := make([]string, 0, len(accessRules))
	for _, accessRule := range accessRules {
		mapping := map[string]interface{}{
			"access_rule_id":  accessRule.RuleId,
			"auth_client_ip":  accessRule.AuthClientIp,
			"rw_permission":   accessRule.RWPermission,
			"user_permission": accessRule.UserPermission,
			"priority":        accessRule.Priority,
		}
		accessRuleList = append(accessRuleList, mapping)
		ids = append(ids, *accessRule.RuleId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("access_rule_list", accessRuleList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set cfs access rule list fail, reason:%s\n ", logId, err.Error())
		return err
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), accessRuleList); err != nil {
			return err
		}
	}
	return nil
}
