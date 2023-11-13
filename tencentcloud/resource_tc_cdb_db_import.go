/*
Provides a resource to create a cdb db_import

Example Usage

```hcl
resource "tencentcloud_cdb_db_import" "db_import" {
  instance_id = ""
  user = ""
  file_name = ""
  password = ""
  db_name = ""
  cos_url = ""
}
```

Import

cdb db_import can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_db_import.db_import db_import_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCdbDbImport() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbDbImportCreate,
		Read:   resourceTencentCloudCdbDbImportRead,
		Delete: resourceTencentCloudCdbDbImportDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.",
			},

			"user": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The username of the cloud database.",
			},

			"file_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "File name. This file refers to the file that the user has uploaded to Tencent Cloud, and only .sql files are supported.",
			},

			"password": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The password of the user account of the cloud database instance.",
			},

			"db_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of the imported target database, if it is not passed, it means that no database is specified.",
			},

			"cos_url": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of the imported target database, if it is not passed, it means that no database is specified.",
			},
		},
	}
}

func resourceTencentCloudCdbDbImportCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_db_import.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewCreateDBImportJobRequest()
		response   = cdb.NewCreateDBImportJobResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user"); ok {
		request.User = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_name"); ok {
		request.FileName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		request.DbName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_url"); ok {
		request.CosUrl = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().CreateDBImportJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdb dbImport failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbDbImportStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbDbImportRead(d, meta)
}

func resourceTencentCloudCdbDbImportRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_db_import.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	dbImportId := d.Id()

	dbImport, err := service.DescribeCdbDbImportById(ctx, instanceId)
	if err != nil {
		return err
	}

	if dbImport == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbDbImport` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dbImport.InstanceId != nil {
		_ = d.Set("instance_id", dbImport.InstanceId)
	}

	if dbImport.User != nil {
		_ = d.Set("user", dbImport.User)
	}

	if dbImport.FileName != nil {
		_ = d.Set("file_name", dbImport.FileName)
	}

	if dbImport.Password != nil {
		_ = d.Set("password", dbImport.Password)
	}

	if dbImport.DbName != nil {
		_ = d.Set("db_name", dbImport.DbName)
	}

	if dbImport.CosUrl != nil {
		_ = d.Set("cos_url", dbImport.CosUrl)
	}

	return nil
}

func resourceTencentCloudCdbDbImportDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_db_import.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	dbImportId := d.Id()

	if err := service.DeleteCdbDbImportById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
