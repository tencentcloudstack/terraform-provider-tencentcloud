package as

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
)

func ResourceTencentCloudAsScalingGroupStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScalingGroupStatusCreate,
		Read:   resourceTencentCloudAsScalingGroupStatusRead,
		Update: resourceTencentCloudAsScalingGroupStatusUpdate,
		Delete: resourceTencentCloudAsScalingGroupStatusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Scaling group ID.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "If enable auto scaling group.",
			},
		},
	}
}

func resourceTencentCloudAsScalingGroupStatusCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group_status.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	autoScalingGroupId := d.Get("auto_scaling_group_id").(string)

	d.SetId(autoScalingGroupId)

	return resourceTencentCloudAsScalingGroupStatusUpdate(d, meta)
}

func resourceTencentCloudAsScalingGroupStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := AsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	autoScalingGroupId := d.Id()

	scalingGroup, has, err := service.DescribeAutoScalingGroupById(ctx, autoScalingGroupId)
	if err != nil {
		return err
	}

	if has == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `AsScalingGroupStatus` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if scalingGroup.AutoScalingGroupId != nil {
		_ = d.Set("auto_scaling_group_id", scalingGroup.AutoScalingGroupId)
	}

	if scalingGroup.EnabledStatus != nil {
		if *scalingGroup.EnabledStatus == "ENABLED" {
			_ = d.Set("enable", true)
		} else {
			_ = d.Set("enable", false)
		}
	}

	return nil
}

func resourceTencentCloudAsScalingGroupStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group_status.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		enable           bool
		enableAsRequest  = as.NewEnableAutoScalingGroupRequest()
		disableAsRequest = as.NewDisableAutoScalingGroupRequest()
	)

	autoScalingGroupId := d.Id()

	if v, ok := d.GetOkExists("enable"); ok {
		enable = v.(bool)
	}

	if enable {
		enableAsRequest.AutoScalingGroupId = &autoScalingGroupId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().EnableAutoScalingGroup(enableAsRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableAsRequest.GetAction(), enableAsRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s enable vpc snapshotPolicyConfig failed, reason:%+v", logId, err)
			return err
		}
	} else {
		disableAsRequest.AutoScalingGroupId = &autoScalingGroupId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().DisableAutoScalingGroup(disableAsRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableAsRequest.GetAction(), enableAsRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s disable vpc snapshotPolicyConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudAsScalingGroupStatusRead(d, meta)
}

func resourceTencentCloudAsScalingGroupStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group_status.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
