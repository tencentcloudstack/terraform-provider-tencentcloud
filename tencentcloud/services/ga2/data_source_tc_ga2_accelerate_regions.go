package ga2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGa2AccelerateRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGa2AccelerateRegionsRead,
		Schema: map[string]*schema.Schema{
			"accelerator_region_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Accelerate region list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region Chinese name.",
						},
						"is_available": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Availability status. 0: unavailable, 1: available.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region identifier.",
						},
						"area_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Area name.",
						},
						"is_china_mainland": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is a China mainland region.",
						},
						"support_isp_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Supported ISP types.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"is_tencent_region": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is a Tencent region.",
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

func dataSourceTencentCloudGa2AccelerateRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ga2_accelerate_regions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(nil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	var respData []*ga2v20250115.AcceleratorRegionSet
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		request := ga2v20250115.NewDescribeAccelerateRegionsRequest()
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DescribeAccelerateRegionsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeAccelerateRegions failed, Response is nil"))
		}

		respData = result.Response.AcceleratorRegionSet
		return nil
	})

	if err != nil {
		return err
	}

	regionSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, regionItem := range respData {
			regionMap := map[string]interface{}{}
			if regionItem.Name != nil {
				regionMap["name"] = regionItem.Name
			}

			if regionItem.IsAvailable != nil {
				regionMap["is_available"] = regionItem.IsAvailable
			}

			if regionItem.Region != nil {
				regionMap["region"] = regionItem.Region
			}

			if regionItem.AreaName != nil {
				regionMap["area_name"] = regionItem.AreaName
			}

			if regionItem.IsChinaMainland != nil {
				regionMap["is_china_mainland"] = regionItem.IsChinaMainland
			}

			if regionItem.SupportIspType != nil {
				ispTypes := make([]string, 0, len(regionItem.SupportIspType))
				for _, isp := range regionItem.SupportIspType {
					if isp != nil {
						ispTypes = append(ispTypes, *isp)
					}
				}
				regionMap["support_isp_type"] = ispTypes
			}

			if regionItem.IsTencentRegion != nil {
				regionMap["is_tencent_region"] = regionItem.IsTencentRegion
			}

			regionSetList = append(regionSetList, regionMap)
		}
	}

	_ = d.Set("accelerator_region_set", regionSetList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), regionSetList); e != nil {
			return e
		}
	}

	return nil
}
