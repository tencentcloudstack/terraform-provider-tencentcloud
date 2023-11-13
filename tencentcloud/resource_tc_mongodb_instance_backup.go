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

Import

mongodb instance_backup can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance_backup.instance_backup instance_backup_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudMongodbInstanceBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceBackupCreate,
		Read:   resourceTencentCloudMongodbInstanceBackupRead,
		Delete: resourceTencentCloudMongodbInstanceBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Description: "Backup notes.",
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_backup.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mongodb.NewCreateBackupDBInstanceRequest()
		response   = mongodb.NewCreateBackupDBInstanceResponse()
		instanceId string
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
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 100*readRetryTimeout, time.Second, service.MongodbInstanceBackupStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
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
