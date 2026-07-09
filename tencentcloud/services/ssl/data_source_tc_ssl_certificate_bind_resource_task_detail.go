package ssl

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslCertificateBindResourceTaskDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslCertificateBindResourceTaskDetailRead,
		Schema: map[string]*schema.Schema{
			"task_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task ID, query the bind cloud resource result according to the task ID obtained by CreateCertificateBindResourceSyncTask.",
			},

			"resource_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Query the result details of the resource type. If not passed, query all. Supported values: clb, cdn, ddos, live, vod, waf, apigateway, teo, tke, cos, tse, tcb, tdmq, mqtt, gaap, scf.",
			},

			"regions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of regions to query. clb, tke, waf, apigateway, tcb, cos, tse support region query.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Async query result status of associated cloud resources: 0 means querying in progress, 1 means query success, 2 means query exception.",
			},

			"cache_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Current result cache time.",
			},

			"clb": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated CLB resource detail list. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of CLB instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CLB instance detail list.",
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
										Description: "CLB instance name.",
									},
									"forward": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Load balancer type, 0 traditional load balancer; 1 application load balancer.",
									},
									"listeners": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "CLB listener list.",
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
													Description: "Listener name.",
												},
												"sni_switch": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether to enable SNI, 1 for enable, 0 for disable.",
												},
												"protocol": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Listener protocol type, HTTPS|TCP_SSL.",
												},
												"certificate": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Certificate data bound to the listener.",
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
																Description: "Domain names bound to the certificate.",
															},
															"cert_ca_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Root certificate ID.",
															},
															"s_s_l_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Certificate authentication mode: UNIDIRECTIONAL one-way authentication, MUTUAL two-way authentication.",
															},
														},
													},
												},
												"rules": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Listener rule list.",
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
																Description: "Domain name bound to the rule.",
															},
															"is_match": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the rule matches the domain name of the certificate to be bound.",
															},
															"certificate": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Certificate data bound to the rule.",
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
																			Description: "Domain names bound to the certificate.",
																		},
																		"cert_ca_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Root certificate ID.",
																		},
																		"s_s_l_mode": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Certificate authentication mode: UNIDIRECTIONAL one-way authentication, MUTUAL two-way authentication.",
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
																Description: "List of non-matching domains.",
															},
															"url": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Path bound to the rule.",
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
													Description: "List of non-matching domains.",
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

			"cdn": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated CDN resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of CDN domains in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CDN domain detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deployed certificate ID.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain status. rejected: domain review not passed, processing: deploying, online: started, offline: closed.",
									},
									"https_billing_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain billing status, on means enabled, off means disabled.",
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
				Description: "Associated WAF resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of WAF instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "WAF instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate ID.",
									},
									"keepalive": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether to keep long connection, 1 yes, 0 no.",
									},
								},
							},
						},
					},
				},
			},

			"ddos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated DDoS resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of DDoS domains in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "DDoS instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol type.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate ID.",
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
				Description: "Associated Live resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of Live instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Live instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bound certificate ID.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "-1: domain not associated with certificate. 1: domain https enabled. 0: domain https disabled.",
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
				Description: "Associated VOD resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of VOD instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "VOD instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate ID.",
									},
								},
							},
						},
					},
				},
			},

			"tke": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated TKE resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of TKE instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TKE instance detail list.",
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
									"cluster_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster type.",
									},
									"cluster_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster version.",
									},
									"namespace_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Cluster namespace list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Namespace name.",
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
																			Description: "TLS domain list.",
																		},
																		"domains": {
																			Type: schema.TypeSet,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																			Computed:    true,
																			Description: "Ingress domain list.",
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
																Description: "List of domains that do not match the new certificate.",
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

			"apigateway": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated APIGateway resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of APIGateway instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "APIGateway instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service ID.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service name.",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate ID.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol.",
									},
								},
							},
						},
					},
				},
			},

			"tcb": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated TCB resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"environments": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TCB environment instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "TCB environment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique ID.",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Source.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name.",
												},
												"status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Status.",
												},
											},
										},
									},
									"access_service": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Access service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Access service instance list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"domain": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Domain name.",
															},
															"status": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Status.",
															},
															"union_status": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Unified domain status.",
															},
															"is_preempted": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether preempted.",
															},
															"i_c_p_status": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ICP blacklist status, 0-not banned, 1-banned.",
															},
															"old_certificate_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Bound certificate ID.",
															},
														},
													},
												},
												"total_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total count.",
												},
											},
										},
									},
									"host_service": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Host service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Host service instance list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"domain": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Domain name.",
															},
															"status": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Status.",
															},
															"d_n_s_status": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "DNS status.",
															},
															"old_certificate_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Bound certificate ID.",
															},
														},
													},
												},
												"total_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total count.",
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
				Description: "Associated TEO resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of TEO instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TEO instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate ID.",
									},
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zone ID.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain status. deployed: deployed; processing: deploying; applying: applying; failed: apply failed; issued: bind failed.",
									},
									"algorithm": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate encryption algorithm.",
									},
								},
							},
						},
					},
				},
			},

			"tse": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated TSE resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of TSE instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TSE instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway ID.",
									},
									"gateway_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway name.",
									},
									"certificate_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Gateway certificate list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Gateway certificate ID.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Gateway certificate name.",
												},
												"bind_domains": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Bound domains.",
												},
												"cert_source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Certificate source.",
												},
												"cert_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Currently bound SSL certificate ID.",
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

			"cos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated COS resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of COS instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "COS instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bound certificate ID.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ENABLED: domain online status; DISABLED: domain offline status.",
									},
									"bucket": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bucket name.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bucket region.",
									},
								},
							},
						},
					},
				},
			},

			"tdmq": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated TDMQ resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of TDMQ instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TDMQ instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance name.",
									},
									"instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance status.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Server certificate ID.",
									},
									"ca_cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CA certificate ID.",
									},
									"no_match_domains": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "List of non-matching domains.",
									},
								},
							},
						},
					},
				},
			},

			"mqtt": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated MQTT resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of MQTT instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "MQTT instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance name.",
									},
									"instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance status.",
									},
									"no_match_domains": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "List of non-matching domains.",
									},
									"server_cert_id_list": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Server certificate list.",
									},
									"ca_cert_id_list": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "CA certificate list.",
									},
								},
							},
						},
					},
				},
			},

			"gaap": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated GAAP resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of GAAP instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "GAAP instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance name.",
									},
									"listener_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Listener list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"listener_status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Listener status.",
												},
												"listener_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Listener ID.",
												},
												"listener_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Listener name.",
												},
												"no_match_domains": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "List of non-matching domains.",
												},
												"cert_id_list": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Certificate ID list bound to the instance.",
												},
												"protocol": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Listener protocol.",
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

			"scf": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated SCF resource detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of SCF instances in the region.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query error message.",
						},
						"instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "SCF instance detail list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate ID.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol.",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region.",
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

func dataSourceTencentCloudSslCertificateBindResourceTaskDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_certificate_bind_resource_task_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("task_id"); ok {
		paramMap["TaskId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_types"); ok {
		resourceTypesSet := v.(*schema.Set).List()
		tmpSet := make([]*string, 0, len(resourceTypesSet))
		for i := range resourceTypesSet {
			resourceType := resourceTypesSet[i].(string)
			tmpSet = append(tmpSet, helper.String(resourceType))
		}
		paramMap["ResourceTypes"] = tmpSet
	}

	if v, ok := d.GetOk("regions"); ok {
		regionsSet := v.(*schema.Set).List()
		tmpSet := make([]*string, 0, len(regionsSet))
		for i := range regionsSet {
			region := regionsSet[i].(string)
			tmpSet = append(tmpSet, helper.String(region))
		}
		paramMap["Regions"] = tmpSet
	}

	service := SslService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var resp *ssl.DescribeCertificateBindResourceTaskDetailResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeCertificateBindResourceTaskDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil {
			log.Printf("[DATASOURCE] read empty, skip SetId")
			return resource.NonRetryableError(fmt.Errorf("ssl_certificate_bind_resource_task_detail read response is nil"))
		}
		resp = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if resp.Status != nil {
		_ = d.Set("status", resp.Status)
	}

	if resp.CacheTime != nil {
		_ = d.Set("cache_time", resp.CacheTime)
	}

	if resp.CLB != nil {
		clbList := make([]map[string]interface{}, 0, len(resp.CLB))
		for _, clbInstanceList := range resp.CLB {
			clbInstanceListMap := map[string]interface{}{}
			if clbInstanceList.Region != nil {
				clbInstanceListMap["region"] = clbInstanceList.Region
			}
			if clbInstanceList.TotalCount != nil {
				clbInstanceListMap["total_count"] = clbInstanceList.TotalCount
			}
			if clbInstanceList.Error != nil {
				clbInstanceListMap["error"] = clbInstanceList.Error
			}
			if clbInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(clbInstanceList.InstanceList))
				for _, clbInstanceDetail := range clbInstanceList.InstanceList {
					clbInstanceDetailMap := map[string]interface{}{}
					if clbInstanceDetail.LoadBalancerId != nil {
						clbInstanceDetailMap["load_balancer_id"] = clbInstanceDetail.LoadBalancerId
					}
					if clbInstanceDetail.LoadBalancerName != nil {
						clbInstanceDetailMap["load_balancer_name"] = clbInstanceDetail.LoadBalancerName
					}
					if clbInstanceDetail.Forward != nil {
						clbInstanceDetailMap["forward"] = clbInstanceDetail.Forward
					}
					if clbInstanceDetail.Listeners != nil {
						listenersList := make([]map[string]interface{}, 0, len(clbInstanceDetail.Listeners))
						for _, listener := range clbInstanceDetail.Listeners {
							listenerMap := map[string]interface{}{}
							if listener.ListenerId != nil {
								listenerMap["listener_id"] = listener.ListenerId
							}
							if listener.ListenerName != nil {
								listenerMap["listener_name"] = listener.ListenerName
							}
							if listener.SniSwitch != nil {
								listenerMap["sni_switch"] = listener.SniSwitch
							}
							if listener.Protocol != nil {
								listenerMap["protocol"] = listener.Protocol
							}
							if listener.Certificate != nil {
								certificateMap := map[string]interface{}{}
								if listener.Certificate.CertId != nil {
									certificateMap["cert_id"] = listener.Certificate.CertId
								}
								if listener.Certificate.DnsNames != nil {
									certificateMap["dns_names"] = listener.Certificate.DnsNames
								}
								if listener.Certificate.CertCaId != nil {
									certificateMap["cert_ca_id"] = listener.Certificate.CertCaId
								}
								if listener.Certificate.SSLMode != nil {
									certificateMap["s_s_l_mode"] = listener.Certificate.SSLMode
								}
								listenerMap["certificate"] = []interface{}{certificateMap}
							}
							if listener.Rules != nil {
								rulesList := make([]map[string]interface{}, 0, len(listener.Rules))
								for _, rule := range listener.Rules {
									ruleMap := map[string]interface{}{}
									if rule.LocationId != nil {
										ruleMap["location_id"] = rule.LocationId
									}
									if rule.Domain != nil {
										ruleMap["domain"] = rule.Domain
									}
									if rule.IsMatch != nil {
										ruleMap["is_match"] = rule.IsMatch
									}
									if rule.Certificate != nil {
										certificateMap := map[string]interface{}{}
										if rule.Certificate.CertId != nil {
											certificateMap["cert_id"] = rule.Certificate.CertId
										}
										if rule.Certificate.DnsNames != nil {
											certificateMap["dns_names"] = rule.Certificate.DnsNames
										}
										if rule.Certificate.CertCaId != nil {
											certificateMap["cert_ca_id"] = rule.Certificate.CertCaId
										}
										if rule.Certificate.SSLMode != nil {
											certificateMap["s_s_l_mode"] = rule.Certificate.SSLMode
										}
										ruleMap["certificate"] = []interface{}{certificateMap}
									}
									if rule.NoMatchDomains != nil {
										ruleMap["no_match_domains"] = rule.NoMatchDomains
									}
									if rule.Url != nil {
										ruleMap["url"] = rule.Url
									}
									rulesList = append(rulesList, ruleMap)
								}
								listenerMap["rules"] = rulesList
							}
							if listener.NoMatchDomains != nil {
								listenerMap["no_match_domains"] = listener.NoMatchDomains
							}
							listenersList = append(listenersList, listenerMap)
						}
						clbInstanceDetailMap["listeners"] = listenersList
					}
					instanceList = append(instanceList, clbInstanceDetailMap)
				}
				clbInstanceListMap["instance_list"] = instanceList
			}
			clbList = append(clbList, clbInstanceListMap)
		}
		_ = d.Set("clb", clbList)
	}

	if resp.CDN != nil {
		cdnList := make([]map[string]interface{}, 0, len(resp.CDN))
		for _, cdnInstanceList := range resp.CDN {
			cdnInstanceListMap := map[string]interface{}{}
			if cdnInstanceList.TotalCount != nil {
				cdnInstanceListMap["total_count"] = cdnInstanceList.TotalCount
			}
			if cdnInstanceList.Error != nil {
				cdnInstanceListMap["error"] = cdnInstanceList.Error
			}
			if cdnInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(cdnInstanceList.InstanceList))
				for _, cdnInstanceDetail := range cdnInstanceList.InstanceList {
					cdnInstanceDetailMap := map[string]interface{}{}
					if cdnInstanceDetail.Domain != nil {
						cdnInstanceDetailMap["domain"] = cdnInstanceDetail.Domain
					}
					if cdnInstanceDetail.CertId != nil {
						cdnInstanceDetailMap["cert_id"] = cdnInstanceDetail.CertId
					}
					if cdnInstanceDetail.Status != nil {
						cdnInstanceDetailMap["status"] = cdnInstanceDetail.Status
					}
					if cdnInstanceDetail.HttpsBillingSwitch != nil {
						cdnInstanceDetailMap["https_billing_switch"] = cdnInstanceDetail.HttpsBillingSwitch
					}
					instanceList = append(instanceList, cdnInstanceDetailMap)
				}
				cdnInstanceListMap["instance_list"] = instanceList
			}
			cdnList = append(cdnList, cdnInstanceListMap)
		}
		_ = d.Set("cdn", cdnList)
	}

	if resp.WAF != nil {
		wafList := make([]map[string]interface{}, 0, len(resp.WAF))
		for _, wafInstanceList := range resp.WAF {
			wafInstanceListMap := map[string]interface{}{}
			if wafInstanceList.Region != nil {
				wafInstanceListMap["region"] = wafInstanceList.Region
			}
			if wafInstanceList.TotalCount != nil {
				wafInstanceListMap["total_count"] = wafInstanceList.TotalCount
			}
			if wafInstanceList.Error != nil {
				wafInstanceListMap["error"] = wafInstanceList.Error
			}
			if wafInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(wafInstanceList.InstanceList))
				for _, wafInstanceDetail := range wafInstanceList.InstanceList {
					wafInstanceDetailMap := map[string]interface{}{}
					if wafInstanceDetail.Domain != nil {
						wafInstanceDetailMap["domain"] = wafInstanceDetail.Domain
					}
					if wafInstanceDetail.CertId != nil {
						wafInstanceDetailMap["cert_id"] = wafInstanceDetail.CertId
					}
					if wafInstanceDetail.Keepalive != nil {
						wafInstanceDetailMap["keepalive"] = wafInstanceDetail.Keepalive
					}
					instanceList = append(instanceList, wafInstanceDetailMap)
				}
				wafInstanceListMap["instance_list"] = instanceList
			}
			wafList = append(wafList, wafInstanceListMap)
		}
		_ = d.Set("waf", wafList)
	}

	if resp.DDOS != nil {
		ddosList := make([]map[string]interface{}, 0, len(resp.DDOS))
		for _, ddosInstanceList := range resp.DDOS {
			ddosInstanceListMap := map[string]interface{}{}
			if ddosInstanceList.TotalCount != nil {
				ddosInstanceListMap["total_count"] = ddosInstanceList.TotalCount
			}
			if ddosInstanceList.Error != nil {
				ddosInstanceListMap["error"] = ddosInstanceList.Error
			}
			if ddosInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(ddosInstanceList.InstanceList))
				for _, ddosInstanceDetail := range ddosInstanceList.InstanceList {
					ddosInstanceDetailMap := map[string]interface{}{}
					if ddosInstanceDetail.Domain != nil {
						ddosInstanceDetailMap["domain"] = ddosInstanceDetail.Domain
					}
					if ddosInstanceDetail.InstanceId != nil {
						ddosInstanceDetailMap["instance_id"] = ddosInstanceDetail.InstanceId
					}
					if ddosInstanceDetail.Protocol != nil {
						ddosInstanceDetailMap["protocol"] = ddosInstanceDetail.Protocol
					}
					if ddosInstanceDetail.CertId != nil {
						ddosInstanceDetailMap["cert_id"] = ddosInstanceDetail.CertId
					}
					if ddosInstanceDetail.VirtualPort != nil {
						ddosInstanceDetailMap["virtual_port"] = ddosInstanceDetail.VirtualPort
					}
					instanceList = append(instanceList, ddosInstanceDetailMap)
				}
				ddosInstanceListMap["instance_list"] = instanceList
			}
			ddosList = append(ddosList, ddosInstanceListMap)
		}
		_ = d.Set("ddos", ddosList)
	}

	if resp.LIVE != nil {
		liveList := make([]map[string]interface{}, 0, len(resp.LIVE))
		for _, liveInstanceList := range resp.LIVE {
			liveInstanceListMap := map[string]interface{}{}
			if liveInstanceList.TotalCount != nil {
				liveInstanceListMap["total_count"] = liveInstanceList.TotalCount
			}
			if liveInstanceList.Error != nil {
				liveInstanceListMap["error"] = liveInstanceList.Error
			}
			if liveInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(liveInstanceList.InstanceList))
				for _, liveInstanceDetail := range liveInstanceList.InstanceList {
					liveInstanceDetailMap := map[string]interface{}{}
					if liveInstanceDetail.Domain != nil {
						liveInstanceDetailMap["domain"] = liveInstanceDetail.Domain
					}
					if liveInstanceDetail.CertId != nil {
						liveInstanceDetailMap["cert_id"] = liveInstanceDetail.CertId
					}
					if liveInstanceDetail.Status != nil {
						liveInstanceDetailMap["status"] = liveInstanceDetail.Status
					}
					instanceList = append(instanceList, liveInstanceDetailMap)
				}
				liveInstanceListMap["instance_list"] = instanceList
			}
			liveList = append(liveList, liveInstanceListMap)
		}
		_ = d.Set("live", liveList)
	}

	if resp.VOD != nil {
		vodList := make([]map[string]interface{}, 0, len(resp.VOD))
		for _, vodInstanceList := range resp.VOD {
			vodInstanceListMap := map[string]interface{}{}
			if vodInstanceList.TotalCount != nil {
				vodInstanceListMap["total_count"] = vodInstanceList.TotalCount
			}
			if vodInstanceList.Error != nil {
				vodInstanceListMap["error"] = vodInstanceList.Error
			}
			if vodInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(vodInstanceList.InstanceList))
				for _, vodInstanceDetail := range vodInstanceList.InstanceList {
					vodInstanceDetailMap := map[string]interface{}{}
					if vodInstanceDetail.Domain != nil {
						vodInstanceDetailMap["domain"] = vodInstanceDetail.Domain
					}
					if vodInstanceDetail.CertId != nil {
						vodInstanceDetailMap["cert_id"] = vodInstanceDetail.CertId
					}
					instanceList = append(instanceList, vodInstanceDetailMap)
				}
				vodInstanceListMap["instance_list"] = instanceList
			}
			vodList = append(vodList, vodInstanceListMap)
		}
		_ = d.Set("vod", vodList)
	}

	if resp.TKE != nil {
		tkeList := make([]map[string]interface{}, 0, len(resp.TKE))
		for _, tkeInstanceList := range resp.TKE {
			tkeInstanceListMap := map[string]interface{}{}
			if tkeInstanceList.Region != nil {
				tkeInstanceListMap["region"] = tkeInstanceList.Region
			}
			if tkeInstanceList.TotalCount != nil {
				tkeInstanceListMap["total_count"] = tkeInstanceList.TotalCount
			}
			if tkeInstanceList.Error != nil {
				tkeInstanceListMap["error"] = tkeInstanceList.Error
			}
			if tkeInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(tkeInstanceList.InstanceList))
				for _, tkeInstanceDetail := range tkeInstanceList.InstanceList {
					tkeInstanceDetailMap := map[string]interface{}{}
					if tkeInstanceDetail.ClusterId != nil {
						tkeInstanceDetailMap["cluster_id"] = tkeInstanceDetail.ClusterId
					}
					if tkeInstanceDetail.ClusterName != nil {
						tkeInstanceDetailMap["cluster_name"] = tkeInstanceDetail.ClusterName
					}
					if tkeInstanceDetail.ClusterType != nil {
						tkeInstanceDetailMap["cluster_type"] = tkeInstanceDetail.ClusterType
					}
					if tkeInstanceDetail.ClusterVersion != nil {
						tkeInstanceDetailMap["cluster_version"] = tkeInstanceDetail.ClusterVersion
					}
					if tkeInstanceDetail.NamespaceList != nil {
						namespaceList := make([]map[string]interface{}, 0, len(tkeInstanceDetail.NamespaceList))
						for _, namespaceDetail := range tkeInstanceDetail.NamespaceList {
							namespaceMap := map[string]interface{}{}
							if namespaceDetail.Name != nil {
								namespaceMap["name"] = namespaceDetail.Name
							}
							if namespaceDetail.SecretList != nil {
								secretList := make([]map[string]interface{}, 0, len(namespaceDetail.SecretList))
								for _, secretDetail := range namespaceDetail.SecretList {
									secretMap := map[string]interface{}{}
									if secretDetail.Name != nil {
										secretMap["name"] = secretDetail.Name
									}
									if secretDetail.CertId != nil {
										secretMap["cert_id"] = secretDetail.CertId
									}
									if secretDetail.IngressList != nil {
										ingressList := make([]map[string]interface{}, 0, len(secretDetail.IngressList))
										for _, ingressDetail := range secretDetail.IngressList {
											ingressMap := map[string]interface{}{}
											if ingressDetail.IngressName != nil {
												ingressMap["ingress_name"] = ingressDetail.IngressName
											}
											if ingressDetail.TlsDomains != nil {
												ingressMap["tls_domains"] = ingressDetail.TlsDomains
											}
											if ingressDetail.Domains != nil {
												ingressMap["domains"] = ingressDetail.Domains
											}
											ingressList = append(ingressList, ingressMap)
										}
										secretMap["ingress_list"] = ingressList
									}
									if secretDetail.NoMatchDomains != nil {
										secretMap["no_match_domains"] = secretDetail.NoMatchDomains
									}
									secretList = append(secretList, secretMap)
								}
								namespaceMap["secret_list"] = secretList
							}
							namespaceList = append(namespaceList, namespaceMap)
						}
						tkeInstanceDetailMap["namespace_list"] = namespaceList
					}
					instanceList = append(instanceList, tkeInstanceDetailMap)
				}
				tkeInstanceListMap["instance_list"] = instanceList
			}
			tkeList = append(tkeList, tkeInstanceListMap)
		}
		_ = d.Set("tke", tkeList)
	}

	if resp.APIGATEWAY != nil {
		apigatewayList := make([]map[string]interface{}, 0, len(resp.APIGATEWAY))
		for _, apiGatewayInstanceList := range resp.APIGATEWAY {
			apiGatewayInstanceListMap := map[string]interface{}{}
			if apiGatewayInstanceList.Region != nil {
				apiGatewayInstanceListMap["region"] = apiGatewayInstanceList.Region
			}
			if apiGatewayInstanceList.TotalCount != nil {
				apiGatewayInstanceListMap["total_count"] = apiGatewayInstanceList.TotalCount
			}
			if apiGatewayInstanceList.Error != nil {
				apiGatewayInstanceListMap["error"] = apiGatewayInstanceList.Error
			}
			if apiGatewayInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(apiGatewayInstanceList.InstanceList))
				for _, apiGatewayInstanceDetail := range apiGatewayInstanceList.InstanceList {
					apiGatewayInstanceDetailMap := map[string]interface{}{}
					if apiGatewayInstanceDetail.ServiceId != nil {
						apiGatewayInstanceDetailMap["service_id"] = apiGatewayInstanceDetail.ServiceId
					}
					if apiGatewayInstanceDetail.ServiceName != nil {
						apiGatewayInstanceDetailMap["service_name"] = apiGatewayInstanceDetail.ServiceName
					}
					if apiGatewayInstanceDetail.Domain != nil {
						apiGatewayInstanceDetailMap["domain"] = apiGatewayInstanceDetail.Domain
					}
					if apiGatewayInstanceDetail.CertId != nil {
						apiGatewayInstanceDetailMap["cert_id"] = apiGatewayInstanceDetail.CertId
					}
					if apiGatewayInstanceDetail.Protocol != nil {
						apiGatewayInstanceDetailMap["protocol"] = apiGatewayInstanceDetail.Protocol
					}
					instanceList = append(instanceList, apiGatewayInstanceDetailMap)
				}
				apiGatewayInstanceListMap["instance_list"] = instanceList
			}
			apigatewayList = append(apigatewayList, apiGatewayInstanceListMap)
		}
		_ = d.Set("apigateway", apigatewayList)
	}

	if resp.TCB != nil {
		tcbList := make([]map[string]interface{}, 0, len(resp.TCB))
		for _, tcbInstanceList := range resp.TCB {
			tcbInstanceListMap := map[string]interface{}{}
			if tcbInstanceList.Region != nil {
				tcbInstanceListMap["region"] = tcbInstanceList.Region
			}
			if tcbInstanceList.Error != nil {
				tcbInstanceListMap["error"] = tcbInstanceList.Error
			}
			if tcbInstanceList.Environments != nil {
				environmentsList := make([]map[string]interface{}, 0, len(tcbInstanceList.Environments))
				for _, environments := range tcbInstanceList.Environments {
					environmentsMap := map[string]interface{}{}
					if environments.Environment != nil {
						environmentMap := map[string]interface{}{}
						if environments.Environment.ID != nil {
							environmentMap["id"] = environments.Environment.ID
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
							accessInstanceList := make([]map[string]interface{}, 0, len(environments.AccessService.InstanceList))
							for _, accessInstance := range environments.AccessService.InstanceList {
								accessInstanceMap := map[string]interface{}{}
								if accessInstance.Domain != nil {
									accessInstanceMap["domain"] = accessInstance.Domain
								}
								if accessInstance.Status != nil {
									accessInstanceMap["status"] = accessInstance.Status
								}
								if accessInstance.UnionStatus != nil {
									accessInstanceMap["union_status"] = accessInstance.UnionStatus
								}
								if accessInstance.IsPreempted != nil {
									accessInstanceMap["is_preempted"] = accessInstance.IsPreempted
								}
								if accessInstance.ICPStatus != nil {
									accessInstanceMap["i_c_p_status"] = accessInstance.ICPStatus
								}
								if accessInstance.OldCertificateId != nil {
									accessInstanceMap["old_certificate_id"] = accessInstance.OldCertificateId
								}
								accessInstanceList = append(accessInstanceList, accessInstanceMap)
							}
							accessServiceMap["instance_list"] = accessInstanceList
						}
						if environments.AccessService.TotalCount != nil {
							accessServiceMap["total_count"] = environments.AccessService.TotalCount
						}
						environmentsMap["access_service"] = []interface{}{accessServiceMap}
					}
					if environments.HostService != nil {
						hostServiceMap := map[string]interface{}{}
						if environments.HostService.InstanceList != nil {
							hostInstanceList := make([]map[string]interface{}, 0, len(environments.HostService.InstanceList))
							for _, hostInstance := range environments.HostService.InstanceList {
								hostInstanceMap := map[string]interface{}{}
								if hostInstance.Domain != nil {
									hostInstanceMap["domain"] = hostInstance.Domain
								}
								if hostInstance.Status != nil {
									hostInstanceMap["status"] = hostInstance.Status
								}
								if hostInstance.DNSStatus != nil {
									hostInstanceMap["d_n_s_status"] = hostInstance.DNSStatus
								}
								if hostInstance.OldCertificateId != nil {
									hostInstanceMap["old_certificate_id"] = hostInstance.OldCertificateId
								}
								hostInstanceList = append(hostInstanceList, hostInstanceMap)
							}
							hostServiceMap["instance_list"] = hostInstanceList
						}
						if environments.HostService.TotalCount != nil {
							hostServiceMap["total_count"] = environments.HostService.TotalCount
						}
						environmentsMap["host_service"] = []interface{}{hostServiceMap}
					}
					environmentsList = append(environmentsList, environmentsMap)
				}
				tcbInstanceListMap["environments"] = environmentsList
			}
			tcbList = append(tcbList, tcbInstanceListMap)
		}
		_ = d.Set("tcb", tcbList)
	}

	if resp.TEO != nil {
		teoList := make([]map[string]interface{}, 0, len(resp.TEO))
		for _, teoInstanceList := range resp.TEO {
			teoInstanceListMap := map[string]interface{}{}
			if teoInstanceList.TotalCount != nil {
				teoInstanceListMap["total_count"] = teoInstanceList.TotalCount
			}
			if teoInstanceList.Error != nil {
				teoInstanceListMap["error"] = teoInstanceList.Error
			}
			if teoInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(teoInstanceList.InstanceList))
				for _, teoInstanceDetail := range teoInstanceList.InstanceList {
					teoInstanceDetailMap := map[string]interface{}{}
					if teoInstanceDetail.Host != nil {
						teoInstanceDetailMap["host"] = teoInstanceDetail.Host
					}
					if teoInstanceDetail.CertId != nil {
						teoInstanceDetailMap["cert_id"] = teoInstanceDetail.CertId
					}
					if teoInstanceDetail.ZoneId != nil {
						teoInstanceDetailMap["zone_id"] = teoInstanceDetail.ZoneId
					}
					if teoInstanceDetail.Status != nil {
						teoInstanceDetailMap["status"] = teoInstanceDetail.Status
					}
					if teoInstanceDetail.Algorithm != nil {
						teoInstanceDetailMap["algorithm"] = teoInstanceDetail.Algorithm
					}
					instanceList = append(instanceList, teoInstanceDetailMap)
				}
				teoInstanceListMap["instance_list"] = instanceList
			}
			teoList = append(teoList, teoInstanceListMap)
		}
		_ = d.Set("teo", teoList)
	}

	if resp.TSE != nil {
		tseList := make([]map[string]interface{}, 0, len(resp.TSE))
		for _, tseInstanceList := range resp.TSE {
			tseInstanceListMap := map[string]interface{}{}
			if tseInstanceList.Region != nil {
				tseInstanceListMap["region"] = tseInstanceList.Region
			}
			if tseInstanceList.TotalCount != nil {
				tseInstanceListMap["total_count"] = tseInstanceList.TotalCount
			}
			if tseInstanceList.Error != nil {
				tseInstanceListMap["error"] = tseInstanceList.Error
			}
			if tseInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(tseInstanceList.InstanceList))
				for _, tseInstanceDetail := range tseInstanceList.InstanceList {
					tseInstanceDetailMap := map[string]interface{}{}
					if tseInstanceDetail.GatewayId != nil {
						tseInstanceDetailMap["gateway_id"] = tseInstanceDetail.GatewayId
					}
					if tseInstanceDetail.GatewayName != nil {
						tseInstanceDetailMap["gateway_name"] = tseInstanceDetail.GatewayName
					}
					if tseInstanceDetail.CertificateList != nil {
						certificateList := make([]map[string]interface{}, 0, len(tseInstanceDetail.CertificateList))
						for _, gatewayCertificate := range tseInstanceDetail.CertificateList {
							gatewayCertificateMap := map[string]interface{}{}
							if gatewayCertificate.Id != nil {
								gatewayCertificateMap["id"] = gatewayCertificate.Id
							}
							if gatewayCertificate.Name != nil {
								gatewayCertificateMap["name"] = gatewayCertificate.Name
							}
							if gatewayCertificate.BindDomains != nil {
								gatewayCertificateMap["bind_domains"] = gatewayCertificate.BindDomains
							}
							if gatewayCertificate.CertSource != nil {
								gatewayCertificateMap["cert_source"] = gatewayCertificate.CertSource
							}
							if gatewayCertificate.CertId != nil {
								gatewayCertificateMap["cert_id"] = gatewayCertificate.CertId
							}
							certificateList = append(certificateList, gatewayCertificateMap)
						}
						tseInstanceDetailMap["certificate_list"] = certificateList
					}
					instanceList = append(instanceList, tseInstanceDetailMap)
				}
				tseInstanceListMap["instance_list"] = instanceList
			}
			tseList = append(tseList, tseInstanceListMap)
		}
		_ = d.Set("tse", tseList)
	}

	if resp.COS != nil {
		cosList := make([]map[string]interface{}, 0, len(resp.COS))
		for _, cosInstanceList := range resp.COS {
			cosInstanceListMap := map[string]interface{}{}
			if cosInstanceList.Region != nil {
				cosInstanceListMap["region"] = cosInstanceList.Region
			}
			if cosInstanceList.TotalCount != nil {
				cosInstanceListMap["total_count"] = cosInstanceList.TotalCount
			}
			if cosInstanceList.Error != nil {
				cosInstanceListMap["error"] = cosInstanceList.Error
			}
			if cosInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(cosInstanceList.InstanceList))
				for _, cosInstanceDetail := range cosInstanceList.InstanceList {
					cosInstanceDetailMap := map[string]interface{}{}
					if cosInstanceDetail.Domain != nil {
						cosInstanceDetailMap["domain"] = cosInstanceDetail.Domain
					}
					if cosInstanceDetail.CertId != nil {
						cosInstanceDetailMap["cert_id"] = cosInstanceDetail.CertId
					}
					if cosInstanceDetail.Status != nil {
						cosInstanceDetailMap["status"] = cosInstanceDetail.Status
					}
					if cosInstanceDetail.Bucket != nil {
						cosInstanceDetailMap["bucket"] = cosInstanceDetail.Bucket
					}
					if cosInstanceDetail.Region != nil {
						cosInstanceDetailMap["region"] = cosInstanceDetail.Region
					}
					instanceList = append(instanceList, cosInstanceDetailMap)
				}
				cosInstanceListMap["instance_list"] = instanceList
			}
			cosList = append(cosList, cosInstanceListMap)
		}
		_ = d.Set("cos", cosList)
	}

	if resp.TDMQ != nil {
		tdmqList := make([]map[string]interface{}, 0, len(resp.TDMQ))
		for _, tdmqInstanceList := range resp.TDMQ {
			tdmqInstanceListMap := map[string]interface{}{}
			if tdmqInstanceList.Region != nil {
				tdmqInstanceListMap["region"] = tdmqInstanceList.Region
			}
			if tdmqInstanceList.TotalCount != nil {
				tdmqInstanceListMap["total_count"] = tdmqInstanceList.TotalCount
			}
			if tdmqInstanceList.Error != nil {
				tdmqInstanceListMap["error"] = tdmqInstanceList.Error
			}
			if tdmqInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(tdmqInstanceList.InstanceList))
				for _, tdmqInstanceDetail := range tdmqInstanceList.InstanceList {
					tdmqInstanceDetailMap := map[string]interface{}{}
					if tdmqInstanceDetail.InstanceId != nil {
						tdmqInstanceDetailMap["instance_id"] = tdmqInstanceDetail.InstanceId
					}
					if tdmqInstanceDetail.InstanceName != nil {
						tdmqInstanceDetailMap["instance_name"] = tdmqInstanceDetail.InstanceName
					}
					if tdmqInstanceDetail.InstanceStatus != nil {
						tdmqInstanceDetailMap["instance_status"] = tdmqInstanceDetail.InstanceStatus
					}
					if tdmqInstanceDetail.CertId != nil {
						tdmqInstanceDetailMap["cert_id"] = tdmqInstanceDetail.CertId
					}
					if tdmqInstanceDetail.CaCertId != nil {
						tdmqInstanceDetailMap["ca_cert_id"] = tdmqInstanceDetail.CaCertId
					}
					if tdmqInstanceDetail.NoMatchDomains != nil {
						tdmqInstanceDetailMap["no_match_domains"] = tdmqInstanceDetail.NoMatchDomains
					}
					instanceList = append(instanceList, tdmqInstanceDetailMap)
				}
				tdmqInstanceListMap["instance_list"] = instanceList
			}
			tdmqList = append(tdmqList, tdmqInstanceListMap)
		}
		_ = d.Set("tdmq", tdmqList)
	}

	if resp.MQTT != nil {
		mqttList := make([]map[string]interface{}, 0, len(resp.MQTT))
		for _, mqttInstanceList := range resp.MQTT {
			mqttInstanceListMap := map[string]interface{}{}
			if mqttInstanceList.Region != nil {
				mqttInstanceListMap["region"] = mqttInstanceList.Region
			}
			if mqttInstanceList.TotalCount != nil {
				mqttInstanceListMap["total_count"] = mqttInstanceList.TotalCount
			}
			if mqttInstanceList.Error != nil {
				mqttInstanceListMap["error"] = mqttInstanceList.Error
			}
			if mqttInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(mqttInstanceList.InstanceList))
				for _, mqttInstanceDetail := range mqttInstanceList.InstanceList {
					mqttInstanceDetailMap := map[string]interface{}{}
					if mqttInstanceDetail.InstanceId != nil {
						mqttInstanceDetailMap["instance_id"] = mqttInstanceDetail.InstanceId
					}
					if mqttInstanceDetail.InstanceName != nil {
						mqttInstanceDetailMap["instance_name"] = mqttInstanceDetail.InstanceName
					}
					if mqttInstanceDetail.InstanceStatus != nil {
						mqttInstanceDetailMap["instance_status"] = mqttInstanceDetail.InstanceStatus
					}
					if mqttInstanceDetail.NoMatchDomains != nil {
						mqttInstanceDetailMap["no_match_domains"] = mqttInstanceDetail.NoMatchDomains
					}
					if mqttInstanceDetail.ServerCertIdList != nil {
						mqttInstanceDetailMap["server_cert_id_list"] = mqttInstanceDetail.ServerCertIdList
					}
					if mqttInstanceDetail.CaCertIdList != nil {
						mqttInstanceDetailMap["ca_cert_id_list"] = mqttInstanceDetail.CaCertIdList
					}
					instanceList = append(instanceList, mqttInstanceDetailMap)
				}
				mqttInstanceListMap["instance_list"] = instanceList
			}
			mqttList = append(mqttList, mqttInstanceListMap)
		}
		_ = d.Set("mqtt", mqttList)
	}

	if resp.GAAP != nil {
		gaapList := make([]map[string]interface{}, 0, len(resp.GAAP))
		for _, gaapInstanceList := range resp.GAAP {
			gaapInstanceListMap := map[string]interface{}{}
			if gaapInstanceList.TotalCount != nil {
				gaapInstanceListMap["total_count"] = gaapInstanceList.TotalCount
			}
			if gaapInstanceList.Error != nil {
				gaapInstanceListMap["error"] = gaapInstanceList.Error
			}
			if gaapInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(gaapInstanceList.InstanceList))
				for _, gaapInstanceDetail := range gaapInstanceList.InstanceList {
					gaapInstanceDetailMap := map[string]interface{}{}
					if gaapInstanceDetail.InstanceId != nil {
						gaapInstanceDetailMap["instance_id"] = gaapInstanceDetail.InstanceId
					}
					if gaapInstanceDetail.InstanceName != nil {
						gaapInstanceDetailMap["instance_name"] = gaapInstanceDetail.InstanceName
					}
					if gaapInstanceDetail.ListenerList != nil {
						listenerList := make([]map[string]interface{}, 0, len(gaapInstanceDetail.ListenerList))
						for _, gaapListenerDetail := range gaapInstanceDetail.ListenerList {
							gaapListenerDetailMap := map[string]interface{}{}
							if gaapListenerDetail.ListenerStatus != nil {
								gaapListenerDetailMap["listener_status"] = gaapListenerDetail.ListenerStatus
							}
							if gaapListenerDetail.ListenerId != nil {
								gaapListenerDetailMap["listener_id"] = gaapListenerDetail.ListenerId
							}
							if gaapListenerDetail.ListenerName != nil {
								gaapListenerDetailMap["listener_name"] = gaapListenerDetail.ListenerName
							}
							if gaapListenerDetail.NoMatchDomains != nil {
								gaapListenerDetailMap["no_match_domains"] = gaapListenerDetail.NoMatchDomains
							}
							if gaapListenerDetail.CertIdList != nil {
								gaapListenerDetailMap["cert_id_list"] = gaapListenerDetail.CertIdList
							}
							if gaapListenerDetail.Protocol != nil {
								gaapListenerDetailMap["protocol"] = gaapListenerDetail.Protocol
							}
							listenerList = append(listenerList, gaapListenerDetailMap)
						}
						gaapInstanceDetailMap["listener_list"] = listenerList
					}
					instanceList = append(instanceList, gaapInstanceDetailMap)
				}
				gaapInstanceListMap["instance_list"] = instanceList
			}
			gaapList = append(gaapList, gaapInstanceListMap)
		}
		_ = d.Set("gaap", gaapList)
	}

	if resp.SCF != nil {
		scfList := make([]map[string]interface{}, 0, len(resp.SCF))
		for _, scfInstanceList := range resp.SCF {
			scfInstanceListMap := map[string]interface{}{}
			if scfInstanceList.Region != nil {
				scfInstanceListMap["region"] = scfInstanceList.Region
			}
			if scfInstanceList.TotalCount != nil {
				scfInstanceListMap["total_count"] = scfInstanceList.TotalCount
			}
			if scfInstanceList.Error != nil {
				scfInstanceListMap["error"] = scfInstanceList.Error
			}
			if scfInstanceList.InstanceList != nil {
				instanceList := make([]map[string]interface{}, 0, len(scfInstanceList.InstanceList))
				for _, scfInstanceDetail := range scfInstanceList.InstanceList {
					scfInstanceDetailMap := map[string]interface{}{}
					if scfInstanceDetail.CertificateId != nil {
						scfInstanceDetailMap["certificate_id"] = scfInstanceDetail.CertificateId
					}
					if scfInstanceDetail.Protocol != nil {
						scfInstanceDetailMap["protocol"] = scfInstanceDetail.Protocol
					}
					if scfInstanceDetail.Domain != nil {
						scfInstanceDetailMap["domain"] = scfInstanceDetail.Domain
					}
					if scfInstanceDetail.Region != nil {
						scfInstanceDetailMap["region"] = scfInstanceDetail.Region
					}
					instanceList = append(instanceList, scfInstanceDetailMap)
				}
				scfInstanceListMap["instance_list"] = instanceList
			}
			scfList = append(scfList, scfInstanceListMap)
		}
		_ = d.Set("scf", scfList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
