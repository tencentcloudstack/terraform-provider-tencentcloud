/*
Provides a resource to create a postgres base_backup

Example Usage

```hcl
resource "tencentcloud_postgres_base_backup" "base_backup" {
  d_b_instance_id = ""
  base_backup_id = ""
  new_expire_time = ""
}
```

Import

postgres base_backup can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_base_backup.base_backup base_backup_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"log"
	"strings"
)

func resourceTencentCloudPostgresBaseBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresBaseBackupCreate,
		Read:   resourceTencentCloudPostgresBaseBackupRead,
		Update: resourceTencentCloudPostgresBaseBackupUpdate,
		Delete: resourceTencentCloudPostgresBaseBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"base_backup_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Base backup ID.",
			},

			"new_expire_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "New expiration time.",
			},
		},
	}
}

func resourceTencentCloudPostgresBaseBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_base_backup.create")()
	defer inconsistentCheck(d, meta)()

	var dBInstanceId string
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
	}

	var baseBackupId string
	if v, ok := d.GetOk("base_backup_id"); ok {
		baseBackupId = v.(string)
	}

	d.SetId(strings.Join([]string{dBInstanceId, baseBackupId}, FILED_SP))

	return resourceTencentCloudPostgresBaseBackupUpdate(d, meta)
}

func resourceTencentCloudPostgresBaseBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_base_backup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dBInstanceId := idSplit[0]
	baseBackupId := idSplit[1]

	BaseBackup, err := service.DescribePostgresBaseBackupById(ctx, dBInstanceId, baseBackupId)
	if err != nil {
		return err
	}

	if BaseBackup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresBaseBackup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if BaseBackup.DBInstanceId != nil {
		_ = d.Set("d_b_instance_id", BaseBackup.DBInstanceId)
	}

	if BaseBackup.BaseBackupId != nil {
		_ = d.Set("base_backup_id", BaseBackup.BaseBackupId)
	}

	if BaseBackup.NewExpireTime != nil {
		_ = d.Set("new_expire_time", BaseBackup.NewExpireTime)
	}

	return nil
}

func resourceTencentCloudPostgresBaseBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_base_backup.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgres.NewModifyBaseBackupExpireTimeRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dBInstanceId := idSplit[0]
	baseBackupId := idSplit[1]

	request.DBInstanceId = &dBInstanceId
	request.BaseBackupId = &baseBackupId

	immutableArgs := []string{"d_b_instance_id", "base_backup_id", "new_expire_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifyBaseBackupExpireTime(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgres BaseBackup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresBaseBackupRead(d, meta)
}

func resourceTencentCloudPostgresBaseBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_base_backup.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
