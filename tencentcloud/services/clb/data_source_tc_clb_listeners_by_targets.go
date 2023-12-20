package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbListenersByTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbListenersByTargetsRead,
		Schema: map[string]*schema.Schema{
			"backends": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "List of private network IPs to be queried.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID.",
						},
						"private_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Private network IP to be queried, which can be of the CVM or ENI.",
						},
					},
				},
			},

			"load_balancers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Detail of the CLB instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "String ID of the CLB instance.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VIP of the CLB instance.",
						},
						"listeners": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Listener rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener ID.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener protocol.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Listener port.",
									},
									"rules": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Bound rule. Note: this field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"location_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Rule ID.",
												},
												"domain": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Domain name.",
												},
												"url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "url.",
												},
												"targets": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Object bound to the real server.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Private network IP type, which can be cvm or eni.",
															},
															"private_ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Private network IP of the real server.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Port bound to the real server.",
															},
															"vpc_id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "VPC ID of the real server. Note: this field may return null, indicating that no valid values can be obtained.",
															},
															"weight": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Weight of the real server. Note: this field may return null, indicating that no valid values can be obtained.",
															},
														},
													},
												},
											},
										},
									},
									"targets": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Object bound to the layer-4 listener. Note: this field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Private network IP type, which can be cvm or eni.",
												},
												"private_ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Private network IP of the real server.",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Port bound to the real server.",
												},
												"vpc_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "VPC ID of the real server. Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"weight": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Weight of the real server. Note: this field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"end_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "End port of the listener. Note: this field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance region.",
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

func dataSourceTencentCloudClbListenersByTargetsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_listeners_by_targets.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("backends"); ok {
		backendsSet := v.([]interface{})
		tmpSet := make([]*clb.LbRsItem, 0, len(backendsSet))

		for _, item := range backendsSet {
			lbRsItem := clb.LbRsItem{}
			lbRsItemMap := item.(map[string]interface{})

			if v, ok := lbRsItemMap["vpc_id"]; ok {
				lbRsItem.VpcId = helper.String(v.(string))
			}
			if v, ok := lbRsItemMap["private_ip"]; ok {
				lbRsItem.PrivateIp = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &lbRsItem)
		}
		paramMap["Backends"] = tmpSet
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var loadBalancers []*clb.LBItem

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbListenersByTargets(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		loadBalancers = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(loadBalancers))
	tmpList := make([]map[string]interface{}, 0, len(loadBalancers))

	if loadBalancers != nil {
		for _, lBItem := range loadBalancers {
			lBItemMap := map[string]interface{}{}

			if lBItem.LoadBalancerId != nil {
				lBItemMap["load_balancer_id"] = lBItem.LoadBalancerId
			}

			if lBItem.Vip != nil {
				lBItemMap["vip"] = lBItem.Vip
			}

			if lBItem.Listeners != nil {
				listenersList := []interface{}{}
				for _, listeners := range lBItem.Listeners {
					listenersMap := map[string]interface{}{}

					if listeners.ListenerId != nil {
						listenersMap["listener_id"] = listeners.ListenerId
					}

					if listeners.Protocol != nil {
						listenersMap["protocol"] = listeners.Protocol
					}

					if listeners.Port != nil {
						listenersMap["port"] = listeners.Port
					}

					if listeners.Rules != nil {
						rulesList := []interface{}{}
						for _, rules := range listeners.Rules {
							rulesMap := map[string]interface{}{}

							if rules.LocationId != nil {
								rulesMap["location_id"] = rules.LocationId
							}

							if rules.Domain != nil {
								rulesMap["domain"] = rules.Domain
							}

							if rules.Url != nil {
								rulesMap["url"] = rules.Url
							}

							if rules.Targets != nil {
								targetsList := []interface{}{}
								for _, targets := range rules.Targets {
									targetsMap := map[string]interface{}{}

									if targets.Type != nil {
										targetsMap["type"] = targets.Type
									}

									if targets.PrivateIp != nil {
										targetsMap["private_ip"] = targets.PrivateIp
									}

									if targets.Port != nil {
										targetsMap["port"] = targets.Port
									}

									if targets.VpcId != nil {
										targetsMap["vpc_id"] = targets.VpcId
									}

									if targets.Weight != nil {
										targetsMap["weight"] = targets.Weight
									}

									targetsList = append(targetsList, targetsMap)
								}

								rulesMap["targets"] = targetsList
							}

							rulesList = append(rulesList, rulesMap)
						}

						listenersMap["rules"] = rulesList
					}

					if listeners.Targets != nil {
						targetsList := []interface{}{}
						for _, targets := range listeners.Targets {
							targetsMap := map[string]interface{}{}

							if targets.Type != nil {
								targetsMap["type"] = targets.Type
							}

							if targets.PrivateIp != nil {
								targetsMap["private_ip"] = targets.PrivateIp
							}

							if targets.Port != nil {
								targetsMap["port"] = targets.Port
							}

							if targets.VpcId != nil {
								targetsMap["vpc_id"] = targets.VpcId
							}

							if targets.Weight != nil {
								targetsMap["weight"] = targets.Weight
							}

							targetsList = append(targetsList, targetsMap)
						}

						listenersMap["targets"] = targetsList
					}

					if listeners.EndPort != nil {
						listenersMap["end_port"] = listeners.EndPort
					}

					listenersList = append(listenersList, listenersMap)
				}

				lBItemMap["listeners"] = listenersList
			}

			if lBItem.Region != nil {
				lBItemMap["region"] = lBItem.Region
			}

			ids = append(ids, *lBItem.LoadBalancerId)
			tmpList = append(tmpList, lBItemMap)
		}

		_ = d.Set("load_balancers", tmpList)
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
