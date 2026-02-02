package cvm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	svccbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"cpu_core_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of CPU cores of the instance.",
			},
			"gpu_core_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of GPU cores of the instance.",
			},
			"memory_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Instance memory capacity, unit in GB.",
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"filter"},
				Description:   "The available zone that the CVM instance locates at. This field is conflict with `filter`.",
			},
			"filter": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      10,
				ConflictsWith: []string{"availability_zone"},
				Description:   "One or more name/value pairs to filter. This field is conflict with `availability_zone`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The filter name. Valid values: `zone`, `instance-family` and `instance-charge-type`.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The filter values.",
						},
					},
				},
			},
			"cbs_filter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cbs filter.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: "Hard disk media type. Value range:\n" +
								"	- CLOUD_BASIC: Represents ordinary Cloud Block Storage;\n" +
								"	- CLOUD_PREMIUM: Represents high-performance Cloud Block Storage;\n" +
								"	- CLOUD_SSD: Represents SSD Cloud Block Storage;\n" +
								"	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.",
						},
						"disk_charge_type": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Payment model. Value range:\n" +
								"	- PREPAID: Prepaid;\n" +
								"	- POSTPAID_BY_HOUR: Post-payment.",
						},
						"disk_usage": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "System disk or data disk. Value range:\n" +
								"	- SYSTEM_DISK: Represents the system disk;\n" +
								"	- DATA_DISK: Represents the data disk.",
						},
					},
				},
			},
			"exclude_sold_out": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate to filter instances types that is sold out or not, default is false.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values.
			"instance_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of cvm instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The available zone that the CVM instance locates at.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the instance.",
						},
						"cpu_core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of CPU cores of the instance.",
						},
						"gpu_core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of GPU cores of the instance.",
						},
						"memory_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance memory capacity, unit in GB.",
						},
						"family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type series of the instance.",
						},
						"instance_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge type of the instance.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sell status of the instance.",
						},
						"network_card": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Network card type, for example: 25 represents 25G network card.",
						},
						"type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type display name.",
						},
						"sold_out_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reason for sold out status.",
						},
						"instance_bandwidth": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Internal network bandwidth, unit: Gbps.",
						},
						"instance_pps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Network packet forwarding capacity, unit: 10K PPS.",
						},
						"storage_block_amount": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of local storage blocks.",
						},
						"cpu_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Processor model.",
						},
						"fpga": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of FPGA cores.",
						},
						"gpu_count": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Physical GPU card count mapped to instance. vGPU type is less than 1, direct-attach GPU type is greater than or equal to 1.",
						},
						"frequency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CPU frequency information.",
						},
						"status_category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stock status category. Valid values: EnoughStock, NormalStock, UnderStock, WithoutStock.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance remark information.",
						},
						"cbs_configs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CBS config. The cbs_configs is populated when the cbs_filter is added.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"available": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the configuration is available.",
									},
									"disk_charge_type": {
										Type:     schema.TypeString,
										Computed: true,
										Description: "Payment model. Value range:\n" +
											"	- PREPAID: Prepaid;\n" +
											"	- POSTPAID_BY_HOUR: Post-payment.",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The availability zone to which the Cloud Block Storage belongs.",
									},
									"instance_family": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance family.",
									},
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
										Description: "Hard disk media type. Value range:\n" +
											"	- CLOUD_BASIC: Represents ordinary Cloud Block Storage;\n" +
											"	- CLOUD_PREMIUM: Represents high-performance Cloud Block Storage;\n" +
											"	- CLOUD_SSD: Represents SSD Cloud Block Storage;\n" +
											"	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.",
									},
									"step_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum step size change in cloud disk size, in GB.",
									},
									"extra_performance_range": {
										Computed:    true,
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "Extra performance range.",
									},
									"device_class": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Device class.",
									},
									"disk_usage": {
										Type:     schema.TypeString,
										Computed: true,
										Description: "Cloud disk type. Value range:\n" +
											"	- SYSTEM_DISK: Represents the system disk;\n" +
											"	- DATA_DISK: Represents the data disk.",
									},
									"min_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum configurable cloud disk size, in GB.",
									},
									"max_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum configurable cloud disk size, in GB.",
									},
								},
							},
						},
						"local_disk_type_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of local disk specifications. Empty if instance type does not support local disks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Local disk type.",
									},
									"partition_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Local disk partition type.",
									},
									"min_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum size of local disk, in GB.",
									},
									"max_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum size of local disk, in GB.",
									},
									"required": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether local disk is required when purchasing. Valid values: REQUIRED, OPTIONAL.",
									},
								},
							},
						},
						"price": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance pricing information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unit_price": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Subsequent unit price, used in postpaid mode, unit: CNY.",
									},
									"charge_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subsequent billing unit. Valid values: HOUR, GB.",
									},
									"original_price": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Original price for prepaid mode, unit: CNY.",
									},
									"discount_price": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Discount price for prepaid mode, unit: CNY.",
									},
									"discount": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Discount rate. For example, 20.0 means 20% off.",
									},
									"unit_price_discount": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Subsequent discount unit price, used in postpaid mode, unit: CNY.",
									},
									"unit_price_second_step": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Subsequent unit price for time range (96, 360) hours in postpaid mode, unit: CNY.",
									},
									"unit_price_discount_second_step": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Subsequent discount unit price for time range (96, 360) hours in postpaid mode, unit: CNY.",
									},
									"unit_price_third_step": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Specifies the original price of subsequent total costs with a usage time interval exceeding 360 hr in postpaid billing mode. measurement unit: usd.",
									},
									"unit_price_discount_third_step": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Discounted price of subsequent total cost for usage time interval exceeding 360 hr in postpaid billing mode. measurement unit: usd.",
									},
								},
							},
						},
						"externals": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Extended attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"release_address": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to release address.",
									},
									"unsupport_networks": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Unsupported network types. Valid values: BASIC (basic network), VPC1.0 (VPC 1.0).",
									},
									"storage_block_attr": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "HDD local storage attributes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Storage block type.",
												},
												"min_size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Minimum size of storage block, in GB.",
												},
												"max_size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum size of storage block, in GB.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_instance_types.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	isExcludeSoldOut := d.Get("exclude_sold_out").(bool)
	cpu, cpuOk := d.GetOk("cpu_core_count")
	gpu, gpuOk := d.GetOk("gpu_core_count")
	memory, memoryOk := d.GetOk("memory_size")
	var instanceSellTypes []*cvm.InstanceTypeQuotaItem
	var errRet error
	var err error
	typeList := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	var zone string
	var zone_in = 0
	if v, ok := d.GetOk("availability_zone"); ok {
		zone = v.(string)
		zone_in = 1
	}
	filters := d.Get("filter").(*schema.Set).List()
	filterMap := make(map[string][]string, len(filters)+zone_in)
	for _, v := range filters {
		item := v.(map[string]interface{})
		name := item["name"].(string)
		values := item["values"].([]interface{})
		filterValues := make([]string, 0, len(values))
		for _, value := range values {
			filterValues = append(filterValues, value.(string))
		}
		filterMap[name] = filterValues
	}
	if zone != "" {
		filterMap["zone"] = []string{zone}
	}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instanceSellTypes, errRet = cvmService.DescribeInstancesSellTypeByFilter(ctx, filterMap)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, instanceType := range instanceSellTypes {
		flag := true
		if cpuOk && int64(cpu.(int)) != *instanceType.Cpu {
			flag = false
		}
		if gpuOk && int64(gpu.(int)) != *instanceType.Gpu {
			flag = false
		}
		if memoryOk && int64(memory.(int)) != *instanceType.Memory {
			flag = false
		}
		if isExcludeSoldOut && CVM_SOLD_OUT_STATUS == *instanceType.Status {
			flag = false
		}

		if flag {
			mapping := map[string]interface{}{
				"availability_zone":    instanceType.Zone,
				"cpu_core_count":       instanceType.Cpu,
				"gpu_core_count":       instanceType.Gpu,
				"memory_size":          instanceType.Memory,
				"family":               instanceType.InstanceFamily,
				"instance_type":        instanceType.InstanceType,
				"instance_charge_type": instanceType.InstanceChargeType,
				"status":               instanceType.Status,
				"network_card":         instanceType.NetworkCard,
				"type_name":            instanceType.TypeName,
				"sold_out_reason":      instanceType.SoldOutReason,
				"instance_bandwidth":   instanceType.InstanceBandwidth,
				"instance_pps":         instanceType.InstancePps,
				"storage_block_amount": instanceType.StorageBlockAmount,
				"cpu_type":             instanceType.CpuType,
				"fpga":                 instanceType.Fpga,
				"gpu_count":            instanceType.GpuCount,
				"frequency":            instanceType.Frequency,
				"status_category":      instanceType.StatusCategory,
				"remark":               instanceType.Remark,
			}

			// Map local_disk_type_list
			localDiskList := make([]interface{}, 0)
			for _, localDisk := range instanceType.LocalDiskTypeList {
				localDiskList = append(localDiskList, map[string]interface{}{
					"type":           localDisk.Type,
					"partition_type": localDisk.PartitionType,
					"min_size":       localDisk.MinSize,
					"max_size":       localDisk.MaxSize,
					"required":       localDisk.Required,
				})
			}
			mapping["local_disk_type_list"] = localDiskList

			// Map price
			priceList := make([]interface{}, 0)
			if instanceType.Price != nil {
				priceList = append(priceList, map[string]interface{}{
					"unit_price":                      instanceType.Price.UnitPrice,
					"charge_unit":                     instanceType.Price.ChargeUnit,
					"original_price":                  instanceType.Price.OriginalPrice,
					"discount_price":                  instanceType.Price.DiscountPrice,
					"discount":                        instanceType.Price.Discount,
					"unit_price_discount":             instanceType.Price.UnitPriceDiscount,
					"unit_price_second_step":          instanceType.Price.UnitPriceSecondStep,
					"unit_price_discount_second_step": instanceType.Price.UnitPriceDiscountSecondStep,
					"unit_price_third_step":           instanceType.Price.UnitPriceThirdStep,
					"unit_price_discount_third_step":  instanceType.Price.UnitPriceDiscountThirdStep,
				})
			}
			mapping["price"] = priceList

			// Map externals
			externalsList := make([]interface{}, 0)
			if instanceType.Externals != nil {
				externalsMap := map[string]interface{}{
					"release_address":    instanceType.Externals.ReleaseAddress,
					"unsupport_networks": instanceType.Externals.UnsupportNetworks,
				}

				// Map storage_block_attr
				storageBlockList := make([]interface{}, 0)
				if instanceType.Externals.StorageBlockAttr != nil {
					storageBlockList = append(storageBlockList, map[string]interface{}{
						"type":     instanceType.Externals.StorageBlockAttr.Type,
						"min_size": instanceType.Externals.StorageBlockAttr.MinSize,
						"max_size": instanceType.Externals.StorageBlockAttr.MaxSize,
					})
				}
				externalsMap["storage_block_attr"] = storageBlockList
				externalsList = append(externalsList, externalsMap)
			}
			mapping["externals"] = externalsList

			typeList = append(typeList, mapping)
			ids = append(ids, *instanceType.InstanceType)
		}
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cbsService := svccbs.NewCbsService(client)
	cbsFilterParams := make(map[string]interface{})
	var hasCbsFilter bool
	if dMap, ok := helper.InterfacesHeadMap(d, "cbs_filter"); ok {
		if v, ok := dMap["disk_types"].([]interface{}); ok && len(v) > 0 {
			cbsFilterParams["disk_types"] = helper.InterfacesStrings(v)
		}
		if v, ok := dMap["disk_charge_type"].(string); ok && v != "" {
			cbsFilterParams["disk_charge_type"] = v
		}
		if v, ok := dMap["disk_usage"].(string); ok && v != "" {
			cbsFilterParams["disk_usage"] = v
		}
		hasCbsFilter = true
	}
	if hasCbsFilter {
		for idx, t := range typeList {
			filterParams := make(map[string]interface{})
			for k, v := range cbsFilterParams {
				filterParams[k] = v
			}

			if v, ok := t["availability_zone"].(*string); ok && v != nil {
				filterParams["availability_zone"] = *v
			}
			if v, ok := t["cpu_core_count"].(*int64); ok && v != nil {
				filterParams["cpu_core_count"] = *v
			}
			if v, ok := t["memory_size"].(*int64); ok && v != nil {
				filterParams["memory_size"] = *v
			}
			if v, ok := t["family"].(*string); ok && v != nil {
				filterParams["family"] = *v
			}
			diskConfigSet, err := cbsService.DescribeDiskConfigQuota(ctx, filterParams)
			if err != nil {
				return err
			}
			cbsConfigList := make([]interface{}, 0)
			for _, diskConfig := range diskConfigSet {
				cbsConfigList = append(cbsConfigList, map[string]interface{}{
					"available":               diskConfig.Available,
					"disk_charge_type":        diskConfig.DiskChargeType,
					"zone":                    diskConfig.Zone,
					"instance_family":         diskConfig.InstanceFamily,
					"disk_type":               diskConfig.DiskType,
					"step_size":               diskConfig.StepSize,
					"extra_performance_range": diskConfig.ExtraPerformanceRange,
					"device_class":            diskConfig.DeviceClass,
					"disk_usage":              diskConfig.DiskUsage,
					"min_disk_size":           diskConfig.MinDiskSize,
					"max_disk_size":           diskConfig.MaxDiskSize,
				})
			}
			typeList[idx]["cbs_configs"] = cbsConfigList
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_types", typeList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance type list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), typeList); err != nil {
			return err
		}
	}
	return nil
}
