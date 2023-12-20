package as

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsStopInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsStopInstancesCreate,
		Read:   resourceTencentCloudAsStopInstancesRead,
		Delete: resourceTencentCloudAsStopInstancesDelete,
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
				Description: "List of cvm instances to stop.",
			},

			"stopped_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Billing method of a pay-as-you-go instance after shutdown. Available values: `KEEP_CHARGING`,`STOP_CHARGING`. Default `KEEP_CHARGING`.",
			},
		},
	}
}

func resourceTencentCloudAsStopInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_as_stop_instances.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = as.NewStopAutoScalingInstancesRequest()
		response   = as.NewStopAutoScalingInstancesResponse()
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

	if v, ok := d.GetOk("stopped_mode"); ok {
		request.StoppedMode = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().StopAutoScalingInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate as stopInstances failed, reason:%+v", logId, err)
		return nil
	}

	activityId = *response.Response.ActivityId

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err = resource.Retry(4*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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

	return resourceTencentCloudAsStopInstancesRead(d, meta)
}

func resourceTencentCloudAsStopInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_stop_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAsStopInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_stop_instances.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
