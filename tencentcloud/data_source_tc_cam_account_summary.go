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
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
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

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	AccountData := &cam.GetAccountSummaryResponseParams{}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCamAccountSummaryByFilter(ctx)
		if e != nil {
			return retryError(e)
		}
		AccountData = result
		return nil
	})
	if err != nil {
		return err
	}
	template := make(map[string]interface{}, 0)

	if AccountData.Policies != nil {
		_ = d.Set("policies", AccountData.Policies)
		template["policies"] = AccountData.Policies
	}

	if AccountData.Roles != nil {
		_ = d.Set("roles", AccountData.Roles)
		template["roles"] = AccountData.Roles
	}

	if AccountData.User != nil {
		_ = d.Set("user", AccountData.User)
		template["user"] = AccountData.User
	}

	if AccountData.Group != nil {
		_ = d.Set("group", AccountData.Group)
		template["group"] = AccountData.Group
	}

	if AccountData.Member != nil {
		_ = d.Set("member", AccountData.Member)
		template["member"] = AccountData.Member
	}

	if AccountData.IdentityProviders != nil {
		_ = d.Set("identity_providers", AccountData.IdentityProviders)
		template["identity_providers"] = AccountData.IdentityProviders
	}
	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), template); e != nil {
			return e
		}
	}
	return nil
}
