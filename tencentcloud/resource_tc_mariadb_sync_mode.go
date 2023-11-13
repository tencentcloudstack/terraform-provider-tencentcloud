/*
Provides a resource to create a mariadb sync_mode

Example Usage

```hcl
resource "tencentcloud_mariadb_sync_mode" "sync_mode" {
  instance_id = ""
  sync_mode =
}
```

Import

mariadb sync_mode can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_sync_mode.sync_mode sync_mode_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMariadbSyncMode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbSyncModeCreate,
		Read:   resourceTencentCloudMariadbSyncModeRead,
		Delete: resourceTencentCloudMariadbSyncModeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the instance for which to modify the sync mode. The ID is in the format of `tdsql-ow728lmc`.",
			},

			"sync_mode": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Sync mode. Valid values: `0` (async), `1` (strong sync), `2` (downgradable strong sync).",
			},
		},
	}
}

func resourceTencentCloudMariadbSyncModeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_sync_mode.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewModifyDBSyncModeRequest()
		response   = mariadb.NewModifyDBSyncModeResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("sync_mode"); v != nil {
		request.SyncMode = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyDBSyncMode(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb syncMode failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudMariadbSyncModeRead(d, meta)
}

func resourceTencentCloudMariadbSyncModeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_sync_mode.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbSyncModeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_sync_mode.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
