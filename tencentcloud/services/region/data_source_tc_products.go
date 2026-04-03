package region

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	regionv20220627 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudProductsRead,
		Schema: map[string]*schema.Schema{
			"product_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Product list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product name, e.g. `cvm`.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudProductsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_products.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = RegionService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})

	var respData []*regionv20220627.RegionProduct
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeProductsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	productListData := make([]map[string]interface{}, 0, len(respData))
	for _, product := range respData {
		productMap := map[string]interface{}{}
		if product.Name != nil {
			productMap["name"] = product.Name
		}
		productListData = append(productListData, productMap)
	}

	_ = d.Set("product_list", productListData)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
