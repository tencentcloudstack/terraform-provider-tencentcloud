package cls

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsConfigAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsConfigAttachmentCreate,
		Read:   resourceTencentCloudClsConfigAttachmentRead,
		Delete: resourceTencentCloudClsConfigAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Collection configuration id.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Machine group id.",
			},
		},
	}
}

func resourceTencentCloudClsConfigAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config_attachment.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cls.NewApplyConfigToMachineGroupRequest()
		configId string
		groupId  string
	)

	if v, ok := d.GetOk("config_id"); ok {
		configId = v.(string)
		request.ConfigId = helper.String(configId)

	}
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(groupId)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ApplyConfigToMachineGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	instanceId := helper.IdFormat(configId, groupId)
	d.SetId(instanceId)

	return nil
}

func resourceTencentCloudClsConfigAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config_attachment.read")()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	configId := idSplit[0]
	groupId := idSplit[1]

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	machineGroup, err := service.DescribeClsMachineGroupByConfigId(ctx, configId, groupId)
	if err != nil {
		return err
	}

	if machineGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsConfigAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("config_id", configId)
	_ = d.Set("group_id", machineGroup.GroupId)

	return nil
}

func resourceTencentCloudClsConfigAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config_attachment.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := cls.NewDeleteConfigFromMachineGroupRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	configId := idSplit[0]
	groupId := idSplit[1]
	request.GroupId = helper.String(groupId)
	request.ConfigId = helper.String(configId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DeleteConfigFromMachineGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
