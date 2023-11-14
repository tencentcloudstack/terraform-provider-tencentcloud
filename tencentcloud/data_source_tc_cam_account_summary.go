/*
Use this data source to query detailed information of cam account_summary

Example Usage

```hcl
data "tencentcloud_cam_account_summary" "account_summary" {
            }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCamAccountSummary() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamAccountSummaryRead,
		Schema: map[string]*schema.Schema{
			"policies": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of policy.",
			},

			"roles": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of role.",
			},

			"user": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of Sub-user.",
			},

			"group": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of Group.",
			},

			"member": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of grouped users.",
			},

			"identity_providers": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of identity provider.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCamAccountSummaryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_account_summary.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCamAccountSummaryByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		policies = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(policies))
	if policies != nil {
		_ = d.Set("policies", policies)
	}

	if roles != nil {
		_ = d.Set("roles", roles)
	}

	if user != nil {
		_ = d.Set("user", user)
	}

	if group != nil {
		_ = d.Set("group", group)
	}

	if member != nil {
		_ = d.Set("member", member)
	}

	if identityProviders != nil {
		_ = d.Set("identity_providers", identityProviders)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
