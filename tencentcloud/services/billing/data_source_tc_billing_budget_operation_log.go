package billing

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	billingv20180709 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudBillingBudgetOperationLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudBillingBudgetOperationLogRead,
		Schema: map[string]*schema.Schema{
			"budget_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Budget id.",
			},

			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Query data list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payer_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Payer uin.",
						},
						"owner_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Owner uin.",
						},
						"operate_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Operate uin.",
						},
						"bill_day": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bill day.",
						},
						"bill_month": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bill month.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification type: ADD, UPDATE.",
						},
						"diff_value": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "change information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"property": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Change attributes.",
									},
									"before": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Content before change.",
									},
									"after": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Content after change.",
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"operation_channel": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operation channel.",
						},
						"budget_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Budget item id.",
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

func dataSourceTencentCloudBillingBudgetOperationLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_billing_budget_operation_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := BillingService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var budgetId string

	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("budget_id"); ok {
		budgetId = v.(string)
		paramMap["BudgetId"] = helper.String(budgetId)
	}

	var respData []*billingv20180709.BudgetOperationLogEntity
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBillingBudgetOperationLogByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	if len(respData) > 0 {
		recordsList := make([]map[string]interface{}, 0, len(respData))

		for _, records := range respData {
			recordsMap := map[string]interface{}{}

			if records.PayerUin != nil {
				recordsMap["payer_uin"] = records.PayerUin
			}

			if records.OwnerUin != nil {
				recordsMap["owner_uin"] = records.OwnerUin
			}

			if records.OperateUin != nil {
				recordsMap["operate_uin"] = records.OperateUin
			}

			if records.BillDay != nil {
				recordsMap["bill_day"] = records.BillDay
			}

			if records.BillMonth != nil {
				recordsMap["bill_month"] = records.BillMonth
			}

			if records.Action != nil {
				recordsMap["action"] = records.Action
			}

			diffValueList := make([]map[string]interface{}, 0, len(records.DiffValue))
			if records.DiffValue != nil {
				for _, diffValue := range records.DiffValue {
					diffValueMap := map[string]interface{}{}

					if diffValue.Property != nil {
						diffValueMap["property"] = diffValue.Property
					}

					if diffValue.Before != nil {
						diffValueMap["before"] = diffValue.Before
					}

					if diffValue.After != nil {
						diffValueMap["after"] = diffValue.After
					}

					diffValueList = append(diffValueList, diffValueMap)
				}

				recordsMap["diff_value"] = diffValueList
			}
			if records.CreateTime != nil {
				recordsMap["create_time"] = records.CreateTime
			}

			if records.UpdateTime != nil {
				recordsMap["update_time"] = records.UpdateTime
			}

			if records.OperationChannel != nil {
				recordsMap["operation_channel"] = records.OperationChannel
			}

			if records.BudgetId != nil {
				recordsMap["budget_id"] = records.BudgetId
			}

			recordsList = append(recordsList, recordsMap)
		}

		_ = d.Set("records", recordsList)
	}

	d.SetId(budgetId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
