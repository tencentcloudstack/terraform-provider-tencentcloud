package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudEbPlatformProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbPlatformProductsRead,
		Schema: map[string]*schema.Schema{
			"platform_products": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Platform product list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform product name.",
						},
						"product_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform product type.",
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

func dataSourceTencentCloudEbPlatformProductsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eb_platform_products.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var platformProducts []*eb.PlatformProduct
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEbPlatformProductsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		platformProducts = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(platformProducts))
	tmpList := make([]map[string]interface{}, 0, len(platformProducts))

	if platformProducts != nil {
		for _, platformProduct := range platformProducts {
			platformProductMap := map[string]interface{}{}

			if platformProduct.ProductName != nil {
				platformProductMap["product_name"] = platformProduct.ProductName
			}

			if platformProduct.ProductType != nil {
				platformProductMap["product_type"] = platformProduct.ProductType
			}

			ids = append(ids, *platformProduct.ProductName)
			tmpList = append(tmpList, platformProductMap)
		}

		_ = d.Set("platform_products", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
