package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapDestRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapDestRegionsRead,
		Schema: map[string]*schema.Schema{
			"dest_region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Source Site Area Details List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region ID.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region name.",
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
							Description: "Property bitmap, where each bit represents a property, where:0, indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"support_feature": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Ability to access regional supportNote: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_type": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "A list of network types supported by the access area, with normal indicating support for regular BGP, cn2 indicating premium BGP, triple indicating three networks, and secure_EIP represents a custom secure EIP.",
									},
								},
							},
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

func dataSourceTencentCloudGaapDestRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_dest_regions.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	var destRegionSet []*gaap.RegionDetail

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapDestRegions(ctx)
		if e != nil {
			return retryError(e)
		}
		destRegionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(destRegionSet))
	tmpList := make([]map[string]interface{}, 0, len(destRegionSet))

	if destRegionSet != nil {
		for _, regionDetail := range destRegionSet {
			regionDetailMap := map[string]interface{}{}

			if regionDetail.RegionId != nil {
				regionDetailMap["region_id"] = regionDetail.RegionId
			}

			if regionDetail.RegionName != nil {
				regionDetailMap["region_name"] = regionDetail.RegionName
			}

			if regionDetail.RegionArea != nil {
				regionDetailMap["region_area"] = regionDetail.RegionArea
			}

			if regionDetail.RegionAreaName != nil {
				regionDetailMap["region_area_name"] = regionDetail.RegionAreaName
			}

			if regionDetail.IDCType != nil {
				regionDetailMap["idc_type"] = regionDetail.IDCType
			}

			if regionDetail.FeatureBitmap != nil {
				regionDetailMap["feature_bitmap"] = regionDetail.FeatureBitmap
			}

			if regionDetail.SupportFeature != nil {
				supportFeatureMap := map[string]interface{}{}

				if regionDetail.SupportFeature.NetworkType != nil {
					supportFeatureMap["network_type"] = regionDetail.SupportFeature.NetworkType
				}

				regionDetailMap["support_feature"] = []interface{}{supportFeatureMap}
			}

			ids = append(ids, *regionDetail.RegionId)
			tmpList = append(tmpList, regionDetailMap)
		}

		_ = d.Set("dest_region_set", tmpList)
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
