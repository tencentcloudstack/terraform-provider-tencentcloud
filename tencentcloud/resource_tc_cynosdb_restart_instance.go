/*
Provides a resource to create a cynosdb restart_instance

Example Usage

```hcl
resource "tencentcloud_cynosdb_restart_instance" "restart_instance" {
  instance_id = "cynosdbmysql-ins-xxxxxxxx"
}
```

Import

cynosdb restart_instance can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_restart_instance.restart_instance restart_instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCynosdbRestartInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbRestartInstanceCreate,
		Read:   resourceTencentCloudCynosdbRestartInstanceRead,
		Delete: resourceTencentCloudCynosdbRestartInstanceDelete,
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
		},
	}
}

func resourceTencentCloudCynosdbRestartInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_restart_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cynosdb.NewRestartInstanceRequest()
		response   = cynosdb.NewRestartInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().RestartInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb restartInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbRestartInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbRestartInstanceRead(d, meta)
}

func resourceTencentCloudCynosdbRestartInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_restart_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbRestartInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_restart_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
