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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCynosdbInstanceName() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbInstanceNameCreate,
		Read:   resourceTencentCloudCynosdbInstanceNameRead,
		Delete: resourceTencentCloudCynosdbInstanceNameDelete,
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

			"instance_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance Name.",
			},
		},
	}
}

func resourceTencentCloudCynosdbInstanceNameCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_name.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cynosdb.NewModifyInstanceNameRequest()
		response   = cynosdb.NewModifyInstanceNameResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb instanceName failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudCynosdbInstanceNameRead(d, meta)
}

func resourceTencentCloudCynosdbInstanceNameRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_name.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbInstanceNameDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_name.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
