package as

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsScaleInInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScaleInInstancesCreate,
		Read:   resourceTencentCloudAsScaleInInstancesRead,
		Delete: resourceTencentCloudAsScaleInInstancesDelete,
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Scaling group ID.",
			},

			"scale_in_number": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Number of instances to be reduced.",
			},
		},
	}
}

func resourceTencentCloudAsScaleInInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scale_in_instances.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = as.NewScaleInInstancesRequest()
		response   = as.NewScaleInInstancesResponse()
		activityId string
	)
	if v, ok := d.GetOk("auto_scaling_group_id"); ok {
		request.AutoScalingGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("scale_in_number"); ok {
		request.ScaleInNumber = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().ScaleInInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate as scaleInInstances failed, reason:%+v", logId, err)
		return err
	}

	activityId = *response.Response.ActivityId
	d.SetId(activityId)

	return resourceTencentCloudAsScaleInInstancesRead(d, meta)
}

func resourceTencentCloudAsScaleInInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scale_in_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAsScaleInInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scale_in_instances.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
