/*
Provides a resource to create a as protect_instances

Example Usage

```hcl
resource "tencentcloud_as_protect_instances" "protect_instances" {
  auto_scaling_group_id = ""
  instance_ids = ""
  protected_from_scale_in = ""
}
```

*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAsProtectInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsProtectInstancesCreate,
		Read:   resourceTencentCloudAsProtectInstancesRead,
		Delete: resourceTencentCloudAsProtectInstancesDelete,
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
				Description: "List of cvm instances to remove.",
			},

			"protected_from_scale_in": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "If instances need protect.",
			},
		},
	}
}

func resourceTencentCloudAsProtectInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_as_protect_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = as.NewSetInstancesProtectionRequest()
		response = as.NewSetInstancesProtectionResponse()
		//instanceIds string
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

	if v, _ := d.GetOk("protected_from_scale_in"); v != nil {
		request.ProtectedFromScaleIn = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().SetInstancesProtection(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Println("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Println("[CRITAL]%s operate as protectInstances failed, reason:%+v", logId, err)
		return nil
	}

	// 需要 setId，可以通过InstancesId作为ID
	d.SetId("")

	return resourceTencentCloudAsProtectInstancesRead(d, meta)
}

func resourceTencentCloudAsProtectInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_protect_instances.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAsProtectInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_protect_instances.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
