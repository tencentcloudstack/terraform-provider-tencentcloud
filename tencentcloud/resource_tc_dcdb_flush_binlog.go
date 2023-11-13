/*
Provides a resource to create a dcdb flush_binlog

Example Usage

```hcl
resource "tencentcloud_dcdb_flush_binlog" "flush_binlog" {
  instance_id = ""
}
```

Import

dcdb flush_binlog can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_flush_binlog.flush_binlog flush_binlog_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDcdbFlushBinlog() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbFlushBinlogCreate,
		Read:   resourceTencentCloudDcdbFlushBinlogRead,
		Delete: resourceTencentCloudDcdbFlushBinlogDelete,
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

func resourceTencentCloudDcdbFlushBinlogCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_flush_binlog.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewFlushBinlogRequest()
		response   = dcdb.NewFlushBinlogResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().FlushBinlog(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb flushBinlog failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudDcdbFlushBinlogRead(d, meta)
}

func resourceTencentCloudDcdbFlushBinlogRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_flush_binlog.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbFlushBinlogDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_flush_binlog.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
