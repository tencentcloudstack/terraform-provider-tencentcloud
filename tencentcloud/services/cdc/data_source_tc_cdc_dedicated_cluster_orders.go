package cdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc/v20201214"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCdcDedicatedClusterOrders() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceTencentCloudCdcDedicatedClusterOrdersRead,
		Schema: map[string]*schema.Schema{
			"dedicated_cluster_ids": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter by Dedicated Cluster ID.",
			},
			//"dedicated_cluster_order_id": {
			//	Optional:    true,
			//	Type:        schema.TypeString,
			//	Description: "Filter by Dedicated Cluster Order ID.",
			//},
			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter by Dedicated Cluster Order Status. Allow filter value: PENDING, INCONSTRUCTION, DELIVERING, DELIVERED, EXPIRED, CANCELLED, OFFLINE.",
			},
			"action_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter by Dedicated Cluster Order Action Type. Allow filter value: CREATE, EXTEND.",
			},
			// computed
			"dedicated_cluster_order_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Filter by Dedicated Cluster Order.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dedicated_cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster ID.",
						},
						"dedicated_cluster_type_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Type ID.",
						},
						"supported_storage_type": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Dedicated Cluster Storage Type.",
						},
						"supported_uplink_speed": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "Dedicated Cluster Supported Uplink Speed.",
						},
						"supported_instance_family": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Dedicated Cluster Supported Instance Family.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dedicated Cluster Supported Weight.",
						},
						"power_draw": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Dedicated Cluster Supported PowerDraw.",
						},
						"order_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Order Status.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Order Create time.",
						},
						"dedicated_cluster_order_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Order ID.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Order Action Type.",
						},
						"dedicated_cluster_order_items": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Dedicated Cluster Order Item List.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dedicated_cluster_type_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dedicated Cluster ID.",
									},
									"supported_storage_type": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Dedicated Cluster Storage Type.",
									},
									"supported_uplink_speed": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "Dedicated Cluster Supported Uplink Speed.",
									},
									"supported_instance_family": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Dedicated Cluster Supported Instance Family.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Dedicated Cluster Supported Weight.",
									},
									"power_draw": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Dedicated Cluster Supported PowerDraw.",
									},
									"sub_order_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dedicated Cluster Order Status.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dedicated Cluster Order Create time.",
									},
									"sub_order_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dedicated Cluster SubOrder ID.",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Dedicated Cluster SubOrder Count.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dedicated Cluster Type Name.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dedicated Cluster Type Description.",
									},
									"total_cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Dedicated Cluster Total CPU.",
									},
									"total_mem": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Dedicated Cluster Total Memory.",
									},
									"total_gpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Dedicated Cluster Total GPU.",
									},
									"type_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dedicated Cluster Type Name.",
									},
									"compute_format": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dedicated Cluster Compute Format.",
									},
									"type_family": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dedicated Cluster Type Family.",
									},
									"sub_order_pay_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Dedicated Cluster SubOrder Pay Status.",
									},
								},
							},
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dedicated Cluster CPU.",
						},
						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dedicated Cluster Memory.",
						},
						"gpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dedicated Cluster GPU.",
						},
						"pay_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dedicated Cluster Order Pay Status.",
						},
						"pay_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Order Pay Type.",
						},
						"time_unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Order Pay Time Unit.",
						},
						"time_span": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dedicated Cluster Order Pay Time Span.",
						},
						"order_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Order Type.",
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

func DataSourceTencentCloudCdcDedicatedClusterOrdersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cdc_dedicated_cluster_orders.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                    = tccommon.GetLogId(tccommon.ContextNil)
		ctx                      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service                  = CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dedicatedClusterOrderSet []*cdc.DedicatedClusterOrder
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("dedicated_cluster_ids"); ok {
		dedicatedClusterIdsSet := v.(*schema.Set).List()
		paramMap["DedicatedClusterIds"] = helper.InterfacesStringsPoint(dedicatedClusterIdsSet)
	}

	//if v, ok := d.GetOk("dedicated_cluster_order_id"); ok {
	//	paramMap["DedicatedClusterOrderIds"] = helper.String(v.(string))
	//}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("action_type"); ok {
		paramMap["ActionType"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdcDedicatedClusterOrdersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		dedicatedClusterOrderSet = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(dedicatedClusterOrderSet))
	tmpList := make([]map[string]interface{}, 0, len(dedicatedClusterOrderSet))

	if dedicatedClusterOrderSet != nil {
		for _, dedicatedClusterOrder := range dedicatedClusterOrderSet {
			dedicatedClusterOrderMap := map[string]interface{}{}

			if dedicatedClusterOrder.DedicatedClusterId != nil {
				dedicatedClusterOrderMap["dedicated_cluster_id"] = dedicatedClusterOrder.DedicatedClusterId
			}

			if dedicatedClusterOrder.DedicatedClusterTypeId != nil {
				dedicatedClusterOrderMap["dedicated_cluster_type_id"] = dedicatedClusterOrder.DedicatedClusterTypeId
			}

			if dedicatedClusterOrder.SupportedStorageType != nil {
				dedicatedClusterOrderMap["supported_storage_type"] = dedicatedClusterOrder.SupportedStorageType
			}

			if dedicatedClusterOrder.SupportedUplinkSpeed != nil {
				dedicatedClusterOrderMap["supported_uplink_speed"] = dedicatedClusterOrder.SupportedUplinkSpeed
			}

			if dedicatedClusterOrder.SupportedInstanceFamily != nil {
				dedicatedClusterOrderMap["supported_instance_family"] = dedicatedClusterOrder.SupportedInstanceFamily
			}

			if dedicatedClusterOrder.Weight != nil {
				dedicatedClusterOrderMap["weight"] = dedicatedClusterOrder.Weight
			}

			if dedicatedClusterOrder.PowerDraw != nil {
				dedicatedClusterOrderMap["power_draw"] = dedicatedClusterOrder.PowerDraw
			}

			if dedicatedClusterOrder.OrderStatus != nil {
				dedicatedClusterOrderMap["order_status"] = dedicatedClusterOrder.OrderStatus
			}

			if dedicatedClusterOrder.CreateTime != nil {
				dedicatedClusterOrderMap["create_time"] = dedicatedClusterOrder.CreateTime
			}

			if dedicatedClusterOrder.DedicatedClusterOrderId != nil {
				dedicatedClusterOrderMap["dedicated_cluster_order_id"] = dedicatedClusterOrder.DedicatedClusterOrderId
			}

			if dedicatedClusterOrder.Action != nil {
				dedicatedClusterOrderMap["action"] = dedicatedClusterOrder.Action
			}

			if dedicatedClusterOrder.DedicatedClusterOrderItems != nil {
				dedicatedClusterOrderItemsList := []interface{}{}
				for _, dedicatedClusterOrderItems := range dedicatedClusterOrder.DedicatedClusterOrderItems {
					dedicatedClusterOrderItemsMap := map[string]interface{}{}

					if dedicatedClusterOrderItems.DedicatedClusterTypeId != nil {
						dedicatedClusterOrderItemsMap["dedicated_cluster_type_id"] = dedicatedClusterOrderItems.DedicatedClusterTypeId
					}

					if dedicatedClusterOrderItems.SupportedStorageType != nil {
						dedicatedClusterOrderItemsMap["supported_storage_type"] = dedicatedClusterOrderItems.SupportedStorageType
					}

					if dedicatedClusterOrderItems.SupportedUplinkSpeed != nil {
						dedicatedClusterOrderItemsMap["supported_uplink_speed"] = dedicatedClusterOrderItems.SupportedUplinkSpeed
					}

					if dedicatedClusterOrderItems.SupportedInstanceFamily != nil {
						dedicatedClusterOrderItemsMap["supported_instance_family"] = dedicatedClusterOrderItems.SupportedInstanceFamily
					}

					if dedicatedClusterOrderItems.Weight != nil {
						dedicatedClusterOrderItemsMap["weight"] = dedicatedClusterOrderItems.Weight
					}

					if dedicatedClusterOrderItems.PowerDraw != nil {
						dedicatedClusterOrderItemsMap["power_draw"] = dedicatedClusterOrderItems.PowerDraw
					}

					if dedicatedClusterOrderItems.SubOrderStatus != nil {
						dedicatedClusterOrderItemsMap["sub_order_status"] = dedicatedClusterOrderItems.SubOrderStatus
					}

					if dedicatedClusterOrderItems.CreateTime != nil {
						dedicatedClusterOrderItemsMap["create_time"] = dedicatedClusterOrderItems.CreateTime
					}

					if dedicatedClusterOrderItems.SubOrderId != nil {
						dedicatedClusterOrderItemsMap["sub_order_id"] = dedicatedClusterOrderItems.SubOrderId
					}

					if dedicatedClusterOrderItems.Count != nil {
						dedicatedClusterOrderItemsMap["count"] = dedicatedClusterOrderItems.Count
					}

					if dedicatedClusterOrderItems.Name != nil {
						dedicatedClusterOrderItemsMap["name"] = dedicatedClusterOrderItems.Name
					}

					if dedicatedClusterOrderItems.Description != nil {
						dedicatedClusterOrderItemsMap["description"] = dedicatedClusterOrderItems.Description
					}

					if dedicatedClusterOrderItems.TotalCpu != nil {
						dedicatedClusterOrderItemsMap["total_cpu"] = dedicatedClusterOrderItems.TotalCpu
					}

					if dedicatedClusterOrderItems.TotalMem != nil {
						dedicatedClusterOrderItemsMap["total_mem"] = dedicatedClusterOrderItems.TotalMem
					}

					if dedicatedClusterOrderItems.TotalGpu != nil {
						dedicatedClusterOrderItemsMap["total_gpu"] = dedicatedClusterOrderItems.TotalGpu
					}

					if dedicatedClusterOrderItems.TypeName != nil {
						dedicatedClusterOrderItemsMap["type_name"] = dedicatedClusterOrderItems.TypeName
					}

					if dedicatedClusterOrderItems.ComputeFormat != nil {
						dedicatedClusterOrderItemsMap["compute_format"] = dedicatedClusterOrderItems.ComputeFormat
					}

					if dedicatedClusterOrderItems.TypeFamily != nil {
						dedicatedClusterOrderItemsMap["type_family"] = dedicatedClusterOrderItems.TypeFamily
					}

					if dedicatedClusterOrderItems.SubOrderPayStatus != nil {
						dedicatedClusterOrderItemsMap["sub_order_pay_status"] = dedicatedClusterOrderItems.SubOrderPayStatus
					}

					dedicatedClusterOrderItemsList = append(dedicatedClusterOrderItemsList, dedicatedClusterOrderItemsMap)
				}

				dedicatedClusterOrderMap["dedicated_cluster_order_items"] = dedicatedClusterOrderItemsList
			}

			if dedicatedClusterOrder.Cpu != nil {
				dedicatedClusterOrderMap["cpu"] = dedicatedClusterOrder.Cpu
			}

			if dedicatedClusterOrder.Mem != nil {
				dedicatedClusterOrderMap["mem"] = dedicatedClusterOrder.Mem
			}

			if dedicatedClusterOrder.Gpu != nil {
				dedicatedClusterOrderMap["gpu"] = dedicatedClusterOrder.Gpu
			}

			if dedicatedClusterOrder.PayStatus != nil {
				dedicatedClusterOrderMap["pay_status"] = dedicatedClusterOrder.PayStatus
			}

			if dedicatedClusterOrder.PayType != nil {
				dedicatedClusterOrderMap["pay_type"] = dedicatedClusterOrder.PayType
			}

			if dedicatedClusterOrder.TimeUnit != nil {
				dedicatedClusterOrderMap["time_unit"] = dedicatedClusterOrder.TimeUnit
			}

			if dedicatedClusterOrder.TimeSpan != nil {
				dedicatedClusterOrderMap["time_span"] = dedicatedClusterOrder.TimeSpan
			}

			if dedicatedClusterOrder.OrderType != nil {
				dedicatedClusterOrderMap["order_type"] = dedicatedClusterOrder.OrderType
			}

			ids = append(ids, *dedicatedClusterOrder.DedicatedClusterId)
			tmpList = append(tmpList, dedicatedClusterOrderMap)
		}

		_ = d.Set("dedicated_cluster_order_set", tmpList)
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
