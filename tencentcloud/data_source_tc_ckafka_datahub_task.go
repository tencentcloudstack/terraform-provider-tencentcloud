/*
Use this data source to query detailed information of ckafka datahub_task

Example Usage

```hcl
data "tencentcloud_ckafka_datahub_task" "datahub_task" {
  limit = 20
  offset = 0
  search_word = "SearchWord"
  target_type = "CKafka"
  task_type = "SOURCE"
  source_type = "CKafka"
  resource = "Resource"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCkafkaDatahubTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaDatahubTaskRead,
		Schema: map[string]*schema.Schema{
			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return number, the default is 20, the maximum is 100.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset, default is 0.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search key.",
			},

			"target_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Destination type of dump.",
			},

			"task_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Task type, SOURCE|SINK.",
			},

			"source_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The source type.",
			},

			"resource": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Resource.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Query result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total Count.",
						},
						"task_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Datahub task information list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"task_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Task ID.",
									},
									"task_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TaskName.",
									},
									"task_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TaskType，SOURCE|SINK.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Status, -1 failed to create, 0 to create, 1 to run, 2 to delete, 3 to deleted, 4 to delete failed, 5 to pause, 6 to pause, 7 to pause, 8 to resume, 9 to resume failed.",
									},
									"source_resource": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Data resource.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource type.",
												},
												"kafka_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Ckafka configuration，required when Type is KAFKA.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"self_built": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the cluster is built by yourself instead of cloud product.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance resource.",
															},
															"topic": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Topic name，use “,” when more than 1 topic.",
															},
															"offset_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Offset type，from beginning:earliest，from latest:latest，from specific time:timestamp.",
															},
															"start_time": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "When Offset type timestamp is required.",
															},
															"resource_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance name.",
															},
															"zone_id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Zone ID.",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Topic的Id.",
															},
															"partition_num": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The partition num of the topic.",
															},
															"enable_toleration": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Enable dead letter queue.",
															},
															"qps_limit": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Qps(query per seconds) limit.",
															},
															"table_mappings": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Maps of table to topic, required when multi topic is selected.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"database": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Database name.",
																		},
																		"table": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Table name,use , to seperate.",
																		},
																		"topic": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Topic name.",
																		},
																		"topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Topic ID.",
																		},
																	},
																},
															},
															"use_table_mapping": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to use multi table。.",
															},
															"use_auto_create_topic": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Does the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).",
															},
															"compression_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.",
															},
															"msg_multiple": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
															},
															"connector_sync_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ConnectorSyncType.",
															},
															"keep_partition": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "KeepPartition.",
															},
														},
													},
												},
												"event_bus_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "EB configuration，required when type is EB.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource type。EB_COS/EB_ES/EB_CLS.",
															},
															"self_built": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether it is a self-built cluster.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance id.",
															},
															"namespace": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF namespace.",
															},
															"function_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF function name.",
															},
															"qualifier": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF version and alias.",
															},
														},
													},
												},
												"mongo_d_b_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "MongoDB config，Required when Type is MONGODB.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MongoDB database name.",
															},
															"collection": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MongoDB collection.",
															},
															"copy_existing": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to copy the stock data, the default parameter is true.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource id.",
															},
															"ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Mongo DB connection ip.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "MongoDB connection port.",
															},
															"user_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MongoDB database user name.",
															},
															"password": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MongoDB database password.",
															},
															"listening_event": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Listening event type, if it is empty, it means select all. Values include insert, update, replace, delete, invalidate, drop, dropdatabase, rename, used between multiple types, separated by commas.",
															},
															"read_preference": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Master-slave priority, default master node.",
															},
															"pipeline": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Aggregation pipeline.",
															},
															"self_built": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether it is a self-built cluster.",
															},
														},
													},
												},
												"es_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Es configuration, required when Type is ES.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Es connection port.",
															},
															"user_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Es UserName.",
															},
															"password": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Es Password.",
															},
															"self_built": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether it is a self-built cluster.",
															},
															"service_vip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance vip.",
															},
															"uniq_vpc_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance vpc id.",
															},
															"drop_invalid_message": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether Es discards the message of parsing failure.",
															},
															"index": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Es index name.",
															},
															"date_format": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Es date suffix.",
															},
															"content_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Key for data in non-json format.",
															},
															"drop_invalid_json_message": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether Es discards messages in non-json format.",
															},
															"document_id_field": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The field name of the document ID value dumped into Es.",
															},
															"index_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Es custom index name type, STRING, JSONPATH, the default is STRING.",
															},
															"drop_cls": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"drop_invalid_message_to_cls": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether to deliver to cls.",
																		},
																		"drop_cls_region": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The region where the cls is delivered.",
																		},
																		"drop_cls_owneruin": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Delivery account of cls.",
																		},
																		"drop_cls_topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Topic of cls.",
																		},
																		"drop_cls_log_set": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls log set.",
																		},
																	},
																},
															},
															"database_primary_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "When the message dumped to ES is the binlog of Database, if you need to synchronize database operations, that is, fill in the primary key of the database table when adding, deleting, and modifying operations to ES.",
															},
															"drop_dlq": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Dead letter queue.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type，DLQ dead letter queue，IGNORE_ERROR|DROP.",
																		},
																		"kafka_param": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Ckafka type dlq.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"self_built": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Whether it is a self-built cluster.",
																					},
																					"resource": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Resource id.",
																					},
																					"topic": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Topic name, multiple separated by ,.",
																					},
																					"offset_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
																					},
																					"start_time": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
																					},
																					"resource_name": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Resource id name.",
																					},
																					"zone_id": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Zone ID.",
																					},
																					"topic_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Topic Id.",
																					},
																					"partition_num": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Topic的分区数.",
																					},
																					"enable_toleration": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Enable the fault-tolerant instance and enable the dead-letter queue.",
																					},
																					"qps_limit": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Qps limit.",
																					},
																					"table_mappings": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"database": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Database name.",
																								},
																								"table": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Table name, multiple tables, separated by (commas).",
																								},
																								"topic": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Topic name.",
																								},
																								"topic_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Topic ID.",
																								},
																							},
																						},
																					},
																					"use_table_mapping": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Distribute to multiple topics switch, the default is false.",
																					},
																					"use_auto_create_topic": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).",
																					},
																					"compression_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open。.",
																					},
																					"msg_multiple": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
																					},
																					"connector_sync_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "ConnectorSyncType.",
																					},
																					"keep_partition": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "KeepPartition.",
																					},
																				},
																			},
																		},
																		"retry_interval": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Retry interval.",
																		},
																		"max_retry_attempts": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Retry times.",
																		},
																		"topic_param": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "DIP Topic type dead letter queue.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"resource": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The topic name of the topic sold separately.",
																					},
																					"offset_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
																					},
																					"start_time": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
																					},
																					"topic_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "TopicId.",
																					},
																					"compression_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.",
																					},
																					"use_auto_create_topic": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).",
																					},
																					"msg_multiple": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
																					},
																				},
																			},
																		},
																		"dlq_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Dlq type，CKAFKA|TOPIC.",
																		},
																	},
																},
															},
														},
													},
												},
												"tdw_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Tdw configuration, required when Type is TDW.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bid": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Tdw bid.",
															},
															"tid": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Tdw tid.",
															},
															"is_domestic": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Default true.",
															},
															"tdw_host": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "TDW address，defalt tl-tdbank-tdmanager.tencent-distribute.com.",
															},
															"tdw_port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "TDW port，default 8099.",
															},
														},
													},
												},
												"dts_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Dts configuration, required when Type is DTS.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts instance Id.",
															},
															"ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts connection ip.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Dts connection port.",
															},
															"topic": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts topic.",
															},
															"group_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts consumer group Id.",
															},
															"group_user": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts account.",
															},
															"group_password": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts consumer group passwd.",
															},
															"tran_sql": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "False to synchronize the original data, true to synchronize the parsed json format data, the default is true.",
															},
														},
													},
												},
												"click_house_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "ClickHouse config，Type CLICKHOUSE requierd.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse cluster.",
															},
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse database name.",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse table.",
															},
															"schema": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "ClickHouse schema.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"column_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Column name.",
																		},
																		"json_key": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The json Key name corresponding to this column.",
																		},
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of table column.",
																		},
																		"allow_null": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether the column item is allowed to be empty.",
																		},
																	},
																},
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource id.",
															},
															"ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse ip.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ClickHouse port.",
															},
															"user_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse user name.",
															},
															"password": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse passwd.",
															},
															"service_vip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance vip.",
															},
															"uniq_vpc_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance vpc id.",
															},
															"self_built": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether it is a self-built cluster.",
															},
															"drop_invalid_message": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether ClickHouse discards the message that fails to parse, the default is true.",
															},
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse type，emr-clickhouse : emr;cdw-clickhouse : cdwch;selfBuilt : .",
															},
															"drop_cls": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"drop_invalid_message_to_cls": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether to deliver to cls.",
																		},
																		"drop_cls_region": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls region.",
																		},
																		"drop_cls_owneruin": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls account.",
																		},
																		"drop_cls_topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls topicId.",
																		},
																		"drop_cls_log_set": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls LogSet id.",
																		},
																	},
																},
															},
														},
													},
												},
												"cls_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Cls configuration，Required when Type is CLS.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"decode_json": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the produced information is in json format.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cls id.",
															},
															"log_set": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "LogSet id.",
															},
															"content_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Required when Decode Json is false.",
															},
															"time_field": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specify the content of a field in the message as the time of the cls log. The format of the field content needs to be a second-level timestamp.",
															},
														},
													},
												},
												"cos_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Cos configuration, required when Type is COS.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cos bucket name.",
															},
															"region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Region code.",
															},
															"object_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ObjectKey.",
															},
															"aggregate_batch_size": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The size of aggregated messages MB.",
															},
															"aggregate_interval": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Time interval.",
															},
															"format_output_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The file format after message aggregation csv|json.",
															},
															"object_key_prefix": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dumped object directory prefix.",
															},
															"directory_time_format": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Partition format formatted according to strptime time.",
															},
														},
													},
												},
												"my_s_q_l_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "MySQL configuration，Required when Type is MYSQL.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MySQL database name，* is the whole database.",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name of the MySQL data table,  is the non-system table in all the monitored databases, which can be separated by , to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name , when a regular expression needs to be filled in, the format is data database name.data table name.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MySQL connection Id.",
															},
															"snapshot_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to Copy inventory information (schema_only does not copy, initial full amount), the default is initial.",
															},
															"ddl_topic": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The Topic that stores the Ddl information of My SQL, if it is empty, it will not be stored by default.",
															},
															"data_source_monitor_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "TABLE indicates that the read item is a table, QUERY indicates that the read item is a query.",
															},
															"data_source_monitor_resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "When DataMonitorMode=TABLE, pass in the Table that needs to be read; when DataMonitorMode=QUERY, pass in the query sql statement that needs to be read.",
															},
															"data_source_increment_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "TIMESTAMP indicates that the incremental column is of timestamp type, INCREMENT indicates that the incremental column is of self-incrementing id type&amp;#39;.",
															},
															"data_source_increment_column": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name of the column to be monitored.",
															},
															"data_source_start_from": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "HEAD means copy stock + incremental data, TAIL means copy only incremental data.",
															},
															"data_target_insert_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.",
															},
															"data_target_primary_key_field": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.",
															},
															"data_target_record_mapping": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Mapping relationship between tables and messages.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"json_key": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The key name of the message.",
																		},
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Message type.",
																		},
																		"allow_null": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether the message is allowed to be empty.",
																		},
																		"column_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Corresponding mapping column name.",
																		},
																		"extra_info": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Database table extra fields.",
																		},
																		"column_size": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Current column size.",
																		},
																		"decimal_digits": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Current column precision.",
																		},
																		"auto_increment": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether it is an auto-increment column.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Database table default parameters.",
																		},
																	},
																},
															},
															"topic_regex": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regular expression for routing events to specific topics, defaults to (.*).",
															},
															"topic_replacement": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "TopicRegex ，$1、$2.",
															},
															"key_columns": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Format library1.table1 field 1,field 2;library 2.table2 field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
															},
															"drop_invalid_message": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to discard messages that fail to parse, the default is true.",
															},
															"drop_cls": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"drop_invalid_message_to_cls": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether to deliver to cls.",
																		},
																		"drop_cls_region": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The region where the cls is delivered.",
																		},
																		"drop_cls_owneruin": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Account.",
																		},
																		"drop_cls_topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls topic.",
																		},
																		"drop_cls_log_set": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls LogSet id.",
																		},
																	},
																},
															},
															"output_format": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Output format，DEFAULT、CANAL_1、CANAL_2.",
															},
															"is_table_prefix": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "When the Table input is a prefix, the value of this item is true, otherwise it is false.",
															},
															"include_content_changes": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.",
															},
															"include_query": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.",
															},
															"record_with_schema": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
															},
															"signal_database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Database name of signal table.",
															},
															"is_table_regular": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the input table is a regular expression, if this option and Is Table Prefix are true at the same time, the judgment priority of this option is higher than Is Table Prefix.",
															},
														},
													},
												},
												"postgre_s_q_l_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "PostgreSQL configuration，Required when Type is POSTGRESQL or TDSQL C_POSTGRESQL.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "PostgreSQL database name.",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "PostgreSQL tableName，* is the non-system table in all the monitored databases, you can use , to monitor multiple data tables, but the data table needs to be filled in the format of Schema name.Data table name, and you need to fill in a regular expression When, the format is Schema name.data table name.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "PostgreSQL connection Id.",
															},
															"plugin_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "(decoderbufs/pgoutput)，default decoderbufs.",
															},
															"snapshot_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Never|initial, default initial.",
															},
															"data_format": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Upstream data format (JSON|Debezium), required when the database synchronization mode matches the default field.",
															},
															"data_target_insert_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.",
															},
															"data_target_primary_key_field": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.",
															},
															"data_target_record_mapping": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Mapping relationship between tables and messages.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"json_key": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The key name of the message.",
																		},
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Message type.",
																		},
																		"allow_null": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether the message is allowed to be empty.",
																		},
																		"column_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Column Name.",
																		},
																		"extra_info": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Database table extra fields.",
																		},
																		"column_size": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Current ColumnSize.",
																		},
																		"decimal_digits": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Current Column DecimalDigits.",
																		},
																		"auto_increment": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether it is an auto-increment column.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Database table default parameters.",
																		},
																	},
																},
															},
															"drop_invalid_message": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to discard messages that fail to parse, the default is true.",
															},
															"is_table_regular": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the input table is a regular expression.",
															},
															"key_columns": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Format  library1.table1:field 1,field2;library2.table2:field2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
															},
															"record_with_schema": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
															},
														},
													},
												},
												"topic_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Topic configuration，Required when Type is Topic.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The topic name of the topic sold separately.",
															},
															"offset_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
															},
															"start_time": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Topic TopicId.",
															},
															"compression_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.",
															},
															"use_auto_create_topic": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).",
															},
															"msg_multiple": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
															},
														},
													},
												},
												"maria_d_b_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "MariaDB configuration，Required when Type is MARIADB.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MariaDB database name，* for all database.",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MariaDB db name，*is the non-system table in all the monitored databases, you can use , to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MariaDB connection Id.",
															},
															"snapshot_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Schema_only|initial，default initial.",
															},
															"key_columns": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Format  library 1. table 1: field 1, field 2; library 2. table 2: field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
															},
															"is_table_prefix": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "When the Table input is a prefix, the value of this item is true, otherwise it is false.",
															},
															"output_format": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Output format，DEFAULT、CANAL_1、CANAL_2.",
															},
															"include_content_changes": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.",
															},
															"include_query": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.",
															},
															"record_with_schema": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
															},
														},
													},
												},
												"s_q_l_server_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "SQLServer configuration，Required when Type is SQLSERVER.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SQLServer database name.",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SQLServer table，*is the non-system table in all the monitored databases, you can use , to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SQLServer connection Id.",
															},
															"snapshot_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Schema_only|initial default initial.",
															},
														},
													},
												},
												"ctsdb_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Ctsdb configuration，Required when Type is CTSDB.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource id.",
															},
															"ctsdb_metric": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Ctsdb metric.",
															},
														},
													},
												},
												"scf_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Scf configuration，Required when Type is SCF.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"function_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF function name.",
															},
															"namespace": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF cloud function namespace, the default is default.",
															},
															"qualifier": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF cloud function version and alias, the default is DEFAULT.",
															},
															"batch_size": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The maximum number of messages sent in each batch, the default is 1000.",
															},
															"max_retries": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The number of retries after the SCF call fails, the default is 5.",
															},
														},
													},
												},
											},
										},
									},
									"target_resource": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Target Resource.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource Type.",
												},
												"kafka_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Ckafka configuration，required when Type is KAFKA.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"self_built": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the cluster is built by yourself instead of cloud product.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance resource.",
															},
															"topic": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Topic name，use “,” when more than 1 topic.",
															},
															"offset_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Offset type，from beginning:earliest，from latest:latest，from specific time:timestamp.",
															},
															"start_time": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "When Offset type timestamp is required.",
															},
															"resource_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance name.",
															},
															"zone_id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Zone ID.",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Topic的Id.",
															},
															"partition_num": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The partition num of the topic.",
															},
															"enable_toleration": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Enable dead letter queue.",
															},
															"qps_limit": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Qps(query per seconds) limit.",
															},
															"table_mappings": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Maps of table to topic, required when multi topic is selected.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"database": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Database name.",
																		},
																		"table": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Table name,use , to seperate.",
																		},
																		"topic": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Topic name.",
																		},
																		"topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Topic ID.",
																		},
																	},
																},
															},
															"use_table_mapping": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to use multi table。.",
															},
															"use_auto_create_topic": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Does the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).",
															},
															"compression_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.",
															},
															"msg_multiple": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
															},
															"connector_sync_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ConnectorSyncType.",
															},
															"keep_partition": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "KeepPartition.",
															},
														},
													},
												},
												"event_bus_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "EB configuration，required when type is EB.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource type。EB_COS/EB_ES/EB_CLS.",
															},
															"self_built": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether it is a self-built cluster.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance id.",
															},
															"namespace": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF namespace.",
															},
															"function_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF function name.",
															},
															"qualifier": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF version and alias.",
															},
														},
													},
												},
												"mongo_d_b_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "MongoDB config，Required when Type is MONGODB.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MongoDB database name.",
															},
															"collection": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MongoDB collection.",
															},
															"copy_existing": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to copy the stock data, the default parameter is true.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource id.",
															},
															"ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Mongo DB connection ip.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "MongoDB connection port.",
															},
															"user_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MongoDB database user name.",
															},
															"password": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MongoDB database password.",
															},
															"listening_event": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Listening event type, if it is empty, it means select all. Values include insert, update, replace, delete, invalidate, drop, dropdatabase, rename, used between multiple types, separated by commas.",
															},
															"read_preference": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Master-slave priority, default master node.",
															},
															"pipeline": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Aggregation pipeline.",
															},
															"self_built": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether it is a self-built cluster.",
															},
														},
													},
												},
												"es_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Es configuration, required when Type is ES.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Es connection port.",
															},
															"user_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Es UserName.",
															},
															"password": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Es Password.",
															},
															"self_built": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether it is a self-built cluster.",
															},
															"service_vip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance vip.",
															},
															"uniq_vpc_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance vpc id.",
															},
															"drop_invalid_message": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether Es discards the message of parsing failure.",
															},
															"index": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Es index name.",
															},
															"date_format": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Es date suffix.",
															},
															"content_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Key for data in non-json format.",
															},
															"drop_invalid_json_message": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether Es discards messages in non-json format.",
															},
															"document_id_field": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The field name of the document ID value dumped into Es.",
															},
															"index_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Es custom index name type, STRING, JSONPATH, the default is STRING.",
															},
															"drop_cls": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"drop_invalid_message_to_cls": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether to deliver to cls.",
																		},
																		"drop_cls_region": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The region where the cls is delivered.",
																		},
																		"drop_cls_owneruin": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Delivery account of cls.",
																		},
																		"drop_cls_topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Topic of cls.",
																		},
																		"drop_cls_log_set": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls log set.",
																		},
																	},
																},
															},
															"database_primary_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "When the message dumped to ES is the binlog of Database, if you need to synchronize database operations, that is, fill in the primary key of the database table when adding, deleting, and modifying operations to ES.",
															},
															"drop_dlq": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Dead letter queue.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type，DLQ dead letter queue，IGNORE_ERROR|DROP.",
																		},
																		"kafka_param": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Ckafka type dlq.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"self_built": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Whether it is a self-built cluster.",
																					},
																					"resource": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Resource id.",
																					},
																					"topic": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Topic name, multiple separated by ,.",
																					},
																					"offset_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
																					},
																					"start_time": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
																					},
																					"resource_name": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Resource id name.",
																					},
																					"zone_id": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Zone ID.",
																					},
																					"topic_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Topic Id.",
																					},
																					"partition_num": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Topic的分区数.",
																					},
																					"enable_toleration": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Enable the fault-tolerant instance and enable the dead-letter queue.",
																					},
																					"qps_limit": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Qps limit.",
																					},
																					"table_mappings": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"database": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Database name.",
																								},
																								"table": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Table name, multiple tables, separated by (commas).",
																								},
																								"topic": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Topic name.",
																								},
																								"topic_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Topic ID.",
																								},
																							},
																						},
																					},
																					"use_table_mapping": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Distribute to multiple topics switch, the default is false.",
																					},
																					"use_auto_create_topic": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).",
																					},
																					"compression_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open。.",
																					},
																					"msg_multiple": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
																					},
																					"connector_sync_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "ConnectorSyncType.",
																					},
																					"keep_partition": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "KeepPartition.",
																					},
																				},
																			},
																		},
																		"retry_interval": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Retry interval.",
																		},
																		"max_retry_attempts": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Retry times.",
																		},
																		"topic_param": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "DIP Topic type dead letter queue.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"resource": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The topic name of the topic sold separately.",
																					},
																					"offset_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
																					},
																					"start_time": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
																					},
																					"topic_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "TopicId.",
																					},
																					"compression_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.",
																					},
																					"use_auto_create_topic": {
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).",
																					},
																					"msg_multiple": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
																					},
																				},
																			},
																		},
																		"dlq_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Dlq type，CKAFKA|TOPIC.",
																		},
																	},
																},
															},
														},
													},
												},
												"tdw_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Tdw configuration, required when Type is TDW.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bid": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Tdw bid.",
															},
															"tid": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Tdw tid.",
															},
															"is_domestic": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Default true.",
															},
															"tdw_host": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "TDW address，defalt tl-tdbank-tdmanager.tencent-distribute.com.",
															},
															"tdw_port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "TDW port，default 8099.",
															},
														},
													},
												},
												"dts_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Dts configuration, required when Type is DTS.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts instance Id.",
															},
															"ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts connection ip.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Dts connection port.",
															},
															"topic": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts topic.",
															},
															"group_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts consumer group Id.",
															},
															"group_user": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts account.",
															},
															"group_password": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dts consumer group passwd.",
															},
															"tran_sql": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "False to synchronize the original data, true to synchronize the parsed json format data, the default is true.",
															},
														},
													},
												},
												"click_house_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "ClickHouse config，Type CLICKHOUSE requierd.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse cluster.",
															},
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse database name.",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse table.",
															},
															"schema": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "ClickHouse schema.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"column_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Column name.",
																		},
																		"json_key": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The json Key name corresponding to this column.",
																		},
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of table column.",
																		},
																		"allow_null": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether the column item is allowed to be empty.",
																		},
																	},
																},
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource id.",
															},
															"ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse ip.",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ClickHouse port.",
															},
															"user_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse user name.",
															},
															"password": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse passwd.",
															},
															"service_vip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance vip.",
															},
															"uniq_vpc_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Instance vpc id.",
															},
															"self_built": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether it is a self-built cluster.",
															},
															"drop_invalid_message": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether ClickHouse discards the message that fails to parse, the default is true.",
															},
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ClickHouse type，emr-clickhouse : emr;cdw-clickhouse : cdwch;selfBuilt : .",
															},
															"drop_cls": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"drop_invalid_message_to_cls": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether to deliver to cls.",
																		},
																		"drop_cls_region": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls region.",
																		},
																		"drop_cls_owneruin": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls account.",
																		},
																		"drop_cls_topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls topicId.",
																		},
																		"drop_cls_log_set": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls LogSet id.",
																		},
																	},
																},
															},
														},
													},
												},
												"cls_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Cls configuration，Required when Type is CLS.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"decode_json": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the produced information is in json format.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cls id.",
															},
															"log_set": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "LogSet id.",
															},
															"content_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Required when Decode Json is false.",
															},
															"time_field": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specify the content of a field in the message as the time of the cls log. The format of the field content needs to be a second-level timestamp.",
															},
														},
													},
												},
												"cos_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Cos configuration, required when Type is COS.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cos bucket name.",
															},
															"region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Region code.",
															},
															"object_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ObjectKey.",
															},
															"aggregate_batch_size": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The size of aggregated messages MB.",
															},
															"aggregate_interval": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Time interval.",
															},
															"format_output_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The file format after message aggregation csv|json.",
															},
															"object_key_prefix": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dumped object directory prefix.",
															},
															"directory_time_format": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Partition format formatted according to strptime time.",
															},
														},
													},
												},
												"my_s_q_l_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "MySQL configuration，Required when Type is MYSQL.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MySQL database name，* is the whole database.",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name of the MySQL data table,  is the non-system table in all the monitored databases, which can be separated by , to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name , when a regular expression needs to be filled in, the format is data database name.data table name.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MySQL connection Id.",
															},
															"snapshot_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to Copy inventory information (schema_only does not copy, initial full amount), the default is initial.",
															},
															"ddl_topic": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The Topic that stores the Ddl information of My SQL, if it is empty, it will not be stored by default.",
															},
															"data_source_monitor_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "TABLE indicates that the read item is a table, QUERY indicates that the read item is a query.",
															},
															"data_source_monitor_resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "When DataMonitorMode=TABLE, pass in the Table that needs to be read; when DataMonitorMode=QUERY, pass in the query sql statement that needs to be read.",
															},
															"data_source_increment_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "TIMESTAMP indicates that the incremental column is of timestamp type, INCREMENT indicates that the incremental column is of self-incrementing id type.",
															},
															"data_source_increment_column": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name of the column to be monitored.",
															},
															"data_source_start_from": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "HEAD means copy stock + incremental data, TAIL means copy only incremental data.",
															},
															"data_target_insert_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.",
															},
															"data_target_primary_key_field": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.",
															},
															"data_target_record_mapping": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Mapping relationship between tables and messages.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"json_key": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The key name of the message.",
																		},
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Message type.",
																		},
																		"allow_null": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether the message is allowed to be empty.",
																		},
																		"column_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Corresponding mapping column name.",
																		},
																		"extra_info": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Database table extra fields.",
																		},
																		"column_size": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Current column size.",
																		},
																		"decimal_digits": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Current column precision.",
																		},
																		"auto_increment": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether it is an auto-increment column.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Database table default parameters.",
																		},
																	},
																},
															},
															"topic_regex": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regular expression for routing events to specific topics, defaults to (.*).",
															},
															"topic_replacement": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "TopicRegex ，$1、$2.",
															},
															"key_columns": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Format library1.table1 field 1,field 2;library 2.table2 field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
															},
															"drop_invalid_message": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to discard messages that fail to parse, the default is true.",
															},
															"drop_cls": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"drop_invalid_message_to_cls": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether to deliver to cls.",
																		},
																		"drop_cls_region": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The region where the cls is delivered.",
																		},
																		"drop_cls_owneruin": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Account.",
																		},
																		"drop_cls_topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls topic.",
																		},
																		"drop_cls_log_set": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cls LogSet id.",
																		},
																	},
																},
															},
															"output_format": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Output format，DEFAULT、CANAL_1、CANAL_2.",
															},
															"is_table_prefix": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "When the Table input is a prefix, the value of this item is true, otherwise it is false.",
															},
															"include_content_changes": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.",
															},
															"include_query": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.",
															},
															"record_with_schema": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
															},
															"signal_database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Database name of signal table.",
															},
															"is_table_regular": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the input table is a regular expression, if this option and Is Table Prefix are true at the same time, the judgment priority of this option is higher than Is Table Prefix.",
															},
														},
													},
												},
												"postgre_s_q_l_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "PostgreSQL configuration，Required when Type is POSTGRESQL or TDSQL C_POSTGRESQL.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "PostgreSQL database name.",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "PostgreSQL tableName，* is the non-system table in all the monitored databases, you can use , to monitor multiple data tables, but the data table needs to be filled in the format of Schema name.Data table name, and you need to fill in a regular expression When, the format is Schema name.data table name.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "PostgreSQL connection Id.",
															},
															"plugin_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "(decoderbufs/pgoutput)，default decoderbufs.",
															},
															"snapshot_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Never|initial, default initial.",
															},
															"data_format": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Upstream data format (JSON|Debezium), required when the database synchronization mode matches the default field.",
															},
															"data_target_insert_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.",
															},
															"data_target_primary_key_field": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.",
															},
															"data_target_record_mapping": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Mapping relationship between tables and messages.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"json_key": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The key name of the message.",
																		},
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Message type.",
																		},
																		"allow_null": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether the message is allowed to be empty.",
																		},
																		"column_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Column Name.",
																		},
																		"extra_info": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Database table extra fields.",
																		},
																		"column_size": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Current ColumnSize.",
																		},
																		"decimal_digits": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Current Column DecimalDigits.",
																		},
																		"auto_increment": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether it is an auto-increment column.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Database table default parameters.",
																		},
																	},
																},
															},
															"drop_invalid_message": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to discard messages that fail to parse, the default is true.",
															},
															"is_table_regular": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the input table is a regular expression.",
															},
															"key_columns": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Format  library1.table1:field 1,field2;library2.table2:field2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
															},
															"record_with_schema": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
															},
														},
													},
												},
												"topic_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Topic configuration，Required when Type is Topic.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The topic name of the topic sold separately.",
															},
															"offset_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
															},
															"start_time": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Topic TopicId.",
															},
															"compression_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.",
															},
															"use_auto_create_topic": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).",
															},
															"msg_multiple": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
															},
														},
													},
												},
												"maria_d_b_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "MariaDB configuration，Required when Type is MARIADB.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MariaDB database name，* for all database.",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MariaDB db name，*is the non-system table in all the monitored databases, you can use , to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "MariaDB connection Id.",
															},
															"snapshot_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Schema_only|initial，default initial.",
															},
															"key_columns": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Format  library 1. table 1: field 1, field 2; library 2. table 2: field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
															},
															"is_table_prefix": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "When the Table input is a prefix, the value of this item is true, otherwise it is false.",
															},
															"output_format": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Output format，DEFAULT、CANAL_1、CANAL_2.",
															},
															"include_content_changes": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.",
															},
															"include_query": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.",
															},
															"record_with_schema": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
															},
														},
													},
												},
												"s_q_l_server_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "SQLServer configuration，Required when Type is SQLSERVER.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SQLServer database name.",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SQLServer table，*is the non-system table in all the monitored databases, you can use , to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.",
															},
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SQLServer connection Id.",
															},
															"snapshot_mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Schema_only|initial default initial.",
															},
														},
													},
												},
												"ctsdb_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Ctsdb configuration，Required when Type is CTSDB.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource id.",
															},
															"ctsdb_metric": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Ctsdb的metric.",
															},
														},
													},
												},
												"scf_param": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Scf configuration，Required when Type is SCF.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"function_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF function name.",
															},
															"namespace": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF cloud function namespace, the default is default.",
															},
															"qualifier": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "SCF cloud function version and alias, the default is DEFAULT.",
															},
															"batch_size": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The maximum number of messages sent in each batch, the default is 1000.",
															},
															"max_retries": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The number of retries after the SCF call fails, the default is 5.",
															},
														},
													},
												},
											},
										},
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CreateTime.",
									},
									"error_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ErrorMessage.",
									},
									"task_progress": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Creation progress percentage.",
									},
									"task_current_step": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Task Current Step.",
									},
									"datahub_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datahub Id.",
									},
									"step_list": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "StepList.",
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

func dataSourceTencentCloudCkafkaDatahubTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_datahub_task.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		paramMap["TargetType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_type"); ok {
		paramMap["TaskType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		paramMap["SourceType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource"); ok {
		paramMap["Resource"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.DescribeDatahubTasksRes

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaDatahubTaskByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		describeDatahubTasksResMap := map[string]interface{}{}

		if result.TotalCount != nil {
			describeDatahubTasksResMap["total_count"] = result.TotalCount
		}

		if result.TaskList != nil {
			taskListList := []interface{}{}
			for _, taskList := range result.TaskList {
				taskListMap := map[string]interface{}{}

				if taskList.TaskId != nil {
					taskListMap["task_id"] = taskList.TaskId
				}

				if taskList.TaskName != nil {
					taskListMap["task_name"] = taskList.TaskName
				}

				if taskList.TaskType != nil {
					taskListMap["task_type"] = taskList.TaskType
				}

				if taskList.Status != nil {
					taskListMap["status"] = taskList.Status
				}

				if taskList.SourceResource != nil {
					sourceResourceMap := map[string]interface{}{}

					if taskList.SourceResource.Type != nil {
						sourceResourceMap["type"] = taskList.SourceResource.Type
					}

					if taskList.SourceResource.KafkaParam != nil {
						kafkaParamMap := map[string]interface{}{}

						if taskList.SourceResource.KafkaParam.SelfBuilt != nil {
							kafkaParamMap["self_built"] = taskList.SourceResource.KafkaParam.SelfBuilt
						}

						if taskList.SourceResource.KafkaParam.Resource != nil {
							kafkaParamMap["resource"] = taskList.SourceResource.KafkaParam.Resource
						}

						if taskList.SourceResource.KafkaParam.Topic != nil {
							kafkaParamMap["topic"] = taskList.SourceResource.KafkaParam.Topic
						}

						if taskList.SourceResource.KafkaParam.OffsetType != nil {
							kafkaParamMap["offset_type"] = taskList.SourceResource.KafkaParam.OffsetType
						}

						if taskList.SourceResource.KafkaParam.StartTime != nil {
							kafkaParamMap["start_time"] = taskList.SourceResource.KafkaParam.StartTime
						}

						if taskList.SourceResource.KafkaParam.ResourceName != nil {
							kafkaParamMap["resource_name"] = taskList.SourceResource.KafkaParam.ResourceName
						}

						if taskList.SourceResource.KafkaParam.ZoneId != nil {
							kafkaParamMap["zone_id"] = taskList.SourceResource.KafkaParam.ZoneId
						}

						if taskList.SourceResource.KafkaParam.TopicId != nil {
							kafkaParamMap["topic_id"] = taskList.SourceResource.KafkaParam.TopicId
						}

						if taskList.SourceResource.KafkaParam.PartitionNum != nil {
							kafkaParamMap["partition_num"] = taskList.SourceResource.KafkaParam.PartitionNum
						}

						if taskList.SourceResource.KafkaParam.EnableToleration != nil {
							kafkaParamMap["enable_toleration"] = taskList.SourceResource.KafkaParam.EnableToleration
						}

						if taskList.SourceResource.KafkaParam.QpsLimit != nil {
							kafkaParamMap["qps_limit"] = taskList.SourceResource.KafkaParam.QpsLimit
						}

						if taskList.SourceResource.KafkaParam.TableMappings != nil {
							tableMappingsList := []interface{}{}
							for _, tableMappings := range taskList.SourceResource.KafkaParam.TableMappings {
								tableMappingsMap := map[string]interface{}{}

								if tableMappings.Database != nil {
									tableMappingsMap["database"] = tableMappings.Database
								}

								if tableMappings.Table != nil {
									tableMappingsMap["table"] = tableMappings.Table
								}

								if tableMappings.Topic != nil {
									tableMappingsMap["topic"] = tableMappings.Topic
								}

								if tableMappings.TopicId != nil {
									tableMappingsMap["topic_id"] = tableMappings.TopicId
								}

								tableMappingsList = append(tableMappingsList, tableMappingsMap)
							}

							kafkaParamMap["table_mappings"] = []interface{}{tableMappingsList}
						}

						if taskList.SourceResource.KafkaParam.UseTableMapping != nil {
							kafkaParamMap["use_table_mapping"] = taskList.SourceResource.KafkaParam.UseTableMapping
						}

						if taskList.SourceResource.KafkaParam.UseAutoCreateTopic != nil {
							kafkaParamMap["use_auto_create_topic"] = taskList.SourceResource.KafkaParam.UseAutoCreateTopic
						}

						if taskList.SourceResource.KafkaParam.CompressionType != nil {
							kafkaParamMap["compression_type"] = taskList.SourceResource.KafkaParam.CompressionType
						}

						if taskList.SourceResource.KafkaParam.MsgMultiple != nil {
							kafkaParamMap["msg_multiple"] = taskList.SourceResource.KafkaParam.MsgMultiple
						}

						if taskList.SourceResource.KafkaParam.ConnectorSyncType != nil {
							kafkaParamMap["connector_sync_type"] = taskList.SourceResource.KafkaParam.ConnectorSyncType
						}

						if taskList.SourceResource.KafkaParam.KeepPartition != nil {
							kafkaParamMap["keep_partition"] = taskList.SourceResource.KafkaParam.KeepPartition
						}

						sourceResourceMap["kafka_param"] = []interface{}{kafkaParamMap}
					}

					if taskList.SourceResource.EventBusParam != nil {
						eventBusParamMap := map[string]interface{}{}

						if taskList.SourceResource.EventBusParam.Type != nil {
							eventBusParamMap["type"] = taskList.SourceResource.EventBusParam.Type
						}

						if taskList.SourceResource.EventBusParam.SelfBuilt != nil {
							eventBusParamMap["self_built"] = taskList.SourceResource.EventBusParam.SelfBuilt
						}

						if taskList.SourceResource.EventBusParam.Resource != nil {
							eventBusParamMap["resource"] = taskList.SourceResource.EventBusParam.Resource
						}

						if taskList.SourceResource.EventBusParam.Namespace != nil {
							eventBusParamMap["namespace"] = taskList.SourceResource.EventBusParam.Namespace
						}

						if taskList.SourceResource.EventBusParam.FunctionName != nil {
							eventBusParamMap["function_name"] = taskList.SourceResource.EventBusParam.FunctionName
						}

						if taskList.SourceResource.EventBusParam.Qualifier != nil {
							eventBusParamMap["qualifier"] = taskList.SourceResource.EventBusParam.Qualifier
						}

						sourceResourceMap["event_bus_param"] = []interface{}{eventBusParamMap}
					}

					if taskList.SourceResource.MongoDBParam != nil {
						mongoDBParamMap := map[string]interface{}{}

						if taskList.SourceResource.MongoDBParam.Database != nil {
							mongoDBParamMap["database"] = taskList.SourceResource.MongoDBParam.Database
						}

						if taskList.SourceResource.MongoDBParam.Collection != nil {
							mongoDBParamMap["collection"] = taskList.SourceResource.MongoDBParam.Collection
						}

						if taskList.SourceResource.MongoDBParam.CopyExisting != nil {
							mongoDBParamMap["copy_existing"] = taskList.SourceResource.MongoDBParam.CopyExisting
						}

						if taskList.SourceResource.MongoDBParam.Resource != nil {
							mongoDBParamMap["resource"] = taskList.SourceResource.MongoDBParam.Resource
						}

						if taskList.SourceResource.MongoDBParam.Ip != nil {
							mongoDBParamMap["ip"] = taskList.SourceResource.MongoDBParam.Ip
						}

						if taskList.SourceResource.MongoDBParam.Port != nil {
							mongoDBParamMap["port"] = taskList.SourceResource.MongoDBParam.Port
						}

						if taskList.SourceResource.MongoDBParam.UserName != nil {
							mongoDBParamMap["user_name"] = taskList.SourceResource.MongoDBParam.UserName
						}

						if taskList.SourceResource.MongoDBParam.Password != nil {
							mongoDBParamMap["password"] = taskList.SourceResource.MongoDBParam.Password
						}

						if taskList.SourceResource.MongoDBParam.ListeningEvent != nil {
							mongoDBParamMap["listening_event"] = taskList.SourceResource.MongoDBParam.ListeningEvent
						}

						if taskList.SourceResource.MongoDBParam.ReadPreference != nil {
							mongoDBParamMap["read_preference"] = taskList.SourceResource.MongoDBParam.ReadPreference
						}

						if taskList.SourceResource.MongoDBParam.Pipeline != nil {
							mongoDBParamMap["pipeline"] = taskList.SourceResource.MongoDBParam.Pipeline
						}

						if taskList.SourceResource.MongoDBParam.SelfBuilt != nil {
							mongoDBParamMap["self_built"] = taskList.SourceResource.MongoDBParam.SelfBuilt
						}

						sourceResourceMap["mongo_d_b_param"] = []interface{}{mongoDBParamMap}
					}

					if taskList.SourceResource.EsParam != nil {
						esParamMap := map[string]interface{}{}

						if taskList.SourceResource.EsParam.Resource != nil {
							esParamMap["resource"] = taskList.SourceResource.EsParam.Resource
						}

						if taskList.SourceResource.EsParam.Port != nil {
							esParamMap["port"] = taskList.SourceResource.EsParam.Port
						}

						if taskList.SourceResource.EsParam.UserName != nil {
							esParamMap["user_name"] = taskList.SourceResource.EsParam.UserName
						}

						if taskList.SourceResource.EsParam.Password != nil {
							esParamMap["password"] = taskList.SourceResource.EsParam.Password
						}

						if taskList.SourceResource.EsParam.SelfBuilt != nil {
							esParamMap["self_built"] = taskList.SourceResource.EsParam.SelfBuilt
						}

						if taskList.SourceResource.EsParam.ServiceVip != nil {
							esParamMap["service_vip"] = taskList.SourceResource.EsParam.ServiceVip
						}

						if taskList.SourceResource.EsParam.UniqVpcId != nil {
							esParamMap["uniq_vpc_id"] = taskList.SourceResource.EsParam.UniqVpcId
						}

						if taskList.SourceResource.EsParam.DropInvalidMessage != nil {
							esParamMap["drop_invalid_message"] = taskList.SourceResource.EsParam.DropInvalidMessage
						}

						if taskList.SourceResource.EsParam.Index != nil {
							esParamMap["index"] = taskList.SourceResource.EsParam.Index
						}

						if taskList.SourceResource.EsParam.DateFormat != nil {
							esParamMap["date_format"] = taskList.SourceResource.EsParam.DateFormat
						}

						if taskList.SourceResource.EsParam.ContentKey != nil {
							esParamMap["content_key"] = taskList.SourceResource.EsParam.ContentKey
						}

						if taskList.SourceResource.EsParam.DropInvalidJsonMessage != nil {
							esParamMap["drop_invalid_json_message"] = taskList.SourceResource.EsParam.DropInvalidJsonMessage
						}

						if taskList.SourceResource.EsParam.DocumentIdField != nil {
							esParamMap["document_id_field"] = taskList.SourceResource.EsParam.DocumentIdField
						}

						if taskList.SourceResource.EsParam.IndexType != nil {
							esParamMap["index_type"] = taskList.SourceResource.EsParam.IndexType
						}

						if taskList.SourceResource.EsParam.DropCls != nil {
							dropClsMap := map[string]interface{}{}

							if taskList.SourceResource.EsParam.DropCls.DropInvalidMessageToCls != nil {
								dropClsMap["drop_invalid_message_to_cls"] = taskList.SourceResource.EsParam.DropCls.DropInvalidMessageToCls
							}

							if taskList.SourceResource.EsParam.DropCls.DropClsRegion != nil {
								dropClsMap["drop_cls_region"] = taskList.SourceResource.EsParam.DropCls.DropClsRegion
							}

							if taskList.SourceResource.EsParam.DropCls.DropClsOwneruin != nil {
								dropClsMap["drop_cls_owneruin"] = taskList.SourceResource.EsParam.DropCls.DropClsOwneruin
							}

							if taskList.SourceResource.EsParam.DropCls.DropClsTopicId != nil {
								dropClsMap["drop_cls_topic_id"] = taskList.SourceResource.EsParam.DropCls.DropClsTopicId
							}

							if taskList.SourceResource.EsParam.DropCls.DropClsLogSet != nil {
								dropClsMap["drop_cls_log_set"] = taskList.SourceResource.EsParam.DropCls.DropClsLogSet
							}

							esParamMap["drop_cls"] = []interface{}{dropClsMap}
						}

						if taskList.SourceResource.EsParam.DatabasePrimaryKey != nil {
							esParamMap["database_primary_key"] = taskList.SourceResource.EsParam.DatabasePrimaryKey
						}

						if taskList.SourceResource.EsParam.DropDlq != nil {
							dropDlqMap := map[string]interface{}{}

							if taskList.SourceResource.EsParam.DropDlq.Type != nil {
								dropDlqMap["type"] = taskList.SourceResource.EsParam.DropDlq.Type
							}

							if taskList.SourceResource.EsParam.DropDlq.KafkaParam != nil {
								kafkaParamMap := map[string]interface{}{}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.SelfBuilt != nil {
									kafkaParamMap["self_built"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.SelfBuilt
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.Resource != nil {
									kafkaParamMap["resource"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.Resource
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.Topic != nil {
									kafkaParamMap["topic"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.Topic
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.OffsetType != nil {
									kafkaParamMap["offset_type"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.OffsetType
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.StartTime != nil {
									kafkaParamMap["start_time"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.StartTime
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.ResourceName != nil {
									kafkaParamMap["resource_name"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.ResourceName
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.ZoneId != nil {
									kafkaParamMap["zone_id"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.ZoneId
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.TopicId != nil {
									kafkaParamMap["topic_id"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.TopicId
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.PartitionNum != nil {
									kafkaParamMap["partition_num"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.PartitionNum
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.EnableToleration != nil {
									kafkaParamMap["enable_toleration"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.EnableToleration
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.QpsLimit != nil {
									kafkaParamMap["qps_limit"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.QpsLimit
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.TableMappings != nil {
									tableMappingsList := []interface{}{}
									for _, tableMappings := range taskList.SourceResource.EsParam.DropDlq.KafkaParam.TableMappings {
										tableMappingsMap := map[string]interface{}{}

										if tableMappings.Database != nil {
											tableMappingsMap["database"] = tableMappings.Database
										}

										if tableMappings.Table != nil {
											tableMappingsMap["table"] = tableMappings.Table
										}

										if tableMappings.Topic != nil {
											tableMappingsMap["topic"] = tableMappings.Topic
										}

										if tableMappings.TopicId != nil {
											tableMappingsMap["topic_id"] = tableMappings.TopicId
										}

										tableMappingsList = append(tableMappingsList, tableMappingsMap)
									}

									kafkaParamMap["table_mappings"] = []interface{}{tableMappingsList}
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.UseTableMapping != nil {
									kafkaParamMap["use_table_mapping"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.UseTableMapping
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic != nil {
									kafkaParamMap["use_auto_create_topic"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.CompressionType != nil {
									kafkaParamMap["compression_type"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.CompressionType
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.MsgMultiple != nil {
									kafkaParamMap["msg_multiple"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.MsgMultiple
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.ConnectorSyncType != nil {
									kafkaParamMap["connector_sync_type"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.ConnectorSyncType
								}

								if taskList.SourceResource.EsParam.DropDlq.KafkaParam.KeepPartition != nil {
									kafkaParamMap["keep_partition"] = taskList.SourceResource.EsParam.DropDlq.KafkaParam.KeepPartition
								}

								dropDlqMap["kafka_param"] = []interface{}{kafkaParamMap}
							}

							if taskList.SourceResource.EsParam.DropDlq.RetryInterval != nil {
								dropDlqMap["retry_interval"] = taskList.SourceResource.EsParam.DropDlq.RetryInterval
							}

							if taskList.SourceResource.EsParam.DropDlq.MaxRetryAttempts != nil {
								dropDlqMap["max_retry_attempts"] = taskList.SourceResource.EsParam.DropDlq.MaxRetryAttempts
							}

							if taskList.SourceResource.EsParam.DropDlq.TopicParam != nil {
								topicParamMap := map[string]interface{}{}

								if taskList.SourceResource.EsParam.DropDlq.TopicParam.Resource != nil {
									topicParamMap["resource"] = taskList.SourceResource.EsParam.DropDlq.TopicParam.Resource
								}

								if taskList.SourceResource.EsParam.DropDlq.TopicParam.OffsetType != nil {
									topicParamMap["offset_type"] = taskList.SourceResource.EsParam.DropDlq.TopicParam.OffsetType
								}

								if taskList.SourceResource.EsParam.DropDlq.TopicParam.StartTime != nil {
									topicParamMap["start_time"] = taskList.SourceResource.EsParam.DropDlq.TopicParam.StartTime
								}

								if taskList.SourceResource.EsParam.DropDlq.TopicParam.TopicId != nil {
									topicParamMap["topic_id"] = taskList.SourceResource.EsParam.DropDlq.TopicParam.TopicId
								}

								if taskList.SourceResource.EsParam.DropDlq.TopicParam.CompressionType != nil {
									topicParamMap["compression_type"] = taskList.SourceResource.EsParam.DropDlq.TopicParam.CompressionType
								}

								if taskList.SourceResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic != nil {
									topicParamMap["use_auto_create_topic"] = taskList.SourceResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic
								}

								if taskList.SourceResource.EsParam.DropDlq.TopicParam.MsgMultiple != nil {
									topicParamMap["msg_multiple"] = taskList.SourceResource.EsParam.DropDlq.TopicParam.MsgMultiple
								}

								dropDlqMap["topic_param"] = []interface{}{topicParamMap}
							}

							if taskList.SourceResource.EsParam.DropDlq.DlqType != nil {
								dropDlqMap["dlq_type"] = taskList.SourceResource.EsParam.DropDlq.DlqType
							}

							esParamMap["drop_dlq"] = []interface{}{dropDlqMap}
						}

						sourceResourceMap["es_param"] = []interface{}{esParamMap}
					}

					if taskList.SourceResource.TdwParam != nil {
						tdwParamMap := map[string]interface{}{}

						if taskList.SourceResource.TdwParam.Bid != nil {
							tdwParamMap["bid"] = taskList.SourceResource.TdwParam.Bid
						}

						if taskList.SourceResource.TdwParam.Tid != nil {
							tdwParamMap["tid"] = taskList.SourceResource.TdwParam.Tid
						}

						if taskList.SourceResource.TdwParam.IsDomestic != nil {
							tdwParamMap["is_domestic"] = taskList.SourceResource.TdwParam.IsDomestic
						}

						if taskList.SourceResource.TdwParam.TdwHost != nil {
							tdwParamMap["tdw_host"] = taskList.SourceResource.TdwParam.TdwHost
						}

						if taskList.SourceResource.TdwParam.TdwPort != nil {
							tdwParamMap["tdw_port"] = taskList.SourceResource.TdwParam.TdwPort
						}

						sourceResourceMap["tdw_param"] = []interface{}{tdwParamMap}
					}

					if taskList.SourceResource.DtsParam != nil {
						dtsParamMap := map[string]interface{}{}

						if taskList.SourceResource.DtsParam.Resource != nil {
							dtsParamMap["resource"] = taskList.SourceResource.DtsParam.Resource
						}

						if taskList.SourceResource.DtsParam.Ip != nil {
							dtsParamMap["ip"] = taskList.SourceResource.DtsParam.Ip
						}

						if taskList.SourceResource.DtsParam.Port != nil {
							dtsParamMap["port"] = taskList.SourceResource.DtsParam.Port
						}

						if taskList.SourceResource.DtsParam.Topic != nil {
							dtsParamMap["topic"] = taskList.SourceResource.DtsParam.Topic
						}

						if taskList.SourceResource.DtsParam.GroupId != nil {
							dtsParamMap["group_id"] = taskList.SourceResource.DtsParam.GroupId
						}

						if taskList.SourceResource.DtsParam.GroupUser != nil {
							dtsParamMap["group_user"] = taskList.SourceResource.DtsParam.GroupUser
						}

						if taskList.SourceResource.DtsParam.GroupPassword != nil {
							dtsParamMap["group_password"] = taskList.SourceResource.DtsParam.GroupPassword
						}

						if taskList.SourceResource.DtsParam.TranSql != nil {
							dtsParamMap["tran_sql"] = taskList.SourceResource.DtsParam.TranSql
						}

						sourceResourceMap["dts_param"] = []interface{}{dtsParamMap}
					}

					if taskList.SourceResource.ClickHouseParam != nil {
						clickHouseParamMap := map[string]interface{}{}

						if taskList.SourceResource.ClickHouseParam.Cluster != nil {
							clickHouseParamMap["cluster"] = taskList.SourceResource.ClickHouseParam.Cluster
						}

						if taskList.SourceResource.ClickHouseParam.Database != nil {
							clickHouseParamMap["database"] = taskList.SourceResource.ClickHouseParam.Database
						}

						if taskList.SourceResource.ClickHouseParam.Table != nil {
							clickHouseParamMap["table"] = taskList.SourceResource.ClickHouseParam.Table
						}

						if taskList.SourceResource.ClickHouseParam.Schema != nil {
							schemaList := []interface{}{}
							for _, schema := range taskList.SourceResource.ClickHouseParam.Schema {
								schemaMap := map[string]interface{}{}

								if schema.ColumnName != nil {
									schemaMap["column_name"] = schema.ColumnName
								}

								if schema.JsonKey != nil {
									schemaMap["json_key"] = schema.JsonKey
								}

								if schema.Type != nil {
									schemaMap["type"] = schema.Type
								}

								if schema.AllowNull != nil {
									schemaMap["allow_null"] = schema.AllowNull
								}

								schemaList = append(schemaList, schemaMap)
							}

							clickHouseParamMap["schema"] = []interface{}{schemaList}
						}

						if taskList.SourceResource.ClickHouseParam.Resource != nil {
							clickHouseParamMap["resource"] = taskList.SourceResource.ClickHouseParam.Resource
						}

						if taskList.SourceResource.ClickHouseParam.Ip != nil {
							clickHouseParamMap["ip"] = taskList.SourceResource.ClickHouseParam.Ip
						}

						if taskList.SourceResource.ClickHouseParam.Port != nil {
							clickHouseParamMap["port"] = taskList.SourceResource.ClickHouseParam.Port
						}

						if taskList.SourceResource.ClickHouseParam.UserName != nil {
							clickHouseParamMap["user_name"] = taskList.SourceResource.ClickHouseParam.UserName
						}

						if taskList.SourceResource.ClickHouseParam.Password != nil {
							clickHouseParamMap["password"] = taskList.SourceResource.ClickHouseParam.Password
						}

						if taskList.SourceResource.ClickHouseParam.ServiceVip != nil {
							clickHouseParamMap["service_vip"] = taskList.SourceResource.ClickHouseParam.ServiceVip
						}

						if taskList.SourceResource.ClickHouseParam.UniqVpcId != nil {
							clickHouseParamMap["uniq_vpc_id"] = taskList.SourceResource.ClickHouseParam.UniqVpcId
						}

						if taskList.SourceResource.ClickHouseParam.SelfBuilt != nil {
							clickHouseParamMap["self_built"] = taskList.SourceResource.ClickHouseParam.SelfBuilt
						}

						if taskList.SourceResource.ClickHouseParam.DropInvalidMessage != nil {
							clickHouseParamMap["drop_invalid_message"] = taskList.SourceResource.ClickHouseParam.DropInvalidMessage
						}

						if taskList.SourceResource.ClickHouseParam.Type != nil {
							clickHouseParamMap["type"] = taskList.SourceResource.ClickHouseParam.Type
						}

						if taskList.SourceResource.ClickHouseParam.DropCls != nil {
							dropClsMap := map[string]interface{}{}

							if taskList.SourceResource.ClickHouseParam.DropCls.DropInvalidMessageToCls != nil {
								dropClsMap["drop_invalid_message_to_cls"] = taskList.SourceResource.ClickHouseParam.DropCls.DropInvalidMessageToCls
							}

							if taskList.SourceResource.ClickHouseParam.DropCls.DropClsRegion != nil {
								dropClsMap["drop_cls_region"] = taskList.SourceResource.ClickHouseParam.DropCls.DropClsRegion
							}

							if taskList.SourceResource.ClickHouseParam.DropCls.DropClsOwneruin != nil {
								dropClsMap["drop_cls_owneruin"] = taskList.SourceResource.ClickHouseParam.DropCls.DropClsOwneruin
							}

							if taskList.SourceResource.ClickHouseParam.DropCls.DropClsTopicId != nil {
								dropClsMap["drop_cls_topic_id"] = taskList.SourceResource.ClickHouseParam.DropCls.DropClsTopicId
							}

							if taskList.SourceResource.ClickHouseParam.DropCls.DropClsLogSet != nil {
								dropClsMap["drop_cls_log_set"] = taskList.SourceResource.ClickHouseParam.DropCls.DropClsLogSet
							}

							clickHouseParamMap["drop_cls"] = []interface{}{dropClsMap}
						}

						sourceResourceMap["click_house_param"] = []interface{}{clickHouseParamMap}
					}

					if taskList.SourceResource.ClsParam != nil {
						clsParamMap := map[string]interface{}{}

						if taskList.SourceResource.ClsParam.DecodeJson != nil {
							clsParamMap["decode_json"] = taskList.SourceResource.ClsParam.DecodeJson
						}

						if taskList.SourceResource.ClsParam.Resource != nil {
							clsParamMap["resource"] = taskList.SourceResource.ClsParam.Resource
						}

						if taskList.SourceResource.ClsParam.LogSet != nil {
							clsParamMap["log_set"] = taskList.SourceResource.ClsParam.LogSet
						}

						if taskList.SourceResource.ClsParam.ContentKey != nil {
							clsParamMap["content_key"] = taskList.SourceResource.ClsParam.ContentKey
						}

						if taskList.SourceResource.ClsParam.TimeField != nil {
							clsParamMap["time_field"] = taskList.SourceResource.ClsParam.TimeField
						}

						sourceResourceMap["cls_param"] = []interface{}{clsParamMap}
					}

					if taskList.SourceResource.CosParam != nil {
						cosParamMap := map[string]interface{}{}

						if taskList.SourceResource.CosParam.BucketName != nil {
							cosParamMap["bucket_name"] = taskList.SourceResource.CosParam.BucketName
						}

						if taskList.SourceResource.CosParam.Region != nil {
							cosParamMap["region"] = taskList.SourceResource.CosParam.Region
						}

						if taskList.SourceResource.CosParam.ObjectKey != nil {
							cosParamMap["object_key"] = taskList.SourceResource.CosParam.ObjectKey
						}

						if taskList.SourceResource.CosParam.AggregateBatchSize != nil {
							cosParamMap["aggregate_batch_size"] = taskList.SourceResource.CosParam.AggregateBatchSize
						}

						if taskList.SourceResource.CosParam.AggregateInterval != nil {
							cosParamMap["aggregate_interval"] = taskList.SourceResource.CosParam.AggregateInterval
						}

						if taskList.SourceResource.CosParam.FormatOutputType != nil {
							cosParamMap["format_output_type"] = taskList.SourceResource.CosParam.FormatOutputType
						}

						if taskList.SourceResource.CosParam.ObjectKeyPrefix != nil {
							cosParamMap["object_key_prefix"] = taskList.SourceResource.CosParam.ObjectKeyPrefix
						}

						if taskList.SourceResource.CosParam.DirectoryTimeFormat != nil {
							cosParamMap["directory_time_format"] = taskList.SourceResource.CosParam.DirectoryTimeFormat
						}

						sourceResourceMap["cos_param"] = []interface{}{cosParamMap}
					}

					if taskList.SourceResource.MySQLParam != nil {
						mySQLParamMap := map[string]interface{}{}

						if taskList.SourceResource.MySQLParam.Database != nil {
							mySQLParamMap["database"] = taskList.SourceResource.MySQLParam.Database
						}

						if taskList.SourceResource.MySQLParam.Table != nil {
							mySQLParamMap["table"] = taskList.SourceResource.MySQLParam.Table
						}

						if taskList.SourceResource.MySQLParam.Resource != nil {
							mySQLParamMap["resource"] = taskList.SourceResource.MySQLParam.Resource
						}

						if taskList.SourceResource.MySQLParam.SnapshotMode != nil {
							mySQLParamMap["snapshot_mode"] = taskList.SourceResource.MySQLParam.SnapshotMode
						}

						if taskList.SourceResource.MySQLParam.DdlTopic != nil {
							mySQLParamMap["ddl_topic"] = taskList.SourceResource.MySQLParam.DdlTopic
						}

						if taskList.SourceResource.MySQLParam.DataSourceMonitorMode != nil {
							mySQLParamMap["data_source_monitor_mode"] = taskList.SourceResource.MySQLParam.DataSourceMonitorMode
						}

						if taskList.SourceResource.MySQLParam.DataSourceMonitorResource != nil {
							mySQLParamMap["data_source_monitor_resource"] = taskList.SourceResource.MySQLParam.DataSourceMonitorResource
						}

						if taskList.SourceResource.MySQLParam.DataSourceIncrementMode != nil {
							mySQLParamMap["data_source_increment_mode"] = taskList.SourceResource.MySQLParam.DataSourceIncrementMode
						}

						if taskList.SourceResource.MySQLParam.DataSourceIncrementColumn != nil {
							mySQLParamMap["data_source_increment_column"] = taskList.SourceResource.MySQLParam.DataSourceIncrementColumn
						}

						if taskList.SourceResource.MySQLParam.DataSourceStartFrom != nil {
							mySQLParamMap["data_source_start_from"] = taskList.SourceResource.MySQLParam.DataSourceStartFrom
						}

						if taskList.SourceResource.MySQLParam.DataTargetInsertMode != nil {
							mySQLParamMap["data_target_insert_mode"] = taskList.SourceResource.MySQLParam.DataTargetInsertMode
						}

						if taskList.SourceResource.MySQLParam.DataTargetPrimaryKeyField != nil {
							mySQLParamMap["data_target_primary_key_field"] = taskList.SourceResource.MySQLParam.DataTargetPrimaryKeyField
						}

						if taskList.SourceResource.MySQLParam.DataTargetRecordMapping != nil {
							dataTargetRecordMappingList := []interface{}{}
							for _, dataTargetRecordMapping := range taskList.SourceResource.MySQLParam.DataTargetRecordMapping {
								dataTargetRecordMappingMap := map[string]interface{}{}

								if dataTargetRecordMapping.JsonKey != nil {
									dataTargetRecordMappingMap["json_key"] = dataTargetRecordMapping.JsonKey
								}

								if dataTargetRecordMapping.Type != nil {
									dataTargetRecordMappingMap["type"] = dataTargetRecordMapping.Type
								}

								if dataTargetRecordMapping.AllowNull != nil {
									dataTargetRecordMappingMap["allow_null"] = dataTargetRecordMapping.AllowNull
								}

								if dataTargetRecordMapping.ColumnName != nil {
									dataTargetRecordMappingMap["column_name"] = dataTargetRecordMapping.ColumnName
								}

								if dataTargetRecordMapping.ExtraInfo != nil {
									dataTargetRecordMappingMap["extra_info"] = dataTargetRecordMapping.ExtraInfo
								}

								if dataTargetRecordMapping.ColumnSize != nil {
									dataTargetRecordMappingMap["column_size"] = dataTargetRecordMapping.ColumnSize
								}

								if dataTargetRecordMapping.DecimalDigits != nil {
									dataTargetRecordMappingMap["decimal_digits"] = dataTargetRecordMapping.DecimalDigits
								}

								if dataTargetRecordMapping.AutoIncrement != nil {
									dataTargetRecordMappingMap["auto_increment"] = dataTargetRecordMapping.AutoIncrement
								}

								if dataTargetRecordMapping.DefaultValue != nil {
									dataTargetRecordMappingMap["default_value"] = dataTargetRecordMapping.DefaultValue
								}

								dataTargetRecordMappingList = append(dataTargetRecordMappingList, dataTargetRecordMappingMap)
							}

							mySQLParamMap["data_target_record_mapping"] = []interface{}{dataTargetRecordMappingList}
						}

						if taskList.SourceResource.MySQLParam.TopicRegex != nil {
							mySQLParamMap["topic_regex"] = taskList.SourceResource.MySQLParam.TopicRegex
						}

						if taskList.SourceResource.MySQLParam.TopicReplacement != nil {
							mySQLParamMap["topic_replacement"] = taskList.SourceResource.MySQLParam.TopicReplacement
						}

						if taskList.SourceResource.MySQLParam.KeyColumns != nil {
							mySQLParamMap["key_columns"] = taskList.SourceResource.MySQLParam.KeyColumns
						}

						if taskList.SourceResource.MySQLParam.DropInvalidMessage != nil {
							mySQLParamMap["drop_invalid_message"] = taskList.SourceResource.MySQLParam.DropInvalidMessage
						}

						if taskList.SourceResource.MySQLParam.DropCls != nil {
							dropClsMap := map[string]interface{}{}

							if taskList.SourceResource.MySQLParam.DropCls.DropInvalidMessageToCls != nil {
								dropClsMap["drop_invalid_message_to_cls"] = taskList.SourceResource.MySQLParam.DropCls.DropInvalidMessageToCls
							}

							if taskList.SourceResource.MySQLParam.DropCls.DropClsRegion != nil {
								dropClsMap["drop_cls_region"] = taskList.SourceResource.MySQLParam.DropCls.DropClsRegion
							}

							if taskList.SourceResource.MySQLParam.DropCls.DropClsOwneruin != nil {
								dropClsMap["drop_cls_owneruin"] = taskList.SourceResource.MySQLParam.DropCls.DropClsOwneruin
							}

							if taskList.SourceResource.MySQLParam.DropCls.DropClsTopicId != nil {
								dropClsMap["drop_cls_topic_id"] = taskList.SourceResource.MySQLParam.DropCls.DropClsTopicId
							}

							if taskList.SourceResource.MySQLParam.DropCls.DropClsLogSet != nil {
								dropClsMap["drop_cls_log_set"] = taskList.SourceResource.MySQLParam.DropCls.DropClsLogSet
							}

							mySQLParamMap["drop_cls"] = []interface{}{dropClsMap}
						}

						if taskList.SourceResource.MySQLParam.OutputFormat != nil {
							mySQLParamMap["output_format"] = taskList.SourceResource.MySQLParam.OutputFormat
						}

						if taskList.SourceResource.MySQLParam.IsTablePrefix != nil {
							mySQLParamMap["is_table_prefix"] = taskList.SourceResource.MySQLParam.IsTablePrefix
						}

						if taskList.SourceResource.MySQLParam.IncludeContentChanges != nil {
							mySQLParamMap["include_content_changes"] = taskList.SourceResource.MySQLParam.IncludeContentChanges
						}

						if taskList.SourceResource.MySQLParam.IncludeQuery != nil {
							mySQLParamMap["include_query"] = taskList.SourceResource.MySQLParam.IncludeQuery
						}

						if taskList.SourceResource.MySQLParam.RecordWithSchema != nil {
							mySQLParamMap["record_with_schema"] = taskList.SourceResource.MySQLParam.RecordWithSchema
						}

						if taskList.SourceResource.MySQLParam.SignalDatabase != nil {
							mySQLParamMap["signal_database"] = taskList.SourceResource.MySQLParam.SignalDatabase
						}

						if taskList.SourceResource.MySQLParam.IsTableRegular != nil {
							mySQLParamMap["is_table_regular"] = taskList.SourceResource.MySQLParam.IsTableRegular
						}

						sourceResourceMap["my_s_q_l_param"] = []interface{}{mySQLParamMap}
					}

					if taskList.SourceResource.PostgreSQLParam != nil {
						postgreSQLParamMap := map[string]interface{}{}

						if taskList.SourceResource.PostgreSQLParam.Database != nil {
							postgreSQLParamMap["database"] = taskList.SourceResource.PostgreSQLParam.Database
						}

						if taskList.SourceResource.PostgreSQLParam.Table != nil {
							postgreSQLParamMap["table"] = taskList.SourceResource.PostgreSQLParam.Table
						}

						if taskList.SourceResource.PostgreSQLParam.Resource != nil {
							postgreSQLParamMap["resource"] = taskList.SourceResource.PostgreSQLParam.Resource
						}

						if taskList.SourceResource.PostgreSQLParam.PluginName != nil {
							postgreSQLParamMap["plugin_name"] = taskList.SourceResource.PostgreSQLParam.PluginName
						}

						if taskList.SourceResource.PostgreSQLParam.SnapshotMode != nil {
							postgreSQLParamMap["snapshot_mode"] = taskList.SourceResource.PostgreSQLParam.SnapshotMode
						}

						if taskList.SourceResource.PostgreSQLParam.DataFormat != nil {
							postgreSQLParamMap["data_format"] = taskList.SourceResource.PostgreSQLParam.DataFormat
						}

						if taskList.SourceResource.PostgreSQLParam.DataTargetInsertMode != nil {
							postgreSQLParamMap["data_target_insert_mode"] = taskList.SourceResource.PostgreSQLParam.DataTargetInsertMode
						}

						if taskList.SourceResource.PostgreSQLParam.DataTargetPrimaryKeyField != nil {
							postgreSQLParamMap["data_target_primary_key_field"] = taskList.SourceResource.PostgreSQLParam.DataTargetPrimaryKeyField
						}

						if taskList.SourceResource.PostgreSQLParam.DataTargetRecordMapping != nil {
							dataTargetRecordMappingList := []interface{}{}
							for _, dataTargetRecordMapping := range taskList.SourceResource.PostgreSQLParam.DataTargetRecordMapping {
								dataTargetRecordMappingMap := map[string]interface{}{}

								if dataTargetRecordMapping.JsonKey != nil {
									dataTargetRecordMappingMap["json_key"] = dataTargetRecordMapping.JsonKey
								}

								if dataTargetRecordMapping.Type != nil {
									dataTargetRecordMappingMap["type"] = dataTargetRecordMapping.Type
								}

								if dataTargetRecordMapping.AllowNull != nil {
									dataTargetRecordMappingMap["allow_null"] = dataTargetRecordMapping.AllowNull
								}

								if dataTargetRecordMapping.ColumnName != nil {
									dataTargetRecordMappingMap["column_name"] = dataTargetRecordMapping.ColumnName
								}

								if dataTargetRecordMapping.ExtraInfo != nil {
									dataTargetRecordMappingMap["extra_info"] = dataTargetRecordMapping.ExtraInfo
								}

								if dataTargetRecordMapping.ColumnSize != nil {
									dataTargetRecordMappingMap["column_size"] = dataTargetRecordMapping.ColumnSize
								}

								if dataTargetRecordMapping.DecimalDigits != nil {
									dataTargetRecordMappingMap["decimal_digits"] = dataTargetRecordMapping.DecimalDigits
								}

								if dataTargetRecordMapping.AutoIncrement != nil {
									dataTargetRecordMappingMap["auto_increment"] = dataTargetRecordMapping.AutoIncrement
								}

								if dataTargetRecordMapping.DefaultValue != nil {
									dataTargetRecordMappingMap["default_value"] = dataTargetRecordMapping.DefaultValue
								}

								dataTargetRecordMappingList = append(dataTargetRecordMappingList, dataTargetRecordMappingMap)
							}

							postgreSQLParamMap["data_target_record_mapping"] = []interface{}{dataTargetRecordMappingList}
						}

						if taskList.SourceResource.PostgreSQLParam.DropInvalidMessage != nil {
							postgreSQLParamMap["drop_invalid_message"] = taskList.SourceResource.PostgreSQLParam.DropInvalidMessage
						}

						if taskList.SourceResource.PostgreSQLParam.IsTableRegular != nil {
							postgreSQLParamMap["is_table_regular"] = taskList.SourceResource.PostgreSQLParam.IsTableRegular
						}

						if taskList.SourceResource.PostgreSQLParam.KeyColumns != nil {
							postgreSQLParamMap["key_columns"] = taskList.SourceResource.PostgreSQLParam.KeyColumns
						}

						if taskList.SourceResource.PostgreSQLParam.RecordWithSchema != nil {
							postgreSQLParamMap["record_with_schema"] = taskList.SourceResource.PostgreSQLParam.RecordWithSchema
						}

						sourceResourceMap["postgre_s_q_l_param"] = []interface{}{postgreSQLParamMap}
					}

					if taskList.SourceResource.TopicParam != nil {
						topicParamMap := map[string]interface{}{}

						if taskList.SourceResource.TopicParam.Resource != nil {
							topicParamMap["resource"] = taskList.SourceResource.TopicParam.Resource
						}

						if taskList.SourceResource.TopicParam.OffsetType != nil {
							topicParamMap["offset_type"] = taskList.SourceResource.TopicParam.OffsetType
						}

						if taskList.SourceResource.TopicParam.StartTime != nil {
							topicParamMap["start_time"] = taskList.SourceResource.TopicParam.StartTime
						}

						if taskList.SourceResource.TopicParam.TopicId != nil {
							topicParamMap["topic_id"] = taskList.SourceResource.TopicParam.TopicId
						}

						if taskList.SourceResource.TopicParam.CompressionType != nil {
							topicParamMap["compression_type"] = taskList.SourceResource.TopicParam.CompressionType
						}

						if taskList.SourceResource.TopicParam.UseAutoCreateTopic != nil {
							topicParamMap["use_auto_create_topic"] = taskList.SourceResource.TopicParam.UseAutoCreateTopic
						}

						if taskList.SourceResource.TopicParam.MsgMultiple != nil {
							topicParamMap["msg_multiple"] = taskList.SourceResource.TopicParam.MsgMultiple
						}

						sourceResourceMap["topic_param"] = []interface{}{topicParamMap}
					}

					if taskList.SourceResource.MariaDBParam != nil {
						mariaDBParamMap := map[string]interface{}{}

						if taskList.SourceResource.MariaDBParam.Database != nil {
							mariaDBParamMap["database"] = taskList.SourceResource.MariaDBParam.Database
						}

						if taskList.SourceResource.MariaDBParam.Table != nil {
							mariaDBParamMap["table"] = taskList.SourceResource.MariaDBParam.Table
						}

						if taskList.SourceResource.MariaDBParam.Resource != nil {
							mariaDBParamMap["resource"] = taskList.SourceResource.MariaDBParam.Resource
						}

						if taskList.SourceResource.MariaDBParam.SnapshotMode != nil {
							mariaDBParamMap["snapshot_mode"] = taskList.SourceResource.MariaDBParam.SnapshotMode
						}

						if taskList.SourceResource.MariaDBParam.KeyColumns != nil {
							mariaDBParamMap["key_columns"] = taskList.SourceResource.MariaDBParam.KeyColumns
						}

						if taskList.SourceResource.MariaDBParam.IsTablePrefix != nil {
							mariaDBParamMap["is_table_prefix"] = taskList.SourceResource.MariaDBParam.IsTablePrefix
						}

						if taskList.SourceResource.MariaDBParam.OutputFormat != nil {
							mariaDBParamMap["output_format"] = taskList.SourceResource.MariaDBParam.OutputFormat
						}

						if taskList.SourceResource.MariaDBParam.IncludeContentChanges != nil {
							mariaDBParamMap["include_content_changes"] = taskList.SourceResource.MariaDBParam.IncludeContentChanges
						}

						if taskList.SourceResource.MariaDBParam.IncludeQuery != nil {
							mariaDBParamMap["include_query"] = taskList.SourceResource.MariaDBParam.IncludeQuery
						}

						if taskList.SourceResource.MariaDBParam.RecordWithSchema != nil {
							mariaDBParamMap["record_with_schema"] = taskList.SourceResource.MariaDBParam.RecordWithSchema
						}

						sourceResourceMap["maria_d_b_param"] = []interface{}{mariaDBParamMap}
					}

					if taskList.SourceResource.SQLServerParam != nil {
						sQLServerParamMap := map[string]interface{}{}

						if taskList.SourceResource.SQLServerParam.Database != nil {
							sQLServerParamMap["database"] = taskList.SourceResource.SQLServerParam.Database
						}

						if taskList.SourceResource.SQLServerParam.Table != nil {
							sQLServerParamMap["table"] = taskList.SourceResource.SQLServerParam.Table
						}

						if taskList.SourceResource.SQLServerParam.Resource != nil {
							sQLServerParamMap["resource"] = taskList.SourceResource.SQLServerParam.Resource
						}

						if taskList.SourceResource.SQLServerParam.SnapshotMode != nil {
							sQLServerParamMap["snapshot_mode"] = taskList.SourceResource.SQLServerParam.SnapshotMode
						}

						sourceResourceMap["s_q_l_server_param"] = []interface{}{sQLServerParamMap}
					}

					if taskList.SourceResource.CtsdbParam != nil {
						ctsdbParamMap := map[string]interface{}{}

						if taskList.SourceResource.CtsdbParam.Resource != nil {
							ctsdbParamMap["resource"] = taskList.SourceResource.CtsdbParam.Resource
						}

						if taskList.SourceResource.CtsdbParam.CtsdbMetric != nil {
							ctsdbParamMap["ctsdb_metric"] = taskList.SourceResource.CtsdbParam.CtsdbMetric
						}

						sourceResourceMap["ctsdb_param"] = []interface{}{ctsdbParamMap}
					}

					if taskList.SourceResource.ScfParam != nil {
						scfParamMap := map[string]interface{}{}

						if taskList.SourceResource.ScfParam.FunctionName != nil {
							scfParamMap["function_name"] = taskList.SourceResource.ScfParam.FunctionName
						}

						if taskList.SourceResource.ScfParam.Namespace != nil {
							scfParamMap["namespace"] = taskList.SourceResource.ScfParam.Namespace
						}

						if taskList.SourceResource.ScfParam.Qualifier != nil {
							scfParamMap["qualifier"] = taskList.SourceResource.ScfParam.Qualifier
						}

						if taskList.SourceResource.ScfParam.BatchSize != nil {
							scfParamMap["batch_size"] = taskList.SourceResource.ScfParam.BatchSize
						}

						if taskList.SourceResource.ScfParam.MaxRetries != nil {
							scfParamMap["max_retries"] = taskList.SourceResource.ScfParam.MaxRetries
						}

						sourceResourceMap["scf_param"] = []interface{}{scfParamMap}
					}

					taskListMap["source_resource"] = []interface{}{sourceResourceMap}
				}

				if taskList.TargetResource != nil {
					targetResourceMap := map[string]interface{}{}

					if taskList.TargetResource.Type != nil {
						targetResourceMap["type"] = taskList.TargetResource.Type
					}

					if taskList.TargetResource.KafkaParam != nil {
						kafkaParamMap := map[string]interface{}{}

						if taskList.TargetResource.KafkaParam.SelfBuilt != nil {
							kafkaParamMap["self_built"] = taskList.TargetResource.KafkaParam.SelfBuilt
						}

						if taskList.TargetResource.KafkaParam.Resource != nil {
							kafkaParamMap["resource"] = taskList.TargetResource.KafkaParam.Resource
						}

						if taskList.TargetResource.KafkaParam.Topic != nil {
							kafkaParamMap["topic"] = taskList.TargetResource.KafkaParam.Topic
						}

						if taskList.TargetResource.KafkaParam.OffsetType != nil {
							kafkaParamMap["offset_type"] = taskList.TargetResource.KafkaParam.OffsetType
						}

						if taskList.TargetResource.KafkaParam.StartTime != nil {
							kafkaParamMap["start_time"] = taskList.TargetResource.KafkaParam.StartTime
						}

						if taskList.TargetResource.KafkaParam.ResourceName != nil {
							kafkaParamMap["resource_name"] = taskList.TargetResource.KafkaParam.ResourceName
						}

						if taskList.TargetResource.KafkaParam.ZoneId != nil {
							kafkaParamMap["zone_id"] = taskList.TargetResource.KafkaParam.ZoneId
						}

						if taskList.TargetResource.KafkaParam.TopicId != nil {
							kafkaParamMap["topic_id"] = taskList.TargetResource.KafkaParam.TopicId
						}

						if taskList.TargetResource.KafkaParam.PartitionNum != nil {
							kafkaParamMap["partition_num"] = taskList.TargetResource.KafkaParam.PartitionNum
						}

						if taskList.TargetResource.KafkaParam.EnableToleration != nil {
							kafkaParamMap["enable_toleration"] = taskList.TargetResource.KafkaParam.EnableToleration
						}

						if taskList.TargetResource.KafkaParam.QpsLimit != nil {
							kafkaParamMap["qps_limit"] = taskList.TargetResource.KafkaParam.QpsLimit
						}

						if taskList.TargetResource.KafkaParam.TableMappings != nil {
							tableMappingsList := []interface{}{}
							for _, tableMappings := range taskList.TargetResource.KafkaParam.TableMappings {
								tableMappingsMap := map[string]interface{}{}

								if tableMappings.Database != nil {
									tableMappingsMap["database"] = tableMappings.Database
								}

								if tableMappings.Table != nil {
									tableMappingsMap["table"] = tableMappings.Table
								}

								if tableMappings.Topic != nil {
									tableMappingsMap["topic"] = tableMappings.Topic
								}

								if tableMappings.TopicId != nil {
									tableMappingsMap["topic_id"] = tableMappings.TopicId
								}

								tableMappingsList = append(tableMappingsList, tableMappingsMap)
							}

							kafkaParamMap["table_mappings"] = []interface{}{tableMappingsList}
						}

						if taskList.TargetResource.KafkaParam.UseTableMapping != nil {
							kafkaParamMap["use_table_mapping"] = taskList.TargetResource.KafkaParam.UseTableMapping
						}

						if taskList.TargetResource.KafkaParam.UseAutoCreateTopic != nil {
							kafkaParamMap["use_auto_create_topic"] = taskList.TargetResource.KafkaParam.UseAutoCreateTopic
						}

						if taskList.TargetResource.KafkaParam.CompressionType != nil {
							kafkaParamMap["compression_type"] = taskList.TargetResource.KafkaParam.CompressionType
						}

						if taskList.TargetResource.KafkaParam.MsgMultiple != nil {
							kafkaParamMap["msg_multiple"] = taskList.TargetResource.KafkaParam.MsgMultiple
						}

						if taskList.TargetResource.KafkaParam.ConnectorSyncType != nil {
							kafkaParamMap["connector_sync_type"] = taskList.TargetResource.KafkaParam.ConnectorSyncType
						}

						if taskList.TargetResource.KafkaParam.KeepPartition != nil {
							kafkaParamMap["keep_partition"] = taskList.TargetResource.KafkaParam.KeepPartition
						}

						targetResourceMap["kafka_param"] = []interface{}{kafkaParamMap}
					}

					if taskList.TargetResource.EventBusParam != nil {
						eventBusParamMap := map[string]interface{}{}

						if taskList.TargetResource.EventBusParam.Type != nil {
							eventBusParamMap["type"] = taskList.TargetResource.EventBusParam.Type
						}

						if taskList.TargetResource.EventBusParam.SelfBuilt != nil {
							eventBusParamMap["self_built"] = taskList.TargetResource.EventBusParam.SelfBuilt
						}

						if taskList.TargetResource.EventBusParam.Resource != nil {
							eventBusParamMap["resource"] = taskList.TargetResource.EventBusParam.Resource
						}

						if taskList.TargetResource.EventBusParam.Namespace != nil {
							eventBusParamMap["namespace"] = taskList.TargetResource.EventBusParam.Namespace
						}

						if taskList.TargetResource.EventBusParam.FunctionName != nil {
							eventBusParamMap["function_name"] = taskList.TargetResource.EventBusParam.FunctionName
						}

						if taskList.TargetResource.EventBusParam.Qualifier != nil {
							eventBusParamMap["qualifier"] = taskList.TargetResource.EventBusParam.Qualifier
						}

						targetResourceMap["event_bus_param"] = []interface{}{eventBusParamMap}
					}

					if taskList.TargetResource.MongoDBParam != nil {
						mongoDBParamMap := map[string]interface{}{}

						if taskList.TargetResource.MongoDBParam.Database != nil {
							mongoDBParamMap["database"] = taskList.TargetResource.MongoDBParam.Database
						}

						if taskList.TargetResource.MongoDBParam.Collection != nil {
							mongoDBParamMap["collection"] = taskList.TargetResource.MongoDBParam.Collection
						}

						if taskList.TargetResource.MongoDBParam.CopyExisting != nil {
							mongoDBParamMap["copy_existing"] = taskList.TargetResource.MongoDBParam.CopyExisting
						}

						if taskList.TargetResource.MongoDBParam.Resource != nil {
							mongoDBParamMap["resource"] = taskList.TargetResource.MongoDBParam.Resource
						}

						if taskList.TargetResource.MongoDBParam.Ip != nil {
							mongoDBParamMap["ip"] = taskList.TargetResource.MongoDBParam.Ip
						}

						if taskList.TargetResource.MongoDBParam.Port != nil {
							mongoDBParamMap["port"] = taskList.TargetResource.MongoDBParam.Port
						}

						if taskList.TargetResource.MongoDBParam.UserName != nil {
							mongoDBParamMap["user_name"] = taskList.TargetResource.MongoDBParam.UserName
						}

						if taskList.TargetResource.MongoDBParam.Password != nil {
							mongoDBParamMap["password"] = taskList.TargetResource.MongoDBParam.Password
						}

						if taskList.TargetResource.MongoDBParam.ListeningEvent != nil {
							mongoDBParamMap["listening_event"] = taskList.TargetResource.MongoDBParam.ListeningEvent
						}

						if taskList.TargetResource.MongoDBParam.ReadPreference != nil {
							mongoDBParamMap["read_preference"] = taskList.TargetResource.MongoDBParam.ReadPreference
						}

						if taskList.TargetResource.MongoDBParam.Pipeline != nil {
							mongoDBParamMap["pipeline"] = taskList.TargetResource.MongoDBParam.Pipeline
						}

						if taskList.TargetResource.MongoDBParam.SelfBuilt != nil {
							mongoDBParamMap["self_built"] = taskList.TargetResource.MongoDBParam.SelfBuilt
						}

						targetResourceMap["mongo_d_b_param"] = []interface{}{mongoDBParamMap}
					}

					if taskList.TargetResource.EsParam != nil {
						esParamMap := map[string]interface{}{}

						if taskList.TargetResource.EsParam.Resource != nil {
							esParamMap["resource"] = taskList.TargetResource.EsParam.Resource
						}

						if taskList.TargetResource.EsParam.Port != nil {
							esParamMap["port"] = taskList.TargetResource.EsParam.Port
						}

						if taskList.TargetResource.EsParam.UserName != nil {
							esParamMap["user_name"] = taskList.TargetResource.EsParam.UserName
						}

						if taskList.TargetResource.EsParam.Password != nil {
							esParamMap["password"] = taskList.TargetResource.EsParam.Password
						}

						if taskList.TargetResource.EsParam.SelfBuilt != nil {
							esParamMap["self_built"] = taskList.TargetResource.EsParam.SelfBuilt
						}

						if taskList.TargetResource.EsParam.ServiceVip != nil {
							esParamMap["service_vip"] = taskList.TargetResource.EsParam.ServiceVip
						}

						if taskList.TargetResource.EsParam.UniqVpcId != nil {
							esParamMap["uniq_vpc_id"] = taskList.TargetResource.EsParam.UniqVpcId
						}

						if taskList.TargetResource.EsParam.DropInvalidMessage != nil {
							esParamMap["drop_invalid_message"] = taskList.TargetResource.EsParam.DropInvalidMessage
						}

						if taskList.TargetResource.EsParam.Index != nil {
							esParamMap["index"] = taskList.TargetResource.EsParam.Index
						}

						if taskList.TargetResource.EsParam.DateFormat != nil {
							esParamMap["date_format"] = taskList.TargetResource.EsParam.DateFormat
						}

						if taskList.TargetResource.EsParam.ContentKey != nil {
							esParamMap["content_key"] = taskList.TargetResource.EsParam.ContentKey
						}

						if taskList.TargetResource.EsParam.DropInvalidJsonMessage != nil {
							esParamMap["drop_invalid_json_message"] = taskList.TargetResource.EsParam.DropInvalidJsonMessage
						}

						if taskList.TargetResource.EsParam.DocumentIdField != nil {
							esParamMap["document_id_field"] = taskList.TargetResource.EsParam.DocumentIdField
						}

						if taskList.TargetResource.EsParam.IndexType != nil {
							esParamMap["index_type"] = taskList.TargetResource.EsParam.IndexType
						}

						if taskList.TargetResource.EsParam.DropCls != nil {
							dropClsMap := map[string]interface{}{}

							if taskList.TargetResource.EsParam.DropCls.DropInvalidMessageToCls != nil {
								dropClsMap["drop_invalid_message_to_cls"] = taskList.TargetResource.EsParam.DropCls.DropInvalidMessageToCls
							}

							if taskList.TargetResource.EsParam.DropCls.DropClsRegion != nil {
								dropClsMap["drop_cls_region"] = taskList.TargetResource.EsParam.DropCls.DropClsRegion
							}

							if taskList.TargetResource.EsParam.DropCls.DropClsOwneruin != nil {
								dropClsMap["drop_cls_owneruin"] = taskList.TargetResource.EsParam.DropCls.DropClsOwneruin
							}

							if taskList.TargetResource.EsParam.DropCls.DropClsTopicId != nil {
								dropClsMap["drop_cls_topic_id"] = taskList.TargetResource.EsParam.DropCls.DropClsTopicId
							}

							if taskList.TargetResource.EsParam.DropCls.DropClsLogSet != nil {
								dropClsMap["drop_cls_log_set"] = taskList.TargetResource.EsParam.DropCls.DropClsLogSet
							}

							esParamMap["drop_cls"] = []interface{}{dropClsMap}
						}

						if taskList.TargetResource.EsParam.DatabasePrimaryKey != nil {
							esParamMap["database_primary_key"] = taskList.TargetResource.EsParam.DatabasePrimaryKey
						}

						if taskList.TargetResource.EsParam.DropDlq != nil {
							dropDlqMap := map[string]interface{}{}

							if taskList.TargetResource.EsParam.DropDlq.Type != nil {
								dropDlqMap["type"] = taskList.TargetResource.EsParam.DropDlq.Type
							}

							if taskList.TargetResource.EsParam.DropDlq.KafkaParam != nil {
								kafkaParamMap := map[string]interface{}{}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.SelfBuilt != nil {
									kafkaParamMap["self_built"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.SelfBuilt
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.Resource != nil {
									kafkaParamMap["resource"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.Resource
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.Topic != nil {
									kafkaParamMap["topic"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.Topic
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.OffsetType != nil {
									kafkaParamMap["offset_type"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.OffsetType
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.StartTime != nil {
									kafkaParamMap["start_time"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.StartTime
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.ResourceName != nil {
									kafkaParamMap["resource_name"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.ResourceName
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.ZoneId != nil {
									kafkaParamMap["zone_id"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.ZoneId
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.TopicId != nil {
									kafkaParamMap["topic_id"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.TopicId
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.PartitionNum != nil {
									kafkaParamMap["partition_num"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.PartitionNum
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.EnableToleration != nil {
									kafkaParamMap["enable_toleration"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.EnableToleration
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.QpsLimit != nil {
									kafkaParamMap["qps_limit"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.QpsLimit
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.TableMappings != nil {
									tableMappingsList := []interface{}{}
									for _, tableMappings := range taskList.TargetResource.EsParam.DropDlq.KafkaParam.TableMappings {
										tableMappingsMap := map[string]interface{}{}

										if tableMappings.Database != nil {
											tableMappingsMap["database"] = tableMappings.Database
										}

										if tableMappings.Table != nil {
											tableMappingsMap["table"] = tableMappings.Table
										}

										if tableMappings.Topic != nil {
											tableMappingsMap["topic"] = tableMappings.Topic
										}

										if tableMappings.TopicId != nil {
											tableMappingsMap["topic_id"] = tableMappings.TopicId
										}

										tableMappingsList = append(tableMappingsList, tableMappingsMap)
									}

									kafkaParamMap["table_mappings"] = []interface{}{tableMappingsList}
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.UseTableMapping != nil {
									kafkaParamMap["use_table_mapping"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.UseTableMapping
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic != nil {
									kafkaParamMap["use_auto_create_topic"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.CompressionType != nil {
									kafkaParamMap["compression_type"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.CompressionType
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.MsgMultiple != nil {
									kafkaParamMap["msg_multiple"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.MsgMultiple
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.ConnectorSyncType != nil {
									kafkaParamMap["connector_sync_type"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.ConnectorSyncType
								}

								if taskList.TargetResource.EsParam.DropDlq.KafkaParam.KeepPartition != nil {
									kafkaParamMap["keep_partition"] = taskList.TargetResource.EsParam.DropDlq.KafkaParam.KeepPartition
								}

								dropDlqMap["kafka_param"] = []interface{}{kafkaParamMap}
							}

							if taskList.TargetResource.EsParam.DropDlq.RetryInterval != nil {
								dropDlqMap["retry_interval"] = taskList.TargetResource.EsParam.DropDlq.RetryInterval
							}

							if taskList.TargetResource.EsParam.DropDlq.MaxRetryAttempts != nil {
								dropDlqMap["max_retry_attempts"] = taskList.TargetResource.EsParam.DropDlq.MaxRetryAttempts
							}

							if taskList.TargetResource.EsParam.DropDlq.TopicParam != nil {
								topicParamMap := map[string]interface{}{}

								if taskList.TargetResource.EsParam.DropDlq.TopicParam.Resource != nil {
									topicParamMap["resource"] = taskList.TargetResource.EsParam.DropDlq.TopicParam.Resource
								}

								if taskList.TargetResource.EsParam.DropDlq.TopicParam.OffsetType != nil {
									topicParamMap["offset_type"] = taskList.TargetResource.EsParam.DropDlq.TopicParam.OffsetType
								}

								if taskList.TargetResource.EsParam.DropDlq.TopicParam.StartTime != nil {
									topicParamMap["start_time"] = taskList.TargetResource.EsParam.DropDlq.TopicParam.StartTime
								}

								if taskList.TargetResource.EsParam.DropDlq.TopicParam.TopicId != nil {
									topicParamMap["topic_id"] = taskList.TargetResource.EsParam.DropDlq.TopicParam.TopicId
								}

								if taskList.TargetResource.EsParam.DropDlq.TopicParam.CompressionType != nil {
									topicParamMap["compression_type"] = taskList.TargetResource.EsParam.DropDlq.TopicParam.CompressionType
								}

								if taskList.TargetResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic != nil {
									topicParamMap["use_auto_create_topic"] = taskList.TargetResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic
								}

								if taskList.TargetResource.EsParam.DropDlq.TopicParam.MsgMultiple != nil {
									topicParamMap["msg_multiple"] = taskList.TargetResource.EsParam.DropDlq.TopicParam.MsgMultiple
								}

								dropDlqMap["topic_param"] = []interface{}{topicParamMap}
							}

							if taskList.TargetResource.EsParam.DropDlq.DlqType != nil {
								dropDlqMap["dlq_type"] = taskList.TargetResource.EsParam.DropDlq.DlqType
							}

							esParamMap["drop_dlq"] = []interface{}{dropDlqMap}
						}

						targetResourceMap["es_param"] = []interface{}{esParamMap}
					}

					if taskList.TargetResource.TdwParam != nil {
						tdwParamMap := map[string]interface{}{}

						if taskList.TargetResource.TdwParam.Bid != nil {
							tdwParamMap["bid"] = taskList.TargetResource.TdwParam.Bid
						}

						if taskList.TargetResource.TdwParam.Tid != nil {
							tdwParamMap["tid"] = taskList.TargetResource.TdwParam.Tid
						}

						if taskList.TargetResource.TdwParam.IsDomestic != nil {
							tdwParamMap["is_domestic"] = taskList.TargetResource.TdwParam.IsDomestic
						}

						if taskList.TargetResource.TdwParam.TdwHost != nil {
							tdwParamMap["tdw_host"] = taskList.TargetResource.TdwParam.TdwHost
						}

						if taskList.TargetResource.TdwParam.TdwPort != nil {
							tdwParamMap["tdw_port"] = taskList.TargetResource.TdwParam.TdwPort
						}

						targetResourceMap["tdw_param"] = []interface{}{tdwParamMap}
					}

					if taskList.TargetResource.DtsParam != nil {
						dtsParamMap := map[string]interface{}{}

						if taskList.TargetResource.DtsParam.Resource != nil {
							dtsParamMap["resource"] = taskList.TargetResource.DtsParam.Resource
						}

						if taskList.TargetResource.DtsParam.Ip != nil {
							dtsParamMap["ip"] = taskList.TargetResource.DtsParam.Ip
						}

						if taskList.TargetResource.DtsParam.Port != nil {
							dtsParamMap["port"] = taskList.TargetResource.DtsParam.Port
						}

						if taskList.TargetResource.DtsParam.Topic != nil {
							dtsParamMap["topic"] = taskList.TargetResource.DtsParam.Topic
						}

						if taskList.TargetResource.DtsParam.GroupId != nil {
							dtsParamMap["group_id"] = taskList.TargetResource.DtsParam.GroupId
						}

						if taskList.TargetResource.DtsParam.GroupUser != nil {
							dtsParamMap["group_user"] = taskList.TargetResource.DtsParam.GroupUser
						}

						if taskList.TargetResource.DtsParam.GroupPassword != nil {
							dtsParamMap["group_password"] = taskList.TargetResource.DtsParam.GroupPassword
						}

						if taskList.TargetResource.DtsParam.TranSql != nil {
							dtsParamMap["tran_sql"] = taskList.TargetResource.DtsParam.TranSql
						}

						targetResourceMap["dts_param"] = []interface{}{dtsParamMap}
					}

					if taskList.TargetResource.ClickHouseParam != nil {
						clickHouseParamMap := map[string]interface{}{}

						if taskList.TargetResource.ClickHouseParam.Cluster != nil {
							clickHouseParamMap["cluster"] = taskList.TargetResource.ClickHouseParam.Cluster
						}

						if taskList.TargetResource.ClickHouseParam.Database != nil {
							clickHouseParamMap["database"] = taskList.TargetResource.ClickHouseParam.Database
						}

						if taskList.TargetResource.ClickHouseParam.Table != nil {
							clickHouseParamMap["table"] = taskList.TargetResource.ClickHouseParam.Table
						}

						if taskList.TargetResource.ClickHouseParam.Schema != nil {
							schemaList := []interface{}{}
							for _, schema := range taskList.TargetResource.ClickHouseParam.Schema {
								schemaMap := map[string]interface{}{}

								if schema.ColumnName != nil {
									schemaMap["column_name"] = schema.ColumnName
								}

								if schema.JsonKey != nil {
									schemaMap["json_key"] = schema.JsonKey
								}

								if schema.Type != nil {
									schemaMap["type"] = schema.Type
								}

								if schema.AllowNull != nil {
									schemaMap["allow_null"] = schema.AllowNull
								}

								schemaList = append(schemaList, schemaMap)
							}

							clickHouseParamMap["schema"] = []interface{}{schemaList}
						}

						if taskList.TargetResource.ClickHouseParam.Resource != nil {
							clickHouseParamMap["resource"] = taskList.TargetResource.ClickHouseParam.Resource
						}

						if taskList.TargetResource.ClickHouseParam.Ip != nil {
							clickHouseParamMap["ip"] = taskList.TargetResource.ClickHouseParam.Ip
						}

						if taskList.TargetResource.ClickHouseParam.Port != nil {
							clickHouseParamMap["port"] = taskList.TargetResource.ClickHouseParam.Port
						}

						if taskList.TargetResource.ClickHouseParam.UserName != nil {
							clickHouseParamMap["user_name"] = taskList.TargetResource.ClickHouseParam.UserName
						}

						if taskList.TargetResource.ClickHouseParam.Password != nil {
							clickHouseParamMap["password"] = taskList.TargetResource.ClickHouseParam.Password
						}

						if taskList.TargetResource.ClickHouseParam.ServiceVip != nil {
							clickHouseParamMap["service_vip"] = taskList.TargetResource.ClickHouseParam.ServiceVip
						}

						if taskList.TargetResource.ClickHouseParam.UniqVpcId != nil {
							clickHouseParamMap["uniq_vpc_id"] = taskList.TargetResource.ClickHouseParam.UniqVpcId
						}

						if taskList.TargetResource.ClickHouseParam.SelfBuilt != nil {
							clickHouseParamMap["self_built"] = taskList.TargetResource.ClickHouseParam.SelfBuilt
						}

						if taskList.TargetResource.ClickHouseParam.DropInvalidMessage != nil {
							clickHouseParamMap["drop_invalid_message"] = taskList.TargetResource.ClickHouseParam.DropInvalidMessage
						}

						if taskList.TargetResource.ClickHouseParam.Type != nil {
							clickHouseParamMap["type"] = taskList.TargetResource.ClickHouseParam.Type
						}

						if taskList.TargetResource.ClickHouseParam.DropCls != nil {
							dropClsMap := map[string]interface{}{}

							if taskList.TargetResource.ClickHouseParam.DropCls.DropInvalidMessageToCls != nil {
								dropClsMap["drop_invalid_message_to_cls"] = taskList.TargetResource.ClickHouseParam.DropCls.DropInvalidMessageToCls
							}

							if taskList.TargetResource.ClickHouseParam.DropCls.DropClsRegion != nil {
								dropClsMap["drop_cls_region"] = taskList.TargetResource.ClickHouseParam.DropCls.DropClsRegion
							}

							if taskList.TargetResource.ClickHouseParam.DropCls.DropClsOwneruin != nil {
								dropClsMap["drop_cls_owneruin"] = taskList.TargetResource.ClickHouseParam.DropCls.DropClsOwneruin
							}

							if taskList.TargetResource.ClickHouseParam.DropCls.DropClsTopicId != nil {
								dropClsMap["drop_cls_topic_id"] = taskList.TargetResource.ClickHouseParam.DropCls.DropClsTopicId
							}

							if taskList.TargetResource.ClickHouseParam.DropCls.DropClsLogSet != nil {
								dropClsMap["drop_cls_log_set"] = taskList.TargetResource.ClickHouseParam.DropCls.DropClsLogSet
							}

							clickHouseParamMap["drop_cls"] = []interface{}{dropClsMap}
						}

						targetResourceMap["click_house_param"] = []interface{}{clickHouseParamMap}
					}

					if taskList.TargetResource.ClsParam != nil {
						clsParamMap := map[string]interface{}{}

						if taskList.TargetResource.ClsParam.DecodeJson != nil {
							clsParamMap["decode_json"] = taskList.TargetResource.ClsParam.DecodeJson
						}

						if taskList.TargetResource.ClsParam.Resource != nil {
							clsParamMap["resource"] = taskList.TargetResource.ClsParam.Resource
						}

						if taskList.TargetResource.ClsParam.LogSet != nil {
							clsParamMap["log_set"] = taskList.TargetResource.ClsParam.LogSet
						}

						if taskList.TargetResource.ClsParam.ContentKey != nil {
							clsParamMap["content_key"] = taskList.TargetResource.ClsParam.ContentKey
						}

						if taskList.TargetResource.ClsParam.TimeField != nil {
							clsParamMap["time_field"] = taskList.TargetResource.ClsParam.TimeField
						}

						targetResourceMap["cls_param"] = []interface{}{clsParamMap}
					}

					if taskList.TargetResource.CosParam != nil {
						cosParamMap := map[string]interface{}{}

						if taskList.TargetResource.CosParam.BucketName != nil {
							cosParamMap["bucket_name"] = taskList.TargetResource.CosParam.BucketName
						}

						if taskList.TargetResource.CosParam.Region != nil {
							cosParamMap["region"] = taskList.TargetResource.CosParam.Region
						}

						if taskList.TargetResource.CosParam.ObjectKey != nil {
							cosParamMap["object_key"] = taskList.TargetResource.CosParam.ObjectKey
						}

						if taskList.TargetResource.CosParam.AggregateBatchSize != nil {
							cosParamMap["aggregate_batch_size"] = taskList.TargetResource.CosParam.AggregateBatchSize
						}

						if taskList.TargetResource.CosParam.AggregateInterval != nil {
							cosParamMap["aggregate_interval"] = taskList.TargetResource.CosParam.AggregateInterval
						}

						if taskList.TargetResource.CosParam.FormatOutputType != nil {
							cosParamMap["format_output_type"] = taskList.TargetResource.CosParam.FormatOutputType
						}

						if taskList.TargetResource.CosParam.ObjectKeyPrefix != nil {
							cosParamMap["object_key_prefix"] = taskList.TargetResource.CosParam.ObjectKeyPrefix
						}

						if taskList.TargetResource.CosParam.DirectoryTimeFormat != nil {
							cosParamMap["directory_time_format"] = taskList.TargetResource.CosParam.DirectoryTimeFormat
						}

						targetResourceMap["cos_param"] = []interface{}{cosParamMap}
					}

					if taskList.TargetResource.MySQLParam != nil {
						mySQLParamMap := map[string]interface{}{}

						if taskList.TargetResource.MySQLParam.Database != nil {
							mySQLParamMap["database"] = taskList.TargetResource.MySQLParam.Database
						}

						if taskList.TargetResource.MySQLParam.Table != nil {
							mySQLParamMap["table"] = taskList.TargetResource.MySQLParam.Table
						}

						if taskList.TargetResource.MySQLParam.Resource != nil {
							mySQLParamMap["resource"] = taskList.TargetResource.MySQLParam.Resource
						}

						if taskList.TargetResource.MySQLParam.SnapshotMode != nil {
							mySQLParamMap["snapshot_mode"] = taskList.TargetResource.MySQLParam.SnapshotMode
						}

						if taskList.TargetResource.MySQLParam.DdlTopic != nil {
							mySQLParamMap["ddl_topic"] = taskList.TargetResource.MySQLParam.DdlTopic
						}

						if taskList.TargetResource.MySQLParam.DataSourceMonitorMode != nil {
							mySQLParamMap["data_source_monitor_mode"] = taskList.TargetResource.MySQLParam.DataSourceMonitorMode
						}

						if taskList.TargetResource.MySQLParam.DataSourceMonitorResource != nil {
							mySQLParamMap["data_source_monitor_resource"] = taskList.TargetResource.MySQLParam.DataSourceMonitorResource
						}

						if taskList.TargetResource.MySQLParam.DataSourceIncrementMode != nil {
							mySQLParamMap["data_source_increment_mode"] = taskList.TargetResource.MySQLParam.DataSourceIncrementMode
						}

						if taskList.TargetResource.MySQLParam.DataSourceIncrementColumn != nil {
							mySQLParamMap["data_source_increment_column"] = taskList.TargetResource.MySQLParam.DataSourceIncrementColumn
						}

						if taskList.TargetResource.MySQLParam.DataSourceStartFrom != nil {
							mySQLParamMap["data_source_start_from"] = taskList.TargetResource.MySQLParam.DataSourceStartFrom
						}

						if taskList.TargetResource.MySQLParam.DataTargetInsertMode != nil {
							mySQLParamMap["data_target_insert_mode"] = taskList.TargetResource.MySQLParam.DataTargetInsertMode
						}

						if taskList.TargetResource.MySQLParam.DataTargetPrimaryKeyField != nil {
							mySQLParamMap["data_target_primary_key_field"] = taskList.TargetResource.MySQLParam.DataTargetPrimaryKeyField
						}

						if taskList.TargetResource.MySQLParam.DataTargetRecordMapping != nil {
							dataTargetRecordMappingList := []interface{}{}
							for _, dataTargetRecordMapping := range taskList.TargetResource.MySQLParam.DataTargetRecordMapping {
								dataTargetRecordMappingMap := map[string]interface{}{}

								if dataTargetRecordMapping.JsonKey != nil {
									dataTargetRecordMappingMap["json_key"] = dataTargetRecordMapping.JsonKey
								}

								if dataTargetRecordMapping.Type != nil {
									dataTargetRecordMappingMap["type"] = dataTargetRecordMapping.Type
								}

								if dataTargetRecordMapping.AllowNull != nil {
									dataTargetRecordMappingMap["allow_null"] = dataTargetRecordMapping.AllowNull
								}

								if dataTargetRecordMapping.ColumnName != nil {
									dataTargetRecordMappingMap["column_name"] = dataTargetRecordMapping.ColumnName
								}

								if dataTargetRecordMapping.ExtraInfo != nil {
									dataTargetRecordMappingMap["extra_info"] = dataTargetRecordMapping.ExtraInfo
								}

								if dataTargetRecordMapping.ColumnSize != nil {
									dataTargetRecordMappingMap["column_size"] = dataTargetRecordMapping.ColumnSize
								}

								if dataTargetRecordMapping.DecimalDigits != nil {
									dataTargetRecordMappingMap["decimal_digits"] = dataTargetRecordMapping.DecimalDigits
								}

								if dataTargetRecordMapping.AutoIncrement != nil {
									dataTargetRecordMappingMap["auto_increment"] = dataTargetRecordMapping.AutoIncrement
								}

								if dataTargetRecordMapping.DefaultValue != nil {
									dataTargetRecordMappingMap["default_value"] = dataTargetRecordMapping.DefaultValue
								}

								dataTargetRecordMappingList = append(dataTargetRecordMappingList, dataTargetRecordMappingMap)
							}

							mySQLParamMap["data_target_record_mapping"] = []interface{}{dataTargetRecordMappingList}
						}

						if taskList.TargetResource.MySQLParam.TopicRegex != nil {
							mySQLParamMap["topic_regex"] = taskList.TargetResource.MySQLParam.TopicRegex
						}

						if taskList.TargetResource.MySQLParam.TopicReplacement != nil {
							mySQLParamMap["topic_replacement"] = taskList.TargetResource.MySQLParam.TopicReplacement
						}

						if taskList.TargetResource.MySQLParam.KeyColumns != nil {
							mySQLParamMap["key_columns"] = taskList.TargetResource.MySQLParam.KeyColumns
						}

						if taskList.TargetResource.MySQLParam.DropInvalidMessage != nil {
							mySQLParamMap["drop_invalid_message"] = taskList.TargetResource.MySQLParam.DropInvalidMessage
						}

						if taskList.TargetResource.MySQLParam.DropCls != nil {
							dropClsMap := map[string]interface{}{}

							if taskList.TargetResource.MySQLParam.DropCls.DropInvalidMessageToCls != nil {
								dropClsMap["drop_invalid_message_to_cls"] = taskList.TargetResource.MySQLParam.DropCls.DropInvalidMessageToCls
							}

							if taskList.TargetResource.MySQLParam.DropCls.DropClsRegion != nil {
								dropClsMap["drop_cls_region"] = taskList.TargetResource.MySQLParam.DropCls.DropClsRegion
							}

							if taskList.TargetResource.MySQLParam.DropCls.DropClsOwneruin != nil {
								dropClsMap["drop_cls_owneruin"] = taskList.TargetResource.MySQLParam.DropCls.DropClsOwneruin
							}

							if taskList.TargetResource.MySQLParam.DropCls.DropClsTopicId != nil {
								dropClsMap["drop_cls_topic_id"] = taskList.TargetResource.MySQLParam.DropCls.DropClsTopicId
							}

							if taskList.TargetResource.MySQLParam.DropCls.DropClsLogSet != nil {
								dropClsMap["drop_cls_log_set"] = taskList.TargetResource.MySQLParam.DropCls.DropClsLogSet
							}

							mySQLParamMap["drop_cls"] = []interface{}{dropClsMap}
						}

						if taskList.TargetResource.MySQLParam.OutputFormat != nil {
							mySQLParamMap["output_format"] = taskList.TargetResource.MySQLParam.OutputFormat
						}

						if taskList.TargetResource.MySQLParam.IsTablePrefix != nil {
							mySQLParamMap["is_table_prefix"] = taskList.TargetResource.MySQLParam.IsTablePrefix
						}

						if taskList.TargetResource.MySQLParam.IncludeContentChanges != nil {
							mySQLParamMap["include_content_changes"] = taskList.TargetResource.MySQLParam.IncludeContentChanges
						}

						if taskList.TargetResource.MySQLParam.IncludeQuery != nil {
							mySQLParamMap["include_query"] = taskList.TargetResource.MySQLParam.IncludeQuery
						}

						if taskList.TargetResource.MySQLParam.RecordWithSchema != nil {
							mySQLParamMap["record_with_schema"] = taskList.TargetResource.MySQLParam.RecordWithSchema
						}

						if taskList.TargetResource.MySQLParam.SignalDatabase != nil {
							mySQLParamMap["signal_database"] = taskList.TargetResource.MySQLParam.SignalDatabase
						}

						if taskList.TargetResource.MySQLParam.IsTableRegular != nil {
							mySQLParamMap["is_table_regular"] = taskList.TargetResource.MySQLParam.IsTableRegular
						}

						targetResourceMap["my_s_q_l_param"] = []interface{}{mySQLParamMap}
					}

					if taskList.TargetResource.PostgreSQLParam != nil {
						postgreSQLParamMap := map[string]interface{}{}

						if taskList.TargetResource.PostgreSQLParam.Database != nil {
							postgreSQLParamMap["database"] = taskList.TargetResource.PostgreSQLParam.Database
						}

						if taskList.TargetResource.PostgreSQLParam.Table != nil {
							postgreSQLParamMap["table"] = taskList.TargetResource.PostgreSQLParam.Table
						}

						if taskList.TargetResource.PostgreSQLParam.Resource != nil {
							postgreSQLParamMap["resource"] = taskList.TargetResource.PostgreSQLParam.Resource
						}

						if taskList.TargetResource.PostgreSQLParam.PluginName != nil {
							postgreSQLParamMap["plugin_name"] = taskList.TargetResource.PostgreSQLParam.PluginName
						}

						if taskList.TargetResource.PostgreSQLParam.SnapshotMode != nil {
							postgreSQLParamMap["snapshot_mode"] = taskList.TargetResource.PostgreSQLParam.SnapshotMode
						}

						if taskList.TargetResource.PostgreSQLParam.DataFormat != nil {
							postgreSQLParamMap["data_format"] = taskList.TargetResource.PostgreSQLParam.DataFormat
						}

						if taskList.TargetResource.PostgreSQLParam.DataTargetInsertMode != nil {
							postgreSQLParamMap["data_target_insert_mode"] = taskList.TargetResource.PostgreSQLParam.DataTargetInsertMode
						}

						if taskList.TargetResource.PostgreSQLParam.DataTargetPrimaryKeyField != nil {
							postgreSQLParamMap["data_target_primary_key_field"] = taskList.TargetResource.PostgreSQLParam.DataTargetPrimaryKeyField
						}

						if taskList.TargetResource.PostgreSQLParam.DataTargetRecordMapping != nil {
							dataTargetRecordMappingList := []interface{}{}
							for _, dataTargetRecordMapping := range taskList.TargetResource.PostgreSQLParam.DataTargetRecordMapping {
								dataTargetRecordMappingMap := map[string]interface{}{}

								if dataTargetRecordMapping.JsonKey != nil {
									dataTargetRecordMappingMap["json_key"] = dataTargetRecordMapping.JsonKey
								}

								if dataTargetRecordMapping.Type != nil {
									dataTargetRecordMappingMap["type"] = dataTargetRecordMapping.Type
								}

								if dataTargetRecordMapping.AllowNull != nil {
									dataTargetRecordMappingMap["allow_null"] = dataTargetRecordMapping.AllowNull
								}

								if dataTargetRecordMapping.ColumnName != nil {
									dataTargetRecordMappingMap["column_name"] = dataTargetRecordMapping.ColumnName
								}

								if dataTargetRecordMapping.ExtraInfo != nil {
									dataTargetRecordMappingMap["extra_info"] = dataTargetRecordMapping.ExtraInfo
								}

								if dataTargetRecordMapping.ColumnSize != nil {
									dataTargetRecordMappingMap["column_size"] = dataTargetRecordMapping.ColumnSize
								}

								if dataTargetRecordMapping.DecimalDigits != nil {
									dataTargetRecordMappingMap["decimal_digits"] = dataTargetRecordMapping.DecimalDigits
								}

								if dataTargetRecordMapping.AutoIncrement != nil {
									dataTargetRecordMappingMap["auto_increment"] = dataTargetRecordMapping.AutoIncrement
								}

								if dataTargetRecordMapping.DefaultValue != nil {
									dataTargetRecordMappingMap["default_value"] = dataTargetRecordMapping.DefaultValue
								}

								dataTargetRecordMappingList = append(dataTargetRecordMappingList, dataTargetRecordMappingMap)
							}

							postgreSQLParamMap["data_target_record_mapping"] = []interface{}{dataTargetRecordMappingList}
						}

						if taskList.TargetResource.PostgreSQLParam.DropInvalidMessage != nil {
							postgreSQLParamMap["drop_invalid_message"] = taskList.TargetResource.PostgreSQLParam.DropInvalidMessage
						}

						if taskList.TargetResource.PostgreSQLParam.IsTableRegular != nil {
							postgreSQLParamMap["is_table_regular"] = taskList.TargetResource.PostgreSQLParam.IsTableRegular
						}

						if taskList.TargetResource.PostgreSQLParam.KeyColumns != nil {
							postgreSQLParamMap["key_columns"] = taskList.TargetResource.PostgreSQLParam.KeyColumns
						}

						if taskList.TargetResource.PostgreSQLParam.RecordWithSchema != nil {
							postgreSQLParamMap["record_with_schema"] = taskList.TargetResource.PostgreSQLParam.RecordWithSchema
						}

						targetResourceMap["postgre_s_q_l_param"] = []interface{}{postgreSQLParamMap}
					}

					if taskList.TargetResource.TopicParam != nil {
						topicParamMap := map[string]interface{}{}

						if taskList.TargetResource.TopicParam.Resource != nil {
							topicParamMap["resource"] = taskList.TargetResource.TopicParam.Resource
						}

						if taskList.TargetResource.TopicParam.OffsetType != nil {
							topicParamMap["offset_type"] = taskList.TargetResource.TopicParam.OffsetType
						}

						if taskList.TargetResource.TopicParam.StartTime != nil {
							topicParamMap["start_time"] = taskList.TargetResource.TopicParam.StartTime
						}

						if taskList.TargetResource.TopicParam.TopicId != nil {
							topicParamMap["topic_id"] = taskList.TargetResource.TopicParam.TopicId
						}

						if taskList.TargetResource.TopicParam.CompressionType != nil {
							topicParamMap["compression_type"] = taskList.TargetResource.TopicParam.CompressionType
						}

						if taskList.TargetResource.TopicParam.UseAutoCreateTopic != nil {
							topicParamMap["use_auto_create_topic"] = taskList.TargetResource.TopicParam.UseAutoCreateTopic
						}

						if taskList.TargetResource.TopicParam.MsgMultiple != nil {
							topicParamMap["msg_multiple"] = taskList.TargetResource.TopicParam.MsgMultiple
						}

						targetResourceMap["topic_param"] = []interface{}{topicParamMap}
					}

					if taskList.TargetResource.MariaDBParam != nil {
						mariaDBParamMap := map[string]interface{}{}

						if taskList.TargetResource.MariaDBParam.Database != nil {
							mariaDBParamMap["database"] = taskList.TargetResource.MariaDBParam.Database
						}

						if taskList.TargetResource.MariaDBParam.Table != nil {
							mariaDBParamMap["table"] = taskList.TargetResource.MariaDBParam.Table
						}

						if taskList.TargetResource.MariaDBParam.Resource != nil {
							mariaDBParamMap["resource"] = taskList.TargetResource.MariaDBParam.Resource
						}

						if taskList.TargetResource.MariaDBParam.SnapshotMode != nil {
							mariaDBParamMap["snapshot_mode"] = taskList.TargetResource.MariaDBParam.SnapshotMode
						}

						if taskList.TargetResource.MariaDBParam.KeyColumns != nil {
							mariaDBParamMap["key_columns"] = taskList.TargetResource.MariaDBParam.KeyColumns
						}

						if taskList.TargetResource.MariaDBParam.IsTablePrefix != nil {
							mariaDBParamMap["is_table_prefix"] = taskList.TargetResource.MariaDBParam.IsTablePrefix
						}

						if taskList.TargetResource.MariaDBParam.OutputFormat != nil {
							mariaDBParamMap["output_format"] = taskList.TargetResource.MariaDBParam.OutputFormat
						}

						if taskList.TargetResource.MariaDBParam.IncludeContentChanges != nil {
							mariaDBParamMap["include_content_changes"] = taskList.TargetResource.MariaDBParam.IncludeContentChanges
						}

						if taskList.TargetResource.MariaDBParam.IncludeQuery != nil {
							mariaDBParamMap["include_query"] = taskList.TargetResource.MariaDBParam.IncludeQuery
						}

						if taskList.TargetResource.MariaDBParam.RecordWithSchema != nil {
							mariaDBParamMap["record_with_schema"] = taskList.TargetResource.MariaDBParam.RecordWithSchema
						}

						targetResourceMap["maria_d_b_param"] = []interface{}{mariaDBParamMap}
					}

					if taskList.TargetResource.SQLServerParam != nil {
						sQLServerParamMap := map[string]interface{}{}

						if taskList.TargetResource.SQLServerParam.Database != nil {
							sQLServerParamMap["database"] = taskList.TargetResource.SQLServerParam.Database
						}

						if taskList.TargetResource.SQLServerParam.Table != nil {
							sQLServerParamMap["table"] = taskList.TargetResource.SQLServerParam.Table
						}

						if taskList.TargetResource.SQLServerParam.Resource != nil {
							sQLServerParamMap["resource"] = taskList.TargetResource.SQLServerParam.Resource
						}

						if taskList.TargetResource.SQLServerParam.SnapshotMode != nil {
							sQLServerParamMap["snapshot_mode"] = taskList.TargetResource.SQLServerParam.SnapshotMode
						}

						targetResourceMap["s_q_l_server_param"] = []interface{}{sQLServerParamMap}
					}

					if taskList.TargetResource.CtsdbParam != nil {
						ctsdbParamMap := map[string]interface{}{}

						if taskList.TargetResource.CtsdbParam.Resource != nil {
							ctsdbParamMap["resource"] = taskList.TargetResource.CtsdbParam.Resource
						}

						if taskList.TargetResource.CtsdbParam.CtsdbMetric != nil {
							ctsdbParamMap["ctsdb_metric"] = taskList.TargetResource.CtsdbParam.CtsdbMetric
						}

						targetResourceMap["ctsdb_param"] = []interface{}{ctsdbParamMap}
					}

					if taskList.TargetResource.ScfParam != nil {
						scfParamMap := map[string]interface{}{}

						if taskList.TargetResource.ScfParam.FunctionName != nil {
							scfParamMap["function_name"] = taskList.TargetResource.ScfParam.FunctionName
						}

						if taskList.TargetResource.ScfParam.Namespace != nil {
							scfParamMap["namespace"] = taskList.TargetResource.ScfParam.Namespace
						}

						if taskList.TargetResource.ScfParam.Qualifier != nil {
							scfParamMap["qualifier"] = taskList.TargetResource.ScfParam.Qualifier
						}

						if taskList.TargetResource.ScfParam.BatchSize != nil {
							scfParamMap["batch_size"] = taskList.TargetResource.ScfParam.BatchSize
						}

						if taskList.TargetResource.ScfParam.MaxRetries != nil {
							scfParamMap["max_retries"] = taskList.TargetResource.ScfParam.MaxRetries
						}

						targetResourceMap["scf_param"] = []interface{}{scfParamMap}
					}

					taskListMap["target_resource"] = []interface{}{targetResourceMap}
				}

				if taskList.CreateTime != nil {
					taskListMap["create_time"] = taskList.CreateTime
				}

				if taskList.ErrorMessage != nil {
					taskListMap["error_message"] = taskList.ErrorMessage
				}

				if taskList.TaskProgress != nil {
					taskListMap["task_progress"] = taskList.TaskProgress
				}

				if taskList.TaskCurrentStep != nil {
					taskListMap["task_current_step"] = taskList.TaskCurrentStep
				}

				if taskList.DatahubId != nil {
					taskListMap["datahub_id"] = taskList.DatahubId
				}

				if taskList.StepList != nil {
					taskListMap["step_list"] = taskList.StepList
				}

				taskListList = append(taskListList, taskListMap)
			}

			describeDatahubTasksResMap["task_list"] = []interface{}{taskListList}
		}

		ids = append(ids, *result.TaskId)
		_ = d.Set("result", describeDatahubTasksResMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), describeDatahubTasksResMap); e != nil {
			return e
		}
	}
	return nil
}
