/*
Use this data source to query detailed information of dlc describe_user_type

Example Usage

```hcl
data "tencentcloud_dlc_describe_user_type" "describe_user_type" {
  user_id = "127382378"
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

func dataSourceTencentCloudDlcDescribeUserType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeUserTypeRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User id (uin), if left blank, it defaults to the caller's sub-uin.",
			},

			"user_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "User type, only support: ADMIN: ddministrator/COMMON: ordinary user.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDlcDescribeUserTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dlc_describe_user_type.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	var userId string
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		paramMap["UserId"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}
	var userType *string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeUserTypeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		userType = result
		return nil
	})
	if err != nil {
		return err
	}

	if userType != nil {
		_ = d.Set("user_type", userType)
	}

	d.SetId(userId + FILED_SP + *userType)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), userType); e != nil {
			return e
		}
	}
	return nil
}
