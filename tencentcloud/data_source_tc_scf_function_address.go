/*
Use this data source to query detailed information of scf function_address

Example Usage

```hcl
data "tencentcloud_scf_function_address" "function_address" {
  function_name = ""
  qualifier = ""
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

func dataSourceTencentCloudScfFunctionAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfFunctionAddressRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function version.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace.",
			},

			"url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cos address of the function.",
			},

			"code_sha256": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "SHA256 code of the function.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudScfFunctionAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_function_address.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		paramMap["Qualifier"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfFunctionAddressByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		url = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(url))
	if url != nil {
		_ = d.Set("url", url)
	}

	if codeSha256 != nil {
		_ = d.Set("code_sha256", codeSha256)
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
