/*
Provides a resource to create a bi datasource

Example Usage

```hcl
resource "tencentcloud_bi_datasource" "datasource" {
  db_host = "1.2.3.4"
  db_port = 3306
  service_type = "Own"
  db_type = "Database type."
  charset = "utf8"
  db_user = "root"
  db_pwd = "abc"
  db_name = "abc"
  source_name = "abc"
  project_id = 123
  catalog = "presto"
  data_origin = "abc"
  data_origin_project_id = "abc"
  data_origin_datasource_id = "abc"
  extra_param = ""
  uniq_vpc_id = ""
  vip = ""
  vport = ""
  vpc_id = ""
}
```

Import

bi datasource can be imported using the id, e.g.

```
terraform import tencentcloud_bi_datasource.datasource datasource_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
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

			"service_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Own or Cloud.",
			},

			"db_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "MYSQL.",
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

			"extra_param": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Extended parameters.",
			},

			"uniq_vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Tencent cloud private network unified identity.",
			},

			"vip": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Private network ip.",
			},

			"vport": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Private network port.",
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
		request  = bi.NewCreateDatasourceRequest()
		response = bi.NewCreateDatasourceResponse()
		id       int
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

	if v, ok := d.GetOk("extra_param"); ok {
		request.ExtraParam = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
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

	id = *response.Response.Id
	d.SetId(helper.Int64ToStr(int64(id)))

	return resourceTencentCloudBiDatasourceRead(d, meta)
}

func resourceTencentCloudBiDatasourceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_datasource.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}

	datasourceId := d.Id()

	datasource, err := service.DescribeBiDatasourceById(ctx, id)
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

	if datasource.DbPwd != nil {
		_ = d.Set("db_pwd", datasource.DbPwd)
	}

	if datasource.DbName != nil {
		_ = d.Set("db_name", datasource.DbName)
	}

	if datasource.SourceName != nil {
		_ = d.Set("source_name", datasource.SourceName)
	}

	if datasource.ProjectId != nil {
		_ = d.Set("project_id", datasource.ProjectId)
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

	if datasource.ExtraParam != nil {
		_ = d.Set("extra_param", datasource.ExtraParam)
	}

	if datasource.UniqVpcId != nil {
		_ = d.Set("uniq_vpc_id", datasource.UniqVpcId)
	}

	if datasource.Vip != nil {
		_ = d.Set("vip", datasource.Vip)
	}

	if datasource.Vport != nil {
		_ = d.Set("vport", datasource.Vport)
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

	datasourceId := d.Id()

	request.Id = &id

	immutableArgs := []string{"db_host", "db_port", "service_type", "db_type", "charset", "db_user", "db_pwd", "db_name", "source_name", "project_id", "catalog", "data_origin", "data_origin_project_id", "data_origin_datasource_id", "extra_param", "uniq_vpc_id", "vip", "vport", "vpc_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("db_host") {
		if v, ok := d.GetOk("db_host"); ok {
			request.DbHost = helper.String(v.(string))
		}
	}

	if d.HasChange("db_port") {
		if v, ok := d.GetOkExists("db_port"); ok {
			request.DbPort = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("service_type") {
		if v, ok := d.GetOk("service_type"); ok {
			request.ServiceType = helper.String(v.(string))
		}
	}

	if d.HasChange("db_type") {
		if v, ok := d.GetOk("db_type"); ok {
			request.DbType = helper.String(v.(string))
		}
	}

	if d.HasChange("charset") {
		if v, ok := d.GetOk("charset"); ok {
			request.Charset = helper.String(v.(string))
		}
	}

	if d.HasChange("db_user") {
		if v, ok := d.GetOk("db_user"); ok {
			request.DbUser = helper.String(v.(string))
		}
	}

	if d.HasChange("db_pwd") {
		if v, ok := d.GetOk("db_pwd"); ok {
			request.DbPwd = helper.String(v.(string))
		}
	}

	if d.HasChange("db_name") {
		if v, ok := d.GetOk("db_name"); ok {
			request.DbName = helper.String(v.(string))
		}
	}

	if d.HasChange("source_name") {
		if v, ok := d.GetOk("source_name"); ok {
			request.SourceName = helper.String(v.(string))
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOkExists("project_id"); ok {
			request.ProjectId = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("catalog") {
		if v, ok := d.GetOk("catalog"); ok {
			request.Catalog = helper.String(v.(string))
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

	if d.HasChange("extra_param") {
		if v, ok := d.GetOk("extra_param"); ok {
			request.ExtraParam = helper.String(v.(string))
		}
	}

	if d.HasChange("uniq_vpc_id") {
		if v, ok := d.GetOk("uniq_vpc_id"); ok {
			request.UniqVpcId = helper.String(v.(string))
		}
	}

	if d.HasChange("vip") {
		if v, ok := d.GetOk("vip"); ok {
			request.Vip = helper.String(v.(string))
		}
	}

	if d.HasChange("vport") {
		if v, ok := d.GetOk("vport"); ok {
			request.Vport = helper.String(v.(string))
		}
	}

	if d.HasChange("vpc_id") {
		if v, ok := d.GetOk("vpc_id"); ok {
			request.VpcId = helper.String(v.(string))
		}
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
	datasourceId := d.Id()

	if err := service.DeleteBiDatasourceById(ctx, id); err != nil {
		return err
	}

	return nil
}
