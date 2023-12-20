package as

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsProtectInstances() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_as_protect_instances.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = as.NewSetInstancesProtectionRequest()
	)
	if v, ok := d.GetOk("auto_scaling_group_id"); ok {
		request.AutoScalingGroupId = helper.String(v.(string))
	}
	ids := make([]string, 0)
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			ids = append(ids, instanceIds)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	if v, _ := d.GetOk("protected_from_scale_in"); v != nil {
		request.ProtectedFromScaleIn = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().SetInstancesProtection(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate as protectInstances failed, reason:%+v", logId, err)
		return nil
	}

	// 需要 setId，可以通过InstancesId作为ID
	d.SetId(helper.DataResourceIdsHash(ids))

	return resourceTencentCloudAsProtectInstancesRead(d, meta)
}

func resourceTencentCloudAsProtectInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_protect_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAsProtectInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_protect_instances.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
