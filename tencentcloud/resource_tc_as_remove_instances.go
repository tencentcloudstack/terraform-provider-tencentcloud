/*
Provides a resource to create a as remove_instances

Example Usage

```hcl
resource "tencentcloud_as_remove_instances" "remove_instances" {
  auto_scaling_group_id = ""
  instance_ids = ""
}
```

Import

as remove_instances can be imported using the id, e.g.

```
terraform import tencentcloud_as_remove_instances.remove_instances remove_instances_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudAsRemoveInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsRemoveInstancesCreate,
		Read:   resourceTencentCloudAsRemoveInstancesRead,
		Delete: resourceTencentCloudAsRemoveInstancesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Launch configuration ID.",
			},

			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of cvm instances to remove .",
			},
		},
	}
}

func resourceTencentCloudAsRemoveInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_as_remove_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = as.NewRemoveInstancesRequest()
		response   = as.NewRemoveInstancesResponse()
		activityId string
	)
	if v, ok := d.GetOk("auto_scaling_group_id"); ok {
		request.AutoScalingGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().RemoveInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Println("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Println("[CRITAL]%s operate as removeInstances failed, reason:%+v", logId, err)
		return nil
	}

	activityId = *response.Response.ActivityId
	d.SetId(activityId)

	return resourceTencentCloudAsRemoveInstancesRead(d, meta)
}

func resourceTencentCloudAsRemoveInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_remove_instances.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAsRemoveInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_remove_instances.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
