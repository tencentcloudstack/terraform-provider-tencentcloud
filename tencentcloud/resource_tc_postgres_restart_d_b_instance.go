/*
Provides a resource to create a postgres restart_d_b_instance

Example Usage

```hcl
resource "tencentcloud_postgres_restart_d_b_instance" "restart_d_b_instance" {
  d_b_instance_id = "postgres-6r233v55"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgres restart_d_b_instance can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_restart_d_b_instance.restart_d_b_instance restart_d_b_instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudPostgresRestartDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresRestartDBInstanceCreate,
		Read:   resourceTencentCloudPostgresRestartDBInstanceRead,
		Delete: resourceTencentCloudPostgresRestartDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "DbInstance ID.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresRestartDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_restart_d_b_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgres.NewRestartDBInstanceRequest()
		response     = postgres.NewRestartDBInstanceResponse()
		dBInstanceId string
	)
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
		request.DBInstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().RestartDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres RestartDBInstance failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(dBInstanceId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 120*readRetryTimeout, time.Second, service.PostgresRestartDBInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudPostgresRestartDBInstanceRead(d, meta)
}

func resourceTencentCloudPostgresRestartDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_restart_d_b_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresRestartDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_restart_d_b_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
