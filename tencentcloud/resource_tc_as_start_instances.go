/*
Provides a resource to create a as start_instances

Example Usage

```hcl
resource "tencentcloud_as_start_instances" "start_instances" {
  auto_scaling_group_id = ""
  instance_ids = ""
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAsStartInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsStartInstancesCreate,
		Read:   resourceTencentCloudAsStartInstancesRead,
		Delete: resourceTencentCloudAsStartInstancesDelete,
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
				Description: "List of cvm instances to start.",
			},
		},
	}
}

func resourceTencentCloudAsStartInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_as_start_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = as.NewStartAutoScalingInstancesRequest()
		response   = as.NewStartAutoScalingInstancesResponse()
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().StartAutoScalingInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Println("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Println("[CRITAL]%s operate as startInstances failed, reason:%+v", logId, err)
		return nil
	}

	activityId = *response.Response.ActivityId

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err = resource.Retry(4*readRetryTimeout, func() *resource.RetryError {
		status, err := service.DescribeActivityById(ctx, activityId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if status == SCALING_GROUP_ACTIVITY_STATUS_INIT || status == SCALING_GROUP_ACTIVITY_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("remove status is running(%s)", status))
		}
		if status == SCALING_GROUP_ACTIVITY_STATUS_SUCCESSFUL {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("remove status is failed(%s)", status))
	})
	if err != nil {
		return err
	}

	d.SetId(activityId)

	return resourceTencentCloudAsStartInstancesRead(d, meta)
}

func resourceTencentCloudAsStartInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_start_instances.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAsStartInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_start_instances.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
