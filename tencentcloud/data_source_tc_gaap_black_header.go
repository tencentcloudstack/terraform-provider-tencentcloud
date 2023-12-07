package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapBlackHeader() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapBlackHeaderRead,
		Schema: map[string]*schema.Schema{
			"black_headers": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Disabled custom header listNote: This field may return null, indicating that a valid value cannot be obtained.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudGaapBlackHeaderRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_black_header.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	var blackHeaders []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapBlackHeader(ctx)
		if e != nil {
			return retryError(e)
		}
		blackHeaders = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(blackHeaders))
	if blackHeaders != nil {
		_ = d.Set("black_headers", blackHeaders)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), blackHeaders); e != nil {
			return e
		}
	}
	return nil
}
