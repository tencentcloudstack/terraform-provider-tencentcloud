/*
Provides a resource to create a cdb isolate_instance

Example Usage

```hcl
resource "tencentcloud_cdb_isolate_instance" "isolate_instance" {
  instance_id = ""
}
```

Import

cdb isolate_instance can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_isolate_instance.isolate_instance isolate_instance_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCdbIsolateInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbIsolateInstanceCreate,
		Read:   resourceTencentCloudCdbIsolateInstanceRead,
		Delete: resourceTencentCloudCdbIsolateInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, and you can use the [query instance list] (https://cloud.tencent.com/document/api/236/15872) interface Gets the value of the field InstanceId in the output parameter.",
			},
		},
	}
}

func resourceTencentCloudCdbIsolateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_isolate_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewIsolateDBInstanceRequest()
		response   = cdb.NewIsolateDBInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().IsolateDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdb isolateInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbIsolateInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbIsolateInstanceRead(d, meta)
}

func resourceTencentCloudCdbIsolateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_isolate_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	isolateInstanceId := d.Id()

	isolateInstance, err := service.DescribeCdbIsolateInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if isolateInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbIsolateInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if isolateInstance.InstanceId != nil {
		_ = d.Set("instance_id", isolateInstance.InstanceId)
	}

	return nil
}

func resourceTencentCloudCdbIsolateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_isolate_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	isolateInstanceId := d.Id()

	if err := service.DeleteCdbIsolateInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
