package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudConfigResourceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudConfigResourceTypesRead,
		Schema: map[string]*schema.Schema{
			"resource_type_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Supported resource type list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product code (e.g. CAM).",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product name.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type identifier (e.g. QCS::CAM::Group).",
						},
						"resource_type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type name.",
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

func dataSourceTencentCloudConfigResourceTypesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_config_resource_types.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, reqErr := service.DescribeConfigResourceTypes(ctx)
	if reqErr != nil {
		return reqErr
	}

	typeList := flattenConfigResourceTypeList(respData)
	_ = d.Set("resource_type_list", typeList)

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

func flattenConfigResourceTypeList(items []*configv20220802.ConfigResource) []map[string]interface{} {
	typeList := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		itemMap := map[string]interface{}{}

		if item.Product != nil {
			itemMap["product"] = item.Product
		}

		if item.ProductName != nil {
			itemMap["product_name"] = item.ProductName
		}

		if item.ResourceType != nil {
			itemMap["resource_type"] = item.ResourceType
		}

		if item.ResourceTypeName != nil {
			itemMap["resource_type_name"] = item.ResourceTypeName
		}

		typeList = append(typeList, itemMap)
	}

	return typeList
}
