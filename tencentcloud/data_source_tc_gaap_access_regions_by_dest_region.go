/*
Use this data source to query detailed information of gaap access regions by dest region

Example Usage

```hcl
data "tencentcloud_gaap_access_regions_by_dest_region" "access_regions_by_dest_region" {
  dest_region = "SouthChina"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapAccessRegionsByDestRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapAccessRegionsByDestRegionRead,
		Schema: map[string]*schema.Schema{
			"dest_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Origin region.",
			},

			"ip_address_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "IP version, can be taken as IPv4 or IPv6, with a default value of IPv4.",
			},

			"package_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Channel package type, where Thunder represents a standard proxy group, Accelerator represents a game accelerator proxy, and CrossBorder represents a cross-border proxy.",
			},

			"access_region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of available acceleration zone information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region id.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chinese or English name of the region.",
						},
						"concurrent_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Optional concurrency value array.",
						},
						"bandwidth_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Optional bandwidth value array.",
						},
						"region_area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region of the computer room.",
						},
						"region_area_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name of the computer room.",
						},
						"idc_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of computer room, where dc represents the DataCenter data center and ec represents the EdgeComputing edge node.",
						},
						"feature_bitmap": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The type of computer room, where dc represents the DataCenter data center, ec represents the feature bitmap, and each bit represents a feature, where:0, indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained. Edge nodes.",
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

func dataSourceTencentCloudGaapAccessRegionsByDestRegionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_access_regions_by_dest_region.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("dest_region"); ok {
		paramMap["dest_region"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_address_version"); ok {
		paramMap["ip_address_version"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_type"); ok {
		paramMap["package_type"] = helper.String(v.(string))
	}

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	var accessRegionSet []*gaap.AccessRegionDetial

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapAccessRegionsByDestRegionByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		accessRegionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(accessRegionSet))
	tmpList := make([]map[string]interface{}, 0, len(accessRegionSet))

	if accessRegionSet != nil {
		for _, accessRegionDetial := range accessRegionSet {
			accessRegionDetialMap := map[string]interface{}{}

			if accessRegionDetial.RegionId != nil {
				accessRegionDetialMap["region_id"] = accessRegionDetial.RegionId
			}

			if accessRegionDetial.RegionName != nil {
				accessRegionDetialMap["region_name"] = accessRegionDetial.RegionName
			}

			if accessRegionDetial.ConcurrentList != nil {
				accessRegionDetialMap["concurrent_list"] = accessRegionDetial.ConcurrentList
			}

			if accessRegionDetial.BandwidthList != nil {
				accessRegionDetialMap["bandwidth_list"] = accessRegionDetial.BandwidthList
			}

			if accessRegionDetial.RegionArea != nil {
				accessRegionDetialMap["region_area"] = accessRegionDetial.RegionArea
			}

			if accessRegionDetial.RegionAreaName != nil {
				accessRegionDetialMap["region_area_name"] = accessRegionDetial.RegionAreaName
			}

			if accessRegionDetial.IDCType != nil {
				accessRegionDetialMap["idc_type"] = accessRegionDetial.IDCType
			}

			if accessRegionDetial.FeatureBitmap != nil {
				accessRegionDetialMap["feature_bitmap"] = accessRegionDetial.FeatureBitmap
			}

			ids = append(ids, *accessRegionDetial.RegionId)
			tmpList = append(tmpList, accessRegionDetialMap)
		}

		_ = d.Set("access_region_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
