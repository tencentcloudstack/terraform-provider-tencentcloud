/*
Use this data source to query detailed information of ssm describe_rotation_history

Example Usage

```hcl
data "tencentcloud_ssm_describe_rotation_history" "describe_rotation_history" {
  secret_name = ""
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

func dataSourceTencentCloudSsmDescribeRotationHistory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmDescribeRotationHistoryRead,
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Secret name.",
			},

			"version_i_ds": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The number of ersion numbers. The maximum number of version numbers that can be displayed to users is 10.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSsmDescribeRotationHistoryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssm_describe_rotation_history.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("secret_name"); ok {
		paramMap["SecretName"] = helper.String(v.(string))
	}

	service := SsmService{client: meta.(*TencentCloudClient).apiV3Conn}

	var versionIDs []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmDescribeRotationHistoryByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		versionIDs = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(versionIDs))
	if versionIDs != nil {
		_ = d.Set("version_i_ds", versionIDs)
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
