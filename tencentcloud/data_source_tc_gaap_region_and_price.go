/*
Use this data source to query detailed information of gaap region and price

Example Usage

```hcl
data "tencentcloud_gaap_region_and_price" "region_and_price" {
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

func dataSourceTencentCloudGaapRegionAndPrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapRegionAndPriceRead,
		Schema: map[string]*schema.Schema{
			"ip_address_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "IP version. Available values: IPv4, IPv6. Default is IPv4.",
			},

			"package_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Type of channel package. `Thunder` represents standard channel group, `Accelerator` represents game accelerator channel, and `CrossBorder` represents cross-border channel.",
			},

			"dest_region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Source Site Area Details List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region Id.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region Name.",
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
							Description: "Type of computer room, dc represents DataCenter data center, ec represents EdgeComputing edge node.",
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
										Description: "A list of network types supported by the access area, with `normal` indicating support for regular BGP, `cn2` indicating premium BGP, `triple` indicating three networks, and `secure_eip` represents a custom secure EIP.",
									},
								},
							},
						},
					},
				},
			},

			"bandwidth_unit_price": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "proxy bandwidth cost gradient price.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth_range": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Band width Range.",
						},
						"bandwidth_unit_price": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Band width Unit Price, Unit:yuan/Mbps/day.",
						},
						"discount_bandwidth_unit_price": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Bandwidth discount price, unit:yuan/Mbps/day.",
						},
					},
				},
			},

			"currency": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Bandwidth Price Currency Type:CNYUSD.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudGaapRegionAndPriceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_region_and_price.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ip_address_version"); ok {
		paramMap["IPAddressVersion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_type"); ok {
		paramMap["PackageType"] = helper.String(v.(string))
	}

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		destRegionSet      []*gaap.RegionDetail
		bandwidthUnitPrice []*gaap.BandwidthPriceGradient
		currency           *string
	)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		resultRegionAndPrice, resultBandwidthUnitPrice, resultCurrency, e := service.DescribeGaapRegionAndPriceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		destRegionSet = resultRegionAndPrice
		bandwidthUnitPrice = resultBandwidthUnitPrice
		currency = resultCurrency
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(destRegionSet))
	tmpDestRegionSetList := make([]map[string]interface{}, 0, len(destRegionSet))
	tmpbandwidthUnitPriceList := make([]map[string]interface{}, 0, len(bandwidthUnitPrice))

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
			tmpDestRegionSetList = append(tmpDestRegionSetList, regionDetailMap)
		}

		_ = d.Set("dest_region_set", tmpDestRegionSetList)
	}

	if bandwidthUnitPrice != nil {
		for _, bandwidthPriceGradient := range bandwidthUnitPrice {
			bandwidthPriceGradientMap := map[string]interface{}{}

			if bandwidthPriceGradient.BandwidthRange != nil {
				bandwidthPriceGradientMap["bandwidth_range"] = bandwidthPriceGradient.BandwidthRange
			}

			if bandwidthPriceGradient.BandwidthUnitPrice != nil {
				bandwidthPriceGradientMap["bandwidth_unit_price"] = bandwidthPriceGradient.BandwidthUnitPrice
			}

			if bandwidthPriceGradient.DiscountBandwidthUnitPrice != nil {
				bandwidthPriceGradientMap["discount_bandwidth_unit_price"] = bandwidthPriceGradient.DiscountBandwidthUnitPrice
			}

			tmpbandwidthUnitPriceList = append(tmpbandwidthUnitPriceList, bandwidthPriceGradientMap)
		}

		_ = d.Set("bandwidth_unit_price", tmpbandwidthUnitPriceList)
	}

	if currency != nil {
		_ = d.Set("currency", currency)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	result := map[string]interface{}{
		"dest_region_set":      tmpDestRegionSetList,
		"bandwidth_unit_price": tmpbandwidthUnitPriceList,
		"currency":             currency,
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
