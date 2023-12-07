package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseModifyInstanceBundle() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseModifyInstanceBundleRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
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

			"modify_bundle_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Change package details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"modify_price": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Change the price difference to be made up after the instance package.",
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
													Description: "Original unit price of the package.",
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
													Description: "A monetary unit of price. Value range CNY: RMB. USD: us dollar.",
												},
											},
										},
									},
								},
							},
						},
						"modify_bundle_state": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "Change the status of the package. Value:\n" +
								"- SOLD_OUT: the package is sold out;\n" +
								"- AVAILABLE: support package changes;\n" +
								"- UNAVAILABLE: package changes are not supported for the time being.",
						},
						"bundle": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Package information.",
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
									"bundle_type_description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Package type description information.",
									},
									"bundle_display_label": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Package tag.Valid values:ACTIVITY: promotional packageNORMAL: regular packageCAREFREE: carefree package.",
									},
								},
							},
						},
						"not_support_modify_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package change reason information is not supported. When the package status is changed to `AVAILABLE`, the information is empty.",
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

func dataSourceTencentCloudLighthouseModifyInstanceBundleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_modify_instance_bundle.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
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

	var modifyBundleSet []*lighthouse.ModifyBundle

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseModifyInstanceBundleByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		modifyBundleSet = result
		return nil
	})
	if err != nil {
		return err
	}

	if len(modifyBundleSet) == 0 {
		d.SetId("")
		return fmt.Errorf("Response is nil")
	}

	ids := make([]string, 0)
	tmpList := make([]map[string]interface{}, 0)

	for _, modifyBundle := range modifyBundleSet {
		modifyBundleMap := map[string]interface{}{}

		if modifyBundle.ModifyPrice != nil {
			modifyPriceMap := map[string]interface{}{}

			if modifyBundle.ModifyPrice.InstancePrice != nil {
				instancePriceMap := map[string]interface{}{}

				if modifyBundle.ModifyPrice.InstancePrice.OriginalBundlePrice != nil {
					instancePriceMap["original_bundle_price"] = modifyBundle.ModifyPrice.InstancePrice.OriginalBundlePrice
				}

				if modifyBundle.ModifyPrice.InstancePrice.OriginalPrice != nil {
					instancePriceMap["original_price"] = modifyBundle.ModifyPrice.InstancePrice.OriginalPrice
				}

				if modifyBundle.ModifyPrice.InstancePrice.Discount != nil {
					instancePriceMap["discount"] = modifyBundle.ModifyPrice.InstancePrice.Discount
				}

				if modifyBundle.ModifyPrice.InstancePrice.DiscountPrice != nil {
					instancePriceMap["discount_price"] = modifyBundle.ModifyPrice.InstancePrice.DiscountPrice
				}

				if modifyBundle.ModifyPrice.InstancePrice.Currency != nil {
					instancePriceMap["currency"] = modifyBundle.ModifyPrice.InstancePrice.Currency
				}

				modifyPriceMap["instance_price"] = []interface{}{instancePriceMap}
			}

			modifyBundleMap["modify_price"] = []interface{}{modifyPriceMap}
		}

		if modifyBundle.ModifyBundleState != nil {
			modifyBundleMap["modify_bundle_state"] = modifyBundle.ModifyBundleState
		}

		if modifyBundle.Bundle != nil {
			bundleMap := map[string]interface{}{}

			if modifyBundle.Bundle.BundleId != nil {
				bundleMap["bundle_id"] = modifyBundle.Bundle.BundleId
			}
			ids = append(ids, *modifyBundle.Bundle.BundleId)

			if modifyBundle.Bundle.Memory != nil {
				bundleMap["memory"] = modifyBundle.Bundle.Memory
			}

			if modifyBundle.Bundle.SystemDiskType != nil {
				bundleMap["system_disk_type"] = modifyBundle.Bundle.SystemDiskType
			}

			if modifyBundle.Bundle.SystemDiskSize != nil {
				bundleMap["system_disk_size"] = modifyBundle.Bundle.SystemDiskSize
			}

			if modifyBundle.Bundle.MonthlyTraffic != nil {
				bundleMap["monthly_traffic"] = modifyBundle.Bundle.MonthlyTraffic
			}

			if modifyBundle.Bundle.SupportLinuxUnixPlatform != nil {
				bundleMap["support_linux_unix_platform"] = modifyBundle.Bundle.SupportLinuxUnixPlatform
			}

			if modifyBundle.Bundle.SupportWindowsPlatform != nil {
				bundleMap["support_windows_platform"] = modifyBundle.Bundle.SupportWindowsPlatform
			}

			if modifyBundle.Bundle.Price != nil {
				priceMap := map[string]interface{}{}

				if modifyBundle.Bundle.Price.InstancePrice != nil {
					instancePriceMap := map[string]interface{}{}

					if modifyBundle.Bundle.Price.InstancePrice.OriginalBundlePrice != nil {
						instancePriceMap["original_bundle_price"] = modifyBundle.Bundle.Price.InstancePrice.OriginalBundlePrice
					}

					if modifyBundle.Bundle.Price.InstancePrice.OriginalPrice != nil {
						instancePriceMap["original_price"] = modifyBundle.Bundle.Price.InstancePrice.OriginalPrice
					}

					if modifyBundle.Bundle.Price.InstancePrice.Discount != nil {
						instancePriceMap["discount"] = modifyBundle.Bundle.Price.InstancePrice.Discount
					}

					if modifyBundle.Bundle.Price.InstancePrice.DiscountPrice != nil {
						instancePriceMap["discount_price"] = modifyBundle.Bundle.Price.InstancePrice.DiscountPrice
					}

					if modifyBundle.Bundle.Price.InstancePrice.Currency != nil {
						instancePriceMap["currency"] = modifyBundle.Bundle.Price.InstancePrice.Currency
					}

					priceMap["instance_price"] = []interface{}{instancePriceMap}
				}

				bundleMap["price"] = []interface{}{priceMap}
			}

			if modifyBundle.Bundle.CPU != nil {
				bundleMap["cpu"] = modifyBundle.Bundle.CPU
			}

			if modifyBundle.Bundle.InternetMaxBandwidthOut != nil {
				bundleMap["internet_max_bandwidth_out"] = modifyBundle.Bundle.InternetMaxBandwidthOut
			}

			if modifyBundle.Bundle.InternetChargeType != nil {
				bundleMap["internet_charge_type"] = modifyBundle.Bundle.InternetChargeType
			}

			if modifyBundle.Bundle.BundleSalesState != nil {
				bundleMap["bundle_sales_state"] = modifyBundle.Bundle.BundleSalesState
			}

			if modifyBundle.Bundle.BundleType != nil {
				bundleMap["bundle_type"] = modifyBundle.Bundle.BundleType
			}

			if modifyBundle.Bundle.BundleTypeDescription != nil {
				bundleMap["bundle_type_description"] = modifyBundle.Bundle.BundleTypeDescription
			}

			if modifyBundle.Bundle.BundleDisplayLabel != nil {
				bundleMap["bundle_display_label"] = modifyBundle.Bundle.BundleDisplayLabel
			}
			modifyBundleMap["bundle"] = []map[string]interface{}{bundleMap}

		}
		if modifyBundle.NotSupportModifyMessage != nil {
			modifyBundleMap["not_support_modify_message"] = modifyBundle.NotSupportModifyMessage
		}

		tmpList = append(tmpList, modifyBundleMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("modify_bundle_set", tmpList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
