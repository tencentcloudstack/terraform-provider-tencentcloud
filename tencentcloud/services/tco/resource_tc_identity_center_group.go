package tco

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIdentityCenterGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterGroupCreate,
		Read:   resourceTencentCloudIdentityCenterGroupRead,
		Update: resourceTencentCloudIdentityCenterGroupUpdate,
		Delete: resourceTencentCloudIdentityCenterGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone id.",
			},

			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the user group. Format: Allow English letters, numbers and special characters-. Length: Maximum 128 characters.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the user group.",
			},
			"group_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Type of user group. `Manual`: manual creation, `Synchronized`: external import.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time for the user group.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user group.",
			},
			"member_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of team members.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A description of the user group. Length: Maximum 1024 characters.",
			},
		},
	}
}

func resourceTencentCloudIdentityCenterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId  string
		groupId string
	)
	var (
		request  = organization.NewCreateGroupRequest()
		response = organization.NewCreateGroupResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_type"); ok {
		request.GroupType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().CreateGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create identity center group failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.GroupInfo.GroupId

	d.SetId(strings.Join([]string{zoneId, groupId}, tccommon.FILED_SP))

	return resourceTencentCloudIdentityCenterGroupRead(d, meta)
}

func resourceTencentCloudIdentityCenterGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	groupId := idSplit[1]

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeIdentityCenterGroupById(ctx, zoneId, groupId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `identity_center_group` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.GroupName != nil {
		_ = d.Set("group_name", respData.GroupName)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.GroupType != nil {
		_ = d.Set("group_type", respData.GroupType)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	if respData.GroupId != nil {
		_ = d.Set("group_id", respData.GroupId)
	}

	if respData.MemberCount != nil {
		_ = d.Set("member_count", respData.MemberCount)
	}

	return nil
}

func resourceTencentCloudIdentityCenterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"zone_id", "group_type"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	groupId := idSplit[1]

	needChange := false
	mutableArgs := []string{"group_name", "description"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := organization.NewUpdateGroupRequest()

		request.ZoneId = helper.String(zoneId)

		request.GroupId = helper.String(groupId)

		if v, ok := d.GetOk("group_name"); ok {
			request.NewGroupName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			request.NewDescription = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().UpdateGroupWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update identity center group failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudIdentityCenterGroupRead(d, meta)
}

func resourceTencentCloudIdentityCenterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	groupId := idSplit[1]

	var (
		request  = organization.NewDeleteGroupRequest()
		response = organization.NewDeleteGroupResponse()
	)

	request.ZoneId = helper.String(zoneId)

	request.GroupId = helper.String(groupId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().DeleteGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete identity center group failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
