/*
Provides a resource to create a cynosdb instance_name

Example Usage

```hcl
resource "tencentcloud_cynosdb_instance_name" "instance_name" {
  instance_id = "cynosdb-ins-dokydbam"
  instance_name = "newName"
}
```

Import

cynosdb instance_name can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_instance_name.instance_name instance_name_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbInstanceName() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbInstanceNameCreate,
		Read:   resourceTencentCloudCynosdbInstanceNameRead,
		Update: resourceTencentCloudCynosdbInstanceNameUpdate,
		Delete: resourceTencentCloudCynosdbInstanceNameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:	 true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"instance_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance Name.",
			},
		},
	}
}

func resourceTencentCloudCynosdbInstanceNameCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_name.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudCynosdbInstanceNameUpdate(d, meta)
}

func resourceTencentCloudCynosdbInstanceNameRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_name.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()
	_, instance, has, err := service.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if !has {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if instance.InstanceName != nil {
		_ = d.Set("instance_name", instance.InstanceName)
	}

	return nil
}

func resourceTencentCloudCynosdbInstanceNameUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_name.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyInstanceNameRequest()

	instanceId := d.Id()
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyInstanceName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb instanceName failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbInstanceNameRead(d, meta)
}

func resourceTencentCloudCynosdbInstanceNameDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_name.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
