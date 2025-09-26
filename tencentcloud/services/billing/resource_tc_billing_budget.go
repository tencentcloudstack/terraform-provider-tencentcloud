package billing

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	billingv20180709 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBillingBudget() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBillingBudgetCreate,
		Read:   resourceTencentCloudBillingBudgetRead,
		Update: resourceTencentCloudBillingBudgetUpdate,
		Delete: resourceTencentCloudBillingBudgetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"budget_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Budget name.",
			},

			"cycle_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cycle type, valid values: DAY, MONTH, QUARTER, YEAR.",
			},

			"period_begin": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Valid period starting time 2025-01-01(cycle: days) / 2025-01 (cycle: months).",
			},

			"period_end": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Expiration period end time 2025-12-01(cycle: days) / 2025-12 (cycle: months).",
			},

			"plan_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "FIX: fixed budget, CYCLE: planned budget.",
			},

			"budget_quota": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Budget value limit. Transfer fixed value when the budget plan type is FIX(Fixed Budget); Passed when the budget plan type is CYCLE(Planned Budget)[{\"dateDesc\":\"2025-07\",\"quota\":\"1000\"},{\"dateDesc\":\"2025-08\",\"quota\":\"2000\"}].",
			},

			"bill_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "BILL: system bill, CONSUMPTION: consumption bill.",
			},

			"fee_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "COST original price, REAL_COST actual cost, CASH cash, INCENTIVE gift, VOUCHER voucher, TRANSFER share, TAX tax, AMOUNT_BEFORE_TAX cash payment (before tax).",
			},

			"warn_json": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Threshold reminder.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"warn_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ACTUAL: actual amount, FORECAST: forecast amount.",
						},
						"cal_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "PERCENTAGE: Percentage of budget amount, ABS: fixed value.",
						},
						"threshold_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Threshold (greater than or equal to 0).",
						},
					},
				},
			},

			"budget_note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Budget remarks.",
			},

			"dimensions_range": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Budget dimension range conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"business": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Products.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"pay_mode": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Pay mode.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"product_codes": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Sub-product.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"component_codes": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Component codes.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"zone_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Zone ids.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"region_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Region ids.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"project_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Project ids.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"action_types": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Action types.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"consumption_types": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Consumption types.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Tag value.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"payer_uins": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Payer uins.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"owner_uins": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Owner uins.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tree_node_uniq_keys": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Unique key for end-level ledger unit.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"wave_threshold_json": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Volatility reminder.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"warn_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ACTUAL: actual amount, FORECAST: forecast amount.",
						},
						"threshold": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Volatility threshold (greater than or equal to 0).",
						},
						"meta_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alarm type: chain month-on-month, yoy year-on-year, fix fixed value\n (Supported types: daily month-on-month chain day, daily month-on-year chain weekday, daily month-on-year monthly month-on-year fixed value fix day, month-on-month chain month, monthly fixed value fix month).",
						},
						"period_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alarm dimension: day day, month month, weekday week\n (Support types: day-to-day chain day, day-to-year weekly dimension chain weekday, day-to-year monthly dimension yoy day, daily fixed value fix day, month-to-month chain month, monthly fixed value fix month).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudBillingBudgetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_budget.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		budgetId string
	)
	var (
		request  = billingv20180709.NewCreateBudgetRequest()
		response = billingv20180709.NewCreateBudgetResponse()
	)

	if v, ok := d.GetOk("budget_name"); ok {
		request.BudgetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cycle_type"); ok {
		request.CycleType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("period_begin"); ok {
		request.PeriodBegin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("period_end"); ok {
		request.PeriodEnd = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plan_type"); ok {
		request.PlanType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("budget_quota"); ok {
		request.BudgetQuota = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bill_type"); ok {
		request.BillType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fee_type"); ok {
		request.FeeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("warn_json"); ok {
		for _, item := range v.([]interface{}) {
			warnJsonMap := item.(map[string]interface{})
			budgetWarn := billingv20180709.BudgetWarn{}
			if v, ok := warnJsonMap["warn_type"]; ok {
				budgetWarn.WarnType = helper.String(v.(string))
			}
			if v, ok := warnJsonMap["cal_type"]; ok {
				budgetWarn.CalType = helper.String(v.(string))
			}
			if v, ok := warnJsonMap["threshold_value"]; ok {
				budgetWarn.ThresholdValue = helper.String(v.(string))
			}
			request.WarnJson = append(request.WarnJson, &budgetWarn)
		}
	}

	if v, ok := d.GetOk("budget_note"); ok {
		request.BudgetNote = helper.String(v.(string))
	}

	if dimensionsRangeMap, ok := helper.InterfacesHeadMap(d, "dimensions_range"); ok {
		budgetConditionsForm := billingv20180709.BudgetConditionsForm{}
		if v, ok := dimensionsRangeMap["business"]; ok {
			businessSet := v.(*schema.Set).List()
			for i := range businessSet {
				business := businessSet[i].(string)
				budgetConditionsForm.Business = append(budgetConditionsForm.Business, helper.String(business))
			}
		}
		if v, ok := dimensionsRangeMap["pay_mode"]; ok {
			payModeSet := v.(*schema.Set).List()
			for i := range payModeSet {
				payMode := payModeSet[i].(string)
				budgetConditionsForm.PayMode = append(budgetConditionsForm.PayMode, helper.String(payMode))
			}
		}
		if v, ok := dimensionsRangeMap["product_codes"]; ok {
			productCodesSet := v.(*schema.Set).List()
			for i := range productCodesSet {
				productCodes := productCodesSet[i].(string)
				budgetConditionsForm.ProductCodes = append(budgetConditionsForm.ProductCodes, helper.String(productCodes))
			}
		}
		if v, ok := dimensionsRangeMap["component_codes"]; ok {
			componentCodesSet := v.(*schema.Set).List()
			for i := range componentCodesSet {
				componentCodes := componentCodesSet[i].(string)
				budgetConditionsForm.ComponentCodes = append(budgetConditionsForm.ComponentCodes, helper.String(componentCodes))
			}
		}
		if v, ok := dimensionsRangeMap["zone_ids"]; ok {
			zoneIdsSet := v.(*schema.Set).List()
			for i := range zoneIdsSet {
				zoneIds := zoneIdsSet[i].(string)
				budgetConditionsForm.ZoneIds = append(budgetConditionsForm.ZoneIds, helper.String(zoneIds))
			}
		}
		if v, ok := dimensionsRangeMap["region_ids"]; ok {
			regionIdsSet := v.(*schema.Set).List()
			for i := range regionIdsSet {
				regionIds := regionIdsSet[i].(string)
				budgetConditionsForm.RegionIds = append(budgetConditionsForm.RegionIds, helper.String(regionIds))
			}
		}
		if v, ok := dimensionsRangeMap["project_ids"]; ok {
			projectIdsSet := v.(*schema.Set).List()
			for i := range projectIdsSet {
				projectIds := projectIdsSet[i].(string)
				budgetConditionsForm.ProjectIds = append(budgetConditionsForm.ProjectIds, helper.String(projectIds))
			}
		}
		if v, ok := dimensionsRangeMap["action_types"]; ok {
			actionTypesSet := v.(*schema.Set).List()
			for i := range actionTypesSet {
				actionTypes := actionTypesSet[i].(string)
				budgetConditionsForm.ActionTypes = append(budgetConditionsForm.ActionTypes, helper.String(actionTypes))
			}
		}
		if v, ok := dimensionsRangeMap["consumption_types"]; ok {
			consumptionTypesSet := v.(*schema.Set).List()
			for i := range consumptionTypesSet {
				consumptionTypes := consumptionTypesSet[i].(string)
				budgetConditionsForm.ConsumptionTypes = append(budgetConditionsForm.ConsumptionTypes, helper.String(consumptionTypes))
			}
		}
		if v, ok := dimensionsRangeMap["tags"]; ok {
			for _, item := range v.([]interface{}) {
				tagsMap := item.(map[string]interface{})
				tagsForm := billingv20180709.TagsForm{}
				if v, ok := tagsMap["tag_key"]; ok {
					tagsForm.TagKey = helper.String(v.(string))
				}
				if v, ok := tagsMap["tag_value"]; ok {
					tagValueSet := v.(*schema.Set).List()
					for i := range tagValueSet {
						tagValue := tagValueSet[i].(string)
						tagsForm.TagValue = append(tagsForm.TagValue, helper.String(tagValue))
					}
				}
				budgetConditionsForm.Tags = append(budgetConditionsForm.Tags, &tagsForm)
			}
		}
		if v, ok := dimensionsRangeMap["payer_uins"]; ok {
			payerUinsSet := v.(*schema.Set).List()
			for i := range payerUinsSet {
				payerUins := payerUinsSet[i].(string)
				budgetConditionsForm.PayerUins = append(budgetConditionsForm.PayerUins, helper.String(payerUins))
			}
		}
		if v, ok := dimensionsRangeMap["owner_uins"]; ok {
			ownerUinsSet := v.(*schema.Set).List()
			for i := range ownerUinsSet {
				ownerUins := ownerUinsSet[i].(string)
				budgetConditionsForm.OwnerUins = append(budgetConditionsForm.OwnerUins, helper.String(ownerUins))
			}
		}
		if v, ok := dimensionsRangeMap["tree_node_uniq_keys"]; ok {
			treeNodeUniqKeysSet := v.(*schema.Set).List()
			for i := range treeNodeUniqKeysSet {
				treeNodeUniqKeys := treeNodeUniqKeysSet[i].(string)
				budgetConditionsForm.TreeNodeUniqKeys = append(budgetConditionsForm.TreeNodeUniqKeys, helper.String(treeNodeUniqKeys))
			}
		}
		request.DimensionsRange = &budgetConditionsForm
	}

	if v, ok := d.GetOk("wave_threshold_json"); ok {
		for _, item := range v.([]interface{}) {
			waveThresholdJsonMap := item.(map[string]interface{})
			waveThresholdForm := billingv20180709.WaveThresholdForm{}
			if v, ok := waveThresholdJsonMap["warn_type"]; ok {
				waveThresholdForm.WarnType = helper.String(v.(string))
			}
			if v, ok := waveThresholdJsonMap["threshold"]; ok {
				waveThresholdForm.Threshold = helper.String(v.(string))
			}
			if v, ok := waveThresholdJsonMap["meta_type"]; ok {
				waveThresholdForm.MetaType = helper.String(v.(string))
			}
			if v, ok := waveThresholdJsonMap["period_type"]; ok {
				waveThresholdForm.PeriodType = helper.String(v.(string))
			}
			request.WaveThresholdJson = append(request.WaveThresholdJson, &waveThresholdForm)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().CreateBudgetWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create billing budget failed, reason:%+v", logId, err)
		return err
	}

	if response.Response != nil && response.Response.Data.BudgetId != nil {
		budgetId = *response.Response.Data.BudgetId
	}

	d.SetId(budgetId)

	return resourceTencentCloudBillingBudgetRead(d, meta)
}

func resourceTencentCloudBillingBudgetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_budget.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := BillingService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	budgetId := d.Id()

	respData, err := service.DescribeBillingBudgetById(ctx, budgetId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `billing_budget` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.Data == nil || len(respData.Data.Records) == 0 {
		d.SetId("")
		return nil
	}
	record := respData.Data.Records[0]

	if record.BudgetName != nil {
		_ = d.Set("budget_name", record.BudgetName)
	}
	if record.CycleType != nil {
		_ = d.Set("cycle_type", record.CycleType)
	}

	if record.PeriodBegin != nil {
		_ = d.Set("period_begin", record.PeriodBegin)
	}

	if record.PeriodEnd != nil {
		_ = d.Set("period_end", record.PeriodEnd)
	}

	if record.PlanType != nil {
		_ = d.Set("plan_type", record.PlanType)
	}

	if record.BudgetQuota != nil {
		_ = d.Set("budget_quota", record.BudgetQuota)
	}

	if record.BillType != nil {
		_ = d.Set("bill_type", record.BillType)
	}

	if record.FeeType != nil {
		_ = d.Set("fee_type", record.FeeType)
	}

	warnJsonList := make([]map[string]interface{}, 0, len(record.WarnJson))
	if record.WarnJson != nil {
		for _, warnJson := range record.WarnJson {
			warnJsonMap := map[string]interface{}{}

			if warnJson.WarnType != nil {
				warnJsonMap["warn_type"] = warnJson.WarnType
			}

			if warnJson.CalType != nil {
				warnJsonMap["cal_type"] = warnJson.CalType
			}

			if warnJson.ThresholdValue != nil {
				warnJsonMap["threshold_value"] = warnJson.ThresholdValue
			}

			warnJsonList = append(warnJsonList, warnJsonMap)
		}

		_ = d.Set("warn_json", warnJsonList)
	}

	if record.BudgetNote != nil {
		_ = d.Set("budget_note", record.BudgetNote)
	}

	waveThresholdJsonList := make([]map[string]interface{}, 0, len(record.WaveThresholdJson))
	if record.WaveThresholdJson != nil {
		for _, waveThresholdJson := range record.WaveThresholdJson {
			waveThresholdJsonMap := map[string]interface{}{}

			if waveThresholdJson.WarnType != nil {
				waveThresholdJsonMap["warn_type"] = waveThresholdJson.WarnType
			}

			if waveThresholdJson.Threshold != nil {
				waveThresholdJsonMap["threshold"] = waveThresholdJson.Threshold
			}

			if waveThresholdJson.MetaType != nil {
				waveThresholdJsonMap["meta_type"] = waveThresholdJson.MetaType
			}

			if waveThresholdJson.PeriodType != nil {
				waveThresholdJsonMap["period_type"] = waveThresholdJson.PeriodType
			}

			waveThresholdJsonList = append(waveThresholdJsonList, waveThresholdJsonMap)
		}

		_ = d.Set("wave_threshold_json", waveThresholdJsonList)
	}

	if record.DimensionsRange != nil {
		dimensionsRangeMap := make(map[string]interface{})
		if record.DimensionsRange.Business != nil {
			dimensionsRangeMap["business"] = record.DimensionsRange.Business
		}

		if record.DimensionsRange.PayMode != nil {
			dimensionsRangeMap["pay_mode"] = record.DimensionsRange.PayMode
		}

		if record.DimensionsRange.ProductCodes != nil {
			dimensionsRangeMap["product_codes"] = record.DimensionsRange.ProductCodes
		}

		if record.DimensionsRange.ComponentCodes != nil {
			dimensionsRangeMap["component_codes"] = record.DimensionsRange.ComponentCodes
		}

		if record.DimensionsRange.ZoneIds != nil {
			dimensionsRangeMap["zone_ids"] = record.DimensionsRange.ZoneIds
		}

		if record.DimensionsRange.RegionIds != nil {
			dimensionsRangeMap["region_ids"] = record.DimensionsRange.RegionIds
		}

		if record.DimensionsRange.ProjectIds != nil {
			dimensionsRangeMap["project_ids"] = record.DimensionsRange.ProjectIds
		}

		if record.DimensionsRange.ActionTypes != nil {
			dimensionsRangeMap["action_types"] = record.DimensionsRange.ActionTypes
		}

		if record.DimensionsRange.ConsumptionTypes != nil {
			dimensionsRangeMap["consumption_types"] = record.DimensionsRange.ConsumptionTypes
		}

		tagsList := make([]map[string]interface{}, 0, len(record.DimensionsRange.Tags))
		if record.DimensionsRange.Tags != nil {
			for _, tags := range record.DimensionsRange.Tags {
				tagsMap := map[string]interface{}{}

				if tags.TagKey != nil {
					tagsMap["tag_key"] = tags.TagKey
				}

				if tags.TagValue != nil {
					tagsMap["tag_value"] = tags.TagValue
				}

				tagsList = append(tagsList, tagsMap)
			}

			dimensionsRangeMap["tags"] = tagsList
		}
		if record.DimensionsRange.PayerUins != nil {
			dimensionsRangeMap["payer_uins"] = record.DimensionsRange.PayerUins
		}

		if record.DimensionsRange.OwnerUins != nil {
			dimensionsRangeMap["owner_uins"] = record.DimensionsRange.OwnerUins
		}

		if record.DimensionsRange.TreeNodeUniqKeys != nil {
			dimensionsRangeMap["tree_node_uniq_keys"] = record.DimensionsRange.TreeNodeUniqKeys
		}
		hasNotNullItem := false
		for _, v := range dimensionsRangeMap {
			if v != nil {
				hasNotNullItem = true
			}
		}
		if hasNotNullItem {
			_ = d.Set("dimensions_range", []interface{}{dimensionsRangeMap})
		} else {
			_ = d.Set("dimensions_range", []interface{}{})
		}

	}

	return nil
}

func resourceTencentCloudBillingBudgetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_budget.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	budgetId := d.Id()

	needChange := false
	mutableArgs := []string{"budget_id", "budget_name", "cycle_type", "period_begin", "period_end", "plan_type", "budget_quota", "bill_type", "fee_type", "warn_json", "budget_note", "dimensions_range", "wave_threshold_json"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := billingv20180709.NewModifyBudgetRequest()
		request.BudgetId = helper.String(budgetId)
		if v, ok := d.GetOk("budget_id"); ok {
			request.BudgetId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("budget_name"); ok {
			request.BudgetName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cycle_type"); ok {
			request.CycleType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("period_begin"); ok {
			request.PeriodBegin = helper.String(v.(string))
		}

		if v, ok := d.GetOk("period_end"); ok {
			request.PeriodEnd = helper.String(v.(string))
		}

		if v, ok := d.GetOk("plan_type"); ok {
			request.PlanType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("budget_quota"); ok {
			request.BudgetQuota = helper.String(v.(string))
		}

		if v, ok := d.GetOk("bill_type"); ok {
			request.BillType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("fee_type"); ok {
			request.FeeType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("warn_json"); ok {
			for _, item := range v.([]interface{}) {
				warnJsonMap := item.(map[string]interface{})
				budgetWarn := billingv20180709.BudgetWarn{}
				if v, ok := warnJsonMap["warn_type"]; ok {
					budgetWarn.WarnType = helper.String(v.(string))
				}
				if v, ok := warnJsonMap["cal_type"]; ok {
					budgetWarn.CalType = helper.String(v.(string))
				}
				if v, ok := warnJsonMap["threshold_value"]; ok {
					budgetWarn.ThresholdValue = helper.String(v.(string))
				}
				request.WarnJson = append(request.WarnJson, &budgetWarn)
			}
		}

		if v, ok := d.GetOk("budget_note"); ok {
			request.BudgetNote = helper.String(v.(string))
		}

		if dimensionsRangeMap, ok := helper.InterfacesHeadMap(d, "dimensions_range"); ok {
			budgetConditionsForm := billingv20180709.BudgetConditionsForm{}
			if v, ok := dimensionsRangeMap["business"]; ok {
				businessSet := v.(*schema.Set).List()
				for i := range businessSet {
					business := businessSet[i].(string)
					budgetConditionsForm.Business = append(budgetConditionsForm.Business, helper.String(business))
				}
			}
			if v, ok := dimensionsRangeMap["pay_mode"]; ok {
				payModeSet := v.(*schema.Set).List()
				for i := range payModeSet {
					payMode := payModeSet[i].(string)
					budgetConditionsForm.PayMode = append(budgetConditionsForm.PayMode, helper.String(payMode))
				}
			}
			if v, ok := dimensionsRangeMap["product_codes"]; ok {
				productCodesSet := v.(*schema.Set).List()
				for i := range productCodesSet {
					productCodes := productCodesSet[i].(string)
					budgetConditionsForm.ProductCodes = append(budgetConditionsForm.ProductCodes, helper.String(productCodes))
				}
			}
			if v, ok := dimensionsRangeMap["component_codes"]; ok {
				componentCodesSet := v.(*schema.Set).List()
				for i := range componentCodesSet {
					componentCodes := componentCodesSet[i].(string)
					budgetConditionsForm.ComponentCodes = append(budgetConditionsForm.ComponentCodes, helper.String(componentCodes))
				}
			}
			if v, ok := dimensionsRangeMap["zone_ids"]; ok {
				zoneIdsSet := v.(*schema.Set).List()
				for i := range zoneIdsSet {
					zoneIds := zoneIdsSet[i].(string)
					budgetConditionsForm.ZoneIds = append(budgetConditionsForm.ZoneIds, helper.String(zoneIds))
				}
			}
			if v, ok := dimensionsRangeMap["region_ids"]; ok {
				regionIdsSet := v.(*schema.Set).List()
				for i := range regionIdsSet {
					regionIds := regionIdsSet[i].(string)
					budgetConditionsForm.RegionIds = append(budgetConditionsForm.RegionIds, helper.String(regionIds))
				}
			}
			if v, ok := dimensionsRangeMap["project_ids"]; ok {
				projectIdsSet := v.(*schema.Set).List()
				for i := range projectIdsSet {
					projectIds := projectIdsSet[i].(string)
					budgetConditionsForm.ProjectIds = append(budgetConditionsForm.ProjectIds, helper.String(projectIds))
				}
			}
			if v, ok := dimensionsRangeMap["action_types"]; ok {
				actionTypesSet := v.(*schema.Set).List()
				for i := range actionTypesSet {
					actionTypes := actionTypesSet[i].(string)
					budgetConditionsForm.ActionTypes = append(budgetConditionsForm.ActionTypes, helper.String(actionTypes))
				}
			}
			if v, ok := dimensionsRangeMap["consumption_types"]; ok {
				consumptionTypesSet := v.(*schema.Set).List()
				for i := range consumptionTypesSet {
					consumptionTypes := consumptionTypesSet[i].(string)
					budgetConditionsForm.ConsumptionTypes = append(budgetConditionsForm.ConsumptionTypes, helper.String(consumptionTypes))
				}
			}
			if v, ok := dimensionsRangeMap["tags"]; ok {
				for _, item := range v.([]interface{}) {
					tagsMap := item.(map[string]interface{})
					tagsForm := billingv20180709.TagsForm{}
					if v, ok := tagsMap["tag_key"]; ok {
						tagsForm.TagKey = helper.String(v.(string))
					}
					if v, ok := tagsMap["tag_value"]; ok {
						tagValueSet := v.(*schema.Set).List()
						for i := range tagValueSet {
							tagValue := tagValueSet[i].(string)
							tagsForm.TagValue = append(tagsForm.TagValue, helper.String(tagValue))
						}
					}
					budgetConditionsForm.Tags = append(budgetConditionsForm.Tags, &tagsForm)
				}
			}
			if v, ok := dimensionsRangeMap["payer_uins"]; ok {
				payerUinsSet := v.(*schema.Set).List()
				for i := range payerUinsSet {
					payerUins := payerUinsSet[i].(string)
					budgetConditionsForm.PayerUins = append(budgetConditionsForm.PayerUins, helper.String(payerUins))
				}
			}
			if v, ok := dimensionsRangeMap["owner_uins"]; ok {
				ownerUinsSet := v.(*schema.Set).List()
				for i := range ownerUinsSet {
					ownerUins := ownerUinsSet[i].(string)
					budgetConditionsForm.OwnerUins = append(budgetConditionsForm.OwnerUins, helper.String(ownerUins))
				}
			}
			if v, ok := dimensionsRangeMap["tree_node_uniq_keys"]; ok {
				treeNodeUniqKeysSet := v.(*schema.Set).List()
				for i := range treeNodeUniqKeysSet {
					treeNodeUniqKeys := treeNodeUniqKeysSet[i].(string)
					budgetConditionsForm.TreeNodeUniqKeys = append(budgetConditionsForm.TreeNodeUniqKeys, helper.String(treeNodeUniqKeys))
				}
			}
			request.DimensionsRange = &budgetConditionsForm
		}

		if v, ok := d.GetOk("wave_threshold_json"); ok {
			for _, item := range v.([]interface{}) {
				waveThresholdJsonMap := item.(map[string]interface{})
				waveThresholdForm := billingv20180709.WaveThresholdForm{}
				if v, ok := waveThresholdJsonMap["warn_type"]; ok {
					waveThresholdForm.WarnType = helper.String(v.(string))
				}
				if v, ok := waveThresholdJsonMap["threshold"]; ok {
					waveThresholdForm.Threshold = helper.String(v.(string))
				}
				if v, ok := waveThresholdJsonMap["meta_type"]; ok {
					waveThresholdForm.MetaType = helper.String(v.(string))
				}
				if v, ok := waveThresholdJsonMap["period_type"]; ok {
					waveThresholdForm.PeriodType = helper.String(v.(string))
				}
				request.WaveThresholdJson = append(request.WaveThresholdJson, &waveThresholdForm)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().ModifyBudgetWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update billing budget failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudBillingBudgetRead(d, meta)
}

func resourceTencentCloudBillingBudgetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_budget.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	budgetId := d.Id()

	var (
		request  = billingv20180709.NewDeleteBudgetRequest()
		response = billingv20180709.NewDeleteBudgetResponse()
	)

	request.BudgetIds = helper.Strings([]string{budgetId})
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().DeleteBudgetWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete billing budget failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
