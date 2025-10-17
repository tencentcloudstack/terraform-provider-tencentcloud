package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataProjectMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataProjectMemberCreate,
		Read:   resourceTencentCloudWedataProjectMemberRead,
		Update: resourceTencentCloudWedataProjectMemberUpdate,
		Delete: resourceTencentCloudWedataProjectMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID.",
			},

			"user_uin": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User ID.",
			},

			"role_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Role ID.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudWedataProjectMemberCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_project_member.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewCreateProjectMemberRequest()
		projectId string
		userUin   string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("user_uin"); ok {
		request.UserUins = append(request.UserUins, helper.String(v.(string)))
		userUin = v.(string)
	}

	if v, ok := d.GetOk("role_ids"); ok {
		roleIdsSet := v.(*schema.Set).List()
		for i := range roleIdsSet {
			roleIds := roleIdsSet[i].(string)
			request.RoleIds = append(request.RoleIds, helper.String(roleIds))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateProjectMemberWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata project member failed, Response is nil."))
		}

		if !*result.Response.Data.Status {
			return resource.NonRetryableError(fmt.Errorf("Create wedata project member failed, Status is false"))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata project member failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{projectId, userUin}, tccommon.FILED_SP))
	return resourceTencentCloudWedataProjectMemberRead(d, meta)
}

func resourceTencentCloudWedataProjectMemberRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_project_member.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	userUin := idSplit[1]

	respData, err := service.DescribeWedataProjectMemberById(ctx, projectId, userUin)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_project_member` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("user_uin", userUin)
	roleList := make([]interface{}, 0)
	for _, items := range respData {
		if items.Roles != nil {
			for _, roles := range items.Roles {
				if roles.RoleId != nil {
					roleList = append(roleList, roles.RoleId)
				}
			}
		}
	}

	_ = d.Set("role_ids", roleList)

	return nil
}

func resourceTencentCloudWedataProjectMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_project_member.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	userUin := idSplit[1]

	if d.HasChange("role_ids") {
		oldInterface, newInterface := d.GetChange("role_ids")
		oldInstances := oldInterface.(*schema.Set)
		newInstances := newInterface.(*schema.Set)
		remove := oldInstances.Difference(newInstances).List()
		add := newInstances.Difference(oldInstances).List()

		if len(add) > 0 {
			request := wedatav20250806.NewGrantMemberProjectRoleRequest()
			request.ProjectId = &projectId
			request.UserUin = &userUin
			for _, item := range add {
				request.RoleIds = append(request.RoleIds, helper.String(item.(string)))
			}

			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().GrantMemberProjectRoleWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update wedata project member failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}

		if len(remove) > 0 {
			request := wedatav20250806.NewRemoveMemberProjectRoleRequest()
			request.ProjectId = &projectId
			request.UserUin = &userUin
			for _, item := range remove {
				request.RoleIds = append(request.RoleIds, helper.String(item.(string)))
			}

			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().RemoveMemberProjectRoleWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update wedata project member failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	}

	return resourceTencentCloudWedataProjectMemberRead(d, meta)
}

func resourceTencentCloudWedataProjectMemberDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_project_member.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wedatav20250806.NewDeleteProjectMemberRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	userUin := idSplit[1]

	request.ProjectId = &projectId
	request.UserUins = append(request.UserUins, &userUin)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteProjectMemberWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata project member failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
