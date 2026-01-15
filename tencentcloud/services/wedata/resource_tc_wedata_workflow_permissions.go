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

func ResourceTencentCloudWedataWorkflowPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataWorkflowPermissionsCreate,
		Read:   resourceTencentCloudWedataWorkflowPermissionsRead,
		Update: resourceTencentCloudWedataWorkflowPermissionsUpdate,
		Delete: resourceTencentCloudWedataWorkflowPermissionsDelete,
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

			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Authorization entity ID.",
			},

			"entity_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Authorization entity type, folder/workflow.",
			},

			"permission_list": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Authorization information array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_target_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Authorization target type (user: user, role: role).",
						},
						"permission_target_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Authorization target ID array (userId/roleId).",
						},
						"permission_type_list": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Authorization permission type array (CAN_VIEW/CAN_RUN/CAN_EDIT/CAN_MANAGE, currently only supports CAN_MANAGE).",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudWedataWorkflowPermissionsCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow_permissions.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = wedatav20250806.NewCreateWorkflowPermissionsRequest()
		response   = wedatav20250806.NewCreateWorkflowPermissionsResponse()
		projectId  string
		entityId   string
		entityType string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("entity_id"); ok {
		request.EntityId = helper.String(v.(string))
		entityId = v.(string)
	}

	if v, ok := d.GetOk("entity_type"); ok {
		request.EntityType = helper.String(v.(string))
		entityType = v.(string)
	}

	if v, ok := d.GetOk("permission_list"); ok {
		for _, item := range v.(*schema.Set).List() {
			permissionListMap := item.(map[string]interface{})
			workflowPermission := wedatav20250806.WorkflowPermission{}
			if v, ok := permissionListMap["permission_target_type"].(string); ok && v != "" {
				workflowPermission.PermissionTargetType = helper.String(v)
			}

			if v, ok := permissionListMap["permission_target_id"].(string); ok && v != "" {
				workflowPermission.PermissionTargetId = helper.String(v)
			}

			if v, ok := permissionListMap["permission_type_list"]; ok {
				permissionTypeListSet := v.(*schema.Set).List()
				for i := range permissionTypeListSet {
					permissionTypeList := permissionTypeListSet[i].(string)
					workflowPermission.PermissionTypeList = append(workflowPermission.PermissionTypeList, helper.String(permissionTypeList))
				}
			}

			request.PermissionList = append(request.PermissionList, &workflowPermission)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateWorkflowPermissionsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata workflow permissions failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata workflow permissions failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.Status == nil || !*response.Response.Data.Status {
		return fmt.Errorf("Create wedata workflow permissions failed, Status is false")
	}

	d.SetId(strings.Join([]string{projectId, entityId, entityType}, tccommon.FILED_SP))
	return resourceTencentCloudWedataWorkflowPermissionsRead(d, meta)
}

func resourceTencentCloudWedataWorkflowPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow_permissions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	entityId := idSplit[1]
	entityType := idSplit[2]

	respData, err := service.DescribeWedataWorkflowPermissionsById(ctx, projectId, entityId, entityType)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_workflow_permissions` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("entity_id", entityId)
	_ = d.Set("entity_type", entityType)

	itemsList := make([]map[string]interface{}, 0, len(respData))
	for _, items := range respData {
		itemsMap := map[string]interface{}{}
		if items.PermissionTargetType != nil {
			itemsMap["permission_target_type"] = items.PermissionTargetType
		}

		if items.PermissionTargetId != nil {
			itemsMap["permission_target_id"] = items.PermissionTargetId
		}

		if items.PermissionTypeList != nil {
			itemsMap["permission_type_list"] = items.PermissionTypeList
		}

		itemsList = append(itemsList, itemsMap)
	}

	_ = d.Set("permission_list", itemsList)
	return nil
}

func resourceTencentCloudWedataWorkflowPermissionsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow_permissions.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	entityId := idSplit[1]
	entityType := idSplit[2]

	if d.HasChange("permission_list") {
		oldInterface, newInterface := d.GetChange("permission_list")
		olds := oldInterface.(*schema.Set)
		news := newInterface.(*schema.Set)
		remove := olds.Difference(news).List()
		add := news.Difference(olds).List()
		if len(remove) > 0 {
			request := wedatav20250806.NewDeleteWorkflowPermissionsRequest()
			response := wedatav20250806.NewDeleteWorkflowPermissionsResponse()
			for _, item := range remove {
				permissionListMap := item.(map[string]interface{})
				workflowPermission := wedatav20250806.DeleteWorkflowPermission{}
				if v, ok := permissionListMap["permission_target_type"].(string); ok && v != "" {
					workflowPermission.PermissionTargetType = helper.String(v)
				}

				if v, ok := permissionListMap["permission_target_id"].(string); ok && v != "" {
					workflowPermission.PermissionTargetId = helper.String(v)
				}

				if v, ok := permissionListMap["permission_type_list"]; ok {
					permissionTypeListSet := v.(*schema.Set).List()
					for i := range permissionTypeListSet {
						permissionTypeList := permissionTypeListSet[i].(string)
						workflowPermission.PermissionTypeList = append(workflowPermission.PermissionTypeList, helper.String(permissionTypeList))
					}
				}

				request.DeletePermissionList = append(request.DeletePermissionList, &workflowPermission)
			}

			request.ProjectId = &projectId
			request.EntityId = &entityId
			request.EntityType = &entityType
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteWorkflowPermissionsWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.Data == nil {
					return resource.NonRetryableError(fmt.Errorf("Delete wedata workflow permissions failed, Response is nil."))
				}

				response = result
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s delete wedata workflow permissions failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			if response.Response.Data.Status == nil || !*response.Response.Data.Status {
				return fmt.Errorf("Delete wedata workflow permissions failed, Status is false")
			}
		}

		if len(add) > 0 {
			request := wedatav20250806.NewCreateWorkflowPermissionsRequest()
			response := wedatav20250806.NewCreateWorkflowPermissionsResponse()
			for _, item := range add {
				permissionListMap := item.(map[string]interface{})
				workflowPermission := wedatav20250806.WorkflowPermission{}
				if v, ok := permissionListMap["permission_target_type"].(string); ok && v != "" {
					workflowPermission.PermissionTargetType = helper.String(v)
				}

				if v, ok := permissionListMap["permission_target_id"].(string); ok && v != "" {
					workflowPermission.PermissionTargetId = helper.String(v)
				}

				if v, ok := permissionListMap["permission_type_list"]; ok {
					permissionTypeListSet := v.(*schema.Set).List()
					for i := range permissionTypeListSet {
						permissionTypeList := permissionTypeListSet[i].(string)
						workflowPermission.PermissionTypeList = append(workflowPermission.PermissionTypeList, helper.String(permissionTypeList))
					}
				}

				request.PermissionList = append(request.PermissionList, &workflowPermission)
			}

			request.ProjectId = &projectId
			request.EntityId = &entityId
			request.EntityType = &entityType
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateWorkflowPermissionsWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.Data == nil {
					return resource.NonRetryableError(fmt.Errorf("Create wedata workflow permissions failed, Response is nil."))
				}

				response = result
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s create wedata workflow permissions failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			if response.Response.Data.Status == nil || !*response.Response.Data.Status {
				return fmt.Errorf("Create wedata workflow permissions failed, Status is false")
			}
		}
	}

	return resourceTencentCloudWedataWorkflowPermissionsRead(d, meta)
}

func resourceTencentCloudWedataWorkflowPermissionsDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow_permissions.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request  = wedatav20250806.NewDeleteWorkflowPermissionsRequest()
		response = wedatav20250806.NewDeleteWorkflowPermissionsResponse()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	entityId := idSplit[1]
	entityType := idSplit[2]

	// get all permissions
	respData, err := service.DescribeWedataWorkflowPermissionsById(ctx, projectId, entityId, entityType)
	if err != nil {
		return err
	}

	if respData == nil || len(respData) == 0 {
		return nil
	}

	for _, item := range respData {
		workflowPermission := wedatav20250806.DeleteWorkflowPermission{}
		if item.PermissionTargetType != nil {
			workflowPermission.PermissionTargetType = item.PermissionTargetType
		}

		if item.PermissionTargetId != nil {
			workflowPermission.PermissionTargetId = item.PermissionTargetId
		}

		if item.PermissionTypeList != nil {
			workflowPermission.PermissionTypeList = item.PermissionTypeList
		}

		request.DeletePermissionList = append(request.DeletePermissionList, &workflowPermission)
	}

	request.ProjectId = &projectId
	request.EntityId = &entityId
	request.EntityType = &entityType
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteWorkflowPermissionsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata workflow permissions failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata workflow permissions failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.Status == nil || !*response.Response.Data.Status {
		return fmt.Errorf("Delete wedata workflow permissions failed, Status is false")
	}

	return nil
}
