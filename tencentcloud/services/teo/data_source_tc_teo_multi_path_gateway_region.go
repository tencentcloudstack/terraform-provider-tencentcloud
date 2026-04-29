package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoMultiPathGatewayRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoMultiPathGatewayRegionRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"gateway_regions": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of available gateway regions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID.",
						},
						"cn_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chinese name of the region.",
						},
						"en_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "English name of the region.",
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

func dataSourceTencentCloudTeoMultiPathGatewayRegionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_multi_path_gateway_region.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	gatewayRegions, err := service.DescribeTeoMultiPathGatewayRegionByFilter(ctx, paramMap)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(gatewayRegions))
	tmpList := make([]map[string]interface{}, 0, len(gatewayRegions))
	if gatewayRegions != nil {
		for _, region := range gatewayRegions {
			regionMap := map[string]interface{}{}
			if region.RegionId != nil {
				regionMap["region_id"] = region.RegionId
				ids = append(ids, *region.RegionId)
			}
			if region.CNName != nil {
				regionMap["cn_name"] = region.CNName
			}
			if region.ENName != nil {
				regionMap["en_name"] = region.ENName
			}
			tmpList = append(tmpList, regionMap)
		}
		_ = d.Set("gateway_regions", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
