package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCvmChcDeniedActions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmChcDeniedActionsRead,
		Schema: map[string]*schema.Schema{
			"chc_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "CHC host IDs.",
			},

			"chc_host_denied_action_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Actions not allowed for the CHC instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CHC instance ID.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CHC instance status.",
						},
						"deny_actions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Actions not allowed for the current CHC instance.",
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

func dataSourceTencentCloudCvmChcDeniedActionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cvm_chc_denied_actions.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("chc_ids"); ok {
		chcIdsSet := v.(*schema.Set).List()
		paramMap["chc_ids"] = helper.InterfacesStrings(chcIdsSet)
	}

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	var chcHostDeniedActionSet []*cvm.ChcHostDeniedActions

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmChcDeniedActionsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		chcHostDeniedActionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(chcHostDeniedActionSet))
	tmpList := make([]map[string]interface{}, 0, len(chcHostDeniedActionSet))

	if len(chcHostDeniedActionSet) > 0 {
		for _, chcHostDeniedActions := range chcHostDeniedActionSet {
			chcHostDeniedActionsMap := map[string]interface{}{}

			if chcHostDeniedActions.ChcId != nil {
				chcHostDeniedActionsMap["chc_id"] = chcHostDeniedActions.ChcId
			}

			if chcHostDeniedActions.State != nil {
				chcHostDeniedActionsMap["state"] = chcHostDeniedActions.State
			}

			if chcHostDeniedActions.DenyActions != nil {
				chcHostDeniedActionsMap["deny_actions"] = chcHostDeniedActions.DenyActions
			}

			ids = append(ids, *chcHostDeniedActions.ChcId)
			tmpList = append(tmpList, chcHostDeniedActionsMap)
		}

		_ = d.Set("chc_host_denied_action_set", tmpList)
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
