/*
Provides a resource to create a sqlserver business_intelligence_file

Example Usage

```hcl
resource "tencentcloud_sqlserver_business_intelligence_file" "business_intelligence_file" {
  instance_id = "mssql-zjaha891"
  file_u_r_l = ""
  file_type = ""
  remark = ""
}
```

Import

sqlserver business_intelligence_file can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_business_intelligence_file.business_intelligence_file business_intelligence_file_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSqlserverBusinessIntelligenceFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverBusinessIntelligenceFileCreate,
		Read:   resourceTencentCloudSqlserverBusinessIntelligenceFileRead,
		Update: resourceTencentCloudSqlserverBusinessIntelligenceFileUpdate,
		Delete: resourceTencentCloudSqlserverBusinessIntelligenceFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"file_u_r_l": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "COS_URL.",
			},

			"file_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File type. Valid values: FLAT (flat file as data source), SSIS (.ispac SSIS package file).",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remarks.",
			},
		},
	}
}

func resourceTencentCloudSqlserverBusinessIntelligenceFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_file.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sqlserver.NewCreateBusinessIntelligenceFileRequest()
		response   = sqlserver.NewCreateBusinessIntelligenceFileResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_u_r_l"); ok {
		request.FileURL = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_type"); ok {
		request.FileType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CreateBusinessIntelligenceFile(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver businessIntelligenceFile failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudSqlserverBusinessIntelligenceFileRead(d, meta)
}

func resourceTencentCloudSqlserverBusinessIntelligenceFileRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_file.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	businessIntelligenceFileId := d.Id()

	businessIntelligenceFile, err := service.DescribeSqlserverBusinessIntelligenceFileById(ctx, instanceId)
	if err != nil {
		return err
	}

	if businessIntelligenceFile == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverBusinessIntelligenceFile` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if businessIntelligenceFile.InstanceId != nil {
		_ = d.Set("instance_id", businessIntelligenceFile.InstanceId)
	}

	if businessIntelligenceFile.FileURL != nil {
		_ = d.Set("file_u_r_l", businessIntelligenceFile.FileURL)
	}

	if businessIntelligenceFile.FileType != nil {
		_ = d.Set("file_type", businessIntelligenceFile.FileType)
	}

	if businessIntelligenceFile.Remark != nil {
		_ = d.Set("remark", businessIntelligenceFile.Remark)
	}

	return nil
}

func resourceTencentCloudSqlserverBusinessIntelligenceFileUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_file.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"instance_id", "file_u_r_l", "file_type", "remark"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudSqlserverBusinessIntelligenceFileRead(d, meta)
}

func resourceTencentCloudSqlserverBusinessIntelligenceFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_file.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	businessIntelligenceFileId := d.Id()

	if err := service.DeleteSqlserverBusinessIntelligenceFileById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
