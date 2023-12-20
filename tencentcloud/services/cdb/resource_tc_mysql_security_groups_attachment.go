package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlSecurityGroupsAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlSecurityGroupsAttachmentCreate,
		Read:   resourceTencentCloudMysqlSecurityGroupsAttachmentRead,
		Delete: resourceTencentCloudMysqlSecurityGroupsAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of security group.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The id of instance.",
			},
		},
	}
}

func resourceTencentCloudMysqlSecurityGroupsAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_security_groups_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request         = mysql.NewAssociateSecurityGroupsRequest()
		securityGroupId string
		instanceId      string
	)
	if v, ok := d.GetOk("security_group_id"); ok {
		securityGroupId = v.(string)
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{helper.String(v.(string))}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().AssociateSecurityGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql securityGroupsAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(securityGroupId + tccommon.FILED_SP + instanceId)

	return resourceTencentCloudMysqlSecurityGroupsAttachmentRead(d, meta)
}

func resourceTencentCloudMysqlSecurityGroupsAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_security_groups_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroupId := idSplit[0]
	instanceId := idSplit[1]

	securityGroupsAttachment, err := service.DescribeMysqlSecurityGroupsAttachmentById(ctx, securityGroupId, instanceId)
	if err != nil {
		return err
	}

	if securityGroupsAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlSecurityGroupsAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil

	}
	_ = d.Set("instance_id", instanceId)
	if securityGroupsAttachment.SecurityGroupId != nil {
		_ = d.Set("security_group_id", securityGroupsAttachment.SecurityGroupId)
	}

	return nil
}

func resourceTencentCloudMysqlSecurityGroupsAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_security_groups_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroupId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteMysqlSecurityGroupsAttachmentById(ctx, securityGroupId, instanceId); err != nil {
		return err
	}

	return nil
}
