/*
Provides a resource to create a postgresql isolate_db_instance_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_isolate_db_instance_operation" "isolate_db_instance_operation" {
  db_instance_id_set = [local.pgsql_id]
}
```
*/
package tencentcloud

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
)

func resourceTencentCloudPostgresqlIsolateDbInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlIsolateDbInstanceOperationCreate,
		Read:   resourceTencentCloudPostgresqlIsolateDbInstanceOperationRead,
		Delete: resourceTencentCloudPostgresqlIsolateDbInstanceOperationDelete,
		Schema: map[string]*schema.Schema{
			"db_instance_id_set": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of resource IDs. Note that currently you cannot isolate multiple instances at the same time. Only one instance ID can be passed in here.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlIsolateDbInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_isolate_db_instance_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = postgresql.NewIsolateDBInstancesRequest()
		ids             []string
		firstInstanceId string
	)
	if v, ok := d.GetOk("db_instance_id_set"); ok {
		dBInstanceIdSetSet := v.(*schema.Set).List()
		for i := range dBInstanceIdSetSet {
			if dBInstanceIdSetSet[i] != nil {
				dBInstanceIdSet := dBInstanceIdSetSet[i].(string)
				request.DBInstanceIdSet = append(request.DBInstanceIdSet, &dBInstanceIdSet)
				ids = append(ids, dBInstanceIdSet)
			}
		}
		firstInstanceId = dBInstanceIdSetSet[0].(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().IsolateDBInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql IsolateDbInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join(ids, FILED_SP))

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"isolated"}, 10*readRetryTimeout, 10*time.Second, service.PostgresqlDbInstanceOperationStateRefreshFunc(firstInstanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudPostgresqlIsolateDbInstanceOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlIsolateDbInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_isolate_db_instance_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlIsolateDbInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_isolate_db_instance_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
