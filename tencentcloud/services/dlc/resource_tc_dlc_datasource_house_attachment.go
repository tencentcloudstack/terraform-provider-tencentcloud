package dlc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcDatasourceHouseAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcDatasourceHouseAttachmentCreate,
		Read:   resourceTencentCloudDlcDatasourceHouseAttachmentRead,
		Update: resourceTencentCloudDlcDatasourceHouseAttachmentUpdate,
		Delete: resourceTencentCloudDlcDatasourceHouseAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"datasource_connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network configuration name.",
			},

			"datasource_connection_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Data source type. Allow value: Mysql, HiveCos, HiveHdfs, HiveCHdfs, Kafka, OtherDatasourceConnection, PostgreSql, SqlServer, ClickHouse, Elasticsearch, TDSQLPostgreSql, TCHouseD, TccHive.",
			},

			"datasource_connection_config": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Data source network configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mysql": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Properties of MySQL data source connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// "jdbc_url": {
									// 	Type:        schema.TypeString,
									// 	Required:    true,
									// 	ForceNew:    true,
									// 	Description: "JDBC URL for connecting to MySQL.",
									// },
									// "user": {
									// 	Type:        schema.TypeString,
									// 	Required:    true,
									// 	ForceNew:    true,
									// 	Description: "Username.",
									// },
									// "password": {
									// 	Type:        schema.TypeString,
									// 	Required:    true,
									// 	ForceNew:    true,
									// 	Description: "MySQL password.",
									// },
									"location": {
										Type:        schema.TypeList,
										Required:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "Network information for MySQL data source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC IPv4 CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet IPv4 CIDR.",
												},
											},
										},
									},
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Database name.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Database instance ID, consistent with the database side.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Database instance name, consistent with the database side.",
									},
								},
							},
						},
						"hive": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Properties of Hive data source connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"meta_store_url": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Address of Hive metastore.",
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Hive data source type, representing data storage location, COS or HDFS.",
									},
									"location": {
										Type:        schema.TypeList,
										Required:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "Private network information where the data source is located.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC IPv4 CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet IPv4 CIDR.",
												},
											},
										},
									},
									// "user": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "If the type is HDFS, a username is required.",
									// },
									"high_availability": {
										Type:        schema.TypeBool,
										Optional:    true,
										ForceNew:    true,
										Description: "If the type is HDFS, high availability needs to be selected.",
									},
									"bucket_url": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "If the type is COS, COS bucket connection needs to be filled in.",
									},
									"hdfs_properties": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "JSON string. If the type is HDFS, this field needs to be filled in.",
									},
									"mysql": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "Metadata database information for Hive.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												// "jdbc_url": {
												// 	Type:        schema.TypeString,
												// 	Required:    true,
												// 	ForceNew:    true,
												// 	Description: "JDBC URL for connecting to MySQL.",
												// },
												// "user": {
												// 	Type:        schema.TypeString,
												// 	Required:    true,
												// 	ForceNew:    true,
												// 	Description: "Username.",
												// },
												// "password": {
												// 	Type:        schema.TypeString,
												// 	Required:    true,
												// 	ForceNew:    true,
												// 	Description: "MySQL password.",
												// },
												"location": {
													Type:        schema.TypeList,
													Required:    true,
													ForceNew:    true,
													MaxItems:    1,
													Description: "Network information for MySQL data source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"vpc_id": {
																Type:        schema.TypeString,
																Required:    true,
																ForceNew:    true,
																Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
															},
															"vpc_cidr_block": {
																Type:        schema.TypeString,
																Required:    true,
																ForceNew:    true,
																Description: "VPC IPv4 CIDR.",
															},
															"subnet_id": {
																Type:        schema.TypeString,
																Required:    true,
																ForceNew:    true,
																Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
															},
															"subnet_cidr_block": {
																Type:        schema.TypeString,
																Required:    true,
																ForceNew:    true,
																Description: "Subnet IPv4 CIDR.",
															},
														},
													},
												},
												"db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Database name.",
												},
												"instance_id": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Database instance ID, consistent with the database side.",
												},
												"instance_name": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Database instance name, consistent with the database side.",
												},
											},
										},
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "EMR cluster ID.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "EMR cluster name.",
									},
									"hive_version": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Version number of Hive component in EMR cluster.",
									},
									"kerberos_info": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "Kerberos details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"krb5_conf": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Krb5Conf file value.",
												},
												"key_tab": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "KeyTab file value.",
												},
												"service_principal": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Service principal.",
												},
											},
										},
									},
									"kerberos_enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										ForceNew:    true,
										Description: "Whether to enable Kerberos.",
									},
								},
							},
						},
						"kafka": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Properties of Kafka data source connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Kafka instance ID.",
									},
									"location": {
										Type:        schema.TypeList,
										Required:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "Network information for Kafka data source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC IPv4 CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet IPv4 CIDR.",
												},
											},
										},
									},
								},
							},
						},
						"other_datasource_connection": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Properties of other data source connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"location": {
										Type:        schema.TypeList,
										Required:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "Network parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC IPv4 CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet IPv4 CIDR.",
												},
											},
										},
									},
								},
							},
						},
						"postgre_sql": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Properties of PostgreSQL data source connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Unique ID of the data source instance.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Name of the data source.",
									},
									// "jdbc_url": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "JDBC access link for the data source.",
									// },
									// "user": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Username for accessing the data source.",
									// },
									// "password": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Data source access password, requires base64 encoding.",
									// },
									"location": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "VPC and subnet information for the data source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC IPv4 CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet IPv4 CIDR.",
												},
											},
										},
									},
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Default database name.",
									},
								},
							},
						},
						"sql_server": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Properties of SQLServer data source connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Unique ID of the data source instance.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Name of the data source.",
									},
									// "jdbc_url": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "JDBC access link for the data source.",
									// },
									// "user": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Username for accessing the data source.",
									// },
									// "password": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Data source access password, requires base64 encoding.",
									// },
									"location": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "VPC and subnet information for the data source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC IPv4 CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet IPv4 CIDR.",
												},
											},
										},
									},
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Default database name.",
									},
								},
							},
						},
						"click_house": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Properties of ClickHouse data source connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Unique ID of the data source instance.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Name of the data source.",
									},
									// "jdbc_url": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "JDBC access link for the data source.",
									// },
									// "user": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Username for accessing the data source.",
									// },
									// "password": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Data source access password, requires base64 encoding.",
									// },
									"location": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "VPC and subnet information for the data source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC IPv4 CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet IPv4 CIDR.",
												},
											},
										},
									},
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Default database name.",
									},
								},
							},
						},
						"elasticsearch": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Properties of Elasticsearch data source connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Data source ID.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Data source name.",
									},
									// "user": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Username.",
									// },
									// "password": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Password, requires base64 encoding.",
									// },
									"location": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "VPC and subnet information for the data source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC IPv4 CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet IPv4 CIDR.",
												},
											},
										},
									},
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Default database name.",
									},
									"service_info": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										Description: "IP and port information for accessing Elasticsearch.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "IP information.",
												},
												"port": {
													Type:        schema.TypeInt,
													Optional:    true,
													ForceNew:    true,
													Description: "Port information.",
												},
											},
										},
									},
								},
							},
						},
						"tdsql_postgre_sql": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Properties of TDSQL-PostgreSQL data source connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Unique ID of the data source instance.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Name of the data source.",
									},
									// "jdbc_url": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "JDBC access link for the data source.",
									// },
									// "user": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Username for accessing the data source.",
									// },
									// "password": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Data source access password, requires base64 encoding.",
									// },
									"location": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "VPC and subnet information for the data source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC IPv4 CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet IPv4 CIDR.",
												},
											},
										},
									},
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Default database name.",
									},
								},
							},
						},
						"tc_house_d": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Properties of Doris data source connection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Unique ID of the data source instance.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Data source name.",
									},
									// "jdbc_url": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "JDBC of the data source.",
									// },
									// "user": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "User for accessing the data source.",
									// },
									// "password": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	ForceNew:    true,
									// 	Description: "Data source access password, requires base64 encoding.",
									// },
									"location": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "VPC and subnet information for the data source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "VPC IPv4 CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Subnet IPv4 CIDR.",
												},
											},
										},
									},
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Default database name.",
									},
									"access_info": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Access information.",
									},
								},
							},
						},
						"tcc_hive": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "TccHive data catalog connection information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Instance ID.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Instance name.",
									},
									"endpoint_service_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Endpoint service ID.",
									},
									"meta_store_url": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Thrift connection address.",
									},
									"hive_version": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Hive version.",
									},
									"tcc_connection": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "Network information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"clb_ip": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Service CLB IP.",
												},
												"clb_port": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Service CLB port.",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "VPC instance ID.",
												},
												"vpc_cidr_block": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "VPC CIDR.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Subnet instance ID.",
												},
												"subnet_cidr_block": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "Subnet CIDR.",
												},
											},
										},
									},
									"hms_endpoint_service_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "HMS endpoint service ID.",
									},
								},
							},
						},
					},
				},
			},

			"data_engine_names": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Engine name, only one engine can be bound.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"network_connection_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Network type, 2-cross-source type, 4-enhanced type.",
			},

			"network_connection_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network configuration description.",
			},
		},
	}
}

func resourceTencentCloudDlcDatasourceHouseAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_datasource_house_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                    = tccommon.GetLogId(tccommon.ContextNil)
		ctx                      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request                  = dlcv20210125.NewAssociateDatasourceHouseRequest()
		datasourceConnectionName string
	)

	if v, ok := d.GetOk("datasource_connection_name"); ok {
		request.DatasourceConnectionName = helper.String(v.(string))
		datasourceConnectionName = v.(string)
	}

	if v, ok := d.GetOk("datasource_connection_type"); ok {
		request.DatasourceConnectionType = helper.String(v.(string))
	}

	if datasourceConnectionConfigMap, ok := helper.InterfacesHeadMap(d, "datasource_connection_config"); ok {
		datasourceConnectionConfig := dlcv20210125.DatasourceConnectionConfig{}
		if mysqlMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["mysql"]); ok {
			mysqlInfo := dlcv20210125.MysqlInfo{}
			mysqlInfo.JdbcUrl = helper.String("")
			mysqlInfo.User = helper.String("")
			mysqlInfo.Password = helper.String("")

			if locationMap, ok := helper.ConvertInterfacesHeadToMap(mysqlMap["location"]); ok {
				datasourceConnectionLocation := dlcv20210125.DatasourceConnectionLocation{}
				if v, ok := locationMap["vpc_id"].(string); ok {
					datasourceConnectionLocation.VpcId = helper.String(v)
				}

				if v, ok := locationMap["vpc_cidr_block"].(string); ok {
					datasourceConnectionLocation.VpcCidrBlock = helper.String(v)
				}

				if v, ok := locationMap["subnet_id"].(string); ok {
					datasourceConnectionLocation.SubnetId = helper.String(v)
				}

				if v, ok := locationMap["subnet_cidr_block"].(string); ok {
					datasourceConnectionLocation.SubnetCidrBlock = helper.String(v)
				}

				mysqlInfo.Location = &datasourceConnectionLocation
			}

			if v, ok := mysqlMap["db_name"].(string); ok {
				mysqlInfo.DbName = helper.String(v)
			}

			if v, ok := mysqlMap["instance_id"].(string); ok {
				mysqlInfo.InstanceId = helper.String(v)
			}

			if v, ok := mysqlMap["instance_name"].(string); ok {
				mysqlInfo.InstanceName = helper.String(v)
			}

			datasourceConnectionConfig.Mysql = &mysqlInfo
		}

		if hiveMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["hive"]); ok {
			hiveInfo := dlcv20210125.HiveInfo{}
			if v, ok := hiveMap["meta_store_url"].(string); ok {
				hiveInfo.MetaStoreUrl = helper.String(v)
			}

			if v, ok := hiveMap["type"].(string); ok {
				hiveInfo.Type = helper.String(v)
			}

			if locationMap, ok := helper.ConvertInterfacesHeadToMap(hiveMap["location"]); ok {
				datasourceConnectionLocation2 := dlcv20210125.DatasourceConnectionLocation{}
				if v, ok := locationMap["vpc_id"].(string); ok {
					datasourceConnectionLocation2.VpcId = helper.String(v)
				}

				if v, ok := locationMap["vpc_cidr_block"].(string); ok {
					datasourceConnectionLocation2.VpcCidrBlock = helper.String(v)
				}

				if v, ok := locationMap["subnet_id"].(string); ok {
					datasourceConnectionLocation2.SubnetId = helper.String(v)
				}

				if v, ok := locationMap["subnet_cidr_block"].(string); ok {
					datasourceConnectionLocation2.SubnetCidrBlock = helper.String(v)
				}

				hiveInfo.Location = &datasourceConnectionLocation2
			}

			hiveInfo.User = helper.String("")

			if v, ok := hiveMap["high_availability"].(bool); ok {
				hiveInfo.HighAvailability = helper.Bool(v)
			}

			if v, ok := hiveMap["bucket_url"].(string); ok {
				hiveInfo.BucketUrl = helper.String(v)
			}

			if v, ok := hiveMap["hdfs_properties"].(string); ok {
				hiveInfo.HdfsProperties = helper.String(v)
			}

			if mysqlMap, ok := helper.ConvertInterfacesHeadToMap(hiveMap["mysql"]); ok {
				mysqlInfo2 := dlcv20210125.MysqlInfo{}
				mysqlInfo2.JdbcUrl = helper.String("")
				mysqlInfo2.User = helper.String("")
				mysqlInfo2.Password = helper.String("")

				if locationMap, ok := helper.ConvertInterfacesHeadToMap(mysqlMap["location"]); ok {
					datasourceConnectionLocation3 := dlcv20210125.DatasourceConnectionLocation{}
					if v, ok := locationMap["vpc_id"].(string); ok {
						datasourceConnectionLocation3.VpcId = helper.String(v)
					}

					if v, ok := locationMap["vpc_cidr_block"].(string); ok {
						datasourceConnectionLocation3.VpcCidrBlock = helper.String(v)
					}

					if v, ok := locationMap["subnet_id"].(string); ok {
						datasourceConnectionLocation3.SubnetId = helper.String(v)
					}

					if v, ok := locationMap["subnet_cidr_block"].(string); ok {
						datasourceConnectionLocation3.SubnetCidrBlock = helper.String(v)
					}

					mysqlInfo2.Location = &datasourceConnectionLocation3
				}

				if v, ok := mysqlMap["db_name"].(string); ok {
					mysqlInfo2.DbName = helper.String(v)
				}

				if v, ok := mysqlMap["instance_id"].(string); ok {
					mysqlInfo2.InstanceId = helper.String(v)
				}

				if v, ok := mysqlMap["instance_name"].(string); ok {
					mysqlInfo2.InstanceName = helper.String(v)
				}

				hiveInfo.Mysql = &mysqlInfo2
			}

			if v, ok := hiveMap["instance_id"].(string); ok {
				hiveInfo.InstanceId = helper.String(v)
			}

			if v, ok := hiveMap["instance_name"].(string); ok {
				hiveInfo.InstanceName = helper.String(v)
			}

			if v, ok := hiveMap["hive_version"].(string); ok {
				hiveInfo.HiveVersion = helper.String(v)
			}

			if kerberosInfoMap, ok := helper.ConvertInterfacesHeadToMap(hiveMap["kerberos_info"]); ok {
				kerberosInfo := dlcv20210125.KerberosInfo{}
				if v, ok := kerberosInfoMap["krb5_conf"].(string); ok {
					kerberosInfo.Krb5Conf = helper.String(v)
				}

				if v, ok := kerberosInfoMap["key_tab"].(string); ok {
					kerberosInfo.KeyTab = helper.String(v)
				}

				if v, ok := kerberosInfoMap["service_principal"].(string); ok {
					kerberosInfo.ServicePrincipal = helper.String(v)
				}

				hiveInfo.KerberosInfo = &kerberosInfo
			}

			if v, ok := hiveMap["kerberos_enable"].(bool); ok {
				hiveInfo.KerberosEnable = helper.Bool(v)
			}

			datasourceConnectionConfig.Hive = &hiveInfo
		}

		if kafkaMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["kafka"]); ok {
			kafkaInfo := dlcv20210125.KafkaInfo{}
			if v, ok := kafkaMap["instance_id"].(string); ok {
				kafkaInfo.InstanceId = helper.String(v)
			}

			if locationMap, ok := helper.ConvertInterfacesHeadToMap(kafkaMap["location"]); ok {
				datasourceConnectionLocation4 := dlcv20210125.DatasourceConnectionLocation{}
				if v, ok := locationMap["vpc_id"].(string); ok {
					datasourceConnectionLocation4.VpcId = helper.String(v)
				}

				if v, ok := locationMap["vpc_cidr_block"].(string); ok {
					datasourceConnectionLocation4.VpcCidrBlock = helper.String(v)
				}

				if v, ok := locationMap["subnet_id"].(string); ok {
					datasourceConnectionLocation4.SubnetId = helper.String(v)
				}

				if v, ok := locationMap["subnet_cidr_block"].(string); ok {
					datasourceConnectionLocation4.SubnetCidrBlock = helper.String(v)
				}

				kafkaInfo.Location = &datasourceConnectionLocation4
			}

			datasourceConnectionConfig.Kafka = &kafkaInfo
		}

		if otherDatasourceConnectionMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["other_datasource_connection"]); ok {
			otherDatasourceConnection := dlcv20210125.OtherDatasourceConnection{}
			if locationMap, ok := helper.ConvertInterfacesHeadToMap(otherDatasourceConnectionMap["location"]); ok {
				datasourceConnectionLocation5 := dlcv20210125.DatasourceConnectionLocation{}
				if v, ok := locationMap["vpc_id"].(string); ok {
					datasourceConnectionLocation5.VpcId = helper.String(v)
				}

				if v, ok := locationMap["vpc_cidr_block"].(string); ok {
					datasourceConnectionLocation5.VpcCidrBlock = helper.String(v)
				}

				if v, ok := locationMap["subnet_id"].(string); ok {
					datasourceConnectionLocation5.SubnetId = helper.String(v)
				}

				if v, ok := locationMap["subnet_cidr_block"].(string); ok {
					datasourceConnectionLocation5.SubnetCidrBlock = helper.String(v)
				}

				otherDatasourceConnection.Location = &datasourceConnectionLocation5
			}

			datasourceConnectionConfig.OtherDatasourceConnection = &otherDatasourceConnection
		}

		if postgreSqlMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["postgre_sql"]); ok {
			dataSourceInfo := dlcv20210125.DataSourceInfo{}
			if v, ok := postgreSqlMap["instance_id"].(string); ok {
				dataSourceInfo.InstanceId = helper.String(v)
			}

			if v, ok := postgreSqlMap["instance_name"].(string); ok {
				dataSourceInfo.InstanceName = helper.String(v)
			}

			dataSourceInfo.JdbcUrl = helper.String("")
			dataSourceInfo.User = helper.String("")
			dataSourceInfo.Password = helper.String("")

			if locationMap, ok := helper.ConvertInterfacesHeadToMap(postgreSqlMap["location"]); ok {
				datasourceConnectionLocation6 := dlcv20210125.DatasourceConnectionLocation{}
				if v, ok := locationMap["vpc_id"].(string); ok {
					datasourceConnectionLocation6.VpcId = helper.String(v)
				}

				if v, ok := locationMap["vpc_cidr_block"].(string); ok {
					datasourceConnectionLocation6.VpcCidrBlock = helper.String(v)
				}

				if v, ok := locationMap["subnet_id"].(string); ok {
					datasourceConnectionLocation6.SubnetId = helper.String(v)
				}

				if v, ok := locationMap["subnet_cidr_block"].(string); ok {
					datasourceConnectionLocation6.SubnetCidrBlock = helper.String(v)
				}

				dataSourceInfo.Location = &datasourceConnectionLocation6
			}

			if v, ok := postgreSqlMap["db_name"].(string); ok {
				dataSourceInfo.DbName = helper.String(v)
			}

			datasourceConnectionConfig.PostgreSql = &dataSourceInfo
		}

		if sqlServerMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["sql_server"]); ok {
			dataSourceInfo2 := dlcv20210125.DataSourceInfo{}
			if v, ok := sqlServerMap["instance_id"].(string); ok {
				dataSourceInfo2.InstanceId = helper.String(v)
			}

			if v, ok := sqlServerMap["instance_name"].(string); ok {
				dataSourceInfo2.InstanceName = helper.String(v)
			}

			dataSourceInfo2.JdbcUrl = helper.String("")
			dataSourceInfo2.User = helper.String("")
			dataSourceInfo2.Password = helper.String("")

			if locationMap, ok := helper.ConvertInterfacesHeadToMap(sqlServerMap["location"]); ok {
				datasourceConnectionLocation7 := dlcv20210125.DatasourceConnectionLocation{}
				if v, ok := locationMap["vpc_id"].(string); ok {
					datasourceConnectionLocation7.VpcId = helper.String(v)
				}

				if v, ok := locationMap["vpc_cidr_block"].(string); ok {
					datasourceConnectionLocation7.VpcCidrBlock = helper.String(v)
				}

				if v, ok := locationMap["subnet_id"].(string); ok {
					datasourceConnectionLocation7.SubnetId = helper.String(v)
				}

				if v, ok := locationMap["subnet_cidr_block"].(string); ok {
					datasourceConnectionLocation7.SubnetCidrBlock = helper.String(v)
				}

				dataSourceInfo2.Location = &datasourceConnectionLocation7
			}

			if v, ok := sqlServerMap["db_name"].(string); ok {
				dataSourceInfo2.DbName = helper.String(v)
			}

			datasourceConnectionConfig.SqlServer = &dataSourceInfo2
		}

		if clickHouseMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["click_house"]); ok {
			dataSourceInfo3 := dlcv20210125.DataSourceInfo{}
			if v, ok := clickHouseMap["instance_id"].(string); ok {
				dataSourceInfo3.InstanceId = helper.String(v)
			}

			if v, ok := clickHouseMap["instance_name"].(string); ok {
				dataSourceInfo3.InstanceName = helper.String(v)
			}

			dataSourceInfo3.JdbcUrl = helper.String("")
			dataSourceInfo3.User = helper.String("")
			dataSourceInfo3.Password = helper.String("")

			if locationMap, ok := helper.ConvertInterfacesHeadToMap(clickHouseMap["location"]); ok {
				datasourceConnectionLocation8 := dlcv20210125.DatasourceConnectionLocation{}
				if v, ok := locationMap["vpc_id"].(string); ok {
					datasourceConnectionLocation8.VpcId = helper.String(v)
				}

				if v, ok := locationMap["vpc_cidr_block"].(string); ok {
					datasourceConnectionLocation8.VpcCidrBlock = helper.String(v)
				}

				if v, ok := locationMap["subnet_id"].(string); ok {
					datasourceConnectionLocation8.SubnetId = helper.String(v)
				}

				if v, ok := locationMap["subnet_cidr_block"].(string); ok {
					datasourceConnectionLocation8.SubnetCidrBlock = helper.String(v)
				}

				dataSourceInfo3.Location = &datasourceConnectionLocation8
			}

			if v, ok := clickHouseMap["db_name"].(string); ok {
				dataSourceInfo3.DbName = helper.String(v)
			}

			datasourceConnectionConfig.ClickHouse = &dataSourceInfo3
		}

		if elasticsearchMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["elasticsearch"]); ok {
			elasticsearchInfo := dlcv20210125.ElasticsearchInfo{}
			if v, ok := elasticsearchMap["instance_id"].(string); ok {
				elasticsearchInfo.InstanceId = helper.String(v)
			}

			if v, ok := elasticsearchMap["instance_name"].(string); ok {
				elasticsearchInfo.InstanceName = helper.String(v)
			}

			elasticsearchInfo.User = helper.String("")
			elasticsearchInfo.Password = helper.String("")

			if locationMap, ok := helper.ConvertInterfacesHeadToMap(elasticsearchMap["location"]); ok {
				datasourceConnectionLocation9 := dlcv20210125.DatasourceConnectionLocation{}
				if v, ok := locationMap["vpc_id"].(string); ok {
					datasourceConnectionLocation9.VpcId = helper.String(v)
				}

				if v, ok := locationMap["vpc_cidr_block"].(string); ok {
					datasourceConnectionLocation9.VpcCidrBlock = helper.String(v)
				}

				if v, ok := locationMap["subnet_id"].(string); ok {
					datasourceConnectionLocation9.SubnetId = helper.String(v)
				}

				if v, ok := locationMap["subnet_cidr_block"].(string); ok {
					datasourceConnectionLocation9.SubnetCidrBlock = helper.String(v)
				}

				elasticsearchInfo.Location = &datasourceConnectionLocation9
			}

			if v, ok := elasticsearchMap["db_name"].(string); ok {
				elasticsearchInfo.DbName = helper.String(v)
			}

			if v, ok := elasticsearchMap["service_info"]; ok {
				for _, item := range v.([]interface{}) {
					serviceInfoMap := item.(map[string]interface{})
					ipPortPair := dlcv20210125.IpPortPair{}
					if v, ok := serviceInfoMap["ip"].(string); ok {
						ipPortPair.Ip = helper.String(v)
					}

					if v, ok := serviceInfoMap["port"].(int); ok {
						ipPortPair.Port = helper.IntInt64(v)
					}

					elasticsearchInfo.ServiceInfo = append(elasticsearchInfo.ServiceInfo, &ipPortPair)
				}
			}

			datasourceConnectionConfig.Elasticsearch = &elasticsearchInfo
		}

		if tDSQLPostgreSqlMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["tdsql_postgre_sql"]); ok {
			dataSourceInfo4 := dlcv20210125.DataSourceInfo{}
			if v, ok := tDSQLPostgreSqlMap["instance_id"].(string); ok {
				dataSourceInfo4.InstanceId = helper.String(v)
			}

			if v, ok := tDSQLPostgreSqlMap["instance_name"].(string); ok {
				dataSourceInfo4.InstanceName = helper.String(v)
			}

			dataSourceInfo4.JdbcUrl = helper.String("")
			dataSourceInfo4.User = helper.String("")
			dataSourceInfo4.Password = helper.String("")

			if locationMap, ok := helper.ConvertInterfacesHeadToMap(tDSQLPostgreSqlMap["location"]); ok {
				datasourceConnectionLocation10 := dlcv20210125.DatasourceConnectionLocation{}
				if v, ok := locationMap["vpc_id"].(string); ok {
					datasourceConnectionLocation10.VpcId = helper.String(v)
				}

				if v, ok := locationMap["vpc_cidr_block"].(string); ok {
					datasourceConnectionLocation10.VpcCidrBlock = helper.String(v)
				}

				if v, ok := locationMap["subnet_id"].(string); ok {
					datasourceConnectionLocation10.SubnetId = helper.String(v)
				}

				if v, ok := locationMap["subnet_cidr_block"].(string); ok {
					datasourceConnectionLocation10.SubnetCidrBlock = helper.String(v)
				}

				dataSourceInfo4.Location = &datasourceConnectionLocation10
			}

			if v, ok := tDSQLPostgreSqlMap["db_name"].(string); ok {
				dataSourceInfo4.DbName = helper.String(v)
			}

			datasourceConnectionConfig.TDSQLPostgreSql = &dataSourceInfo4
		}

		if tCHouseDMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["tc_house_d"]); ok {
			tCHouseD := dlcv20210125.TCHouseD{}
			if v, ok := tCHouseDMap["instance_id"].(string); ok {
				tCHouseD.InstanceId = helper.String(v)
			}

			if v, ok := tCHouseDMap["instance_name"].(string); ok {
				tCHouseD.InstanceName = helper.String(v)
			}

			tCHouseD.JdbcUrl = helper.String("")
			tCHouseD.User = helper.String("")
			tCHouseD.Password = helper.String("")

			if locationMap, ok := helper.ConvertInterfacesHeadToMap(tCHouseDMap["location"]); ok {
				datasourceConnectionLocation11 := dlcv20210125.DatasourceConnectionLocation{}
				if v, ok := locationMap["vpc_id"].(string); ok {
					datasourceConnectionLocation11.VpcId = helper.String(v)
				}

				if v, ok := locationMap["vpc_cidr_block"].(string); ok {
					datasourceConnectionLocation11.VpcCidrBlock = helper.String(v)
				}

				if v, ok := locationMap["subnet_id"].(string); ok {
					datasourceConnectionLocation11.SubnetId = helper.String(v)
				}

				if v, ok := locationMap["subnet_cidr_block"].(string); ok {
					datasourceConnectionLocation11.SubnetCidrBlock = helper.String(v)
				}

				tCHouseD.Location = &datasourceConnectionLocation11
			}

			if v, ok := tCHouseDMap["db_name"].(string); ok {
				tCHouseD.DbName = helper.String(v)
			}

			if v, ok := tCHouseDMap["access_info"].(string); ok {
				tCHouseD.AccessInfo = helper.String(v)
			}

			datasourceConnectionConfig.TCHouseD = &tCHouseD
		}

		if tccHiveMap, ok := helper.ConvertInterfacesHeadToMap(datasourceConnectionConfigMap["tcc_hive"]); ok {
			tccHive := dlcv20210125.TccHive{}
			if v, ok := tccHiveMap["instance_id"].(string); ok {
				tccHive.InstanceId = helper.String(v)
			}

			if v, ok := tccHiveMap["instance_name"].(string); ok {
				tccHive.InstanceName = helper.String(v)
			}

			if v, ok := tccHiveMap["endpoint_service_id"].(string); ok {
				tccHive.EndpointServiceId = helper.String(v)
			}

			if v, ok := tccHiveMap["meta_store_url"].(string); ok {
				tccHive.MetaStoreUrl = helper.String(v)
			}

			if v, ok := tccHiveMap["hive_version"].(string); ok {
				tccHive.HiveVersion = helper.String(v)
			}

			if tccConnectionMap, ok := helper.ConvertInterfacesHeadToMap(tccHiveMap["tcc_connection"]); ok {
				netWork := dlcv20210125.NetWork{}
				if v, ok := tccConnectionMap["clb_ip"].(string); ok {
					netWork.ClbIp = helper.String(v)
				}

				if v, ok := tccConnectionMap["clb_port"].(string); ok {
					netWork.ClbPort = helper.String(v)
				}

				if v, ok := tccConnectionMap["vpc_id"].(string); ok {
					netWork.VpcId = helper.String(v)
				}

				if v, ok := tccConnectionMap["vpc_cidr_block"].(string); ok {
					netWork.VpcCidrBlock = helper.String(v)
				}

				if v, ok := tccConnectionMap["subnet_id"].(string); ok {
					netWork.SubnetId = helper.String(v)
				}

				if v, ok := tccConnectionMap["subnet_cidr_block"].(string); ok {
					netWork.SubnetCidrBlock = helper.String(v)
				}

				tccHive.TccConnection = &netWork
			}

			if v, ok := tccHiveMap["hms_endpoint_service_id"].(string); ok {
				tccHive.HmsEndpointServiceId = helper.String(v)
			}

			datasourceConnectionConfig.TccHive = &tccHive
		}

		request.DatasourceConnectionConfig = &datasourceConnectionConfig
	}

	if v, ok := d.GetOk("data_engine_names"); ok {
		dataEngineNamesSet := v.(*schema.Set).List()
		for i := range dataEngineNamesSet {
			dataEngineNames := dataEngineNamesSet[i].(string)
			request.DataEngineNames = append(request.DataEngineNames, helper.String(dataEngineNames))
		}
	}

	if v, ok := d.GetOkExists("network_connection_type"); ok {
		request.NetworkConnectionType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("network_connection_desc"); ok {
		request.NetworkConnectionDesc = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().AssociateDatasourceHouseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc datasource house attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(datasourceConnectionName)

	// wait
	waitReq := dlcv20210125.NewDescribeNetworkConnectionsRequest()
	waitReq.NetworkConnectionName = &datasourceConnectionName
	reqErr = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeNetworkConnectionsWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.NetworkConnectionSet == nil || len(result.Response.NetworkConnectionSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc datasource house attachment failed, Response is nil."))
		}

		if result.Response.NetworkConnectionSet[0].State != nil && *result.Response.NetworkConnectionSet[0].State == 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("DLC datasource house attachment is still running..."))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc datasource house attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudDlcDatasourceHouseAttachmentRead(d, meta)
}

func resourceTencentCloudDlcDatasourceHouseAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_datasource_house_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                    = tccommon.GetLogId(tccommon.ContextNil)
		ctx                      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service                  = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		datasourceConnectionName = d.Id()
	)

	respData, err := service.DescribeDlcDatasourceHouseAttachmentById(ctx, datasourceConnectionName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dlc_datasource_house_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.DatasourceConnectionName != nil {
		_ = d.Set("datasource_connection_name", respData.DatasourceConnectionName)
	}

	if respData.HouseName != nil {
		_ = d.Set("data_engine_names", []string{*respData.HouseName})
	}

	if respData.NetworkConnectionType != nil {
		_ = d.Set("network_connection_type", respData.NetworkConnectionType)
	}

	if respData.NetworkConnectionDesc != nil {
		_ = d.Set("network_connection_desc", respData.NetworkConnectionDesc)
	}

	return nil
}

func resourceTencentCloudDlcDatasourceHouseAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_datasource_house_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                    = tccommon.GetLogId(tccommon.ContextNil)
		ctx                      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		datasourceConnectionName = d.Id()
	)

	if d.HasChange("network_connection_desc") {
		request := dlcv20210125.NewUpdateNetworkConnectionRequest()
		if v, ok := d.GetOk("network_connection_desc"); ok {
			request.NetworkConnectionDesc = helper.String(v.(string))
		}

		request.NetworkConnectionName = &datasourceConnectionName
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateNetworkConnectionWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dlc datasource house attachment failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudDlcDatasourceHouseAttachmentRead(d, meta)
}

func resourceTencentCloudDlcDatasourceHouseAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_datasource_house_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                    = tccommon.GetLogId(tccommon.ContextNil)
		ctx                      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request                  = dlcv20210125.NewUnboundDatasourceHouseRequest()
		datasourceConnectionName = d.Id()
	)

	request.NetworkConnectionName = &datasourceConnectionName
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UnboundDatasourceHouseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dlc datasource house attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
