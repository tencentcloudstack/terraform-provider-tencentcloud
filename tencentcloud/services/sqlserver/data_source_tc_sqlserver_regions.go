package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverRegionsRead,
		Schema: map[string]*schema.Schema{
			"region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Region information array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID in the format of ap-guangzhou.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Numeric ID of region.",
						},
						"region_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current purchasability of this region. UNAVAILABLE: not purchasable, AVAILABLE: purchasable.",
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

func dataSourceTencentCloudSqlserverRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_regions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		regionSet []*sqlserver.RegionInfo
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverDatasourceRegionsByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}

		regionSet = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(regionSet))
	tmpList := make([]map[string]interface{}, 0, len(regionSet))

	if regionSet != nil {
		for _, regionInfo := range regionSet {
			regionInfoMap := map[string]interface{}{}

			if regionInfo.Region != nil {
				regionInfoMap["region"] = regionInfo.Region
			}

			if regionInfo.RegionName != nil {
				regionInfoMap["region_name"] = regionInfo.RegionName
			}

			if regionInfo.RegionId != nil {
				regionInfoMap["region_id"] = regionInfo.RegionId
			}

			if regionInfo.RegionState != nil {
				regionInfoMap["region_state"] = regionInfo.RegionState
			}

			ids = append(ids, *regionInfo.Region)
			tmpList = append(tmpList, regionInfoMap)
		}

		_ = d.Set("region_set", tmpList)
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
