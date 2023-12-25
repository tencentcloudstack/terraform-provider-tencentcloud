package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfDeliveryConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfDeliveryConfigsRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search word.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "deploy group information about the deployment group associated with a delivery item.Note: This field may return null, which means that no valid value was obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "total count. Note: This field may return null, which means that no valid value was obtained.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "content. Note: This field may return null, which means that no valid value was obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "config id.",
									},
									"config_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "config name.",
									},
									"collect_path": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "harvest log path. Note: This field may return null, which means that no valid value was obtained.",
									},
									"groups": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Associated deployment group information.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Group Id.",
												},
												"group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Group Name.",
												},
												"cluster_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cluster type.",
												},
												"cluster_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cluster ID. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"cluster_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cluster Name. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"namespace_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Namespace Name. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"associate_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Associate Time. Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation time.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"kafka_v_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Kafka VIP. Note: This field may return null, which means that no valid value was obtained.",
									},
									"kafka_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "KafkaAddress refers to the address of a Kafka server.Note: This field may return null, which means that no valid value was obtained.",
									},
									"kafka_v_port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Kafka VPort. Note: This field may return null, which means that no valid value was obtained.",
									},
									"topic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Topic. Note: This field may return null, which means that no valid value was obtained.",
									},
									"line_rule": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Line Rule for log. Note: This field may return null, which means that no valid value was obtained.",
									},
									"custom_rule": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CustomRule specifies a custom line separator rule.Note: This field may return null, which means that no valid value was obtained.",
									},
									"enable_global_line_rule": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether a single row rule should be applied.Note: This field may return null, which means that no valid value was obtained.",
									},
									"enable_auth": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "whether use auth for kafka. Note: This field may return null, which means that no valid value was obtained.",
									},
									"username": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "user Name. Note: This field may return null, which means that no valid value was obtained.",
									},
									"password": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password. Note: This field may return null, which means that no valid value was obtained.",
									},
									"kafka_infos": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Kafka Infos. Note: This field may return null, which means that no valid value was obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"topic": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Kafka topic. Note: This field may return null, which means that no valid value was obtained.",
												},
												"path": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "harvest log path. Note: This field may return null, which means that no valid value was obtained.",
												},
												"line_rule": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Line rule specifies the type of line separator used in a file. It can have one of the following values: default: The default line separator is used to separate lines in the file. time: The lines in the file are separated based on time. custom: A custom line separator is used. In this case, the CustomRule field should be filled with the specific custom value. Note: This field may return null, which means that no valid value was obtained.",
												},
												"custom_rule": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Custom Line Rule.",
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

func dataSourceTencentCloudTsfDeliveryConfigsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_delivery_configs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var bindGroups *tsf.DeliveryConfigBindGroups

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfDeliveryConfigsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		bindGroups = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(bindGroups.Content))
	deliveryConfigBindGroupsMap := map[string]interface{}{}
	if bindGroups != nil {
		if bindGroups.TotalCount != nil {
			deliveryConfigBindGroupsMap["total_count"] = bindGroups.TotalCount
		}

		if bindGroups.Content != nil {
			contentList := []interface{}{}
			for _, content := range bindGroups.Content {
				contentMap := map[string]interface{}{}

				if content.ConfigId != nil {
					contentMap["config_id"] = *content.ConfigId
				}

				if content.ConfigName != nil {
					contentMap["config_name"] = content.ConfigName
				}

				if content.CollectPath != nil {
					contentMap["collect_path"] = content.CollectPath
				}

				if content.Groups != nil {
					groupsList := []interface{}{}
					for _, groups := range content.Groups {
						groupsMap := map[string]interface{}{}

						if groups.GroupId != nil {
							groupsMap["group_id"] = groups.GroupId
						}

						if groups.GroupName != nil {
							groupsMap["group_name"] = groups.GroupName
						}

						if groups.ClusterType != nil {
							groupsMap["cluster_type"] = groups.ClusterType
						}

						if groups.ClusterId != nil {
							groupsMap["cluster_id"] = groups.ClusterId
						}

						if groups.ClusterName != nil {
							groupsMap["cluster_name"] = groups.ClusterName
						}

						if groups.NamespaceName != nil {
							groupsMap["namespace_name"] = groups.NamespaceName
						}

						if groups.AssociateTime != nil {
							groupsMap["associate_time"] = groups.AssociateTime
						}

						groupsList = append(groupsList, groupsMap)
					}

					contentMap["groups"] = groupsList
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.KafkaVIp != nil {
					contentMap["kafka_v_ip"] = content.KafkaVIp
				}

				if content.KafkaAddress != nil {
					contentMap["kafka_address"] = content.KafkaAddress
				}

				if content.KafkaVPort != nil {
					contentMap["kafka_v_port"] = content.KafkaVPort
				}

				if content.Topic != nil {
					contentMap["topic"] = content.Topic
				}

				if content.LineRule != nil {
					contentMap["line_rule"] = content.LineRule
				}

				if content.CustomRule != nil {
					contentMap["custom_rule"] = content.CustomRule
				}

				if content.EnableGlobalLineRule != nil {
					contentMap["enable_global_line_rule"] = content.EnableGlobalLineRule
				}

				if content.EnableAuth != nil {
					contentMap["enable_auth"] = content.EnableAuth
				}

				if content.Username != nil {
					contentMap["username"] = content.Username
				}

				if content.Password != nil {
					contentMap["password"] = content.Password
				}

				if content.KafkaInfos != nil {
					kafkaInfosList := []interface{}{}
					for _, kafkaInfos := range content.KafkaInfos {
						kafkaInfosMap := map[string]interface{}{}

						if kafkaInfos.Topic != nil {
							kafkaInfosMap["topic"] = kafkaInfos.Topic
						}

						if kafkaInfos.Path != nil {
							kafkaInfosMap["path"] = kafkaInfos.Path
						}

						if kafkaInfos.LineRule != nil {
							kafkaInfosMap["line_rule"] = kafkaInfos.LineRule
						}

						if kafkaInfos.CustomRule != nil {
							kafkaInfosMap["custom_rule"] = kafkaInfos.CustomRule
						}

						kafkaInfosList = append(kafkaInfosList, kafkaInfosMap)
					}

					contentMap["kafka_infos"] = kafkaInfosList
					ids = append(ids, *content.ConfigId)
				}

				contentList = append(contentList, contentMap)
			}

			deliveryConfigBindGroupsMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{deliveryConfigBindGroupsMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), deliveryConfigBindGroupsMap); e != nil {
			return e
		}
	}
	return nil
}
