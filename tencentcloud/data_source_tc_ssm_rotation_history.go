/*
Use this data source to query detailed information of ssm rotation_history

Example Usage

```hcl
data "tencentcloud_ssm_rotation_history" "example" {
  secret_name = "keep_terraform"
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

func dataSourceTencentCloudSsmRotationHistory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmRotationHistoryRead,
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Secret name.",
			},
			"version_ids": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The number of version numbers. The maximum number of version numbers that can be displayed to users is 10.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSsmRotationHistoryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssm_rotation_history.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SsmService{client: meta.(*TencentCloudClient).apiV3Conn}
		versionIDs []*string
		secretName string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("secret_name"); ok {
		paramMap["SecretName"] = helper.String(v.(string))
		secretName = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmRotationHistoryByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		versionIDs = result
		return nil
	})

	if err != nil {
		return err
	}

	if versionIDs != nil {
		_ = d.Set("version_ids", versionIDs)
	}

	d.SetId(secretName)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
