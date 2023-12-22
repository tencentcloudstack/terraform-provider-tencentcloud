package postgresql

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlZonesRead,
		Schema: map[string]*schema.Schema{
			"zone_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "AZ information set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AZ abbreviation.",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AZ name.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "AZ number.",
						},
						"zone_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability status. Valid values:`UNAVAILABLE`.`AVAILABLE`.`SELLOUT`.`SUPPORTMODIFYONLY` (supports configuration adjustment).",
						},
						"zone_support_ipv6": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the AZ supports IPv6 address access.",
						},
						"standby_zone_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "AZs that can be used as standby when this AZ is primaryNote: this field may return `null`, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudPostgresqlZonesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_zones.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var zoneSet []*postgresql.ZoneInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlZonesByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		zoneSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(zoneSet))
	tmpList := make([]map[string]interface{}, 0, len(zoneSet))

	if zoneSet != nil {
		for _, zoneInfo := range zoneSet {
			zoneInfoMap := map[string]interface{}{}

			if zoneInfo.Zone != nil {
				zoneInfoMap["zone"] = zoneInfo.Zone
			}

			if zoneInfo.ZoneName != nil {
				zoneInfoMap["zone_name"] = zoneInfo.ZoneName
			}

			if zoneInfo.ZoneId != nil {
				zoneInfoMap["zone_id"] = zoneInfo.ZoneId
			}

			if zoneInfo.ZoneState != nil {
				zoneInfoMap["zone_state"] = zoneInfo.ZoneState
			}

			if zoneInfo.ZoneSupportIpv6 != nil {
				zoneInfoMap["zone_support_ipv6"] = zoneInfo.ZoneSupportIpv6
			}

			if zoneInfo.StandbyZoneSet != nil {
				zoneInfoMap["standby_zone_set"] = helper.StringsInterfaces(zoneInfo.StandbyZoneSet)
			}

			ids = append(ids, *zoneInfo.Zone)
			tmpList = append(tmpList, zoneInfoMap)
		}

		_ = d.Set("zone_set", tmpList)
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
