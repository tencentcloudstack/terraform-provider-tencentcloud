package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudBiDatasource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBiDatasourceCreate,
		Read:   resourceTencentCloudBiDatasourceRead,
		Update: resourceTencentCloudBiDatasourceUpdate,
		Delete: resourceTencentCloudBiDatasourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_host": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Host.",
			},

			"db_port": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Port.",
			},

			"db_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`MYSQL`, `MSSQL`, `POSTGRE`, `ORACLE`, `CLICKHOUSE`, `TIDB`, `HIVE`, `PRESTO`.",
			},

			"charset": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Charset.",
			},

			"db_user": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User name.",
			},

			"db_pwd": {
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "Password.",
			},

			"db_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},

			"source_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Datasource name in BI.",
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project id.",
			},

			"service_type": {
				Optional:    true,
				Default:     "{\"Type\":\"Own\"}",
				Type:        schema.TypeString,
				Description: "Own or Cloud, default: `Own`.",
			},

			"catalog": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Catalog.",
			},

			"data_origin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Third-party datasource identification, this parameter can be ignored.",
			},

			"data_origin_project_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Third-party datasource project id, this parameter can be ignored.",
			},

			"data_origin_datasource_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Third-party datasource project id, this parameter can be ignored.",
			},

			"uniq_vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Tencent cloud private network unified identity.",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Tencent cloud private network identity.",
			},
		},
	}
}

func resourceTencentCloudBiDatasourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_datasource.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = bi.NewCreateDatasourceRequest()
		response  = bi.NewCreateDatasourceResponse()
		projectId int
		id        int64
	)
	if v, ok := d.GetOk("db_host"); ok {
		request.DbHost = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("db_port"); ok {
		request.DbPort = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("service_type"); ok {
		request.ServiceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_type"); ok {
		request.DbType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("charset"); ok {
		request.Charset = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_user"); ok {
		request.DbUser = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_pwd"); ok {
		request.DbPwd = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		request.DbName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_name"); ok {
		request.SourceName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		projectId = v.(int)
		request.ProjectId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("catalog"); ok {
		request.Catalog = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_origin"); ok {
		request.DataOrigin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_origin_project_id"); ok {
		request.DataOriginProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_origin_datasource_id"); ok {
		request.DataOriginDatasourceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().CreateDatasource(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create bi datasource failed, reason:%+v", logId, err)
		return err
	}

	id = *response.Response.Data.Id
	d.SetId(strings.Join([]string{helper.Int64ToStr(int64(projectId)), helper.Int64ToStr(id)}, FILED_SP))

	return resourceTencentCloudBiDatasourceRead(d, meta)
}

func resourceTencentCloudBiDatasourceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_datasource.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	projectIdInt, _ := strconv.ParseInt(projectId, 10, 64)
	id := idSplit[1]
	idInt, _ := strconv.ParseInt(id, 10, 64)

	datasource, err := service.DescribeBiDatasourceById(ctx, uint64(projectIdInt), uint64(idInt))
	if err != nil {
		return err
	}

	if datasource == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `BiDatasource` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if datasource.DbHost != nil {
		_ = d.Set("db_host", datasource.DbHost)
	}

	if datasource.DbPort != nil {
		_ = d.Set("db_port", datasource.DbPort)
	}

	if datasource.ServiceType != nil {
		_ = d.Set("service_type", datasource.ServiceType)
	}

	if datasource.DbType != nil {
		_ = d.Set("db_type", datasource.DbType)
	}

	if datasource.Charset != nil {
		_ = d.Set("charset", datasource.Charset)
	}

	if datasource.DbUser != nil {
		_ = d.Set("db_user", datasource.DbUser)
	}

	// if datasource.DbPwd != nil {
	// 	_ = d.Set("db_pwd", datasource.DbPwd)
	// }

	if datasource.DbName != nil {
		_ = d.Set("db_name", datasource.DbName)
	}

	if datasource.SourceName != nil {
		_ = d.Set("source_name", datasource.SourceName)
	}

	if datasource.ProjectId != nil {
		projectIdInt, err := strconv.ParseInt(*datasource.ProjectId, 10, 64)
		if err != nil {
			return err
		}
		_ = d.Set("project_id", projectIdInt)
	}

	if datasource.Catalog != nil {
		_ = d.Set("catalog", datasource.Catalog)
	}

	if datasource.DataOrigin != nil {
		_ = d.Set("data_origin", datasource.DataOrigin)
	}

	if datasource.DataOriginProjectId != nil {
		_ = d.Set("data_origin_project_id", datasource.DataOriginProjectId)
	}

	if datasource.DataOriginDatasourceId != nil {
		_ = d.Set("data_origin_datasource_id", datasource.DataOriginDatasourceId)
	}

	if datasource.UniqVpcId != nil {
		_ = d.Set("uniq_vpc_id", datasource.UniqVpcId)
	}

	if datasource.VpcId != nil {
		_ = d.Set("vpc_id", datasource.VpcId)
	}

	return nil
}

func resourceTencentCloudBiDatasourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_datasource.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := bi.NewModifyDatasourceRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	projectIdInt, _ := strconv.ParseInt(projectId, 10, 64)
	id := idSplit[1]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	idUint := uint64(projectIdInt)

	request.ProjectId = &idUint
	request.Id = &idInt

	if v, ok := d.GetOk("db_host"); ok {
		request.DbHost = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("db_port"); ok {
		request.DbPort = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("service_type"); ok {
		request.ServiceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_type"); ok {
		request.DbType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("charset"); ok {
		request.Charset = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_user"); ok {
		request.DbUser = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_pwd"); ok {
		request.DbPwd = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		request.DbName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_name"); ok {
		request.SourceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("catalog"); ok {
		request.Catalog = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_origin"); ok {
		request.DataOrigin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_origin_project_id"); ok {
		request.DataOriginProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_origin_datasource_id"); ok {
		request.DataOriginDatasourceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().ModifyDatasource(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update bi datasource failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudBiDatasourceRead(d, meta)
}

func resourceTencentCloudBiDatasourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_datasource.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	projectIdInt, _ := strconv.ParseInt(projectId, 10, 64)
	id := idSplit[1]
	idInt, _ := strconv.ParseInt(id, 10, 64)

	if err := service.DeleteBiDatasourceById(ctx, uint64(projectIdInt), uint64(idInt)); err != nil {
		return err
	}

	return nil
}
