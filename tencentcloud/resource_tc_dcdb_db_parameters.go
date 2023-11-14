/*
Provides a resource to create a dcdb db_parameters

Example Usage

```hcl
resource "tencentcloud_dcdb_db_parameters" "db_parameters" {
  instance_id = ""
  sync_mode =
}
```

Import

dcdb db_parameters can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_db_parameters.db_parameters db_parameters_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"log"
)

func resourceTencentCloudDcdbDbParameters() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbDbParametersCreate,
		Read:   resourceTencentCloudDcdbDbParametersRead,
		Update: resourceTencentCloudDcdbDbParametersUpdate,
		Delete: resourceTencentCloudDcdbDbParametersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the instance for which to modify the sync mode. The ID is in the format of `tdsql-ow728lmc`.",
			},

			"sync_mode": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Sync mode. Valid values: `0` (async), `1` (strong sync), `2` (downgradable strong sync).",
			},
		},
	}
}

func resourceTencentCloudDcdbDbParametersCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_parameters.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudDcdbDbParametersUpdate(d, meta)
}

func resourceTencentCloudDcdbDbParametersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_parameters.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	dbParametersId := d.Id()

	dbParameters, err := service.DescribeDcdbDbParametersById(ctx, instanceId)
	if err != nil {
		return err
	}

	if dbParameters == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbDbParameters` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dbParameters.InstanceId != nil {
		_ = d.Set("instance_id", dbParameters.InstanceId)
	}

	if dbParameters.SyncMode != nil {
		_ = d.Set("sync_mode", dbParameters.SyncMode)
	}

	return nil
}

func resourceTencentCloudDcdbDbParametersUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_parameters.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dcdb.NewModifyDBSyncModeRequest()

	dbParametersId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "sync_mode"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyDBSyncMode(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dcdb dbParameters failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbDbParametersRead(d, meta)
}

func resourceTencentCloudDcdbDbParametersDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_parameters.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
