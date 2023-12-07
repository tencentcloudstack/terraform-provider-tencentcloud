package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSsmProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmProductsRead,
		Schema: map[string]*schema.Schema{
			"products": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of supported services.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSsmProductsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssm_products.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = SsmService{client: meta.(*TencentCloudClient).apiV3Conn}
		products []*string
	)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmProductsByFilter(ctx)
		if e != nil {
			return retryError(e)
		}

		products = result
		return nil
	})

	if err != nil {
		return err
	}

	if products != nil {
		_ = d.Set("products", products)
	}

	d.SetId(helper.StrListToStr(products))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), products); e != nil {
			return e
		}
	}

	return nil
}
