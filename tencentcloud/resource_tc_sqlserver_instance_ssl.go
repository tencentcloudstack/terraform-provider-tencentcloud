/*
Provides a resource to create a sqlserver instance_ssl

Example Usage

```hcl
resource "tencentcloud_sqlserver_instance_ssl" "instance_ssl" {
  instance_id = "mssql-i1z41iwd"
  type = "enable"
  wait_switch = 0
}
```

Import

sqlserver instance_ssl can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_instance_ssl.instance_ssl instance_ssl_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSqlserverInstanceSsl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverInstanceSslCreate,
		Read:   resourceTencentCloudSqlserverInstanceSslRead,
		Update: resourceTencentCloudSqlserverInstanceSslUpdate,
		Delete: resourceTencentCloudSqlserverInstanceSslDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Operation type. enable: turn on SSL, disable: turn off SSL, renew: update the certificate validity period.",
			},
		},
	}
}

func resourceTencentCloudSqlserverInstanceSslCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_ssl.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverInstanceSslUpdate(d, meta)
}

func resourceTencentCloudSqlserverInstanceSslRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_ssl.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	instanceSsl, err := service.DescribeSqlserverInstanceSslById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceSsl == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverInstanceSsl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceSsl.InstanceId != nil {
		_ = d.Set("instance_id", instanceSsl.InstanceId)
	}

	if instanceSsl.Type != nil {
		_ = d.Set("type", instanceSsl.Type)
	}

	return nil
}

func resourceTencentCloudSqlserverInstanceSslUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_ssl.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = sqlserver.NewModifyDBInstanceSSLRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	request.WaitSwitch = helper.IntUint64(0)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBInstanceSSL(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver instanceSsl failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverInstanceSslRead(d, meta)
}

func resourceTencentCloudSqlserverInstanceSslDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_instance_ssl.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
