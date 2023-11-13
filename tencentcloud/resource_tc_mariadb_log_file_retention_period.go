/*
Provides a resource to create a mariadb log_file_retention_period

Example Usage

```hcl
resource "tencentcloud_mariadb_log_file_retention_period" "log_file_retention_period" {
  instance_id = &lt;nil&gt;
  days = &lt;nil&gt;
}
```

Import

mariadb log_file_retention_period can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_log_file_retention_period.log_file_retention_period log_file_retention_period_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"log"
)

func resourceTencentCloudMariadbLogFileRetentionPeriod() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbLogFileRetentionPeriodCreate,
		Read:   resourceTencentCloudMariadbLogFileRetentionPeriodRead,
		Update: resourceTencentCloudMariadbLogFileRetentionPeriodUpdate,
		Delete: resourceTencentCloudMariadbLogFileRetentionPeriodDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Retention days.",
			},
		},
	}
}

func resourceTencentCloudMariadbLogFileRetentionPeriodCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_log_file_retention_period.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbLogFileRetentionPeriodUpdate(d, meta)
}

func resourceTencentCloudMariadbLogFileRetentionPeriodRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_log_file_retention_period.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	logFileRetentionPeriodId := d.Id()

	logFileRetentionPeriod, err := service.DescribeMariadbLogFileRetentionPeriodById(ctx, instanceId)
	if err != nil {
		return err
	}

	if logFileRetentionPeriod == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbLogFileRetentionPeriod` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if logFileRetentionPeriod.InstanceId != nil {
		_ = d.Set("instance_id", logFileRetentionPeriod.InstanceId)
	}

	if logFileRetentionPeriod.Days != nil {
		_ = d.Set("days", logFileRetentionPeriod.Days)
	}

	return nil
}

func resourceTencentCloudMariadbLogFileRetentionPeriodUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_log_file_retention_period.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mariadb.NewModifyLogFileRetentionPeriodRequest()

	logFileRetentionPeriodId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "days"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyLogFileRetentionPeriod(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mariadb logFileRetentionPeriod failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbLogFileRetentionPeriodRead(d, meta)
}

func resourceTencentCloudMariadbLogFileRetentionPeriodDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_log_file_retention_period.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
