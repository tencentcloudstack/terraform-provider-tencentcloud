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
				Optional: true,
				Type:     schema.TypeList,
				Description: "Filter list.\n" +
					"- `bundle-id`: filter by the bundle ID.\n" +
					"- `support-platform-type`: filter by system type, valid values: `LINUX_UNIX`, `WINDOWS`.\n" +
					"- `bundle-type`: filter according to package type, valid values: `GENERAL_BUNDLE`, `STORAGE_BUNDLE`, `ENTERPRISE_BUNDLE`, `EXCLUSIVE_BUNDLE`, `BEFAST_BUNDLE`.\n" +
					"- `bundle-state`: filter according to package status, valid values: `ONLINE`, `OFFLINE`.\n" +
					"NOTE: The upper limit of Filters per request is 10. The upper limit of Filter.Values is 5. Parameter does not support specifying both BundleIds and Filters.",
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
						"cpu": {
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
		bundleIds := make([]string, 0)
		for _, bundleId := range v.(*schema.Set).List() {
			bundleIds = append(bundleIds, bundleId.(string))
		}
		paramMap["bundle_ids"] = bundleIds
	}

	if v, ok := d.GetOk("offset"); ok {
		paramMap["offset"] = v.(int)
	}

	if v, ok := d.GetOk("limit"); ok {
		paramMap["limit"] = v.(int)
	}

	if v, ok := d.GetOk("zones"); ok {
		zones := make([]string, 0)
		for _, zone := range v.(*schema.Set).List() {
			zones = append(zones, zone.(string))
		}
		paramMap["zones"] = zones
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

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

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
				bundleMap["cpu"] = bundle.CPU
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
