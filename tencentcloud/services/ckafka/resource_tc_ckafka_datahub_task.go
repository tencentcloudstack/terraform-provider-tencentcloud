package ckafka

import (
	"context"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCkafkaDatahubTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaDatahubTaskCreate,
		Read:   resourceTencentCloudCkafkaDatahubTaskRead,
		Update: resourceTencentCloudCkafkaDatahubTaskUpdate,
		Delete: resourceTencentCloudCkafkaDatahubTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"task_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "name of the task.",
			},

			"task_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "type of the task, SOURCE(data input), SINK(data output).",
			},

			"source_resource": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "data resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "resource type.",
						},
						"kafka_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "ckafka configuration, required when Type is KAFKA.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"self_built": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "whether the cluster is built by yourself instead of cloud product.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "instance resource.",
									},
									"topic": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Topic name, use `,` when more than 1 topic.",
									},
									"offset_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Offset type, from beginning:earliest, from latest:latest, from specific time:timestamp.",
									},
									"start_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "when Offset type timestamp is required.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "instance name.",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Zone ID.",
									},
									"topic_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Topic id.",
									},
									"partition_num": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "the partition num of the topic.",
									},
									"enable_toleration": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "enable dead letter queue.",
									},
									"qps_limit": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Qps(query per seconds) limit.",
									},
									"table_mappings": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "maps of table to topic, required when multi topic is selected.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "database name.",
												},
												"table": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "table name,use, to separate.",
												},
												"topic": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Topic name.",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Topic ID.",
												},
											},
										},
									},
									"use_table_mapping": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "whether to use multi table.",
									},
									"use_auto_create_topic": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Does the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).",
									},
									"compression_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.",
									},
									"msg_multiple": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
									},
								},
							},
						},
						"event_bus_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "EB configuration, required when type is EB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "resource type, EB_COS/EB_ES/EB_CLS.",
									},
									"self_built": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether it is a self-built cluster.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "instance id.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SCF namespace.",
									},
									"function_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SCF function name.",
									},
									"qualifier": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SCF version and alias.",
									},
								},
							},
						},
						"mongo_db_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "MongoDB config, Required when Type is MONGODB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MongoDB database name.",
									},
									"collection": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MongoDB collection.",
									},
									"copy_existing": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether to copy the stock data, the default parameter is true.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "resource id.",
									},
									"ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Mongo DB connection ip.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "MongoDB connection port.",
									},
									"user_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "MongoDB database user name.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "MongoDB database password.",
									},
									"listening_event": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Listening event type, if it is empty, it means select all. Values include insert, update, replace, delete, invalidate, drop, dropdatabase, rename, used between multiple types, separated by commas.",
									},
									"read_preference": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Master-slave priority, default master node.",
									},
									"pipeline": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "aggregation pipeline.",
									},
									"self_built": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether it is a self-built cluster.",
									},
								},
							},
						},
						"es_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Es configuration, required when Type is ES.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Es connection port.",
									},
									"user_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Es UserName.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Es Password.",
									},
									"self_built": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether it is a self-built cluster.",
									},
									"service_vip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "instance vip.",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "instance vpc id.",
									},
									"drop_invalid_message": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether Es discards the message of parsing failure.",
									},
									"index": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Es index name.",
									},
									"date_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Es date suffix.",
									},
									"content_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "key for data in non-json format.",
									},
									"drop_invalid_json_message": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether Es discards messages in non-json format.",
									},
									"document_id_field": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The field name of the document ID value dumped into Es.",
									},
									"index_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Es custom index name type, STRING, JSONPATH, the default is STRING.",
									},
									"drop_cls": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"drop_invalid_message_to_cls": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to deliver to cls.",
												},
												"drop_cls_region": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The region where the cls is delivered.",
												},
												"drop_cls_owneruin": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Delivery account of cls.",
												},
												"drop_cls_topic_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "topic of cls.",
												},
												"drop_cls_log_set": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls log set.",
												},
											},
										},
									},
									"database_primary_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the message dumped to ES is the binlog of Database, if you need to synchronize database operations, that is, fill in the primary key of the database table when adding, deleting, and modifying operations to ES.",
									},
									"drop_dlq": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "dead letter queue.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "type, DLQ dead letter queue, IGNORE_ERROR|DROP.",
												},
												"kafka_param": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Ckafka type dlq.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"self_built": {
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Whether it is a self-built cluster.",
															},
															"resource": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "resource id.",
															},
															"topic": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Topic name, multiple separated by `,`.",
															},
															"offset_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
															},
															"start_time": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
															},
															"resource_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "resource id name.",
															},
															"zone_id": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Zone ID.",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Topic Id.",
															},
															"partition_num": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Partition num.",
															},
															"enable_toleration": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Enable the fault-tolerant instance and enable the dead-letter queue.",
															},
															"qps_limit": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Qps limit.",
															},
															"table_mappings": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"database": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "database name.",
																		},
																		"table": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Table name, multiple tables, separated by (commas).",
																		},
																		"topic": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Topic name.",
																		},
																		"topic_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Topic ID.",
																		},
																	},
																},
															},
															"use_table_mapping": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Distribute to multiple topics switch, the default is false.",
															},
															"use_auto_create_topic": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).",
															},
															"compression_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.",
															},
															"msg_multiple": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
															},
														},
													},
												},
												"retry_interval": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "retry interval.",
												},
												"max_retry_attempts": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "retry times.",
												},
												"topic_param": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "DIP Topic type dead letter queue.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The topic name of the topic sold separately.",
															},
															"offset_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
															},
															"start_time": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "TopicId.",
															},
															"compression_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.",
															},
															"use_auto_create_topic": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).",
															},
															"msg_multiple": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
															},
														},
													},
												},
												"dlq_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "dlq type, CKAFKA|TOPIC.",
												},
											},
										},
									},
								},
							},
						},
						"tdw_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Tdw configuration, required when Type is TDW.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bid": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tdw bid.",
									},
									"tid": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tdw tid.",
									},
									"is_domestic": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "default true.",
									},
									"tdw_host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TDW address, defalt tl-tdbank-tdmanager.tencent-distribute.com.",
									},
									"tdw_port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "TDW port, default 8099.",
									},
								},
							},
						},
						"dts_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Dts configuration, required when Type is DTS.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Dts instance Id.",
									},
									"ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dts connection ip.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Dts connection port.",
									},
									"topic": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dts topic.",
									},
									"group_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dts consumer group Id.",
									},
									"group_user": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dts account.",
									},
									"group_password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dts consumer group passwd.",
									},
									"tran_sql": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "False to synchronize the original data, true to synchronize the parsed json format data, the default is true.",
									},
								},
							},
						},
						"click_house_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "ClickHouse config, Type CLICKHOUSE requierd.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "ClickHouse cluster.",
									},
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "ClickHouse database name.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "ClickHouse table.",
									},
									"schema": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "ClickHouse schema.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"column_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "column name.",
												},
												"json_key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The json Key name corresponding to this column.",
												},
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "type of table column.",
												},
												"allow_null": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether the column item is allowed to be empty.",
												},
											},
										},
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "resource id.",
									},
									"ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ClickHouse ip.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "ClickHouse port.",
									},
									"user_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ClickHouse user name.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ClickHouse passwd.",
									},
									"service_vip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "instance vip.",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "instance vpc id.",
									},
									"self_built": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether it is a self-built cluster.",
									},
									"drop_invalid_message": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether ClickHouse discards the message that fails to parse, the default is true.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ClickHouse type, emr-clickhouse: emr;cdw-clickhouse: cdwch; selfBuilt: ``.",
									},
									"drop_cls": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"drop_invalid_message_to_cls": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to deliver to cls.",
												},
												"drop_cls_region": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls region.",
												},
												"drop_cls_owneruin": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls account.",
												},
												"drop_cls_topic_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls topicId.",
												},
												"drop_cls_log_set": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls LogSet id.",
												},
											},
										},
									},
								},
							},
						},
						"cls_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Cls configuration, Required when Type is CLS.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"decode_json": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the produced information is in json format.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "cls id.",
									},
									"log_set": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "LogSet id.",
									},
									"content_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Required when Decode Json is false.",
									},
									"time_field": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specify the content of a field in the message as the time of the cls log. The format of the field content needs to be a second-level timestamp.",
									},
								},
							},
						},
						"cos_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Cos configuration, required when Type is COS.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "cos bucket name.",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "region code.",
									},
									"object_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ObjectKey.",
									},
									"aggregate_batch_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The size of aggregated messages MB.",
									},
									"aggregate_interval": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "time interval.",
									},
									"format_output_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The file format after message aggregation csv|json.",
									},
									"object_key_prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dumped object directory prefix.",
									},
									"directory_time_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Partition format formatted according to strptime time.",
									},
								},
							},
						},
						"my_sql_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "MySQL configuration, Required when Type is MYSQL.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MySQL database name, * is the whole database.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the MySQL data table,  is the non-system table in all the monitored databases, which can be separated by, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name, when a regular expression needs to be filled in, the format is data database name.data table name.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MySQL connection Id.",
									},
									"snapshot_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "whether to Copy inventory information (schema_only does not copy, initial full amount), the default is initial.",
									},
									"ddl_topic": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Topic that stores the Ddl information of My SQL, if it is empty, it will not be stored by default.",
									},
									"data_source_monitor_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TABLE indicates that the read item is a table, QUERY indicates that the read item is a query.",
									},
									"data_source_monitor_resource": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When DataMonitorMode=TABLE, pass in the Table that needs to be read; when DataMonitorMode=QUERY, pass in the query sql statement that needs to be read.",
									},
									"data_source_increment_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TIMESTAMP indicates that the incremental column is of timestamp type, INCREMENT indicates that the incremental column is of self-incrementing id type&#39;.",
									},
									"data_source_increment_column": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the column to be monitored.",
									},
									"data_source_start_from": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "HEAD means copy stock + incremental data, TAIL means copy only incremental data.",
									},
									"data_target_insert_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.",
									},
									"data_target_primary_key_field": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.",
									},
									"data_target_record_mapping": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Mapping relationship between tables and messages.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"json_key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The key name of the message.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "message type.",
												},
												"allow_null": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether the message is allowed to be empty.",
												},
												"column_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Corresponding mapping column name.",
												},
												"extra_info": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database table extra fields.",
												},
												"column_size": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "current column size.",
												},
												"decimal_digits": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "current column precision.",
												},
												"auto_increment": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether it is an auto-increment column.",
												},
												"default_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database table default parameters.",
												},
											},
										},
									},
									"topic_regex": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Regular expression for routing events to specific topics, defaults to (.*).",
									},
									"topic_replacement": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TopicRegex, $1, $2.",
									},
									"key_columns": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Format library1.table1 field 1,field 2;library 2.table2 field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
									},
									"drop_invalid_message": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to discard messages that fail to parse, the default is true.",
									},
									"drop_cls": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"drop_invalid_message_to_cls": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to deliver to cls.",
												},
												"drop_cls_region": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The region where the cls is delivered.",
												},
												"drop_cls_owneruin": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "account.",
												},
												"drop_cls_topic_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls topic.",
												},
												"drop_cls_log_set": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls LogSet id.",
												},
											},
										},
									},
									"output_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "output format, DEFAULT, CANAL_1, CANAL_2.",
									},
									"is_table_prefix": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "When the Table input is a prefix, the value of this item is true, otherwise it is false.",
									},
									"include_content_changes": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.",
									},
									"include_query": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.",
									},
									"record_with_schema": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
									},
									"signal_database": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "database name of signal table.",
									},
									"is_table_regular": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the input table is a regular expression, if this option and Is Table Prefix are true at the same time, the judgment priority of this option is higher than Is Table Prefix.",
									},
								},
							},
						},
						"postgre_sql_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "PostgreSQL configuration, Required when Type is POSTGRESQL or TDSQL C_POSTGRESQL.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "PostgreSQL database name.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "PostgreSQL tableName, * is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of Schema name.Data table name, and you need to fill in a regular expression When, the format is Schema name.data table name.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "PostgreSQL connection Id.",
									},
									"plugin_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "(decoderbufs/pgoutput), default decoderbufs.",
									},
									"snapshot_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "never|initial, default initial.",
									},
									"data_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Upstream data format (JSON|Debezium), required when the database synchronization mode matches the default field.",
									},
									"data_target_insert_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.",
									},
									"data_target_primary_key_field": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.",
									},
									"data_target_record_mapping": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Mapping relationship between tables and messages.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"json_key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The key name of the message.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "message type.",
												},
												"allow_null": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether the message is allowed to be empty.",
												},
												"column_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Column Name.",
												},
												"extra_info": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database table extra fields.",
												},
												"column_size": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "current ColumnSize.",
												},
												"decimal_digits": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "current Column DecimalDigits.",
												},
												"auto_increment": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether it is an auto-increment column.",
												},
												"default_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database table default parameters.",
												},
											},
										},
									},
									"drop_invalid_message": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to discard messages that fail to parse, the default is true.",
									},
									"is_table_regular": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the input table is a regular expression.",
									},
									"key_columns": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Format  library1.table1:field 1,field2;library2.table2:field2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
									},
									"record_with_schema": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
									},
								},
							},
						},
						"topic_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Topic configuration, Required when Type is Topic.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The topic name of the topic sold separately.",
									},
									"offset_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
									},
									"start_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
									},
									"topic_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Topic TopicId.",
									},
									"compression_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.",
									},
									"use_auto_create_topic": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).",
									},
									"msg_multiple": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
									},
								},
							},
						},
						"maria_db_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "MariaDB configuration, Required when Type is MARIADB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MariaDB database name, * for all database.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MariaDB db name, *is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MariaDB connection Id.",
									},
									"snapshot_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "schema_only|initial, default initial.",
									},
									"key_columns": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Format  library 1. table 1: field 1, field 2; library 2. table 2: field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
									},
									"is_table_prefix": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "When the Table input is a prefix, the value of this item is true, otherwise it is false.",
									},
									"output_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "output format, DEFAULT, CANAL_1, CANAL_2.",
									},
									"include_content_changes": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.",
									},
									"include_query": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.",
									},
									"record_with_schema": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
									},
								},
							},
						},
						"sql_server_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "SQLServer configuration, Required when Type is SQLSERVER.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SQLServer database name.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SQLServer table, *is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SQLServer connection Id.",
									},
									"snapshot_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "schema_only|initial default initial.",
									},
								},
							},
						},
						"ctsdb_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Ctsdb configuration, Required when Type is CTSDB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "resource id.",
									},
									"ctsdb_metric": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Ctsdb metric.",
									},
								},
							},
						},
						"scf_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Scf configuration, Required when Type is SCF.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"function_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SCF function name.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SCF cloud function namespace, the default is default.",
									},
									"qualifier": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SCF cloud function version and alias, the default is DEFAULT.",
									},
									"batch_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum number of messages sent in each batch, the default is 1000.",
									},
									"max_retries": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of retries after the SCF call fails, the default is 5.",
									},
								},
							},
						},
					},
				},
			},

			"target_resource": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Target Resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource Type.",
						},
						"kafka_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "ckafka configuration, required when Type is KAFKA.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"self_built": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "whether the cluster is built by yourself instead of cloud product.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "instance resource.",
									},
									"topic": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Topic name, use `,` when more than 1 topic.",
									},
									"offset_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Offset type, from beginning:earliest, from latest:latest, from specific time:timestamp.",
									},
									"start_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "when Offset type timestamp is required.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "instance name.",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Zone ID.",
									},
									"topic_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Topic id.",
									},
									"partition_num": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "the partition num of the topic.",
									},
									"enable_toleration": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "enable dead letter queue.",
									},
									"qps_limit": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Qps(query per seconds) limit.",
									},
									"table_mappings": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "maps of table to topic, required when multi topic is selected.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "database name.",
												},
												"table": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "table name,use, to separate.",
												},
												"topic": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Topic name.",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Topic ID.",
												},
											},
										},
									},
									"use_table_mapping": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "whether to use multi table.",
									},
									"use_auto_create_topic": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Does the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).",
									},
									"compression_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.",
									},
									"msg_multiple": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
									},
								},
							},
						},
						"event_bus_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "EB configuration, required when type is EB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "resource type. EB_COS/EB_ES/EB_CLS.",
									},
									"self_built": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether it is a self-built cluster.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "instance id.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SCF namespace.",
									},
									"function_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SCF function name.",
									},
									"qualifier": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SCF version and alias.",
									},
								},
							},
						},
						"mongo_db_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "MongoDB config, Required when Type is MONGODB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MongoDB database name.",
									},
									"collection": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MongoDB collection.",
									},
									"copy_existing": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether to copy the stock data, the default parameter is true.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "resource id.",
									},
									"ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Mongo DB connection ip.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "MongoDB connection port.",
									},
									"user_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "MongoDB database user name.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "MongoDB database password.",
									},
									"listening_event": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Listening event type, if it is empty, it means select all. Values include insert, update, replace, delete, invalidate, drop, dropdatabase, rename, used between multiple types, separated by commas.",
									},
									"read_preference": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Master-slave priority, default master node.",
									},
									"pipeline": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "aggregation pipeline.",
									},
									"self_built": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether it is a self-built cluster.",
									},
								},
							},
						},
						"es_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Es configuration, required when Type is ES.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Es connection port.",
									},
									"user_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Es UserName.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Es Password.",
									},
									"self_built": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether it is a self-built cluster.",
									},
									"service_vip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "instance vip.",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "instance vpc id.",
									},
									"drop_invalid_message": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether Es discards the message of parsing failure.",
									},
									"index": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Es index name.",
									},
									"date_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Es date suffix.",
									},
									"content_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "key for data in non-json format.",
									},
									"drop_invalid_json_message": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether Es discards messages in non-json format.",
									},
									"document_id_field": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The field name of the document ID value dumped into Es.",
									},
									"index_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Es custom index name type, STRING, JSONPATH, the default is STRING.",
									},
									"drop_cls": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"drop_invalid_message_to_cls": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to deliver to cls.",
												},
												"drop_cls_region": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The region where the cls is delivered.",
												},
												"drop_cls_owneruin": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Delivery account of cls.",
												},
												"drop_cls_topic_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "topic of cls.",
												},
												"drop_cls_log_set": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls log set.",
												},
											},
										},
									},
									"database_primary_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the message dumped to ES is the binlog of Database, if you need to synchronize database operations, that is, fill in the primary key of the database table when adding, deleting, and modifying operations to ES.",
									},
									"drop_dlq": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "dead letter queue.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "type, DLQ dead letter queue, IGNORE_ERROR|DROP.",
												},
												"kafka_param": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Ckafka type dlq.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"self_built": {
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Whether it is a self-built cluster.",
															},
															"resource": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "resource id.",
															},
															"topic": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Topic name, multiple separated by,.",
															},
															"offset_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
															},
															"start_time": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
															},
															"resource_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "resource id name.",
															},
															"zone_id": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Zone ID.",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Topic Id.",
															},
															"partition_num": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Partition num.",
															},
															"enable_toleration": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Enable the fault-tolerant instance and enable the dead-letter queue.",
															},
															"qps_limit": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Qps limit.",
															},
															"table_mappings": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"database": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "database name.",
																		},
																		"table": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Table name, multiple tables, separated by (commas).",
																		},
																		"topic": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Topic name.",
																		},
																		"topic_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Topic ID.",
																		},
																	},
																},
															},
															"use_table_mapping": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Distribute to multiple topics switch, the default is false.",
															},
															"use_auto_create_topic": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).",
															},
															"compression_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.",
															},
															"msg_multiple": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
															},
														},
													},
												},
												"retry_interval": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "retry interval.",
												},
												"max_retry_attempts": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "retry times.",
												},
												"topic_param": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "DIP Topic type dead letter queue.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The topic name of the topic sold separately.",
															},
															"offset_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
															},
															"start_time": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "TopicId.",
															},
															"compression_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.",
															},
															"use_auto_create_topic": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).",
															},
															"msg_multiple": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
															},
														},
													},
												},
												"dlq_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "dlq type, CKAFKA|TOPIC.",
												},
											},
										},
									},
								},
							},
						},
						"tdw_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Tdw configuration, required when Type is TDW.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bid": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tdw bid.",
									},
									"tid": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tdw tid.",
									},
									"is_domestic": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "default true.",
									},
									"tdw_host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TDW address, defalt tl-tdbank-tdmanager.tencent-distribute.com.",
									},
									"tdw_port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "TDW port, default 8099.",
									},
								},
							},
						},
						"dts_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Dts configuration, required when Type is DTS.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Dts instance Id.",
									},
									"ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dts connection ip.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Dts connection port.",
									},
									"topic": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dts topic.",
									},
									"group_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dts consumer group Id.",
									},
									"group_user": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dts account.",
									},
									"group_password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dts consumer group passwd.",
									},
									"tran_sql": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "False to synchronize the original data, true to synchronize the parsed json format data, the default is true.",
									},
								},
							},
						},
						"click_house_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "ClickHouse config, Type CLICKHOUSE requierd.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "ClickHouse cluster.",
									},
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "ClickHouse database name.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "ClickHouse table.",
									},
									"schema": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "ClickHouse schema.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"column_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "column name.",
												},
												"json_key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The json Key name corresponding to this column.",
												},
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "type of table column.",
												},
												"allow_null": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether the column item is allowed to be empty.",
												},
											},
										},
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "resource id.",
									},
									"ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ClickHouse ip.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "ClickHouse port.",
									},
									"user_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ClickHouse user name.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ClickHouse passwd.",
									},
									"service_vip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "instance vip.",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "instance vpc id.",
									},
									"self_built": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether it is a self-built cluster.",
									},
									"drop_invalid_message": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether ClickHouse discards the message that fails to parse, the default is true.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ClickHouse type, emr-clickhouse: emr;cdw-clickhouse: cdwch;selfBuilt: ``.",
									},
									"drop_cls": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"drop_invalid_message_to_cls": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to deliver to cls.",
												},
												"drop_cls_region": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls region.",
												},
												"drop_cls_owneruin": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls account.",
												},
												"drop_cls_topic_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls topicId.",
												},
												"drop_cls_log_set": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls LogSet id.",
												},
											},
										},
									},
								},
							},
						},
						"cls_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Cls configuration, Required when Type is CLS.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"decode_json": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the produced information is in json format.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "cls id.",
									},
									"log_set": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "LogSet id.",
									},
									"content_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Required when Decode Json is false.",
									},
									"time_field": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specify the content of a field in the message as the time of the cls log. The format of the field content needs to be a second-level timestamp.",
									},
								},
							},
						},
						"cos_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Cos configuration, required when Type is COS.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "cos bucket name.",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "region code.",
									},
									"object_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ObjectKey.",
									},
									"aggregate_batch_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The size of aggregated messages MB.",
									},
									"aggregate_interval": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "time interval.",
									},
									"format_output_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The file format after message aggregation csv|json.",
									},
									"object_key_prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dumped object directory prefix.",
									},
									"directory_time_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Partition format formatted according to strptime time.",
									},
								},
							},
						},
						"my_sql_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "MySQL configuration, Required when Type is MYSQL.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MySQL database name, * is the whole database.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the MySQL data table,  is the non-system table in all the monitored databases, which can be separated by, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name, when a regular expression needs to be filled in, the format is data database name.data table name.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MySQL connection Id.",
									},
									"snapshot_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "whether to Copy inventory information (schema_only does not copy, initial full amount), the default is initial.",
									},
									"ddl_topic": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Topic that stores the Ddl information of My SQL, if it is empty, it will not be stored by default.",
									},
									"data_source_monitor_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TABLE indicates that the read item is a table, QUERY indicates that the read item is a query.",
									},
									"data_source_monitor_resource": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When DataMonitorMode=TABLE, pass in the Table that needs to be read; when DataMonitorMode=QUERY, pass in the query sql statement that needs to be read.",
									},
									"data_source_increment_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TIMESTAMP indicates that the incremental column is of timestamp type, INCREMENT indicates that the incremental column is of self-incrementing id type.",
									},
									"data_source_increment_column": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "the name of the column to be monitored.",
									},
									"data_source_start_from": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "HEAD means copy stock + incremental data, TAIL means copy only incremental data.",
									},
									"data_target_insert_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.",
									},
									"data_target_primary_key_field": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.",
									},
									"data_target_record_mapping": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Mapping relationship between tables and messages.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"json_key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The key name of the message.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "message type.",
												},
												"allow_null": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether the message is allowed to be empty.",
												},
												"column_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Corresponding mapping column name.",
												},
												"extra_info": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database table extra fields.",
												},
												"column_size": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "current column size.",
												},
												"decimal_digits": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "current column precision.",
												},
												"auto_increment": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether it is an auto-increment column.",
												},
												"default_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database table default parameters.",
												},
											},
										},
									},
									"topic_regex": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Regular expression for routing events to specific topics, defaults to (.*).",
									},
									"topic_replacement": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TopicRegex, $1, $2.",
									},
									"key_columns": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Format library1.table1 field 1,field 2;library 2.table2 field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
									},
									"drop_invalid_message": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to discard messages that fail to parse, the default is true.",
									},
									"drop_cls": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"drop_invalid_message_to_cls": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to deliver to cls.",
												},
												"drop_cls_region": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The region where the cls is delivered.",
												},
												"drop_cls_owneruin": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "account.",
												},
												"drop_cls_topic_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls topic.",
												},
												"drop_cls_log_set": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cls LogSet id.",
												},
											},
										},
									},
									"output_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "output format, DEFAULT, CANAL_1, CANAL_2.",
									},
									"is_table_prefix": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "When the Table input is a prefix, the value of this item is true, otherwise it is false.",
									},
									"include_content_changes": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.",
									},
									"include_query": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.",
									},
									"record_with_schema": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
									},
									"signal_database": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "database name of signal table.",
									},
									"is_table_regular": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the input table is a regular expression, if this option and Is Table Prefix are true at the same time, the judgment priority of this option is higher than Is Table Prefix.",
									},
								},
							},
						},
						"postgre_sql_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "PostgreSQL configuration, Required when Type is POSTGRESQL or TDSQL C_POSTGRESQL.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "PostgreSQL database name.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "PostgreSQL tableName, * is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of Schema name.Data table name, and you need to fill in a regular expression When, the format is Schema name.data table name.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "PostgreSQL connection Id.",
									},
									"plugin_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "(decoderbufs/pgoutput), default decoderbufs.",
									},
									"snapshot_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "never|initial, default initial.",
									},
									"data_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Upstream data format (JSON|Debezium), required when the database synchronization mode matches the default field.",
									},
									"data_target_insert_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.",
									},
									"data_target_primary_key_field": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.",
									},
									"data_target_record_mapping": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Mapping relationship between tables and messages.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"json_key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The key name of the message.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "message type.",
												},
												"allow_null": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether the message is allowed to be empty.",
												},
												"column_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Column Name.",
												},
												"extra_info": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database table extra fields.",
												},
												"column_size": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "current ColumnSize.",
												},
												"decimal_digits": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "current Column DecimalDigits.",
												},
												"auto_increment": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether it is an auto-increment column.",
												},
												"default_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database table default parameters.",
												},
											},
										},
									},
									"drop_invalid_message": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to discard messages that fail to parse, the default is true.",
									},
									"is_table_regular": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the input table is a regular expression.",
									},
									"key_columns": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Format  library1.table1:field 1,field2;library2.table2:field2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
									},
									"record_with_schema": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
									},
								},
							},
						},
						"topic_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Topic configuration, Required when Type is Topic.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The topic name of the topic sold separately.",
									},
									"offset_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
									},
									"start_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
									},
									"topic_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Topic TopicId.",
									},
									"compression_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.",
									},
									"use_auto_create_topic": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).",
									},
									"msg_multiple": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
									},
								},
							},
						},
						"maria_db_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "MariaDB configuration, Required when Type is MARIADB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MariaDB database name, * for all database.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MariaDB db name, *is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MariaDB connection Id.",
									},
									"snapshot_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "schema_only|initial, default initial.",
									},
									"key_columns": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Format  library 1. table 1: field 1, field 2; library 2. table 2: field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.",
									},
									"is_table_prefix": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "When the Table input is a prefix, the value of this item is true, otherwise it is false.",
									},
									"output_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "output format, DEFAULT, CANAL_1, CANAL_2.",
									},
									"include_content_changes": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.",
									},
									"include_query": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.",
									},
									"record_with_schema": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.",
									},
								},
							},
						},
						"sql_server_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "SQLServer configuration, Required when Type is SQLSERVER.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SQLServer database name.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SQLServer table, *is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.",
									},
									"resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SQLServer connection Id.",
									},
									"snapshot_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "schema_only|initial default initial.",
									},
								},
							},
						},
						"ctsdb_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Ctsdb configuration, Required when Type is CTSDB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "resource id.",
									},
									"ctsdb_metric": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Ctsdb metric.",
									},
								},
							},
						},
						"scf_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Scf configuration, Required when Type is SCF.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"function_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SCF function name.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SCF cloud function namespace, the default is default.",
									},
									"qualifier": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SCF cloud function version and alias, the default is DEFAULT.",
									},
									"batch_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum number of messages sent in each batch, the default is 1000.",
									},
									"max_retries": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of retries after the SCF call fails, the default is 5.",
									},
								},
							},
						},
					},
				},
			},

			"transform_param": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Data Processing Rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"analysis_format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "parsing format, JSON | DELIMITER| REGULAR.",
						},
						"output_format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "output format.",
						},
						"failure_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Whether to keep parsing failure data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "type, DLQ dead letter queue, IGNORE_ERROR|DROP.",
									},
									"kafka_param": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Ckafka type dlq.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"self_built": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether it is a self-built cluster.",
												},
												"resource": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "resource id.",
												},
												"topic": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Topic name, multiple separated by,.",
												},
												"offset_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
												},
												"resource_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "resource id name.",
												},
												"zone_id": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Zone ID.",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Topic Id.",
												},
												"partition_num": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Partition num.",
												},
												"enable_toleration": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Enable the fault-tolerant instance and enable the dead-letter queue.",
												},
												"qps_limit": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Qps limit.",
												},
												"table_mappings": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "database name.",
															},
															"table": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Table name, multiple tables, separated by (commas).",
															},
															"topic": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Topic name.",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Topic ID.",
															},
														},
													},
												},
												"use_table_mapping": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Distribute to multiple topics switch, the default is false.",
												},
												"use_auto_create_topic": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).",
												},
												"compression_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.",
												},
												"msg_multiple": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
												},
											},
										},
									},
									"retry_interval": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "retry interval.",
									},
									"max_retry_attempts": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "retry times.",
									},
									"topic_param": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "DIP Topic type dead letter queue.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The topic name of the topic sold separately.",
												},
												"offset_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "TopicId.",
												},
												"compression_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.",
												},
												"use_auto_create_topic": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).",
												},
												"msg_multiple": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
												},
											},
										},
									},
									"dlq_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "dlq type, CKAFKA|TOPIC.",
									},
								},
							},
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Raw data.",
						},
						"source_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Data source, TOPIC pulls from the source topic, CUSTOMIZE custom.",
						},
						"regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "delimiter, regular expression.",
						},
						"map_param": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Map.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "key.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Type, DEFAULT default, DATE system default - timestamp, CUSTOMIZE custom, MAPPING mapping.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "value.",
									},
								},
							},
						},
						"filter_param": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "filter.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Key.",
									},
									"match_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Matching mode, prefix matches PREFIX, suffix matches SUFFIX, contains matches CONTAINS, except matches EXCEPT, value matches NUMBER, IP matches IP.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Value.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "REGULAR.",
									},
								},
							},
						},
						"result": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Test Results.",
						},
						"analyse_result": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Analysis result.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "key.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Type, DEFAULT default, DATE system default - timestamp, CUSTOMIZE custom, MAPPING mapping.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "value.",
									},
								},
							},
						},
						"use_event_bus": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the underlying engine uses eb.",
						},
					},
				},
			},

			"schema_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "SchemaId.",
			},

			"transforms_param": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Data processing rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Raw data.",
						},
						"field_chain": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "processing chain.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"analyse": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: "analyze.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"format": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Parsing format, JSON, DELIMITER delimiter, REGULAR regular extraction, SOURCE processing all results of the upper layer.",
												},
												"regex": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "delimiter, regular expression.",
												},
												"input_value_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "KEY to be processed again - mode.",
												},
												"input_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "KEY to be processed again - KEY expression.",
												},
											},
										},
									},
									"secondary_analyse": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "secondary analysis.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"regex": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "delimiter.",
												},
											},
										},
									},
									"s_m_t": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "data processing.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "KEY.",
												},
												"operate": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Operation, DATE system preset - timestamp, CUSTOMIZE customization, MAPPING mapping, JSONPATH.",
												},
												"scheme_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "data type, ORIGINAL, STRING, INT64, FLOAT64, BOOLEAN, MAP, ARRAY.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "VALUE.",
												},
												"value_operate": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "VALUE process.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Processing mode, REPLACE replacement, SUBSTR interception, DATE date conversion, TRIM removal of leading and trailing spaces, REGEX REPLACE regular replacement, URL DECODE, LOWERCASE conversion to lowercase.",
															},
															"replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "replace, TYPE=REPLACE is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "been replaced value.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"substr": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Substr, TYPE=SUBSTR is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "interception starting position.",
																		},
																		"end": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "cut-off position.",
																		},
																	},
																},
															},
															"date": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Time conversion, required when TYPE=DATE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Time format.",
																		},
																		"target_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "input type, string|unix.",
																		},
																		"time_zone": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "default GMT+8.",
																		},
																	},
																},
															},
															"regex_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Regular replacement, required when TYPE=REGEX REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Regular.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"split": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The value supports one split and multiple values, required when TYPE=SPLIT.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																	},
																},
															},
															"k_v": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Key-value secondary analysis, must be passed when TYPE=KV.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"delimiter": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Key-value secondary analysis delimiter.",
																		},
																		"keep_original_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Keep the source Key, the default is false not to keep.",
																		},
																	},
																},
															},
															"result": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "result.",
															},
															"json_path_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Json Path replacement, must pass when TYPE=JSON PATH REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replaced value, Jsonpath expression.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replacement value, Jsonpath expression or string.",
																		},
																	},
																},
															},
															"url_decode": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Url parsing.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"charset_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "code.",
																		},
																	},
																},
															},
														},
													},
												},
												"original_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "OriginalValue.",
												},
												"value_operates": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "VALUE process chain.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Processing mode, REPLACE replacement, SUBSTR interception, DATE date conversion, TRIM removal of leading and trailing spaces, REGEX REPLACE regular replacement, URL DECODE, LOWERCASE conversion to lowercase.",
															},
															"replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "replace, TYPE=REPLACE is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "been replaced value.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"substr": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Substr, TYPE=SUBSTR is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "interception starting position.",
																		},
																		"end": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "cut-off position.",
																		},
																	},
																},
															},
															"date": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Time conversion, required when TYPE=DATE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Time format.",
																		},
																		"target_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "input type, string|unix.",
																		},
																		"time_zone": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "default GMT+8.",
																		},
																	},
																},
															},
															"regex_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Regular replacement, required when TYPE=REGEX REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Regular.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"split": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The value supports one split and multiple values, required when TYPE=SPLIT.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																	},
																},
															},
															"k_v": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Key-value secondary analysis, must be passed when TYPE=KV.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"delimiter": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Key-value secondary analysis delimiter.",
																		},
																		"keep_original_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Keep the source Key, the default is false not to keep.",
																		},
																	},
																},
															},
															"result": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "result.",
															},
															"json_path_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Json Path replacement, must pass when TYPE=JSON PATH REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replaced value, Jsonpath expression.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replacement value, Jsonpath expression or string.",
																		},
																	},
																},
															},
															"url_decode": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Url parsing.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"charset_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "code.",
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
									"result": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Test Results.",
									},
									"analyse_result": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Analysis result.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "KEY.",
												},
												"operate": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Operation, DATE system preset - timestamp, CUSTOMIZE customization, MAPPING mapping, JSONPATH.",
												},
												"scheme_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "data type, ORIGINAL, STRING, INT64, FLOAT64, BOOLEAN, MAP, ARRAY.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "VALUE.",
												},
												"value_operate": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "VALUE process.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Processing mode, REPLACE replacement, SUBSTR interception, DATE date conversion, TRIM removal of leading and trailing spaces, REGEX REPLACE regular replacement, URL DECODE, LOWERCASE conversion to lowercase.",
															},
															"replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "replace, TYPE=REPLACE is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "been replaced value.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"substr": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Substr, TYPE=SUBSTR is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "interception starting position.",
																		},
																		"end": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "cut-off position.",
																		},
																	},
																},
															},
															"date": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Time conversion, required when TYPE=DATE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Time format.",
																		},
																		"target_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "input type, string|unix.",
																		},
																		"time_zone": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "default GMT+8.",
																		},
																	},
																},
															},
															"regex_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Regular replacement, required when TYPE=REGEX REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Regular.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"split": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The value supports one split and multiple values, required when TYPE=SPLIT.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																	},
																},
															},
															"k_v": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Key-value secondary analysis, must be passed when TYPE=KV.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"delimiter": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Key-value secondary analysis delimiter.",
																		},
																		"keep_original_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Keep the source Key, the default is false not to keep.",
																		},
																	},
																},
															},
															"result": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "result.",
															},
															"json_path_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Json Path replacement, must pass when TYPE=JSON PATH REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replaced value, Jsonpath expression.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replacement value, Jsonpath expression or string.",
																		},
																	},
																},
															},
															"url_decode": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Url parsing.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"charset_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "code.",
																		},
																	},
																},
															},
														},
													},
												},
												"original_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "OriginalValue.",
												},
												"value_operates": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "VALUE process chain.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Processing mode, REPLACE replacement, SUBSTR interception, DATE date conversion, TRIM removal of leading and trailing spaces, REGEX REPLACE regular replacement, URL DECODE, LOWERCASE conversion to lowercase.",
															},
															"replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "replace, TYPE=REPLACE is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "been replaced value.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"substr": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Substr, TYPE=SUBSTR is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "interception starting position.",
																		},
																		"end": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "cut-off position.",
																		},
																	},
																},
															},
															"date": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Time conversion, required when TYPE=DATE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Time format.",
																		},
																		"target_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "input type, string|unix.",
																		},
																		"time_zone": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "default GMT+8.",
																		},
																	},
																},
															},
															"regex_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Regular replacement, required when TYPE=REGEX REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Regular.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"split": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The value supports one split and multiple values, required when TYPE=SPLIT.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																	},
																},
															},
															"k_v": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Key-value secondary analysis, must be passed when TYPE=KV.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"delimiter": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Key-value secondary analysis delimiter.",
																		},
																		"keep_original_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Keep the source Key, the default is false not to keep.",
																		},
																	},
																},
															},
															"result": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "result.",
															},
															"json_path_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Json Path replacement, must pass when TYPE=JSON PATH REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replaced value, Jsonpath expression.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replacement value, Jsonpath expression or string.",
																		},
																	},
																},
															},
															"url_decode": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Url parsing.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"charset_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "code.",
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
									"secondary_analyse_result": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Secondary Analysis Results.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "KEY.",
												},
												"operate": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Operation, DATE system preset - timestamp, CUSTOMIZE customization, MAPPING mapping, JSONPATH.",
												},
												"scheme_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "data type, ORIGINAL, STRING, INT64, FLOAT64, BOOLEAN, MAP, ARRAY.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "VALUE.",
												},
												"value_operate": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "VALUE process.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Processing mode, REPLACE replacement, SUBSTR interception, DATE date conversion, TRIM removal of leading and trailing spaces, REGEX REPLACE regular replacement, URL DECODE, LOWERCASE conversion to lowercase.",
															},
															"replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "replace, TYPE=REPLACE is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "been replaced value.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"substr": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Substr, TYPE=SUBSTR is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "interception starting position.",
																		},
																		"end": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "cut-off position.",
																		},
																	},
																},
															},
															"date": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Time conversion, required when TYPE=DATE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Time format.",
																		},
																		"target_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "input type, string|unix.",
																		},
																		"time_zone": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "default GMT+8.",
																		},
																	},
																},
															},
															"regex_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Regular replacement, required when TYPE=REGEX REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Regular.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"split": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The value supports one split and multiple values, required when TYPE=SPLIT.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																	},
																},
															},
															"k_v": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Key-value secondary analysis, must be passed when TYPE=KV.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"delimiter": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Key-value secondary analysis delimiter.",
																		},
																		"keep_original_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Keep the source Key, the default is false not to keep.",
																		},
																	},
																},
															},
															"result": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "result.",
															},
															"json_path_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Json Path replacement, must pass when TYPE=JSON PATH REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replaced value, Jsonpath expression.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replacement value, Jsonpath expression or string.",
																		},
																	},
																},
															},
															"url_decode": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Url parsing.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"charset_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "code.",
																		},
																	},
																},
															},
														},
													},
												},
												"original_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "OriginalValue.",
												},
												"value_operates": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "VALUE process chain.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Processing mode, REPLACE replacement, SUBSTR interception, DATE date conversion, TRIM removal of leading and trailing spaces, REGEX REPLACE regular replacement, URL DECODE, LOWERCASE conversion to lowercase.",
															},
															"replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "replace, TYPE=REPLACE is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "been replaced value.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"substr": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Substr, TYPE=SUBSTR is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "interception starting position.",
																		},
																		"end": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "cut-off position.",
																		},
																	},
																},
															},
															"date": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Time conversion, required when TYPE=DATE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Time format.",
																		},
																		"target_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "input type, string|unix.",
																		},
																		"time_zone": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "default GMT+8.",
																		},
																	},
																},
															},
															"regex_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Regular replacement, required when TYPE=REGEX REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Regular.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "new value.",
																		},
																	},
																},
															},
															"split": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The value supports one split and multiple values, required when TYPE=SPLIT.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																	},
																},
															},
															"k_v": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Key-value secondary analysis, must be passed when TYPE=KV.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"delimiter": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "delimiter.",
																		},
																		"regex": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Key-value secondary analysis delimiter.",
																		},
																		"keep_original_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Keep the source Key, the default is false not to keep.",
																		},
																	},
																},
															},
															"result": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "result.",
															},
															"json_path_replace": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Json Path replacement, must pass when TYPE=JSON PATH REPLACE.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"old_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replaced value, Jsonpath expression.",
																		},
																		"new_value": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Replacement value, Jsonpath expression or string.",
																		},
																	},
																},
															},
															"url_decode": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Url parsing.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"charset_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "code.",
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
									"analyse_json_result": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Parsing results in JSON format.",
									},
									"secondary_analyse_json_result": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secondary parsing results in JSON format.",
									},
								},
							},
						},
						"filter_param": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "filter.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Key.",
									},
									"match_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Matching mode, prefix matches PREFIX, suffix matches SUFFIX, contains matches CONTAINS, except matches EXCEPT, value matches NUMBER, IP matches IP.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Value.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "REGULAR.",
									},
								},
							},
						},
						"failure_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "fail process.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "type, DLQ dead letter queue, IGNORE_ERROR|DROP.",
									},
									"kafka_param": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Ckafka type dlq.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"self_built": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether it is a self-built cluster.",
												},
												"resource": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "resource id.",
												},
												"topic": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Topic name, multiple separated by,.",
												},
												"offset_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
												},
												"resource_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "resource id name.",
												},
												"zone_id": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Zone ID.",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Topic Id.",
												},
												"partition_num": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Partition num.",
												},
												"enable_toleration": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Enable the fault-tolerant instance and enable the dead-letter queue.",
												},
												"qps_limit": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Qps limit.",
												},
												"table_mappings": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "database name.",
															},
															"table": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Table name, multiple tables, separated by (commas).",
															},
															"topic": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Topic name.",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Topic ID.",
															},
														},
													},
												},
												"use_table_mapping": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Distribute to multiple topics switch, the default is false.",
												},
												"use_auto_create_topic": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).",
												},
												"compression_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.",
												},
												"msg_multiple": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
												},
											},
										},
									},
									"retry_interval": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "retry interval.",
									},
									"max_retry_attempts": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "retry times.",
									},
									"topic_param": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "DIP Topic type dead letter queue.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The topic name of the topic sold separately.",
												},
												"offset_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Offset type, initial position earliest, latest position latest, time point position timestamp.",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "TopicId.",
												},
												"compression_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.",
												},
												"use_auto_create_topic": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).",
												},
												"msg_multiple": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).",
												},
											},
										},
									},
									"dlq_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "dlq type, CKAFKA|TOPIC.",
									},
								},
							},
						},
						"result": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "result.",
						},
						"source_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "data source.",
						},
						"output_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "output format, JSON, ROW, default JSON.",
						},
						"row_param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The output format is ROW Required.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"row_content": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "row content, KEY_VALUE, VALUE.",
									},
									"key_value_delimiter": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "key, value delimiter.",
									},
									"entry_delimiter": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "delimiter.",
									},
								},
							},
						},
						"keep_metadata": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to keep the data source Topic metadata information (source Topic, Partition, Offset), the default is false.",
						},
						"batch_analyse": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "data process.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"format": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "ONE BY ONE single output, MERGE combined output.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCkafkaDatahubTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_datahub_task.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = ckafka.NewCreateDatahubTaskRequest()
		response = ckafka.NewCreateDatahubTaskResponse()
		taskId   string
	)
	if v, ok := d.GetOk("task_name"); ok {
		request.TaskName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_type"); ok {
		request.TaskType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "source_resource"); ok {
		datahubResource := ckafka.DatahubResource{}
		if v, ok := dMap["type"]; ok {
			datahubResource.Type = helper.String(v.(string))
		}
		if kafkaParamMap, ok := helper.InterfaceToMap(dMap, "kafka_param"); ok {
			kafkaParam := ckafka.KafkaParam{}
			if v, ok := kafkaParamMap["self_built"]; ok {
				kafkaParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := kafkaParamMap["resource"]; ok {
				kafkaParam.Resource = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["topic"]; ok {
				kafkaParam.Topic = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["offset_type"]; ok {
				kafkaParam.OffsetType = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["start_time"]; ok {
				kafkaParam.StartTime = helper.IntUint64(v.(int))
			}
			if v, ok := kafkaParamMap["resource_name"]; ok {
				kafkaParam.ResourceName = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["zone_id"]; ok {
				kafkaParam.ZoneId = helper.IntInt64(v.(int))
			}
			if v, ok := kafkaParamMap["topic_id"]; ok {
				kafkaParam.TopicId = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["partition_num"]; ok {
				kafkaParam.PartitionNum = helper.IntInt64(v.(int))
			}
			if v, ok := kafkaParamMap["enable_toleration"]; ok {
				kafkaParam.EnableToleration = helper.Bool(v.(bool))
			}
			if v, ok := kafkaParamMap["qps_limit"]; ok {
				kafkaParam.QpsLimit = helper.IntUint64(v.(int))
			}
			if v, ok := kafkaParamMap["table_mappings"]; ok {
				for _, item := range v.([]interface{}) {
					tableMappingsMap := item.(map[string]interface{})
					tableMapping := ckafka.TableMapping{}
					if v, ok := tableMappingsMap["database"]; ok {
						tableMapping.Database = helper.String(v.(string))
					}
					if v, ok := tableMappingsMap["table"]; ok {
						tableMapping.Table = helper.String(v.(string))
					}
					if v, ok := tableMappingsMap["topic"]; ok {
						tableMapping.Topic = helper.String(v.(string))
					}
					if v, ok := tableMappingsMap["topic_id"]; ok {
						tableMapping.TopicId = helper.String(v.(string))
					}
					kafkaParam.TableMappings = append(kafkaParam.TableMappings, &tableMapping)
				}
			}
			if v, ok := kafkaParamMap["use_table_mapping"]; ok {
				kafkaParam.UseTableMapping = helper.Bool(v.(bool))
			}
			if v, ok := kafkaParamMap["use_auto_create_topic"]; ok {
				kafkaParam.UseAutoCreateTopic = helper.Bool(v.(bool))
			}
			if v, ok := kafkaParamMap["compression_type"]; ok {
				kafkaParam.CompressionType = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["msg_multiple"]; ok {
				kafkaParam.MsgMultiple = helper.IntInt64(v.(int))
			}
			datahubResource.KafkaParam = &kafkaParam
		}
		if eventBusParamMap, ok := helper.InterfaceToMap(dMap, "event_bus_param"); ok {
			eventBusParam := ckafka.EventBusParam{}
			if v, ok := eventBusParamMap["type"]; ok {
				eventBusParam.Type = helper.String(v.(string))
			}
			if v, ok := eventBusParamMap["self_built"]; ok {
				eventBusParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := eventBusParamMap["resource"]; ok {
				eventBusParam.Resource = helper.String(v.(string))
			}
			if v, ok := eventBusParamMap["namespace"]; ok {
				eventBusParam.Namespace = helper.String(v.(string))
			}
			if v, ok := eventBusParamMap["function_name"]; ok {
				eventBusParam.FunctionName = helper.String(v.(string))
			}
			if v, ok := eventBusParamMap["qualifier"]; ok {
				eventBusParam.Qualifier = helper.String(v.(string))
			}
			datahubResource.EventBusParam = &eventBusParam
		}
		if mongoDBParamMap, ok := helper.InterfaceToMap(dMap, "mongo_db_param"); ok {
			mongoDBParam := ckafka.MongoDBParam{}
			if v, ok := mongoDBParamMap["database"]; ok {
				mongoDBParam.Database = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["collection"]; ok {
				mongoDBParam.Collection = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["copy_existing"]; ok {
				mongoDBParam.CopyExisting = helper.Bool(v.(bool))
			}
			if v, ok := mongoDBParamMap["resource"]; ok {
				mongoDBParam.Resource = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["ip"]; ok {
				mongoDBParam.Ip = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["port"]; ok {
				mongoDBParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := mongoDBParamMap["user_name"]; ok {
				mongoDBParam.UserName = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["password"]; ok {
				mongoDBParam.Password = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["listening_event"]; ok {
				mongoDBParam.ListeningEvent = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["read_preference"]; ok {
				mongoDBParam.ReadPreference = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["pipeline"]; ok {
				mongoDBParam.Pipeline = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["self_built"]; ok {
				mongoDBParam.SelfBuilt = helper.Bool(v.(bool))
			}
			datahubResource.MongoDBParam = &mongoDBParam
		}
		if esParamMap, ok := helper.InterfaceToMap(dMap, "es_param"); ok {
			esParam := ckafka.EsParam{}
			if v, ok := esParamMap["resource"]; ok {
				esParam.Resource = helper.String(v.(string))
			}
			if v, ok := esParamMap["port"]; ok {
				esParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := esParamMap["user_name"]; ok {
				esParam.UserName = helper.String(v.(string))
			}
			if v, ok := esParamMap["password"]; ok {
				esParam.Password = helper.String(v.(string))
			}
			if v, ok := esParamMap["self_built"]; ok {
				esParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := esParamMap["service_vip"]; ok {
				esParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := esParamMap["uniq_vpc_id"]; ok {
				esParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := esParamMap["drop_invalid_message"]; ok {
				esParam.DropInvalidMessage = helper.Bool(v.(bool))
			}
			if v, ok := esParamMap["index"]; ok {
				esParam.Index = helper.String(v.(string))
			}
			if v, ok := esParamMap["date_format"]; ok {
				esParam.DateFormat = helper.String(v.(string))
			}
			if v, ok := esParamMap["content_key"]; ok {
				esParam.ContentKey = helper.String(v.(string))
			}
			if v, ok := esParamMap["drop_invalid_json_message"]; ok {
				esParam.DropInvalidJsonMessage = helper.Bool(v.(bool))
			}
			if v, ok := esParamMap["document_id_field"]; ok {
				esParam.DocumentIdField = helper.String(v.(string))
			}
			if v, ok := esParamMap["index_type"]; ok {
				esParam.IndexType = helper.String(v.(string))
			}
			if dropClsMap, ok := helper.InterfaceToMap(esParamMap, "drop_cls"); ok {
				dropCls := ckafka.DropCls{}
				if v, ok := dropClsMap["drop_invalid_message_to_cls"]; ok {
					dropCls.DropInvalidMessageToCls = helper.Bool(v.(bool))
				}
				if v, ok := dropClsMap["drop_cls_region"]; ok {
					dropCls.DropClsRegion = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_owneruin"]; ok {
					dropCls.DropClsOwneruin = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_topic_id"]; ok {
					dropCls.DropClsTopicId = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_log_set"]; ok {
					dropCls.DropClsLogSet = helper.String(v.(string))
				}
				esParam.DropCls = &dropCls
			}
			if v, ok := esParamMap["database_primary_key"]; ok {
				esParam.DatabasePrimaryKey = helper.String(v.(string))
			}
			if dropDlqMap, ok := helper.InterfaceToMap(esParamMap, "drop_dlq"); ok {
				failureParam := ckafka.FailureParam{}
				if v, ok := dropDlqMap["type"]; ok {
					failureParam.Type = helper.String(v.(string))
				}
				if kafkaParamMap, ok := helper.InterfaceToMap(dropDlqMap, "kafka_param"); ok {
					kafkaParam := ckafka.KafkaParam{}
					if v, ok := kafkaParamMap["self_built"]; ok {
						kafkaParam.SelfBuilt = helper.Bool(v.(bool))
					}
					if v, ok := kafkaParamMap["resource"]; ok {
						kafkaParam.Resource = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["topic"]; ok {
						kafkaParam.Topic = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["offset_type"]; ok {
						kafkaParam.OffsetType = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["start_time"]; ok {
						kafkaParam.StartTime = helper.IntUint64(v.(int))
					}
					if v, ok := kafkaParamMap["resource_name"]; ok {
						kafkaParam.ResourceName = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["zone_id"]; ok {
						kafkaParam.ZoneId = helper.IntInt64(v.(int))
					}
					if v, ok := kafkaParamMap["topic_id"]; ok {
						kafkaParam.TopicId = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["partition_num"]; ok {
						kafkaParam.PartitionNum = helper.IntInt64(v.(int))
					}
					if v, ok := kafkaParamMap["enable_toleration"]; ok {
						kafkaParam.EnableToleration = helper.Bool(v.(bool))
					}
					if v, ok := kafkaParamMap["qps_limit"]; ok {
						kafkaParam.QpsLimit = helper.IntUint64(v.(int))
					}
					if v, ok := kafkaParamMap["table_mappings"]; ok {
						for _, item := range v.([]interface{}) {
							tableMappingsMap := item.(map[string]interface{})
							tableMapping := ckafka.TableMapping{}
							if v, ok := tableMappingsMap["database"]; ok {
								tableMapping.Database = helper.String(v.(string))
							}
							if v, ok := tableMappingsMap["table"]; ok {
								tableMapping.Table = helper.String(v.(string))
							}
							if v, ok := tableMappingsMap["topic"]; ok {
								tableMapping.Topic = helper.String(v.(string))
							}
							if v, ok := tableMappingsMap["topic_id"]; ok {
								tableMapping.TopicId = helper.String(v.(string))
							}
							kafkaParam.TableMappings = append(kafkaParam.TableMappings, &tableMapping)
						}
					}
					if v, ok := kafkaParamMap["use_table_mapping"]; ok {
						kafkaParam.UseTableMapping = helper.Bool(v.(bool))
					}
					if v, ok := kafkaParamMap["use_auto_create_topic"]; ok {
						kafkaParam.UseAutoCreateTopic = helper.Bool(v.(bool))
					}
					if v, ok := kafkaParamMap["compression_type"]; ok {
						kafkaParam.CompressionType = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["msg_multiple"]; ok {
						kafkaParam.MsgMultiple = helper.IntInt64(v.(int))
					}
					failureParam.KafkaParam = &kafkaParam
				}
				if v, ok := dropDlqMap["retry_interval"]; ok {
					failureParam.RetryInterval = helper.IntUint64(v.(int))
				}
				if v, ok := dropDlqMap["max_retry_attempts"]; ok {
					failureParam.MaxRetryAttempts = helper.IntUint64(v.(int))
				}
				if topicParamMap, ok := helper.InterfaceToMap(dropDlqMap, "topic_param"); ok {
					topicParam := ckafka.TopicParam{}
					if v, ok := topicParamMap["resource"]; ok {
						topicParam.Resource = helper.String(v.(string))
					}
					if v, ok := topicParamMap["offset_type"]; ok {
						topicParam.OffsetType = helper.String(v.(string))
					}
					if v, ok := topicParamMap["start_time"]; ok {
						topicParam.StartTime = helper.IntUint64(v.(int))
					}
					if v, ok := topicParamMap["topic_id"]; ok {
						topicParam.TopicId = helper.String(v.(string))
					}
					if v, ok := topicParamMap["compression_type"]; ok {
						topicParam.CompressionType = helper.String(v.(string))
					}
					if v, ok := topicParamMap["use_auto_create_topic"]; ok {
						topicParam.UseAutoCreateTopic = helper.Bool(v.(bool))
					}
					if v, ok := topicParamMap["msg_multiple"]; ok {
						topicParam.MsgMultiple = helper.IntInt64(v.(int))
					}
					failureParam.TopicParam = &topicParam
				}
				if v, ok := dropDlqMap["dlq_type"]; ok {
					failureParam.DlqType = helper.String(v.(string))
				}
				esParam.DropDlq = &failureParam
			}
			datahubResource.EsParam = &esParam
		}
		if tdwParamMap, ok := helper.InterfaceToMap(dMap, "tdw_param"); ok {
			tdwParam := ckafka.TdwParam{}
			if v, ok := tdwParamMap["bid"]; ok {
				tdwParam.Bid = helper.String(v.(string))
			}
			if v, ok := tdwParamMap["tid"]; ok {
				tdwParam.Tid = helper.String(v.(string))
			}
			if v, ok := tdwParamMap["is_domestic"]; ok {
				tdwParam.IsDomestic = helper.Bool(v.(bool))
			}
			if v, ok := tdwParamMap["tdw_host"]; ok {
				tdwParam.TdwHost = helper.String(v.(string))
			}
			if v, ok := tdwParamMap["tdw_port"]; ok {
				tdwParam.TdwPort = helper.IntInt64(v.(int))
			}
			datahubResource.TdwParam = &tdwParam
		}
		if dtsParamMap, ok := helper.InterfaceToMap(dMap, "dts_param"); ok {
			dtsParam := ckafka.DtsParam{}
			if v, ok := dtsParamMap["resource"]; ok {
				dtsParam.Resource = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["ip"]; ok {
				dtsParam.Ip = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["port"]; ok {
				dtsParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dtsParamMap["topic"]; ok {
				dtsParam.Topic = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["group_id"]; ok {
				dtsParam.GroupId = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["group_user"]; ok {
				dtsParam.GroupUser = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["group_password"]; ok {
				dtsParam.GroupPassword = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["tran_sql"]; ok {
				dtsParam.TranSql = helper.Bool(v.(bool))
			}
			datahubResource.DtsParam = &dtsParam
		}
		if clickHouseParamMap, ok := helper.InterfaceToMap(dMap, "click_house_param"); ok {
			clickHouseParam := ckafka.ClickHouseParam{}
			if v, ok := clickHouseParamMap["cluster"]; ok {
				clickHouseParam.Cluster = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["database"]; ok {
				clickHouseParam.Database = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["table"]; ok {
				clickHouseParam.Table = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["schema"]; ok {
				for _, item := range v.([]interface{}) {
					schemaMap := item.(map[string]interface{})
					clickHouseSchema := ckafka.ClickHouseSchema{}
					if v, ok := schemaMap["column_name"]; ok {
						clickHouseSchema.ColumnName = helper.String(v.(string))
					}
					if v, ok := schemaMap["json_key"]; ok {
						clickHouseSchema.JsonKey = helper.String(v.(string))
					}
					if v, ok := schemaMap["type"]; ok {
						clickHouseSchema.Type = helper.String(v.(string))
					}
					if v, ok := schemaMap["allow_null"]; ok {
						clickHouseSchema.AllowNull = helper.Bool(v.(bool))
					}
					clickHouseParam.Schema = append(clickHouseParam.Schema, &clickHouseSchema)
				}
			}
			if v, ok := clickHouseParamMap["resource"]; ok {
				clickHouseParam.Resource = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["ip"]; ok {
				clickHouseParam.Ip = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["port"]; ok {
				clickHouseParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := clickHouseParamMap["user_name"]; ok {
				clickHouseParam.UserName = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["password"]; ok {
				clickHouseParam.Password = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["service_vip"]; ok {
				clickHouseParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["uniq_vpc_id"]; ok {
				clickHouseParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["self_built"]; ok {
				clickHouseParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := clickHouseParamMap["drop_invalid_message"]; ok {
				clickHouseParam.DropInvalidMessage = helper.Bool(v.(bool))
			}
			if v, ok := clickHouseParamMap["type"]; ok {
				clickHouseParam.Type = helper.String(v.(string))
			}
			if dropClsMap, ok := helper.InterfaceToMap(clickHouseParamMap, "drop_cls"); ok {
				dropCls := ckafka.DropCls{}
				if v, ok := dropClsMap["drop_invalid_message_to_cls"]; ok {
					dropCls.DropInvalidMessageToCls = helper.Bool(v.(bool))
				}
				if v, ok := dropClsMap["drop_cls_region"]; ok {
					dropCls.DropClsRegion = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_owneruin"]; ok {
					dropCls.DropClsOwneruin = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_topic_id"]; ok {
					dropCls.DropClsTopicId = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_log_set"]; ok {
					dropCls.DropClsLogSet = helper.String(v.(string))
				}
				clickHouseParam.DropCls = &dropCls
			}
			datahubResource.ClickHouseParam = &clickHouseParam
		}
		if clsParamMap, ok := helper.InterfaceToMap(dMap, "cls_param"); ok {
			clsParam := ckafka.ClsParam{}
			if v, ok := clsParamMap["decode_json"]; ok {
				clsParam.DecodeJson = helper.Bool(v.(bool))
			}
			if v, ok := clsParamMap["resource"]; ok {
				clsParam.Resource = helper.String(v.(string))
			}
			if v, ok := clsParamMap["log_set"]; ok {
				clsParam.LogSet = helper.String(v.(string))
			}
			if v, ok := clsParamMap["content_key"]; ok {
				clsParam.ContentKey = helper.String(v.(string))
			}
			if v, ok := clsParamMap["time_field"]; ok {
				clsParam.TimeField = helper.String(v.(string))
			}
			datahubResource.ClsParam = &clsParam
		}
		if cosParamMap, ok := helper.InterfaceToMap(dMap, "cos_param"); ok {
			cosParam := ckafka.CosParam{}
			if v, ok := cosParamMap["bucket_name"]; ok {
				cosParam.BucketName = helper.String(v.(string))
			}
			if v, ok := cosParamMap["region"]; ok {
				cosParam.Region = helper.String(v.(string))
			}
			if v, ok := cosParamMap["object_key"]; ok {
				cosParam.ObjectKey = helper.String(v.(string))
			}
			if v, ok := cosParamMap["aggregate_batch_size"]; ok {
				cosParam.AggregateBatchSize = helper.IntUint64(v.(int))
			}
			if v, ok := cosParamMap["aggregate_interval"]; ok {
				cosParam.AggregateInterval = helper.IntUint64(v.(int))
			}
			if v, ok := cosParamMap["format_output_type"]; ok {
				cosParam.FormatOutputType = helper.String(v.(string))
			}
			if v, ok := cosParamMap["object_key_prefix"]; ok {
				cosParam.ObjectKeyPrefix = helper.String(v.(string))
			}
			if v, ok := cosParamMap["directory_time_format"]; ok {
				cosParam.DirectoryTimeFormat = helper.String(v.(string))
			}
			datahubResource.CosParam = &cosParam
		}
		if mySQLParamMap, ok := helper.InterfaceToMap(dMap, "my_sql_param"); ok {
			mySQLParam := ckafka.MySQLParam{}
			if v, ok := mySQLParamMap["database"]; ok {
				mySQLParam.Database = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["table"]; ok {
				mySQLParam.Table = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["resource"]; ok {
				mySQLParam.Resource = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["snapshot_mode"]; ok {
				mySQLParam.SnapshotMode = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["ddl_topic"]; ok {
				mySQLParam.DdlTopic = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_source_monitor_mode"]; ok {
				mySQLParam.DataSourceMonitorMode = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_source_monitor_resource"]; ok {
				mySQLParam.DataSourceMonitorResource = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_source_increment_mode"]; ok {
				mySQLParam.DataSourceIncrementMode = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_source_increment_column"]; ok {
				mySQLParam.DataSourceIncrementColumn = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_source_start_from"]; ok {
				mySQLParam.DataSourceStartFrom = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_target_insert_mode"]; ok {
				mySQLParam.DataTargetInsertMode = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_target_primary_key_field"]; ok {
				mySQLParam.DataTargetPrimaryKeyField = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_target_record_mapping"]; ok {
				for _, item := range v.([]interface{}) {
					dataTargetRecordMappingMap := item.(map[string]interface{})
					recordMapping := ckafka.RecordMapping{}
					if v, ok := dataTargetRecordMappingMap["json_key"]; ok {
						recordMapping.JsonKey = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["type"]; ok {
						recordMapping.Type = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["allow_null"]; ok {
						recordMapping.AllowNull = helper.Bool(v.(bool))
					}
					if v, ok := dataTargetRecordMappingMap["column_name"]; ok {
						recordMapping.ColumnName = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["extra_info"]; ok {
						recordMapping.ExtraInfo = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["column_size"]; ok {
						recordMapping.ColumnSize = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["decimal_digits"]; ok {
						recordMapping.DecimalDigits = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["auto_increment"]; ok {
						recordMapping.AutoIncrement = helper.Bool(v.(bool))
					}
					if v, ok := dataTargetRecordMappingMap["default_value"]; ok {
						recordMapping.DefaultValue = helper.String(v.(string))
					}
					mySQLParam.DataTargetRecordMapping = append(mySQLParam.DataTargetRecordMapping, &recordMapping)
				}
			}
			if v, ok := mySQLParamMap["topic_regex"]; ok {
				mySQLParam.TopicRegex = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["topic_replacement"]; ok {
				mySQLParam.TopicReplacement = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["key_columns"]; ok {
				mySQLParam.KeyColumns = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["drop_invalid_message"]; ok {
				mySQLParam.DropInvalidMessage = helper.Bool(v.(bool))
			}
			if dropClsMap, ok := helper.InterfaceToMap(mySQLParamMap, "drop_cls"); ok {
				dropCls := ckafka.DropCls{}
				if v, ok := dropClsMap["drop_invalid_message_to_cls"]; ok {
					dropCls.DropInvalidMessageToCls = helper.Bool(v.(bool))
				}
				if v, ok := dropClsMap["drop_cls_region"]; ok {
					dropCls.DropClsRegion = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_owneruin"]; ok {
					dropCls.DropClsOwneruin = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_topic_id"]; ok {
					dropCls.DropClsTopicId = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_log_set"]; ok {
					dropCls.DropClsLogSet = helper.String(v.(string))
				}
				mySQLParam.DropCls = &dropCls
			}
			if v, ok := mySQLParamMap["output_format"]; ok {
				mySQLParam.OutputFormat = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["is_table_prefix"]; ok {
				mySQLParam.IsTablePrefix = helper.Bool(v.(bool))
			}
			if v, ok := mySQLParamMap["include_content_changes"]; ok {
				mySQLParam.IncludeContentChanges = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["include_query"]; ok {
				mySQLParam.IncludeQuery = helper.Bool(v.(bool))
			}
			if v, ok := mySQLParamMap["record_with_schema"]; ok {
				mySQLParam.RecordWithSchema = helper.Bool(v.(bool))
			}
			if v, ok := mySQLParamMap["signal_database"]; ok {
				mySQLParam.SignalDatabase = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["is_table_regular"]; ok {
				mySQLParam.IsTableRegular = helper.Bool(v.(bool))
			}
			datahubResource.MySQLParam = &mySQLParam
		}
		if postgreSQLParamMap, ok := helper.InterfaceToMap(dMap, "postgre_sql_param"); ok {
			postgreSQLParam := ckafka.PostgreSQLParam{}
			if v, ok := postgreSQLParamMap["database"]; ok {
				postgreSQLParam.Database = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["table"]; ok {
				postgreSQLParam.Table = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["resource"]; ok {
				postgreSQLParam.Resource = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["plugin_name"]; ok {
				postgreSQLParam.PluginName = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["snapshot_mode"]; ok {
				postgreSQLParam.SnapshotMode = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["data_format"]; ok {
				postgreSQLParam.DataFormat = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["data_target_insert_mode"]; ok {
				postgreSQLParam.DataTargetInsertMode = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["data_target_primary_key_field"]; ok {
				postgreSQLParam.DataTargetPrimaryKeyField = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["data_target_record_mapping"]; ok {
				for _, item := range v.([]interface{}) {
					dataTargetRecordMappingMap := item.(map[string]interface{})
					recordMapping := ckafka.RecordMapping{}
					if v, ok := dataTargetRecordMappingMap["json_key"]; ok {
						recordMapping.JsonKey = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["type"]; ok {
						recordMapping.Type = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["allow_null"]; ok {
						recordMapping.AllowNull = helper.Bool(v.(bool))
					}
					if v, ok := dataTargetRecordMappingMap["column_name"]; ok {
						recordMapping.ColumnName = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["extra_info"]; ok {
						recordMapping.ExtraInfo = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["column_size"]; ok {
						recordMapping.ColumnSize = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["decimal_digits"]; ok {
						recordMapping.DecimalDigits = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["auto_increment"]; ok {
						recordMapping.AutoIncrement = helper.Bool(v.(bool))
					}
					if v, ok := dataTargetRecordMappingMap["default_value"]; ok {
						recordMapping.DefaultValue = helper.String(v.(string))
					}
					postgreSQLParam.DataTargetRecordMapping = append(postgreSQLParam.DataTargetRecordMapping, &recordMapping)
				}
			}
			if v, ok := postgreSQLParamMap["drop_invalid_message"]; ok {
				postgreSQLParam.DropInvalidMessage = helper.Bool(v.(bool))
			}
			if v, ok := postgreSQLParamMap["is_table_regular"]; ok {
				postgreSQLParam.IsTableRegular = helper.Bool(v.(bool))
			}
			if v, ok := postgreSQLParamMap["key_columns"]; ok {
				postgreSQLParam.KeyColumns = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["record_with_schema"]; ok {
				postgreSQLParam.RecordWithSchema = helper.Bool(v.(bool))
			}
			datahubResource.PostgreSQLParam = &postgreSQLParam
		}
		if topicParamMap, ok := helper.InterfaceToMap(dMap, "topic_param"); ok {
			topicParam := ckafka.TopicParam{}
			if v, ok := topicParamMap["resource"]; ok {
				topicParam.Resource = helper.String(v.(string))
			}
			if v, ok := topicParamMap["offset_type"]; ok {
				topicParam.OffsetType = helper.String(v.(string))
			}
			if v, ok := topicParamMap["start_time"]; ok {
				topicParam.StartTime = helper.IntUint64(v.(int))
			}
			if v, ok := topicParamMap["topic_id"]; ok {
				topicParam.TopicId = helper.String(v.(string))
			}
			if v, ok := topicParamMap["compression_type"]; ok {
				topicParam.CompressionType = helper.String(v.(string))
			}
			if v, ok := topicParamMap["use_auto_create_topic"]; ok {
				topicParam.UseAutoCreateTopic = helper.Bool(v.(bool))
			}
			if v, ok := topicParamMap["msg_multiple"]; ok {
				topicParam.MsgMultiple = helper.IntInt64(v.(int))
			}
			datahubResource.TopicParam = &topicParam
		}
		if mariaDBParamMap, ok := helper.InterfaceToMap(dMap, "maria_db_param"); ok {
			mariaDBParam := ckafka.MariaDBParam{}
			if v, ok := mariaDBParamMap["database"]; ok {
				mariaDBParam.Database = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["table"]; ok {
				mariaDBParam.Table = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["resource"]; ok {
				mariaDBParam.Resource = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["snapshot_mode"]; ok {
				mariaDBParam.SnapshotMode = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["key_columns"]; ok {
				mariaDBParam.KeyColumns = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["is_table_prefix"]; ok {
				mariaDBParam.IsTablePrefix = helper.Bool(v.(bool))
			}
			if v, ok := mariaDBParamMap["output_format"]; ok {
				mariaDBParam.OutputFormat = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["include_content_changes"]; ok {
				mariaDBParam.IncludeContentChanges = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["include_query"]; ok {
				mariaDBParam.IncludeQuery = helper.Bool(v.(bool))
			}
			if v, ok := mariaDBParamMap["record_with_schema"]; ok {
				mariaDBParam.RecordWithSchema = helper.Bool(v.(bool))
			}
			datahubResource.MariaDBParam = &mariaDBParam
		}
		if sQLServerParamMap, ok := helper.InterfaceToMap(dMap, "sql_server_param"); ok {
			sQLServerParam := ckafka.SQLServerParam{}
			if v, ok := sQLServerParamMap["database"]; ok {
				sQLServerParam.Database = helper.String(v.(string))
			}
			if v, ok := sQLServerParamMap["table"]; ok {
				sQLServerParam.Table = helper.String(v.(string))
			}
			if v, ok := sQLServerParamMap["resource"]; ok {
				sQLServerParam.Resource = helper.String(v.(string))
			}
			if v, ok := sQLServerParamMap["snapshot_mode"]; ok {
				sQLServerParam.SnapshotMode = helper.String(v.(string))
			}
			datahubResource.SQLServerParam = &sQLServerParam
		}
		if ctsdbParamMap, ok := helper.InterfaceToMap(dMap, "ctsdb_param"); ok {
			ctsdbParam := ckafka.CtsdbParam{}
			if v, ok := ctsdbParamMap["resource"]; ok {
				ctsdbParam.Resource = helper.String(v.(string))
			}
			if v, ok := ctsdbParamMap["ctsdb_metric"]; ok {
				ctsdbParam.CtsdbMetric = helper.String(v.(string))
			}
			datahubResource.CtsdbParam = &ctsdbParam
		}
		if scfParamMap, ok := helper.InterfaceToMap(dMap, "scf_param"); ok {
			scfParam := ckafka.ScfParam{}
			if v, ok := scfParamMap["function_name"]; ok {
				scfParam.FunctionName = helper.String(v.(string))
			}
			if v, ok := scfParamMap["namespace"]; ok {
				scfParam.Namespace = helper.String(v.(string))
			}
			if v, ok := scfParamMap["qualifier"]; ok {
				scfParam.Qualifier = helper.String(v.(string))
			}
			if v, ok := scfParamMap["batch_size"]; ok {
				scfParam.BatchSize = helper.IntInt64(v.(int))
			}
			if v, ok := scfParamMap["max_retries"]; ok {
				scfParam.MaxRetries = helper.IntInt64(v.(int))
			}
			datahubResource.ScfParam = &scfParam
		}
		request.SourceResource = &datahubResource
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "target_resource"); ok {
		datahubResource := ckafka.DatahubResource{}
		if v, ok := dMap["type"]; ok {
			datahubResource.Type = helper.String(v.(string))
		}
		if kafkaParamMap, ok := helper.InterfaceToMap(dMap, "kafka_param"); ok {
			kafkaParam := ckafka.KafkaParam{}
			if v, ok := kafkaParamMap["self_built"]; ok {
				kafkaParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := kafkaParamMap["resource"]; ok {
				kafkaParam.Resource = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["topic"]; ok {
				kafkaParam.Topic = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["offset_type"]; ok {
				kafkaParam.OffsetType = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["start_time"]; ok {
				kafkaParam.StartTime = helper.IntUint64(v.(int))
			}
			if v, ok := kafkaParamMap["resource_name"]; ok {
				kafkaParam.ResourceName = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["zone_id"]; ok {
				kafkaParam.ZoneId = helper.IntInt64(v.(int))
			}
			if v, ok := kafkaParamMap["topic_id"]; ok {
				kafkaParam.TopicId = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["partition_num"]; ok {
				kafkaParam.PartitionNum = helper.IntInt64(v.(int))
			}
			if v, ok := kafkaParamMap["enable_toleration"]; ok {
				kafkaParam.EnableToleration = helper.Bool(v.(bool))
			}
			if v, ok := kafkaParamMap["qps_limit"]; ok {
				kafkaParam.QpsLimit = helper.IntUint64(v.(int))
			}
			if v, ok := kafkaParamMap["table_mappings"]; ok {
				for _, item := range v.([]interface{}) {
					tableMappingsMap := item.(map[string]interface{})
					tableMapping := ckafka.TableMapping{}
					if v, ok := tableMappingsMap["database"]; ok {
						tableMapping.Database = helper.String(v.(string))
					}
					if v, ok := tableMappingsMap["table"]; ok {
						tableMapping.Table = helper.String(v.(string))
					}
					if v, ok := tableMappingsMap["topic"]; ok {
						tableMapping.Topic = helper.String(v.(string))
					}
					if v, ok := tableMappingsMap["topic_id"]; ok {
						tableMapping.TopicId = helper.String(v.(string))
					}
					kafkaParam.TableMappings = append(kafkaParam.TableMappings, &tableMapping)
				}
			}
			if v, ok := kafkaParamMap["use_table_mapping"]; ok {
				kafkaParam.UseTableMapping = helper.Bool(v.(bool))
			}
			if v, ok := kafkaParamMap["use_auto_create_topic"]; ok {
				kafkaParam.UseAutoCreateTopic = helper.Bool(v.(bool))
			}
			if v, ok := kafkaParamMap["compression_type"]; ok {
				kafkaParam.CompressionType = helper.String(v.(string))
			}
			if v, ok := kafkaParamMap["msg_multiple"]; ok {
				kafkaParam.MsgMultiple = helper.IntInt64(v.(int))
			}
			datahubResource.KafkaParam = &kafkaParam
		}
		if eventBusParamMap, ok := helper.InterfaceToMap(dMap, "event_bus_param"); ok {
			eventBusParam := ckafka.EventBusParam{}
			if v, ok := eventBusParamMap["type"]; ok {
				eventBusParam.Type = helper.String(v.(string))
			}
			if v, ok := eventBusParamMap["self_built"]; ok {
				eventBusParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := eventBusParamMap["resource"]; ok {
				eventBusParam.Resource = helper.String(v.(string))
			}
			if v, ok := eventBusParamMap["namespace"]; ok {
				eventBusParam.Namespace = helper.String(v.(string))
			}
			if v, ok := eventBusParamMap["function_name"]; ok {
				eventBusParam.FunctionName = helper.String(v.(string))
			}
			if v, ok := eventBusParamMap["qualifier"]; ok {
				eventBusParam.Qualifier = helper.String(v.(string))
			}
			datahubResource.EventBusParam = &eventBusParam
		}
		if mongoDBParamMap, ok := helper.InterfaceToMap(dMap, "mongo_db_param"); ok {
			mongoDBParam := ckafka.MongoDBParam{}
			if v, ok := mongoDBParamMap["database"]; ok {
				mongoDBParam.Database = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["collection"]; ok {
				mongoDBParam.Collection = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["copy_existing"]; ok {
				mongoDBParam.CopyExisting = helper.Bool(v.(bool))
			}
			if v, ok := mongoDBParamMap["resource"]; ok {
				mongoDBParam.Resource = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["ip"]; ok {
				mongoDBParam.Ip = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["port"]; ok {
				mongoDBParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := mongoDBParamMap["user_name"]; ok {
				mongoDBParam.UserName = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["password"]; ok {
				mongoDBParam.Password = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["listening_event"]; ok {
				mongoDBParam.ListeningEvent = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["read_preference"]; ok {
				mongoDBParam.ReadPreference = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["pipeline"]; ok {
				mongoDBParam.Pipeline = helper.String(v.(string))
			}
			if v, ok := mongoDBParamMap["self_built"]; ok {
				mongoDBParam.SelfBuilt = helper.Bool(v.(bool))
			}
			datahubResource.MongoDBParam = &mongoDBParam
		}
		if esParamMap, ok := helper.InterfaceToMap(dMap, "es_param"); ok {
			esParam := ckafka.EsParam{}
			if v, ok := esParamMap["resource"]; ok {
				esParam.Resource = helper.String(v.(string))
			}
			if v, ok := esParamMap["port"]; ok {
				esParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := esParamMap["user_name"]; ok {
				esParam.UserName = helper.String(v.(string))
			}
			if v, ok := esParamMap["password"]; ok {
				esParam.Password = helper.String(v.(string))
			}
			if v, ok := esParamMap["self_built"]; ok {
				esParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := esParamMap["service_vip"]; ok {
				esParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := esParamMap["uniq_vpc_id"]; ok {
				esParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := esParamMap["drop_invalid_message"]; ok {
				esParam.DropInvalidMessage = helper.Bool(v.(bool))
			}
			if v, ok := esParamMap["index"]; ok {
				esParam.Index = helper.String(v.(string))
			}
			if v, ok := esParamMap["date_format"]; ok {
				esParam.DateFormat = helper.String(v.(string))
			}
			if v, ok := esParamMap["content_key"]; ok {
				esParam.ContentKey = helper.String(v.(string))
			}
			if v, ok := esParamMap["drop_invalid_json_message"]; ok {
				esParam.DropInvalidJsonMessage = helper.Bool(v.(bool))
			}
			if v, ok := esParamMap["document_id_field"]; ok {
				esParam.DocumentIdField = helper.String(v.(string))
			}
			if v, ok := esParamMap["index_type"]; ok {
				esParam.IndexType = helper.String(v.(string))
			}
			if dropClsMap, ok := helper.InterfaceToMap(esParamMap, "drop_cls"); ok {
				dropCls := ckafka.DropCls{}
				if v, ok := dropClsMap["drop_invalid_message_to_cls"]; ok {
					dropCls.DropInvalidMessageToCls = helper.Bool(v.(bool))
				}
				if v, ok := dropClsMap["drop_cls_region"]; ok {
					dropCls.DropClsRegion = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_owneruin"]; ok {
					dropCls.DropClsOwneruin = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_topic_id"]; ok {
					dropCls.DropClsTopicId = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_log_set"]; ok {
					dropCls.DropClsLogSet = helper.String(v.(string))
				}
				esParam.DropCls = &dropCls
			}
			if v, ok := esParamMap["database_primary_key"]; ok {
				esParam.DatabasePrimaryKey = helper.String(v.(string))
			}
			if dropDlqMap, ok := helper.InterfaceToMap(esParamMap, "drop_dlq"); ok {
				failureParam := ckafka.FailureParam{}
				if v, ok := dropDlqMap["type"]; ok {
					failureParam.Type = helper.String(v.(string))
				}
				if kafkaParamMap, ok := helper.InterfaceToMap(dropDlqMap, "kafka_param"); ok {
					kafkaParam := ckafka.KafkaParam{}
					if v, ok := kafkaParamMap["self_built"]; ok {
						kafkaParam.SelfBuilt = helper.Bool(v.(bool))
					}
					if v, ok := kafkaParamMap["resource"]; ok {
						kafkaParam.Resource = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["topic"]; ok {
						kafkaParam.Topic = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["offset_type"]; ok {
						kafkaParam.OffsetType = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["start_time"]; ok {
						kafkaParam.StartTime = helper.IntUint64(v.(int))
					}
					if v, ok := kafkaParamMap["resource_name"]; ok {
						kafkaParam.ResourceName = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["zone_id"]; ok {
						kafkaParam.ZoneId = helper.IntInt64(v.(int))
					}
					if v, ok := kafkaParamMap["topic_id"]; ok {
						kafkaParam.TopicId = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["partition_num"]; ok {
						kafkaParam.PartitionNum = helper.IntInt64(v.(int))
					}
					if v, ok := kafkaParamMap["enable_toleration"]; ok {
						kafkaParam.EnableToleration = helper.Bool(v.(bool))
					}
					if v, ok := kafkaParamMap["qps_limit"]; ok {
						kafkaParam.QpsLimit = helper.IntUint64(v.(int))
					}
					if v, ok := kafkaParamMap["table_mappings"]; ok {
						for _, item := range v.([]interface{}) {
							tableMappingsMap := item.(map[string]interface{})
							tableMapping := ckafka.TableMapping{}
							if v, ok := tableMappingsMap["database"]; ok {
								tableMapping.Database = helper.String(v.(string))
							}
							if v, ok := tableMappingsMap["table"]; ok {
								tableMapping.Table = helper.String(v.(string))
							}
							if v, ok := tableMappingsMap["topic"]; ok {
								tableMapping.Topic = helper.String(v.(string))
							}
							if v, ok := tableMappingsMap["topic_id"]; ok {
								tableMapping.TopicId = helper.String(v.(string))
							}
							kafkaParam.TableMappings = append(kafkaParam.TableMappings, &tableMapping)
						}
					}
					if v, ok := kafkaParamMap["use_table_mapping"]; ok {
						kafkaParam.UseTableMapping = helper.Bool(v.(bool))
					}
					if v, ok := kafkaParamMap["use_auto_create_topic"]; ok {
						kafkaParam.UseAutoCreateTopic = helper.Bool(v.(bool))
					}
					if v, ok := kafkaParamMap["compression_type"]; ok {
						kafkaParam.CompressionType = helper.String(v.(string))
					}
					if v, ok := kafkaParamMap["msg_multiple"]; ok {
						kafkaParam.MsgMultiple = helper.IntInt64(v.(int))
					}
					failureParam.KafkaParam = &kafkaParam
				}
				if v, ok := dropDlqMap["retry_interval"]; ok {
					failureParam.RetryInterval = helper.IntUint64(v.(int))
				}
				if v, ok := dropDlqMap["max_retry_attempts"]; ok {
					failureParam.MaxRetryAttempts = helper.IntUint64(v.(int))
				}
				if topicParamMap, ok := helper.InterfaceToMap(dropDlqMap, "topic_param"); ok {
					topicParam := ckafka.TopicParam{}
					if v, ok := topicParamMap["resource"]; ok {
						topicParam.Resource = helper.String(v.(string))
					}
					if v, ok := topicParamMap["offset_type"]; ok {
						topicParam.OffsetType = helper.String(v.(string))
					}
					if v, ok := topicParamMap["start_time"]; ok {
						topicParam.StartTime = helper.IntUint64(v.(int))
					}
					if v, ok := topicParamMap["topic_id"]; ok {
						topicParam.TopicId = helper.String(v.(string))
					}
					if v, ok := topicParamMap["compression_type"]; ok {
						topicParam.CompressionType = helper.String(v.(string))
					}
					if v, ok := topicParamMap["use_auto_create_topic"]; ok {
						topicParam.UseAutoCreateTopic = helper.Bool(v.(bool))
					}
					if v, ok := topicParamMap["msg_multiple"]; ok {
						topicParam.MsgMultiple = helper.IntInt64(v.(int))
					}
					failureParam.TopicParam = &topicParam
				}
				if v, ok := dropDlqMap["dlq_type"]; ok {
					failureParam.DlqType = helper.String(v.(string))
				}
				esParam.DropDlq = &failureParam
			}
			datahubResource.EsParam = &esParam
		}
		if tdwParamMap, ok := helper.InterfaceToMap(dMap, "tdw_param"); ok {
			tdwParam := ckafka.TdwParam{}
			if v, ok := tdwParamMap["bid"]; ok {
				tdwParam.Bid = helper.String(v.(string))
			}
			if v, ok := tdwParamMap["tid"]; ok {
				tdwParam.Tid = helper.String(v.(string))
			}
			if v, ok := tdwParamMap["is_domestic"]; ok {
				tdwParam.IsDomestic = helper.Bool(v.(bool))
			}
			if v, ok := tdwParamMap["tdw_host"]; ok {
				tdwParam.TdwHost = helper.String(v.(string))
			}
			if v, ok := tdwParamMap["tdw_port"]; ok {
				tdwParam.TdwPort = helper.IntInt64(v.(int))
			}
			datahubResource.TdwParam = &tdwParam
		}
		if dtsParamMap, ok := helper.InterfaceToMap(dMap, "dts_param"); ok {
			dtsParam := ckafka.DtsParam{}
			if v, ok := dtsParamMap["resource"]; ok {
				dtsParam.Resource = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["ip"]; ok {
				dtsParam.Ip = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["port"]; ok {
				dtsParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dtsParamMap["topic"]; ok {
				dtsParam.Topic = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["group_id"]; ok {
				dtsParam.GroupId = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["group_user"]; ok {
				dtsParam.GroupUser = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["group_password"]; ok {
				dtsParam.GroupPassword = helper.String(v.(string))
			}
			if v, ok := dtsParamMap["tran_sql"]; ok {
				dtsParam.TranSql = helper.Bool(v.(bool))
			}
			datahubResource.DtsParam = &dtsParam
		}
		if clickHouseParamMap, ok := helper.InterfaceToMap(dMap, "click_house_param"); ok {
			clickHouseParam := ckafka.ClickHouseParam{}
			if v, ok := clickHouseParamMap["cluster"]; ok {
				clickHouseParam.Cluster = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["database"]; ok {
				clickHouseParam.Database = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["table"]; ok {
				clickHouseParam.Table = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["schema"]; ok {
				for _, item := range v.([]interface{}) {
					schemaMap := item.(map[string]interface{})
					clickHouseSchema := ckafka.ClickHouseSchema{}
					if v, ok := schemaMap["column_name"]; ok {
						clickHouseSchema.ColumnName = helper.String(v.(string))
					}
					if v, ok := schemaMap["json_key"]; ok {
						clickHouseSchema.JsonKey = helper.String(v.(string))
					}
					if v, ok := schemaMap["type"]; ok {
						clickHouseSchema.Type = helper.String(v.(string))
					}
					if v, ok := schemaMap["allow_null"]; ok {
						clickHouseSchema.AllowNull = helper.Bool(v.(bool))
					}
					clickHouseParam.Schema = append(clickHouseParam.Schema, &clickHouseSchema)
				}
			}
			if v, ok := clickHouseParamMap["resource"]; ok {
				clickHouseParam.Resource = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["ip"]; ok {
				clickHouseParam.Ip = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["port"]; ok {
				clickHouseParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := clickHouseParamMap["user_name"]; ok {
				clickHouseParam.UserName = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["password"]; ok {
				clickHouseParam.Password = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["service_vip"]; ok {
				clickHouseParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["uniq_vpc_id"]; ok {
				clickHouseParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := clickHouseParamMap["self_built"]; ok {
				clickHouseParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := clickHouseParamMap["drop_invalid_message"]; ok {
				clickHouseParam.DropInvalidMessage = helper.Bool(v.(bool))
			}
			if v, ok := clickHouseParamMap["type"]; ok {
				clickHouseParam.Type = helper.String(v.(string))
			}
			if dropClsMap, ok := helper.InterfaceToMap(clickHouseParamMap, "drop_cls"); ok {
				dropCls := ckafka.DropCls{}
				if v, ok := dropClsMap["drop_invalid_message_to_cls"]; ok {
					dropCls.DropInvalidMessageToCls = helper.Bool(v.(bool))
				}
				if v, ok := dropClsMap["drop_cls_region"]; ok {
					dropCls.DropClsRegion = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_owneruin"]; ok {
					dropCls.DropClsOwneruin = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_topic_id"]; ok {
					dropCls.DropClsTopicId = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_log_set"]; ok {
					dropCls.DropClsLogSet = helper.String(v.(string))
				}
				clickHouseParam.DropCls = &dropCls
			}
			datahubResource.ClickHouseParam = &clickHouseParam
		}
		if clsParamMap, ok := helper.InterfaceToMap(dMap, "cls_param"); ok {
			clsParam := ckafka.ClsParam{}
			if v, ok := clsParamMap["decode_json"]; ok {
				clsParam.DecodeJson = helper.Bool(v.(bool))
			}
			if v, ok := clsParamMap["resource"]; ok {
				clsParam.Resource = helper.String(v.(string))
			}
			if v, ok := clsParamMap["log_set"]; ok {
				clsParam.LogSet = helper.String(v.(string))
			}
			if v, ok := clsParamMap["content_key"]; ok {
				clsParam.ContentKey = helper.String(v.(string))
			}
			if v, ok := clsParamMap["time_field"]; ok {
				clsParam.TimeField = helper.String(v.(string))
			}
			datahubResource.ClsParam = &clsParam
		}
		if cosParamMap, ok := helper.InterfaceToMap(dMap, "cos_param"); ok {
			cosParam := ckafka.CosParam{}
			if v, ok := cosParamMap["bucket_name"]; ok {
				cosParam.BucketName = helper.String(v.(string))
			}
			if v, ok := cosParamMap["region"]; ok {
				cosParam.Region = helper.String(v.(string))
			}
			if v, ok := cosParamMap["object_key"]; ok {
				cosParam.ObjectKey = helper.String(v.(string))
			}
			if v, ok := cosParamMap["aggregate_batch_size"]; ok {
				cosParam.AggregateBatchSize = helper.IntUint64(v.(int))
			}
			if v, ok := cosParamMap["aggregate_interval"]; ok {
				cosParam.AggregateInterval = helper.IntUint64(v.(int))
			}
			if v, ok := cosParamMap["format_output_type"]; ok {
				cosParam.FormatOutputType = helper.String(v.(string))
			}
			if v, ok := cosParamMap["object_key_prefix"]; ok {
				cosParam.ObjectKeyPrefix = helper.String(v.(string))
			}
			if v, ok := cosParamMap["directory_time_format"]; ok {
				cosParam.DirectoryTimeFormat = helper.String(v.(string))
			}
			datahubResource.CosParam = &cosParam
		}
		if mySQLParamMap, ok := helper.InterfaceToMap(dMap, "my_sql_param"); ok {
			mySQLParam := ckafka.MySQLParam{}
			if v, ok := mySQLParamMap["database"]; ok {
				mySQLParam.Database = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["table"]; ok {
				mySQLParam.Table = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["resource"]; ok {
				mySQLParam.Resource = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["snapshot_mode"]; ok {
				mySQLParam.SnapshotMode = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["ddl_topic"]; ok {
				mySQLParam.DdlTopic = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_source_monitor_mode"]; ok {
				mySQLParam.DataSourceMonitorMode = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_source_monitor_resource"]; ok {
				mySQLParam.DataSourceMonitorResource = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_source_increment_mode"]; ok {
				mySQLParam.DataSourceIncrementMode = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_source_increment_column"]; ok {
				mySQLParam.DataSourceIncrementColumn = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_source_start_from"]; ok {
				mySQLParam.DataSourceStartFrom = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_target_insert_mode"]; ok {
				mySQLParam.DataTargetInsertMode = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_target_primary_key_field"]; ok {
				mySQLParam.DataTargetPrimaryKeyField = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["data_target_record_mapping"]; ok {
				for _, item := range v.([]interface{}) {
					dataTargetRecordMappingMap := item.(map[string]interface{})
					recordMapping := ckafka.RecordMapping{}
					if v, ok := dataTargetRecordMappingMap["json_key"]; ok {
						recordMapping.JsonKey = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["type"]; ok {
						recordMapping.Type = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["allow_null"]; ok {
						recordMapping.AllowNull = helper.Bool(v.(bool))
					}
					if v, ok := dataTargetRecordMappingMap["column_name"]; ok {
						recordMapping.ColumnName = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["extra_info"]; ok {
						recordMapping.ExtraInfo = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["column_size"]; ok {
						recordMapping.ColumnSize = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["decimal_digits"]; ok {
						recordMapping.DecimalDigits = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["auto_increment"]; ok {
						recordMapping.AutoIncrement = helper.Bool(v.(bool))
					}
					if v, ok := dataTargetRecordMappingMap["default_value"]; ok {
						recordMapping.DefaultValue = helper.String(v.(string))
					}
					mySQLParam.DataTargetRecordMapping = append(mySQLParam.DataTargetRecordMapping, &recordMapping)
				}
			}
			if v, ok := mySQLParamMap["topic_regex"]; ok {
				mySQLParam.TopicRegex = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["topic_replacement"]; ok {
				mySQLParam.TopicReplacement = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["key_columns"]; ok {
				mySQLParam.KeyColumns = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["drop_invalid_message"]; ok {
				mySQLParam.DropInvalidMessage = helper.Bool(v.(bool))
			}
			if dropClsMap, ok := helper.InterfaceToMap(mySQLParamMap, "drop_cls"); ok {
				dropCls := ckafka.DropCls{}
				if v, ok := dropClsMap["drop_invalid_message_to_cls"]; ok {
					dropCls.DropInvalidMessageToCls = helper.Bool(v.(bool))
				}
				if v, ok := dropClsMap["drop_cls_region"]; ok {
					dropCls.DropClsRegion = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_owneruin"]; ok {
					dropCls.DropClsOwneruin = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_topic_id"]; ok {
					dropCls.DropClsTopicId = helper.String(v.(string))
				}
				if v, ok := dropClsMap["drop_cls_log_set"]; ok {
					dropCls.DropClsLogSet = helper.String(v.(string))
				}
				mySQLParam.DropCls = &dropCls
			}
			if v, ok := mySQLParamMap["output_format"]; ok {
				mySQLParam.OutputFormat = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["is_table_prefix"]; ok {
				mySQLParam.IsTablePrefix = helper.Bool(v.(bool))
			}
			if v, ok := mySQLParamMap["include_content_changes"]; ok {
				mySQLParam.IncludeContentChanges = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["include_query"]; ok {
				mySQLParam.IncludeQuery = helper.Bool(v.(bool))
			}
			if v, ok := mySQLParamMap["record_with_schema"]; ok {
				mySQLParam.RecordWithSchema = helper.Bool(v.(bool))
			}
			if v, ok := mySQLParamMap["signal_database"]; ok {
				mySQLParam.SignalDatabase = helper.String(v.(string))
			}
			if v, ok := mySQLParamMap["is_table_regular"]; ok {
				mySQLParam.IsTableRegular = helper.Bool(v.(bool))
			}
			datahubResource.MySQLParam = &mySQLParam
		}
		if postgreSQLParamMap, ok := helper.InterfaceToMap(dMap, "postgre_sql_param"); ok {
			postgreSQLParam := ckafka.PostgreSQLParam{}
			if v, ok := postgreSQLParamMap["database"]; ok {
				postgreSQLParam.Database = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["table"]; ok {
				postgreSQLParam.Table = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["resource"]; ok {
				postgreSQLParam.Resource = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["plugin_name"]; ok {
				postgreSQLParam.PluginName = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["snapshot_mode"]; ok {
				postgreSQLParam.SnapshotMode = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["data_format"]; ok {
				postgreSQLParam.DataFormat = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["data_target_insert_mode"]; ok {
				postgreSQLParam.DataTargetInsertMode = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["data_target_primary_key_field"]; ok {
				postgreSQLParam.DataTargetPrimaryKeyField = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["data_target_record_mapping"]; ok {
				for _, item := range v.([]interface{}) {
					dataTargetRecordMappingMap := item.(map[string]interface{})
					recordMapping := ckafka.RecordMapping{}
					if v, ok := dataTargetRecordMappingMap["json_key"]; ok {
						recordMapping.JsonKey = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["type"]; ok {
						recordMapping.Type = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["allow_null"]; ok {
						recordMapping.AllowNull = helper.Bool(v.(bool))
					}
					if v, ok := dataTargetRecordMappingMap["column_name"]; ok {
						recordMapping.ColumnName = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["extra_info"]; ok {
						recordMapping.ExtraInfo = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["column_size"]; ok {
						recordMapping.ColumnSize = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["decimal_digits"]; ok {
						recordMapping.DecimalDigits = helper.String(v.(string))
					}
					if v, ok := dataTargetRecordMappingMap["auto_increment"]; ok {
						recordMapping.AutoIncrement = helper.Bool(v.(bool))
					}
					if v, ok := dataTargetRecordMappingMap["default_value"]; ok {
						recordMapping.DefaultValue = helper.String(v.(string))
					}
					postgreSQLParam.DataTargetRecordMapping = append(postgreSQLParam.DataTargetRecordMapping, &recordMapping)
				}
			}
			if v, ok := postgreSQLParamMap["drop_invalid_message"]; ok {
				postgreSQLParam.DropInvalidMessage = helper.Bool(v.(bool))
			}
			if v, ok := postgreSQLParamMap["is_table_regular"]; ok {
				postgreSQLParam.IsTableRegular = helper.Bool(v.(bool))
			}
			if v, ok := postgreSQLParamMap["key_columns"]; ok {
				postgreSQLParam.KeyColumns = helper.String(v.(string))
			}
			if v, ok := postgreSQLParamMap["record_with_schema"]; ok {
				postgreSQLParam.RecordWithSchema = helper.Bool(v.(bool))
			}
			datahubResource.PostgreSQLParam = &postgreSQLParam
		}
		if topicParamMap, ok := helper.InterfaceToMap(dMap, "topic_param"); ok {
			topicParam := ckafka.TopicParam{}
			if v, ok := topicParamMap["resource"]; ok {
				topicParam.Resource = helper.String(v.(string))
			}
			if v, ok := topicParamMap["offset_type"]; ok {
				topicParam.OffsetType = helper.String(v.(string))
			}
			if v, ok := topicParamMap["start_time"]; ok {
				topicParam.StartTime = helper.IntUint64(v.(int))
			}
			if v, ok := topicParamMap["topic_id"]; ok {
				topicParam.TopicId = helper.String(v.(string))
			}
			if v, ok := topicParamMap["compression_type"]; ok {
				topicParam.CompressionType = helper.String(v.(string))
			}
			if v, ok := topicParamMap["use_auto_create_topic"]; ok {
				topicParam.UseAutoCreateTopic = helper.Bool(v.(bool))
			}
			if v, ok := topicParamMap["msg_multiple"]; ok {
				topicParam.MsgMultiple = helper.IntInt64(v.(int))
			}
			datahubResource.TopicParam = &topicParam
		}
		if mariaDBParamMap, ok := helper.InterfaceToMap(dMap, "maria_db_param"); ok {
			mariaDBParam := ckafka.MariaDBParam{}
			if v, ok := mariaDBParamMap["database"]; ok {
				mariaDBParam.Database = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["table"]; ok {
				mariaDBParam.Table = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["resource"]; ok {
				mariaDBParam.Resource = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["snapshot_mode"]; ok {
				mariaDBParam.SnapshotMode = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["key_columns"]; ok {
				mariaDBParam.KeyColumns = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["is_table_prefix"]; ok {
				mariaDBParam.IsTablePrefix = helper.Bool(v.(bool))
			}
			if v, ok := mariaDBParamMap["output_format"]; ok {
				mariaDBParam.OutputFormat = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["include_content_changes"]; ok {
				mariaDBParam.IncludeContentChanges = helper.String(v.(string))
			}
			if v, ok := mariaDBParamMap["include_query"]; ok {
				mariaDBParam.IncludeQuery = helper.Bool(v.(bool))
			}
			if v, ok := mariaDBParamMap["record_with_schema"]; ok {
				mariaDBParam.RecordWithSchema = helper.Bool(v.(bool))
			}
			datahubResource.MariaDBParam = &mariaDBParam
		}
		if sQLServerParamMap, ok := helper.InterfaceToMap(dMap, "sql_server_param"); ok {
			sQLServerParam := ckafka.SQLServerParam{}
			if v, ok := sQLServerParamMap["database"]; ok {
				sQLServerParam.Database = helper.String(v.(string))
			}
			if v, ok := sQLServerParamMap["table"]; ok {
				sQLServerParam.Table = helper.String(v.(string))
			}
			if v, ok := sQLServerParamMap["resource"]; ok {
				sQLServerParam.Resource = helper.String(v.(string))
			}
			if v, ok := sQLServerParamMap["snapshot_mode"]; ok {
				sQLServerParam.SnapshotMode = helper.String(v.(string))
			}
			datahubResource.SQLServerParam = &sQLServerParam
		}
		if ctsdbParamMap, ok := helper.InterfaceToMap(dMap, "ctsdb_param"); ok {
			ctsdbParam := ckafka.CtsdbParam{}
			if v, ok := ctsdbParamMap["resource"]; ok {
				ctsdbParam.Resource = helper.String(v.(string))
			}
			if v, ok := ctsdbParamMap["ctsdb_metric"]; ok {
				ctsdbParam.CtsdbMetric = helper.String(v.(string))
			}
			datahubResource.CtsdbParam = &ctsdbParam
		}
		if scfParamMap, ok := helper.InterfaceToMap(dMap, "scf_param"); ok {
			scfParam := ckafka.ScfParam{}
			if v, ok := scfParamMap["function_name"]; ok {
				scfParam.FunctionName = helper.String(v.(string))
			}
			if v, ok := scfParamMap["namespace"]; ok {
				scfParam.Namespace = helper.String(v.(string))
			}
			if v, ok := scfParamMap["qualifier"]; ok {
				scfParam.Qualifier = helper.String(v.(string))
			}
			if v, ok := scfParamMap["batch_size"]; ok {
				scfParam.BatchSize = helper.IntInt64(v.(int))
			}
			if v, ok := scfParamMap["max_retries"]; ok {
				scfParam.MaxRetries = helper.IntInt64(v.(int))
			}
			datahubResource.ScfParam = &scfParam
		}
		request.TargetResource = &datahubResource
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "transform_param"); ok {
		transformParam := ckafka.TransformParam{}
		if v, ok := dMap["analysis_format"]; ok {
			transformParam.AnalysisFormat = helper.String(v.(string))
		}
		if v, ok := dMap["output_format"]; ok {
			transformParam.OutputFormat = helper.String(v.(string))
		}
		if failureParamMap, ok := helper.InterfaceToMap(dMap, "failure_param"); ok {
			failureParam := ckafka.FailureParam{}
			if v, ok := failureParamMap["type"]; ok {
				failureParam.Type = helper.String(v.(string))
			}
			if kafkaParamMap, ok := helper.InterfaceToMap(failureParamMap, "kafka_param"); ok {
				kafkaParam := ckafka.KafkaParam{}
				if v, ok := kafkaParamMap["self_built"]; ok {
					kafkaParam.SelfBuilt = helper.Bool(v.(bool))
				}
				if v, ok := kafkaParamMap["resource"]; ok {
					kafkaParam.Resource = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["topic"]; ok {
					kafkaParam.Topic = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["offset_type"]; ok {
					kafkaParam.OffsetType = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["start_time"]; ok {
					kafkaParam.StartTime = helper.IntUint64(v.(int))
				}
				if v, ok := kafkaParamMap["resource_name"]; ok {
					kafkaParam.ResourceName = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["zone_id"]; ok {
					kafkaParam.ZoneId = helper.IntInt64(v.(int))
				}
				if v, ok := kafkaParamMap["topic_id"]; ok {
					kafkaParam.TopicId = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["partition_num"]; ok {
					kafkaParam.PartitionNum = helper.IntInt64(v.(int))
				}
				if v, ok := kafkaParamMap["enable_toleration"]; ok {
					kafkaParam.EnableToleration = helper.Bool(v.(bool))
				}
				if v, ok := kafkaParamMap["qps_limit"]; ok {
					kafkaParam.QpsLimit = helper.IntUint64(v.(int))
				}
				if v, ok := kafkaParamMap["table_mappings"]; ok {
					for _, item := range v.([]interface{}) {
						tableMappingsMap := item.(map[string]interface{})
						tableMapping := ckafka.TableMapping{}
						if v, ok := tableMappingsMap["database"]; ok {
							tableMapping.Database = helper.String(v.(string))
						}
						if v, ok := tableMappingsMap["table"]; ok {
							tableMapping.Table = helper.String(v.(string))
						}
						if v, ok := tableMappingsMap["topic"]; ok {
							tableMapping.Topic = helper.String(v.(string))
						}
						if v, ok := tableMappingsMap["topic_id"]; ok {
							tableMapping.TopicId = helper.String(v.(string))
						}
						kafkaParam.TableMappings = append(kafkaParam.TableMappings, &tableMapping)
					}
				}
				if v, ok := kafkaParamMap["use_table_mapping"]; ok {
					kafkaParam.UseTableMapping = helper.Bool(v.(bool))
				}
				if v, ok := kafkaParamMap["use_auto_create_topic"]; ok {
					kafkaParam.UseAutoCreateTopic = helper.Bool(v.(bool))
				}
				if v, ok := kafkaParamMap["compression_type"]; ok {
					kafkaParam.CompressionType = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["msg_multiple"]; ok {
					kafkaParam.MsgMultiple = helper.IntInt64(v.(int))
				}
				failureParam.KafkaParam = &kafkaParam
			}
			if v, ok := failureParamMap["retry_interval"]; ok {
				failureParam.RetryInterval = helper.IntUint64(v.(int))
			}
			if v, ok := failureParamMap["max_retry_attempts"]; ok {
				failureParam.MaxRetryAttempts = helper.IntUint64(v.(int))
			}
			if topicParamMap, ok := helper.InterfaceToMap(failureParamMap, "topic_param"); ok {
				topicParam := ckafka.TopicParam{}
				if v, ok := topicParamMap["resource"]; ok {
					topicParam.Resource = helper.String(v.(string))
				}
				if v, ok := topicParamMap["offset_type"]; ok {
					topicParam.OffsetType = helper.String(v.(string))
				}
				if v, ok := topicParamMap["start_time"]; ok {
					topicParam.StartTime = helper.IntUint64(v.(int))
				}
				if v, ok := topicParamMap["topic_id"]; ok {
					topicParam.TopicId = helper.String(v.(string))
				}
				if v, ok := topicParamMap["compression_type"]; ok {
					topicParam.CompressionType = helper.String(v.(string))
				}
				if v, ok := topicParamMap["use_auto_create_topic"]; ok {
					topicParam.UseAutoCreateTopic = helper.Bool(v.(bool))
				}
				if v, ok := topicParamMap["msg_multiple"]; ok {
					topicParam.MsgMultiple = helper.IntInt64(v.(int))
				}
				failureParam.TopicParam = &topicParam
			}
			if v, ok := failureParamMap["dlq_type"]; ok {
				failureParam.DlqType = helper.String(v.(string))
			}
			transformParam.FailureParam = &failureParam
		}
		if v, ok := dMap["content"]; ok {
			transformParam.Content = helper.String(v.(string))
		}
		if v, ok := dMap["source_type"]; ok {
			transformParam.SourceType = helper.String(v.(string))
		}
		if v, ok := dMap["regex"]; ok {
			transformParam.Regex = helper.String(v.(string))
		}
		if v, ok := dMap["map_param"]; ok {
			for _, item := range v.([]interface{}) {
				mapParamMap := item.(map[string]interface{})
				mapParam := ckafka.MapParam{}
				if v, ok := mapParamMap["key"]; ok {
					mapParam.Key = helper.String(v.(string))
				}
				if v, ok := mapParamMap["type"]; ok {
					mapParam.Type = helper.String(v.(string))
				}
				if v, ok := mapParamMap["value"]; ok {
					mapParam.Value = helper.String(v.(string))
				}
				transformParam.MapParam = append(transformParam.MapParam, &mapParam)
			}
		}
		if v, ok := dMap["filter_param"]; ok {
			for _, item := range v.([]interface{}) {
				filterParamMap := item.(map[string]interface{})
				filterMapParam := ckafka.FilterMapParam{}
				if v, ok := filterParamMap["key"]; ok {
					filterMapParam.Key = helper.String(v.(string))
				}
				if v, ok := filterParamMap["match_mode"]; ok {
					filterMapParam.MatchMode = helper.String(v.(string))
				}
				if v, ok := filterParamMap["value"]; ok {
					filterMapParam.Value = helper.String(v.(string))
				}
				if v, ok := filterParamMap["type"]; ok {
					filterMapParam.Type = helper.String(v.(string))
				}
				transformParam.FilterParam = append(transformParam.FilterParam, &filterMapParam)
			}
		}
		if v, ok := dMap["result"]; ok {
			transformParam.Result = helper.String(v.(string))
		}
		if v, ok := dMap["analyse_result"]; ok {
			for _, item := range v.([]interface{}) {
				analyseResultMap := item.(map[string]interface{})
				mapParam := ckafka.MapParam{}
				if v, ok := analyseResultMap["key"]; ok {
					mapParam.Key = helper.String(v.(string))
				}
				if v, ok := analyseResultMap["type"]; ok {
					mapParam.Type = helper.String(v.(string))
				}
				if v, ok := analyseResultMap["value"]; ok {
					mapParam.Value = helper.String(v.(string))
				}
				transformParam.AnalyseResult = append(transformParam.AnalyseResult, &mapParam)
			}
		}
		if v, ok := dMap["use_event_bus"]; ok {
			transformParam.UseEventBus = helper.Bool(v.(bool))
		}
		request.TransformParam = &transformParam
	}

	if v, ok := d.GetOk("schema_id"); ok {
		request.SchemaId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "transforms_param"); ok {
		transformsParam := ckafka.TransformsParam{}
		if v, ok := dMap["content"]; ok {
			transformsParam.Content = helper.String(v.(string))
		}
		if v, ok := dMap["field_chain"]; ok {
			for _, item := range v.([]interface{}) {
				fieldChainMap := item.(map[string]interface{})
				fieldParam := ckafka.FieldParam{}
				if analyseMap, ok := helper.InterfaceToMap(fieldChainMap, "analyse"); ok {
					analyseParam := ckafka.AnalyseParam{}
					if v, ok := analyseMap["format"]; ok {
						analyseParam.Format = helper.String(v.(string))
					}
					if v, ok := analyseMap["regex"]; ok {
						analyseParam.Regex = helper.String(v.(string))
					}
					if v, ok := analyseMap["input_value_type"]; ok && v != "" {
						analyseParam.InputValueType = helper.String(v.(string))
					}
					if v, ok := analyseMap["input_value"]; ok && v != "" {
						analyseParam.InputValue = helper.String(v.(string))
					}
					fieldParam.Analyse = &analyseParam
				}
				if secondaryAnalyseMap, ok := helper.InterfaceToMap(fieldChainMap, "secondary_analyse"); ok {
					secondaryAnalyseParam := ckafka.SecondaryAnalyseParam{}
					if v, ok := secondaryAnalyseMap["regex"]; ok {
						secondaryAnalyseParam.Regex = helper.String(v.(string))
					}
					fieldParam.SecondaryAnalyse = &secondaryAnalyseParam
				}
				if v, ok := fieldChainMap["s_m_t"]; ok {
					for _, item := range v.([]interface{}) {
						sMTMap := item.(map[string]interface{})
						sMTParam := ckafka.SMTParam{}
						if v, ok := sMTMap["key"]; ok {
							sMTParam.Key = helper.String(v.(string))
						}
						if v, ok := sMTMap["operate"]; ok {
							sMTParam.Operate = helper.String(v.(string))
						}
						if v, ok := sMTMap["scheme_type"]; ok {
							sMTParam.SchemeType = helper.String(v.(string))
						}
						if v, ok := sMTMap["value"]; ok {
							sMTParam.Value = helper.String(v.(string))
						}
						if valueOperateMap, ok := helper.InterfaceToMap(sMTMap, "value_operate"); ok {
							valueParam := ckafka.ValueParam{}
							if v, ok := valueOperateMap["type"]; ok {
								valueParam.Type = helper.String(v.(string))
							}
							if replaceMap, ok := helper.InterfaceToMap(valueOperateMap, "replace"); ok {
								replaceParam := ckafka.ReplaceParam{}
								if v, ok := replaceMap["old_value"]; ok {
									replaceParam.OldValue = helper.String(v.(string))
								}
								if v, ok := replaceMap["new_value"]; ok {
									replaceParam.NewValue = helper.String(v.(string))
								}
								valueParam.Replace = &replaceParam
							}
							if substrMap, ok := helper.InterfaceToMap(valueOperateMap, "substr"); ok {
								substrParam := ckafka.SubstrParam{}
								if v, ok := substrMap["start"]; ok {
									substrParam.Start = helper.IntInt64(v.(int))
								}
								if v, ok := substrMap["end"]; ok {
									substrParam.End = helper.IntInt64(v.(int))
								}
								valueParam.Substr = &substrParam
							}
							if dateMap, ok := helper.InterfaceToMap(valueOperateMap, "date"); ok {
								dateParam := ckafka.DateParam{}
								if v, ok := dateMap["format"]; ok {
									dateParam.Format = helper.String(v.(string))
								}
								if v, ok := dateMap["target_type"]; ok {
									dateParam.TargetType = helper.String(v.(string))
								}
								if v, ok := dateMap["time_zone"]; ok {
									dateParam.TimeZone = helper.String(v.(string))
								}
								valueParam.Date = &dateParam
							}
							if regexReplaceMap, ok := helper.InterfaceToMap(valueOperateMap, "regex_replace"); ok {
								regexReplaceParam := ckafka.RegexReplaceParam{}
								if v, ok := regexReplaceMap["regex"]; ok {
									regexReplaceParam.Regex = helper.String(v.(string))
								}
								if v, ok := regexReplaceMap["new_value"]; ok {
									regexReplaceParam.NewValue = helper.String(v.(string))
								}
								valueParam.RegexReplace = &regexReplaceParam
							}
							if splitMap, ok := helper.InterfaceToMap(valueOperateMap, "split"); ok {
								splitParam := ckafka.SplitParam{}
								if v, ok := splitMap["regex"]; ok {
									splitParam.Regex = helper.String(v.(string))
								}
								valueParam.Split = &splitParam
							}
							if kVMap, ok := helper.InterfaceToMap(valueOperateMap, "k_v"); ok {
								kVParam := ckafka.KVParam{}
								if v, ok := kVMap["delimiter"]; ok {
									kVParam.Delimiter = helper.String(v.(string))
								}
								if v, ok := kVMap["regex"]; ok {
									kVParam.Regex = helper.String(v.(string))
								}
								if v, ok := kVMap["keep_original_key"]; ok {
									kVParam.KeepOriginalKey = helper.String(v.(string))
								}
								valueParam.KV = &kVParam
							}
							if v, ok := valueOperateMap["result"]; ok {
								valueParam.Result = helper.String(v.(string))
							}
							if jsonPathReplaceMap, ok := helper.InterfaceToMap(valueOperateMap, "json_path_replace"); ok {
								jsonPathReplaceParam := ckafka.JsonPathReplaceParam{}
								if v, ok := jsonPathReplaceMap["old_value"]; ok {
									jsonPathReplaceParam.OldValue = helper.String(v.(string))
								}
								if v, ok := jsonPathReplaceMap["new_value"]; ok {
									jsonPathReplaceParam.NewValue = helper.String(v.(string))
								}
								valueParam.JsonPathReplace = &jsonPathReplaceParam
							}
							if urlDecodeMap, ok := helper.InterfaceToMap(valueOperateMap, "url_decode"); ok {
								urlDecodeParam := ckafka.UrlDecodeParam{}
								if v, ok := urlDecodeMap["charset_name"]; ok {
									urlDecodeParam.CharsetName = helper.String(v.(string))
								}
								valueParam.UrlDecode = &urlDecodeParam
							}
							sMTParam.ValueOperate = &valueParam
						}
						if v, ok := sMTMap["original_value"]; ok {
							sMTParam.OriginalValue = helper.String(v.(string))
						}
						if v, ok := sMTMap["value_operates"]; ok {
							for _, item := range v.([]interface{}) {
								valueOperatesMap := item.(map[string]interface{})
								valueParam := ckafka.ValueParam{}
								if v, ok := valueOperatesMap["type"]; ok {
									valueParam.Type = helper.String(v.(string))
								}
								if replaceMap, ok := helper.InterfaceToMap(valueOperatesMap, "replace"); ok {
									replaceParam := ckafka.ReplaceParam{}
									if v, ok := replaceMap["old_value"]; ok {
										replaceParam.OldValue = helper.String(v.(string))
									}
									if v, ok := replaceMap["new_value"]; ok {
										replaceParam.NewValue = helper.String(v.(string))
									}
									valueParam.Replace = &replaceParam
								}
								if substrMap, ok := helper.InterfaceToMap(valueOperatesMap, "substr"); ok {
									substrParam := ckafka.SubstrParam{}
									if v, ok := substrMap["start"]; ok {
										substrParam.Start = helper.IntInt64(v.(int))
									}
									if v, ok := substrMap["end"]; ok {
										substrParam.End = helper.IntInt64(v.(int))
									}
									valueParam.Substr = &substrParam
								}
								if dateMap, ok := helper.InterfaceToMap(valueOperatesMap, "date"); ok {
									dateParam := ckafka.DateParam{}
									if v, ok := dateMap["format"]; ok {
										dateParam.Format = helper.String(v.(string))
									}
									if v, ok := dateMap["target_type"]; ok {
										dateParam.TargetType = helper.String(v.(string))
									}
									if v, ok := dateMap["time_zone"]; ok {
										dateParam.TimeZone = helper.String(v.(string))
									}
									valueParam.Date = &dateParam
								}
								if regexReplaceMap, ok := helper.InterfaceToMap(valueOperatesMap, "regex_replace"); ok {
									regexReplaceParam := ckafka.RegexReplaceParam{}
									if v, ok := regexReplaceMap["regex"]; ok {
										regexReplaceParam.Regex = helper.String(v.(string))
									}
									if v, ok := regexReplaceMap["new_value"]; ok {
										regexReplaceParam.NewValue = helper.String(v.(string))
									}
									valueParam.RegexReplace = &regexReplaceParam
								}
								if splitMap, ok := helper.InterfaceToMap(valueOperatesMap, "split"); ok {
									splitParam := ckafka.SplitParam{}
									if v, ok := splitMap["regex"]; ok {
										splitParam.Regex = helper.String(v.(string))
									}
									valueParam.Split = &splitParam
								}
								if kVMap, ok := helper.InterfaceToMap(valueOperatesMap, "k_v"); ok {
									kVParam := ckafka.KVParam{}
									if v, ok := kVMap["delimiter"]; ok {
										kVParam.Delimiter = helper.String(v.(string))
									}
									if v, ok := kVMap["regex"]; ok {
										kVParam.Regex = helper.String(v.(string))
									}
									if v, ok := kVMap["keep_original_key"]; ok {
										kVParam.KeepOriginalKey = helper.String(v.(string))
									}
									valueParam.KV = &kVParam
								}
								if v, ok := valueOperatesMap["result"]; ok {
									valueParam.Result = helper.String(v.(string))
								}
								if jsonPathReplaceMap, ok := helper.InterfaceToMap(valueOperatesMap, "json_path_replace"); ok {
									jsonPathReplaceParam := ckafka.JsonPathReplaceParam{}
									if v, ok := jsonPathReplaceMap["old_value"]; ok {
										jsonPathReplaceParam.OldValue = helper.String(v.(string))
									}
									if v, ok := jsonPathReplaceMap["new_value"]; ok {
										jsonPathReplaceParam.NewValue = helper.String(v.(string))
									}
									valueParam.JsonPathReplace = &jsonPathReplaceParam
								}
								if urlDecodeMap, ok := helper.InterfaceToMap(valueOperatesMap, "url_decode"); ok {
									urlDecodeParam := ckafka.UrlDecodeParam{}
									if v, ok := urlDecodeMap["charset_name"]; ok {
										urlDecodeParam.CharsetName = helper.String(v.(string))
									}
									valueParam.UrlDecode = &urlDecodeParam
								}
								sMTParam.ValueOperates = append(sMTParam.ValueOperates, &valueParam)
							}
						}
						fieldParam.SMT = append(fieldParam.SMT, &sMTParam)
					}
				}
				if v, ok := fieldChainMap["result"]; ok {
					fieldParam.Result = helper.String(v.(string))
				}
				if v, ok := fieldChainMap["analyse_result"]; ok {
					for _, item := range v.([]interface{}) {
						analyseResultMap := item.(map[string]interface{})
						sMTParam := ckafka.SMTParam{}
						if v, ok := analyseResultMap["key"]; ok {
							sMTParam.Key = helper.String(v.(string))
						}
						if v, ok := analyseResultMap["operate"]; ok {
							sMTParam.Operate = helper.String(v.(string))
						}
						if v, ok := analyseResultMap["scheme_type"]; ok {
							sMTParam.SchemeType = helper.String(v.(string))
						}
						if v, ok := analyseResultMap["value"]; ok {
							sMTParam.Value = helper.String(v.(string))
						}
						if valueOperateMap, ok := helper.InterfaceToMap(analyseResultMap, "value_operate"); ok {
							valueParam := ckafka.ValueParam{}
							if v, ok := valueOperateMap["type"]; ok {
								valueParam.Type = helper.String(v.(string))
							}
							if replaceMap, ok := helper.InterfaceToMap(valueOperateMap, "replace"); ok {
								replaceParam := ckafka.ReplaceParam{}
								if v, ok := replaceMap["old_value"]; ok {
									replaceParam.OldValue = helper.String(v.(string))
								}
								if v, ok := replaceMap["new_value"]; ok {
									replaceParam.NewValue = helper.String(v.(string))
								}
								valueParam.Replace = &replaceParam
							}
							if substrMap, ok := helper.InterfaceToMap(valueOperateMap, "substr"); ok {
								substrParam := ckafka.SubstrParam{}
								if v, ok := substrMap["start"]; ok {
									substrParam.Start = helper.IntInt64(v.(int))
								}
								if v, ok := substrMap["end"]; ok {
									substrParam.End = helper.IntInt64(v.(int))
								}
								valueParam.Substr = &substrParam
							}
							if dateMap, ok := helper.InterfaceToMap(valueOperateMap, "date"); ok {
								dateParam := ckafka.DateParam{}
								if v, ok := dateMap["format"]; ok {
									dateParam.Format = helper.String(v.(string))
								}
								if v, ok := dateMap["target_type"]; ok {
									dateParam.TargetType = helper.String(v.(string))
								}
								if v, ok := dateMap["time_zone"]; ok {
									dateParam.TimeZone = helper.String(v.(string))
								}
								valueParam.Date = &dateParam
							}
							if regexReplaceMap, ok := helper.InterfaceToMap(valueOperateMap, "regex_replace"); ok {
								regexReplaceParam := ckafka.RegexReplaceParam{}
								if v, ok := regexReplaceMap["regex"]; ok {
									regexReplaceParam.Regex = helper.String(v.(string))
								}
								if v, ok := regexReplaceMap["new_value"]; ok {
									regexReplaceParam.NewValue = helper.String(v.(string))
								}
								valueParam.RegexReplace = &regexReplaceParam
							}
							if splitMap, ok := helper.InterfaceToMap(valueOperateMap, "split"); ok {
								splitParam := ckafka.SplitParam{}
								if v, ok := splitMap["regex"]; ok {
									splitParam.Regex = helper.String(v.(string))
								}
								valueParam.Split = &splitParam
							}
							if kVMap, ok := helper.InterfaceToMap(valueOperateMap, "k_v"); ok {
								kVParam := ckafka.KVParam{}
								if v, ok := kVMap["delimiter"]; ok {
									kVParam.Delimiter = helper.String(v.(string))
								}
								if v, ok := kVMap["regex"]; ok {
									kVParam.Regex = helper.String(v.(string))
								}
								if v, ok := kVMap["keep_original_key"]; ok {
									kVParam.KeepOriginalKey = helper.String(v.(string))
								}
								valueParam.KV = &kVParam
							}
							if v, ok := valueOperateMap["result"]; ok {
								valueParam.Result = helper.String(v.(string))
							}
							if jsonPathReplaceMap, ok := helper.InterfaceToMap(valueOperateMap, "json_path_replace"); ok {
								jsonPathReplaceParam := ckafka.JsonPathReplaceParam{}
								if v, ok := jsonPathReplaceMap["old_value"]; ok {
									jsonPathReplaceParam.OldValue = helper.String(v.(string))
								}
								if v, ok := jsonPathReplaceMap["new_value"]; ok {
									jsonPathReplaceParam.NewValue = helper.String(v.(string))
								}
								valueParam.JsonPathReplace = &jsonPathReplaceParam
							}
							if urlDecodeMap, ok := helper.InterfaceToMap(valueOperateMap, "url_decode"); ok {
								urlDecodeParam := ckafka.UrlDecodeParam{}
								if v, ok := urlDecodeMap["charset_name"]; ok {
									urlDecodeParam.CharsetName = helper.String(v.(string))
								}
								valueParam.UrlDecode = &urlDecodeParam
							}
							sMTParam.ValueOperate = &valueParam
						}
						if v, ok := analyseResultMap["original_value"]; ok {
							sMTParam.OriginalValue = helper.String(v.(string))
						}
						if v, ok := analyseResultMap["value_operates"]; ok {
							for _, item := range v.([]interface{}) {
								valueOperatesMap := item.(map[string]interface{})
								valueParam := ckafka.ValueParam{}
								if v, ok := valueOperatesMap["type"]; ok {
									valueParam.Type = helper.String(v.(string))
								}
								if replaceMap, ok := helper.InterfaceToMap(valueOperatesMap, "replace"); ok {
									replaceParam := ckafka.ReplaceParam{}
									if v, ok := replaceMap["old_value"]; ok {
										replaceParam.OldValue = helper.String(v.(string))
									}
									if v, ok := replaceMap["new_value"]; ok {
										replaceParam.NewValue = helper.String(v.(string))
									}
									valueParam.Replace = &replaceParam
								}
								if substrMap, ok := helper.InterfaceToMap(valueOperatesMap, "substr"); ok {
									substrParam := ckafka.SubstrParam{}
									if v, ok := substrMap["start"]; ok {
										substrParam.Start = helper.IntInt64(v.(int))
									}
									if v, ok := substrMap["end"]; ok {
										substrParam.End = helper.IntInt64(v.(int))
									}
									valueParam.Substr = &substrParam
								}
								if dateMap, ok := helper.InterfaceToMap(valueOperatesMap, "date"); ok {
									dateParam := ckafka.DateParam{}
									if v, ok := dateMap["format"]; ok {
										dateParam.Format = helper.String(v.(string))
									}
									if v, ok := dateMap["target_type"]; ok {
										dateParam.TargetType = helper.String(v.(string))
									}
									if v, ok := dateMap["time_zone"]; ok {
										dateParam.TimeZone = helper.String(v.(string))
									}
									valueParam.Date = &dateParam
								}
								if regexReplaceMap, ok := helper.InterfaceToMap(valueOperatesMap, "regex_replace"); ok {
									regexReplaceParam := ckafka.RegexReplaceParam{}
									if v, ok := regexReplaceMap["regex"]; ok {
										regexReplaceParam.Regex = helper.String(v.(string))
									}
									if v, ok := regexReplaceMap["new_value"]; ok {
										regexReplaceParam.NewValue = helper.String(v.(string))
									}
									valueParam.RegexReplace = &regexReplaceParam
								}
								if splitMap, ok := helper.InterfaceToMap(valueOperatesMap, "split"); ok {
									splitParam := ckafka.SplitParam{}
									if v, ok := splitMap["regex"]; ok {
										splitParam.Regex = helper.String(v.(string))
									}
									valueParam.Split = &splitParam
								}
								if kVMap, ok := helper.InterfaceToMap(valueOperatesMap, "k_v"); ok {
									kVParam := ckafka.KVParam{}
									if v, ok := kVMap["delimiter"]; ok {
										kVParam.Delimiter = helper.String(v.(string))
									}
									if v, ok := kVMap["regex"]; ok {
										kVParam.Regex = helper.String(v.(string))
									}
									if v, ok := kVMap["keep_original_key"]; ok {
										kVParam.KeepOriginalKey = helper.String(v.(string))
									}
									valueParam.KV = &kVParam
								}
								if v, ok := valueOperatesMap["result"]; ok {
									valueParam.Result = helper.String(v.(string))
								}
								if jsonPathReplaceMap, ok := helper.InterfaceToMap(valueOperatesMap, "json_path_replace"); ok {
									jsonPathReplaceParam := ckafka.JsonPathReplaceParam{}
									if v, ok := jsonPathReplaceMap["old_value"]; ok {
										jsonPathReplaceParam.OldValue = helper.String(v.(string))
									}
									if v, ok := jsonPathReplaceMap["new_value"]; ok {
										jsonPathReplaceParam.NewValue = helper.String(v.(string))
									}
									valueParam.JsonPathReplace = &jsonPathReplaceParam
								}
								if urlDecodeMap, ok := helper.InterfaceToMap(valueOperatesMap, "url_decode"); ok {
									urlDecodeParam := ckafka.UrlDecodeParam{}
									if v, ok := urlDecodeMap["charset_name"]; ok {
										urlDecodeParam.CharsetName = helper.String(v.(string))
									}
									valueParam.UrlDecode = &urlDecodeParam
								}
								sMTParam.ValueOperates = append(sMTParam.ValueOperates, &valueParam)
							}
						}
						fieldParam.AnalyseResult = append(fieldParam.AnalyseResult, &sMTParam)
					}
				}
				if v, ok := fieldChainMap["secondary_analyse_result"]; ok {
					for _, item := range v.([]interface{}) {
						secondaryAnalyseResultMap := item.(map[string]interface{})
						sMTParam := ckafka.SMTParam{}
						if v, ok := secondaryAnalyseResultMap["key"]; ok {
							sMTParam.Key = helper.String(v.(string))
						}
						if v, ok := secondaryAnalyseResultMap["operate"]; ok {
							sMTParam.Operate = helper.String(v.(string))
						}
						if v, ok := secondaryAnalyseResultMap["scheme_type"]; ok {
							sMTParam.SchemeType = helper.String(v.(string))
						}
						if v, ok := secondaryAnalyseResultMap["value"]; ok {
							sMTParam.Value = helper.String(v.(string))
						}
						if valueOperateMap, ok := helper.InterfaceToMap(secondaryAnalyseResultMap, "value_operate"); ok {
							valueParam := ckafka.ValueParam{}
							if v, ok := valueOperateMap["type"]; ok {
								valueParam.Type = helper.String(v.(string))
							}
							if replaceMap, ok := helper.InterfaceToMap(valueOperateMap, "replace"); ok {
								replaceParam := ckafka.ReplaceParam{}
								if v, ok := replaceMap["old_value"]; ok {
									replaceParam.OldValue = helper.String(v.(string))
								}
								if v, ok := replaceMap["new_value"]; ok {
									replaceParam.NewValue = helper.String(v.(string))
								}
								valueParam.Replace = &replaceParam
							}
							if substrMap, ok := helper.InterfaceToMap(valueOperateMap, "substr"); ok {
								substrParam := ckafka.SubstrParam{}
								if v, ok := substrMap["start"]; ok {
									substrParam.Start = helper.IntInt64(v.(int))
								}
								if v, ok := substrMap["end"]; ok {
									substrParam.End = helper.IntInt64(v.(int))
								}
								valueParam.Substr = &substrParam
							}
							if dateMap, ok := helper.InterfaceToMap(valueOperateMap, "date"); ok {
								dateParam := ckafka.DateParam{}
								if v, ok := dateMap["format"]; ok {
									dateParam.Format = helper.String(v.(string))
								}
								if v, ok := dateMap["target_type"]; ok {
									dateParam.TargetType = helper.String(v.(string))
								}
								if v, ok := dateMap["time_zone"]; ok {
									dateParam.TimeZone = helper.String(v.(string))
								}
								valueParam.Date = &dateParam
							}
							if regexReplaceMap, ok := helper.InterfaceToMap(valueOperateMap, "regex_replace"); ok {
								regexReplaceParam := ckafka.RegexReplaceParam{}
								if v, ok := regexReplaceMap["regex"]; ok {
									regexReplaceParam.Regex = helper.String(v.(string))
								}
								if v, ok := regexReplaceMap["new_value"]; ok {
									regexReplaceParam.NewValue = helper.String(v.(string))
								}
								valueParam.RegexReplace = &regexReplaceParam
							}
							if splitMap, ok := helper.InterfaceToMap(valueOperateMap, "split"); ok {
								splitParam := ckafka.SplitParam{}
								if v, ok := splitMap["regex"]; ok {
									splitParam.Regex = helper.String(v.(string))
								}
								valueParam.Split = &splitParam
							}
							if kVMap, ok := helper.InterfaceToMap(valueOperateMap, "k_v"); ok {
								kVParam := ckafka.KVParam{}
								if v, ok := kVMap["delimiter"]; ok {
									kVParam.Delimiter = helper.String(v.(string))
								}
								if v, ok := kVMap["regex"]; ok {
									kVParam.Regex = helper.String(v.(string))
								}
								if v, ok := kVMap["keep_original_key"]; ok {
									kVParam.KeepOriginalKey = helper.String(v.(string))
								}
								valueParam.KV = &kVParam
							}
							if v, ok := valueOperateMap["result"]; ok {
								valueParam.Result = helper.String(v.(string))
							}
							if jsonPathReplaceMap, ok := helper.InterfaceToMap(valueOperateMap, "json_path_replace"); ok {
								jsonPathReplaceParam := ckafka.JsonPathReplaceParam{}
								if v, ok := jsonPathReplaceMap["old_value"]; ok {
									jsonPathReplaceParam.OldValue = helper.String(v.(string))
								}
								if v, ok := jsonPathReplaceMap["new_value"]; ok {
									jsonPathReplaceParam.NewValue = helper.String(v.(string))
								}
								valueParam.JsonPathReplace = &jsonPathReplaceParam
							}
							if urlDecodeMap, ok := helper.InterfaceToMap(valueOperateMap, "url_decode"); ok {
								urlDecodeParam := ckafka.UrlDecodeParam{}
								if v, ok := urlDecodeMap["charset_name"]; ok {
									urlDecodeParam.CharsetName = helper.String(v.(string))
								}
								valueParam.UrlDecode = &urlDecodeParam
							}
							sMTParam.ValueOperate = &valueParam
						}
						if v, ok := secondaryAnalyseResultMap["original_value"]; ok {
							sMTParam.OriginalValue = helper.String(v.(string))
						}
						if v, ok := secondaryAnalyseResultMap["value_operates"]; ok {
							for _, item := range v.([]interface{}) {
								valueOperatesMap := item.(map[string]interface{})
								valueParam := ckafka.ValueParam{}
								if v, ok := valueOperatesMap["type"]; ok {
									valueParam.Type = helper.String(v.(string))
								}
								if replaceMap, ok := helper.InterfaceToMap(valueOperatesMap, "replace"); ok {
									replaceParam := ckafka.ReplaceParam{}
									if v, ok := replaceMap["old_value"]; ok {
										replaceParam.OldValue = helper.String(v.(string))
									}
									if v, ok := replaceMap["new_value"]; ok {
										replaceParam.NewValue = helper.String(v.(string))
									}
									valueParam.Replace = &replaceParam
								}
								if substrMap, ok := helper.InterfaceToMap(valueOperatesMap, "substr"); ok {
									substrParam := ckafka.SubstrParam{}
									if v, ok := substrMap["start"]; ok {
										substrParam.Start = helper.IntInt64(v.(int))
									}
									if v, ok := substrMap["end"]; ok {
										substrParam.End = helper.IntInt64(v.(int))
									}
									valueParam.Substr = &substrParam
								}
								if dateMap, ok := helper.InterfaceToMap(valueOperatesMap, "date"); ok {
									dateParam := ckafka.DateParam{}
									if v, ok := dateMap["format"]; ok {
										dateParam.Format = helper.String(v.(string))
									}
									if v, ok := dateMap["target_type"]; ok {
										dateParam.TargetType = helper.String(v.(string))
									}
									if v, ok := dateMap["time_zone"]; ok {
										dateParam.TimeZone = helper.String(v.(string))
									}
									valueParam.Date = &dateParam
								}
								if regexReplaceMap, ok := helper.InterfaceToMap(valueOperatesMap, "regex_replace"); ok {
									regexReplaceParam := ckafka.RegexReplaceParam{}
									if v, ok := regexReplaceMap["regex"]; ok {
										regexReplaceParam.Regex = helper.String(v.(string))
									}
									if v, ok := regexReplaceMap["new_value"]; ok {
										regexReplaceParam.NewValue = helper.String(v.(string))
									}
									valueParam.RegexReplace = &regexReplaceParam
								}
								if splitMap, ok := helper.InterfaceToMap(valueOperatesMap, "split"); ok {
									splitParam := ckafka.SplitParam{}
									if v, ok := splitMap["regex"]; ok {
										splitParam.Regex = helper.String(v.(string))
									}
									valueParam.Split = &splitParam
								}
								if kVMap, ok := helper.InterfaceToMap(valueOperatesMap, "k_v"); ok {
									kVParam := ckafka.KVParam{}
									if v, ok := kVMap["delimiter"]; ok {
										kVParam.Delimiter = helper.String(v.(string))
									}
									if v, ok := kVMap["regex"]; ok {
										kVParam.Regex = helper.String(v.(string))
									}
									if v, ok := kVMap["keep_original_key"]; ok {
										kVParam.KeepOriginalKey = helper.String(v.(string))
									}
									valueParam.KV = &kVParam
								}
								if v, ok := valueOperatesMap["result"]; ok {
									valueParam.Result = helper.String(v.(string))
								}
								if jsonPathReplaceMap, ok := helper.InterfaceToMap(valueOperatesMap, "json_path_replace"); ok {
									jsonPathReplaceParam := ckafka.JsonPathReplaceParam{}
									if v, ok := jsonPathReplaceMap["old_value"]; ok {
										jsonPathReplaceParam.OldValue = helper.String(v.(string))
									}
									if v, ok := jsonPathReplaceMap["new_value"]; ok {
										jsonPathReplaceParam.NewValue = helper.String(v.(string))
									}
									valueParam.JsonPathReplace = &jsonPathReplaceParam
								}
								if urlDecodeMap, ok := helper.InterfaceToMap(valueOperatesMap, "url_decode"); ok {
									urlDecodeParam := ckafka.UrlDecodeParam{}
									if v, ok := urlDecodeMap["charset_name"]; ok {
										urlDecodeParam.CharsetName = helper.String(v.(string))
									}
									valueParam.UrlDecode = &urlDecodeParam
								}
								sMTParam.ValueOperates = append(sMTParam.ValueOperates, &valueParam)
							}
						}
						fieldParam.SecondaryAnalyseResult = append(fieldParam.SecondaryAnalyseResult, &sMTParam)
					}
				}
				if v, ok := fieldChainMap["analyse_json_result"]; ok {
					fieldParam.AnalyseJsonResult = helper.String(v.(string))
				}
				if v, ok := fieldChainMap["secondary_analyse_json_result"]; ok {
					fieldParam.SecondaryAnalyseJsonResult = helper.String(v.(string))
				}
				transformsParam.FieldChain = append(transformsParam.FieldChain, &fieldParam)
			}
		}
		if v, ok := dMap["filter_param"]; ok {
			for _, item := range v.([]interface{}) {
				filterParamMap := item.(map[string]interface{})
				filterMapParam := ckafka.FilterMapParam{}
				if v, ok := filterParamMap["key"]; ok {
					filterMapParam.Key = helper.String(v.(string))
				}
				if v, ok := filterParamMap["match_mode"]; ok {
					filterMapParam.MatchMode = helper.String(v.(string))
				}
				if v, ok := filterParamMap["value"]; ok {
					filterMapParam.Value = helper.String(v.(string))
				}
				if v, ok := filterParamMap["type"]; ok {
					filterMapParam.Type = helper.String(v.(string))
				}
				transformsParam.FilterParam = append(transformsParam.FilterParam, &filterMapParam)
			}
		}
		if failureParamMap, ok := helper.InterfaceToMap(dMap, "failure_param"); ok {
			failureParam := ckafka.FailureParam{}
			if v, ok := failureParamMap["type"]; ok {
				failureParam.Type = helper.String(v.(string))
			}
			if kafkaParamMap, ok := helper.InterfaceToMap(failureParamMap, "kafka_param"); ok {
				kafkaParam := ckafka.KafkaParam{}
				if v, ok := kafkaParamMap["self_built"]; ok {
					kafkaParam.SelfBuilt = helper.Bool(v.(bool))
				}
				if v, ok := kafkaParamMap["resource"]; ok {
					kafkaParam.Resource = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["topic"]; ok {
					kafkaParam.Topic = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["offset_type"]; ok {
					kafkaParam.OffsetType = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["start_time"]; ok {
					kafkaParam.StartTime = helper.IntUint64(v.(int))
				}
				if v, ok := kafkaParamMap["resource_name"]; ok {
					kafkaParam.ResourceName = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["zone_id"]; ok {
					kafkaParam.ZoneId = helper.IntInt64(v.(int))
				}
				if v, ok := kafkaParamMap["topic_id"]; ok {
					kafkaParam.TopicId = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["partition_num"]; ok {
					kafkaParam.PartitionNum = helper.IntInt64(v.(int))
				}
				if v, ok := kafkaParamMap["enable_toleration"]; ok {
					kafkaParam.EnableToleration = helper.Bool(v.(bool))
				}
				if v, ok := kafkaParamMap["qps_limit"]; ok {
					kafkaParam.QpsLimit = helper.IntUint64(v.(int))
				}
				if v, ok := kafkaParamMap["table_mappings"]; ok {
					for _, item := range v.([]interface{}) {
						tableMappingsMap := item.(map[string]interface{})
						tableMapping := ckafka.TableMapping{}
						if v, ok := tableMappingsMap["database"]; ok {
							tableMapping.Database = helper.String(v.(string))
						}
						if v, ok := tableMappingsMap["table"]; ok {
							tableMapping.Table = helper.String(v.(string))
						}
						if v, ok := tableMappingsMap["topic"]; ok {
							tableMapping.Topic = helper.String(v.(string))
						}
						if v, ok := tableMappingsMap["topic_id"]; ok {
							tableMapping.TopicId = helper.String(v.(string))
						}
						kafkaParam.TableMappings = append(kafkaParam.TableMappings, &tableMapping)
					}
				}
				if v, ok := kafkaParamMap["use_table_mapping"]; ok {
					kafkaParam.UseTableMapping = helper.Bool(v.(bool))
				}
				if v, ok := kafkaParamMap["use_auto_create_topic"]; ok {
					kafkaParam.UseAutoCreateTopic = helper.Bool(v.(bool))
				}
				if v, ok := kafkaParamMap["compression_type"]; ok {
					kafkaParam.CompressionType = helper.String(v.(string))
				}
				if v, ok := kafkaParamMap["msg_multiple"]; ok {
					kafkaParam.MsgMultiple = helper.IntInt64(v.(int))
				}
				failureParam.KafkaParam = &kafkaParam
			}
			if v, ok := failureParamMap["retry_interval"]; ok {
				failureParam.RetryInterval = helper.IntUint64(v.(int))
			}
			if v, ok := failureParamMap["max_retry_attempts"]; ok {
				failureParam.MaxRetryAttempts = helper.IntUint64(v.(int))
			}
			if topicParamMap, ok := helper.InterfaceToMap(failureParamMap, "topic_param"); ok {
				topicParam := ckafka.TopicParam{}
				if v, ok := topicParamMap["resource"]; ok {
					topicParam.Resource = helper.String(v.(string))
				}
				if v, ok := topicParamMap["offset_type"]; ok {
					topicParam.OffsetType = helper.String(v.(string))
				}
				if v, ok := topicParamMap["start_time"]; ok {
					topicParam.StartTime = helper.IntUint64(v.(int))
				}
				if v, ok := topicParamMap["topic_id"]; ok {
					topicParam.TopicId = helper.String(v.(string))
				}
				if v, ok := topicParamMap["compression_type"]; ok {
					topicParam.CompressionType = helper.String(v.(string))
				}
				if v, ok := topicParamMap["use_auto_create_topic"]; ok {
					topicParam.UseAutoCreateTopic = helper.Bool(v.(bool))
				}
				if v, ok := topicParamMap["msg_multiple"]; ok {
					topicParam.MsgMultiple = helper.IntInt64(v.(int))
				}
				failureParam.TopicParam = &topicParam
			}
			if v, ok := failureParamMap["dlq_type"]; ok {
				failureParam.DlqType = helper.String(v.(string))
			}
			transformsParam.FailureParam = &failureParam
		}
		if v, ok := dMap["result"]; ok {
			transformsParam.Result = helper.String(v.(string))
		}
		if v, ok := dMap["source_type"]; ok {
			transformsParam.SourceType = helper.String(v.(string))
		}
		if v, ok := dMap["output_format"]; ok {
			transformsParam.OutputFormat = helper.String(v.(string))
		}
		if rowParamMap, ok := helper.InterfaceToMap(dMap, "row_param"); ok {
			rowParam := ckafka.RowParam{}
			if v, ok := rowParamMap["row_content"]; ok {
				rowParam.RowContent = helper.String(v.(string))
			}
			if v, ok := rowParamMap["key_value_delimiter"]; ok {
				rowParam.KeyValueDelimiter = helper.String(v.(string))
			}
			if v, ok := rowParamMap["entry_delimiter"]; ok {
				rowParam.EntryDelimiter = helper.String(v.(string))
			}
			transformsParam.RowParam = &rowParam
		}
		if v, ok := dMap["keep_metadata"]; ok {
			transformsParam.KeepMetadata = helper.Bool(v.(bool))
		}
		if batchAnalyseMap, ok := helper.InterfaceToMap(dMap, "batch_analyse"); ok {
			batchAnalyseParam := ckafka.BatchAnalyseParam{}
			if v, ok := batchAnalyseMap["format"]; ok {
				batchAnalyseParam.Format = helper.String(v.(string))
			}
			transformsParam.BatchAnalyse = &batchAnalyseParam
		}
		request.TransformsParam = &transformsParam
	}

	if v, ok := d.GetOk("task_id"); ok {
		request.TaskId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().CreateDatahubTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ckafka datahubTask failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.Result.TaskId
	d.SetId(taskId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 2*tccommon.ReadRetryTimeout, time.Second, service.CkafkaDatahubTaskStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCkafkaDatahubTaskRead(d, meta)
}

func resourceTencentCloudCkafkaDatahubTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_datahub_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	taskId := d.Id()

	datahubTask, err := service.DescribeDatahubTask(ctx, taskId)
	if err != nil {
		return err
	}

	if datahubTask == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CkafkaDatahubTask` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if datahubTask.TaskName != nil {
		_ = d.Set("task_name", datahubTask.TaskName)
	}

	if datahubTask.TaskType != nil {
		_ = d.Set("task_type", datahubTask.TaskType)
	}

	if datahubTask.SourceResource != nil {
		sourceResourceMap := map[string]interface{}{}

		if datahubTask.SourceResource.Type != nil {
			sourceResourceMap["type"] = datahubTask.SourceResource.Type
		}

		if datahubTask.SourceResource.KafkaParam != nil {
			kafkaParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.KafkaParam.SelfBuilt != nil {
				kafkaParamMap["self_built"] = datahubTask.SourceResource.KafkaParam.SelfBuilt
			}

			if datahubTask.SourceResource.KafkaParam.Resource != nil {
				kafkaParamMap["resource"] = datahubTask.SourceResource.KafkaParam.Resource
			}

			if datahubTask.SourceResource.KafkaParam.Topic != nil {
				kafkaParamMap["topic"] = datahubTask.SourceResource.KafkaParam.Topic
			}

			if datahubTask.SourceResource.KafkaParam.OffsetType != nil {
				kafkaParamMap["offset_type"] = datahubTask.SourceResource.KafkaParam.OffsetType
			}

			if datahubTask.SourceResource.KafkaParam.StartTime != nil {
				kafkaParamMap["start_time"] = datahubTask.SourceResource.KafkaParam.StartTime
			}

			if datahubTask.SourceResource.KafkaParam.ResourceName != nil {
				kafkaParamMap["resource_name"] = datahubTask.SourceResource.KafkaParam.ResourceName
			}

			if datahubTask.SourceResource.KafkaParam.ZoneId != nil {
				kafkaParamMap["zone_id"] = datahubTask.SourceResource.KafkaParam.ZoneId
			}

			if datahubTask.SourceResource.KafkaParam.TopicId != nil {
				kafkaParamMap["topic_id"] = datahubTask.SourceResource.KafkaParam.TopicId
			}

			if datahubTask.SourceResource.KafkaParam.PartitionNum != nil {
				kafkaParamMap["partition_num"] = datahubTask.SourceResource.KafkaParam.PartitionNum
			}

			if datahubTask.SourceResource.KafkaParam.EnableToleration != nil {
				kafkaParamMap["enable_toleration"] = datahubTask.SourceResource.KafkaParam.EnableToleration
			}

			if datahubTask.SourceResource.KafkaParam.QpsLimit != nil {
				kafkaParamMap["qps_limit"] = datahubTask.SourceResource.KafkaParam.QpsLimit
			}

			if datahubTask.SourceResource.KafkaParam.TableMappings != nil {
				tableMappingsList := []interface{}{}
				for _, tableMappings := range datahubTask.SourceResource.KafkaParam.TableMappings {
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

			if datahubTask.SourceResource.KafkaParam.UseTableMapping != nil {
				kafkaParamMap["use_table_mapping"] = datahubTask.SourceResource.KafkaParam.UseTableMapping
			}

			if datahubTask.SourceResource.KafkaParam.UseAutoCreateTopic != nil {
				kafkaParamMap["use_auto_create_topic"] = datahubTask.SourceResource.KafkaParam.UseAutoCreateTopic
			}

			if datahubTask.SourceResource.KafkaParam.CompressionType != nil {
				kafkaParamMap["compression_type"] = datahubTask.SourceResource.KafkaParam.CompressionType
			}

			if datahubTask.SourceResource.KafkaParam.MsgMultiple != nil {
				kafkaParamMap["msg_multiple"] = datahubTask.SourceResource.KafkaParam.MsgMultiple
			}

			sourceResourceMap["kafka_param"] = []interface{}{kafkaParamMap}
		}

		if datahubTask.SourceResource.EventBusParam != nil {
			eventBusParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.EventBusParam.Type != nil {
				eventBusParamMap["type"] = datahubTask.SourceResource.EventBusParam.Type
			}

			if datahubTask.SourceResource.EventBusParam.SelfBuilt != nil {
				eventBusParamMap["self_built"] = datahubTask.SourceResource.EventBusParam.SelfBuilt
			}

			if datahubTask.SourceResource.EventBusParam.Resource != nil {
				eventBusParamMap["resource"] = datahubTask.SourceResource.EventBusParam.Resource
			}

			if datahubTask.SourceResource.EventBusParam.Namespace != nil {
				eventBusParamMap["namespace"] = datahubTask.SourceResource.EventBusParam.Namespace
			}

			if datahubTask.SourceResource.EventBusParam.FunctionName != nil {
				eventBusParamMap["function_name"] = datahubTask.SourceResource.EventBusParam.FunctionName
			}

			if datahubTask.SourceResource.EventBusParam.Qualifier != nil {
				eventBusParamMap["qualifier"] = datahubTask.SourceResource.EventBusParam.Qualifier
			}

			sourceResourceMap["event_bus_param"] = []interface{}{eventBusParamMap}
		}

		if datahubTask.SourceResource.MongoDBParam != nil {
			mongoDBParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.MongoDBParam.Database != nil {
				mongoDBParamMap["database"] = datahubTask.SourceResource.MongoDBParam.Database
			}

			if datahubTask.SourceResource.MongoDBParam.Collection != nil {
				mongoDBParamMap["collection"] = datahubTask.SourceResource.MongoDBParam.Collection
			}

			if datahubTask.SourceResource.MongoDBParam.CopyExisting != nil {
				mongoDBParamMap["copy_existing"] = datahubTask.SourceResource.MongoDBParam.CopyExisting
			}

			if datahubTask.SourceResource.MongoDBParam.Resource != nil {
				mongoDBParamMap["resource"] = datahubTask.SourceResource.MongoDBParam.Resource
			}

			if datahubTask.SourceResource.MongoDBParam.Ip != nil {
				mongoDBParamMap["ip"] = datahubTask.SourceResource.MongoDBParam.Ip
			}

			if datahubTask.SourceResource.MongoDBParam.Port != nil {
				mongoDBParamMap["port"] = datahubTask.SourceResource.MongoDBParam.Port
			}

			if datahubTask.SourceResource.MongoDBParam.UserName != nil {
				mongoDBParamMap["user_name"] = datahubTask.SourceResource.MongoDBParam.UserName
			}

			if datahubTask.SourceResource.MongoDBParam.Password != nil {
				mongoDBParamMap["password"] = datahubTask.SourceResource.MongoDBParam.Password
			}

			if datahubTask.SourceResource.MongoDBParam.ListeningEvent != nil {
				mongoDBParamMap["listening_event"] = datahubTask.SourceResource.MongoDBParam.ListeningEvent
			}

			if datahubTask.SourceResource.MongoDBParam.ReadPreference != nil {
				mongoDBParamMap["read_preference"] = datahubTask.SourceResource.MongoDBParam.ReadPreference
			}

			if datahubTask.SourceResource.MongoDBParam.Pipeline != nil {
				mongoDBParamMap["pipeline"] = datahubTask.SourceResource.MongoDBParam.Pipeline
			}

			if datahubTask.SourceResource.MongoDBParam.SelfBuilt != nil {
				mongoDBParamMap["self_built"] = datahubTask.SourceResource.MongoDBParam.SelfBuilt
			}

			sourceResourceMap["mongo_db_param"] = []interface{}{mongoDBParamMap}
		}

		if datahubTask.SourceResource.EsParam != nil {
			esParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.EsParam.Resource != nil {
				esParamMap["resource"] = datahubTask.SourceResource.EsParam.Resource
			}

			if datahubTask.SourceResource.EsParam.Port != nil {
				esParamMap["port"] = datahubTask.SourceResource.EsParam.Port
			}

			if datahubTask.SourceResource.EsParam.UserName != nil {
				esParamMap["user_name"] = datahubTask.SourceResource.EsParam.UserName
			}

			if datahubTask.SourceResource.EsParam.Password != nil {
				esParamMap["password"] = datahubTask.SourceResource.EsParam.Password
			}

			if datahubTask.SourceResource.EsParam.SelfBuilt != nil {
				esParamMap["self_built"] = datahubTask.SourceResource.EsParam.SelfBuilt
			}

			if datahubTask.SourceResource.EsParam.ServiceVip != nil {
				esParamMap["service_vip"] = datahubTask.SourceResource.EsParam.ServiceVip
			}

			if datahubTask.SourceResource.EsParam.UniqVpcId != nil {
				esParamMap["uniq_vpc_id"] = datahubTask.SourceResource.EsParam.UniqVpcId
			}

			if datahubTask.SourceResource.EsParam.DropInvalidMessage != nil {
				esParamMap["drop_invalid_message"] = datahubTask.SourceResource.EsParam.DropInvalidMessage
			}

			if datahubTask.SourceResource.EsParam.Index != nil {
				esParamMap["index"] = datahubTask.SourceResource.EsParam.Index
			}

			if datahubTask.SourceResource.EsParam.DateFormat != nil {
				esParamMap["date_format"] = datahubTask.SourceResource.EsParam.DateFormat
			}

			if datahubTask.SourceResource.EsParam.ContentKey != nil {
				esParamMap["content_key"] = datahubTask.SourceResource.EsParam.ContentKey
			}

			if datahubTask.SourceResource.EsParam.DropInvalidJsonMessage != nil {
				esParamMap["drop_invalid_json_message"] = datahubTask.SourceResource.EsParam.DropInvalidJsonMessage
			}

			if datahubTask.SourceResource.EsParam.DocumentIdField != nil {
				esParamMap["document_id_field"] = datahubTask.SourceResource.EsParam.DocumentIdField
			}

			if datahubTask.SourceResource.EsParam.IndexType != nil {
				esParamMap["index_type"] = datahubTask.SourceResource.EsParam.IndexType
			}

			if datahubTask.SourceResource.EsParam.DropCls != nil {
				dropClsMap := map[string]interface{}{}

				if datahubTask.SourceResource.EsParam.DropCls.DropInvalidMessageToCls != nil {
					dropClsMap["drop_invalid_message_to_cls"] = datahubTask.SourceResource.EsParam.DropCls.DropInvalidMessageToCls
				}

				if datahubTask.SourceResource.EsParam.DropCls.DropClsRegion != nil {
					dropClsMap["drop_cls_region"] = datahubTask.SourceResource.EsParam.DropCls.DropClsRegion
				}

				if datahubTask.SourceResource.EsParam.DropCls.DropClsOwneruin != nil {
					dropClsMap["drop_cls_owneruin"] = datahubTask.SourceResource.EsParam.DropCls.DropClsOwneruin
				}

				if datahubTask.SourceResource.EsParam.DropCls.DropClsTopicId != nil {
					dropClsMap["drop_cls_topic_id"] = datahubTask.SourceResource.EsParam.DropCls.DropClsTopicId
				}

				if datahubTask.SourceResource.EsParam.DropCls.DropClsLogSet != nil {
					dropClsMap["drop_cls_log_set"] = datahubTask.SourceResource.EsParam.DropCls.DropClsLogSet
				}

				esParamMap["drop_cls"] = []interface{}{dropClsMap}
			}

			if datahubTask.SourceResource.EsParam.DatabasePrimaryKey != nil {
				esParamMap["database_primary_key"] = datahubTask.SourceResource.EsParam.DatabasePrimaryKey
			}

			if datahubTask.SourceResource.EsParam.DropDlq != nil {
				dropDlqMap := map[string]interface{}{}

				if datahubTask.SourceResource.EsParam.DropDlq.Type != nil {
					dropDlqMap["type"] = datahubTask.SourceResource.EsParam.DropDlq.Type
				}

				if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam != nil {
					kafkaParamMap := map[string]interface{}{}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.SelfBuilt != nil {
						kafkaParamMap["self_built"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.SelfBuilt
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.Resource != nil {
						kafkaParamMap["resource"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.Resource
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.Topic != nil {
						kafkaParamMap["topic"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.Topic
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.OffsetType != nil {
						kafkaParamMap["offset_type"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.OffsetType
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.StartTime != nil {
						kafkaParamMap["start_time"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.StartTime
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.ResourceName != nil {
						kafkaParamMap["resource_name"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.ResourceName
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.ZoneId != nil {
						kafkaParamMap["zone_id"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.ZoneId
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.TopicId != nil {
						kafkaParamMap["topic_id"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.TopicId
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.PartitionNum != nil {
						kafkaParamMap["partition_num"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.PartitionNum
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.EnableToleration != nil {
						kafkaParamMap["enable_toleration"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.EnableToleration
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.QpsLimit != nil {
						kafkaParamMap["qps_limit"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.QpsLimit
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.TableMappings != nil {
						tableMappingsList := []interface{}{}
						for _, tableMappings := range datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.TableMappings {
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

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.UseTableMapping != nil {
						kafkaParamMap["use_table_mapping"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.UseTableMapping
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic != nil {
						kafkaParamMap["use_auto_create_topic"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.CompressionType != nil {
						kafkaParamMap["compression_type"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.CompressionType
					}

					if datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.MsgMultiple != nil {
						kafkaParamMap["msg_multiple"] = datahubTask.SourceResource.EsParam.DropDlq.KafkaParam.MsgMultiple
					}

					dropDlqMap["kafka_param"] = []interface{}{kafkaParamMap}
				}

				if datahubTask.SourceResource.EsParam.DropDlq.RetryInterval != nil {
					dropDlqMap["retry_interval"] = datahubTask.SourceResource.EsParam.DropDlq.RetryInterval
				}

				if datahubTask.SourceResource.EsParam.DropDlq.MaxRetryAttempts != nil {
					dropDlqMap["max_retry_attempts"] = datahubTask.SourceResource.EsParam.DropDlq.MaxRetryAttempts
				}

				if datahubTask.SourceResource.EsParam.DropDlq.TopicParam != nil {
					topicParamMap := map[string]interface{}{}

					if datahubTask.SourceResource.EsParam.DropDlq.TopicParam.Resource != nil {
						topicParamMap["resource"] = datahubTask.SourceResource.EsParam.DropDlq.TopicParam.Resource
					}

					if datahubTask.SourceResource.EsParam.DropDlq.TopicParam.OffsetType != nil {
						topicParamMap["offset_type"] = datahubTask.SourceResource.EsParam.DropDlq.TopicParam.OffsetType
					}

					if datahubTask.SourceResource.EsParam.DropDlq.TopicParam.StartTime != nil {
						topicParamMap["start_time"] = datahubTask.SourceResource.EsParam.DropDlq.TopicParam.StartTime
					}

					if datahubTask.SourceResource.EsParam.DropDlq.TopicParam.TopicId != nil {
						topicParamMap["topic_id"] = datahubTask.SourceResource.EsParam.DropDlq.TopicParam.TopicId
					}

					if datahubTask.SourceResource.EsParam.DropDlq.TopicParam.CompressionType != nil {
						topicParamMap["compression_type"] = datahubTask.SourceResource.EsParam.DropDlq.TopicParam.CompressionType
					}

					if datahubTask.SourceResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic != nil {
						topicParamMap["use_auto_create_topic"] = datahubTask.SourceResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic
					}

					if datahubTask.SourceResource.EsParam.DropDlq.TopicParam.MsgMultiple != nil {
						topicParamMap["msg_multiple"] = datahubTask.SourceResource.EsParam.DropDlq.TopicParam.MsgMultiple
					}

					dropDlqMap["topic_param"] = []interface{}{topicParamMap}
				}

				if datahubTask.SourceResource.EsParam.DropDlq.DlqType != nil {
					dropDlqMap["dlq_type"] = datahubTask.SourceResource.EsParam.DropDlq.DlqType
				}

				esParamMap["drop_dlq"] = []interface{}{dropDlqMap}
			}

			sourceResourceMap["es_param"] = []interface{}{esParamMap}
		}

		if datahubTask.SourceResource.TdwParam != nil {
			tdwParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.TdwParam.Bid != nil {
				tdwParamMap["bid"] = datahubTask.SourceResource.TdwParam.Bid
			}

			if datahubTask.SourceResource.TdwParam.Tid != nil {
				tdwParamMap["tid"] = datahubTask.SourceResource.TdwParam.Tid
			}

			if datahubTask.SourceResource.TdwParam.IsDomestic != nil {
				tdwParamMap["is_domestic"] = datahubTask.SourceResource.TdwParam.IsDomestic
			}

			if datahubTask.SourceResource.TdwParam.TdwHost != nil {
				tdwParamMap["tdw_host"] = datahubTask.SourceResource.TdwParam.TdwHost
			}

			if datahubTask.SourceResource.TdwParam.TdwPort != nil {
				tdwParamMap["tdw_port"] = datahubTask.SourceResource.TdwParam.TdwPort
			}

			sourceResourceMap["tdw_param"] = []interface{}{tdwParamMap}
		}

		if datahubTask.SourceResource.DtsParam != nil {
			dtsParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.DtsParam.Resource != nil {
				dtsParamMap["resource"] = datahubTask.SourceResource.DtsParam.Resource
			}

			if datahubTask.SourceResource.DtsParam.Ip != nil {
				dtsParamMap["ip"] = datahubTask.SourceResource.DtsParam.Ip
			}

			if datahubTask.SourceResource.DtsParam.Port != nil {
				dtsParamMap["port"] = datahubTask.SourceResource.DtsParam.Port
			}

			if datahubTask.SourceResource.DtsParam.Topic != nil {
				dtsParamMap["topic"] = datahubTask.SourceResource.DtsParam.Topic
			}

			if datahubTask.SourceResource.DtsParam.GroupId != nil {
				dtsParamMap["group_id"] = datahubTask.SourceResource.DtsParam.GroupId
			}

			if datahubTask.SourceResource.DtsParam.GroupUser != nil {
				dtsParamMap["group_user"] = datahubTask.SourceResource.DtsParam.GroupUser
			}

			if datahubTask.SourceResource.DtsParam.GroupPassword != nil {
				dtsParamMap["group_password"] = datahubTask.SourceResource.DtsParam.GroupPassword
			}

			if datahubTask.SourceResource.DtsParam.TranSql != nil {
				dtsParamMap["tran_sql"] = datahubTask.SourceResource.DtsParam.TranSql
			}

			sourceResourceMap["dts_param"] = []interface{}{dtsParamMap}
		}

		if datahubTask.SourceResource.ClickHouseParam != nil {
			clickHouseParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.ClickHouseParam.Cluster != nil {
				clickHouseParamMap["cluster"] = datahubTask.SourceResource.ClickHouseParam.Cluster
			}

			if datahubTask.SourceResource.ClickHouseParam.Database != nil {
				clickHouseParamMap["database"] = datahubTask.SourceResource.ClickHouseParam.Database
			}

			if datahubTask.SourceResource.ClickHouseParam.Table != nil {
				clickHouseParamMap["table"] = datahubTask.SourceResource.ClickHouseParam.Table
			}

			if datahubTask.SourceResource.ClickHouseParam.Schema != nil {
				schemaList := []interface{}{}
				for _, schema := range datahubTask.SourceResource.ClickHouseParam.Schema {
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

			if datahubTask.SourceResource.ClickHouseParam.Resource != nil {
				clickHouseParamMap["resource"] = datahubTask.SourceResource.ClickHouseParam.Resource
			}

			if datahubTask.SourceResource.ClickHouseParam.Ip != nil {
				clickHouseParamMap["ip"] = datahubTask.SourceResource.ClickHouseParam.Ip
			}

			if datahubTask.SourceResource.ClickHouseParam.Port != nil {
				clickHouseParamMap["port"] = datahubTask.SourceResource.ClickHouseParam.Port
			}

			if datahubTask.SourceResource.ClickHouseParam.UserName != nil {
				clickHouseParamMap["user_name"] = datahubTask.SourceResource.ClickHouseParam.UserName
			}

			if datahubTask.SourceResource.ClickHouseParam.Password != nil {
				clickHouseParamMap["password"] = datahubTask.SourceResource.ClickHouseParam.Password
			}

			if datahubTask.SourceResource.ClickHouseParam.ServiceVip != nil {
				clickHouseParamMap["service_vip"] = datahubTask.SourceResource.ClickHouseParam.ServiceVip
			}

			if datahubTask.SourceResource.ClickHouseParam.UniqVpcId != nil {
				clickHouseParamMap["uniq_vpc_id"] = datahubTask.SourceResource.ClickHouseParam.UniqVpcId
			}

			if datahubTask.SourceResource.ClickHouseParam.SelfBuilt != nil {
				clickHouseParamMap["self_built"] = datahubTask.SourceResource.ClickHouseParam.SelfBuilt
			}

			if datahubTask.SourceResource.ClickHouseParam.DropInvalidMessage != nil {
				clickHouseParamMap["drop_invalid_message"] = datahubTask.SourceResource.ClickHouseParam.DropInvalidMessage
			}

			if datahubTask.SourceResource.ClickHouseParam.Type != nil {
				clickHouseParamMap["type"] = datahubTask.SourceResource.ClickHouseParam.Type
			}

			if datahubTask.SourceResource.ClickHouseParam.DropCls != nil {
				dropClsMap := map[string]interface{}{}

				if datahubTask.SourceResource.ClickHouseParam.DropCls.DropInvalidMessageToCls != nil {
					dropClsMap["drop_invalid_message_to_cls"] = datahubTask.SourceResource.ClickHouseParam.DropCls.DropInvalidMessageToCls
				}

				if datahubTask.SourceResource.ClickHouseParam.DropCls.DropClsRegion != nil {
					dropClsMap["drop_cls_region"] = datahubTask.SourceResource.ClickHouseParam.DropCls.DropClsRegion
				}

				if datahubTask.SourceResource.ClickHouseParam.DropCls.DropClsOwneruin != nil {
					dropClsMap["drop_cls_owneruin"] = datahubTask.SourceResource.ClickHouseParam.DropCls.DropClsOwneruin
				}

				if datahubTask.SourceResource.ClickHouseParam.DropCls.DropClsTopicId != nil {
					dropClsMap["drop_cls_topic_id"] = datahubTask.SourceResource.ClickHouseParam.DropCls.DropClsTopicId
				}

				if datahubTask.SourceResource.ClickHouseParam.DropCls.DropClsLogSet != nil {
					dropClsMap["drop_cls_log_set"] = datahubTask.SourceResource.ClickHouseParam.DropCls.DropClsLogSet
				}

				clickHouseParamMap["drop_cls"] = []interface{}{dropClsMap}
			}

			sourceResourceMap["click_house_param"] = []interface{}{clickHouseParamMap}
		}

		if datahubTask.SourceResource.ClsParam != nil {
			clsParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.ClsParam.DecodeJson != nil {
				clsParamMap["decode_json"] = datahubTask.SourceResource.ClsParam.DecodeJson
			}

			if datahubTask.SourceResource.ClsParam.Resource != nil {
				clsParamMap["resource"] = datahubTask.SourceResource.ClsParam.Resource
			}

			if datahubTask.SourceResource.ClsParam.LogSet != nil {
				clsParamMap["log_set"] = datahubTask.SourceResource.ClsParam.LogSet
			}

			if datahubTask.SourceResource.ClsParam.ContentKey != nil {
				clsParamMap["content_key"] = datahubTask.SourceResource.ClsParam.ContentKey
			}

			if datahubTask.SourceResource.ClsParam.TimeField != nil {
				clsParamMap["time_field"] = datahubTask.SourceResource.ClsParam.TimeField
			}

			sourceResourceMap["cls_param"] = []interface{}{clsParamMap}
		}

		if datahubTask.SourceResource.CosParam != nil {
			cosParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.CosParam.BucketName != nil {
				cosParamMap["bucket_name"] = datahubTask.SourceResource.CosParam.BucketName
			}

			if datahubTask.SourceResource.CosParam.Region != nil {
				cosParamMap["region"] = datahubTask.SourceResource.CosParam.Region
			}

			if datahubTask.SourceResource.CosParam.ObjectKey != nil {
				cosParamMap["object_key"] = datahubTask.SourceResource.CosParam.ObjectKey
			}

			if datahubTask.SourceResource.CosParam.AggregateBatchSize != nil {
				cosParamMap["aggregate_batch_size"] = datahubTask.SourceResource.CosParam.AggregateBatchSize
			}

			if datahubTask.SourceResource.CosParam.AggregateInterval != nil {
				cosParamMap["aggregate_interval"] = datahubTask.SourceResource.CosParam.AggregateInterval
			}

			if datahubTask.SourceResource.CosParam.FormatOutputType != nil {
				cosParamMap["format_output_type"] = datahubTask.SourceResource.CosParam.FormatOutputType
			}

			if datahubTask.SourceResource.CosParam.ObjectKeyPrefix != nil {
				cosParamMap["object_key_prefix"] = datahubTask.SourceResource.CosParam.ObjectKeyPrefix
			}

			if datahubTask.SourceResource.CosParam.DirectoryTimeFormat != nil {
				cosParamMap["directory_time_format"] = datahubTask.SourceResource.CosParam.DirectoryTimeFormat
			}

			sourceResourceMap["cos_param"] = []interface{}{cosParamMap}
		}

		if datahubTask.SourceResource.MySQLParam != nil {
			mySQLParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.MySQLParam.Database != nil {
				mySQLParamMap["database"] = datahubTask.SourceResource.MySQLParam.Database
			}

			if datahubTask.SourceResource.MySQLParam.Table != nil {
				mySQLParamMap["table"] = datahubTask.SourceResource.MySQLParam.Table
			}

			if datahubTask.SourceResource.MySQLParam.Resource != nil {
				mySQLParamMap["resource"] = datahubTask.SourceResource.MySQLParam.Resource
			}

			if datahubTask.SourceResource.MySQLParam.SnapshotMode != nil {
				mySQLParamMap["snapshot_mode"] = datahubTask.SourceResource.MySQLParam.SnapshotMode
			}

			if datahubTask.SourceResource.MySQLParam.DdlTopic != nil {
				mySQLParamMap["ddl_topic"] = datahubTask.SourceResource.MySQLParam.DdlTopic
			}

			if datahubTask.SourceResource.MySQLParam.DataSourceMonitorMode != nil {
				mySQLParamMap["data_source_monitor_mode"] = datahubTask.SourceResource.MySQLParam.DataSourceMonitorMode
			}

			if datahubTask.SourceResource.MySQLParam.DataSourceMonitorResource != nil {
				mySQLParamMap["data_source_monitor_resource"] = datahubTask.SourceResource.MySQLParam.DataSourceMonitorResource
			}

			if datahubTask.SourceResource.MySQLParam.DataSourceIncrementMode != nil {
				mySQLParamMap["data_source_increment_mode"] = datahubTask.SourceResource.MySQLParam.DataSourceIncrementMode
			}

			if datahubTask.SourceResource.MySQLParam.DataSourceIncrementColumn != nil {
				mySQLParamMap["data_source_increment_column"] = datahubTask.SourceResource.MySQLParam.DataSourceIncrementColumn
			}

			if datahubTask.SourceResource.MySQLParam.DataSourceStartFrom != nil {
				mySQLParamMap["data_source_start_from"] = datahubTask.SourceResource.MySQLParam.DataSourceStartFrom
			}

			if datahubTask.SourceResource.MySQLParam.DataTargetInsertMode != nil {
				mySQLParamMap["data_target_insert_mode"] = datahubTask.SourceResource.MySQLParam.DataTargetInsertMode
			}

			if datahubTask.SourceResource.MySQLParam.DataTargetPrimaryKeyField != nil {
				mySQLParamMap["data_target_primary_key_field"] = datahubTask.SourceResource.MySQLParam.DataTargetPrimaryKeyField
			}

			if datahubTask.SourceResource.MySQLParam.DataTargetRecordMapping != nil {
				dataTargetRecordMappingList := []interface{}{}
				for _, dataTargetRecordMapping := range datahubTask.SourceResource.MySQLParam.DataTargetRecordMapping {
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

			if datahubTask.SourceResource.MySQLParam.TopicRegex != nil {
				mySQLParamMap["topic_regex"] = datahubTask.SourceResource.MySQLParam.TopicRegex
			}

			if datahubTask.SourceResource.MySQLParam.TopicReplacement != nil {
				mySQLParamMap["topic_replacement"] = datahubTask.SourceResource.MySQLParam.TopicReplacement
			}

			if datahubTask.SourceResource.MySQLParam.KeyColumns != nil {
				mySQLParamMap["key_columns"] = datahubTask.SourceResource.MySQLParam.KeyColumns
			}

			if datahubTask.SourceResource.MySQLParam.DropInvalidMessage != nil {
				mySQLParamMap["drop_invalid_message"] = datahubTask.SourceResource.MySQLParam.DropInvalidMessage
			}

			if datahubTask.SourceResource.MySQLParam.DropCls != nil {
				dropClsMap := map[string]interface{}{}

				if datahubTask.SourceResource.MySQLParam.DropCls.DropInvalidMessageToCls != nil {
					dropClsMap["drop_invalid_message_to_cls"] = datahubTask.SourceResource.MySQLParam.DropCls.DropInvalidMessageToCls
				}

				if datahubTask.SourceResource.MySQLParam.DropCls.DropClsRegion != nil {
					dropClsMap["drop_cls_region"] = datahubTask.SourceResource.MySQLParam.DropCls.DropClsRegion
				}

				if datahubTask.SourceResource.MySQLParam.DropCls.DropClsOwneruin != nil {
					dropClsMap["drop_cls_owneruin"] = datahubTask.SourceResource.MySQLParam.DropCls.DropClsOwneruin
				}

				if datahubTask.SourceResource.MySQLParam.DropCls.DropClsTopicId != nil {
					dropClsMap["drop_cls_topic_id"] = datahubTask.SourceResource.MySQLParam.DropCls.DropClsTopicId
				}

				if datahubTask.SourceResource.MySQLParam.DropCls.DropClsLogSet != nil {
					dropClsMap["drop_cls_log_set"] = datahubTask.SourceResource.MySQLParam.DropCls.DropClsLogSet
				}

				mySQLParamMap["drop_cls"] = []interface{}{dropClsMap}
			}

			if datahubTask.SourceResource.MySQLParam.OutputFormat != nil {
				mySQLParamMap["output_format"] = datahubTask.SourceResource.MySQLParam.OutputFormat
			}

			if datahubTask.SourceResource.MySQLParam.IsTablePrefix != nil {
				mySQLParamMap["is_table_prefix"] = datahubTask.SourceResource.MySQLParam.IsTablePrefix
			}

			if datahubTask.SourceResource.MySQLParam.IncludeContentChanges != nil {
				mySQLParamMap["include_content_changes"] = datahubTask.SourceResource.MySQLParam.IncludeContentChanges
			}

			if datahubTask.SourceResource.MySQLParam.IncludeQuery != nil {
				mySQLParamMap["include_query"] = datahubTask.SourceResource.MySQLParam.IncludeQuery
			}

			if datahubTask.SourceResource.MySQLParam.RecordWithSchema != nil {
				mySQLParamMap["record_with_schema"] = datahubTask.SourceResource.MySQLParam.RecordWithSchema
			}

			if datahubTask.SourceResource.MySQLParam.SignalDatabase != nil {
				mySQLParamMap["signal_database"] = datahubTask.SourceResource.MySQLParam.SignalDatabase
			}

			if datahubTask.SourceResource.MySQLParam.IsTableRegular != nil {
				mySQLParamMap["is_table_regular"] = datahubTask.SourceResource.MySQLParam.IsTableRegular
			}

			sourceResourceMap["my_sql_param"] = []interface{}{mySQLParamMap}
		}

		if datahubTask.SourceResource.PostgreSQLParam != nil {
			postgreSQLParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.PostgreSQLParam.Database != nil {
				postgreSQLParamMap["database"] = datahubTask.SourceResource.PostgreSQLParam.Database
			}

			if datahubTask.SourceResource.PostgreSQLParam.Table != nil {
				postgreSQLParamMap["table"] = datahubTask.SourceResource.PostgreSQLParam.Table
			}

			if datahubTask.SourceResource.PostgreSQLParam.Resource != nil {
				postgreSQLParamMap["resource"] = datahubTask.SourceResource.PostgreSQLParam.Resource
			}

			if datahubTask.SourceResource.PostgreSQLParam.PluginName != nil {
				postgreSQLParamMap["plugin_name"] = datahubTask.SourceResource.PostgreSQLParam.PluginName
			}

			if datahubTask.SourceResource.PostgreSQLParam.SnapshotMode != nil {
				postgreSQLParamMap["snapshot_mode"] = datahubTask.SourceResource.PostgreSQLParam.SnapshotMode
			}

			if datahubTask.SourceResource.PostgreSQLParam.DataFormat != nil {
				postgreSQLParamMap["data_format"] = datahubTask.SourceResource.PostgreSQLParam.DataFormat
			}

			if datahubTask.SourceResource.PostgreSQLParam.DataTargetInsertMode != nil {
				postgreSQLParamMap["data_target_insert_mode"] = datahubTask.SourceResource.PostgreSQLParam.DataTargetInsertMode
			}

			if datahubTask.SourceResource.PostgreSQLParam.DataTargetPrimaryKeyField != nil {
				postgreSQLParamMap["data_target_primary_key_field"] = datahubTask.SourceResource.PostgreSQLParam.DataTargetPrimaryKeyField
			}

			if datahubTask.SourceResource.PostgreSQLParam.DataTargetRecordMapping != nil {
				dataTargetRecordMappingList := []interface{}{}
				for _, dataTargetRecordMapping := range datahubTask.SourceResource.PostgreSQLParam.DataTargetRecordMapping {
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

			if datahubTask.SourceResource.PostgreSQLParam.DropInvalidMessage != nil {
				postgreSQLParamMap["drop_invalid_message"] = datahubTask.SourceResource.PostgreSQLParam.DropInvalidMessage
			}

			if datahubTask.SourceResource.PostgreSQLParam.IsTableRegular != nil {
				postgreSQLParamMap["is_table_regular"] = datahubTask.SourceResource.PostgreSQLParam.IsTableRegular
			}

			if datahubTask.SourceResource.PostgreSQLParam.KeyColumns != nil {
				postgreSQLParamMap["key_columns"] = datahubTask.SourceResource.PostgreSQLParam.KeyColumns
			}

			if datahubTask.SourceResource.PostgreSQLParam.RecordWithSchema != nil {
				postgreSQLParamMap["record_with_schema"] = datahubTask.SourceResource.PostgreSQLParam.RecordWithSchema
			}

			sourceResourceMap["postgre_sql_param"] = []interface{}{postgreSQLParamMap}
		}

		if datahubTask.SourceResource.TopicParam != nil {
			topicParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.TopicParam.Resource != nil {
				topicParamMap["resource"] = datahubTask.SourceResource.TopicParam.Resource
			}

			if datahubTask.SourceResource.TopicParam.OffsetType != nil {
				topicParamMap["offset_type"] = datahubTask.SourceResource.TopicParam.OffsetType
			}

			if datahubTask.SourceResource.TopicParam.StartTime != nil {
				topicParamMap["start_time"] = datahubTask.SourceResource.TopicParam.StartTime
			}

			if datahubTask.SourceResource.TopicParam.TopicId != nil {
				topicParamMap["topic_id"] = datahubTask.SourceResource.TopicParam.TopicId
			}

			if datahubTask.SourceResource.TopicParam.CompressionType != nil {
				topicParamMap["compression_type"] = datahubTask.SourceResource.TopicParam.CompressionType
			}

			if datahubTask.SourceResource.TopicParam.UseAutoCreateTopic != nil {
				topicParamMap["use_auto_create_topic"] = datahubTask.SourceResource.TopicParam.UseAutoCreateTopic
			}

			if datahubTask.SourceResource.TopicParam.MsgMultiple != nil {
				topicParamMap["msg_multiple"] = datahubTask.SourceResource.TopicParam.MsgMultiple
			}

			sourceResourceMap["topic_param"] = []interface{}{topicParamMap}
		}

		if datahubTask.SourceResource.MariaDBParam != nil {
			mariaDBParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.MariaDBParam.Database != nil {
				mariaDBParamMap["database"] = datahubTask.SourceResource.MariaDBParam.Database
			}

			if datahubTask.SourceResource.MariaDBParam.Table != nil {
				mariaDBParamMap["table"] = datahubTask.SourceResource.MariaDBParam.Table
			}

			if datahubTask.SourceResource.MariaDBParam.Resource != nil {
				mariaDBParamMap["resource"] = datahubTask.SourceResource.MariaDBParam.Resource
			}

			if datahubTask.SourceResource.MariaDBParam.SnapshotMode != nil {
				mariaDBParamMap["snapshot_mode"] = datahubTask.SourceResource.MariaDBParam.SnapshotMode
			}

			if datahubTask.SourceResource.MariaDBParam.KeyColumns != nil {
				mariaDBParamMap["key_columns"] = datahubTask.SourceResource.MariaDBParam.KeyColumns
			}

			if datahubTask.SourceResource.MariaDBParam.IsTablePrefix != nil {
				mariaDBParamMap["is_table_prefix"] = datahubTask.SourceResource.MariaDBParam.IsTablePrefix
			}

			if datahubTask.SourceResource.MariaDBParam.OutputFormat != nil {
				mariaDBParamMap["output_format"] = datahubTask.SourceResource.MariaDBParam.OutputFormat
			}

			if datahubTask.SourceResource.MariaDBParam.IncludeContentChanges != nil {
				mariaDBParamMap["include_content_changes"] = datahubTask.SourceResource.MariaDBParam.IncludeContentChanges
			}

			if datahubTask.SourceResource.MariaDBParam.IncludeQuery != nil {
				mariaDBParamMap["include_query"] = datahubTask.SourceResource.MariaDBParam.IncludeQuery
			}

			if datahubTask.SourceResource.MariaDBParam.RecordWithSchema != nil {
				mariaDBParamMap["record_with_schema"] = datahubTask.SourceResource.MariaDBParam.RecordWithSchema
			}

			sourceResourceMap["maria_db_param"] = []interface{}{mariaDBParamMap}
		}

		if datahubTask.SourceResource.SQLServerParam != nil {
			sQLServerParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.SQLServerParam.Database != nil {
				sQLServerParamMap["database"] = datahubTask.SourceResource.SQLServerParam.Database
			}

			if datahubTask.SourceResource.SQLServerParam.Table != nil {
				sQLServerParamMap["table"] = datahubTask.SourceResource.SQLServerParam.Table
			}

			if datahubTask.SourceResource.SQLServerParam.Resource != nil {
				sQLServerParamMap["resource"] = datahubTask.SourceResource.SQLServerParam.Resource
			}

			if datahubTask.SourceResource.SQLServerParam.SnapshotMode != nil {
				sQLServerParamMap["snapshot_mode"] = datahubTask.SourceResource.SQLServerParam.SnapshotMode
			}

			sourceResourceMap["sql_server_param"] = []interface{}{sQLServerParamMap}
		}

		if datahubTask.SourceResource.CtsdbParam != nil {
			ctsdbParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.CtsdbParam.Resource != nil {
				ctsdbParamMap["resource"] = datahubTask.SourceResource.CtsdbParam.Resource
			}

			if datahubTask.SourceResource.CtsdbParam.CtsdbMetric != nil {
				ctsdbParamMap["ctsdb_metric"] = datahubTask.SourceResource.CtsdbParam.CtsdbMetric
			}

			sourceResourceMap["ctsdb_param"] = []interface{}{ctsdbParamMap}
		}

		if datahubTask.SourceResource.ScfParam != nil {
			scfParamMap := map[string]interface{}{}

			if datahubTask.SourceResource.ScfParam.FunctionName != nil {
				scfParamMap["function_name"] = datahubTask.SourceResource.ScfParam.FunctionName
			}

			if datahubTask.SourceResource.ScfParam.Namespace != nil {
				scfParamMap["namespace"] = datahubTask.SourceResource.ScfParam.Namespace
			}

			if datahubTask.SourceResource.ScfParam.Qualifier != nil {
				scfParamMap["qualifier"] = datahubTask.SourceResource.ScfParam.Qualifier
			}

			if datahubTask.SourceResource.ScfParam.BatchSize != nil {
				scfParamMap["batch_size"] = datahubTask.SourceResource.ScfParam.BatchSize
			}

			if datahubTask.SourceResource.ScfParam.MaxRetries != nil {
				scfParamMap["max_retries"] = datahubTask.SourceResource.ScfParam.MaxRetries
			}

			sourceResourceMap["scf_param"] = []interface{}{scfParamMap}
		}

		_ = d.Set("source_resource", []interface{}{sourceResourceMap})
	}

	if datahubTask.TargetResource != nil {
		targetResourceMap := map[string]interface{}{}

		if datahubTask.TargetResource.Type != nil {
			targetResourceMap["type"] = datahubTask.TargetResource.Type
		}

		if datahubTask.TargetResource.KafkaParam != nil {
			kafkaParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.KafkaParam.SelfBuilt != nil {
				kafkaParamMap["self_built"] = datahubTask.TargetResource.KafkaParam.SelfBuilt
			}

			if datahubTask.TargetResource.KafkaParam.Resource != nil {
				kafkaParamMap["resource"] = datahubTask.TargetResource.KafkaParam.Resource
			}

			if datahubTask.TargetResource.KafkaParam.Topic != nil {
				kafkaParamMap["topic"] = datahubTask.TargetResource.KafkaParam.Topic
			}

			if datahubTask.TargetResource.KafkaParam.OffsetType != nil {
				kafkaParamMap["offset_type"] = datahubTask.TargetResource.KafkaParam.OffsetType
			}

			if datahubTask.TargetResource.KafkaParam.StartTime != nil {
				kafkaParamMap["start_time"] = datahubTask.TargetResource.KafkaParam.StartTime
			}

			if datahubTask.TargetResource.KafkaParam.ResourceName != nil {
				kafkaParamMap["resource_name"] = datahubTask.TargetResource.KafkaParam.ResourceName
			}

			if datahubTask.TargetResource.KafkaParam.ZoneId != nil {
				kafkaParamMap["zone_id"] = datahubTask.TargetResource.KafkaParam.ZoneId
			}

			if datahubTask.TargetResource.KafkaParam.TopicId != nil {
				kafkaParamMap["topic_id"] = datahubTask.TargetResource.KafkaParam.TopicId
			}

			if datahubTask.TargetResource.KafkaParam.PartitionNum != nil {
				kafkaParamMap["partition_num"] = datahubTask.TargetResource.KafkaParam.PartitionNum
			}

			if datahubTask.TargetResource.KafkaParam.EnableToleration != nil {
				kafkaParamMap["enable_toleration"] = datahubTask.TargetResource.KafkaParam.EnableToleration
			}

			if datahubTask.TargetResource.KafkaParam.QpsLimit != nil {
				kafkaParamMap["qps_limit"] = datahubTask.TargetResource.KafkaParam.QpsLimit
			}

			if datahubTask.TargetResource.KafkaParam.TableMappings != nil {
				tableMappingsList := []interface{}{}
				for _, tableMappings := range datahubTask.TargetResource.KafkaParam.TableMappings {
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

			if datahubTask.TargetResource.KafkaParam.UseTableMapping != nil {
				kafkaParamMap["use_table_mapping"] = datahubTask.TargetResource.KafkaParam.UseTableMapping
			}

			if datahubTask.TargetResource.KafkaParam.UseAutoCreateTopic != nil {
				kafkaParamMap["use_auto_create_topic"] = datahubTask.TargetResource.KafkaParam.UseAutoCreateTopic
			}

			if datahubTask.TargetResource.KafkaParam.CompressionType != nil {
				kafkaParamMap["compression_type"] = datahubTask.TargetResource.KafkaParam.CompressionType
			}

			if datahubTask.TargetResource.KafkaParam.MsgMultiple != nil {
				kafkaParamMap["msg_multiple"] = datahubTask.TargetResource.KafkaParam.MsgMultiple
			}

			targetResourceMap["kafka_param"] = []interface{}{kafkaParamMap}
		}

		if datahubTask.TargetResource.EventBusParam != nil {
			eventBusParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.EventBusParam.Type != nil {
				eventBusParamMap["type"] = datahubTask.TargetResource.EventBusParam.Type
			}

			if datahubTask.TargetResource.EventBusParam.SelfBuilt != nil {
				eventBusParamMap["self_built"] = datahubTask.TargetResource.EventBusParam.SelfBuilt
			}

			if datahubTask.TargetResource.EventBusParam.Resource != nil {
				eventBusParamMap["resource"] = datahubTask.TargetResource.EventBusParam.Resource
			}

			if datahubTask.TargetResource.EventBusParam.Namespace != nil {
				eventBusParamMap["namespace"] = datahubTask.TargetResource.EventBusParam.Namespace
			}

			if datahubTask.TargetResource.EventBusParam.FunctionName != nil {
				eventBusParamMap["function_name"] = datahubTask.TargetResource.EventBusParam.FunctionName
			}

			if datahubTask.TargetResource.EventBusParam.Qualifier != nil {
				eventBusParamMap["qualifier"] = datahubTask.TargetResource.EventBusParam.Qualifier
			}

			targetResourceMap["event_bus_param"] = []interface{}{eventBusParamMap}
		}

		if datahubTask.TargetResource.MongoDBParam != nil {
			mongoDBParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.MongoDBParam.Database != nil {
				mongoDBParamMap["database"] = datahubTask.TargetResource.MongoDBParam.Database
			}

			if datahubTask.TargetResource.MongoDBParam.Collection != nil {
				mongoDBParamMap["collection"] = datahubTask.TargetResource.MongoDBParam.Collection
			}

			if datahubTask.TargetResource.MongoDBParam.CopyExisting != nil {
				mongoDBParamMap["copy_existing"] = datahubTask.TargetResource.MongoDBParam.CopyExisting
			}

			if datahubTask.TargetResource.MongoDBParam.Resource != nil {
				mongoDBParamMap["resource"] = datahubTask.TargetResource.MongoDBParam.Resource
			}

			if datahubTask.TargetResource.MongoDBParam.Ip != nil {
				mongoDBParamMap["ip"] = datahubTask.TargetResource.MongoDBParam.Ip
			}

			if datahubTask.TargetResource.MongoDBParam.Port != nil {
				mongoDBParamMap["port"] = datahubTask.TargetResource.MongoDBParam.Port
			}

			if datahubTask.TargetResource.MongoDBParam.UserName != nil {
				mongoDBParamMap["user_name"] = datahubTask.TargetResource.MongoDBParam.UserName
			}

			if datahubTask.TargetResource.MongoDBParam.Password != nil {
				mongoDBParamMap["password"] = datahubTask.TargetResource.MongoDBParam.Password
			}

			if datahubTask.TargetResource.MongoDBParam.ListeningEvent != nil {
				mongoDBParamMap["listening_event"] = datahubTask.TargetResource.MongoDBParam.ListeningEvent
			}

			if datahubTask.TargetResource.MongoDBParam.ReadPreference != nil {
				mongoDBParamMap["read_preference"] = datahubTask.TargetResource.MongoDBParam.ReadPreference
			}

			if datahubTask.TargetResource.MongoDBParam.Pipeline != nil {
				mongoDBParamMap["pipeline"] = datahubTask.TargetResource.MongoDBParam.Pipeline
			}

			if datahubTask.TargetResource.MongoDBParam.SelfBuilt != nil {
				mongoDBParamMap["self_built"] = datahubTask.TargetResource.MongoDBParam.SelfBuilt
			}

			targetResourceMap["mongo_db_param"] = []interface{}{mongoDBParamMap}
		}

		if datahubTask.TargetResource.EsParam != nil {
			esParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.EsParam.Resource != nil {
				esParamMap["resource"] = datahubTask.TargetResource.EsParam.Resource
			}

			if datahubTask.TargetResource.EsParam.Port != nil {
				esParamMap["port"] = datahubTask.TargetResource.EsParam.Port
			}

			if datahubTask.TargetResource.EsParam.UserName != nil {
				esParamMap["user_name"] = datahubTask.TargetResource.EsParam.UserName
			}

			if datahubTask.TargetResource.EsParam.Password != nil {
				esParamMap["password"] = datahubTask.TargetResource.EsParam.Password
			}

			if datahubTask.TargetResource.EsParam.SelfBuilt != nil {
				esParamMap["self_built"] = datahubTask.TargetResource.EsParam.SelfBuilt
			}

			if datahubTask.TargetResource.EsParam.ServiceVip != nil {
				esParamMap["service_vip"] = datahubTask.TargetResource.EsParam.ServiceVip
			}

			if datahubTask.TargetResource.EsParam.UniqVpcId != nil {
				esParamMap["uniq_vpc_id"] = datahubTask.TargetResource.EsParam.UniqVpcId
			}

			if datahubTask.TargetResource.EsParam.DropInvalidMessage != nil {
				esParamMap["drop_invalid_message"] = datahubTask.TargetResource.EsParam.DropInvalidMessage
			}

			if datahubTask.TargetResource.EsParam.Index != nil {
				esParamMap["index"] = datahubTask.TargetResource.EsParam.Index
			}

			if datahubTask.TargetResource.EsParam.DateFormat != nil {
				esParamMap["date_format"] = datahubTask.TargetResource.EsParam.DateFormat
			}

			if datahubTask.TargetResource.EsParam.ContentKey != nil {
				esParamMap["content_key"] = datahubTask.TargetResource.EsParam.ContentKey
			}

			if datahubTask.TargetResource.EsParam.DropInvalidJsonMessage != nil {
				esParamMap["drop_invalid_json_message"] = datahubTask.TargetResource.EsParam.DropInvalidJsonMessage
			}

			if datahubTask.TargetResource.EsParam.DocumentIdField != nil {
				esParamMap["document_id_field"] = datahubTask.TargetResource.EsParam.DocumentIdField
			}

			if datahubTask.TargetResource.EsParam.IndexType != nil {
				esParamMap["index_type"] = datahubTask.TargetResource.EsParam.IndexType
			}

			if datahubTask.TargetResource.EsParam.DropCls != nil {
				dropClsMap := map[string]interface{}{}

				if datahubTask.TargetResource.EsParam.DropCls.DropInvalidMessageToCls != nil {
					dropClsMap["drop_invalid_message_to_cls"] = datahubTask.TargetResource.EsParam.DropCls.DropInvalidMessageToCls
				}

				if datahubTask.TargetResource.EsParam.DropCls.DropClsRegion != nil {
					dropClsMap["drop_cls_region"] = datahubTask.TargetResource.EsParam.DropCls.DropClsRegion
				}

				if datahubTask.TargetResource.EsParam.DropCls.DropClsOwneruin != nil {
					dropClsMap["drop_cls_owneruin"] = datahubTask.TargetResource.EsParam.DropCls.DropClsOwneruin
				}

				if datahubTask.TargetResource.EsParam.DropCls.DropClsTopicId != nil {
					dropClsMap["drop_cls_topic_id"] = datahubTask.TargetResource.EsParam.DropCls.DropClsTopicId
				}

				if datahubTask.TargetResource.EsParam.DropCls.DropClsLogSet != nil {
					dropClsMap["drop_cls_log_set"] = datahubTask.TargetResource.EsParam.DropCls.DropClsLogSet
				}

				esParamMap["drop_cls"] = []interface{}{dropClsMap}
			}

			if datahubTask.TargetResource.EsParam.DatabasePrimaryKey != nil {
				esParamMap["database_primary_key"] = datahubTask.TargetResource.EsParam.DatabasePrimaryKey
			}

			if datahubTask.TargetResource.EsParam.DropDlq != nil {
				dropDlqMap := map[string]interface{}{}

				if datahubTask.TargetResource.EsParam.DropDlq.Type != nil {
					dropDlqMap["type"] = datahubTask.TargetResource.EsParam.DropDlq.Type
				}

				if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam != nil {
					kafkaParamMap := map[string]interface{}{}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.SelfBuilt != nil {
						kafkaParamMap["self_built"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.SelfBuilt
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.Resource != nil {
						kafkaParamMap["resource"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.Resource
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.Topic != nil {
						kafkaParamMap["topic"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.Topic
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.OffsetType != nil {
						kafkaParamMap["offset_type"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.OffsetType
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.StartTime != nil {
						kafkaParamMap["start_time"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.StartTime
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.ResourceName != nil {
						kafkaParamMap["resource_name"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.ResourceName
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.ZoneId != nil {
						kafkaParamMap["zone_id"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.ZoneId
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.TopicId != nil {
						kafkaParamMap["topic_id"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.TopicId
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.PartitionNum != nil {
						kafkaParamMap["partition_num"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.PartitionNum
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.EnableToleration != nil {
						kafkaParamMap["enable_toleration"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.EnableToleration
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.QpsLimit != nil {
						kafkaParamMap["qps_limit"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.QpsLimit
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.TableMappings != nil {
						tableMappingsList := []interface{}{}
						for _, tableMappings := range datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.TableMappings {
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

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.UseTableMapping != nil {
						kafkaParamMap["use_table_mapping"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.UseTableMapping
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic != nil {
						kafkaParamMap["use_auto_create_topic"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.CompressionType != nil {
						kafkaParamMap["compression_type"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.CompressionType
					}

					if datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.MsgMultiple != nil {
						kafkaParamMap["msg_multiple"] = datahubTask.TargetResource.EsParam.DropDlq.KafkaParam.MsgMultiple
					}

					dropDlqMap["kafka_param"] = []interface{}{kafkaParamMap}
				}

				if datahubTask.TargetResource.EsParam.DropDlq.RetryInterval != nil {
					dropDlqMap["retry_interval"] = datahubTask.TargetResource.EsParam.DropDlq.RetryInterval
				}

				if datahubTask.TargetResource.EsParam.DropDlq.MaxRetryAttempts != nil {
					dropDlqMap["max_retry_attempts"] = datahubTask.TargetResource.EsParam.DropDlq.MaxRetryAttempts
				}

				if datahubTask.TargetResource.EsParam.DropDlq.TopicParam != nil {
					topicParamMap := map[string]interface{}{}

					if datahubTask.TargetResource.EsParam.DropDlq.TopicParam.Resource != nil {
						topicParamMap["resource"] = datahubTask.TargetResource.EsParam.DropDlq.TopicParam.Resource
					}

					if datahubTask.TargetResource.EsParam.DropDlq.TopicParam.OffsetType != nil {
						topicParamMap["offset_type"] = datahubTask.TargetResource.EsParam.DropDlq.TopicParam.OffsetType
					}

					if datahubTask.TargetResource.EsParam.DropDlq.TopicParam.StartTime != nil {
						topicParamMap["start_time"] = datahubTask.TargetResource.EsParam.DropDlq.TopicParam.StartTime
					}

					if datahubTask.TargetResource.EsParam.DropDlq.TopicParam.TopicId != nil {
						topicParamMap["topic_id"] = datahubTask.TargetResource.EsParam.DropDlq.TopicParam.TopicId
					}

					if datahubTask.TargetResource.EsParam.DropDlq.TopicParam.CompressionType != nil {
						topicParamMap["compression_type"] = datahubTask.TargetResource.EsParam.DropDlq.TopicParam.CompressionType
					}

					if datahubTask.TargetResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic != nil {
						topicParamMap["use_auto_create_topic"] = datahubTask.TargetResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic
					}

					if datahubTask.TargetResource.EsParam.DropDlq.TopicParam.MsgMultiple != nil {
						topicParamMap["msg_multiple"] = datahubTask.TargetResource.EsParam.DropDlq.TopicParam.MsgMultiple
					}

					dropDlqMap["topic_param"] = []interface{}{topicParamMap}
				}

				if datahubTask.TargetResource.EsParam.DropDlq.DlqType != nil {
					dropDlqMap["dlq_type"] = datahubTask.TargetResource.EsParam.DropDlq.DlqType
				}

				esParamMap["drop_dlq"] = []interface{}{dropDlqMap}
			}

			targetResourceMap["es_param"] = []interface{}{esParamMap}
		}

		if datahubTask.TargetResource.TdwParam != nil {
			tdwParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.TdwParam.Bid != nil {
				tdwParamMap["bid"] = datahubTask.TargetResource.TdwParam.Bid
			}

			if datahubTask.TargetResource.TdwParam.Tid != nil {
				tdwParamMap["tid"] = datahubTask.TargetResource.TdwParam.Tid
			}

			if datahubTask.TargetResource.TdwParam.IsDomestic != nil {
				tdwParamMap["is_domestic"] = datahubTask.TargetResource.TdwParam.IsDomestic
			}

			if datahubTask.TargetResource.TdwParam.TdwHost != nil {
				tdwParamMap["tdw_host"] = datahubTask.TargetResource.TdwParam.TdwHost
			}

			if datahubTask.TargetResource.TdwParam.TdwPort != nil {
				tdwParamMap["tdw_port"] = datahubTask.TargetResource.TdwParam.TdwPort
			}

			targetResourceMap["tdw_param"] = []interface{}{tdwParamMap}
		}

		if datahubTask.TargetResource.DtsParam != nil {
			dtsParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.DtsParam.Resource != nil {
				dtsParamMap["resource"] = datahubTask.TargetResource.DtsParam.Resource
			}

			if datahubTask.TargetResource.DtsParam.Ip != nil {
				dtsParamMap["ip"] = datahubTask.TargetResource.DtsParam.Ip
			}

			if datahubTask.TargetResource.DtsParam.Port != nil {
				dtsParamMap["port"] = datahubTask.TargetResource.DtsParam.Port
			}

			if datahubTask.TargetResource.DtsParam.Topic != nil {
				dtsParamMap["topic"] = datahubTask.TargetResource.DtsParam.Topic
			}

			if datahubTask.TargetResource.DtsParam.GroupId != nil {
				dtsParamMap["group_id"] = datahubTask.TargetResource.DtsParam.GroupId
			}

			if datahubTask.TargetResource.DtsParam.GroupUser != nil {
				dtsParamMap["group_user"] = datahubTask.TargetResource.DtsParam.GroupUser
			}

			if datahubTask.TargetResource.DtsParam.GroupPassword != nil {
				dtsParamMap["group_password"] = datahubTask.TargetResource.DtsParam.GroupPassword
			}

			if datahubTask.TargetResource.DtsParam.TranSql != nil {
				dtsParamMap["tran_sql"] = datahubTask.TargetResource.DtsParam.TranSql
			}

			targetResourceMap["dts_param"] = []interface{}{dtsParamMap}
		}

		if datahubTask.TargetResource.ClickHouseParam != nil {
			clickHouseParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.ClickHouseParam.Cluster != nil {
				clickHouseParamMap["cluster"] = datahubTask.TargetResource.ClickHouseParam.Cluster
			}

			if datahubTask.TargetResource.ClickHouseParam.Database != nil {
				clickHouseParamMap["database"] = datahubTask.TargetResource.ClickHouseParam.Database
			}

			if datahubTask.TargetResource.ClickHouseParam.Table != nil {
				clickHouseParamMap["table"] = datahubTask.TargetResource.ClickHouseParam.Table
			}

			if datahubTask.TargetResource.ClickHouseParam.Schema != nil {
				schemaList := []interface{}{}
				for _, schema := range datahubTask.TargetResource.ClickHouseParam.Schema {
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

			if datahubTask.TargetResource.ClickHouseParam.Resource != nil {
				clickHouseParamMap["resource"] = datahubTask.TargetResource.ClickHouseParam.Resource
			}

			if datahubTask.TargetResource.ClickHouseParam.Ip != nil {
				clickHouseParamMap["ip"] = datahubTask.TargetResource.ClickHouseParam.Ip
			}

			if datahubTask.TargetResource.ClickHouseParam.Port != nil {
				clickHouseParamMap["port"] = datahubTask.TargetResource.ClickHouseParam.Port
			}

			if datahubTask.TargetResource.ClickHouseParam.UserName != nil {
				clickHouseParamMap["user_name"] = datahubTask.TargetResource.ClickHouseParam.UserName
			}

			if datahubTask.TargetResource.ClickHouseParam.Password != nil {
				clickHouseParamMap["password"] = datahubTask.TargetResource.ClickHouseParam.Password
			}

			if datahubTask.TargetResource.ClickHouseParam.ServiceVip != nil {
				clickHouseParamMap["service_vip"] = datahubTask.TargetResource.ClickHouseParam.ServiceVip
			}

			if datahubTask.TargetResource.ClickHouseParam.UniqVpcId != nil {
				clickHouseParamMap["uniq_vpc_id"] = datahubTask.TargetResource.ClickHouseParam.UniqVpcId
			}

			if datahubTask.TargetResource.ClickHouseParam.SelfBuilt != nil {
				clickHouseParamMap["self_built"] = datahubTask.TargetResource.ClickHouseParam.SelfBuilt
			}

			if datahubTask.TargetResource.ClickHouseParam.DropInvalidMessage != nil {
				clickHouseParamMap["drop_invalid_message"] = datahubTask.TargetResource.ClickHouseParam.DropInvalidMessage
			}

			if datahubTask.TargetResource.ClickHouseParam.Type != nil {
				clickHouseParamMap["type"] = datahubTask.TargetResource.ClickHouseParam.Type
			}

			if datahubTask.TargetResource.ClickHouseParam.DropCls != nil {
				dropClsMap := map[string]interface{}{}

				if datahubTask.TargetResource.ClickHouseParam.DropCls.DropInvalidMessageToCls != nil {
					dropClsMap["drop_invalid_message_to_cls"] = datahubTask.TargetResource.ClickHouseParam.DropCls.DropInvalidMessageToCls
				}

				if datahubTask.TargetResource.ClickHouseParam.DropCls.DropClsRegion != nil {
					dropClsMap["drop_cls_region"] = datahubTask.TargetResource.ClickHouseParam.DropCls.DropClsRegion
				}

				if datahubTask.TargetResource.ClickHouseParam.DropCls.DropClsOwneruin != nil {
					dropClsMap["drop_cls_owneruin"] = datahubTask.TargetResource.ClickHouseParam.DropCls.DropClsOwneruin
				}

				if datahubTask.TargetResource.ClickHouseParam.DropCls.DropClsTopicId != nil {
					dropClsMap["drop_cls_topic_id"] = datahubTask.TargetResource.ClickHouseParam.DropCls.DropClsTopicId
				}

				if datahubTask.TargetResource.ClickHouseParam.DropCls.DropClsLogSet != nil {
					dropClsMap["drop_cls_log_set"] = datahubTask.TargetResource.ClickHouseParam.DropCls.DropClsLogSet
				}

				clickHouseParamMap["drop_cls"] = []interface{}{dropClsMap}
			}

			targetResourceMap["click_house_param"] = []interface{}{clickHouseParamMap}
		}

		if datahubTask.TargetResource.ClsParam != nil {
			clsParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.ClsParam.DecodeJson != nil {
				clsParamMap["decode_json"] = datahubTask.TargetResource.ClsParam.DecodeJson
			}

			if datahubTask.TargetResource.ClsParam.Resource != nil {
				clsParamMap["resource"] = datahubTask.TargetResource.ClsParam.Resource
			}

			if datahubTask.TargetResource.ClsParam.LogSet != nil {
				clsParamMap["log_set"] = datahubTask.TargetResource.ClsParam.LogSet
			}

			if datahubTask.TargetResource.ClsParam.ContentKey != nil {
				clsParamMap["content_key"] = datahubTask.TargetResource.ClsParam.ContentKey
			}

			if datahubTask.TargetResource.ClsParam.TimeField != nil {
				clsParamMap["time_field"] = datahubTask.TargetResource.ClsParam.TimeField
			}

			targetResourceMap["cls_param"] = []interface{}{clsParamMap}
		}

		if datahubTask.TargetResource.CosParam != nil {
			cosParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.CosParam.BucketName != nil {
				cosParamMap["bucket_name"] = datahubTask.TargetResource.CosParam.BucketName
			}

			if datahubTask.TargetResource.CosParam.Region != nil {
				cosParamMap["region"] = datahubTask.TargetResource.CosParam.Region
			}

			if datahubTask.TargetResource.CosParam.ObjectKey != nil {
				cosParamMap["object_key"] = datahubTask.TargetResource.CosParam.ObjectKey
			}

			if datahubTask.TargetResource.CosParam.AggregateBatchSize != nil {
				cosParamMap["aggregate_batch_size"] = datahubTask.TargetResource.CosParam.AggregateBatchSize
			}

			if datahubTask.TargetResource.CosParam.AggregateInterval != nil {
				cosParamMap["aggregate_interval"] = datahubTask.TargetResource.CosParam.AggregateInterval
			}

			if datahubTask.TargetResource.CosParam.FormatOutputType != nil {
				cosParamMap["format_output_type"] = datahubTask.TargetResource.CosParam.FormatOutputType
			}

			if datahubTask.TargetResource.CosParam.ObjectKeyPrefix != nil {
				cosParamMap["object_key_prefix"] = datahubTask.TargetResource.CosParam.ObjectKeyPrefix
			}

			if datahubTask.TargetResource.CosParam.DirectoryTimeFormat != nil {
				cosParamMap["directory_time_format"] = datahubTask.TargetResource.CosParam.DirectoryTimeFormat
			}

			targetResourceMap["cos_param"] = []interface{}{cosParamMap}
		}

		if datahubTask.TargetResource.MySQLParam != nil {
			mySQLParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.MySQLParam.Database != nil {
				mySQLParamMap["database"] = datahubTask.TargetResource.MySQLParam.Database
			}

			if datahubTask.TargetResource.MySQLParam.Table != nil {
				mySQLParamMap["table"] = datahubTask.TargetResource.MySQLParam.Table
			}

			if datahubTask.TargetResource.MySQLParam.Resource != nil {
				mySQLParamMap["resource"] = datahubTask.TargetResource.MySQLParam.Resource
			}

			if datahubTask.TargetResource.MySQLParam.SnapshotMode != nil {
				mySQLParamMap["snapshot_mode"] = datahubTask.TargetResource.MySQLParam.SnapshotMode
			}

			if datahubTask.TargetResource.MySQLParam.DdlTopic != nil {
				mySQLParamMap["ddl_topic"] = datahubTask.TargetResource.MySQLParam.DdlTopic
			}

			if datahubTask.TargetResource.MySQLParam.DataSourceMonitorMode != nil {
				mySQLParamMap["data_source_monitor_mode"] = datahubTask.TargetResource.MySQLParam.DataSourceMonitorMode
			}

			if datahubTask.TargetResource.MySQLParam.DataSourceMonitorResource != nil {
				mySQLParamMap["data_source_monitor_resource"] = datahubTask.TargetResource.MySQLParam.DataSourceMonitorResource
			}

			if datahubTask.TargetResource.MySQLParam.DataSourceIncrementMode != nil {
				mySQLParamMap["data_source_increment_mode"] = datahubTask.TargetResource.MySQLParam.DataSourceIncrementMode
			}

			if datahubTask.TargetResource.MySQLParam.DataSourceIncrementColumn != nil {
				mySQLParamMap["data_source_increment_column"] = datahubTask.TargetResource.MySQLParam.DataSourceIncrementColumn
			}

			if datahubTask.TargetResource.MySQLParam.DataSourceStartFrom != nil {
				mySQLParamMap["data_source_start_from"] = datahubTask.TargetResource.MySQLParam.DataSourceStartFrom
			}

			if datahubTask.TargetResource.MySQLParam.DataTargetInsertMode != nil {
				mySQLParamMap["data_target_insert_mode"] = datahubTask.TargetResource.MySQLParam.DataTargetInsertMode
			}

			if datahubTask.TargetResource.MySQLParam.DataTargetPrimaryKeyField != nil {
				mySQLParamMap["data_target_primary_key_field"] = datahubTask.TargetResource.MySQLParam.DataTargetPrimaryKeyField
			}

			if datahubTask.TargetResource.MySQLParam.DataTargetRecordMapping != nil {
				dataTargetRecordMappingList := []interface{}{}
				for _, dataTargetRecordMapping := range datahubTask.TargetResource.MySQLParam.DataTargetRecordMapping {
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

			if datahubTask.TargetResource.MySQLParam.TopicRegex != nil {
				mySQLParamMap["topic_regex"] = datahubTask.TargetResource.MySQLParam.TopicRegex
			}

			if datahubTask.TargetResource.MySQLParam.TopicReplacement != nil {
				mySQLParamMap["topic_replacement"] = datahubTask.TargetResource.MySQLParam.TopicReplacement
			}

			if datahubTask.TargetResource.MySQLParam.KeyColumns != nil {
				mySQLParamMap["key_columns"] = datahubTask.TargetResource.MySQLParam.KeyColumns
			}

			if datahubTask.TargetResource.MySQLParam.DropInvalidMessage != nil {
				mySQLParamMap["drop_invalid_message"] = datahubTask.TargetResource.MySQLParam.DropInvalidMessage
			}

			if datahubTask.TargetResource.MySQLParam.DropCls != nil {
				dropClsMap := map[string]interface{}{}

				if datahubTask.TargetResource.MySQLParam.DropCls.DropInvalidMessageToCls != nil {
					dropClsMap["drop_invalid_message_to_cls"] = datahubTask.TargetResource.MySQLParam.DropCls.DropInvalidMessageToCls
				}

				if datahubTask.TargetResource.MySQLParam.DropCls.DropClsRegion != nil {
					dropClsMap["drop_cls_region"] = datahubTask.TargetResource.MySQLParam.DropCls.DropClsRegion
				}

				if datahubTask.TargetResource.MySQLParam.DropCls.DropClsOwneruin != nil {
					dropClsMap["drop_cls_owneruin"] = datahubTask.TargetResource.MySQLParam.DropCls.DropClsOwneruin
				}

				if datahubTask.TargetResource.MySQLParam.DropCls.DropClsTopicId != nil {
					dropClsMap["drop_cls_topic_id"] = datahubTask.TargetResource.MySQLParam.DropCls.DropClsTopicId
				}

				if datahubTask.TargetResource.MySQLParam.DropCls.DropClsLogSet != nil {
					dropClsMap["drop_cls_log_set"] = datahubTask.TargetResource.MySQLParam.DropCls.DropClsLogSet
				}

				mySQLParamMap["drop_cls"] = []interface{}{dropClsMap}
			}

			if datahubTask.TargetResource.MySQLParam.OutputFormat != nil {
				mySQLParamMap["output_format"] = datahubTask.TargetResource.MySQLParam.OutputFormat
			}

			if datahubTask.TargetResource.MySQLParam.IsTablePrefix != nil {
				mySQLParamMap["is_table_prefix"] = datahubTask.TargetResource.MySQLParam.IsTablePrefix
			}

			if datahubTask.TargetResource.MySQLParam.IncludeContentChanges != nil {
				mySQLParamMap["include_content_changes"] = datahubTask.TargetResource.MySQLParam.IncludeContentChanges
			}

			if datahubTask.TargetResource.MySQLParam.IncludeQuery != nil {
				mySQLParamMap["include_query"] = datahubTask.TargetResource.MySQLParam.IncludeQuery
			}

			if datahubTask.TargetResource.MySQLParam.RecordWithSchema != nil {
				mySQLParamMap["record_with_schema"] = datahubTask.TargetResource.MySQLParam.RecordWithSchema
			}

			if datahubTask.TargetResource.MySQLParam.SignalDatabase != nil {
				mySQLParamMap["signal_database"] = datahubTask.TargetResource.MySQLParam.SignalDatabase
			}

			if datahubTask.TargetResource.MySQLParam.IsTableRegular != nil {
				mySQLParamMap["is_table_regular"] = datahubTask.TargetResource.MySQLParam.IsTableRegular
			}

			targetResourceMap["my_sql_param"] = []interface{}{mySQLParamMap}
		}

		if datahubTask.TargetResource.PostgreSQLParam != nil {
			postgreSQLParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.PostgreSQLParam.Database != nil {
				postgreSQLParamMap["database"] = datahubTask.TargetResource.PostgreSQLParam.Database
			}

			if datahubTask.TargetResource.PostgreSQLParam.Table != nil {
				postgreSQLParamMap["table"] = datahubTask.TargetResource.PostgreSQLParam.Table
			}

			if datahubTask.TargetResource.PostgreSQLParam.Resource != nil {
				postgreSQLParamMap["resource"] = datahubTask.TargetResource.PostgreSQLParam.Resource
			}

			if datahubTask.TargetResource.PostgreSQLParam.PluginName != nil {
				postgreSQLParamMap["plugin_name"] = datahubTask.TargetResource.PostgreSQLParam.PluginName
			}

			if datahubTask.TargetResource.PostgreSQLParam.SnapshotMode != nil {
				postgreSQLParamMap["snapshot_mode"] = datahubTask.TargetResource.PostgreSQLParam.SnapshotMode
			}

			if datahubTask.TargetResource.PostgreSQLParam.DataFormat != nil {
				postgreSQLParamMap["data_format"] = datahubTask.TargetResource.PostgreSQLParam.DataFormat
			}

			if datahubTask.TargetResource.PostgreSQLParam.DataTargetInsertMode != nil {
				postgreSQLParamMap["data_target_insert_mode"] = datahubTask.TargetResource.PostgreSQLParam.DataTargetInsertMode
			}

			if datahubTask.TargetResource.PostgreSQLParam.DataTargetPrimaryKeyField != nil {
				postgreSQLParamMap["data_target_primary_key_field"] = datahubTask.TargetResource.PostgreSQLParam.DataTargetPrimaryKeyField
			}

			if datahubTask.TargetResource.PostgreSQLParam.DataTargetRecordMapping != nil {
				dataTargetRecordMappingList := []interface{}{}
				for _, dataTargetRecordMapping := range datahubTask.TargetResource.PostgreSQLParam.DataTargetRecordMapping {
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

			if datahubTask.TargetResource.PostgreSQLParam.DropInvalidMessage != nil {
				postgreSQLParamMap["drop_invalid_message"] = datahubTask.TargetResource.PostgreSQLParam.DropInvalidMessage
			}

			if datahubTask.TargetResource.PostgreSQLParam.IsTableRegular != nil {
				postgreSQLParamMap["is_table_regular"] = datahubTask.TargetResource.PostgreSQLParam.IsTableRegular
			}

			if datahubTask.TargetResource.PostgreSQLParam.KeyColumns != nil {
				postgreSQLParamMap["key_columns"] = datahubTask.TargetResource.PostgreSQLParam.KeyColumns
			}

			if datahubTask.TargetResource.PostgreSQLParam.RecordWithSchema != nil {
				postgreSQLParamMap["record_with_schema"] = datahubTask.TargetResource.PostgreSQLParam.RecordWithSchema
			}

			targetResourceMap["postgre_sql_param"] = []interface{}{postgreSQLParamMap}
		}

		if datahubTask.TargetResource.TopicParam != nil {
			topicParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.TopicParam.Resource != nil {
				topicParamMap["resource"] = datahubTask.TargetResource.TopicParam.Resource
			}

			if datahubTask.TargetResource.TopicParam.OffsetType != nil {
				topicParamMap["offset_type"] = datahubTask.TargetResource.TopicParam.OffsetType
			}

			if datahubTask.TargetResource.TopicParam.StartTime != nil {
				topicParamMap["start_time"] = datahubTask.TargetResource.TopicParam.StartTime
			}

			if datahubTask.TargetResource.TopicParam.TopicId != nil {
				topicParamMap["topic_id"] = datahubTask.TargetResource.TopicParam.TopicId
			}

			if datahubTask.TargetResource.TopicParam.CompressionType != nil {
				topicParamMap["compression_type"] = datahubTask.TargetResource.TopicParam.CompressionType
			}

			if datahubTask.TargetResource.TopicParam.UseAutoCreateTopic != nil {
				topicParamMap["use_auto_create_topic"] = datahubTask.TargetResource.TopicParam.UseAutoCreateTopic
			}

			if datahubTask.TargetResource.TopicParam.MsgMultiple != nil {
				topicParamMap["msg_multiple"] = datahubTask.TargetResource.TopicParam.MsgMultiple
			}

			targetResourceMap["topic_param"] = []interface{}{topicParamMap}
		}

		if datahubTask.TargetResource.MariaDBParam != nil {
			mariaDBParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.MariaDBParam.Database != nil {
				mariaDBParamMap["database"] = datahubTask.TargetResource.MariaDBParam.Database
			}

			if datahubTask.TargetResource.MariaDBParam.Table != nil {
				mariaDBParamMap["table"] = datahubTask.TargetResource.MariaDBParam.Table
			}

			if datahubTask.TargetResource.MariaDBParam.Resource != nil {
				mariaDBParamMap["resource"] = datahubTask.TargetResource.MariaDBParam.Resource
			}

			if datahubTask.TargetResource.MariaDBParam.SnapshotMode != nil {
				mariaDBParamMap["snapshot_mode"] = datahubTask.TargetResource.MariaDBParam.SnapshotMode
			}

			if datahubTask.TargetResource.MariaDBParam.KeyColumns != nil {
				mariaDBParamMap["key_columns"] = datahubTask.TargetResource.MariaDBParam.KeyColumns
			}

			if datahubTask.TargetResource.MariaDBParam.IsTablePrefix != nil {
				mariaDBParamMap["is_table_prefix"] = datahubTask.TargetResource.MariaDBParam.IsTablePrefix
			}

			if datahubTask.TargetResource.MariaDBParam.OutputFormat != nil {
				mariaDBParamMap["output_format"] = datahubTask.TargetResource.MariaDBParam.OutputFormat
			}

			if datahubTask.TargetResource.MariaDBParam.IncludeContentChanges != nil {
				mariaDBParamMap["include_content_changes"] = datahubTask.TargetResource.MariaDBParam.IncludeContentChanges
			}

			if datahubTask.TargetResource.MariaDBParam.IncludeQuery != nil {
				mariaDBParamMap["include_query"] = datahubTask.TargetResource.MariaDBParam.IncludeQuery
			}

			if datahubTask.TargetResource.MariaDBParam.RecordWithSchema != nil {
				mariaDBParamMap["record_with_schema"] = datahubTask.TargetResource.MariaDBParam.RecordWithSchema
			}

			targetResourceMap["maria_db_param"] = []interface{}{mariaDBParamMap}
		}

		if datahubTask.TargetResource.SQLServerParam != nil {
			sQLServerParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.SQLServerParam.Database != nil {
				sQLServerParamMap["database"] = datahubTask.TargetResource.SQLServerParam.Database
			}

			if datahubTask.TargetResource.SQLServerParam.Table != nil {
				sQLServerParamMap["table"] = datahubTask.TargetResource.SQLServerParam.Table
			}

			if datahubTask.TargetResource.SQLServerParam.Resource != nil {
				sQLServerParamMap["resource"] = datahubTask.TargetResource.SQLServerParam.Resource
			}

			if datahubTask.TargetResource.SQLServerParam.SnapshotMode != nil {
				sQLServerParamMap["snapshot_mode"] = datahubTask.TargetResource.SQLServerParam.SnapshotMode
			}

			targetResourceMap["sql_server_param"] = []interface{}{sQLServerParamMap}
		}

		if datahubTask.TargetResource.CtsdbParam != nil {
			ctsdbParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.CtsdbParam.Resource != nil {
				ctsdbParamMap["resource"] = datahubTask.TargetResource.CtsdbParam.Resource
			}

			if datahubTask.TargetResource.CtsdbParam.CtsdbMetric != nil {
				ctsdbParamMap["ctsdb_metric"] = datahubTask.TargetResource.CtsdbParam.CtsdbMetric
			}

			targetResourceMap["ctsdb_param"] = []interface{}{ctsdbParamMap}
		}

		if datahubTask.TargetResource.ScfParam != nil {
			scfParamMap := map[string]interface{}{}

			if datahubTask.TargetResource.ScfParam.FunctionName != nil {
				scfParamMap["function_name"] = datahubTask.TargetResource.ScfParam.FunctionName
			}

			if datahubTask.TargetResource.ScfParam.Namespace != nil {
				scfParamMap["namespace"] = datahubTask.TargetResource.ScfParam.Namespace
			}

			if datahubTask.TargetResource.ScfParam.Qualifier != nil {
				scfParamMap["qualifier"] = datahubTask.TargetResource.ScfParam.Qualifier
			}

			if datahubTask.TargetResource.ScfParam.BatchSize != nil {
				scfParamMap["batch_size"] = datahubTask.TargetResource.ScfParam.BatchSize
			}

			if datahubTask.TargetResource.ScfParam.MaxRetries != nil {
				scfParamMap["max_retries"] = datahubTask.TargetResource.ScfParam.MaxRetries
			}

			targetResourceMap["scf_param"] = []interface{}{scfParamMap}
		}

		_ = d.Set("target_resource", []interface{}{targetResourceMap})
	}

	if datahubTask.TransformParam != nil {
		transformParamMap := map[string]interface{}{}

		if datahubTask.TransformParam.AnalysisFormat != nil {
			transformParamMap["analysis_format"] = datahubTask.TransformParam.AnalysisFormat
		}

		if datahubTask.TransformParam.OutputFormat != nil {
			transformParamMap["output_format"] = datahubTask.TransformParam.OutputFormat
		}

		if datahubTask.TransformParam.FailureParam != nil {
			failureParamMap := map[string]interface{}{}

			if datahubTask.TransformParam.FailureParam.Type != nil {
				failureParamMap["type"] = datahubTask.TransformParam.FailureParam.Type
			}

			if datahubTask.TransformParam.FailureParam.KafkaParam != nil {
				kafkaParamMap := map[string]interface{}{}

				if datahubTask.TransformParam.FailureParam.KafkaParam.SelfBuilt != nil {
					kafkaParamMap["self_built"] = datahubTask.TransformParam.FailureParam.KafkaParam.SelfBuilt
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.Resource != nil {
					kafkaParamMap["resource"] = datahubTask.TransformParam.FailureParam.KafkaParam.Resource
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.Topic != nil {
					kafkaParamMap["topic"] = datahubTask.TransformParam.FailureParam.KafkaParam.Topic
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.OffsetType != nil {
					kafkaParamMap["offset_type"] = datahubTask.TransformParam.FailureParam.KafkaParam.OffsetType
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.StartTime != nil {
					kafkaParamMap["start_time"] = datahubTask.TransformParam.FailureParam.KafkaParam.StartTime
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.ResourceName != nil {
					kafkaParamMap["resource_name"] = datahubTask.TransformParam.FailureParam.KafkaParam.ResourceName
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.ZoneId != nil {
					kafkaParamMap["zone_id"] = datahubTask.TransformParam.FailureParam.KafkaParam.ZoneId
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.TopicId != nil {
					kafkaParamMap["topic_id"] = datahubTask.TransformParam.FailureParam.KafkaParam.TopicId
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.PartitionNum != nil {
					kafkaParamMap["partition_num"] = datahubTask.TransformParam.FailureParam.KafkaParam.PartitionNum
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.EnableToleration != nil {
					kafkaParamMap["enable_toleration"] = datahubTask.TransformParam.FailureParam.KafkaParam.EnableToleration
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.QpsLimit != nil {
					kafkaParamMap["qps_limit"] = datahubTask.TransformParam.FailureParam.KafkaParam.QpsLimit
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.TableMappings != nil {
					tableMappingsList := []interface{}{}
					for _, tableMappings := range datahubTask.TransformParam.FailureParam.KafkaParam.TableMappings {
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

				if datahubTask.TransformParam.FailureParam.KafkaParam.UseTableMapping != nil {
					kafkaParamMap["use_table_mapping"] = datahubTask.TransformParam.FailureParam.KafkaParam.UseTableMapping
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.UseAutoCreateTopic != nil {
					kafkaParamMap["use_auto_create_topic"] = datahubTask.TransformParam.FailureParam.KafkaParam.UseAutoCreateTopic
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.CompressionType != nil {
					kafkaParamMap["compression_type"] = datahubTask.TransformParam.FailureParam.KafkaParam.CompressionType
				}

				if datahubTask.TransformParam.FailureParam.KafkaParam.MsgMultiple != nil {
					kafkaParamMap["msg_multiple"] = datahubTask.TransformParam.FailureParam.KafkaParam.MsgMultiple
				}

				failureParamMap["kafka_param"] = []interface{}{kafkaParamMap}
			}

			if datahubTask.TransformParam.FailureParam.RetryInterval != nil {
				failureParamMap["retry_interval"] = datahubTask.TransformParam.FailureParam.RetryInterval
			}

			if datahubTask.TransformParam.FailureParam.MaxRetryAttempts != nil {
				failureParamMap["max_retry_attempts"] = datahubTask.TransformParam.FailureParam.MaxRetryAttempts
			}

			if datahubTask.TransformParam.FailureParam.TopicParam != nil {
				topicParamMap := map[string]interface{}{}

				if datahubTask.TransformParam.FailureParam.TopicParam.Resource != nil {
					topicParamMap["resource"] = datahubTask.TransformParam.FailureParam.TopicParam.Resource
				}

				if datahubTask.TransformParam.FailureParam.TopicParam.OffsetType != nil {
					topicParamMap["offset_type"] = datahubTask.TransformParam.FailureParam.TopicParam.OffsetType
				}

				if datahubTask.TransformParam.FailureParam.TopicParam.StartTime != nil {
					topicParamMap["start_time"] = datahubTask.TransformParam.FailureParam.TopicParam.StartTime
				}

				if datahubTask.TransformParam.FailureParam.TopicParam.TopicId != nil {
					topicParamMap["topic_id"] = datahubTask.TransformParam.FailureParam.TopicParam.TopicId
				}

				if datahubTask.TransformParam.FailureParam.TopicParam.CompressionType != nil {
					topicParamMap["compression_type"] = datahubTask.TransformParam.FailureParam.TopicParam.CompressionType
				}

				if datahubTask.TransformParam.FailureParam.TopicParam.UseAutoCreateTopic != nil {
					topicParamMap["use_auto_create_topic"] = datahubTask.TransformParam.FailureParam.TopicParam.UseAutoCreateTopic
				}

				if datahubTask.TransformParam.FailureParam.TopicParam.MsgMultiple != nil {
					topicParamMap["msg_multiple"] = datahubTask.TransformParam.FailureParam.TopicParam.MsgMultiple
				}

				failureParamMap["topic_param"] = []interface{}{topicParamMap}
			}

			if datahubTask.TransformParam.FailureParam.DlqType != nil {
				failureParamMap["dlq_type"] = datahubTask.TransformParam.FailureParam.DlqType
			}

			transformParamMap["failure_param"] = []interface{}{failureParamMap}
		}

		if datahubTask.TransformParam.Content != nil {
			transformParamMap["content"] = datahubTask.TransformParam.Content
		}

		if datahubTask.TransformParam.SourceType != nil {
			transformParamMap["source_type"] = datahubTask.TransformParam.SourceType
		}

		if datahubTask.TransformParam.Regex != nil {
			transformParamMap["regex"] = datahubTask.TransformParam.Regex
		}

		if datahubTask.TransformParam.MapParam != nil {
			mapParamList := []interface{}{}
			for _, mapParam := range datahubTask.TransformParam.MapParam {
				mapParamMap := map[string]interface{}{}

				if mapParam.Key != nil {
					mapParamMap["key"] = mapParam.Key
				}

				if mapParam.Type != nil {
					mapParamMap["type"] = mapParam.Type
				}

				if mapParam.Value != nil {
					mapParamMap["value"] = mapParam.Value
				}

				mapParamList = append(mapParamList, mapParamMap)
			}

			transformParamMap["map_param"] = []interface{}{mapParamList}
		}

		if datahubTask.TransformParam.FilterParam != nil {
			filterParamList := []interface{}{}
			for _, filterParam := range datahubTask.TransformParam.FilterParam {
				filterParamMap := map[string]interface{}{}

				if filterParam.Key != nil {
					filterParamMap["key"] = filterParam.Key
				}

				if filterParam.MatchMode != nil {
					filterParamMap["match_mode"] = filterParam.MatchMode
				}

				if filterParam.Value != nil {
					filterParamMap["value"] = filterParam.Value
				}

				if filterParam.Type != nil {
					filterParamMap["type"] = filterParam.Type
				}

				filterParamList = append(filterParamList, filterParamMap)
			}

			transformParamMap["filter_param"] = []interface{}{filterParamList}
		}

		if datahubTask.TransformParam.Result != nil {
			transformParamMap["result"] = datahubTask.TransformParam.Result
		}

		if datahubTask.TransformParam.AnalyseResult != nil {
			analyseResultList := []interface{}{}
			for _, analyseResult := range datahubTask.TransformParam.AnalyseResult {
				analyseResultMap := map[string]interface{}{}

				if analyseResult.Key != nil {
					analyseResultMap["key"] = analyseResult.Key
				}

				if analyseResult.Type != nil {
					analyseResultMap["type"] = analyseResult.Type
				}

				if analyseResult.Value != nil {
					analyseResultMap["value"] = analyseResult.Value
				}

				analyseResultList = append(analyseResultList, analyseResultMap)
			}

			transformParamMap["analyse_result"] = []interface{}{analyseResultList}
		}

		if datahubTask.TransformParam.UseEventBus != nil {
			transformParamMap["use_event_bus"] = datahubTask.TransformParam.UseEventBus
		}

		_ = d.Set("transform_param", []interface{}{transformParamMap})
	}

	if datahubTask.SchemaId != nil {
		_ = d.Set("schema_id", datahubTask.SchemaId)
	}

	if datahubTask.TransformsParam != nil {
		transformsParamMap := map[string]interface{}{}

		if datahubTask.TransformsParam.Content != nil {
			transformsParamMap["content"] = datahubTask.TransformsParam.Content
		}

		if datahubTask.TransformsParam.FieldChain != nil {
			fieldChainList := []interface{}{}
			for _, fieldChain := range datahubTask.TransformsParam.FieldChain {
				fieldChainMap := map[string]interface{}{}

				if fieldChain.Analyse != nil {
					analyseMap := map[string]interface{}{}

					if fieldChain.Analyse.Format != nil {
						analyseMap["format"] = fieldChain.Analyse.Format
					}

					if fieldChain.Analyse.Regex != nil {
						analyseMap["regex"] = fieldChain.Analyse.Regex
					}

					if fieldChain.Analyse.InputValueType != nil {
						analyseMap["input_value_type"] = fieldChain.Analyse.InputValueType
					}

					if fieldChain.Analyse.InputValue != nil {
						analyseMap["input_value"] = fieldChain.Analyse.InputValue
					}

					fieldChainMap["analyse"] = []interface{}{analyseMap}
				}

				if fieldChain.SecondaryAnalyse != nil {
					secondaryAnalyseMap := map[string]interface{}{}

					if fieldChain.SecondaryAnalyse.Regex != nil {
						secondaryAnalyseMap["regex"] = fieldChain.SecondaryAnalyse.Regex
					}

					fieldChainMap["secondary_analyse"] = []interface{}{secondaryAnalyseMap}
				}

				if fieldChain.SMT != nil {
					sMTList := []interface{}{}
					for _, sMT := range fieldChain.SMT {
						sMTMap := map[string]interface{}{}

						if sMT.Key != nil {
							sMTMap["key"] = sMT.Key
						}

						if sMT.Operate != nil {
							sMTMap["operate"] = sMT.Operate
						}

						if sMT.SchemeType != nil {
							sMTMap["scheme_type"] = sMT.SchemeType
						}

						if sMT.Value != nil {
							sMTMap["value"] = sMT.Value
						}

						if sMT.ValueOperate != nil {
							valueOperateMap := map[string]interface{}{}

							if sMT.ValueOperate.Type != nil {
								valueOperateMap["type"] = sMT.ValueOperate.Type
							}

							if sMT.ValueOperate.Replace != nil {
								replaceMap := map[string]interface{}{}

								if sMT.ValueOperate.Replace.OldValue != nil {
									replaceMap["old_value"] = sMT.ValueOperate.Replace.OldValue
								}

								if sMT.ValueOperate.Replace.NewValue != nil {
									replaceMap["new_value"] = sMT.ValueOperate.Replace.NewValue
								}

								valueOperateMap["replace"] = []interface{}{replaceMap}
							}

							if sMT.ValueOperate.Substr != nil {
								substrMap := map[string]interface{}{}

								if sMT.ValueOperate.Substr.Start != nil {
									substrMap["start"] = sMT.ValueOperate.Substr.Start
								}

								if sMT.ValueOperate.Substr.End != nil {
									substrMap["end"] = sMT.ValueOperate.Substr.End
								}

								valueOperateMap["substr"] = []interface{}{substrMap}
							}

							if sMT.ValueOperate.Date != nil {
								dateMap := map[string]interface{}{}

								if sMT.ValueOperate.Date.Format != nil {
									dateMap["format"] = sMT.ValueOperate.Date.Format
								}

								if sMT.ValueOperate.Date.TargetType != nil {
									dateMap["target_type"] = sMT.ValueOperate.Date.TargetType
								}

								if sMT.ValueOperate.Date.TimeZone != nil {
									dateMap["time_zone"] = sMT.ValueOperate.Date.TimeZone
								}

								valueOperateMap["date"] = []interface{}{dateMap}
							}

							if sMT.ValueOperate.RegexReplace != nil {
								regexReplaceMap := map[string]interface{}{}

								if sMT.ValueOperate.RegexReplace.Regex != nil {
									regexReplaceMap["regex"] = sMT.ValueOperate.RegexReplace.Regex
								}

								if sMT.ValueOperate.RegexReplace.NewValue != nil {
									regexReplaceMap["new_value"] = sMT.ValueOperate.RegexReplace.NewValue
								}

								valueOperateMap["regex_replace"] = []interface{}{regexReplaceMap}
							}

							if sMT.ValueOperate.Split != nil {
								splitMap := map[string]interface{}{}

								if sMT.ValueOperate.Split.Regex != nil {
									splitMap["regex"] = sMT.ValueOperate.Split.Regex
								}

								valueOperateMap["split"] = []interface{}{splitMap}
							}

							if sMT.ValueOperate.KV != nil {
								kVMap := map[string]interface{}{}

								if sMT.ValueOperate.KV.Delimiter != nil {
									kVMap["delimiter"] = sMT.ValueOperate.KV.Delimiter
								}

								if sMT.ValueOperate.KV.Regex != nil {
									kVMap["regex"] = sMT.ValueOperate.KV.Regex
								}

								if sMT.ValueOperate.KV.KeepOriginalKey != nil {
									kVMap["keep_original_key"] = sMT.ValueOperate.KV.KeepOriginalKey
								}

								valueOperateMap["k_v"] = []interface{}{kVMap}
							}

							if sMT.ValueOperate.Result != nil {
								valueOperateMap["result"] = sMT.ValueOperate.Result
							}

							if sMT.ValueOperate.JsonPathReplace != nil {
								jsonPathReplaceMap := map[string]interface{}{}

								if sMT.ValueOperate.JsonPathReplace.OldValue != nil {
									jsonPathReplaceMap["old_value"] = sMT.ValueOperate.JsonPathReplace.OldValue
								}

								if sMT.ValueOperate.JsonPathReplace.NewValue != nil {
									jsonPathReplaceMap["new_value"] = sMT.ValueOperate.JsonPathReplace.NewValue
								}

								valueOperateMap["json_path_replace"] = []interface{}{jsonPathReplaceMap}
							}

							if sMT.ValueOperate.UrlDecode != nil {
								urlDecodeMap := map[string]interface{}{}

								if sMT.ValueOperate.UrlDecode.CharsetName != nil {
									urlDecodeMap["charset_name"] = sMT.ValueOperate.UrlDecode.CharsetName
								}

								valueOperateMap["url_decode"] = []interface{}{urlDecodeMap}
							}
							sMTMap["value_operate"] = []interface{}{valueOperateMap}
						}

						if sMT.OriginalValue != nil {
							sMTMap["original_value"] = sMT.OriginalValue
						}

						if sMT.ValueOperates != nil {
							valueOperatesList := []interface{}{}
							for _, valueOperates := range sMT.ValueOperates {
								valueOperatesMap := map[string]interface{}{}

								if valueOperates.Type != nil {
									valueOperatesMap["type"] = valueOperates.Type
								}

								if valueOperates.Replace != nil {
									replaceMap := map[string]interface{}{}

									if valueOperates.Replace.OldValue != nil {
										replaceMap["old_value"] = valueOperates.Replace.OldValue
									}

									if valueOperates.Replace.NewValue != nil {
										replaceMap["new_value"] = valueOperates.Replace.NewValue
									}

									valueOperatesMap["replace"] = []interface{}{replaceMap}
								}

								if valueOperates.Substr != nil {
									substrMap := map[string]interface{}{}

									if valueOperates.Substr.Start != nil {
										substrMap["start"] = valueOperates.Substr.Start
									}

									if valueOperates.Substr.End != nil {
										substrMap["end"] = valueOperates.Substr.End
									}

									valueOperatesMap["substr"] = []interface{}{substrMap}
								}

								if valueOperates.Date != nil {
									dateMap := map[string]interface{}{}

									if valueOperates.Date.Format != nil {
										dateMap["format"] = valueOperates.Date.Format
									}

									if valueOperates.Date.TargetType != nil {
										dateMap["target_type"] = valueOperates.Date.TargetType
									}

									if valueOperates.Date.TimeZone != nil {
										dateMap["time_zone"] = valueOperates.Date.TimeZone
									}

									valueOperatesMap["date"] = []interface{}{dateMap}
								}

								if valueOperates.RegexReplace != nil {
									regexReplaceMap := map[string]interface{}{}

									if valueOperates.RegexReplace.Regex != nil {
										regexReplaceMap["regex"] = valueOperates.RegexReplace.Regex
									}

									if valueOperates.RegexReplace.NewValue != nil {
										regexReplaceMap["new_value"] = valueOperates.RegexReplace.NewValue
									}

									valueOperatesMap["regex_replace"] = []interface{}{regexReplaceMap}
								}

								if valueOperates.Split != nil {
									splitMap := map[string]interface{}{}

									if valueOperates.Split.Regex != nil {
										splitMap["regex"] = valueOperates.Split.Regex
									}

									valueOperatesMap["split"] = []interface{}{splitMap}
								}

								if valueOperates.KV != nil {
									kVMap := map[string]interface{}{}

									if valueOperates.KV.Delimiter != nil {
										kVMap["delimiter"] = valueOperates.KV.Delimiter
									}

									if valueOperates.KV.Regex != nil {
										kVMap["regex"] = valueOperates.KV.Regex
									}

									if valueOperates.KV.KeepOriginalKey != nil {
										kVMap["keep_original_key"] = valueOperates.KV.KeepOriginalKey
									}

									valueOperatesMap["k_v"] = []interface{}{kVMap}
								}

								if valueOperates.Result != nil {
									valueOperatesMap["result"] = valueOperates.Result
								}

								if valueOperates.JsonPathReplace != nil {
									jsonPathReplaceMap := map[string]interface{}{}

									if valueOperates.JsonPathReplace.OldValue != nil {
										jsonPathReplaceMap["old_value"] = valueOperates.JsonPathReplace.OldValue
									}

									if valueOperates.JsonPathReplace.NewValue != nil {
										jsonPathReplaceMap["new_value"] = valueOperates.JsonPathReplace.NewValue
									}

									valueOperatesMap["json_path_replace"] = []interface{}{jsonPathReplaceMap}
								}

								if valueOperates.UrlDecode != nil {
									urlDecodeMap := map[string]interface{}{}

									if valueOperates.UrlDecode.CharsetName != nil {
										urlDecodeMap["charset_name"] = valueOperates.UrlDecode.CharsetName
									}

									valueOperatesMap["url_decode"] = []interface{}{urlDecodeMap}
								}
								valueOperatesList = append(valueOperatesList, valueOperatesMap)
							}

							sMTMap["value_operates"] = []interface{}{valueOperatesList}
						}

						sMTList = append(sMTList, sMTMap)
					}

					fieldChainMap["s_m_t"] = []interface{}{sMTList}
				}

				if fieldChain.Result != nil {
					fieldChainMap["result"] = fieldChain.Result
				}

				if fieldChain.AnalyseResult != nil {
					analyseResultList := []interface{}{}
					for _, analyseResult := range fieldChain.AnalyseResult {
						analyseResultMap := map[string]interface{}{}

						if analyseResult.Key != nil {
							analyseResultMap["key"] = analyseResult.Key
						}

						if analyseResult.Operate != nil {
							analyseResultMap["operate"] = analyseResult.Operate
						}

						if analyseResult.SchemeType != nil {
							analyseResultMap["scheme_type"] = analyseResult.SchemeType
						}

						if analyseResult.Value != nil {
							analyseResultMap["value"] = analyseResult.Value
						}

						if analyseResult.ValueOperate != nil {
							valueOperateMap := map[string]interface{}{}

							if analyseResult.ValueOperate.Type != nil {
								valueOperateMap["type"] = analyseResult.ValueOperate.Type
							}

							if analyseResult.ValueOperate.Replace != nil {
								replaceMap := map[string]interface{}{}

								if analyseResult.ValueOperate.Replace.OldValue != nil {
									replaceMap["old_value"] = analyseResult.ValueOperate.Replace.OldValue
								}

								if analyseResult.ValueOperate.Replace.NewValue != nil {
									replaceMap["new_value"] = analyseResult.ValueOperate.Replace.NewValue
								}

								valueOperateMap["replace"] = []interface{}{replaceMap}
							}

							if analyseResult.ValueOperate.Substr != nil {
								substrMap := map[string]interface{}{}

								if analyseResult.ValueOperate.Substr.Start != nil {
									substrMap["start"] = analyseResult.ValueOperate.Substr.Start
								}

								if analyseResult.ValueOperate.Substr.End != nil {
									substrMap["end"] = analyseResult.ValueOperate.Substr.End
								}

								valueOperateMap["substr"] = []interface{}{substrMap}
							}

							if analyseResult.ValueOperate.Date != nil {
								dateMap := map[string]interface{}{}

								if analyseResult.ValueOperate.Date.Format != nil {
									dateMap["format"] = analyseResult.ValueOperate.Date.Format
								}

								if analyseResult.ValueOperate.Date.TargetType != nil {
									dateMap["target_type"] = analyseResult.ValueOperate.Date.TargetType
								}

								if analyseResult.ValueOperate.Date.TimeZone != nil {
									dateMap["time_zone"] = analyseResult.ValueOperate.Date.TimeZone
								}

								valueOperateMap["date"] = []interface{}{dateMap}
							}

							if analyseResult.ValueOperate.RegexReplace != nil {
								regexReplaceMap := map[string]interface{}{}

								if analyseResult.ValueOperate.RegexReplace.Regex != nil {
									regexReplaceMap["regex"] = analyseResult.ValueOperate.RegexReplace.Regex
								}

								if analyseResult.ValueOperate.RegexReplace.NewValue != nil {
									regexReplaceMap["new_value"] = analyseResult.ValueOperate.RegexReplace.NewValue
								}

								valueOperateMap["regex_replace"] = []interface{}{regexReplaceMap}
							}

							if analyseResult.ValueOperate.Split != nil {
								splitMap := map[string]interface{}{}

								if analyseResult.ValueOperate.Split.Regex != nil {
									splitMap["regex"] = analyseResult.ValueOperate.Split.Regex
								}

								valueOperateMap["split"] = []interface{}{splitMap}
							}

							if analyseResult.ValueOperate.KV != nil {
								kVMap := map[string]interface{}{}

								if analyseResult.ValueOperate.KV.Delimiter != nil {
									kVMap["delimiter"] = analyseResult.ValueOperate.KV.Delimiter
								}

								if analyseResult.ValueOperate.KV.Regex != nil {
									kVMap["regex"] = analyseResult.ValueOperate.KV.Regex
								}

								if analyseResult.ValueOperate.KV.KeepOriginalKey != nil {
									kVMap["keep_original_key"] = analyseResult.ValueOperate.KV.KeepOriginalKey
								}

								valueOperateMap["k_v"] = []interface{}{kVMap}
							}

							if analyseResult.ValueOperate.Result != nil {
								valueOperateMap["result"] = analyseResult.ValueOperate.Result
							}

							if analyseResult.ValueOperate.JsonPathReplace != nil {
								jsonPathReplaceMap := map[string]interface{}{}

								if analyseResult.ValueOperate.JsonPathReplace.OldValue != nil {
									jsonPathReplaceMap["old_value"] = analyseResult.ValueOperate.JsonPathReplace.OldValue
								}

								if analyseResult.ValueOperate.JsonPathReplace.NewValue != nil {
									jsonPathReplaceMap["new_value"] = analyseResult.ValueOperate.JsonPathReplace.NewValue
								}

								valueOperateMap["json_path_replace"] = []interface{}{jsonPathReplaceMap}
							}

							if analyseResult.ValueOperate.UrlDecode != nil {
								urlDecodeMap := map[string]interface{}{}

								if analyseResult.ValueOperate.UrlDecode.CharsetName != nil {
									urlDecodeMap["charset_name"] = analyseResult.ValueOperate.UrlDecode.CharsetName
								}

								valueOperateMap["url_decode"] = []interface{}{urlDecodeMap}
							}

							analyseResultMap["value_operate"] = []interface{}{valueOperateMap}
						}

						if analyseResult.OriginalValue != nil {
							analyseResultMap["original_value"] = analyseResult.OriginalValue
						}

						if analyseResult.ValueOperates != nil {
							valueOperatesList := []interface{}{}
							for _, valueOperates := range analyseResult.ValueOperates {
								valueOperatesMap := map[string]interface{}{}

								if valueOperates.Type != nil {
									valueOperatesMap["type"] = valueOperates.Type
								}

								if valueOperates.Replace != nil {
									replaceMap := map[string]interface{}{}

									if valueOperates.Replace.OldValue != nil {
										replaceMap["old_value"] = valueOperates.Replace.OldValue
									}

									if valueOperates.Replace.NewValue != nil {
										replaceMap["new_value"] = valueOperates.Replace.NewValue
									}

									valueOperatesMap["replace"] = []interface{}{replaceMap}
								}

								if valueOperates.Substr != nil {
									substrMap := map[string]interface{}{}

									if valueOperates.Substr.Start != nil {
										substrMap["start"] = valueOperates.Substr.Start
									}

									if valueOperates.Substr.End != nil {
										substrMap["end"] = valueOperates.Substr.End
									}

									valueOperatesMap["substr"] = []interface{}{substrMap}
								}

								if valueOperates.Date != nil {
									dateMap := map[string]interface{}{}

									if valueOperates.Date.Format != nil {
										dateMap["format"] = valueOperates.Date.Format
									}

									if valueOperates.Date.TargetType != nil {
										dateMap["target_type"] = valueOperates.Date.TargetType
									}

									if valueOperates.Date.TimeZone != nil {
										dateMap["time_zone"] = valueOperates.Date.TimeZone
									}

									valueOperatesMap["date"] = []interface{}{dateMap}
								}

								if valueOperates.RegexReplace != nil {
									regexReplaceMap := map[string]interface{}{}

									if valueOperates.RegexReplace.Regex != nil {
										regexReplaceMap["regex"] = valueOperates.RegexReplace.Regex
									}

									if valueOperates.RegexReplace.NewValue != nil {
										regexReplaceMap["new_value"] = valueOperates.RegexReplace.NewValue
									}

									valueOperatesMap["regex_replace"] = []interface{}{regexReplaceMap}
								}

								if valueOperates.Split != nil {
									splitMap := map[string]interface{}{}

									if valueOperates.Split.Regex != nil {
										splitMap["regex"] = valueOperates.Split.Regex
									}

									valueOperatesMap["split"] = []interface{}{splitMap}
								}

								if valueOperates.KV != nil {
									kVMap := map[string]interface{}{}

									if valueOperates.KV.Delimiter != nil {
										kVMap["delimiter"] = valueOperates.KV.Delimiter
									}

									if valueOperates.KV.Regex != nil {
										kVMap["regex"] = valueOperates.KV.Regex
									}

									if valueOperates.KV.KeepOriginalKey != nil {
										kVMap["keep_original_key"] = valueOperates.KV.KeepOriginalKey
									}

									valueOperatesMap["k_v"] = []interface{}{kVMap}
								}

								if valueOperates.Result != nil {
									valueOperatesMap["result"] = valueOperates.Result
								}

								if valueOperates.JsonPathReplace != nil {
									jsonPathReplaceMap := map[string]interface{}{}

									if valueOperates.JsonPathReplace.OldValue != nil {
										jsonPathReplaceMap["old_value"] = valueOperates.JsonPathReplace.OldValue
									}

									if valueOperates.JsonPathReplace.NewValue != nil {
										jsonPathReplaceMap["new_value"] = valueOperates.JsonPathReplace.NewValue
									}

									valueOperatesMap["json_path_replace"] = []interface{}{jsonPathReplaceMap}
								}

								if valueOperates.UrlDecode != nil {
									urlDecodeMap := map[string]interface{}{}

									if valueOperates.UrlDecode.CharsetName != nil {
										urlDecodeMap["charset_name"] = valueOperates.UrlDecode.CharsetName
									}

									valueOperatesMap["url_decode"] = []interface{}{urlDecodeMap}
								}

								valueOperatesList = append(valueOperatesList, valueOperatesMap)
							}

							analyseResultMap["value_operates"] = []interface{}{valueOperatesList}
						}

						analyseResultList = append(analyseResultList, analyseResultMap)
					}

					fieldChainMap["analyse_result"] = []interface{}{analyseResultList}
				}

				if fieldChain.SecondaryAnalyseResult != nil {
					secondaryAnalyseResultList := []interface{}{}
					for _, secondaryAnalyseResult := range fieldChain.SecondaryAnalyseResult {
						secondaryAnalyseResultMap := map[string]interface{}{}

						if secondaryAnalyseResult.Key != nil {
							secondaryAnalyseResultMap["key"] = secondaryAnalyseResult.Key
						}

						if secondaryAnalyseResult.Operate != nil {
							secondaryAnalyseResultMap["operate"] = secondaryAnalyseResult.Operate
						}

						if secondaryAnalyseResult.SchemeType != nil {
							secondaryAnalyseResultMap["scheme_type"] = secondaryAnalyseResult.SchemeType
						}

						if secondaryAnalyseResult.Value != nil {
							secondaryAnalyseResultMap["value"] = secondaryAnalyseResult.Value
						}

						if secondaryAnalyseResult.ValueOperate != nil {
							valueOperateMap := map[string]interface{}{}

							if secondaryAnalyseResult.ValueOperate.Type != nil {
								valueOperateMap["type"] = secondaryAnalyseResult.ValueOperate.Type
							}

							if secondaryAnalyseResult.ValueOperate.Replace != nil {
								replaceMap := map[string]interface{}{}

								if secondaryAnalyseResult.ValueOperate.Replace.OldValue != nil {
									replaceMap["old_value"] = secondaryAnalyseResult.ValueOperate.Replace.OldValue
								}

								if secondaryAnalyseResult.ValueOperate.Replace.NewValue != nil {
									replaceMap["new_value"] = secondaryAnalyseResult.ValueOperate.Replace.NewValue
								}

								valueOperateMap["replace"] = []interface{}{replaceMap}
							}

							if secondaryAnalyseResult.ValueOperate.Substr != nil {
								substrMap := map[string]interface{}{}

								if secondaryAnalyseResult.ValueOperate.Substr.Start != nil {
									substrMap["start"] = secondaryAnalyseResult.ValueOperate.Substr.Start
								}

								if secondaryAnalyseResult.ValueOperate.Substr.End != nil {
									substrMap["end"] = secondaryAnalyseResult.ValueOperate.Substr.End
								}

								valueOperateMap["substr"] = []interface{}{substrMap}
							}

							if secondaryAnalyseResult.ValueOperate.Date != nil {
								dateMap := map[string]interface{}{}

								if secondaryAnalyseResult.ValueOperate.Date.Format != nil {
									dateMap["format"] = secondaryAnalyseResult.ValueOperate.Date.Format
								}

								if secondaryAnalyseResult.ValueOperate.Date.TargetType != nil {
									dateMap["target_type"] = secondaryAnalyseResult.ValueOperate.Date.TargetType
								}

								if secondaryAnalyseResult.ValueOperate.Date.TimeZone != nil {
									dateMap["time_zone"] = secondaryAnalyseResult.ValueOperate.Date.TimeZone
								}

								valueOperateMap["date"] = []interface{}{dateMap}
							}

							if secondaryAnalyseResult.ValueOperate.RegexReplace != nil {
								regexReplaceMap := map[string]interface{}{}

								if secondaryAnalyseResult.ValueOperate.RegexReplace.Regex != nil {
									regexReplaceMap["regex"] = secondaryAnalyseResult.ValueOperate.RegexReplace.Regex
								}

								if secondaryAnalyseResult.ValueOperate.RegexReplace.NewValue != nil {
									regexReplaceMap["new_value"] = secondaryAnalyseResult.ValueOperate.RegexReplace.NewValue
								}

								valueOperateMap["regex_replace"] = []interface{}{regexReplaceMap}
							}

							if secondaryAnalyseResult.ValueOperate.Split != nil {
								splitMap := map[string]interface{}{}

								if secondaryAnalyseResult.ValueOperate.Split.Regex != nil {
									splitMap["regex"] = secondaryAnalyseResult.ValueOperate.Split.Regex
								}

								valueOperateMap["split"] = []interface{}{splitMap}
							}

							if secondaryAnalyseResult.ValueOperate.KV != nil {
								kVMap := map[string]interface{}{}

								if secondaryAnalyseResult.ValueOperate.KV.Delimiter != nil {
									kVMap["delimiter"] = secondaryAnalyseResult.ValueOperate.KV.Delimiter
								}

								if secondaryAnalyseResult.ValueOperate.KV.Regex != nil {
									kVMap["regex"] = secondaryAnalyseResult.ValueOperate.KV.Regex
								}

								if secondaryAnalyseResult.ValueOperate.KV.KeepOriginalKey != nil {
									kVMap["keep_original_key"] = secondaryAnalyseResult.ValueOperate.KV.KeepOriginalKey
								}

								valueOperateMap["k_v"] = []interface{}{kVMap}
							}

							if secondaryAnalyseResult.ValueOperate.Result != nil {
								valueOperateMap["result"] = secondaryAnalyseResult.ValueOperate.Result
							}

							if secondaryAnalyseResult.ValueOperate.JsonPathReplace != nil {
								jsonPathReplaceMap := map[string]interface{}{}

								if secondaryAnalyseResult.ValueOperate.JsonPathReplace.OldValue != nil {
									jsonPathReplaceMap["old_value"] = secondaryAnalyseResult.ValueOperate.JsonPathReplace.OldValue
								}

								if secondaryAnalyseResult.ValueOperate.JsonPathReplace.NewValue != nil {
									jsonPathReplaceMap["new_value"] = secondaryAnalyseResult.ValueOperate.JsonPathReplace.NewValue
								}

								valueOperateMap["json_path_replace"] = []interface{}{jsonPathReplaceMap}
							}

							if secondaryAnalyseResult.ValueOperate.UrlDecode != nil {
								urlDecodeMap := map[string]interface{}{}

								if secondaryAnalyseResult.ValueOperate.UrlDecode.CharsetName != nil {
									urlDecodeMap["charset_name"] = secondaryAnalyseResult.ValueOperate.UrlDecode.CharsetName
								}

								valueOperateMap["url_decode"] = []interface{}{urlDecodeMap}
							}

							secondaryAnalyseResultMap["value_operate"] = []interface{}{valueOperateMap}
						}

						if secondaryAnalyseResult.OriginalValue != nil {
							secondaryAnalyseResultMap["original_value"] = secondaryAnalyseResult.OriginalValue
						}

						if secondaryAnalyseResult.ValueOperates != nil {
							valueOperatesList := []interface{}{}
							for _, valueOperates := range secondaryAnalyseResult.ValueOperates {
								valueOperatesMap := map[string]interface{}{}

								if valueOperates.Type != nil {
									valueOperatesMap["type"] = valueOperates.Type
								}

								if valueOperates.Replace != nil {
									replaceMap := map[string]interface{}{}

									if valueOperates.Replace.OldValue != nil {
										replaceMap["old_value"] = valueOperates.Replace.OldValue
									}

									if valueOperates.Replace.NewValue != nil {
										replaceMap["new_value"] = valueOperates.Replace.NewValue
									}

									valueOperatesMap["replace"] = []interface{}{replaceMap}
								}

								if valueOperates.Substr != nil {
									substrMap := map[string]interface{}{}

									if valueOperates.Substr.Start != nil {
										substrMap["start"] = valueOperates.Substr.Start
									}

									if valueOperates.Substr.End != nil {
										substrMap["end"] = valueOperates.Substr.End
									}

									valueOperatesMap["substr"] = []interface{}{substrMap}
								}

								if valueOperates.Date != nil {
									dateMap := map[string]interface{}{}

									if valueOperates.Date.Format != nil {
										dateMap["format"] = valueOperates.Date.Format
									}

									if valueOperates.Date.TargetType != nil {
										dateMap["target_type"] = valueOperates.Date.TargetType
									}

									if valueOperates.Date.TimeZone != nil {
										dateMap["time_zone"] = valueOperates.Date.TimeZone
									}

									valueOperatesMap["date"] = []interface{}{dateMap}
								}

								if valueOperates.RegexReplace != nil {
									regexReplaceMap := map[string]interface{}{}

									if valueOperates.RegexReplace.Regex != nil {
										regexReplaceMap["regex"] = valueOperates.RegexReplace.Regex
									}

									if valueOperates.RegexReplace.NewValue != nil {
										regexReplaceMap["new_value"] = valueOperates.RegexReplace.NewValue
									}

									valueOperatesMap["regex_replace"] = []interface{}{regexReplaceMap}
								}

								if valueOperates.Split != nil {
									splitMap := map[string]interface{}{}

									if valueOperates.Split.Regex != nil {
										splitMap["regex"] = valueOperates.Split.Regex
									}

									valueOperatesMap["split"] = []interface{}{splitMap}
								}

								if valueOperates.KV != nil {
									kVMap := map[string]interface{}{}

									if valueOperates.KV.Delimiter != nil {
										kVMap["delimiter"] = valueOperates.KV.Delimiter
									}

									if valueOperates.KV.Regex != nil {
										kVMap["regex"] = valueOperates.KV.Regex
									}

									if valueOperates.KV.KeepOriginalKey != nil {
										kVMap["keep_original_key"] = valueOperates.KV.KeepOriginalKey
									}

									valueOperatesMap["k_v"] = []interface{}{kVMap}
								}

								if valueOperates.Result != nil {
									valueOperatesMap["result"] = valueOperates.Result
								}

								if valueOperates.JsonPathReplace != nil {
									jsonPathReplaceMap := map[string]interface{}{}

									if valueOperates.JsonPathReplace.OldValue != nil {
										jsonPathReplaceMap["old_value"] = valueOperates.JsonPathReplace.OldValue
									}

									if valueOperates.JsonPathReplace.NewValue != nil {
										jsonPathReplaceMap["new_value"] = valueOperates.JsonPathReplace.NewValue
									}

									valueOperatesMap["json_path_replace"] = []interface{}{jsonPathReplaceMap}
								}

								if valueOperates.UrlDecode != nil {
									urlDecodeMap := map[string]interface{}{}

									if valueOperates.UrlDecode.CharsetName != nil {
										urlDecodeMap["charset_name"] = valueOperates.UrlDecode.CharsetName
									}

									valueOperatesMap["url_decode"] = []interface{}{urlDecodeMap}
								}

								valueOperatesList = append(valueOperatesList, valueOperatesMap)
							}

							secondaryAnalyseResultMap["value_operates"] = []interface{}{valueOperatesList}
						}

						secondaryAnalyseResultList = append(secondaryAnalyseResultList, secondaryAnalyseResultMap)
					}

					fieldChainMap["secondary_analyse_result"] = []interface{}{secondaryAnalyseResultList}
				}

				if fieldChain.AnalyseJsonResult != nil {
					fieldChainMap["analyse_json_result"] = fieldChain.AnalyseJsonResult
				}

				if fieldChain.SecondaryAnalyseJsonResult != nil {
					fieldChainMap["secondary_analyse_json_result"] = fieldChain.SecondaryAnalyseJsonResult
				}

				fieldChainList = append(fieldChainList, fieldChainMap)
			}

			transformsParamMap["field_chain"] = []interface{}{fieldChainList}
		}

		if datahubTask.TransformsParam.FilterParam != nil {
			filterParamList := []interface{}{}
			for _, filterParam := range datahubTask.TransformsParam.FilterParam {
				filterParamMap := map[string]interface{}{}

				if filterParam.Key != nil {
					filterParamMap["key"] = filterParam.Key
				}

				if filterParam.MatchMode != nil {
					filterParamMap["match_mode"] = filterParam.MatchMode
				}

				if filterParam.Value != nil {
					filterParamMap["value"] = filterParam.Value
				}

				if filterParam.Type != nil {
					filterParamMap["type"] = filterParam.Type
				}

				filterParamList = append(filterParamList, filterParamMap)
			}

			transformsParamMap["filter_param"] = []interface{}{filterParamList}
		}

		if datahubTask.TransformsParam.FailureParam != nil {
			failureParamMap := map[string]interface{}{}

			if datahubTask.TransformsParam.FailureParam.Type != nil {
				failureParamMap["type"] = datahubTask.TransformsParam.FailureParam.Type
			}

			if datahubTask.TransformsParam.FailureParam.KafkaParam != nil {
				kafkaParamMap := map[string]interface{}{}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.SelfBuilt != nil {
					kafkaParamMap["self_built"] = datahubTask.TransformsParam.FailureParam.KafkaParam.SelfBuilt
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.Resource != nil {
					kafkaParamMap["resource"] = datahubTask.TransformsParam.FailureParam.KafkaParam.Resource
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.Topic != nil {
					kafkaParamMap["topic"] = datahubTask.TransformsParam.FailureParam.KafkaParam.Topic
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.OffsetType != nil {
					kafkaParamMap["offset_type"] = datahubTask.TransformsParam.FailureParam.KafkaParam.OffsetType
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.StartTime != nil {
					kafkaParamMap["start_time"] = datahubTask.TransformsParam.FailureParam.KafkaParam.StartTime
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.ResourceName != nil {
					kafkaParamMap["resource_name"] = datahubTask.TransformsParam.FailureParam.KafkaParam.ResourceName
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.ZoneId != nil {
					kafkaParamMap["zone_id"] = datahubTask.TransformsParam.FailureParam.KafkaParam.ZoneId
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.TopicId != nil {
					kafkaParamMap["topic_id"] = datahubTask.TransformsParam.FailureParam.KafkaParam.TopicId
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.PartitionNum != nil {
					kafkaParamMap["partition_num"] = datahubTask.TransformsParam.FailureParam.KafkaParam.PartitionNum
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.EnableToleration != nil {
					kafkaParamMap["enable_toleration"] = datahubTask.TransformsParam.FailureParam.KafkaParam.EnableToleration
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.QpsLimit != nil {
					kafkaParamMap["qps_limit"] = datahubTask.TransformsParam.FailureParam.KafkaParam.QpsLimit
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.TableMappings != nil {
					tableMappingsList := []interface{}{}
					for _, tableMappings := range datahubTask.TransformsParam.FailureParam.KafkaParam.TableMappings {
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

				if datahubTask.TransformsParam.FailureParam.KafkaParam.UseTableMapping != nil {
					kafkaParamMap["use_table_mapping"] = datahubTask.TransformsParam.FailureParam.KafkaParam.UseTableMapping
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.UseAutoCreateTopic != nil {
					kafkaParamMap["use_auto_create_topic"] = datahubTask.TransformsParam.FailureParam.KafkaParam.UseAutoCreateTopic
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.CompressionType != nil {
					kafkaParamMap["compression_type"] = datahubTask.TransformsParam.FailureParam.KafkaParam.CompressionType
				}

				if datahubTask.TransformsParam.FailureParam.KafkaParam.MsgMultiple != nil {
					kafkaParamMap["msg_multiple"] = datahubTask.TransformsParam.FailureParam.KafkaParam.MsgMultiple
				}

				failureParamMap["kafka_param"] = []interface{}{kafkaParamMap}
			}

			if datahubTask.TransformsParam.FailureParam.RetryInterval != nil {
				failureParamMap["retry_interval"] = datahubTask.TransformsParam.FailureParam.RetryInterval
			}

			if datahubTask.TransformsParam.FailureParam.MaxRetryAttempts != nil {
				failureParamMap["max_retry_attempts"] = datahubTask.TransformsParam.FailureParam.MaxRetryAttempts
			}

			if datahubTask.TransformsParam.FailureParam.TopicParam != nil {
				topicParamMap := map[string]interface{}{}

				if datahubTask.TransformsParam.FailureParam.TopicParam.Resource != nil {
					topicParamMap["resource"] = datahubTask.TransformsParam.FailureParam.TopicParam.Resource
				}

				if datahubTask.TransformsParam.FailureParam.TopicParam.OffsetType != nil {
					topicParamMap["offset_type"] = datahubTask.TransformsParam.FailureParam.TopicParam.OffsetType
				}

				if datahubTask.TransformsParam.FailureParam.TopicParam.StartTime != nil {
					topicParamMap["start_time"] = datahubTask.TransformsParam.FailureParam.TopicParam.StartTime
				}

				if datahubTask.TransformsParam.FailureParam.TopicParam.TopicId != nil {
					topicParamMap["topic_id"] = datahubTask.TransformsParam.FailureParam.TopicParam.TopicId
				}

				if datahubTask.TransformsParam.FailureParam.TopicParam.CompressionType != nil {
					topicParamMap["compression_type"] = datahubTask.TransformsParam.FailureParam.TopicParam.CompressionType
				}

				if datahubTask.TransformsParam.FailureParam.TopicParam.UseAutoCreateTopic != nil {
					topicParamMap["use_auto_create_topic"] = datahubTask.TransformsParam.FailureParam.TopicParam.UseAutoCreateTopic
				}

				if datahubTask.TransformsParam.FailureParam.TopicParam.MsgMultiple != nil {
					topicParamMap["msg_multiple"] = datahubTask.TransformsParam.FailureParam.TopicParam.MsgMultiple
				}

				failureParamMap["topic_param"] = []interface{}{topicParamMap}
			}

			if datahubTask.TransformsParam.FailureParam.DlqType != nil {
				failureParamMap["dlq_type"] = datahubTask.TransformsParam.FailureParam.DlqType
			}

			transformsParamMap["failure_param"] = []interface{}{failureParamMap}
		}

		if datahubTask.TransformsParam.Result != nil {
			transformsParamMap["result"] = datahubTask.TransformsParam.Result
		}

		if datahubTask.TransformsParam.SourceType != nil {
			transformsParamMap["source_type"] = datahubTask.TransformsParam.SourceType
		}

		if datahubTask.TransformsParam.OutputFormat != nil {
			transformsParamMap["output_format"] = datahubTask.TransformsParam.OutputFormat
		}

		if datahubTask.TransformsParam.RowParam != nil {
			rowParamMap := map[string]interface{}{}

			if datahubTask.TransformsParam.RowParam.RowContent != nil {
				rowParamMap["row_content"] = datahubTask.TransformsParam.RowParam.RowContent
			}

			if datahubTask.TransformsParam.RowParam.KeyValueDelimiter != nil {
				rowParamMap["key_value_delimiter"] = datahubTask.TransformsParam.RowParam.KeyValueDelimiter
			}

			if datahubTask.TransformsParam.RowParam.EntryDelimiter != nil {
				rowParamMap["entry_delimiter"] = datahubTask.TransformsParam.RowParam.EntryDelimiter
			}

			transformsParamMap["row_param"] = []interface{}{rowParamMap}
		}

		if datahubTask.TransformsParam.KeepMetadata != nil {
			transformsParamMap["keep_metadata"] = datahubTask.TransformsParam.KeepMetadata
		}

		if datahubTask.TransformsParam.BatchAnalyse != nil {
			batchAnalyseMap := map[string]interface{}{}

			if datahubTask.TransformsParam.BatchAnalyse.Format != nil {
				batchAnalyseMap["format"] = datahubTask.TransformsParam.BatchAnalyse.Format
			}

			transformsParamMap["batch_analyse"] = []interface{}{batchAnalyseMap}
		}

		_ = d.Set("transforms_param", []interface{}{transformsParamMap})
	}

	return nil
}

func resourceTencentCloudCkafkaDatahubTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_datahub_task.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := ckafka.NewModifyDatahubTaskRequest()

	taskId := d.Id()
	request.TaskId = &taskId
	var hasChange bool
	if d.HasChange("task_name") {
		if v, ok := d.GetOk("task_name"); ok {
			request.TaskId = helper.String(v.(string))
		}
		hasChange = true
	}

	if !hasChange {
		return resourceTencentCloudCkafkaDatahubTaskRead(d, meta)
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().ModifyDatahubTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ckafka datahubTask failed, reason:%+v", logId, err)
		return err
	}

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 2*tccommon.ReadRetryTimeout, time.Second, service.CkafkaDatahubTaskStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCkafkaDatahubTaskRead(d, meta)
}

func resourceTencentCloudCkafkaDatahubTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_datahub_task.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	taskId := d.Id()

	if err := service.DeleteCkafkaDatahubTaskById(ctx, taskId); err != nil {
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"3"}, 2*tccommon.ReadRetryTimeout, time.Second, service.CkafkaDatahubTaskStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
