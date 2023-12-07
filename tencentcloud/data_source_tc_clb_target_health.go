package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClbTargetHealth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbTargetHealthRead,
		Schema: map[string]*schema.Schema{
			"load_balancer_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of IDs of CLB instances to be queried.",
			},

			"load_balancers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "CLB instance list. Note: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance ID.",
						},
						"load_balancer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance name. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"listeners": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of listeners. Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener ID.",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener name. Note: This field may return null, indicating that no valid values can be obtained.",
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
										Description: "List of forwarding rules of the listener. Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"location_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Forwarding rule ID.",
												},
												"domain": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Domain name of the forwarding rule. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Forwarding rule Url. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"targets": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Health status of the real server bound to this rule. Note: this field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Private IP of the target.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Port bound to the target.",
															},
															"health_status": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Current health status. true: healthy; false: unhealthy.",
															},
															"target_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance ID of the target, such as ins-12345678.",
															},
															"health_status_detail": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Detailed information about the current health status. Alive: healthy; Dead: exceptional; Unknown: check not started/checking/unknown status.",
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

func dataSourceTencentCloudClbTargetHealthRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_target_health.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("load_balancer_ids"); ok {
		loadBalancerIdsSet := v.(*schema.Set).List()
		paramMap["LoadBalancerIds"] = helper.InterfacesStringsPoint(loadBalancerIdsSet)
	}

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var loadBalancers []*clb.LoadBalancerHealth

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbTargetHealthByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
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
		for _, loadBalancerHealth := range loadBalancers {
			loadBalancerHealthMap := map[string]interface{}{}

			if loadBalancerHealth.LoadBalancerId != nil {
				loadBalancerHealthMap["load_balancer_id"] = loadBalancerHealth.LoadBalancerId
			}

			if loadBalancerHealth.LoadBalancerName != nil {
				loadBalancerHealthMap["load_balancer_name"] = loadBalancerHealth.LoadBalancerName
			}

			if loadBalancerHealth.Listeners != nil {
				listenersList := []interface{}{}
				for _, listeners := range loadBalancerHealth.Listeners {
					listenersMap := map[string]interface{}{}

					if listeners.ListenerId != nil {
						listenersMap["listener_id"] = listeners.ListenerId
					}

					if listeners.ListenerName != nil {
						listenersMap["listener_name"] = listeners.ListenerName
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

									if targets.IP != nil {
										targetsMap["ip"] = targets.IP
									}

									if targets.Port != nil {
										targetsMap["port"] = targets.Port
									}

									if targets.HealthStatus != nil {
										targetsMap["health_status"] = targets.HealthStatus
									}

									if targets.TargetId != nil {
										targetsMap["target_id"] = targets.TargetId
									}

									if targets.HealthStatusDetail != nil {
										targetsMap["health_status_detail"] = targets.HealthStatusDetail
									}

									targetsList = append(targetsList, targetsMap)
								}

								rulesMap["targets"] = targetsList
							}

							rulesList = append(rulesList, rulesMap)
						}

						listenersMap["rules"] = rulesList
					}

					listenersList = append(listenersList, listenersMap)
				}

				loadBalancerHealthMap["listeners"] = listenersList
			}

			ids = append(ids, *loadBalancerHealth.LoadBalancerId)
			tmpList = append(tmpList, loadBalancerHealthMap)
		}

		_ = d.Set("load_balancers", tmpList)
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
