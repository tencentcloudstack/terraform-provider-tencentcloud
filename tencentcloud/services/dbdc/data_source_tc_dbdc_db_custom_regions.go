package dbdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbdcDbCustomRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomRegionsRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"region_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DB Custom supported region list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name.",
						},
						"region_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sale status. Values: SELL (normal sale), SOLD_OUT (sold out).",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDbdcDbCustomRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_regions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	var respData []*dbdcv20201029.RegionInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDBCustomRegions(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	regionSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, region := range respData {
			regionMap := map[string]interface{}{}
			if region.Region != nil {
				regionMap["region"] = region.Region
			}

			if region.RegionState != nil {
				regionMap["region_state"] = region.RegionState
			}

			regionSetList = append(regionSetList, regionMap)
		}

		_ = d.Set("region_set", regionSetList)
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
