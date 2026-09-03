package billing

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	billingv20180709 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudBillingBillDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudBillingBillDetailRead,
		Schema: map[string]*schema.Schema{
			"month": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Month, format yyyy-mm. Either Month or BeginTime&EndTime must be provided. If BeginTime&EndTime is provided, Month is ignored. Data within the last 18 months can be pulled at most.",
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Period start time, format yyyy-mm-dd hh:ii:ss. Must be paired with EndTime, and must be in the same month. Cross-month query is not supported, the query result is the whole month data.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Period end time, format yyyy-mm-dd hh:ii:ss. Must be paired with BeginTime, and must be in the same month. Cross-month query is not supported, the query result is the whole month data.",
			},
			"need_record_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to return the total number of records, 1 means yes, 0 means no.",
			},
			"pay_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Payment mode, prePay (prepaid) / postPay (pay-as-you-go).",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource ID.",
			},
			"action_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Transaction type name, such as prepay purchase, prepay renew, postpay deduction, etc.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID.",
			},
			"business_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Product name code.",
			},
			"context": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Context information returned by the last request, used to speed up pagination for Month>=2023-05.",
			},
			"payer_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Payer account ID (the account ID is the user's unique account identifier on Tencent Cloud). Defaults to querying the current account's bill.",
			},

			"detail_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Bill detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"business_code_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product name.",
						},
						"product_code_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sub-product name.",
						},
						"pay_mode_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing mode.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project name.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID.",
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource alias.",
						},
						"action_type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Transaction type.",
						},
						"order_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Order ID.",
						},
						"bill_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Transaction ID.",
						},
						"pay_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deduction time.",
						},
						"fee_begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time of use.",
						},
						"fee_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time of use.",
						},
						"component_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Component list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"component_code_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component type.",
									},
									"item_code_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component name.",
									},
									"single_price": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component list price.",
									},
									"specified_price": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component specified price. Deprecated.",
									},
									"price_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component price unit.",
									},
									"used_amount": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component usage.",
									},
									"used_amount_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component usage unit.",
									},
									"real_total_measure": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Original usage/duration before resource pack deduction.",
									},
									"deducted_measure": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deducted usage/duration (including resource pack).",
									},
									"time_span": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Usage duration.",
									},
									"time_unit_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Duration unit.",
									},
									"cost": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component original price.",
									},
									"discount": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Discount rate.",
									},
									"reduce_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Discount type.",
									},
									"real_cost": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Discounted total price.",
									},
									"voucher_pay_amount": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Voucher payment amount.",
									},
									"cash_pay_amount": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cash account payment amount.",
									},
									"incentive_pay_amount": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Incentive account payment amount.",
									},
									"transfer_pay_amount": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Transfer account payment amount.",
									},
									"item_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component type code.",
									},
									"component_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component name code.",
									},
									"contract_price": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Component unit price after discount.",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance type.",
									},
									"ri_time_span": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Reserved instance deducted usage duration.",
									},
									"original_cost_with_ri": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Reserved instance deducted component original price.",
									},
									"sp_deduction_rate": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Savings plan deduction rate.",
									},
									"sp_deduction": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Savings plan deduction amount. Deprecated.",
									},
									"original_cost_with_sp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Savings plan deducted component original price.",
									},
									"blended_discount": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Blended discount rate.",
									},
									"component_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Component config description list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Config description name.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Config description value.",
												},
											},
										},
									},
								},
							},
						},
						"payer_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Payer UIN.",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Owner UIN.",
						},
						"operate_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operator UIN.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag information list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"business_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product code.",
						},
						"product_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sub-product code.",
						},
						"action_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Transaction type code.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"price_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Price attributes.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"associated_order": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated transaction order.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prepay_purchase": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "New purchase order.",
									},
									"prepay_renew": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Renewal order.",
									},
									"prepay_modify_up": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Upgrade order.",
									},
									"reverse_order": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Reverse order.",
									},
									"new_order": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Adjusted order after discount.",
									},
									"original": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Original order before discount adjustment.",
									},
								},
							},
						},
						"formula": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Calculation description.",
						},
						"formula_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing rule link.",
						},
						"bill_day": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bill day.",
						},
						"bill_month": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bill month.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bill record ID.",
						},
						"region_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domestic/International code.",
						},
						"region_type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domestic/International name.",
						},
						"reserve_detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reserve detail (instance config).",
						},
						"discount_object": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Discount object.",
						},
						"discount_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Discount type.",
						},
						"discount_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Discount content.",
						},
					},
				},
			},

			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total record count.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudBillingBillDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_billing_bill_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BillingService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("month"); ok {
		paramMap["Month"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("begin_time"); ok {
		paramMap["BeginTime"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("need_record_num"); ok {
		paramMap["NeedRecordNum"] = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("pay_mode"); ok {
		paramMap["PayMode"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("resource_id"); ok {
		paramMap["ResourceId"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("action_type"); ok {
		paramMap["ActionType"] = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		paramMap["ProjectId"] = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("business_code"); ok {
		paramMap["BusinessCode"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("context"); ok {
		paramMap["Context"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("payer_uin"); ok {
		paramMap["PayerUin"] = helper.String(v.(string))
	}

	var (
		respData    []*billingv20180709.BillDetail
		total       uint64
		respContext *string
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, totalRet, contextRet, e := service.DescribeBillingBillDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		total = totalRet
		respContext = contextRet
		return nil
	})
	if err != nil {
		return err
	}

	detailSetList := make([]map[string]interface{}, 0, len(respData))
	for _, detail := range respData {
		detailMap := map[string]interface{}{}
		if detail.BusinessCodeName != nil {
			detailMap["business_code_name"] = detail.BusinessCodeName
		}
		if detail.ProductCodeName != nil {
			detailMap["product_code_name"] = detail.ProductCodeName
		}
		if detail.PayModeName != nil {
			detailMap["pay_mode_name"] = detail.PayModeName
		}
		if detail.ProjectName != nil {
			detailMap["project_name"] = detail.ProjectName
		}
		if detail.RegionName != nil {
			detailMap["region_name"] = detail.RegionName
		}
		if detail.ZoneName != nil {
			detailMap["zone_name"] = detail.ZoneName
		}
		if detail.ResourceId != nil {
			detailMap["resource_id"] = detail.ResourceId
		}
		if detail.ResourceName != nil {
			detailMap["resource_name"] = detail.ResourceName
		}
		if detail.ActionTypeName != nil {
			detailMap["action_type_name"] = detail.ActionTypeName
		}
		if detail.OrderId != nil {
			detailMap["order_id"] = detail.OrderId
		}
		if detail.BillId != nil {
			detailMap["bill_id"] = detail.BillId
		}
		if detail.PayTime != nil {
			detailMap["pay_time"] = detail.PayTime
		}
		if detail.FeeBeginTime != nil {
			detailMap["fee_begin_time"] = detail.FeeBeginTime
		}
		if detail.FeeEndTime != nil {
			detailMap["fee_end_time"] = detail.FeeEndTime
		}
		if detail.ComponentSet != nil {
			componentSetList := make([]map[string]interface{}, 0, len(detail.ComponentSet))
			for _, component := range detail.ComponentSet {
				componentMap := map[string]interface{}{}
				if component.ComponentCodeName != nil {
					componentMap["component_code_name"] = component.ComponentCodeName
				}
				if component.ItemCodeName != nil {
					componentMap["item_code_name"] = component.ItemCodeName
				}
				if component.SinglePrice != nil {
					componentMap["single_price"] = component.SinglePrice
				}
				if component.SpecifiedPrice != nil {
					componentMap["specified_price"] = component.SpecifiedPrice
				}
				if component.PriceUnit != nil {
					componentMap["price_unit"] = component.PriceUnit
				}
				if component.UsedAmount != nil {
					componentMap["used_amount"] = component.UsedAmount
				}
				if component.UsedAmountUnit != nil {
					componentMap["used_amount_unit"] = component.UsedAmountUnit
				}
				if component.RealTotalMeasure != nil {
					componentMap["real_total_measure"] = component.RealTotalMeasure
				}
				if component.DeductedMeasure != nil {
					componentMap["deducted_measure"] = component.DeductedMeasure
				}
				if component.TimeSpan != nil {
					componentMap["time_span"] = component.TimeSpan
				}
				if component.TimeUnitName != nil {
					componentMap["time_unit_name"] = component.TimeUnitName
				}
				if component.Cost != nil {
					componentMap["cost"] = component.Cost
				}
				if component.Discount != nil {
					componentMap["discount"] = component.Discount
				}
				if component.ReduceType != nil {
					componentMap["reduce_type"] = component.ReduceType
				}
				if component.RealCost != nil {
					componentMap["real_cost"] = component.RealCost
				}
				if component.VoucherPayAmount != nil {
					componentMap["voucher_pay_amount"] = component.VoucherPayAmount
				}
				if component.CashPayAmount != nil {
					componentMap["cash_pay_amount"] = component.CashPayAmount
				}
				if component.IncentivePayAmount != nil {
					componentMap["incentive_pay_amount"] = component.IncentivePayAmount
				}
				if component.TransferPayAmount != nil {
					componentMap["transfer_pay_amount"] = component.TransferPayAmount
				}
				if component.ItemCode != nil {
					componentMap["item_code"] = component.ItemCode
				}
				if component.ComponentCode != nil {
					componentMap["component_code"] = component.ComponentCode
				}
				if component.ContractPrice != nil {
					componentMap["contract_price"] = component.ContractPrice
				}
				if component.InstanceType != nil {
					componentMap["instance_type"] = component.InstanceType
				}
				if component.RiTimeSpan != nil {
					componentMap["ri_time_span"] = component.RiTimeSpan
				}
				if component.OriginalCostWithRI != nil {
					componentMap["original_cost_with_ri"] = component.OriginalCostWithRI
				}
				if component.SPDeductionRate != nil {
					componentMap["sp_deduction_rate"] = component.SPDeductionRate
				}
				if component.SPDeduction != nil {
					componentMap["sp_deduction"] = component.SPDeduction
				}
				if component.OriginalCostWithSP != nil {
					componentMap["original_cost_with_sp"] = component.OriginalCostWithSP
				}
				if component.BlendedDiscount != nil {
					componentMap["blended_discount"] = component.BlendedDiscount
				}
				if component.ComponentConfig != nil {
					componentConfigList := make([]map[string]interface{}, 0, len(component.ComponentConfig))
					for _, config := range component.ComponentConfig {
						configMap := map[string]interface{}{}
						if config.Name != nil {
							configMap["name"] = config.Name
						}
						if config.Value != nil {
							configMap["value"] = config.Value
						}
						componentConfigList = append(componentConfigList, configMap)
					}
					componentMap["component_config"] = componentConfigList
				}
				componentSetList = append(componentSetList, componentMap)
			}
			detailMap["component_set"] = componentSetList
		}
		if detail.PayerUin != nil {
			detailMap["payer_uin"] = detail.PayerUin
		}
		if detail.OwnerUin != nil {
			detailMap["owner_uin"] = detail.OwnerUin
		}
		if detail.OperateUin != nil {
			detailMap["operate_uin"] = detail.OperateUin
		}
		if detail.Tags != nil {
			tagsList := make([]map[string]interface{}, 0, len(detail.Tags))
			for _, tag := range detail.Tags {
				tagMap := map[string]interface{}{}
				if tag.TagKey != nil {
					tagMap["tag_key"] = tag.TagKey
				}
				if tag.TagValue != nil {
					tagMap["tag_value"] = tag.TagValue
				}
				tagsList = append(tagsList, tagMap)
			}
			detailMap["tags"] = tagsList
		}
		if detail.BusinessCode != nil {
			detailMap["business_code"] = detail.BusinessCode
		}
		if detail.ProductCode != nil {
			detailMap["product_code"] = detail.ProductCode
		}
		if detail.ActionType != nil {
			detailMap["action_type"] = detail.ActionType
		}
		if detail.RegionId != nil {
			detailMap["region_id"] = detail.RegionId
		}
		if detail.ProjectId != nil {
			detailMap["project_id"] = detail.ProjectId
		}
		if detail.PriceInfo != nil {
			priceInfoList := make([]string, 0, len(detail.PriceInfo))
			for _, priceInfo := range detail.PriceInfo {
				if priceInfo != nil {
					priceInfoList = append(priceInfoList, *priceInfo)
				}
			}
			detailMap["price_info"] = priceInfoList
		}
		if detail.AssociatedOrder != nil {
			associatedOrderList := make([]map[string]interface{}, 0, 1)
			associatedOrderMap := map[string]interface{}{}
			if detail.AssociatedOrder.PrepayPurchase != nil {
				associatedOrderMap["prepay_purchase"] = detail.AssociatedOrder.PrepayPurchase
			}
			if detail.AssociatedOrder.PrepayRenew != nil {
				associatedOrderMap["prepay_renew"] = detail.AssociatedOrder.PrepayRenew
			}
			if detail.AssociatedOrder.PrepayModifyUp != nil {
				associatedOrderMap["prepay_modify_up"] = detail.AssociatedOrder.PrepayModifyUp
			}
			if detail.AssociatedOrder.ReverseOrder != nil {
				associatedOrderMap["reverse_order"] = detail.AssociatedOrder.ReverseOrder
			}
			if detail.AssociatedOrder.NewOrder != nil {
				associatedOrderMap["new_order"] = detail.AssociatedOrder.NewOrder
			}
			if detail.AssociatedOrder.Original != nil {
				associatedOrderMap["original"] = detail.AssociatedOrder.Original
			}
			associatedOrderList = append(associatedOrderList, associatedOrderMap)
			detailMap["associated_order"] = associatedOrderList
		}
		if detail.Formula != nil {
			detailMap["formula"] = detail.Formula
		}
		if detail.FormulaUrl != nil {
			detailMap["formula_url"] = detail.FormulaUrl
		}
		if detail.BillDay != nil {
			detailMap["bill_day"] = detail.BillDay
		}
		if detail.BillMonth != nil {
			detailMap["bill_month"] = detail.BillMonth
		}
		if detail.Id != nil {
			detailMap["id"] = detail.Id
		}
		if detail.RegionType != nil {
			detailMap["region_type"] = detail.RegionType
		}
		if detail.RegionTypeName != nil {
			detailMap["region_type_name"] = detail.RegionTypeName
		}
		if detail.ReserveDetail != nil {
			detailMap["reserve_detail"] = detail.ReserveDetail
		}
		if detail.DiscountObject != nil {
			detailMap["discount_object"] = detail.DiscountObject
		}
		if detail.DiscountType != nil {
			detailMap["discount_type"] = detail.DiscountType
		}
		if detail.DiscountContent != nil {
			detailMap["discount_content"] = detail.DiscountContent
		}
		detailSetList = append(detailSetList, detailMap)
	}

	_ = d.Set("detail_set", detailSetList)
	_ = d.Set("total", total)
	if respContext != nil {
		_ = d.Set("context", respContext)
	}

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), detailSetList); e != nil {
			return e
		}
	}

	return nil
}
