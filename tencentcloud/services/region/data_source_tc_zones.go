package region

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	regionv20220627 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudZonesRead,
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

			"zone_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Zone list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone name, e.g. `ap-guangzhou-3`.",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone description, e.g. `Guangzhou Zone 3`.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone ID.",
						},
						"zone_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone status, `AVAILABLE` or `UNAVAILABLE`.",
						},
						"parent_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent zone identifier.",
						},
						"parent_zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent zone ID.",
						},
						"parent_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent zone description.",
						},
						"zone_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone type.",
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

func dataSourceTencentCloudZonesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_zones.read")()
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

	var respData []*regionv20220627.ZoneInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeZonesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	zoneListData := make([]map[string]interface{}, 0, len(respData))
	for _, zoneInfo := range respData {
		zoneMap := map[string]interface{}{}
		if zoneInfo.Zone != nil {
			zoneMap["zone"] = zoneInfo.Zone
		}
		if zoneInfo.ZoneName != nil {
			zoneMap["zone_name"] = zoneInfo.ZoneName
		}
		if zoneInfo.ZoneId != nil {
			zoneMap["zone_id"] = zoneInfo.ZoneId
		}
		if zoneInfo.ZoneState != nil {
			zoneMap["zone_state"] = zoneInfo.ZoneState
		}
		if zoneInfo.ParentZone != nil {
			zoneMap["parent_zone"] = zoneInfo.ParentZone
		}
		if zoneInfo.ParentZoneId != nil {
			zoneMap["parent_zone_id"] = zoneInfo.ParentZoneId
		}
		if zoneInfo.ParentZoneName != nil {
			zoneMap["parent_zone_name"] = zoneInfo.ParentZoneName
		}
		if zoneInfo.ZoneType != nil {
			zoneMap["zone_type"] = zoneInfo.ZoneType
		}
		zoneListData = append(zoneListData, zoneMap)
	}

	_ = d.Set("zone_list", zoneListData)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
