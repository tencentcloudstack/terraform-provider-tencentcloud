/*
Provides a resource to create a postgres create_base_backup

Example Usage

```hcl
resource "tencentcloud_postgres_create_base_backup" "create_base_backup" {
  d_b_instance_id = ""
  switch_tag =
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgres create_base_backup can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_create_base_backup.create_base_backup create_base_backup_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudPostgresCreateBaseBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresCreateBaseBackupCreate,
		Read:   resourceTencentCloudPostgresCreateBaseBackupRead,
		Delete: resourceTencentCloudPostgresCreateBaseBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of the instance waiting for a switch.",
			},

			"switch_tag": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Valid value: `0` (switch immediately).",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresCreateBaseBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_create_base_backup.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgres.NewModifySwitchTimePeriodRequest()
		response     = postgres.NewModifySwitchTimePeriodResponse()
		dBInstanceId string
	)
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
		request.DBInstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("switch_tag"); v != nil {
		request.SwitchTag = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifySwitchTimePeriod(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres CreateBaseBackup failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresCreateBaseBackupRead(d, meta)
}

func resourceTencentCloudPostgresCreateBaseBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_create_base_backup.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresCreateBaseBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_create_base_backup.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
