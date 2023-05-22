/*
Use this data source to query detailed information of cvm image_quota

Example Usage

```hcl
data "tencentcloud_cvm_image_quota" "image_quota" {
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

func dataSourceTencentCloudCvmImageQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmImageQuotaRead,
		Schema: map[string]*schema.Schema{
			"image_num_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The image quota of an account.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCvmImageQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cvm_image_quota.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var imageNumQuota int64
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmImageQuotaByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		imageNumQuota = result
		return nil
	})
	if err != nil {
		return err
	}

	_ = d.Set("image_num_quota", imageNumQuota)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), map[string]interface{}{
			"image_num_quota": imageNumQuota,
		}); e != nil {
			return e
		}
	}
	return nil
}
