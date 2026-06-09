package region

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	regionv20220627 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRegionsRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Product name to query, e.g. `cvm`. Use `tencentcloud_products` to get available product names.",
			},

			"scene": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Scene control parameter. `0` or not set means do not query optional business whitelist; `1` means query optional business whitelist.",
			},

			"region_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Region list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region identifier, e.g. `ap-guangzhou`.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name, e.g. `South China (Guangzhou)`.",
						},
						"region_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region availability status, e.g. `AVAILABLE`.",
						},
						"region_type_m_c": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Console type, null when called via API.",
						},
						"location_m_c": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region description in different languages.",
						},
						"region_name_m_c": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region description displayed in console.",
						},
						"region_id_m_c": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID for console.",
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

func dataSourceTencentCloudRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_regions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = RegionService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("scene"); ok {
		paramMap["Scene"] = helper.IntInt64(v.(int))
	}

	var respData []*regionv20220627.RegionInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRegionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	regionListData := make([]map[string]interface{}, 0, len(respData))
	for _, regionInfo := range respData {
		regionMap := map[string]interface{}{}
		if regionInfo.Region != nil {
			regionMap["region"] = regionInfo.Region
		}
		if regionInfo.RegionName != nil {
			regionMap["region_name"] = regionInfo.RegionName
		}
		if regionInfo.RegionState != nil {
			regionMap["region_state"] = regionInfo.RegionState
		}
		if regionInfo.RegionTypeMC != nil {
			regionMap["region_type_m_c"] = regionInfo.RegionTypeMC
		}
		if regionInfo.LocationMC != nil {
			regionMap["location_m_c"] = regionInfo.LocationMC
		}
		if regionInfo.RegionNameMC != nil {
			regionMap["region_name_m_c"] = regionInfo.RegionNameMC
		}
		if regionInfo.RegionIdMC != nil {
			regionMap["region_id_m_c"] = regionInfo.RegionIdMC
		}
		regionListData = append(regionListData, regionMap)
	}

	_ = d.Set("region_list", regionListData)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
