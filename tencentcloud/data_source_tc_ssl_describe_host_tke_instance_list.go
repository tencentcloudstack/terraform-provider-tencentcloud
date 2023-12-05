package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslDescribeHostTkeInstanceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostTkeInstanceListRead,
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

func dataSourceTencentCloudSslDescribeHostTkeInstanceListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_host_tke_instance_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceList []*ssl.TkeInstanceDetail
	var asyncTotalNum, asyncOffset *int64
	var asyncCacheTime *string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, total, offset, cacheTime, e := service.DescribeSslDescribeHostTkeInstanceListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceList = result
		asyncTotalNum = total
		asyncOffset = offset
		asyncCacheTime = cacheTime
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList))
	tmpList := make([]map[string]interface{}, 0, len(instanceList))

	if instanceList != nil {
		for _, tkeInstanceDetail := range instanceList {
			tkeInstanceDetailMap := map[string]interface{}{}

			if tkeInstanceDetail.ClusterId != nil {
				tkeInstanceDetailMap["cluster_id"] = tkeInstanceDetail.ClusterId
			}

			if tkeInstanceDetail.ClusterName != nil {
				tkeInstanceDetailMap["cluster_name"] = tkeInstanceDetail.ClusterName
			}

			if tkeInstanceDetail.NamespaceList != nil {
				namespaceListList := []interface{}{}
				for _, namespaceList := range tkeInstanceDetail.NamespaceList {
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

				tkeInstanceDetailMap["namespace_list"] = []interface{}{namespaceListList}
			}

			if tkeInstanceDetail.ClusterType != nil {
				tkeInstanceDetailMap["cluster_type"] = tkeInstanceDetail.ClusterType
			}

			if tkeInstanceDetail.ClusterVersion != nil {
				tkeInstanceDetailMap["cluster_version"] = tkeInstanceDetail.ClusterVersion
			}

			ids = append(ids, *tkeInstanceDetail.ClusterId)
			tmpList = append(tmpList, tkeInstanceDetailMap)
		}

		_ = d.Set("instance_list", tmpList)
	}

	if asyncTotalNum != nil {
		_ = d.Set("async_total_num", asyncTotalNum)
	}

	if asyncOffset != nil {
		_ = d.Set("async_offset", asyncOffset)
	}

	if asyncCacheTime != nil {
		_ = d.Set("async_cache_time", asyncCacheTime)
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
