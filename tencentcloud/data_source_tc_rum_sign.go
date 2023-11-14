/*
Use this data source to query detailed information of rum sign

Example Usage

```hcl
data "tencentcloud_rum_sign" "sign" {
  timeout = 1800
  file_type = web
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
				Description: "Bucket type. `web`:&amp;amp;#39;web project&amp;amp;#39;; `app`:&amp;amp;#39;app project&amp;amp;#39; .",
			},

			"secret_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Temporary access key.",
			},

			"secret_i_d": {
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

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRumSignByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		secretKey = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(secretKey))
	if secretKey != nil {
		_ = d.Set("secret_key", secretKey)
	}

	if secretID != nil {
		_ = d.Set("secret_i_d", secretID)
	}

	if sessionToken != nil {
		_ = d.Set("session_token", sessionToken)
	}

	if startTime != nil {
		_ = d.Set("start_time", startTime)
	}

	if expiredTime != nil {
		_ = d.Set("expired_time", expiredTime)
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
