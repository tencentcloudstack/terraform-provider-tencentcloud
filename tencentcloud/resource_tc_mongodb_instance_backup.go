/*
Provides a resource to create a mongodb instance_backup

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_backup" "instance_backup" {
  instance_id = "cmgo-9d0p6umb"
  backup_method = 0
  backup_remark = "my backup"
}
```

*/
package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMongodbInstanceBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceBackupCreate,
		Read:   resourceTencentCloudMongodbInstanceBackupRead,
		Delete: resourceTencentCloudMongodbInstanceBackupDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.",
			},

			"backup_method": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "0:logical backup, 1:physical backup.",
			},

			"backup_remark": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "backup notes.",
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mongodb_instance_backup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mongodb.NewCreateBackupDBInstanceRequest()
		response = mongodb.NewCreateBackupDBInstanceResponse()
		taskId   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("backup_method"); v != nil {
		request.BackupMethod = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("backup_remark"); ok {
		request.BackupRemark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMongodbClient().CreateBackupDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mongodb instanceBackup failed, reason:%+v", logId, err)
		return nil
	}

	taskId = *response.Response.AsyncRequestId
	d.SetId(taskId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}

	timeout := d.Timeout(schema.TimeoutCreate)
	if response != nil && response.Response != nil {
		if err = service.DescribeAsyncRequestInfo(ctx, taskId, timeout); err != nil {
			return err
		}
	}

	return resourceTencentCloudMongodbInstanceBackupRead(d, meta)
}

func resourceTencentCloudMongodbInstanceBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_backup.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMongodbInstanceBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_backup.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
