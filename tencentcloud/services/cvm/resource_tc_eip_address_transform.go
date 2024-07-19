package cvm

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEipAddressTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipAddressTransformCreate,
		Read:   resourceTencentCloudEipAddressTransformRead,
		Delete: resourceTencentCloudEipAddressTransformDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the instance ID of a normal public network IP to be operated. eg:ins-23mk45jn.",
			},
		},
	}
}

func resourceTencentCloudEipAddressTransformCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_address_transform.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId string
	)
	var (
		request  = vpc.NewTransformAddressRequest()
		response = vpc.NewTransformAddressResponse()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().TransformAddressWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create eip address transform failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	if err := resourceTencentCloudEipAddressTransformCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudEipAddressTransformRead(d, meta)
}

func resourceTencentCloudEipAddressTransformRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_address_transform.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudEipAddressTransformDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_address_transform.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
