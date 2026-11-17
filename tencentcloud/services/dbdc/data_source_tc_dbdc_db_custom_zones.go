package dbdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbdcDbCustomZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomZonesRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"zone_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DB Custom available zone list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available zone, such as `ap-guangzhou-3`.",
						},
						"zone_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone sale status. Values: `SELL` (normal sale), `SOLD_OUT` (sold out).",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDbdcDbCustomZonesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_zones.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})

	var respData []*dbdcv20201029.ZoneInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, _, e := service.DescribeDBCustomZonesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	zoneSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, zone := range respData {
			zoneMap := map[string]interface{}{}
			if zone.Zone != nil {
				zoneMap["zone"] = zone.Zone
			}

			if zone.ZoneState != nil {
				zoneMap["zone_state"] = zone.ZoneState
			}

			zoneSetList = append(zoneSetList, zoneMap)
		}

		_ = d.Set("zone_set", zoneSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
