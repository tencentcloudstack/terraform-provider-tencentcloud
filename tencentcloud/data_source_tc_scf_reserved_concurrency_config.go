/*
Use this data source to query detailed information of scf reserved_concurrency_config

Example Usage

```hcl
data "tencentcloud_scf_reserved_concurrency_config" "reserved_concurrency_config" {
  function_name = ""
  namespace = ""
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

func dataSourceTencentCloudScfReservedConcurrencyConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfReservedConcurrencyConfigRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Specifies the function of which you want to obtain the reserved quota.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace. Default value: default.",
			},

			"reserved_mem": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The reserved quota of the functionNote: this field may return `null`, indicating that no valid values can be obtained.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudScfReservedConcurrencyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_reserved_concurrency_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfReservedConcurrencyConfigByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		reservedMem = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(reservedMem))
	if reservedMem != nil {
		_ = d.Set("reserved_mem", reservedMem)
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
