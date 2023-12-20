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

func ResourceTencentCloudCkafkaConnectResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaConnectResourceCreate,
		Read:   resourceTencentCloudCkafkaConnectResourceRead,
		Update: resourceTencentCloudCkafkaConnectResourceUpdate,
		Delete: resourceTencentCloudCkafkaConnectResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "connection source name.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "connection source type.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Connection source description.",
			},

			"dts_connect_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Dts configuration, required when Type is DTS.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Dts port.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Id of the Dts consumption group.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The account number of the Dts consumption group.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The password of the Dts consumption group.",
						},
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Dts instance Id.",
						},
						"topic": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Topic subscribed by Dts.",
						},
						"is_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to update to the associated Datahub task, default: false.",
						},
					},
				},
			},

			"mongodb_connect_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Mongo DB configuration, required when Type is MONGODB.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "MongoDB port.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The username of the Mongo DB connection source.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Password for the source of the Mongo DB connection.",
						},
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance resource of Mongo DB connection source.",
						},
						"self_built": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the Mongo DB connection source is a self-built cluster.",
						},
						"service_vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The instance VIP of the Mongo DB connection source, when it is a Tencent Cloud instance, it is required.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The vpc Id of the Mongo DB connection source, which is required when it is a Tencent Cloud instance.",
						},
						"is_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to update to the associated Datahub task, default: false.",
						},
					},
				},
			},

			"es_connect_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Es configuration, required when Type is ES.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Es port.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Es The username of the connection source.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Es The password of the connection source.",
						},
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance resource of Es connection source.",
						},
						"self_built": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the Es connection source is a self-built cluster.",
						},
						"service_vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The instance vip of the Es connection source, when it is a Tencent Cloud instance, it is required.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The vpc Id of the Es connection source, when it is a Tencent Cloud instance, it is required.",
						},
						"is_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to update to the associated Datahub task, default: false.",
						},
					},
				},
			},

			"clickhouse_connect_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "ClickHouse configuration, required when Type is CLICKHOUSE.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Clickhouse connection port.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The username of the clickhouse connection source.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Password for Clickhouse connection source.",
						},
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance resources for Click House connection sources.",
						},
						"self_built": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the Clickhouse connection source is a self-built cluster.",
						},
						"service_vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance VIP of the ClickHouse connection source, when it is a Tencent Cloud instance, it is required.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The vpc Id of the source of the ClickHouse connection, when it is a Tencent Cloud instance, it is required.",
						},
						"is_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to update to the associated Datahub task, default: false.",
						},
					},
				},
			},

			"mysql_connect_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "MySQL configuration, required when Type is MYSQL or TDSQL C_MYSQL.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "MySQL port.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Username of Mysql connection source.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Mysql connection source password.",
						},
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance resource of My SQL connection source.",
						},
						"service_vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The instance vip of the MySQL connection source, when it is a Tencent Cloud instance, it is required.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The vpc Id of the My SQL connection source, when it is a Tencent Cloud instance, it is required.",
						},
						"is_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to update to the associated Datahub task, default: false.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Required when type is TDSQL C_MYSQL.",
						},
						"self_built": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Mysql Whether the connection source is a self-built cluster, default: false.",
						},
					},
				},
			},

			"postgresql_connect_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Postgresql configuration, required when Type is POSTGRESQL or TDSQL C POSTGRESQL.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "PostgreSQL port.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "PostgreSQL The username of the connection source.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "PostgreSQL password.",
						},
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "PostgreSQL instanceId.",
						},
						"service_vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The instance VIP of the Postgresql connection source, when it is a Tencent Cloud instance, it is required.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The instance vpcId of the Postgresql connection source, when it is a Tencent Cloud instance, it is required.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Required when type is TDSQL C_POSTGRESQL.",
						},
						"is_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to update to the associated Datahub task, default: false.",
						},
						"self_built": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "PostgreSQL Whether the connection source is a self-built cluster, default: false.",
						},
					},
				},
			},

			"mariadb_connect_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Maria DB configuration, required when Type is MARIADB.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "MariaDB port.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "MariaDB The username of the connection source.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "MariaDB password.",
						},
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "MariaDB instanceId.",
						},
						"service_vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The instance vip of the Maria DB connection source, when it is a Tencent Cloud instance, it is required.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "MariaDB vpcId, When it is a Tencent Cloud instance, it is required.",
						},
						"is_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to update to the associated Datahub task, default: false.",
						},
					},
				},
			},

			"sqlserver_connect_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "SQLServer configuration, required when Type is SQLSERVER.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "SQLServer port.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "SQLServer The username of the connection source.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "SQLServer password.",
						},
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "SQLServer instanceId.",
						},
						"service_vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQLServer instance vip, When it is a Tencent Cloud instance, it is required.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQLServer vpcId, When it is a Tencent Cloud instance, it is required.",
						},
						"is_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to update to the associated Dip task, default: false.",
						},
					},
				},
			},

			"doris_connect_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Doris configuration, required when Type is DORIS.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Doris jdbc CLB port, Usually mapped to port 9030 of fe.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Doris  The username of the connection source.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Doris  password.",
						},
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Doris  instanceId.",
						},
						"service_vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Doris vip, When it is a Tencent Cloud instance, it is required.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Doris vpcId, When it is a Tencent Cloud instance, it is required.",
						},
						"is_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to update to the associated Datahub task, default: false.",
						},
						"self_built": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Doris Whether the connection source is a self-built cluster, default: false.",
						},
						"be_port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Doris http CLB port, Usually mapped to port 8040 of be.",
						},
					},
				},
			},

			"kafka_connect_param": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Kafka configuration, required when Type is KAFKA.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Kafka instanceId, When it is a Tencent Cloud instance, it is required.",
						},
						"self_built": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether it is a self-built cluster, default: false.",
						},
						"is_update": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to update to the associated Dip task, default: false.",
						},
						"broker_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Kafka broker ip, Mandatory when self-built.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CKafka instanceId region, Required when crossing regions.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCkafkaConnectResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_connect_resource.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = ckafka.NewCreateConnectResourceRequest()
		response   = ckafka.NewCreateConnectResourceResponse()
		resourceId string
	)
	if v, ok := d.GetOk("resource_name"); ok {
		request.ResourceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "dts_connect_param"); ok {
		dtsConnectParam := ckafka.DtsConnectParam{}
		if v, ok := dMap["port"]; ok {
			dtsConnectParam.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["group_id"]; ok {
			dtsConnectParam.GroupId = helper.String(v.(string))
		}
		if v, ok := dMap["user_name"]; ok {
			dtsConnectParam.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			dtsConnectParam.Password = helper.String(v.(string))
		}
		if v, ok := dMap["resource"]; ok {
			dtsConnectParam.Resource = helper.String(v.(string))
		}
		if v, ok := dMap["topic"]; ok {
			dtsConnectParam.Topic = helper.String(v.(string))
		}
		if v, ok := dMap["is_update"]; ok {
			dtsConnectParam.IsUpdate = helper.Bool(v.(bool))
		}
		request.DtsConnectParam = &dtsConnectParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "mongodb_connect_param"); ok {
		mongoDBConnectParam := ckafka.MongoDBConnectParam{}
		if v, ok := dMap["port"]; ok {
			mongoDBConnectParam.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["user_name"]; ok {
			mongoDBConnectParam.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			mongoDBConnectParam.Password = helper.String(v.(string))
		}
		if v, ok := dMap["resource"]; ok {
			mongoDBConnectParam.Resource = helper.String(v.(string))
		}
		if v, ok := dMap["self_built"]; ok {
			mongoDBConnectParam.SelfBuilt = helper.Bool(v.(bool))
		}
		if v, ok := dMap["service_vip"]; ok {
			mongoDBConnectParam.ServiceVip = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpc_id"]; ok {
			mongoDBConnectParam.UniqVpcId = helper.String(v.(string))
		}
		if v, ok := dMap["is_update"]; ok {
			mongoDBConnectParam.IsUpdate = helper.Bool(v.(bool))
		}
		request.MongoDBConnectParam = &mongoDBConnectParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "es_connect_param"); ok {
		esConnectParam := ckafka.EsConnectParam{}
		if v, ok := dMap["port"]; ok {
			esConnectParam.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["user_name"]; ok {
			esConnectParam.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			esConnectParam.Password = helper.String(v.(string))
		}
		if v, ok := dMap["resource"]; ok {
			esConnectParam.Resource = helper.String(v.(string))
		}
		if v, ok := dMap["self_built"]; ok {
			esConnectParam.SelfBuilt = helper.Bool(v.(bool))
		}
		if v, ok := dMap["service_vip"]; ok {
			esConnectParam.ServiceVip = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpc_id"]; ok {
			esConnectParam.UniqVpcId = helper.String(v.(string))
		}
		if v, ok := dMap["is_update"]; ok {
			esConnectParam.IsUpdate = helper.Bool(v.(bool))
		}
		request.EsConnectParam = &esConnectParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "clickhouse_connect_param"); ok {
		clickHouseConnectParam := ckafka.ClickHouseConnectParam{}
		if v, ok := dMap["port"]; ok {
			clickHouseConnectParam.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["user_name"]; ok {
			clickHouseConnectParam.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			clickHouseConnectParam.Password = helper.String(v.(string))
		}
		if v, ok := dMap["resource"]; ok {
			clickHouseConnectParam.Resource = helper.String(v.(string))
		}
		if v, ok := dMap["self_built"]; ok {
			clickHouseConnectParam.SelfBuilt = helper.Bool(v.(bool))
		}
		if v, ok := dMap["service_vip"]; ok {
			clickHouseConnectParam.ServiceVip = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpc_id"]; ok {
			clickHouseConnectParam.UniqVpcId = helper.String(v.(string))
		}
		if v, ok := dMap["is_update"]; ok {
			clickHouseConnectParam.IsUpdate = helper.Bool(v.(bool))
		}
		request.ClickHouseConnectParam = &clickHouseConnectParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "mysql_connect_param"); ok {
		mySQLConnectParam := ckafka.MySQLConnectParam{}
		if v, ok := dMap["port"]; ok {
			mySQLConnectParam.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["user_name"]; ok {
			mySQLConnectParam.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			mySQLConnectParam.Password = helper.String(v.(string))
		}
		if v, ok := dMap["resource"]; ok {
			mySQLConnectParam.Resource = helper.String(v.(string))
		}
		if v, ok := dMap["service_vip"]; ok {
			mySQLConnectParam.ServiceVip = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpc_id"]; ok {
			mySQLConnectParam.UniqVpcId = helper.String(v.(string))
		}
		if v, ok := dMap["is_update"]; ok {
			mySQLConnectParam.IsUpdate = helper.Bool(v.(bool))
		}
		if v, ok := dMap["cluster_id"]; ok {
			mySQLConnectParam.ClusterId = helper.String(v.(string))
		}
		if v, ok := dMap["self_built"]; ok {
			mySQLConnectParam.SelfBuilt = helper.Bool(v.(bool))
		}
		request.MySQLConnectParam = &mySQLConnectParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "postgresql_connect_param"); ok {
		postgreSQLConnectParam := ckafka.PostgreSQLConnectParam{}
		if v, ok := dMap["port"]; ok {
			postgreSQLConnectParam.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["user_name"]; ok {
			postgreSQLConnectParam.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			postgreSQLConnectParam.Password = helper.String(v.(string))
		}
		if v, ok := dMap["resource"]; ok {
			postgreSQLConnectParam.Resource = helper.String(v.(string))
		}
		if v, ok := dMap["service_vip"]; ok {
			postgreSQLConnectParam.ServiceVip = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpc_id"]; ok {
			postgreSQLConnectParam.UniqVpcId = helper.String(v.(string))
		}
		if v, ok := dMap["cluster_id"]; ok {
			postgreSQLConnectParam.ClusterId = helper.String(v.(string))
		}
		if v, ok := dMap["is_update"]; ok {
			postgreSQLConnectParam.IsUpdate = helper.Bool(v.(bool))
		}
		if v, ok := dMap["self_built"]; ok {
			postgreSQLConnectParam.SelfBuilt = helper.Bool(v.(bool))
		}
		request.PostgreSQLConnectParam = &postgreSQLConnectParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "mariadb_connect_param"); ok {
		mariaDBConnectParam := ckafka.MariaDBConnectParam{}
		if v, ok := dMap["port"]; ok {
			mariaDBConnectParam.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["user_name"]; ok {
			mariaDBConnectParam.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			mariaDBConnectParam.Password = helper.String(v.(string))
		}
		if v, ok := dMap["resource"]; ok {
			mariaDBConnectParam.Resource = helper.String(v.(string))
		}
		if v, ok := dMap["service_vip"]; ok {
			mariaDBConnectParam.ServiceVip = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpc_id"]; ok {
			mariaDBConnectParam.UniqVpcId = helper.String(v.(string))
		}
		if v, ok := dMap["is_update"]; ok {
			mariaDBConnectParam.IsUpdate = helper.Bool(v.(bool))
		}
		request.MariaDBConnectParam = &mariaDBConnectParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "sqlserver_connect_param"); ok {
		sQLServerConnectParam := ckafka.SQLServerConnectParam{}
		if v, ok := dMap["port"]; ok {
			sQLServerConnectParam.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["user_name"]; ok {
			sQLServerConnectParam.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			sQLServerConnectParam.Password = helper.String(v.(string))
		}
		if v, ok := dMap["resource"]; ok {
			sQLServerConnectParam.Resource = helper.String(v.(string))
		}
		if v, ok := dMap["service_vip"]; ok {
			sQLServerConnectParam.ServiceVip = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpc_id"]; ok {
			sQLServerConnectParam.UniqVpcId = helper.String(v.(string))
		}
		if v, ok := dMap["is_update"]; ok {
			sQLServerConnectParam.IsUpdate = helper.Bool(v.(bool))
		}
		request.SQLServerConnectParam = &sQLServerConnectParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "doris_connect_param"); ok {
		dorisConnectParam := ckafka.DorisConnectParam{}
		if v, ok := dMap["port"]; ok {
			dorisConnectParam.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["user_name"]; ok {
			dorisConnectParam.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			dorisConnectParam.Password = helper.String(v.(string))
		}
		if v, ok := dMap["resource"]; ok {
			dorisConnectParam.Resource = helper.String(v.(string))
		}
		if v, ok := dMap["service_vip"]; ok {
			dorisConnectParam.ServiceVip = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpc_id"]; ok {
			dorisConnectParam.UniqVpcId = helper.String(v.(string))
		}
		if v, ok := dMap["is_update"]; ok {
			dorisConnectParam.IsUpdate = helper.Bool(v.(bool))
		}
		if v, ok := dMap["self_built"]; ok {
			dorisConnectParam.SelfBuilt = helper.Bool(v.(bool))
		}
		if v, ok := dMap["be_port"]; ok {
			dorisConnectParam.BePort = helper.IntInt64(v.(int))
		}
		request.DorisConnectParam = &dorisConnectParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "kafka_connect_param"); ok {
		kafkaConnectParam := ckafka.KafkaConnectParam{}
		if v, ok := dMap["resource"]; ok {
			kafkaConnectParam.Resource = helper.String(v.(string))
		}
		if v, ok := dMap["self_built"]; ok {
			kafkaConnectParam.SelfBuilt = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_update"]; ok {
			kafkaConnectParam.IsUpdate = helper.Bool(v.(bool))
		}
		if v, ok := dMap["broker_address"]; ok {
			kafkaConnectParam.BrokerAddress = helper.String(v.(string))
		}
		if v, ok := dMap["region"]; ok {
			kafkaConnectParam.Region = helper.String(v.(string))
		}
		request.KafkaConnectParam = &kafkaConnectParam
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().CreateConnectResource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ckafka connectResource failed, reason:%+v", logId, err)
		return err
	}

	resourceId = *response.Response.Result.ResourceId
	d.SetId(resourceId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 2*tccommon.ReadRetryTimeout, time.Second, service.CkafkaConnectResourceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCkafkaConnectResourceRead(d, meta)
}

func resourceTencentCloudCkafkaConnectResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_connect_resource.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	resourceId := d.Id()

	connectResource, err := service.DescribeCkafkaConnectResourceById(ctx, resourceId)
	if err != nil {
		return err
	}

	if connectResource == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CkafkaConnectResource` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if connectResource.ResourceName != nil {
		_ = d.Set("resource_name", connectResource.ResourceName)
	}

	if connectResource.Type != nil {
		_ = d.Set("type", connectResource.Type)
	}

	if connectResource.Description != nil {
		_ = d.Set("description", connectResource.Description)
	}

	if connectResource.DtsConnectParam != nil {
		dtsConnectParamMap := map[string]interface{}{}

		if connectResource.DtsConnectParam.Port != nil {
			dtsConnectParamMap["port"] = connectResource.DtsConnectParam.Port
		}

		if connectResource.DtsConnectParam.GroupId != nil {
			dtsConnectParamMap["group_id"] = connectResource.DtsConnectParam.GroupId
		}

		if connectResource.DtsConnectParam.UserName != nil {
			dtsConnectParamMap["user_name"] = connectResource.DtsConnectParam.UserName
		}

		password := connectResource.DtsConnectParam.Password
		if password != nil && *password != "" {
			dtsConnectParamMap["password"] = password
		} else {
			key := "dts_connect_param.0.password"
			if v, ok := d.GetOk(key); ok {
				dtsConnectParamMap["password"] = helper.String(v.(string))
			}
		}

		if connectResource.DtsConnectParam.Resource != nil {
			dtsConnectParamMap["resource"] = connectResource.DtsConnectParam.Resource
		}

		if connectResource.DtsConnectParam.Topic != nil {
			dtsConnectParamMap["topic"] = connectResource.DtsConnectParam.Topic
		}

		if connectResource.DtsConnectParam.IsUpdate != nil {
			dtsConnectParamMap["is_update"] = connectResource.DtsConnectParam.IsUpdate
		}

		_ = d.Set("dts_connect_param", []interface{}{dtsConnectParamMap})
	}

	if connectResource.MongoDBConnectParam != nil {
		mongoDBConnectParamMap := map[string]interface{}{}

		if connectResource.MongoDBConnectParam.Port != nil {
			mongoDBConnectParamMap["port"] = connectResource.MongoDBConnectParam.Port
		}

		if connectResource.MongoDBConnectParam.UserName != nil {
			mongoDBConnectParamMap["user_name"] = connectResource.MongoDBConnectParam.UserName
		}

		password := connectResource.MongoDBConnectParam.Password
		if password != nil && *password != "" {
			mongoDBConnectParamMap["password"] = password
		} else {
			key := "mongodb_connect_param.0.password"
			if v, ok := d.GetOk(key); ok {
				mongoDBConnectParamMap["password"] = helper.String(v.(string))
			}
		}

		if connectResource.MongoDBConnectParam.Resource != nil {
			mongoDBConnectParamMap["resource"] = connectResource.MongoDBConnectParam.Resource
		}

		if connectResource.MongoDBConnectParam.SelfBuilt != nil {
			mongoDBConnectParamMap["self_built"] = connectResource.MongoDBConnectParam.SelfBuilt
		}

		if connectResource.MongoDBConnectParam.ServiceVip != nil {
			mongoDBConnectParamMap["service_vip"] = connectResource.MongoDBConnectParam.ServiceVip
		}

		if connectResource.MongoDBConnectParam.UniqVpcId != nil {
			mongoDBConnectParamMap["uniq_vpc_id"] = connectResource.MongoDBConnectParam.UniqVpcId
		}

		if connectResource.MongoDBConnectParam.IsUpdate != nil {
			mongoDBConnectParamMap["is_update"] = connectResource.MongoDBConnectParam.IsUpdate
		}

		_ = d.Set("mongodb_connect_param", []interface{}{mongoDBConnectParamMap})
	}

	if connectResource.EsConnectParam != nil {
		esConnectParamMap := map[string]interface{}{}

		if connectResource.EsConnectParam.Port != nil {
			esConnectParamMap["port"] = connectResource.EsConnectParam.Port
		}

		if connectResource.EsConnectParam.UserName != nil {
			esConnectParamMap["user_name"] = connectResource.EsConnectParam.UserName
		}

		password := connectResource.EsConnectParam.Password
		if password != nil && *password != "" {
			esConnectParamMap["password"] = password
		} else {
			key := "es_connect_param.0.password"
			if v, ok := d.GetOk(key); ok {
				esConnectParamMap["password"] = helper.String(v.(string))
			}
		}

		if connectResource.EsConnectParam.Resource != nil {
			esConnectParamMap["resource"] = connectResource.EsConnectParam.Resource
		}

		if connectResource.EsConnectParam.SelfBuilt != nil {
			esConnectParamMap["self_built"] = connectResource.EsConnectParam.SelfBuilt
		}

		if connectResource.EsConnectParam.ServiceVip != nil {
			esConnectParamMap["service_vip"] = connectResource.EsConnectParam.ServiceVip
		}

		if connectResource.EsConnectParam.UniqVpcId != nil {
			esConnectParamMap["uniq_vpc_id"] = connectResource.EsConnectParam.UniqVpcId
		}

		if connectResource.EsConnectParam.IsUpdate != nil {
			esConnectParamMap["is_update"] = connectResource.EsConnectParam.IsUpdate
		}

		_ = d.Set("es_connect_param", []interface{}{esConnectParamMap})
	}

	if connectResource.ClickHouseConnectParam != nil {
		clickHouseConnectParamMap := map[string]interface{}{}

		if connectResource.ClickHouseConnectParam.Port != nil {
			clickHouseConnectParamMap["port"] = connectResource.ClickHouseConnectParam.Port
		}

		if connectResource.ClickHouseConnectParam.UserName != nil {
			clickHouseConnectParamMap["user_name"] = connectResource.ClickHouseConnectParam.UserName
		}

		password := connectResource.ClickHouseConnectParam.Password
		if password != nil && *password != "" {
			clickHouseConnectParamMap["password"] = password
		} else {
			key := "clickhouse_connect_param.0.password"
			if v, ok := d.GetOk(key); ok {
				clickHouseConnectParamMap["password"] = helper.String(v.(string))
			}
		}

		if connectResource.ClickHouseConnectParam.Resource != nil {
			clickHouseConnectParamMap["resource"] = connectResource.ClickHouseConnectParam.Resource
		}

		if connectResource.ClickHouseConnectParam.SelfBuilt != nil {
			clickHouseConnectParamMap["self_built"] = connectResource.ClickHouseConnectParam.SelfBuilt
		}

		if connectResource.ClickHouseConnectParam.ServiceVip != nil {
			clickHouseConnectParamMap["service_vip"] = connectResource.ClickHouseConnectParam.ServiceVip
		}

		if connectResource.ClickHouseConnectParam.UniqVpcId != nil {
			clickHouseConnectParamMap["uniq_vpc_id"] = connectResource.ClickHouseConnectParam.UniqVpcId
		}

		if connectResource.ClickHouseConnectParam.IsUpdate != nil {
			clickHouseConnectParamMap["is_update"] = connectResource.ClickHouseConnectParam.IsUpdate
		}

		_ = d.Set("clickhouse_connect_param", []interface{}{clickHouseConnectParamMap})
	}

	if connectResource.MySQLConnectParam != nil {
		mySQLConnectParamMap := map[string]interface{}{}

		if connectResource.MySQLConnectParam.Port != nil {
			mySQLConnectParamMap["port"] = connectResource.MySQLConnectParam.Port
		}

		if connectResource.MySQLConnectParam.UserName != nil {
			mySQLConnectParamMap["user_name"] = connectResource.MySQLConnectParam.UserName
		}

		password := connectResource.MySQLConnectParam.Password
		if password != nil && *password != "" {
			mySQLConnectParamMap["password"] = password
		} else {
			key := "mysql_connect_param.0.password"
			if v, ok := d.GetOk(key); ok {
				mySQLConnectParamMap["password"] = helper.String(v.(string))
			}
		}

		if connectResource.MySQLConnectParam.Resource != nil {
			mySQLConnectParamMap["resource"] = connectResource.MySQLConnectParam.Resource
		}

		if connectResource.MySQLConnectParam.ServiceVip != nil {
			mySQLConnectParamMap["service_vip"] = connectResource.MySQLConnectParam.ServiceVip
		}

		if connectResource.MySQLConnectParam.UniqVpcId != nil {
			mySQLConnectParamMap["uniq_vpc_id"] = connectResource.MySQLConnectParam.UniqVpcId
		}

		if connectResource.MySQLConnectParam.IsUpdate != nil {
			mySQLConnectParamMap["is_update"] = connectResource.MySQLConnectParam.IsUpdate
		}

		if connectResource.MySQLConnectParam.ClusterId != nil {
			mySQLConnectParamMap["cluster_id"] = connectResource.MySQLConnectParam.ClusterId
		}

		if connectResource.MySQLConnectParam.SelfBuilt != nil {
			mySQLConnectParamMap["self_built"] = connectResource.MySQLConnectParam.SelfBuilt
		}

		_ = d.Set("mysql_connect_param", []interface{}{mySQLConnectParamMap})
	}

	if connectResource.PostgreSQLConnectParam != nil {
		postgreSQLConnectParamMap := map[string]interface{}{}

		if connectResource.PostgreSQLConnectParam.Port != nil {
			postgreSQLConnectParamMap["port"] = connectResource.PostgreSQLConnectParam.Port
		}

		if connectResource.PostgreSQLConnectParam.UserName != nil {
			postgreSQLConnectParamMap["user_name"] = connectResource.PostgreSQLConnectParam.UserName
		}

		password := connectResource.PostgreSQLConnectParam.Password
		if password != nil && *password != "" {
			postgreSQLConnectParamMap["password"] = password
		} else {
			key := "postgresql_connect_param.0.password"
			if v, ok := d.GetOk(key); ok {
				postgreSQLConnectParamMap["password"] = helper.String(v.(string))
			}
		}

		if connectResource.PostgreSQLConnectParam.Resource != nil {
			postgreSQLConnectParamMap["resource"] = connectResource.PostgreSQLConnectParam.Resource
		}

		if connectResource.PostgreSQLConnectParam.ServiceVip != nil {
			postgreSQLConnectParamMap["service_vip"] = connectResource.PostgreSQLConnectParam.ServiceVip
		}

		if connectResource.PostgreSQLConnectParam.UniqVpcId != nil {
			postgreSQLConnectParamMap["uniq_vpc_id"] = connectResource.PostgreSQLConnectParam.UniqVpcId
		}

		if connectResource.PostgreSQLConnectParam.ClusterId != nil {
			postgreSQLConnectParamMap["cluster_id"] = connectResource.PostgreSQLConnectParam.ClusterId
		}

		if connectResource.PostgreSQLConnectParam.IsUpdate != nil {
			postgreSQLConnectParamMap["is_update"] = connectResource.PostgreSQLConnectParam.IsUpdate
		}

		if connectResource.PostgreSQLConnectParam.SelfBuilt != nil {
			postgreSQLConnectParamMap["self_built"] = connectResource.PostgreSQLConnectParam.SelfBuilt
		}

		_ = d.Set("postgresql_connect_param", []interface{}{postgreSQLConnectParamMap})
	}

	if connectResource.MariaDBConnectParam != nil {
		mariaDBConnectParamMap := map[string]interface{}{}

		if connectResource.MariaDBConnectParam.Port != nil {
			mariaDBConnectParamMap["port"] = connectResource.MariaDBConnectParam.Port
		}

		if connectResource.MariaDBConnectParam.UserName != nil {
			mariaDBConnectParamMap["user_name"] = connectResource.MariaDBConnectParam.UserName
		}

		password := connectResource.MariaDBConnectParam.Password
		if password != nil && *password != "" {
			mariaDBConnectParamMap["password"] = password
		} else {
			key := "mariadb_connect_param.0.password"
			if v, ok := d.GetOk(key); ok {
				mariaDBConnectParamMap["password"] = helper.String(v.(string))
			}
		}

		if connectResource.MariaDBConnectParam.Resource != nil {
			mariaDBConnectParamMap["resource"] = connectResource.MariaDBConnectParam.Resource
		}

		if connectResource.MariaDBConnectParam.ServiceVip != nil {
			mariaDBConnectParamMap["service_vip"] = connectResource.MariaDBConnectParam.ServiceVip
		}

		if connectResource.MariaDBConnectParam.UniqVpcId != nil {
			mariaDBConnectParamMap["uniq_vpc_id"] = connectResource.MariaDBConnectParam.UniqVpcId
		}

		if connectResource.MariaDBConnectParam.IsUpdate != nil {
			mariaDBConnectParamMap["is_update"] = connectResource.MariaDBConnectParam.IsUpdate
		}

		_ = d.Set("mariadb_connect_param", []interface{}{mariaDBConnectParamMap})
	}

	if connectResource.SQLServerConnectParam != nil {
		sQLServerConnectParamMap := map[string]interface{}{}

		if connectResource.SQLServerConnectParam.Port != nil {
			sQLServerConnectParamMap["port"] = connectResource.SQLServerConnectParam.Port
		}

		if connectResource.SQLServerConnectParam.UserName != nil {
			sQLServerConnectParamMap["user_name"] = connectResource.SQLServerConnectParam.UserName
		}

		password := connectResource.SQLServerConnectParam.Password
		if password != nil && *password != "" {
			sQLServerConnectParamMap["password"] = password
		} else {
			key := "sqlserver_connect_param.0.password"
			if v, ok := d.GetOk(key); ok {
				sQLServerConnectParamMap["password"] = helper.String(v.(string))
			}
		}

		if connectResource.SQLServerConnectParam.Resource != nil {
			sQLServerConnectParamMap["resource"] = connectResource.SQLServerConnectParam.Resource
		}

		if connectResource.SQLServerConnectParam.ServiceVip != nil {
			sQLServerConnectParamMap["service_vip"] = connectResource.SQLServerConnectParam.ServiceVip
		}

		if connectResource.SQLServerConnectParam.UniqVpcId != nil {
			sQLServerConnectParamMap["uniq_vpc_id"] = connectResource.SQLServerConnectParam.UniqVpcId
		}

		if connectResource.SQLServerConnectParam.IsUpdate != nil {
			sQLServerConnectParamMap["is_update"] = connectResource.SQLServerConnectParam.IsUpdate
		}

		_ = d.Set("sqlserver_connect_param", []interface{}{sQLServerConnectParamMap})
	}

	if connectResource.DorisConnectParam != nil {
		dorisConnectParamMap := map[string]interface{}{}

		if connectResource.DorisConnectParam.Port != nil {
			dorisConnectParamMap["port"] = connectResource.DorisConnectParam.Port
		}

		if connectResource.DorisConnectParam.UserName != nil {
			dorisConnectParamMap["user_name"] = connectResource.DorisConnectParam.UserName
		}

		password := connectResource.DorisConnectParam.Password
		if password != nil && *password != "" {
			dorisConnectParamMap["password"] = password
		} else {
			key := "doris_connect_param.0.password"
			if v, ok := d.GetOk(key); ok {
				dorisConnectParamMap["password"] = helper.String(v.(string))
			}
		}
		if connectResource.DorisConnectParam.Resource != nil {
			dorisConnectParamMap["resource"] = connectResource.DorisConnectParam.Resource
		}

		if connectResource.DorisConnectParam.ServiceVip != nil {
			dorisConnectParamMap["service_vip"] = connectResource.DorisConnectParam.ServiceVip
		}

		if connectResource.DorisConnectParam.UniqVpcId != nil {
			dorisConnectParamMap["uniq_vpc_id"] = connectResource.DorisConnectParam.UniqVpcId
		}

		if connectResource.DorisConnectParam.IsUpdate != nil {
			dorisConnectParamMap["is_update"] = connectResource.DorisConnectParam.IsUpdate
		}

		if connectResource.DorisConnectParam.SelfBuilt != nil {
			dorisConnectParamMap["self_built"] = connectResource.DorisConnectParam.SelfBuilt
		}

		if connectResource.DorisConnectParam.BePort != nil {
			dorisConnectParamMap["be_port"] = connectResource.DorisConnectParam.BePort
		}

		_ = d.Set("doris_connect_param", []interface{}{dorisConnectParamMap})
	}

	if connectResource.KafkaConnectParam != nil {
		kafkaConnectParamMap := map[string]interface{}{}

		if connectResource.KafkaConnectParam.Resource != nil {
			kafkaConnectParamMap["resource"] = connectResource.KafkaConnectParam.Resource
		}

		if connectResource.KafkaConnectParam.SelfBuilt != nil {
			kafkaConnectParamMap["self_built"] = connectResource.KafkaConnectParam.SelfBuilt
		}

		if connectResource.KafkaConnectParam.IsUpdate != nil {
			kafkaConnectParamMap["is_update"] = connectResource.KafkaConnectParam.IsUpdate
		}

		if connectResource.KafkaConnectParam.BrokerAddress != nil {
			kafkaConnectParamMap["broker_address"] = connectResource.KafkaConnectParam.BrokerAddress
		}

		if connectResource.KafkaConnectParam.Region != nil {
			kafkaConnectParamMap["region"] = connectResource.KafkaConnectParam.Region
		}

		_ = d.Set("kafka_connect_param", []interface{}{kafkaConnectParamMap})
	}

	return nil
}

func resourceTencentCloudCkafkaConnectResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_connect_resource.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := ckafka.NewModifyConnectResourceRequest()

	resourceId := d.Id()

	request.ResourceId = &resourceId

	if d.HasChange("resource_name") {
		if v, ok := d.GetOk("resource_name"); ok {
			request.ResourceName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("dts_connect_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "dts_connect_param"); ok {
			dtsConnectParam := ckafka.DtsModifyConnectParam{}
			if v, ok := dMap["port"]; ok {
				dtsConnectParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["group_id"]; ok {
				dtsConnectParam.GroupId = helper.String(v.(string))
			}
			if v, ok := dMap["user_name"]; ok {
				dtsConnectParam.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				dtsConnectParam.Password = helper.String(v.(string))
			}
			if v, ok := dMap["resource"]; ok {
				dtsConnectParam.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["topic"]; ok {
				dtsConnectParam.Topic = helper.String(v.(string))
			}
			if v, ok := dMap["is_update"]; ok {
				dtsConnectParam.IsUpdate = helper.Bool(v.(bool))
			}
			request.DtsConnectParam = &dtsConnectParam
		}
	}

	if d.HasChange("mongodb_connect_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "mongodb_connect_param"); ok {
			mongoDBConnectParam := ckafka.MongoDBModifyConnectParam{}
			if v, ok := dMap["port"]; ok {
				mongoDBConnectParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["user_name"]; ok {
				mongoDBConnectParam.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				mongoDBConnectParam.Password = helper.String(v.(string))
			}
			if v, ok := dMap["resource"]; ok {
				mongoDBConnectParam.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["self_built"]; ok {
				mongoDBConnectParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := dMap["service_vip"]; ok {
				mongoDBConnectParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				mongoDBConnectParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["is_update"]; ok {
				mongoDBConnectParam.IsUpdate = helper.Bool(v.(bool))
			}
			request.MongoDBConnectParam = &mongoDBConnectParam
		}
	}

	if d.HasChange("es_connect_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "es_connect_param"); ok {
			esConnectParam := ckafka.EsModifyConnectParam{}
			if v, ok := dMap["port"]; ok {
				esConnectParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["user_name"]; ok {
				esConnectParam.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				esConnectParam.Password = helper.String(v.(string))
			}
			if v, ok := dMap["resource"]; ok {
				esConnectParam.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["self_built"]; ok {
				esConnectParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := dMap["service_vip"]; ok {
				esConnectParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				esConnectParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["is_update"]; ok {
				esConnectParam.IsUpdate = helper.Bool(v.(bool))
			}
			request.EsConnectParam = &esConnectParam
		}
	}

	if d.HasChange("clickhouse_connect_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "clickhouse_connect_param"); ok {
			clickHouseConnectParam := ckafka.ClickHouseModifyConnectParam{}
			if v, ok := dMap["port"]; ok {
				clickHouseConnectParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["user_name"]; ok {
				clickHouseConnectParam.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				clickHouseConnectParam.Password = helper.String(v.(string))
			}
			if v, ok := dMap["resource"]; ok {
				clickHouseConnectParam.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["self_built"]; ok {
				clickHouseConnectParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := dMap["service_vip"]; ok {
				clickHouseConnectParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				clickHouseConnectParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["is_update"]; ok {
				clickHouseConnectParam.IsUpdate = helper.Bool(v.(bool))
			}
			request.ClickHouseConnectParam = &clickHouseConnectParam
		}
	}

	if d.HasChange("mysql_connect_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "mysql_connect_param"); ok {
			mySQLConnectParam := ckafka.MySQLModifyConnectParam{}
			if v, ok := dMap["port"]; ok {
				mySQLConnectParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["user_name"]; ok {
				mySQLConnectParam.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				mySQLConnectParam.Password = helper.String(v.(string))
			}
			if v, ok := dMap["resource"]; ok {
				mySQLConnectParam.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["service_vip"]; ok {
				mySQLConnectParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				mySQLConnectParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["is_update"]; ok {
				mySQLConnectParam.IsUpdate = helper.Bool(v.(bool))
			}
			if v, ok := dMap["cluster_id"]; ok {
				mySQLConnectParam.ClusterId = helper.String(v.(string))
			}
			if v, ok := dMap["self_built"]; ok {
				mySQLConnectParam.SelfBuilt = helper.Bool(v.(bool))
			}
			request.MySQLConnectParam = &mySQLConnectParam
		}
	}

	if d.HasChange("postgresql_connect_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "postgresql_connect_param"); ok {
			postgreSQLConnectParam := ckafka.PostgreSQLModifyConnectParam{}
			if v, ok := dMap["port"]; ok {
				postgreSQLConnectParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["user_name"]; ok {
				postgreSQLConnectParam.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				postgreSQLConnectParam.Password = helper.String(v.(string))
			}
			if v, ok := dMap["resource"]; ok {
				postgreSQLConnectParam.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["service_vip"]; ok {
				postgreSQLConnectParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				postgreSQLConnectParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["cluster_id"]; ok {
				postgreSQLConnectParam.ClusterId = helper.String(v.(string))
			}
			if v, ok := dMap["is_update"]; ok {
				postgreSQLConnectParam.IsUpdate = helper.Bool(v.(bool))
			}
			if v, ok := dMap["self_built"]; ok {
				postgreSQLConnectParam.SelfBuilt = helper.Bool(v.(bool))
			}
			request.PostgreSQLConnectParam = &postgreSQLConnectParam
		}
	}

	if d.HasChange("mariadb_connect_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "mariadb_connect_param"); ok {
			mariaDBConnectParam := ckafka.MariaDBModifyConnectParam{}
			if v, ok := dMap["port"]; ok {
				mariaDBConnectParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["user_name"]; ok {
				mariaDBConnectParam.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				mariaDBConnectParam.Password = helper.String(v.(string))
			}
			if v, ok := dMap["resource"]; ok {
				mariaDBConnectParam.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["service_vip"]; ok {
				mariaDBConnectParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				mariaDBConnectParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["is_update"]; ok {
				mariaDBConnectParam.IsUpdate = helper.Bool(v.(bool))
			}
			request.MariaDBConnectParam = &mariaDBConnectParam
		}
	}

	if d.HasChange("sqlserver_connect_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "sqlserver_connect_param"); ok {
			sQLServerConnectParam := ckafka.SQLServerModifyConnectParam{}
			if v, ok := dMap["port"]; ok {
				sQLServerConnectParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["user_name"]; ok {
				sQLServerConnectParam.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				sQLServerConnectParam.Password = helper.String(v.(string))
			}
			if v, ok := dMap["resource"]; ok {
				sQLServerConnectParam.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["service_vip"]; ok {
				sQLServerConnectParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				sQLServerConnectParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["is_update"]; ok {
				sQLServerConnectParam.IsUpdate = helper.Bool(v.(bool))
			}
			request.SQLServerConnectParam = &sQLServerConnectParam
		}
	}

	if d.HasChange("doris_connect_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "doris_connect_param"); ok {
			dorisConnectParam := ckafka.DorisModifyConnectParam{}
			if v, ok := dMap["port"]; ok {
				dorisConnectParam.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["user_name"]; ok {
				dorisConnectParam.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				dorisConnectParam.Password = helper.String(v.(string))
			}
			if v, ok := dMap["resource"]; ok {
				dorisConnectParam.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["service_vip"]; ok {
				dorisConnectParam.ServiceVip = helper.String(v.(string))
			}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				dorisConnectParam.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["is_update"]; ok {
				dorisConnectParam.IsUpdate = helper.Bool(v.(bool))
			}
			if v, ok := dMap["self_built"]; ok {
				dorisConnectParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := dMap["be_port"]; ok {
				dorisConnectParam.BePort = helper.IntInt64(v.(int))
			}
			request.DorisConnectParam = &dorisConnectParam
		}
	}

	if d.HasChange("kafka_connect_param") {
		if dMap, ok := helper.InterfacesHeadMap(d, "kafka_connect_param"); ok {
			kafkaConnectParam := ckafka.KafkaConnectParam{}
			if v, ok := dMap["resource"]; ok {
				kafkaConnectParam.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["self_built"]; ok {
				kafkaConnectParam.SelfBuilt = helper.Bool(v.(bool))
			}
			if v, ok := dMap["is_update"]; ok {
				kafkaConnectParam.IsUpdate = helper.Bool(v.(bool))
			}
			if v, ok := dMap["broker_address"]; ok {
				kafkaConnectParam.BrokerAddress = helper.String(v.(string))
			}
			if v, ok := dMap["region"]; ok {
				kafkaConnectParam.Region = helper.String(v.(string))
			}
			request.KafkaConnectParam = &kafkaConnectParam
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().ModifyConnectResource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ckafka connectResource failed, reason:%+v", logId, err)
		return err
	}

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 2*tccommon.ReadRetryTimeout, time.Second, service.CkafkaConnectResourceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCkafkaConnectResourceRead(d, meta)
}

func resourceTencentCloudCkafkaConnectResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_connect_resource.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	resourceId := d.Id()

	if err := service.DeleteCkafkaConnectResourceById(ctx, resourceId); err != nil {
		return err
	}

	return nil
}
