/*
Use this data source to query detailed information of ssm get_service_status

Example Usage

```hcl
data "tencentcloud_ssm_get_service_status" "get_service_status" {
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

func dataSourceTencentCloudSsmGetServiceStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmGetServiceStatusRead,
		Schema: map[string]*schema.Schema{
			"service_enabled": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "True means the service has been activated, false means the service has not been activated yet.",
			},

			"invalid_type": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Service unavailability type: 0-Not purchased, 1-Normal, 2-Service suspended due to arrears, 3-Resource release.",
			},

			"access_key_escrow_enabled": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "True means that the user can already use the key safe hosting function, false means that the user cannot use the key safe hosting function temporarily.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSsmGetServiceStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssm_get_service_status.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := SsmService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmGetServiceStatusByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		serviceEnabled = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(serviceEnabled))
	if serviceEnabled != nil {
		_ = d.Set("service_enabled", serviceEnabled)
	}

	if invalidType != nil {
		_ = d.Set("invalid_type", invalidType)
	}

	if accessKeyEscrowEnabled != nil {
		_ = d.Set("access_key_escrow_enabled", accessKeyEscrowEnabled)
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
