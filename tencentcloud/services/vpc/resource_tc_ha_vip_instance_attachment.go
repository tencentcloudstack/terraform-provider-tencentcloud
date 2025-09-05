package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudHaVipInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudHaVipInstanceAttachmentCreate,
		Read:   resourceTencentCloudHaVipInstanceAttachmentRead,
		Delete: resourceTencentCloudHaVipInstanceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique ID of the slave machine or network card to which HaVip is bound.",
			},
			"ha_vip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Unique ID of the HaVip instance.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The type of HaVip binding. Values:CVM, ENI.",
			},
		},
	}
}

func resourceTencentCloudHaVipInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ha_vip_instance_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId   string
		haVipId      string
		instanceType string
	)
	var (
		request  = vpc.NewAssociateHaVipInstanceRequest()
		response = vpc.NewAssociateHaVipInstanceResponse()
	)

	instanceId = d.Get("instance_id").(string)
	haVipId = d.Get("ha_vip_id").(string)
	instanceType = d.Get("instance_type").(string)

	request.HaVipAssociationSet = []*vpc.HaVipAssociation{
		{
			HaVipId:      helper.String(haVipId),
			InstanceType: helper.String(instanceType),
			InstanceId:   helper.String(instanceId),
		},
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssociateHaVipInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ha vip instance attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(strings.Join([]string{haVipId, instanceId}, tccommon.FILED_SP))

	return resourceTencentCloudHaVipInstanceAttachmentRead(d, meta)
}

func resourceTencentCloudHaVipInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ha_vip_instance_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	haVipId := idSplit[0]
	instanceId := idSplit[1]

	var haVips []*vpc.HaVip
	filters := map[string]string{
		"havip-id":                      haVipId,
		"havip-association.instance-id": instanceId,
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeHaVipByFilter(ctx, filters)
		if e != nil {
			return tccommon.RetryError(e)
		}
		haVips = result
		return nil
	})
	if err != nil {
		return err
	}

	if len(haVips) == 0 || len(haVips[0].HaVipAssociationSet) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `ha_vip_instance_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("ha_vip_id", haVips[0].HaVipAssociationSet[0].HaVipId)
	_ = d.Set("instance_type", haVips[0].HaVipAssociationSet[0].InstanceType)
	_ = d.Set("instance_id", haVips[0].HaVipAssociationSet[0].InstanceId)

	_ = instanceId
	_ = haVipId
	return nil
}

func resourceTencentCloudHaVipInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ha_vip_instance_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	haVipId := idSplit[0]
	instanceId := idSplit[1]

	var (
		request  = vpc.NewDisassociateHaVipInstanceRequest()
		response = vpc.NewDisassociateHaVipInstanceResponse()
	)

	request.HaVipAssociationSet = []*vpc.HaVipAssociation{
		{
			HaVipId:    helper.String(haVipId),
			InstanceId: helper.String(instanceId),
		},
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DisassociateHaVipInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete ha vip instance attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = instanceId
	_ = haVipId
	return nil
}
