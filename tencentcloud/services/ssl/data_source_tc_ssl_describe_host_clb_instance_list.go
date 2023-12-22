package ssl

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslDescribeHostClbInstanceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostClbInstanceListRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID to be deployed.",
			},

			"is_cache": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to query the cache, 1: Yes; 0: No, the default is the query cache, the cache is half an hour.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of filtering parameters; Filterkey: domainmatch.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter parameter key.",
						},
						"filter_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter parameter value.",
						},
					},
				},
			},

			"async_cache": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to cache asynchronous.",
			},

			"old_certificate_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Original certificate ID.",
			},

			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "CLB instance listener listNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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

			"async_total_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of asynchronous refreshNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"async_offset": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Asynchronous refresh current execution numberNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"async_cache_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Current cache read timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeHostClbInstanceListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_describe_host_clb_instance_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("certificate_id"); ok {
		paramMap["CertificateId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("is_cache"); v != nil {
		paramMap["IsCache"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*ssl.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := ssl.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["filter_key"]; ok {
				filter.FilterKey = helper.String(v.(string))
			}
			if v, ok := filterMap["filter_value"]; ok {
				filter.FilterValue = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, _ := d.GetOk("async_cache"); v != nil {
		paramMap["AsyncCache"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("old_certificate_id"); ok {
		paramMap["OldCertificateId"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceList *ssl.DescribeHostClbInstanceListResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeHostClbInstanceListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList.InstanceList))
	tmpList := make([]map[string]interface{}, 0, len(instanceList.InstanceList))

	if instanceList != nil && instanceList.InstanceList != nil {
		for _, clbInstanceDetail := range instanceList.InstanceList {
			clbInstanceDetailMap := map[string]interface{}{}

			if clbInstanceDetail.LoadBalancerId != nil {
				clbInstanceDetailMap["load_balancer_id"] = clbInstanceDetail.LoadBalancerId
			}

			if clbInstanceDetail.LoadBalancerName != nil {
				clbInstanceDetailMap["load_balancer_name"] = clbInstanceDetail.LoadBalancerName
			}

			if clbInstanceDetail.Listeners != nil {
				listenersList := []interface{}{}
				for _, listeners := range clbInstanceDetail.Listeners {
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

				clbInstanceDetailMap["listeners"] = []interface{}{listenersList}
			}

			ids = append(ids, *clbInstanceDetail.LoadBalancerId)
			tmpList = append(tmpList, clbInstanceDetailMap)
		}

		_ = d.Set("instance_list", tmpList)
	}

	if instanceList.AsyncTotalNum != nil {
		_ = d.Set("async_total_num", instanceList.AsyncTotalNum)
	}

	if instanceList.AsyncOffset != nil {
		_ = d.Set("async_offset", instanceList.AsyncOffset)
	}

	if instanceList.AsyncCacheTime != nil {
		_ = d.Set("async_cache_time", instanceList.AsyncCacheTime)
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
