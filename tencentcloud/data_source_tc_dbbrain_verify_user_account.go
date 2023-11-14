/*
Use this data source to query detailed information of dbbrain verify_user_account

Example Usage

```hcl
data "tencentcloud_dbbrain_verify_user_account" "verify_user_account" {
  instance_id = ""
  user = ""
  password = ""
  product = ""
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

func dataSourceTencentCloudDbbrainVerifyUserAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainVerifyUserAccountRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"user": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database account name.",
			},

			"password": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database account password.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values： mysql - cloud database MySQL; cynosdb - cloud database TDSQL-C for MySQL, the default is mysql.",
			},

			"session_token": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The session token is valid for 5 minutes.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainVerifyUserAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_verify_user_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user"); ok {
		paramMap["User"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		paramMap["Password"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainVerifyUserAccountByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		sessionToken = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(sessionToken))
	if sessionToken != nil {
		_ = d.Set("session_token", sessionToken)
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
