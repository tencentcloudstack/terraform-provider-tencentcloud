package wedata

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataDataSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataDataSourceCreate,
		Read:   resourceTencentCloudWedataDataSourceRead,
		Update: resourceTencentCloudWedataDataSourceUpdate,
		Delete: resourceTencentCloudWedataDataSourceDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Data source project ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Data source name.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Data source type: enumeration values.\n\n- MYSQL\n- TENCENT_MYSQL\n- POSTGRE\n- ORACLE\n- SQLSERVER\n- FTP\n- HIVE\n- HUDI\n- HDFS\n- ICEBERG\n- KAFKA\n- DTS_KAFKA\n- HBASE\n- SPARK\n- TBASE\n- DB2\n- DM\n- GAUSSDB\n- GBASE\n- IMPALA\n- ES\n- TENCENT_ES\n- GREENPLUM\n- SAP_HANA\n- SFTP\n- OCEANBASE\n- CLICKHOUSE\n- KUDU\n- VERTICA\n- REDIS\n- COS\n- DLC\n- DORIS\n- CKAFKA\n- S3_DATAINSIGHT\n- TDSQL\n- TDSQL_MYSQL\n- MONGODB\n- TENCENT_MONGODB\n- REST_API\n- TiDB\n- StarRocks\n- Trino\n- Kyuubi\n- TCHOUSE_X\n- TCHOUSE_P\n- TCHOUSE_C\n- TCHOUSE_D\n- INFLUXDB\n- BIG_QUERY\n- SSH\n- BLOB\n- TDSQL_POSTGRE\n- GDB\n- TDENGINE\n- TDSQLC.",
			},

			"prod_con_properties": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data source configuration information, stored in JSON KV format, with different KV storage information for each data source type.\n\n> deployType: \nCONNSTR_PUBLICDB(Public network instance) \nCONNSTR_CVMDB(Self-built instance)\nINSTANCE(Cloud instance)\n\n```\nmysql: Self-built instance\n{\n    \"deployType\": \"CONNSTR_CVMDB\",\n    \"url\": \"jdbc:mysql://1.1.1.1:1111/database\",\n    \"username\": \"root\",\n    \"password\": \"root\",\n    \"region\": \"ap-shanghai\",\n    \"vpcId\": \"vpc-kprq42yo\",\n    \"type\": \"MYSQL\"\n}\nmysql: Cloud instance\n{\n    \"instanceid\": \"cdb-12uxdo5e\",\n    \"db\": \"db\",\n    \"region\": \"ap-shanghai\",\n    \"username\": \"msyql\",\n    \"password\": \"mysql\",\n    \"deployType\": \"INSTANCE\",\n    \"type\": \"TENCENT_MYSQL\"\n}\nsql_server: \n{\n    \"deployType\": \"CONNSTR_PUBLICDB\",\n    \"url\": \"jdbc:sqlserver://1.1.1.1:223;DatabaseName=database\",\n    \"username\": \"user_1\",\n    \"password\": \"pass_2\",\n    \"type\": \"SQLSERVER\"\n}\nredis:\n    redisType:\n    -NO_ACCOUT(No account)\n    -SELF_ACCOUNT(Custom account)\n{\n    \"deployType\": \"CONNSTR_PUBLICDB\",\n    \"username\":\"\"\n    \"password\": \"pass\",\n    \"ip\": \"1.1.1.1\",\n    \"port\": \"6379\",\n    \"redisType\": \"NO_ACCOUT\",\n    \"type\": \"REDIS\"\n}\noracle: \n{\n    \"deployType\": \"CONNSTR_CVMDB\",\n    \"url\": \"jdbc:oracle:thin:@1.1.1.1:1521:prod\",\n    \"username\": \"oracle\",\n    \"password\": \"pass\",\n    \"region\": \"ap-shanghai\",\n    \"vpcId\": \"vpc-kprq42yo\",\n    \"type\": \"ORACLE\"\n}\nmongodb:\n    advanceParams(Custom parameters, will be appended to the URL)\n{\n    \"advanceParams\": [\n        {\n            \"key\": \"authSource\",\n            \"value\": \"auth\"\n        }\n    ],\n    \"db\": \"admin\",\n    \"deployType\": \"CONNSTR_PUBLICDB\",\n    \"username\": \"user\",\n    \"password\": \"pass\",\n    \"type\": \"MONGODB\",\n    \"host\": \"1.1.1.1:9200\"\n}\npostgresql:\n{\n    \"deployType\": \"CONNSTR_PUBLICDB\",\n    \"url\": \"jdbc:postgresql://1.1.1.1:1921/database\",\n    \"username\": \"user\",\n    \"password\": \"pass\",\n    \"type\": \"POSTGRE\"\n}\nkafka:\n    authType:\n        - sasl\n        - jaas\n        - sasl_plaintext\n        - sasl_ssl\n        - GSSAPI\n    ssl:\n        -PLAIN\n        -GSSAPI\n{\n    \"deployType\": \"CONNSTR_PUBLICDB\",\n    \"host\": \"1.1.1.1:9092\",\n    \"ssl\": \"GSSAPI\",\n    \"authType\": \"sasl\",\n    \"type\": \"KAFKA\",\n    \"principal\": \"aaaa\",\n    \"serviceName\": \"kafka\"\n}\n\ncos:\n{\n    \"region\": \"ap-shanghai\",\n    \"deployType\": \"INSTANCE\",\n    \"secretId\": \"aaaaa\",\n    \"secretKey\": \"sssssss\",\n    \"bucket\": \"aaa\",\n    \"type\": \"COS\"\n}\n\n```.",
			},

			"dev_con_properties": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Development environment data source configuration information, required if the project is in standard mode.",
			},

			"prod_file_upload": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Production environment data source file upload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trust_store": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Truststore authentication file, default filename truststore.jks.",
						},
						"key_store": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Keystore authentication file, default filename keystore.jks.",
						},
						"core_site": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "core-site.xml file.",
						},
						"hdfs_site": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "hdfs-site.xml file.",
						},
						"hive_site": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "hive-site.xml file.",
						},
						"hbase_site": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "hbase-site file.",
						},
						"key_tab": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "keytab file, default filename [data source name].keytab.",
						},
						"krb5_conf": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "krb5.conf file.",
						},
						"private_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Private key, default filename private_key.pem.",
						},
						"public_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public key, default filename public_key.pem.",
						},
					},
				},
			},

			"dev_file_upload": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Development environment data source file upload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trust_store": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Truststore authentication file, default filename truststore.jks.",
						},
						"key_store": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Keystore authentication file, default filename keystore.jks.",
						},
						"core_site": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "core-site.xml file.",
						},
						"hdfs_site": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "hdfs-site.xml file.",
						},
						"hive_site": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "hive-site.xml file.",
						},
						"hbase_site": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "hbase-site file.",
						},
						"key_tab": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "keytab file, default filename [data source name].keytab.",
						},
						"krb5_conf": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "krb5.conf file.",
						},
						"private_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Private key, default filename private_key.pem.",
						},
						"public_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public key, default filename public_key.pem.",
						},
					},
				},
			},

			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data source display name, for visual viewing.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data source description information.",
			},

			// computed
			"data_source_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Data source ID.",
			},
		},
	}
}

func resourceTencentCloudWedataDataSourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_data_source.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = wedatav20250806.NewCreateDataSourceRequest()
		response     = wedatav20250806.NewCreateDataSourceResponse()
		projectId    string
		datasourceId string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("prod_con_properties"); ok {
		request.ProdConProperties = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dev_con_properties"); ok {
		request.DevConProperties = helper.String(v.(string))
	}

	if prodFileUploadMap, ok := helper.InterfacesHeadMap(d, "prod_file_upload"); ok {
		dataSourceFileUpload := wedatav20250806.DataSourceFileUpload{}
		if v, ok := prodFileUploadMap["trust_store"].(string); ok && v != "" {
			dataSourceFileUpload.TrustStore = helper.String(v)
		}

		if v, ok := prodFileUploadMap["key_store"].(string); ok && v != "" {
			dataSourceFileUpload.KeyStore = helper.String(v)
		}

		if v, ok := prodFileUploadMap["core_site"].(string); ok && v != "" {
			dataSourceFileUpload.CoreSite = helper.String(v)
		}

		if v, ok := prodFileUploadMap["hdfs_site"].(string); ok && v != "" {
			dataSourceFileUpload.HdfsSite = helper.String(v)
		}

		if v, ok := prodFileUploadMap["hive_site"].(string); ok && v != "" {
			dataSourceFileUpload.HiveSite = helper.String(v)
		}

		if v, ok := prodFileUploadMap["hbase_site"].(string); ok && v != "" {
			dataSourceFileUpload.HBASESite = helper.String(v)
		}

		if v, ok := prodFileUploadMap["key_tab"].(string); ok && v != "" {
			dataSourceFileUpload.KeyTab = helper.String(v)
		}

		if v, ok := prodFileUploadMap["krb5_conf"].(string); ok && v != "" {
			dataSourceFileUpload.KRB5Conf = helper.String(v)
		}

		if v, ok := prodFileUploadMap["private_key"].(string); ok && v != "" {
			dataSourceFileUpload.PrivateKey = helper.String(v)
		}

		if v, ok := prodFileUploadMap["public_key"].(string); ok && v != "" {
			dataSourceFileUpload.PublicKey = helper.String(v)
		}

		request.ProdFileUpload = &dataSourceFileUpload
	}

	if devFileUploadMap, ok := helper.InterfacesHeadMap(d, "dev_file_upload"); ok {
		dataSourceFileUpload2 := wedatav20250806.DataSourceFileUpload{}
		if v, ok := devFileUploadMap["trust_store"].(string); ok && v != "" {
			dataSourceFileUpload2.TrustStore = helper.String(v)
		}

		if v, ok := devFileUploadMap["key_store"].(string); ok && v != "" {
			dataSourceFileUpload2.KeyStore = helper.String(v)
		}

		if v, ok := devFileUploadMap["core_site"].(string); ok && v != "" {
			dataSourceFileUpload2.CoreSite = helper.String(v)
		}

		if v, ok := devFileUploadMap["hdfs_site"].(string); ok && v != "" {
			dataSourceFileUpload2.HdfsSite = helper.String(v)
		}

		if v, ok := devFileUploadMap["hive_site"].(string); ok && v != "" {
			dataSourceFileUpload2.HiveSite = helper.String(v)
		}

		if v, ok := devFileUploadMap["hbase_site"].(string); ok && v != "" {
			dataSourceFileUpload2.HBASESite = helper.String(v)
		}

		if v, ok := devFileUploadMap["key_tab"].(string); ok && v != "" {
			dataSourceFileUpload2.KeyTab = helper.String(v)
		}

		if v, ok := devFileUploadMap["krb5_conf"].(string); ok && v != "" {
			dataSourceFileUpload2.KRB5Conf = helper.String(v)
		}

		if v, ok := devFileUploadMap["private_key"].(string); ok && v != "" {
			dataSourceFileUpload2.PrivateKey = helper.String(v)
		}

		if v, ok := devFileUploadMap["public_key"].(string); ok && v != "" {
			dataSourceFileUpload2.PublicKey = helper.String(v)
		}

		request.DevFileUpload = &dataSourceFileUpload2
	}

	if v, ok := d.GetOk("display_name"); ok {
		request.DisplayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateDataSourceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata data source failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata data source failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.Status == nil || !*response.Response.Data.Status {
		return fmt.Errorf("Create wedata data source failed, Status is false")
	}

	if response.Response.Data.DataSourceId == nil {
		return fmt.Errorf("DataSourceId is nil.")
	}

	datasourceIdInt64 := *response.Response.Data.DataSourceId
	datasourceId = strconv.FormatInt(datasourceIdInt64, 10)
	d.SetId(strings.Join([]string{projectId, datasourceId}, tccommon.FILED_SP))
	return resourceTencentCloudWedataDataSourceRead(d, meta)
}

func resourceTencentCloudWedataDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_data_source.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	datasourceId := idSplit[1]

	respData, err := service.DescribeWedataDataSourceById(ctx, projectId, datasourceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_data_source` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ProjectId != nil {
		_ = d.Set("project_id", respData.ProjectId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.ProdConProperties != nil {
		_ = d.Set("prod_con_properties", respData.ProdConProperties)
	}

	if respData.DevConProperties != nil {
		_ = d.Set("dev_con_properties", respData.DevConProperties)
	}

	if respData.DisplayName != nil {
		_ = d.Set("display_name", respData.DisplayName)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Id != nil {
		_ = d.Set("data_source_id", respData.Id)
	}

	return nil
}

func resourceTencentCloudWedataDataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_data_source.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	datasourceId := idSplit[1]
	datasourceIdUint64 := helper.StrToUint64Point(datasourceId)

	needChange := false
	mutableArgs := []string{"prod_con_properties", "dev_con_properties", "prod_file_upload", "dev_file_upload", "display_name", "description"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wedatav20250806.NewUpdateDataSourceRequest()
		response := wedatav20250806.NewUpdateDataSourceResponse()
		if v, ok := d.GetOk("prod_con_properties"); ok {
			request.ProdConProperties = helper.String(v.(string))
		}

		if v, ok := d.GetOk("dev_con_properties"); ok {
			request.DevConProperties = helper.String(v.(string))
		}

		if prodFileUploadMap, ok := helper.InterfacesHeadMap(d, "prod_file_upload"); ok {
			dataSourceFileUpload := wedatav20250806.DataSourceFileUpload{}
			if v, ok := prodFileUploadMap["trust_store"].(string); ok && v != "" {
				dataSourceFileUpload.TrustStore = helper.String(v)
			}

			if v, ok := prodFileUploadMap["key_store"].(string); ok && v != "" {
				dataSourceFileUpload.KeyStore = helper.String(v)
			}

			if v, ok := prodFileUploadMap["core_site"].(string); ok && v != "" {
				dataSourceFileUpload.CoreSite = helper.String(v)
			}

			if v, ok := prodFileUploadMap["hdfs_site"].(string); ok && v != "" {
				dataSourceFileUpload.HdfsSite = helper.String(v)
			}

			if v, ok := prodFileUploadMap["hive_site"].(string); ok && v != "" {
				dataSourceFileUpload.HiveSite = helper.String(v)
			}

			if v, ok := prodFileUploadMap["hbase_site"].(string); ok && v != "" {
				dataSourceFileUpload.HBASESite = helper.String(v)
			}

			if v, ok := prodFileUploadMap["key_tab"].(string); ok && v != "" {
				dataSourceFileUpload.KeyTab = helper.String(v)
			}

			if v, ok := prodFileUploadMap["krb5_conf"].(string); ok && v != "" {
				dataSourceFileUpload.KRB5Conf = helper.String(v)
			}

			if v, ok := prodFileUploadMap["private_key"].(string); ok && v != "" {
				dataSourceFileUpload.PrivateKey = helper.String(v)
			}

			if v, ok := prodFileUploadMap["public_key"].(string); ok && v != "" {
				dataSourceFileUpload.PublicKey = helper.String(v)
			}

			request.ProdFileUpload = &dataSourceFileUpload
		}

		if devFileUploadMap, ok := helper.InterfacesHeadMap(d, "dev_file_upload"); ok {
			dataSourceFileUpload2 := wedatav20250806.DataSourceFileUpload{}
			if v, ok := devFileUploadMap["trust_store"].(string); ok && v != "" {
				dataSourceFileUpload2.TrustStore = helper.String(v)
			}

			if v, ok := devFileUploadMap["key_store"].(string); ok && v != "" {
				dataSourceFileUpload2.KeyStore = helper.String(v)
			}

			if v, ok := devFileUploadMap["core_site"].(string); ok && v != "" {
				dataSourceFileUpload2.CoreSite = helper.String(v)
			}

			if v, ok := devFileUploadMap["hdfs_site"].(string); ok && v != "" {
				dataSourceFileUpload2.HdfsSite = helper.String(v)
			}

			if v, ok := devFileUploadMap["hive_site"].(string); ok && v != "" {
				dataSourceFileUpload2.HiveSite = helper.String(v)
			}

			if v, ok := devFileUploadMap["hbase_site"].(string); ok && v != "" {
				dataSourceFileUpload2.HBASESite = helper.String(v)
			}

			if v, ok := devFileUploadMap["key_tab"].(string); ok && v != "" {
				dataSourceFileUpload2.KeyTab = helper.String(v)
			}

			if v, ok := devFileUploadMap["krb5_conf"].(string); ok && v != "" {
				dataSourceFileUpload2.KRB5Conf = helper.String(v)
			}

			if v, ok := devFileUploadMap["private_key"].(string); ok && v != "" {
				dataSourceFileUpload2.PrivateKey = helper.String(v)
			}

			if v, ok := devFileUploadMap["public_key"].(string); ok && v != "" {
				dataSourceFileUpload2.PublicKey = helper.String(v)
			}

			request.DevFileUpload = &dataSourceFileUpload2
		}

		if v, ok := d.GetOk("display_name"); ok {
			request.DisplayName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		request.ProjectId = &projectId
		request.Id = datasourceIdUint64
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateDataSourceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("Update data source failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update wedata data source failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if !*response.Response.Data.Status {
			return fmt.Errorf("Update data source %s failed, Status is false.", datasourceId)
		}
	}

	return resourceTencentCloudWedataDataSourceRead(d, meta)
}

func resourceTencentCloudWedataDataSourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_data_source.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = wedatav20250806.NewDeleteDataSourceRequest()
		response = wedatav20250806.NewDeleteDataSourceResponse()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	datasourceId := idSplit[1]
	datasourceIdUint64 := helper.StrToUint64Point(datasourceId)

	request.ProjectId = &projectId
	request.Id = datasourceIdUint64
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteDataSourceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete data source failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata data source failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if *response.Response.Data.Status {
		return nil
	}

	return fmt.Errorf("Delete data source %s failed, Status is false.", datasourceId)
}
