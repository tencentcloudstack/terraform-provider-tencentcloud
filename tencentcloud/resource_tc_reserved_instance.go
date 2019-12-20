/*
Provides a reserved instance resource.

~> **NOTE:** Reserved instance cannot be deleted and updated. The reserved instance still exist which can be extracted by reserved_instances data source when reserved instance is destroied.

Example Usage

```hcl
resource "tencentcloud_reserved_instance" "ri" {
  config_id      = "469043dd-28b9-4d89-b557-74f6a8326259"
  instance_count = 2
}
```

Import

Reserved instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_reserved_instance.foo 6cc16e7c-47d7-4fae-9b44-ce5c0f59a920
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func resourceTencentCloudReservedInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudReservedInstanceCreate,
		Read:   resourceTencentCloudReservedInstanceRead,
		Update: resourceTencentCloudReservedInstanceUpdate,
		Delete: resourceTencentCloudReservedInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Configuration id of the reserved instance.",
			},
			"instance_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerMin(1),
				Description:  "Number of reserved instances to be purchased.",
			},

			// computed
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Start time of the RI.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiry time of the RI.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the RI at the time of purchase.",
			},
		},
	}
}

func resourceTencentCloudReservedInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_reserved_instance.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	configId := d.Get("config_id").(string)
	count := d.Get("instance_count").(int)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instanceId string
	var errRet error
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		instanceId, errRet = cvmService.CreateReservedInstance(ctx, configId, int64(count))
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(instanceId)

	return resourceTencentCloudReservedInstanceRead(d, meta)
}

func resourceTencentCloudReservedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_reserved_instance.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	filter := map[string]string{
		"reserved-instances-id": id,
	}
	var instances []*cvm.ReservedInstances
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instances, errRet = cvmService.DescribeReservedInstanceByFilter(ctx, filter)
		if errRet != nil {
			return retryError(errRet, "InternalError")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(instances) < 1 {
		d.SetId("")
		return nil
	}
	instance := instances[0]

	_ = d.Set("instance_count", instance.InstanceCount)
	_ = d.Set("start_time", instance.StartTime)
	_ = d.Set("end_time", instance.EndTime)
	_ = d.Set("status", instance.State)

	return nil
}

func resourceTencentCloudReservedInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("reserved instance not allowed to modify")
}

func resourceTencentCloudReservedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
