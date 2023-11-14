/*
Provides a resource to create a wedata datasource

Example Usage

```hcl
resource "tencentcloud_wedata_datasource" "datasource" {
  name = "name"
  category = "DB"
  type = "MYSQL"
  owner_project_id = "110111121"
  owner_project_name = "ownerprojectname"
  owner_project_ident = "OwnerProjectIdent"
  biz_params = "{}"
  params = "{}"
  description = "descr"
  display = "Display"
  database_name = "db"
  instance = "instance"
  status = 1
  cluster_id = "cid"
  collect = "false"
  c_o_s_bucket = "aaaa"
  c_o_s_region = "ap-beijing"
}
```

Import

wedata datasource can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_datasource.datasource datasource_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudWedataDatasource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataDatasourceCreate,
		Read:   resourceTencentCloudWedataDatasourceRead,
		Update: resourceTencentCloudWedataDatasourceUpdate,
		Delete: resourceTencentCloudWedataDatasourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Description: "Ownerprojectname.",
			},

			"owner_project_ident": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "OwnerProjectIdent.",
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

			"c_o_s_bucket": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "COSBucket.",
			},

			"c_o_s_region": {
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

	logId := getLogId(contextNil)

	var (
		request  = wedata.NewCreateDataSourceRequest()
		response = wedata.NewCreateDataSourceResponse()
		iD       int
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

	if v, ok := d.GetOk("c_o_s_bucket"); ok {
		request.COSBucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("c_o_s_region"); ok {
		request.COSRegion = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateDataSource(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata datasource failed, reason:%+v", logId, err)
		return err
	}

	iD = *response.Response.ID
	d.SetId(helper.Int64ToStr(int64(iD)))

	return resourceTencentCloudWedataDatasourceRead(d, meta)
}

func resourceTencentCloudWedataDatasourceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_datasource.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}

	datasourceId := d.Id()

	datasource, err := service.DescribeWedataDatasourceById(ctx, iD)
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

	if datasource.Collect != nil {
		_ = d.Set("collect", datasource.Collect)
	}

	if datasource.COSBucket != nil {
		_ = d.Set("c_o_s_bucket", datasource.COSBucket)
	}

	if datasource.COSRegion != nil {
		_ = d.Set("c_o_s_region", datasource.COSRegion)
	}

	return nil
}

func resourceTencentCloudWedataDatasourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_datasource.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := wedata.NewModifyDataSourceRequest()

	datasourceId := d.Id()

	request.ID = &iD

	immutableArgs := []string{"name", "category", "type", "owner_project_id", "owner_project_name", "owner_project_ident", "biz_params", "params", "description", "display", "database_name", "instance", "status", "cluster_id", "collect", "c_o_s_bucket", "c_o_s_region"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("category") {
		if v, ok := d.GetOk("category"); ok {
			request.Category = helper.String(v.(string))
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("owner_project_id") {
		if v, ok := d.GetOk("owner_project_id"); ok {
			request.OwnerProjectId = helper.String(v.(string))
		}
	}

	if d.HasChange("owner_project_name") {
		if v, ok := d.GetOk("owner_project_name"); ok {
			request.OwnerProjectName = helper.String(v.(string))
		}
	}

	if d.HasChange("owner_project_ident") {
		if v, ok := d.GetOk("owner_project_ident"); ok {
			request.OwnerProjectIdent = helper.String(v.(string))
		}
	}

	if d.HasChange("biz_params") {
		if v, ok := d.GetOk("biz_params"); ok {
			request.BizParams = helper.String(v.(string))
		}
	}

	if d.HasChange("params") {
		if v, ok := d.GetOk("params"); ok {
			request.Params = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("display") {
		if v, ok := d.GetOk("display"); ok {
			request.Display = helper.String(v.(string))
		}
	}

	if d.HasChange("database_name") {
		if v, ok := d.GetOk("database_name"); ok {
			request.DatabaseName = helper.String(v.(string))
		}
	}

	if d.HasChange("instance") {
		if v, ok := d.GetOk("instance"); ok {
			request.Instance = helper.String(v.(string))
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("cluster_id") {
		if v, ok := d.GetOk("cluster_id"); ok {
			request.ClusterId = helper.String(v.(string))
		}
	}

	if d.HasChange("collect") {
		if v, ok := d.GetOk("collect"); ok {
			request.Collect = helper.String(v.(string))
		}
	}

	if d.HasChange("c_o_s_bucket") {
		if v, ok := d.GetOk("c_o_s_bucket"); ok {
			request.COSBucket = helper.String(v.(string))
		}
	}

	if d.HasChange("c_o_s_region") {
		if v, ok := d.GetOk("c_o_s_region"); ok {
			request.COSRegion = helper.String(v.(string))
		}
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

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	datasourceId := d.Id()

	if err := service.DeleteWedataDatasourceById(ctx, iD); err != nil {
		return err
	}

	return nil
}
