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

func resourceTencentCloudBiDatasourceCloud() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBiDatasourceCloudCreate,
		Read:   resourceTencentCloudBiDatasourceCloudRead,
		Update: resourceTencentCloudBiDatasourceCloudUpdate,
		Delete: resourceTencentCloudBiDatasourceCloudDelete,

		Schema: map[string]*schema.Schema{
			"service_type": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Service type, Own or Cloud.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Service type, Cloud.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance Id.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region.",
						},
					},
				},
			},

			"db_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`MYSQL`, `TDSQL-C_MYSQL`, `TDSQL_MYSQL`, `MSSQL`, `POSTGRESQL`, `MARIADB`.",
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
				Type:        schema.TypeString,
				Description: "Project id.",
			},

			"vip": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Public cloud intranet ip.",
			},

			"vport": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Public cloud intranet port.",
			},

			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Vpc identification.",
			},

			"uniq_vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Unified vpc identification.",
			},

			"region_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Region identifier.",
			},

			"extra_param": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Extended parameters.",
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
		},
	}
}

func resourceTencentCloudBiDatasourceCloudCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_datasource_cloud.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = bi.NewCreateDatasourceCloudRequest()
		response  = bi.NewCreateDatasourceCloudResponse()
		projectId string
		id        int64
	)

	if dMap, ok := helper.InterfacesHeadMap(d, "service_type"); ok {
		v, o := helper.MapToString(dMap)
		if !o {
			return fmt.Errorf("ServiceType `%s` format error", dMap)
		}
		request.ServiceType = &v
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

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vport"); ok {
		request.Vport = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("region_id"); ok {
		request.RegionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("extra_param"); ok {
		request.ExtraParam = helper.String(v.(string))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().CreateDatasourceCloud(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create bi datasourceCloud failed, reason:%+v", logId, err)
		return err
	}

	id = *response.Response.Data.Id
	d.SetId(strings.Join([]string{projectId, strconv.FormatInt(id, 10)}, FILED_SP))

	return resourceTencentCloudBiDatasourceCloudRead(d, meta)
}

func resourceTencentCloudBiDatasourceCloudRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_datasource_cloud.read")()
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

	datasourceCloud, err := service.DescribeBiDatasourceCloudById(ctx, uint64(projectIdInt), uint64(idInt))
	if err != nil {
		return err
	}

	if datasourceCloud == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `BiDatasourceCloud` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if datasourceCloud.ServiceType != nil {
		v, err := helper.JsonToMap(*datasourceCloud.ServiceType)
		if err != nil {
			return fmt.Errorf("ServiceType `%v` format error", *datasourceCloud.ServiceType)
		}

		_ = d.Set("service_type", []interface{}{v})
	}

	if datasourceCloud.DbType != nil {
		_ = d.Set("db_type", datasourceCloud.DbType)
	}

	if datasourceCloud.Charset != nil {
		_ = d.Set("charset", datasourceCloud.Charset)
	}

	if datasourceCloud.DbUser != nil {
		_ = d.Set("db_user", datasourceCloud.DbUser)
	}

	if datasourceCloud.DbName != nil {
		_ = d.Set("db_name", datasourceCloud.DbName)
	}

	if datasourceCloud.SourceName != nil {
		_ = d.Set("source_name", datasourceCloud.SourceName)
	}

	if datasourceCloud.ProjectId != nil {
		_ = d.Set("project_id", datasourceCloud.ProjectId)
	}

	// if datasourceCloud.Vip != nil {
	// 	_ = d.Set("vip", datasourceCloud.Vip)
	// }

	// if datasourceCloud.Vport != nil {
	// 	_ = d.Set("vport", datasourceCloud.Vport)
	// }

	if datasourceCloud.VpcId != nil {
		_ = d.Set("vpc_id", datasourceCloud.VpcId)
	}

	if datasourceCloud.UniqVpcId != nil {
		_ = d.Set("uniq_vpc_id", datasourceCloud.UniqVpcId)
	}

	if datasourceCloud.RegionId != nil {
		_ = d.Set("region_id", datasourceCloud.RegionId)
	}

	if datasourceCloud.ExtraParam != nil {
		_ = d.Set("extra_param", datasourceCloud.ExtraParam)
	}

	if datasourceCloud.DataOrigin != nil {
		_ = d.Set("data_origin", datasourceCloud.DataOrigin)
	}

	if datasourceCloud.DataOriginProjectId != nil {
		_ = d.Set("data_origin_project_id", datasourceCloud.DataOriginProjectId)
	}

	if datasourceCloud.DataOriginDatasourceId != nil {
		_ = d.Set("data_origin_datasource_id", datasourceCloud.DataOriginDatasourceId)
	}

	return nil
}

func resourceTencentCloudBiDatasourceCloudUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_datasource_cloud.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := bi.NewModifyDatasourceCloudRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	id := idSplit[1]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	idUint64 := uint64(idInt)

	request.Id = &idUint64
	request.ProjectId = &projectId

	if dMap, ok := helper.InterfacesHeadMap(d, "service_type"); ok {
		v, o := helper.MapToString(dMap)
		if !o {
			return fmt.Errorf("ServiceType `%s` format error", dMap)
		}
		request.ServiceType = &v
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

	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vport"); ok {
		request.Vport = helper.String(v.(string))
	}

	if v, ok := d.GetOk("region_id"); ok {
		request.RegionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if d.HasChange("uniq_vpc_id") {
		if v, ok := d.GetOk("uniq_vpc_id"); ok {
			request.UniqVpcId = helper.String(v.(string))
		}
	}

	if d.HasChange("extra_param") {
		if v, ok := d.GetOk("extra_param"); ok {
			request.ExtraParam = helper.String(v.(string))
		}
	}

	if d.HasChange("data_origin") {
		if v, ok := d.GetOk("data_origin"); ok {
			request.DataOrigin = helper.String(v.(string))
		}
	}

	if d.HasChange("data_origin_project_id") {
		if v, ok := d.GetOk("data_origin_project_id"); ok {
			request.DataOriginProjectId = helper.String(v.(string))
		}
	}

	if d.HasChange("data_origin_datasource_id") {
		if v, ok := d.GetOk("data_origin_datasource_id"); ok {
			request.DataOriginDatasourceId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().ModifyDatasourceCloud(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update bi datasourceCloud failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudBiDatasourceCloudRead(d, meta)
}

func resourceTencentCloudBiDatasourceCloudDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_datasource_cloud.delete")()
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

	if err := service.DeleteBiDatasourceCloudById(ctx, uint64(projectIdInt), uint64(idInt)); err != nil {
		return err
	}

	return nil
}
