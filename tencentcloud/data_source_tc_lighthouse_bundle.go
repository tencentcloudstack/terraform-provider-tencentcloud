/*
Use this data source to query detailed information of lighthouse bundle

Example Usage

```hcl
data "tencentcloud_lighthouse_bundle" "bundle" {
  bundle_ids =
  offset = 0
  limit = 20
  zones =
  filters {
		name = "bundle-id"
		values =

  }
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseBundle() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseBundleRead,
		Schema: map[string]*schema.Schema{
			"bundle_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Bundle ID list.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Offset. Default value is 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of returned results. Default value is 20. Maximum value is 100.",
			},

			"zones": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Zone list, which contains all zones by default.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter listbundle-idFilter by the bundle ID.Type: StringRequired: Nosupport-platform-typeFilter by the OS type.Valid values: LINUX_UNIX (Linux or Unix), WINDOWS (Windows)Type: StringRequired: Nobundle-typeFilter by the bundle type.Valid values: GENERAL_BUNDLE (General bundle), STORAGE_BUNDLE (Storage bundle), ENTERPRISE_BUNDLE (Enterprise bundle), EXCLUSIVE_BUNDLE (Dedicated bundle), BEFAST_BUNDLE (BeFast bundle)Type: StringRequired: Nobundle-stateFilter by the bundle status.Valid values: ONLINE, OFFLINEType: StringRequired: NoEach request can contain up to 10 Filters, and up to 5 Filter.Values for each filter. You cannot specify both BundleIds and Filters at the same time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter value of field.",
						},
					},
				},
			},

			"bundle_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of bundle details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bundle_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package ID.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size in GB.",
						},
						"system_disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System disk type.",
						},
						"system_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "System disk size.",
						},
						"monthly_traffic": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Monthly network traffic in Gb.",
						},
						"support_linux_unix_platform": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether Linux/Unix is supported.",
						},
						"support_windows_platform": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether Windows is supported.",
						},
						"price": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Current package unit price information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_price": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Instance price.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"original_bundle_price": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Original package unit price.",
												},
												"original_price": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Original price.",
												},
												"discount": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Discount.",
												},
												"discount_price": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Discounted price.",
												},
												"currency": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Currency unit. Valid values: CNY and USD.",
												},
											},
										},
									},
								},
							},
						},
						"c_p_u": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU.",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak bandwidth in Mbps.",
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network billing mode.",
						},
						"bundle_sales_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package sale status. Valid values are AVAILABLE, SOLD_OUT.",
						},
						"bundle_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package type.Valid values:GENERAL_BUNDLE: generalSTORAGE_BUNDLE: Storage.",
						},
						"bundle_display_label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package tag.Valid values:ACTIVITY: promotional packageNORMAL: regular packageCAREFREE: carefree package.",
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

func dataSourceTencentCloudLighthouseBundleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_bundle.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("bundle_ids"); ok {
		bundleIdsSet := v.(*schema.Set).List()
		paramMap["BundleIds"] = helper.InterfacesStringsPoint(bundleIdsSet)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		paramMap["Zones"] = helper.InterfacesStringsPoint(zonesSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*lighthouse.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := lighthouse.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var bundleSet []*lighthouse.Bundle

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseBundleByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		bundleSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(bundleSet))
	tmpList := make([]map[string]interface{}, 0, len(bundleSet))

	if bundleSet != nil {
		for _, bundle := range bundleSet {
			bundleMap := map[string]interface{}{}

			if bundle.BundleId != nil {
				bundleMap["bundle_id"] = bundle.BundleId
			}

			if bundle.Memory != nil {
				bundleMap["memory"] = bundle.Memory
			}

			if bundle.SystemDiskType != nil {
				bundleMap["system_disk_type"] = bundle.SystemDiskType
			}

			if bundle.SystemDiskSize != nil {
				bundleMap["system_disk_size"] = bundle.SystemDiskSize
			}

			if bundle.MonthlyTraffic != nil {
				bundleMap["monthly_traffic"] = bundle.MonthlyTraffic
			}

			if bundle.SupportLinuxUnixPlatform != nil {
				bundleMap["support_linux_unix_platform"] = bundle.SupportLinuxUnixPlatform
			}

			if bundle.SupportWindowsPlatform != nil {
				bundleMap["support_windows_platform"] = bundle.SupportWindowsPlatform
			}

			if bundle.Price != nil {
				priceMap := map[string]interface{}{}

				if bundle.Price.InstancePrice != nil {
					instancePriceMap := map[string]interface{}{}

					if bundle.Price.InstancePrice.OriginalBundlePrice != nil {
						instancePriceMap["original_bundle_price"] = bundle.Price.InstancePrice.OriginalBundlePrice
					}

					if bundle.Price.InstancePrice.OriginalPrice != nil {
						instancePriceMap["original_price"] = bundle.Price.InstancePrice.OriginalPrice
					}

					if bundle.Price.InstancePrice.Discount != nil {
						instancePriceMap["discount"] = bundle.Price.InstancePrice.Discount
					}

					if bundle.Price.InstancePrice.DiscountPrice != nil {
						instancePriceMap["discount_price"] = bundle.Price.InstancePrice.DiscountPrice
					}

					if bundle.Price.InstancePrice.Currency != nil {
						instancePriceMap["currency"] = bundle.Price.InstancePrice.Currency
					}

					priceMap["instance_price"] = []interface{}{instancePriceMap}
				}

				bundleMap["price"] = []interface{}{priceMap}
			}

			if bundle.CPU != nil {
				bundleMap["c_p_u"] = bundle.CPU
			}

			if bundle.InternetMaxBandwidthOut != nil {
				bundleMap["internet_max_bandwidth_out"] = bundle.InternetMaxBandwidthOut
			}

			if bundle.InternetChargeType != nil {
				bundleMap["internet_charge_type"] = bundle.InternetChargeType
			}

			if bundle.BundleSalesState != nil {
				bundleMap["bundle_sales_state"] = bundle.BundleSalesState
			}

			if bundle.BundleType != nil {
				bundleMap["bundle_type"] = bundle.BundleType
			}

			if bundle.BundleDisplayLabel != nil {
				bundleMap["bundle_display_label"] = bundle.BundleDisplayLabel
			}

			ids = append(ids, *bundle.BundleId)
			tmpList = append(tmpList, bundleMap)
		}

		_ = d.Set("bundle_set", tmpList)
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
