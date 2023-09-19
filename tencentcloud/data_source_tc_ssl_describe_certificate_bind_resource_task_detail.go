/*
Use this data source to query detailed information of ssl describe_certificate_bind_resource_task_detail

Example Usage

```hcl
data "tencentcloud_ssl_describe_certificate_bind_resource_task_detail" "describe_certificate_bind_resource_task_detail" {
  task_id = ""
  resource_types =
  regions =
                        }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslDescribeCertificateBindResourceTaskDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeCertificateBindResourceTaskDetailRead,
		Schema: map[string]*schema.Schema{
			"task_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task ID, query the results of cloud resources based on the task ID query.",
			},

			"resource_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Query the results of the type of resource type.",
			},

			"regions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Query the data of the regional list, CLB, TKE, WAF, Apigateway, TCB support regional inquiries, other resource types do not support.",
			},

			"clb": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Related CLB resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "area.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CLB instance detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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
										Description: "CLB instance name name.",
									},
									"listeners": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "CLB listener listNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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
													Description: "Name of listeners.",
												},
												"sni_switch": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether to turn on SNI, 1 to open, 0 to close.",
												},
												"protocol": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Type of listener protocol, https | TCP_SSL.",
												},
												"certificate": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Certificate data binding of listenersNote: This field may return NULL, indicating that the valid value cannot be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cert_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Certificate ID.",
															},
															"dns_names": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "Domain name binding of certificates.",
															},
															"cert_ca_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Root certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
															"s_s_l_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Certificate certification mode: unidirectional unidirectional authentication, Mutual two -way certificationNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
														},
													},
												},
												"rules": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of listeners&#39; rulesNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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
																Description: "Domain name binding.",
															},
															"is_match": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the rules match the domain name to be bound to the certificate.",
															},
															"certificate": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Certificate data that has been bound to the rulesNote: This field may return NULL, indicating that the valid value cannot be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cert_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Certificate ID.",
																		},
																		"dns_names": {
																			Type: schema.TypeSet,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																			Computed:    true,
																			Description: "Domain name binding of certificates.",
																		},
																		"cert_ca_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Root certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
																		},
																		"s_s_l_mode": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Certificate certification mode: unidirectional unidirectional authentication, Mutual two -way certificationNote: This field may return NULL, indicating that the valid value cannot be obtained.",
																		},
																	},
																},
															},
															"no_match_domains": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "List of non -matching fieldsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
														},
													},
												},
												"no_match_domains": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "List of non -matching fieldsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
												},
											},
										},
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of CLB instances under this region.",
						},
					},
				},
			},

			"cdn": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Related CDN resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of domain names in the region.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CDN domain name detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deployment certificate ID.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"https_billing_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name billing status.",
									},
								},
							},
						},
					},
				},
			},

			"waf": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Related WAF resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "area.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "WAF instance detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"keepalive": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether to keep a long connection, 1 is, 0 NoNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of WAF instances under this region.",
						},
					},
				},
			},

			"ddos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Related DDOS resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of DDOS domain names under this region.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "DDOS example detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "domain name.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "agreement type.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"virtual_port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Forwarding port.",
									},
								},
							},
						},
					},
				},
			},

			"live": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Related live resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of LIVE instances in this region.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Live instance detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Binded certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "-1: Unrelated certificate of domain name.1: The domain name HTTPS has been opened.0: The domain name HTTPS has been closed.",
									},
								},
							},
						},
					},
				},
			},

			"vod": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Related VOD resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "VOD example detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate ID.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of VOD instances under this region.",
						},
					},
				},
			},

			"tke": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Related TKE resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "area.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TKE instance detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster ID.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster name.",
									},
									"namespace_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Cluster Naming Space List.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "namespace name.",
												},
												"secret_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Secret list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Secret name.",
															},
															"cert_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Certificate ID.",
															},
															"ingress_list": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Ingress list.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ingress_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Ingress name.",
																		},
																		"tls_domains": {
																			Type: schema.TypeSet,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																			Computed:    true,
																			Description: "TLS domain name list.",
																		},
																		"domains": {
																			Type: schema.TypeSet,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																			Computed:    true,
																			Description: "Ingress domain name list.",
																		},
																	},
																},
															},
															"no_match_domains": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "List of domain names that are not matched with the new certificateNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
														},
													},
												},
											},
										},
									},
									"cluster_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster.",
									},
									"cluster_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total TKE instance under this region.",
						},
					},
				},
			},

			"apigateway": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Related ApIgateway Resources DetailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "area.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Apiguateway example detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Example name.",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Use Agreement.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of Apigateway instances under this region.",
						},
					},
				},
			},

			"tcb": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Related TCB resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "area.",
						},
						"environments": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TCB environment example detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "TCB environmentNote: This field may return NULL, indicating that the valid value cannot be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_d": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "sourceNote: This field may return NULL, indicating that the valid value cannot be obtained.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "nameNote: This field may return NULL, indicating that the valid value cannot be obtained.",
												},
												"status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "stateNote: This field may return NULL, indicating that the valid value cannot be obtained.",
												},
											},
										},
									},
									"access_service": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Access serviceNote: This field may return NULL, indicating that the valid value cannot be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Instance listNote: This field may return NULL, indicating that the valid value cannot be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"domain": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "domain nameNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
															"status": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "stateNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
															"union_status": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Unified domain name stateNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
															"is_preempted": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether it is seized, being seized means that the domain name is bound by other environments, and needs to be unbinded or re -bound.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
															"i_c_p_status": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ICP blacklist is banned, 0-unblocks, 1-banNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
															"old_certificate_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Binding certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
														},
													},
												},
												"total_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "quantityNote: This field may return NULL, indicating that the valid value cannot be obtained.",
												},
											},
										},
									},
									"host_service": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Static custodyNote: This field may return NULL, indicating that the valid value cannot be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Instance listNote: This field may return NULL, indicating that the valid value cannot be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"domain": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "domain nameNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
															"status": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "stateNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
															"d_n_s_status": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Analytical statusNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
															"old_certificate_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Binding certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
															},
														},
													},
												},
												"total_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "quantityNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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

			"teo": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Related Teo resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Edgeone example detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate ID.",
									},
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regional IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Edgeone total number.",
						},
					},
				},
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Related Cloud Resources Inquiry results: 0 indicates that in the query, 1 means the query is successful.2 means the query is abnormal; if the status is 1, check the results of bindResourceResult; if the state is 2, check the reason for ERROR.",
			},

			"cache_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Current result cache time.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeCertificateBindResourceTaskDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_certificate_bind_resource_task_detail.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var taskId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("task_id"); ok {
		taskId = v.(string)
		paramMap["TaskId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_types"); ok {
		resourceTypesSet := v.(*schema.Set).List()
		paramMap["ResourceTypes"] = helper.InterfacesStringsPoint(resourceTypesSet)
	}

	if v, ok := d.GetOk("regions"); ok {
		regionsSet := v.(*schema.Set).List()
		paramMap["Regions"] = helper.InterfacesStringsPoint(regionsSet)
	}

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var describeRes *ssl.DescribeCertificateBindResourceTaskDetailResponseParams

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeCertificateBindResourceTaskDetailByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		describeRes = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, 10)

	if describeRes.CLB != nil {
		for _, clbInstanceList := range describeRes.CLB {
			clbInstanceListMap := map[string]interface{}{}

			if clbInstanceList.Region != nil {
				clbInstanceListMap["region"] = clbInstanceList.Region
			}

			if clbInstanceList.InstanceList != nil {
				instanceListList := []interface{}{}
				for _, instanceList := range clbInstanceList.InstanceList {
					instanceListMap := map[string]interface{}{}

					if instanceList.LoadBalancerId != nil {
						instanceListMap["load_balancer_id"] = instanceList.LoadBalancerId
					}

					if instanceList.LoadBalancerName != nil {
						instanceListMap["load_balancer_name"] = instanceList.LoadBalancerName
					}

					if instanceList.Listeners != nil {
						listenersList := []interface{}{}
						for _, listeners := range instanceList.Listeners {
							listenersMap := map[string]interface{}{}

							if listeners.ListenerId != nil {
								listenersMap["listener_id"] = listeners.ListenerId
							}

							if listeners.ListenerName != nil {
								listenersMap["listener_name"] = listeners.ListenerName
							}

							if listeners.SniSwitch != nil {
								listenersMap["sni_switch"] = listeners.SniSwitch
							}

							if listeners.Protocol != nil {
								listenersMap["protocol"] = listeners.Protocol
							}

							if listeners.Certificate != nil {
								certificateMap := map[string]interface{}{}

								if listeners.Certificate.CertId != nil {
									certificateMap["cert_id"] = listeners.Certificate.CertId
								}

								if listeners.Certificate.DnsNames != nil {
									certificateMap["dns_names"] = listeners.Certificate.DnsNames
								}

								if listeners.Certificate.CertCaId != nil {
									certificateMap["cert_ca_id"] = listeners.Certificate.CertCaId
								}

								if listeners.Certificate.SSLMode != nil {
									certificateMap["s_s_l_mode"] = listeners.Certificate.SSLMode
								}

								listenersMap["certificate"] = []interface{}{certificateMap}
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

									if rules.IsMatch != nil {
										rulesMap["is_match"] = rules.IsMatch
									}

									if rules.Certificate != nil {
										certificateMap := map[string]interface{}{}

										if rules.Certificate.CertId != nil {
											certificateMap["cert_id"] = rules.Certificate.CertId
										}

										if rules.Certificate.DnsNames != nil {
											certificateMap["dns_names"] = rules.Certificate.DnsNames
										}

										if rules.Certificate.CertCaId != nil {
											certificateMap["cert_ca_id"] = rules.Certificate.CertCaId
										}

										if rules.Certificate.SSLMode != nil {
											certificateMap["s_s_l_mode"] = rules.Certificate.SSLMode
										}

										rulesMap["certificate"] = []interface{}{certificateMap}
									}

									if rules.NoMatchDomains != nil {
										rulesMap["no_match_domains"] = rules.NoMatchDomains
									}

									rulesList = append(rulesList, rulesMap)
								}

								listenersMap["rules"] = []interface{}{rulesList}
							}

							if listeners.NoMatchDomains != nil {
								listenersMap["no_match_domains"] = listeners.NoMatchDomains
							}

							listenersList = append(listenersList, listenersMap)
						}

						instanceListMap["listeners"] = []interface{}{listenersList}
					}

					instanceListList = append(instanceListList, instanceListMap)
				}

				clbInstanceListMap["instance_list"] = []interface{}{instanceListList}
			}

			if clbInstanceList.TotalCount != nil {
				clbInstanceListMap["total_count"] = clbInstanceList.TotalCount
			}

			tmpList = append(tmpList, clbInstanceListMap)
		}

		_ = d.Set("clb", tmpList)
	}

	if describeRes.CDN != nil {
		for _, cdnInstanceList := range describeRes.CDN {
			cdnInstanceListMap := map[string]interface{}{}

			if cdnInstanceList.TotalCount != nil {
				cdnInstanceListMap["total_count"] = cdnInstanceList.TotalCount
			}

			if cdnInstanceList.InstanceList != nil {
				instanceListList := []interface{}{}
				for _, instanceList := range cdnInstanceList.InstanceList {
					instanceListMap := map[string]interface{}{}

					if instanceList.Domain != nil {
						instanceListMap["domain"] = instanceList.Domain
					}

					if instanceList.CertId != nil {
						instanceListMap["cert_id"] = instanceList.CertId
					}

					if instanceList.Status != nil {
						instanceListMap["status"] = instanceList.Status
					}

					if instanceList.HttpsBillingSwitch != nil {
						instanceListMap["https_billing_switch"] = instanceList.HttpsBillingSwitch
					}

					instanceListList = append(instanceListList, instanceListMap)
				}

				cdnInstanceListMap["instance_list"] = []interface{}{instanceListList}
			}

			tmpList = append(tmpList, cdnInstanceListMap)
		}

		_ = d.Set("cdn", tmpList)
	}

	if describeRes.WAF != nil {
		for _, wafInstanceList := range describeRes.WAF {
			wafInstanceListMap := map[string]interface{}{}

			if wafInstanceList.Region != nil {
				wafInstanceListMap["region"] = wafInstanceList.Region
			}

			if wafInstanceList.InstanceList != nil {
				instanceListList := []interface{}{}
				for _, instanceList := range wafInstanceList.InstanceList {
					instanceListMap := map[string]interface{}{}

					if instanceList.Domain != nil {
						instanceListMap["domain"] = instanceList.Domain
					}

					if instanceList.CertId != nil {
						instanceListMap["cert_id"] = instanceList.CertId
					}

					if instanceList.Keepalive != nil {
						instanceListMap["keepalive"] = instanceList.Keepalive
					}

					instanceListList = append(instanceListList, instanceListMap)
				}

				wafInstanceListMap["instance_list"] = []interface{}{instanceListList}
			}

			if wafInstanceList.TotalCount != nil {
				wafInstanceListMap["total_count"] = wafInstanceList.TotalCount
			}

			tmpList = append(tmpList, wafInstanceListMap)
		}

		_ = d.Set("waf", tmpList)
	}

	if describeRes.DDOS != nil {
		for _, ddosInstanceList := range describeRes.DDOS {
			ddosInstanceListMap := map[string]interface{}{}

			if ddosInstanceList.TotalCount != nil {
				ddosInstanceListMap["total_count"] = ddosInstanceList.TotalCount
			}

			if ddosInstanceList.InstanceList != nil {
				instanceListList := []interface{}{}
				for _, instanceList := range ddosInstanceList.InstanceList {
					instanceListMap := map[string]interface{}{}

					if instanceList.Domain != nil {
						instanceListMap["domain"] = instanceList.Domain
					}

					if instanceList.InstanceId != nil {
						instanceListMap["instance_id"] = instanceList.InstanceId
					}

					if instanceList.Protocol != nil {
						instanceListMap["protocol"] = instanceList.Protocol
					}

					if instanceList.CertId != nil {
						instanceListMap["cert_id"] = instanceList.CertId
					}

					if instanceList.VirtualPort != nil {
						instanceListMap["virtual_port"] = instanceList.VirtualPort
					}

					instanceListList = append(instanceListList, instanceListMap)
				}

				ddosInstanceListMap["instance_list"] = []interface{}{instanceListList}
			}

			tmpList = append(tmpList, ddosInstanceListMap)
		}

		_ = d.Set("ddos", tmpList)
	}

	if describeRes.LIVE != nil {
		for _, liveInstanceList := range describeRes.LIVE {
			liveInstanceListMap := map[string]interface{}{}

			if liveInstanceList.TotalCount != nil {
				liveInstanceListMap["total_count"] = liveInstanceList.TotalCount
			}

			if liveInstanceList.InstanceList != nil {
				instanceListList := []interface{}{}
				for _, instanceList := range liveInstanceList.InstanceList {
					instanceListMap := map[string]interface{}{}

					if instanceList.Domain != nil {
						instanceListMap["domain"] = instanceList.Domain
					}

					if instanceList.CertId != nil {
						instanceListMap["cert_id"] = instanceList.CertId
					}

					if instanceList.Status != nil {
						instanceListMap["status"] = instanceList.Status
					}

					instanceListList = append(instanceListList, instanceListMap)
				}

				liveInstanceListMap["instance_list"] = []interface{}{instanceListList}
			}

			tmpList = append(tmpList, liveInstanceListMap)
		}

		_ = d.Set("live", tmpList)
	}

	if describeRes.VOD != nil {
		for _, vODInstanceList := range describeRes.VOD {
			vODInstanceListMap := map[string]interface{}{}

			if vODInstanceList.InstanceList != nil {
				instanceListList := []interface{}{}
				for _, instanceList := range vODInstanceList.InstanceList {
					instanceListMap := map[string]interface{}{}

					if instanceList.Domain != nil {
						instanceListMap["domain"] = instanceList.Domain
					}

					if instanceList.CertId != nil {
						instanceListMap["cert_id"] = instanceList.CertId
					}

					instanceListList = append(instanceListList, instanceListMap)
				}

				vODInstanceListMap["instance_list"] = []interface{}{instanceListList}
			}

			if vODInstanceList.TotalCount != nil {
				vODInstanceListMap["total_count"] = vODInstanceList.TotalCount
			}

			tmpList = append(tmpList, vODInstanceListMap)
		}

		_ = d.Set("vod", tmpList)
	}

	if describeRes.TKE != nil {
		for _, tkeInstanceList := range describeRes.TKE {
			tkeInstanceListMap := map[string]interface{}{}

			if tkeInstanceList.Region != nil {
				tkeInstanceListMap["region"] = tkeInstanceList.Region
			}

			if tkeInstanceList.InstanceList != nil {
				instanceListList := []interface{}{}
				for _, instanceList := range tkeInstanceList.InstanceList {
					instanceListMap := map[string]interface{}{}

					if instanceList.ClusterId != nil {
						instanceListMap["cluster_id"] = instanceList.ClusterId
					}

					if instanceList.ClusterName != nil {
						instanceListMap["cluster_name"] = instanceList.ClusterName
					}

					if instanceList.NamespaceList != nil {
						namespaceListList := []interface{}{}
						for _, namespaceList := range instanceList.NamespaceList {
							namespaceListMap := map[string]interface{}{}

							if namespaceList.Name != nil {
								namespaceListMap["name"] = namespaceList.Name
							}

							if namespaceList.SecretList != nil {
								secretListList := []interface{}{}
								for _, secretList := range namespaceList.SecretList {
									secretListMap := map[string]interface{}{}

									if secretList.Name != nil {
										secretListMap["name"] = secretList.Name
									}

									if secretList.CertId != nil {
										secretListMap["cert_id"] = secretList.CertId
									}

									if secretList.IngressList != nil {
										ingressListList := []interface{}{}
										for _, ingressList := range secretList.IngressList {
											ingressListMap := map[string]interface{}{}

											if ingressList.IngressName != nil {
												ingressListMap["ingress_name"] = ingressList.IngressName
											}

											if ingressList.TlsDomains != nil {
												ingressListMap["tls_domains"] = ingressList.TlsDomains
											}

											if ingressList.Domains != nil {
												ingressListMap["domains"] = ingressList.Domains
											}

											ingressListList = append(ingressListList, ingressListMap)
										}

										secretListMap["ingress_list"] = []interface{}{ingressListList}
									}

									if secretList.NoMatchDomains != nil {
										secretListMap["no_match_domains"] = secretList.NoMatchDomains
									}

									secretListList = append(secretListList, secretListMap)
								}

								namespaceListMap["secret_list"] = []interface{}{secretListList}
							}

							namespaceListList = append(namespaceListList, namespaceListMap)
						}

						instanceListMap["namespace_list"] = []interface{}{namespaceListList}
					}

					if instanceList.ClusterType != nil {
						instanceListMap["cluster_type"] = instanceList.ClusterType
					}

					if instanceList.ClusterVersion != nil {
						instanceListMap["cluster_version"] = instanceList.ClusterVersion
					}

					instanceListList = append(instanceListList, instanceListMap)
				}

				tkeInstanceListMap["instance_list"] = []interface{}{instanceListList}
			}

			if tkeInstanceList.TotalCount != nil {
				tkeInstanceListMap["total_count"] = tkeInstanceList.TotalCount
			}

			tmpList = append(tmpList, tkeInstanceListMap)
		}

		_ = d.Set("tke", tmpList)
	}

	if describeRes.APIGATEWAY != nil {
		for _, apiGatewayInstanceList := range describeRes.APIGATEWAY {
			apiGatewayInstanceListMap := map[string]interface{}{}

			if apiGatewayInstanceList.Region != nil {
				apiGatewayInstanceListMap["region"] = apiGatewayInstanceList.Region
			}

			if apiGatewayInstanceList.InstanceList != nil {
				instanceListList := []interface{}{}
				for _, instanceList := range apiGatewayInstanceList.InstanceList {
					instanceListMap := map[string]interface{}{}

					if instanceList.ServiceId != nil {
						instanceListMap["service_id"] = instanceList.ServiceId
					}

					if instanceList.ServiceName != nil {
						instanceListMap["service_name"] = instanceList.ServiceName
					}

					if instanceList.Domain != nil {
						instanceListMap["domain"] = instanceList.Domain
					}

					if instanceList.CertId != nil {
						instanceListMap["cert_id"] = instanceList.CertId
					}

					if instanceList.Protocol != nil {
						instanceListMap["protocol"] = instanceList.Protocol
					}

					instanceListList = append(instanceListList, instanceListMap)
				}

				apiGatewayInstanceListMap["instance_list"] = []interface{}{instanceListList}
			}

			if apiGatewayInstanceList.TotalCount != nil {
				apiGatewayInstanceListMap["total_count"] = apiGatewayInstanceList.TotalCount
			}

			tmpList = append(tmpList, apiGatewayInstanceListMap)
		}

		_ = d.Set("apigateway", tmpList)
	}

	if describeRes.TCB != nil {
		for _, tCBInstanceList := range describeRes.TCB {
			tCBInstanceListMap := map[string]interface{}{}

			if tCBInstanceList.Region != nil {
				tCBInstanceListMap["region"] = tCBInstanceList.Region
			}

			if tCBInstanceList.Environments != nil {
				environmentsList := []interface{}{}
				for _, environments := range tCBInstanceList.Environments {
					environmentsMap := map[string]interface{}{}

					if environments.Environment != nil {
						environmentMap := map[string]interface{}{}

						if environments.Environment.ID != nil {
							environmentMap["i_d"] = environments.Environment.ID
						}

						if environments.Environment.Source != nil {
							environmentMap["source"] = environments.Environment.Source
						}

						if environments.Environment.Name != nil {
							environmentMap["name"] = environments.Environment.Name
						}

						if environments.Environment.Status != nil {
							environmentMap["status"] = environments.Environment.Status
						}

						environmentsMap["environment"] = []interface{}{environmentMap}
					}

					if environments.AccessService != nil {
						accessServiceMap := map[string]interface{}{}

						if environments.AccessService.InstanceList != nil {
							instanceListList := []interface{}{}
							for _, instanceList := range environments.AccessService.InstanceList {
								instanceListMap := map[string]interface{}{}

								if instanceList.Domain != nil {
									instanceListMap["domain"] = instanceList.Domain
								}

								if instanceList.Status != nil {
									instanceListMap["status"] = instanceList.Status
								}

								if instanceList.UnionStatus != nil {
									instanceListMap["union_status"] = instanceList.UnionStatus
								}

								if instanceList.IsPreempted != nil {
									instanceListMap["is_preempted"] = instanceList.IsPreempted
								}

								if instanceList.ICPStatus != nil {
									instanceListMap["i_c_p_status"] = instanceList.ICPStatus
								}

								if instanceList.OldCertificateId != nil {
									instanceListMap["old_certificate_id"] = instanceList.OldCertificateId
								}

								instanceListList = append(instanceListList, instanceListMap)
							}

							accessServiceMap["instance_list"] = []interface{}{instanceListList}
						}

						if environments.AccessService.TotalCount != nil {
							accessServiceMap["total_count"] = environments.AccessService.TotalCount
						}

						environmentsMap["access_service"] = []interface{}{accessServiceMap}
					}

					if environments.HostService != nil {
						hostServiceMap := map[string]interface{}{}

						if environments.HostService.InstanceList != nil {
							instanceListList := []interface{}{}
							for _, instanceList := range environments.HostService.InstanceList {
								instanceListMap := map[string]interface{}{}

								if instanceList.Domain != nil {
									instanceListMap["domain"] = instanceList.Domain
								}

								if instanceList.Status != nil {
									instanceListMap["status"] = instanceList.Status
								}

								if instanceList.DNSStatus != nil {
									instanceListMap["d_n_s_status"] = instanceList.DNSStatus
								}

								if instanceList.OldCertificateId != nil {
									instanceListMap["old_certificate_id"] = instanceList.OldCertificateId
								}

								instanceListList = append(instanceListList, instanceListMap)
							}

							hostServiceMap["instance_list"] = []interface{}{instanceListList}
						}

						if environments.HostService.TotalCount != nil {
							hostServiceMap["total_count"] = environments.HostService.TotalCount
						}

						environmentsMap["host_service"] = []interface{}{hostServiceMap}
					}

					environmentsList = append(environmentsList, environmentsMap)
				}

				tCBInstanceListMap["environments"] = []interface{}{environmentsList}
			}

			tmpList = append(tmpList, tCBInstanceListMap)
		}

		_ = d.Set("tcb", tmpList)
	}

	if describeRes.TEO != nil {
		for _, teoInstanceList := range describeRes.TEO {
			teoInstanceListMap := map[string]interface{}{}

			if teoInstanceList.InstanceList != nil {
				instanceListList := []interface{}{}
				for _, instanceList := range teoInstanceList.InstanceList {
					instanceListMap := map[string]interface{}{}

					if instanceList.Host != nil {
						instanceListMap["host"] = instanceList.Host
					}

					if instanceList.CertId != nil {
						instanceListMap["cert_id"] = instanceList.CertId
					}

					if instanceList.ZoneId != nil {
						instanceListMap["zone_id"] = instanceList.ZoneId
					}

					if instanceList.Status != nil {
						instanceListMap["status"] = instanceList.Status
					}

					instanceListList = append(instanceListList, instanceListMap)
				}

				teoInstanceListMap["instance_list"] = []interface{}{instanceListList}
			}

			if teoInstanceList.TotalCount != nil {
				teoInstanceListMap["total_count"] = teoInstanceList.TotalCount
			}

			tmpList = append(tmpList, teoInstanceListMap)
		}

		_ = d.Set("teo", tmpList)
	}

	if describeRes.Status != nil {
		_ = d.Set("status", describeRes.Status)
	}

	if describeRes.CacheTime != nil {
		_ = d.Set("cache_time", describeRes.CacheTime)
	}

	d.SetId(taskId)
	output3, ok := d.GetOk("result_output_file")
	if ok && output3.(string) != "" {
		if e := writeToFile(output3.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
