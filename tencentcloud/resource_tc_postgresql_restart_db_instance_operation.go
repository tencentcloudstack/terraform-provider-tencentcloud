/*
Provides a resource to create a postgresql restart_db_instance_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_restart_db_instance_operation" "restart_db_instance_operation" {
  db_instance_id = local.pgsql_id
}
```
*/
package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlRestartDbInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlRestartDbInstanceOperationCreate,
		Read:   resourceTencentCloudPostgresqlRestartDbInstanceOperationRead,
		Delete: resourceTencentCloudPostgresqlRestartDbInstanceOperationDelete,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "dbInstance ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlRestartDbInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_restart_db_instance_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgresql.NewRestartDBInstanceRequest()
		dBInstanceId string
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dBInstanceId = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().RestartDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql RestartDbInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dBInstanceId)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 2*readRetryTimeout, 10*time.Second, service.PostgresqlDbInstanceOperationStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudPostgresqlRestartDbInstanceOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlRestartDbInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_restart_db_instance_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlRestartDbInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_restart_db_instance_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
