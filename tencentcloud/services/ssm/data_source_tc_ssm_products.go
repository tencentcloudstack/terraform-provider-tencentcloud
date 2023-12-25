package ssm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSsmProducts() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_ssm_products.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		products []*string
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmProductsByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
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
		if e := tccommon.WriteToFile(output.(string), products); e != nil {
			return e
		}
	}

	return nil
}
