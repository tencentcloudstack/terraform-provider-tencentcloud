package vpc

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudVpcAclAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcAclAttachmentCreate,
		Read:   resourceTencentCloudVpcAclAttachmentRead,
		Delete: resourceTencentCloudVpcAclAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "ID of the attached ACL.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Subnet instance ID.",
			},
		},
	}
}

func resourceTencentCloudVpcAclAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_acl_attachment.create")()
	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		subnetIds []string
	)

	aclId := d.Get("acl_id").(string)
	subnetId := d.Get("subnet_id").(string)

	subnetIds = append(subnetIds, subnetId)

	err := service.AssociateAclSubnets(ctx, aclId, subnetIds)
	if err != nil {
		return err
	}

	d.SetId(aclId + tccommon.FILED_SP + subnetId)

	return resourceTencentCloudVpcAclAttachmentRead(d, meta)
}

func resourceTencentCloudVpcAclAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_acl_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service      = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		attachmentId = d.Id()
	)

	idSplit := strings.Split(attachmentId, tccommon.FILED_SP)
	if len(idSplit) < 2 {
		d.SetId("")
		return nil
	}
	aclId := idSplit[0]
	subnetId := idSplit[1]
	results, err := service.DescribeNetWorkAcls(ctx, aclId, "", "")
	if err != nil {
		return err
	}
	if len(results) < 1 || len(results[0].SubnetSet) < 1 {
		d.SetId("")
		return nil
	}

	_ = d.Set("acl_id", aclId)
	_ = d.Set("subnet_id", subnetId)

	return nil

}

func resourceTencentCloudVpcAclAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_acl_attachment.delete")()
	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		attachmentAcl = d.Id()
	)

	err := service.DeleteAclAttachment(ctx, attachmentAcl)
	if err != nil {
		log.Printf("[CRITAL]%s delete ACL attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil

}
