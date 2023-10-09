/*
Use this data source to query detailed information of rum sign

Example Usage

```hcl
data "tencentcloud_rum_sign" "sign" {
  timeout   = 1800
  file_type = 1
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRumSign() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumSignRead,
		Schema: map[string]*schema.Schema{
			"timeout": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Timeout duration.",
			},

			"file_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Bucket type. `1`:web project; `2`:app project.",
			},

			"secret_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Temporary access key.",
			},

			"secret_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Temporary access key ID.",
			},

			"session_token": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Temporary access key token.",
			},

			"start_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Start timestamp.",
			},

			"expired_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Expiration timestamp.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudRumSignRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rum_sign.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("timeout"); v != nil {
		paramMap["Timeout"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("file_type"); v != nil {
		paramMap["FileType"] = helper.IntInt64(v.(int))
	}

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *rum.DescribeReleaseFileSignResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeRumSignByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	var token string
	if result != nil {
		if result.SecretKey != nil {
			_ = d.Set("secret_key", result.SecretKey)
		}

		if result.SecretID != nil {
			_ = d.Set("secret_id", result.SecretID)
		}

		if result.SessionToken != nil {
			token = *result.SessionToken
			_ = d.Set("session_token", result.SessionToken)
		}

		if result.StartTime != nil {
			_ = d.Set("start_time", result.StartTime)
		}

		if result.ExpiredTime != nil {
			_ = d.Set("expired_time", result.ExpiredTime)
		}
	}

	d.SetId(helper.DataResourceIdsHash([]string{token}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
