package cvm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmSecurityGroupAttachmentCreate,
		Read:   resourceTencentCloudCvmSecurityGroupAttachmentRead,
		Delete: resourceTencentCloudCvmSecurityGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the security group to be associated, such as sg-efil73jd. Only one security group can be associated.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID. To obtain the instance IDs, you can call DescribeInstances and look for InstanceId in the response.",
			},
		},
	}
}

func resourceTencentCloudCvmSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_security_group_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId      string
		securityGroupId string
	)
	var (
		request  = cvm.NewAssociateSecurityGroupsRequest()
		response = cvm.NewAssociateSecurityGroupsResponse()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		securityGroupId = v.(string)
	}

	request.SecurityGroupIds = []*string{helper.String(securityGroupId)}

	request.InstanceIds = []*string{helper.String(instanceId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().AssociateSecurityGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm security group attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(strings.Join([]string{instanceId, securityGroupId}, tccommon.FILED_SP))

	return resourceTencentCloudCvmSecurityGroupAttachmentRead(d, meta)
}

func resourceTencentCloudCvmSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_security_group_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	//logId := tccommon.GetLogId(tccommon.ContextNil)
	//
	//ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	//
	//service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	//
	//idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	//if len(idSplit) != 2 {
	//	return fmt.Errorf("id is broken,%s", d.Id())
	//}
	//instanceId := idSplit[0]
	//securityGroupId := idSplit[1]
	//
	//_ = d.Set("instance_id", instanceId)
	//
	//_ = d.Set("security_group_id", securityGroupId)

	//respData, err := service.DescribeCvmSecurityGroupAttachmentById(ctx, instanceId)
	//if err != nil {
	//	return err
	//}
	//
	//if respData == nil {
	//	d.SetId("")
	//	log.Printf("[WARN]%s resource `cvm_security_group_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
	//	return nil
	//}
	//if err := resourceTencentCloudCvmSecurityGroupAttachmentReadPostHandleResponse0(ctx, respData); err != nil {
	//	return err
	//}
	//
	//_ = securityGroupId
	return nil
}

func resourceTencentCloudCvmSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_security_group_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	var (
		request  = cvm.NewDisassociateSecurityGroupsRequest()
		response = cvm.NewDisassociateSecurityGroupsResponse()
	)

	request.SecurityGroupIds = []*string{helper.String(securityGroupId)}

	request.InstanceIds = []*string{helper.String(instanceId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().DisassociateSecurityGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm security group attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
