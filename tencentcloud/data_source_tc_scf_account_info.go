package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudScfAccountInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfAccountInfoRead,
		Schema: map[string]*schema.Schema{
			"account_usage": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Namespace usage information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespaces_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of namespaces.",
						},
						"namespace": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Namespace details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"functions": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Function array.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace name.",
									},
									"functions_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of functions in namespace.",
									},
									"total_concurrency_mem": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total memory quota of the namespace Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"total_allocated_concurrency_mem": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Concurrency usage of the namespace Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"total_allocated_provisioned_mem": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Provisioned concurrency usage of the namespace Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"total_concurrency_mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Upper limit of user concurrency memory in the current region.",
						},
						"total_allocated_concurrency_mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Quota of configured user concurrency memory in the current region.",
						},
						"user_concurrency_mem_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Quota of account concurrency actually configured by user.",
						},
					},
				},
			},

			"account_limit": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Namespace limit information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespaces_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Limit of namespace quantity.",
						},
						"namespace": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Namespace limit information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"functions_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total number of functions.",
									},
									"trigger": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Trigger information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cos": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of COS triggers.",
												},
												"timer": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of timer triggers.",
												},
												"cmq": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of CMQ triggers.",
												},
												"total": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total number of triggers.",
												},
												"ckafka": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of CKafka triggers.",
												},
												"apigw": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of API Gateway triggers.",
												},
												"cls": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of CLS triggers.",
												},
												"clb": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of CLB triggers.",
												},
												"mps": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of MPS triggers.",
												},
												"cm": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of CM triggers.",
												},
												"vod": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of VOD triggers.",
												},
												"eb": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of EventBridge triggers Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"namespace": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace name.",
									},
									"concurrent_executions": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Concurrency.",
									},
									"timeout_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Timeout limit.",
									},
									"test_model_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Test event limit Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"init_timeout_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Initialization timeout limit.",
									},
									"retry_num_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Limit of async retry attempt quantity.",
									},
									"min_msg_ttl": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Lower limit of message retention time for async retry.",
									},
									"max_msg_ttl": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Upper limit of message retention time for async retry.",
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

func dataSourceTencentCloudScfAccountInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_account_info.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var accountInfo *scf.GetAccountResponseParams

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfAccountInfo(ctx)
		if e != nil {
			return retryError(e)
		}
		accountInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)

	accountInfoMap := map[string]interface{}{}

	if accountInfo.AccountUsage != nil {
		usageInfoMap := map[string]interface{}{}
		accountUsage := accountInfo.AccountUsage

		if accountUsage.NamespacesCount != nil {
			usageInfoMap["namespaces_count"] = accountUsage.NamespacesCount
		}

		if accountUsage.Namespace != nil {
			namespaceList := []interface{}{}
			for _, namespace := range accountUsage.Namespace {
				namespaceMap := map[string]interface{}{}

				if namespace.Functions != nil {
					namespaceMap["functions"] = namespace.Functions
				}

				if namespace.Namespace != nil {
					namespaceMap["namespace"] = namespace.Namespace
				}

				if namespace.FunctionsCount != nil {
					namespaceMap["functions_count"] = namespace.FunctionsCount
				}

				if namespace.TotalConcurrencyMem != nil {
					namespaceMap["total_concurrency_mem"] = namespace.TotalConcurrencyMem
				}

				if namespace.TotalAllocatedConcurrencyMem != nil {
					namespaceMap["total_allocated_concurrency_mem"] = namespace.TotalAllocatedConcurrencyMem
				}

				if namespace.TotalAllocatedProvisionedMem != nil {
					namespaceMap["total_allocated_provisioned_mem"] = namespace.TotalAllocatedProvisionedMem
				}

				namespaceList = append(namespaceList, namespaceMap)
			}

			usageInfoMap["namespace"] = namespaceList
		}

		if accountUsage.TotalConcurrencyMem != nil {
			usageInfoMap["total_concurrency_mem"] = accountUsage.TotalConcurrencyMem
		}

		if accountUsage.TotalAllocatedConcurrencyMem != nil {
			usageInfoMap["total_allocated_concurrency_mem"] = accountUsage.TotalAllocatedConcurrencyMem
		}

		if accountUsage.UserConcurrencyMemLimit != nil {
			usageInfoMap["user_concurrency_mem_limit"] = accountUsage.UserConcurrencyMemLimit
		}

		_ = d.Set("account_usage", []interface{}{usageInfoMap})
		accountInfoMap["account_usage"] = []interface{}{usageInfoMap}
	}

	if accountInfo.AccountLimit != nil {
		limitsInfoMap := map[string]interface{}{}
		accountLimit := accountInfo.AccountLimit

		if accountLimit.NamespacesCount != nil {
			limitsInfoMap["namespaces_count"] = accountLimit.NamespacesCount
		}

		if accountLimit.Namespace != nil {
			namespaceList := []interface{}{}
			for _, namespace := range accountLimit.Namespace {
				namespaceMap := map[string]interface{}{}

				if namespace.FunctionsCount != nil {
					namespaceMap["functions_count"] = namespace.FunctionsCount
				}

				if namespace.Trigger != nil {
					triggerMap := map[string]interface{}{}

					if namespace.Trigger.Cos != nil {
						triggerMap["cos"] = namespace.Trigger.Cos
					}

					if namespace.Trigger.Timer != nil {
						triggerMap["timer"] = namespace.Trigger.Timer
					}

					if namespace.Trigger.Cmq != nil {
						triggerMap["cmq"] = namespace.Trigger.Cmq
					}

					if namespace.Trigger.Total != nil {
						triggerMap["total"] = namespace.Trigger.Total
					}

					if namespace.Trigger.Ckafka != nil {
						triggerMap["ckafka"] = namespace.Trigger.Ckafka
					}

					if namespace.Trigger.Apigw != nil {
						triggerMap["apigw"] = namespace.Trigger.Apigw
					}

					if namespace.Trigger.Cls != nil {
						triggerMap["cls"] = namespace.Trigger.Cls
					}

					if namespace.Trigger.Clb != nil {
						triggerMap["clb"] = namespace.Trigger.Clb
					}

					if namespace.Trigger.Mps != nil {
						triggerMap["mps"] = namespace.Trigger.Mps
					}

					if namespace.Trigger.Cm != nil {
						triggerMap["cm"] = namespace.Trigger.Cm
					}

					if namespace.Trigger.Vod != nil {
						triggerMap["vod"] = namespace.Trigger.Vod
					}

					if namespace.Trigger.Eb != nil {
						triggerMap["eb"] = namespace.Trigger.Eb
					}

					namespaceMap["trigger"] = []interface{}{triggerMap}
				}

				if namespace.Namespace != nil {
					namespaceMap["namespace"] = namespace.Namespace
					ids = append(ids, *namespace.Namespace)
				}

				if namespace.ConcurrentExecutions != nil {
					namespaceMap["concurrent_executions"] = namespace.ConcurrentExecutions
				}

				if namespace.TimeoutLimit != nil {
					namespaceMap["timeout_limit"] = namespace.TimeoutLimit
				}

				if namespace.TestModelLimit != nil {
					namespaceMap["test_model_limit"] = namespace.TestModelLimit
				}

				if namespace.InitTimeoutLimit != nil {
					namespaceMap["init_timeout_limit"] = namespace.InitTimeoutLimit
				}

				if namespace.RetryNumLimit != nil {
					namespaceMap["retry_num_limit"] = namespace.RetryNumLimit
				}

				if namespace.MinMsgTTL != nil {
					namespaceMap["min_msg_ttl"] = namespace.MinMsgTTL
				}

				if namespace.MaxMsgTTL != nil {
					namespaceMap["max_msg_ttl"] = namespace.MaxMsgTTL
				}

				namespaceList = append(namespaceList, namespaceMap)
			}

			limitsInfoMap["namespace"] = namespaceList
		}
		_ = d.Set("account_limit", []interface{}{limitsInfoMap})
		accountInfoMap["account_limit"] = []interface{}{limitsInfoMap}
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), accountInfoMap); e != nil {
			return e
		}
	}
	return nil
}
