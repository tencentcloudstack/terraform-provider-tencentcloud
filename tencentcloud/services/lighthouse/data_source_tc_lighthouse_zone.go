package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseZoneRead,
		Schema: map[string]*schema.Schema{
			"order_field": {
				Optional: true,
				Type:     schema.TypeString,
				Description: "Sorting field. Valid values:\n" +
					"- ZONE: Sort by the availability zone.\n" +
					"- INSTANCE_DISPLAY_LABEL: Sort by visibility labels (HIDDEN, NORMAL and SELECTED). Default: [HIDDEN, NORMAL, SELECTED].\n" +
					"Sort by availability zone by default.",
			},

			"order": {
				Optional: true,
				Type:     schema.TypeString,
				Description: "Specifies how availability zones are listed. Valid values:\n" +
					"- ASC: Ascending sort.\n" +
					"- DESC: Descending sort.\n" +
					"The default value is ASC.",
			},
			"zone_info_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of zone info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chinese name of availability zone.",
						},
						"instance_display_label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance purchase page availability zone display label.",
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

func dataSourceTencentCloudLighthouseZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_zone.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("order_field"); ok {
		paramMap["order_field"] = v.(string)
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["order"] = v.(string)
	}

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var zoneInfoSet []*lighthouse.ZoneInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseZoneByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		zoneInfoSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(zoneInfoSet))
	tmpList := make([]map[string]interface{}, 0, len(zoneInfoSet))
	for _, zoneInfo := range zoneInfoSet {
		ids = append(ids, *zoneInfo.Zone)
		tmpList = append(tmpList, map[string]interface{}{
			"zone":                   *zoneInfo.Zone,
			"zone_name":              *zoneInfo.ZoneName,
			"instance_display_label": *zoneInfo.InstanceDisplayLabel,
		})
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("zone_info_set", tmpList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
