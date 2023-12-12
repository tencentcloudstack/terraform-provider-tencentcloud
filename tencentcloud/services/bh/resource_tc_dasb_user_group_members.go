package bh

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudDasbUserGroupMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbUserGroupMembersCreate,
		Read:   resourceTencentCloudDasbUserGroupMembersRead,
		Delete: resourceTencentCloudDasbUserGroupMembersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "User Group ID.",
			},
			"member_id_set": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Collection of member user IDs.",
			},
		},
	}
}

func resourceTencentCloudDasbUserGroupMembersCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user_group_members.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		request        = dasb.NewAddUserGroupMembersRequest()
		userGroupId    string
		memberIdSetStr string
	)

	if v, ok := d.GetOkExists("user_group_id"); ok {
		request.Id = helper.IntUint64(v.(int))
		userGroupIdInt := v.(int)
		userGroupId = strconv.Itoa(userGroupIdInt)
	}

	if v, ok := d.GetOk("member_id_set"); ok {
		memberIdSetSet := v.(*schema.Set).List()
		tmpList := make([]string, 0)
		for i := range memberIdSetSet {
			memberIdSet := memberIdSetSet[i].(int)
			request.MemberIdSet = append(request.MemberIdSet, helper.IntUint64(memberIdSet))
			tmpList = append(tmpList, strconv.Itoa(memberIdSet))
		}

		memberIdSetStr = strings.Join(tmpList, tccommon.COMMA_SP)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().AddUserGroupMembers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("dasb UserGroupMembers not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb UserGroupMembers failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{userGroupId, memberIdSetStr}, tccommon.FILED_SP))

	return resourceTencentCloudDasbUserGroupMembersRead(d, meta)
}

func resourceTencentCloudDasbUserGroupMembersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user_group_members.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	userGroupId := idSplit[0]

	UserGroupMembers, err := service.DescribeDasbUserGroupMembersById(ctx, userGroupId)
	if err != nil {
		return err
	}

	if UserGroupMembers == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DasbUserGroupMembers` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("user_group_id", userGroupId)
	_ = d.Set("member_id_set", UserGroupMembers)

	return nil
}

func resourceTencentCloudDasbUserGroupMembersDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_user_group_members.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	userGroupId := idSplit[0]
	memberIdSetStr := idSplit[1]

	if err := service.DeleteDasbUserGroupMembersById(ctx, userGroupId, memberIdSetStr); err != nil {
		return err
	}

	return nil
}
