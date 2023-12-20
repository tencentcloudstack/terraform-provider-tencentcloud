package ckafka

import (
	"errors"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaRegionRead,
		Schema: map[string]*schema.Schema{
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Return a list of region enumeration results.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "region ID.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "geographical name.",
						},
						"area_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "area name.",
						},
						"region_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region Code.",
						},
						"region_code_v3": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region Code(V3 version).",
						},
						"support": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NONE: The default value does not support any special models CVM: Supports CVM types.",
						},
						"ipv6": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to support ipv6, 0: means not supported, 1: means supported.",
						},
						"multi_zone": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to support cross-availability zones, 0: means not supported, 1: means supported.",
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

func dataSourceTencentCloudCkafkaRegionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_region.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var result []*ckafka.Region
	request := ckafka.NewDescribeRegionRequest()
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().DescribeRegion(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if response == nil || response.Response == nil {
			return tccommon.RetryError(errors.New("Response is null"))
		}
		result = response.Response.Result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	tmpList := make([]map[string]interface{}, 0, len(result))

	if result != nil {
		for _, region := range result {
			regionMap := map[string]interface{}{}

			if region.RegionId != nil {
				regionMap["region_id"] = region.RegionId
			}

			if region.RegionName != nil {
				regionMap["region_name"] = region.RegionName
			}

			if region.AreaName != nil {
				regionMap["area_name"] = region.AreaName
			}

			if region.RegionCode != nil {
				regionMap["region_code"] = region.RegionCode
			}

			if region.RegionCodeV3 != nil {
				regionMap["region_code_v3"] = region.RegionCodeV3
			}

			if region.Support != nil {
				regionMap["support"] = region.Support
			}

			if region.Ipv6 != nil {
				regionMap["ipv6"] = region.Ipv6
			}

			if region.MultiZone != nil {
				regionMap["multi_zone"] = region.MultiZone
			}

			ids = append(ids, strconv.FormatInt(*region.RegionId, 10))
			tmpList = append(tmpList, regionMap)
		}

		_ = d.Set("result", tmpList)
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
