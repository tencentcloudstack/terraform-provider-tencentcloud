package dcdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDcdbSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDcdbSecurityGroupAttachmentRead,
		Create: resourceTencentCloudDcdbSecurityGroupAttachmentCreate,
		Delete: resourceTencentCloudDcdbSecurityGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "security group id.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "attached instance id.",
			},
		},
	}
}

func resourceTencentCloudDcdbSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_security_group_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = dcdb.NewAssociateSecurityGroupsRequest()
		// response        *dcdb.AssociateSecurityGroupsResponse
		instanceId      string
		securityGroupId string
	)

	request.Product = helper.String("dcdb") // api only use this fixed value

	if v, ok := d.GetOk("security_group_id"); ok {
		securityGroupId = v.(string)
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{helper.String(v.(string))}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().AssociateSecurityGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		// response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dcdb securityGroupAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + securityGroupId)
	return resourceTencentCloudDcdbSecurityGroupAttachmentRead(d, meta)
}

func resourceTencentCloudDcdbSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_security_group_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	securityGroupAttachment, err := service.DescribeDcdbSecurityGroup(ctx, instanceId)

	if err != nil {
		return err
	}

	if securityGroupAttachment == nil {
		d.SetId("")
		return fmt.Errorf("resource `securityGroupAttachment` %s does not exist", d.Id())
	}

	for _, sg := range securityGroupAttachment.Groups {
		if securityGroupId == *sg.SecurityGroupId {
			_ = d.Set("instance_id", instanceId)
			_ = d.Set("security_group_id", securityGroupId)
			return nil
		}
	}

	return fmt.Errorf("The security group get from api does not match with current instance %v", d.Id())
}

func resourceTencentCloudDcdbSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_security_group_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	if err := service.DeleteDcdbSecurityGroupAttachmentById(ctx, instanceId, securityGroupId); err != nil {
		return err
	}

	return nil
}
