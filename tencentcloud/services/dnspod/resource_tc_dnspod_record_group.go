package dnspod

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDnspodRecordGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodRecordGroupCreate,
		Read:   resourceTencentCloudDnspodRecordGroupRead,
		Update: resourceTencentCloudDnspodRecordGroupUpdate,
		Delete: resourceTencentCloudDnspodRecordGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Record Group Name.",
			},

			"group_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Group ID.",
			},
		},
	}
}

func resourceTencentCloudDnspodRecordGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_record_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = dnspod.NewCreateRecordGroupRequest()
		response = dnspod.NewCreateRecordGroupResponse()
		groupId  uint64
		domain   string
	)
	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().CreateRecordGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dnspod record_group failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.GroupId
	d.SetId(strings.Join([]string{domain, helper.UInt64ToStr(groupId)}, tccommon.FILED_SP))

	return resourceTencentCloudDnspodRecordGroupRead(d, meta)
}

func resourceTencentCloudDnspodRecordGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_record_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_record_group id is broken, id is %s", d.Id())
	}
	domain := idSplit[0]
	groupId := helper.StrToUInt64(idSplit[1])

	recordGroup, err := service.DescribeDnspodRecordGroupById(ctx, domain, groupId)
	if err != nil {
		return err
	}

	if recordGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DnspodRecordGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)

	if recordGroup.GroupName != nil {
		_ = d.Set("group_name", recordGroup.GroupName)
	}

	if recordGroup.GroupId != nil {
		_ = d.Set("group_id", recordGroup.GroupId)
	}

	return nil
}

func resourceTencentCloudDnspodRecordGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_record_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := dnspod.NewModifyRecordGroupRequest()
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_record_group id is broken, id is %s", d.Id())
	}
	request.Domain = helper.String(idSplit[0])
	request.GroupId = helper.StrToUint64Point(idSplit[1])

	immutableArgs := []string{"domain"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyRecordGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dnspod record_group failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDnspodRecordGroupRead(d, meta)
}

func resourceTencentCloudDnspodRecordGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_record_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_record_group id is broken, id is %s", d.Id())
	}
	domain := idSplit[0]
	groupId := helper.StrToUInt64(idSplit[1])

	if err := service.DeleteDnspodRecordGroupById(ctx, domain, groupId); err != nil {
		return err
	}

	return nil
}
