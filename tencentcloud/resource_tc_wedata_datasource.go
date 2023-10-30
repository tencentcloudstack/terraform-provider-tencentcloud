/*
Provides a resource to create a wedata datasource

Example Usage

```hcl
resource "tencentcloud_wedata_datasource" "example" {
  name                = "tf_example"
  category            = "DB"
  type                = "MYSQL"
  owner_project_id    = "110111121"
  owner_project_name  = "ownerprojectname"
  owner_project_ident = "OwnerProjectIdent"
  biz_params          = "{}"
  params              = "{}"
  description         = "descr"
  display             = "Display"
  database_name       = "db"
  instance            = "instance"
  status              = 1
  cluster_id          = "cid"
  collect             = "false"
  cos_bucket          = "aaaa"
  cos_region          = "ap-guangzhou"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWedataDatasource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataDatasourceCreate,
		Read:   resourceTencentCloudWedataDatasourceRead,
		Update: resourceTencentCloudWedataDatasourceUpdate,
		Delete: resourceTencentCloudWedataDatasourceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "DataSource Name.",
			},
			"category": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "DataSource Category.",
			},
			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "DataSource Type.",
			},
			"owner_project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Owner projectId.",
			},
			"owner_project_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Owner project name.",
			},
			"owner_project_ident": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Owner Project Ident.",
			},
			"biz_params": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "BizParams.",
			},
			"params": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Params.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description.",
			},
			"display": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Display.",
			},
			"database_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Dbname.",
			},
			"instance": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance.",
			},
			"status": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Status.",
			},
			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ClusterId.",
			},
			"collect": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Collect.",
			},
			"cos_bucket": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "COSBucket.",
			},
			"cos_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cos region.",
			},
		},
	}
}

func resourceTencentCloudWedataDatasourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_datasource.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		request  = wedata.NewCreateDataSourceRequest()
		response = wedata.NewCreateDataSourceResponse()
		Id       string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("category"); ok {
		request.Category = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner_project_id"); ok {
		request.OwnerProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner_project_name"); ok {
		request.OwnerProjectName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner_project_ident"); ok {
		request.OwnerProjectIdent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("biz_params"); ok {
		request.BizParams = helper.String(v.(string))
	}

	if v, ok := d.GetOk("params"); ok {
		request.Params = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("display"); ok {
		request.Display = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database_name"); ok {
		request.DatabaseName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance"); ok {
		request.Instance = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("collect"); ok {
		request.Collect = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_bucket"); ok {
		request.COSBucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_region"); ok {
		request.COSRegion = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateDataSource(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf(" wedata datasource not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata datasource failed, reason:%+v", logId, err)
		return err
	}

	DataInt := *response.Response.Data
	Id = strconv.FormatUint(DataInt, 10)
	d.SetId(Id)

	return resourceTencentCloudWedataDatasourceRead(d, meta)
}

func resourceTencentCloudWedataDatasourceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_datasource.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId        = getLogId(contextNil)
		ctx          = context.WithValue(context.TODO(), logIdKey, logId)
		service      = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
		datasourceId = d.Id()
	)

	datasource, err := service.DescribeWedataDatasourceById(ctx, datasourceId)
	if err != nil {
		return err
	}

	if datasource == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataDatasource` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if datasource.Name != nil {
		_ = d.Set("name", datasource.Name)
	}

	if datasource.Category != nil {
		_ = d.Set("category", datasource.Category)
	}

	if datasource.Type != nil {
		_ = d.Set("type", datasource.Type)
	}

	if datasource.OwnerProjectId != nil {
		_ = d.Set("owner_project_id", datasource.OwnerProjectId)
	}

	if datasource.OwnerProjectName != nil {
		_ = d.Set("owner_project_name", datasource.OwnerProjectName)
	}

	if datasource.OwnerProjectIdent != nil {
		_ = d.Set("owner_project_ident", datasource.OwnerProjectIdent)
	}

	if datasource.BizParams != nil {
		_ = d.Set("biz_params", datasource.BizParams)
	}

	if datasource.Params != nil {
		_ = d.Set("params", datasource.Params)
	}

	if datasource.Description != nil {
		_ = d.Set("description", datasource.Description)
	}

	if datasource.Display != nil {
		_ = d.Set("display", datasource.Display)
	}

	if datasource.DatabaseName != nil {
		_ = d.Set("database_name", datasource.DatabaseName)
	}

	if datasource.Instance != nil {
		_ = d.Set("instance", datasource.Instance)
	}

	if datasource.Status != nil {
		_ = d.Set("status", datasource.Status)
	}

	if datasource.ClusterId != nil {
		_ = d.Set("cluster_id", datasource.ClusterId)
	}

	return nil
}

func resourceTencentCloudWedataDatasourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_datasource.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId        = getLogId(contextNil)
		request      = wedata.NewModifyDataSourceRequest()
		datasourceId = d.Id()
	)

	Id, _ := strconv.ParseUint(datasourceId, 10, 64)
	request.ID = &Id

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("category"); ok {
		request.Category = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner_project_id"); ok {
		request.OwnerProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner_project_name"); ok {
		request.OwnerProjectName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner_project_ident"); ok {
		request.OwnerProjectIdent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("biz_params"); ok {
		request.BizParams = helper.String(v.(string))
	}

	if v, ok := d.GetOk("params"); ok {
		request.Params = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("display"); ok {
		request.Display = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database_name"); ok {
		request.DatabaseName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance"); ok {
		request.Instance = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("collect"); ok {
		request.Collect = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_bucket"); ok {
		request.COSBucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_region"); ok {
		request.COSRegion = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().ModifyDataSource(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update wedata datasource failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataDatasourceRead(d, meta)
}

func resourceTencentCloudWedataDatasourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_datasource.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId        = getLogId(contextNil)
		ctx          = context.WithValue(context.TODO(), logIdKey, logId)
		service      = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
		datasourceId = d.Id()
	)

	if err := service.DeleteWedataDatasourceById(ctx, datasourceId); err != nil {
		return err
	}

	return nil
}
