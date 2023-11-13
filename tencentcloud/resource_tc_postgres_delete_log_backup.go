/*
Provides a resource to create a postgres delete_log_backup

Example Usage

```hcl
resource "tencentcloud_postgres_delete_log_backup" "delete_log_backup" {
  d_b_instance_id = ""
  log_backup_id = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgres delete_log_backup can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_delete_log_backup.delete_log_backup delete_log_backup_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudPostgresDeleteLogBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresDeleteLogBackupCreate,
		Read:   resourceTencentCloudPostgresDeleteLogBackupRead,
		Delete: resourceTencentCloudPostgresDeleteLogBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"log_backup_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Log backup ID.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresDeleteLogBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_delete_log_backup.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgres.NewDeleteLogBackupRequest()
		response     = postgres.NewDeleteLogBackupResponse()
		dBInstanceId string
		logBackupId  string
	)
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
		request.DBInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_backup_id"); ok {
		logBackupId = v.(string)
		request.LogBackupId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().DeleteLogBackup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres DeleteLogBackup failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(strings.Join([]string{dBInstanceId, logBackupId}, FILED_SP))

	return resourceTencentCloudPostgresDeleteLogBackupRead(d, meta)
}

func resourceTencentCloudPostgresDeleteLogBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_delete_log_backup.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresDeleteLogBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_delete_log_backup.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
