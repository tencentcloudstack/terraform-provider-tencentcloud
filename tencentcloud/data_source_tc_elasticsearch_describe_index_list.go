package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudElasticsearchDescribeIndexList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticsearchDescribeIndexListRead,
		Schema: map[string]*schema.Schema{
			"index_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Index type. `auto`: Autonomous index; `normal`: General index.",
			},

			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ES cluster id.",
			},

			"index_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Index name. If you fill in the blanks, get all indexes.",
			},

			"username": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster access user name.",
			},

			"password": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster access password.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort field. Support index name: IndexName, index storage: IndexStorage, index creation time: IndexCreateTime.",
			},

			"index_status_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Index status list.",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort order, which supports asc and desc. The default is desc data format asc,desc.",
			},

			"index_meta_fields": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Index metadata field.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Index type.",
						},
						"index_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Index name.",
						},
						"index_meta_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Index meta json.",
						},
						"index_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Index status.",
						},
						"index_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Index storage.",
						},
						"index_create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Index create time.",
						},
						"backing_indices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Backing indices.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Index name.",
									},
									"index_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Index status.",
									},
									"index_storage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Index storage.",
									},
									"index_phrase": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Index phrase.",
									},
									"index_create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Index create time.",
									},
								},
							},
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster id.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster name.",
						},
						"cluster_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster version.",
						},
						"index_policy_field": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Index lifecycle field.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"warm_enable": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to enable warm.",
									},
									"warm_min_age": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Warm phase transition time.",
									},
									"cold_enable": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to enable the cold phase.",
									},
									"cold_min_age": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cold phase transition time.",
									},
									"frozen_enable": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Start frozen phase.",
									},
									"frozen_min_age": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Frozen phase transition time.",
									},
									"cold_action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cold action.",
									},
								},
							},
						},
						"index_options_field": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Index options field.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expire_max_age": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire max age.",
									},
									"expire_max_size": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire max size.",
									},
									"rollover_max_age": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rollover max age.",
									},
									"rollover_dynamic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to turn on dynamic scrolling.",
									},
									"shard_num_dynamic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to enable dynamic slicing.",
									},
									"timestamp_field": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time partition field.",
									},
									"write_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Write mode.",
									},
								},
							},
						},
						"index_settings_field": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Index settings field.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"number_of_shards": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Number of index main fragments.",
									},
									"number_of_replicas": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Number of index copy fragments.",
									},
									"refresh_interval": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Index refresh frequency.",
									},
								},
							},
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "App id.",
						},
						"index_docs": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of indexed documents.",
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

func dataSourceTencentCloudElasticsearchDescribeIndexListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_elasticsearch_describe_index_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("index_type"); ok {
		paramMap["IndexType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("index_name"); ok {
		paramMap["IndexName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		paramMap["Username"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		paramMap["Password"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("index_status_list"); ok {
		indexStatusListSet := v.(*schema.Set).List()
		paramMap["IndexStatusList"] = helper.InterfacesStringsPoint(indexStatusListSet)
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}

	var indexMetaFields []*elasticsearch.IndexMetaField

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeElasticsearchDescribeIndexListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		indexMetaFields = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(indexMetaFields))
	tmpList := make([]map[string]interface{}, 0, len(indexMetaFields))

	if indexMetaFields != nil {
		for _, indexMetaField := range indexMetaFields {
			indexMetaFieldMap := map[string]interface{}{}

			if indexMetaField.IndexType != nil {
				indexMetaFieldMap["index_type"] = indexMetaField.IndexType
			}

			if indexMetaField.IndexName != nil {
				indexMetaFieldMap["index_name"] = indexMetaField.IndexName
			}

			if indexMetaField.IndexMetaJson != nil {
				indexMetaFieldMap["index_meta_json"] = indexMetaField.IndexMetaJson
			}

			if indexMetaField.IndexStatus != nil {
				indexMetaFieldMap["index_status"] = indexMetaField.IndexStatus
			}

			if indexMetaField.IndexStorage != nil {
				indexMetaFieldMap["index_storage"] = indexMetaField.IndexStorage
			}

			if indexMetaField.IndexCreateTime != nil {
				indexMetaFieldMap["index_create_time"] = indexMetaField.IndexCreateTime
			}

			if indexMetaField.BackingIndices != nil {
				backingIndicesList := []interface{}{}
				for _, backingIndices := range indexMetaField.BackingIndices {
					backingIndicesMap := map[string]interface{}{}

					if backingIndices.IndexName != nil {
						backingIndicesMap["index_name"] = backingIndices.IndexName
					}

					if backingIndices.IndexStatus != nil {
						backingIndicesMap["index_status"] = backingIndices.IndexStatus
					}

					if backingIndices.IndexStorage != nil {
						backingIndicesMap["index_storage"] = backingIndices.IndexStorage
					}

					if backingIndices.IndexPhrase != nil {
						backingIndicesMap["index_phrase"] = backingIndices.IndexPhrase
					}

					if backingIndices.IndexCreateTime != nil {
						backingIndicesMap["index_create_time"] = backingIndices.IndexCreateTime
					}

					backingIndicesList = append(backingIndicesList, backingIndicesMap)
				}

				indexMetaFieldMap["backing_indices"] = []interface{}{backingIndicesList}
			}

			if indexMetaField.ClusterId != nil {
				indexMetaFieldMap["cluster_id"] = indexMetaField.ClusterId
			}

			if indexMetaField.ClusterName != nil {
				indexMetaFieldMap["cluster_name"] = indexMetaField.ClusterName
			}

			if indexMetaField.ClusterVersion != nil {
				indexMetaFieldMap["cluster_version"] = indexMetaField.ClusterVersion
			}

			if indexMetaField.IndexPolicyField != nil {
				indexPolicyFieldMap := map[string]interface{}{}

				if indexMetaField.IndexPolicyField.WarmEnable != nil {
					indexPolicyFieldMap["warm_enable"] = indexMetaField.IndexPolicyField.WarmEnable
				}

				if indexMetaField.IndexPolicyField.WarmMinAge != nil {
					indexPolicyFieldMap["warm_min_age"] = indexMetaField.IndexPolicyField.WarmMinAge
				}

				if indexMetaField.IndexPolicyField.ColdEnable != nil {
					indexPolicyFieldMap["cold_enable"] = indexMetaField.IndexPolicyField.ColdEnable
				}

				if indexMetaField.IndexPolicyField.ColdMinAge != nil {
					indexPolicyFieldMap["cold_min_age"] = indexMetaField.IndexPolicyField.ColdMinAge
				}

				if indexMetaField.IndexPolicyField.FrozenEnable != nil {
					indexPolicyFieldMap["frozen_enable"] = indexMetaField.IndexPolicyField.FrozenEnable
				}

				if indexMetaField.IndexPolicyField.FrozenMinAge != nil {
					indexPolicyFieldMap["frozen_min_age"] = indexMetaField.IndexPolicyField.FrozenMinAge
				}

				if indexMetaField.IndexPolicyField.ColdAction != nil {
					indexPolicyFieldMap["cold_action"] = indexMetaField.IndexPolicyField.ColdAction
				}

				indexMetaFieldMap["index_policy_field"] = []interface{}{indexPolicyFieldMap}
			}

			if indexMetaField.IndexOptionsField != nil {
				indexOptionsFieldMap := map[string]interface{}{}

				if indexMetaField.IndexOptionsField.ExpireMaxAge != nil {
					indexOptionsFieldMap["expire_max_age"] = indexMetaField.IndexOptionsField.ExpireMaxAge
				}

				if indexMetaField.IndexOptionsField.ExpireMaxSize != nil {
					indexOptionsFieldMap["expire_max_size"] = indexMetaField.IndexOptionsField.ExpireMaxSize
				}

				if indexMetaField.IndexOptionsField.RolloverMaxAge != nil {
					indexOptionsFieldMap["rollover_max_age"] = indexMetaField.IndexOptionsField.RolloverMaxAge
				}

				if indexMetaField.IndexOptionsField.RolloverDynamic != nil {
					indexOptionsFieldMap["rollover_dynamic"] = indexMetaField.IndexOptionsField.RolloverDynamic
				}

				if indexMetaField.IndexOptionsField.ShardNumDynamic != nil {
					indexOptionsFieldMap["shard_num_dynamic"] = indexMetaField.IndexOptionsField.ShardNumDynamic
				}

				if indexMetaField.IndexOptionsField.TimestampField != nil {
					indexOptionsFieldMap["timestamp_field"] = indexMetaField.IndexOptionsField.TimestampField
				}

				if indexMetaField.IndexOptionsField.WriteMode != nil {
					indexOptionsFieldMap["write_mode"] = indexMetaField.IndexOptionsField.WriteMode
				}

				indexMetaFieldMap["index_options_field"] = []interface{}{indexOptionsFieldMap}
			}

			if indexMetaField.IndexSettingsField != nil {
				indexSettingsFieldMap := map[string]interface{}{}

				if indexMetaField.IndexSettingsField.NumberOfShards != nil {
					indexSettingsFieldMap["number_of_shards"] = indexMetaField.IndexSettingsField.NumberOfShards
				}

				if indexMetaField.IndexSettingsField.NumberOfReplicas != nil {
					indexSettingsFieldMap["number_of_replicas"] = indexMetaField.IndexSettingsField.NumberOfReplicas
				}

				if indexMetaField.IndexSettingsField.RefreshInterval != nil {
					indexSettingsFieldMap["refresh_interval"] = indexMetaField.IndexSettingsField.RefreshInterval
				}

				indexMetaFieldMap["index_settings_field"] = []interface{}{indexSettingsFieldMap}
			}

			if indexMetaField.AppId != nil {
				indexMetaFieldMap["app_id"] = indexMetaField.AppId
			}

			if indexMetaField.IndexDocs != nil {
				indexMetaFieldMap["index_docs"] = indexMetaField.IndexDocs
			}

			ids = append(ids, *indexMetaField.ClusterId+FILED_SP+*indexMetaField.IndexName)
			tmpList = append(tmpList, indexMetaFieldMap)
		}

		_ = d.Set("index_meta_fields", tmpList)
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
