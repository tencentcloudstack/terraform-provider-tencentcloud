/*
Use this data source to query detailed information of ckafka connect_resource

Example Usage

```hcl
data "tencentcloud_ckafka_connect_resource" "connect_resource" {
  type = "DTS"
  search_word = "resourceName"
  offset = 0
  limit = 20
  resource_region = "region"
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

func dataSourceTencentCloudCkafkaConnectResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaConnectResourceRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Connection source type.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Keyword for search.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset, default is 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return the number, the default is 20, the maximum is 100.",
			},

			"resource_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Keyword query of the connection source, query the connection in the connection management list in the local region according to the region (only support the connection source containing the region input).",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Connection source list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of connection sources.",
						},
						"connect_resource_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Resource List.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ResourceId.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ResourceName.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource type.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resource status.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation time.",
									},
									"error_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Error Messages.",
									},
									"datahub_task_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of Datahub tasks associated with this connection source.",
									},
									"current_step": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The current step of the connection source.",
									},
									"task_progress": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Creation progress percentage.",
									},
									"step_list": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Step List.",
									},
									"dts_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Dts configuration, returned when Type is DTS.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Dts port.",
												},
												"group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The id of the Dts consumer group.",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The UserName of the Dts consumer group.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The password of the Dts consumer group.",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts Id.",
												},
												"topic": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Topic subscribed by Dts.",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to update to the associated Datahub task.",
												},
											},
										},
									},
									"mongo_d_b_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Mongo DB configuration, returned when Type is MONGODB.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "MongoDB port.",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB The username of the connection source.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB The password of the connection source.",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB Instance resource of connection source.",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "MongoDBWhether the connection source is a self-built cluster.",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDBInstance VIP of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDBThe vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to update to the associated Datahub task.",
												},
											},
										},
									},
									"es_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Es configuration, return when Type is ES.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Es port.",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EsThe username of the connection source.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es The password of the connection source.",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es Instance resource of connection source.",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "EsWhether the connection source is a self-built cluster.",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EsInstance VIP of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EsThe vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to update to the associated Datahub task.",
												},
											},
										},
									},
									"click_house_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Click House configuration, returned when Type is CLICKHOUSE.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ClickHouse port.",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouseThe username of the connection source.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse The password of the connection source.",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse Instance resource of connection source.",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "ClickHouseWhether the connection source is a self-built cluster.",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouseInstance VIP of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouseThe vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to update to the associated Datahub task.",
												},
											},
										},
									},
									"my_s_q_l_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Mysql configuration, returned when Type is MYSQL or TDSQL C MYSQL.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "MySQL port.",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MySQLThe username of the connection source.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MySQL The password of the connection source.",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MySQL Instance resource of connection source.",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MySQLInstance VIP of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MySQLThe vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to update to the associated Datahub task.",
												},
												"cluster_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Required when type is TDSQL C_MYSQL.",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Mysql Whether the connection source is a self-built cluster.",
												},
											},
										},
									},
									"postgre_s_q_l_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Postgresql configuration, returned when Type is POSTGRESQL or TDSQL C POSTGRESQL.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "PostgreSQL port.",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQLThe username of the connection source.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQL The password of the connection source.",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQL Instance resource of connection source.",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQL Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQL The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"cluster_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Required when type is TDSQL C_POSTGRESQL.",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to update to the associated Datahub task.",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "PostgreSQL Whether the connection source is a self-built cluster.",
												},
											},
										},
									},
									"maria_d_b_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Mariadb configuration, returned when Type is MARIADB.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "MariaDB port.",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB The username of the connection source.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB The password of the connection source.",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB Instance resource of connection source.",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to update to the associated Datahub task.",
												},
											},
										},
									},
									"s_q_l_server_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "SQL Server configuration, returned when Type is SQLSERVER.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "SQLServer port.",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServerThe username of the connection source.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServer The password of the connection source.",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServer Instance resource of connection source.",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServer Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServerThe vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to update to the associated Dip task.",
												},
											},
										},
									},
									"ctsdb_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Ctsdb configuration, returned when Type is CTSDB.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Ctsdb port.",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Ctsdb vip.",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Ctsdb vpcId.",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Ctsdb The username of the connection source.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Ctsdb The password of the connection source.",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Ctsdb Instance resource of connection source.",
												},
											},
										},
									},
									"doris_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Doris Configuration, returned when Type is DORIS.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Doris jdbc Load balancing connection port, usually mapped to port 9030 of fe.",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Doris The username of the connection source.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Doris  The password of the connection source.",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Doris  Instance resource of connection source.",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Doris Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Doris The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to update to the associated Datahub task.",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Doris Whether the connection source is a self-built cluster.",
												},
												"be_port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Doris&amp;#39;s http load balancing connection port, usually mapped to be&amp;#39;s 8040 port        .",
												},
											},
										},
									},
									"kafka_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Kafka configuration, returned when Type is KAFKA.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Instance resource of Kafka connection source, required when not self-built.",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether it is a self-built cluster.",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to update to the associated Dip task.",
												},
												"broker_address": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Broker address for Kafka connection, required for self-build                        .",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Instance resource region of CKafka connection source, required when crossing regions.",
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

func dataSourceTencentCloudCkafkaConnectResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_connect_resource.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("resource_region"); ok {
		paramMap["ResourceRegion"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.DescribeConnectResourcesResp

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaConnectResourceByFilter(ctx, paramMap)
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
		describeConnectResourcesRespMap := map[string]interface{}{}

		if result.TotalCount != nil {
			describeConnectResourcesRespMap["total_count"] = result.TotalCount
		}

		if result.ConnectResourceList != nil {
			connectResourceListList := []interface{}{}
			for _, connectResourceList := range result.ConnectResourceList {
				connectResourceListMap := map[string]interface{}{}

				if connectResourceList.ResourceId != nil {
					connectResourceListMap["resource_id"] = connectResourceList.ResourceId
				}

				if connectResourceList.ResourceName != nil {
					connectResourceListMap["resource_name"] = connectResourceList.ResourceName
				}

				if connectResourceList.Description != nil {
					connectResourceListMap["description"] = connectResourceList.Description
				}

				if connectResourceList.Type != nil {
					connectResourceListMap["type"] = connectResourceList.Type
				}

				if connectResourceList.Status != nil {
					connectResourceListMap["status"] = connectResourceList.Status
				}

				if connectResourceList.CreateTime != nil {
					connectResourceListMap["create_time"] = connectResourceList.CreateTime
				}

				if connectResourceList.ErrorMessage != nil {
					connectResourceListMap["error_message"] = connectResourceList.ErrorMessage
				}

				if connectResourceList.DatahubTaskCount != nil {
					connectResourceListMap["datahub_task_count"] = connectResourceList.DatahubTaskCount
				}

				if connectResourceList.CurrentStep != nil {
					connectResourceListMap["current_step"] = connectResourceList.CurrentStep
				}

				if connectResourceList.TaskProgress != nil {
					connectResourceListMap["task_progress"] = connectResourceList.TaskProgress
				}

				if connectResourceList.StepList != nil {
					connectResourceListMap["step_list"] = connectResourceList.StepList
				}

				if connectResourceList.DtsConnectParam != nil {
					dtsConnectParamMap := map[string]interface{}{}

					if connectResourceList.DtsConnectParam.Port != nil {
						dtsConnectParamMap["port"] = connectResourceList.DtsConnectParam.Port
					}

					if connectResourceList.DtsConnectParam.GroupId != nil {
						dtsConnectParamMap["group_id"] = connectResourceList.DtsConnectParam.GroupId
					}

					if connectResourceList.DtsConnectParam.UserName != nil {
						dtsConnectParamMap["user_name"] = connectResourceList.DtsConnectParam.UserName
					}

					if connectResourceList.DtsConnectParam.Password != nil {
						dtsConnectParamMap["password"] = connectResourceList.DtsConnectParam.Password
					}

					if connectResourceList.DtsConnectParam.Resource != nil {
						dtsConnectParamMap["resource"] = connectResourceList.DtsConnectParam.Resource
					}

					if connectResourceList.DtsConnectParam.Topic != nil {
						dtsConnectParamMap["topic"] = connectResourceList.DtsConnectParam.Topic
					}

					if connectResourceList.DtsConnectParam.IsUpdate != nil {
						dtsConnectParamMap["is_update"] = connectResourceList.DtsConnectParam.IsUpdate
					}

					connectResourceListMap["dts_connect_param"] = []interface{}{dtsConnectParamMap}
				}

				if connectResourceList.MongoDBConnectParam != nil {
					mongoDBConnectParamMap := map[string]interface{}{}

					if connectResourceList.MongoDBConnectParam.Port != nil {
						mongoDBConnectParamMap["port"] = connectResourceList.MongoDBConnectParam.Port
					}

					if connectResourceList.MongoDBConnectParam.UserName != nil {
						mongoDBConnectParamMap["user_name"] = connectResourceList.MongoDBConnectParam.UserName
					}

					if connectResourceList.MongoDBConnectParam.Password != nil {
						mongoDBConnectParamMap["password"] = connectResourceList.MongoDBConnectParam.Password
					}

					if connectResourceList.MongoDBConnectParam.Resource != nil {
						mongoDBConnectParamMap["resource"] = connectResourceList.MongoDBConnectParam.Resource
					}

					if connectResourceList.MongoDBConnectParam.SelfBuilt != nil {
						mongoDBConnectParamMap["self_built"] = connectResourceList.MongoDBConnectParam.SelfBuilt
					}

					if connectResourceList.MongoDBConnectParam.ServiceVip != nil {
						mongoDBConnectParamMap["service_vip"] = connectResourceList.MongoDBConnectParam.ServiceVip
					}

					if connectResourceList.MongoDBConnectParam.UniqVpcId != nil {
						mongoDBConnectParamMap["uniq_vpc_id"] = connectResourceList.MongoDBConnectParam.UniqVpcId
					}

					if connectResourceList.MongoDBConnectParam.IsUpdate != nil {
						mongoDBConnectParamMap["is_update"] = connectResourceList.MongoDBConnectParam.IsUpdate
					}

					connectResourceListMap["mongo_d_b_connect_param"] = []interface{}{mongoDBConnectParamMap}
				}

				if connectResourceList.EsConnectParam != nil {
					esConnectParamMap := map[string]interface{}{}

					if connectResourceList.EsConnectParam.Port != nil {
						esConnectParamMap["port"] = connectResourceList.EsConnectParam.Port
					}

					if connectResourceList.EsConnectParam.UserName != nil {
						esConnectParamMap["user_name"] = connectResourceList.EsConnectParam.UserName
					}

					if connectResourceList.EsConnectParam.Password != nil {
						esConnectParamMap["password"] = connectResourceList.EsConnectParam.Password
					}

					if connectResourceList.EsConnectParam.Resource != nil {
						esConnectParamMap["resource"] = connectResourceList.EsConnectParam.Resource
					}

					if connectResourceList.EsConnectParam.SelfBuilt != nil {
						esConnectParamMap["self_built"] = connectResourceList.EsConnectParam.SelfBuilt
					}

					if connectResourceList.EsConnectParam.ServiceVip != nil {
						esConnectParamMap["service_vip"] = connectResourceList.EsConnectParam.ServiceVip
					}

					if connectResourceList.EsConnectParam.UniqVpcId != nil {
						esConnectParamMap["uniq_vpc_id"] = connectResourceList.EsConnectParam.UniqVpcId
					}

					if connectResourceList.EsConnectParam.IsUpdate != nil {
						esConnectParamMap["is_update"] = connectResourceList.EsConnectParam.IsUpdate
					}

					connectResourceListMap["es_connect_param"] = []interface{}{esConnectParamMap}
				}

				if connectResourceList.ClickHouseConnectParam != nil {
					clickHouseConnectParamMap := map[string]interface{}{}

					if connectResourceList.ClickHouseConnectParam.Port != nil {
						clickHouseConnectParamMap["port"] = connectResourceList.ClickHouseConnectParam.Port
					}

					if connectResourceList.ClickHouseConnectParam.UserName != nil {
						clickHouseConnectParamMap["user_name"] = connectResourceList.ClickHouseConnectParam.UserName
					}

					if connectResourceList.ClickHouseConnectParam.Password != nil {
						clickHouseConnectParamMap["password"] = connectResourceList.ClickHouseConnectParam.Password
					}

					if connectResourceList.ClickHouseConnectParam.Resource != nil {
						clickHouseConnectParamMap["resource"] = connectResourceList.ClickHouseConnectParam.Resource
					}

					if connectResourceList.ClickHouseConnectParam.SelfBuilt != nil {
						clickHouseConnectParamMap["self_built"] = connectResourceList.ClickHouseConnectParam.SelfBuilt
					}

					if connectResourceList.ClickHouseConnectParam.ServiceVip != nil {
						clickHouseConnectParamMap["service_vip"] = connectResourceList.ClickHouseConnectParam.ServiceVip
					}

					if connectResourceList.ClickHouseConnectParam.UniqVpcId != nil {
						clickHouseConnectParamMap["uniq_vpc_id"] = connectResourceList.ClickHouseConnectParam.UniqVpcId
					}

					if connectResourceList.ClickHouseConnectParam.IsUpdate != nil {
						clickHouseConnectParamMap["is_update"] = connectResourceList.ClickHouseConnectParam.IsUpdate
					}

					connectResourceListMap["click_house_connect_param"] = []interface{}{clickHouseConnectParamMap}
				}

				if connectResourceList.MySQLConnectParam != nil {
					mySQLConnectParamMap := map[string]interface{}{}

					if connectResourceList.MySQLConnectParam.Port != nil {
						mySQLConnectParamMap["port"] = connectResourceList.MySQLConnectParam.Port
					}

					if connectResourceList.MySQLConnectParam.UserName != nil {
						mySQLConnectParamMap["user_name"] = connectResourceList.MySQLConnectParam.UserName
					}

					if connectResourceList.MySQLConnectParam.Password != nil {
						mySQLConnectParamMap["password"] = connectResourceList.MySQLConnectParam.Password
					}

					if connectResourceList.MySQLConnectParam.Resource != nil {
						mySQLConnectParamMap["resource"] = connectResourceList.MySQLConnectParam.Resource
					}

					if connectResourceList.MySQLConnectParam.ServiceVip != nil {
						mySQLConnectParamMap["service_vip"] = connectResourceList.MySQLConnectParam.ServiceVip
					}

					if connectResourceList.MySQLConnectParam.UniqVpcId != nil {
						mySQLConnectParamMap["uniq_vpc_id"] = connectResourceList.MySQLConnectParam.UniqVpcId
					}

					if connectResourceList.MySQLConnectParam.IsUpdate != nil {
						mySQLConnectParamMap["is_update"] = connectResourceList.MySQLConnectParam.IsUpdate
					}

					if connectResourceList.MySQLConnectParam.ClusterId != nil {
						mySQLConnectParamMap["cluster_id"] = connectResourceList.MySQLConnectParam.ClusterId
					}

					if connectResourceList.MySQLConnectParam.SelfBuilt != nil {
						mySQLConnectParamMap["self_built"] = connectResourceList.MySQLConnectParam.SelfBuilt
					}

					connectResourceListMap["my_s_q_l_connect_param"] = []interface{}{mySQLConnectParamMap}
				}

				if connectResourceList.PostgreSQLConnectParam != nil {
					postgreSQLConnectParamMap := map[string]interface{}{}

					if connectResourceList.PostgreSQLConnectParam.Port != nil {
						postgreSQLConnectParamMap["port"] = connectResourceList.PostgreSQLConnectParam.Port
					}

					if connectResourceList.PostgreSQLConnectParam.UserName != nil {
						postgreSQLConnectParamMap["user_name"] = connectResourceList.PostgreSQLConnectParam.UserName
					}

					if connectResourceList.PostgreSQLConnectParam.Password != nil {
						postgreSQLConnectParamMap["password"] = connectResourceList.PostgreSQLConnectParam.Password
					}

					if connectResourceList.PostgreSQLConnectParam.Resource != nil {
						postgreSQLConnectParamMap["resource"] = connectResourceList.PostgreSQLConnectParam.Resource
					}

					if connectResourceList.PostgreSQLConnectParam.ServiceVip != nil {
						postgreSQLConnectParamMap["service_vip"] = connectResourceList.PostgreSQLConnectParam.ServiceVip
					}

					if connectResourceList.PostgreSQLConnectParam.UniqVpcId != nil {
						postgreSQLConnectParamMap["uniq_vpc_id"] = connectResourceList.PostgreSQLConnectParam.UniqVpcId
					}

					if connectResourceList.PostgreSQLConnectParam.ClusterId != nil {
						postgreSQLConnectParamMap["cluster_id"] = connectResourceList.PostgreSQLConnectParam.ClusterId
					}

					if connectResourceList.PostgreSQLConnectParam.IsUpdate != nil {
						postgreSQLConnectParamMap["is_update"] = connectResourceList.PostgreSQLConnectParam.IsUpdate
					}

					if connectResourceList.PostgreSQLConnectParam.SelfBuilt != nil {
						postgreSQLConnectParamMap["self_built"] = connectResourceList.PostgreSQLConnectParam.SelfBuilt
					}

					connectResourceListMap["postgre_s_q_l_connect_param"] = []interface{}{postgreSQLConnectParamMap}
				}

				if connectResourceList.MariaDBConnectParam != nil {
					mariaDBConnectParamMap := map[string]interface{}{}

					if connectResourceList.MariaDBConnectParam.Port != nil {
						mariaDBConnectParamMap["port"] = connectResourceList.MariaDBConnectParam.Port
					}

					if connectResourceList.MariaDBConnectParam.UserName != nil {
						mariaDBConnectParamMap["user_name"] = connectResourceList.MariaDBConnectParam.UserName
					}

					if connectResourceList.MariaDBConnectParam.Password != nil {
						mariaDBConnectParamMap["password"] = connectResourceList.MariaDBConnectParam.Password
					}

					if connectResourceList.MariaDBConnectParam.Resource != nil {
						mariaDBConnectParamMap["resource"] = connectResourceList.MariaDBConnectParam.Resource
					}

					if connectResourceList.MariaDBConnectParam.ServiceVip != nil {
						mariaDBConnectParamMap["service_vip"] = connectResourceList.MariaDBConnectParam.ServiceVip
					}

					if connectResourceList.MariaDBConnectParam.UniqVpcId != nil {
						mariaDBConnectParamMap["uniq_vpc_id"] = connectResourceList.MariaDBConnectParam.UniqVpcId
					}

					if connectResourceList.MariaDBConnectParam.IsUpdate != nil {
						mariaDBConnectParamMap["is_update"] = connectResourceList.MariaDBConnectParam.IsUpdate
					}

					connectResourceListMap["maria_d_b_connect_param"] = []interface{}{mariaDBConnectParamMap}
				}

				if connectResourceList.SQLServerConnectParam != nil {
					sQLServerConnectParamMap := map[string]interface{}{}

					if connectResourceList.SQLServerConnectParam.Port != nil {
						sQLServerConnectParamMap["port"] = connectResourceList.SQLServerConnectParam.Port
					}

					if connectResourceList.SQLServerConnectParam.UserName != nil {
						sQLServerConnectParamMap["user_name"] = connectResourceList.SQLServerConnectParam.UserName
					}

					if connectResourceList.SQLServerConnectParam.Password != nil {
						sQLServerConnectParamMap["password"] = connectResourceList.SQLServerConnectParam.Password
					}

					if connectResourceList.SQLServerConnectParam.Resource != nil {
						sQLServerConnectParamMap["resource"] = connectResourceList.SQLServerConnectParam.Resource
					}

					if connectResourceList.SQLServerConnectParam.ServiceVip != nil {
						sQLServerConnectParamMap["service_vip"] = connectResourceList.SQLServerConnectParam.ServiceVip
					}

					if connectResourceList.SQLServerConnectParam.UniqVpcId != nil {
						sQLServerConnectParamMap["uniq_vpc_id"] = connectResourceList.SQLServerConnectParam.UniqVpcId
					}

					if connectResourceList.SQLServerConnectParam.IsUpdate != nil {
						sQLServerConnectParamMap["is_update"] = connectResourceList.SQLServerConnectParam.IsUpdate
					}

					connectResourceListMap["s_q_l_server_connect_param"] = []interface{}{sQLServerConnectParamMap}
				}

				if connectResourceList.CtsdbConnectParam != nil {
					ctsdbConnectParamMap := map[string]interface{}{}

					if connectResourceList.CtsdbConnectParam.Port != nil {
						ctsdbConnectParamMap["port"] = connectResourceList.CtsdbConnectParam.Port
					}

					if connectResourceList.CtsdbConnectParam.ServiceVip != nil {
						ctsdbConnectParamMap["service_vip"] = connectResourceList.CtsdbConnectParam.ServiceVip
					}

					if connectResourceList.CtsdbConnectParam.UniqVpcId != nil {
						ctsdbConnectParamMap["uniq_vpc_id"] = connectResourceList.CtsdbConnectParam.UniqVpcId
					}

					if connectResourceList.CtsdbConnectParam.UserName != nil {
						ctsdbConnectParamMap["user_name"] = connectResourceList.CtsdbConnectParam.UserName
					}

					if connectResourceList.CtsdbConnectParam.Password != nil {
						ctsdbConnectParamMap["password"] = connectResourceList.CtsdbConnectParam.Password
					}

					if connectResourceList.CtsdbConnectParam.Resource != nil {
						ctsdbConnectParamMap["resource"] = connectResourceList.CtsdbConnectParam.Resource
					}

					connectResourceListMap["ctsdb_connect_param"] = []interface{}{ctsdbConnectParamMap}
				}

				if connectResourceList.DorisConnectParam != nil {
					dorisConnectParamMap := map[string]interface{}{}

					if connectResourceList.DorisConnectParam.Port != nil {
						dorisConnectParamMap["port"] = connectResourceList.DorisConnectParam.Port
					}

					if connectResourceList.DorisConnectParam.UserName != nil {
						dorisConnectParamMap["user_name"] = connectResourceList.DorisConnectParam.UserName
					}

					if connectResourceList.DorisConnectParam.Password != nil {
						dorisConnectParamMap["password"] = connectResourceList.DorisConnectParam.Password
					}

					if connectResourceList.DorisConnectParam.Resource != nil {
						dorisConnectParamMap["resource"] = connectResourceList.DorisConnectParam.Resource
					}

					if connectResourceList.DorisConnectParam.ServiceVip != nil {
						dorisConnectParamMap["service_vip"] = connectResourceList.DorisConnectParam.ServiceVip
					}

					if connectResourceList.DorisConnectParam.UniqVpcId != nil {
						dorisConnectParamMap["uniq_vpc_id"] = connectResourceList.DorisConnectParam.UniqVpcId
					}

					if connectResourceList.DorisConnectParam.IsUpdate != nil {
						dorisConnectParamMap["is_update"] = connectResourceList.DorisConnectParam.IsUpdate
					}

					if connectResourceList.DorisConnectParam.SelfBuilt != nil {
						dorisConnectParamMap["self_built"] = connectResourceList.DorisConnectParam.SelfBuilt
					}

					if connectResourceList.DorisConnectParam.BePort != nil {
						dorisConnectParamMap["be_port"] = connectResourceList.DorisConnectParam.BePort
					}

					connectResourceListMap["doris_connect_param"] = []interface{}{dorisConnectParamMap}
				}

				if connectResourceList.KafkaConnectParam != nil {
					kafkaConnectParamMap := map[string]interface{}{}

					if connectResourceList.KafkaConnectParam.Resource != nil {
						kafkaConnectParamMap["resource"] = connectResourceList.KafkaConnectParam.Resource
					}

					if connectResourceList.KafkaConnectParam.SelfBuilt != nil {
						kafkaConnectParamMap["self_built"] = connectResourceList.KafkaConnectParam.SelfBuilt
					}

					if connectResourceList.KafkaConnectParam.IsUpdate != nil {
						kafkaConnectParamMap["is_update"] = connectResourceList.KafkaConnectParam.IsUpdate
					}

					if connectResourceList.KafkaConnectParam.BrokerAddress != nil {
						kafkaConnectParamMap["broker_address"] = connectResourceList.KafkaConnectParam.BrokerAddress
					}

					if connectResourceList.KafkaConnectParam.Region != nil {
						kafkaConnectParamMap["region"] = connectResourceList.KafkaConnectParam.Region
					}

					connectResourceListMap["kafka_connect_param"] = []interface{}{kafkaConnectParamMap}
				}

				connectResourceListList = append(connectResourceListList, connectResourceListMap)
			}

			describeConnectResourcesRespMap["connect_resource_list"] = []interface{}{connectResourceListList}
		}

		ids = append(ids, *result.ResourceId)
		_ = d.Set("result", describeConnectResourcesRespMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), describeConnectResourcesRespMap); e != nil {
			return e
		}
	}
	return nil
}
