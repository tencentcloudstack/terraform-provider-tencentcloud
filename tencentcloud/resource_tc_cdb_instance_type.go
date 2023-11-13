/*
Provides a resource to create a cdb instance_type

Example Usage

```hcl
resource "tencentcloud_cdb_instance_type" "instance_type" {
  instance_id = ""
}
```

Import

cdb instance_type can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_instance_type.instance_type instance_type_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"log"
	"time"
)

func resourceTencentCloudCdbInstanceType() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbInstanceTypeCreate,
		Read:   resourceTencentCloudCdbInstanceTypeRead,
		Update: resourceTencentCloudCdbInstanceTypeUpdate,
		Delete: resourceTencentCloudCdbInstanceTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Disaster recovery instance ID in the format of cdb-c1nl9rpv. It is the same as the instance ID displayed in the TencentDB console.",
			},
		},
	}
}

func resourceTencentCloudCdbInstanceTypeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_instance_type.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudCdbInstanceTypeUpdate(d, meta)
}

func resourceTencentCloudCdbInstanceTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_instance_type.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceTypeId := d.Id()

	instanceType, err := service.DescribeCdbInstanceTypeById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceType == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbInstanceType` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceType.InstanceId != nil {
		_ = d.Set("instance_id", instanceType.InstanceId)
	}

	return nil
}

func resourceTencentCloudCdbInstanceTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_instance_type.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewSwitchDrInstanceToMasterRequest()

	instanceTypeId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().SwitchDrInstanceToMaster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb instanceType failed, reason:%+v", logId, err)
		return err
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbInstanceTypeStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbInstanceTypeRead(d, meta)
}

func resourceTencentCloudCdbInstanceTypeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_instance_type.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
