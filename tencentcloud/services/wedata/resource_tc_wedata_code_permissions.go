package wedata

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataCodePermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataCodePermissionsCreate,
		Read:   resourceTencentCloudWedataCodePermissionsRead,
		Update: resourceTencentCloudWedataCodePermissionsUpdate,
		Delete: resourceTencentCloudWedataCodePermissionsDelete,
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

			"authorize_permission_objects": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Permission operation objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Authorization resource information, including resourceId and resourceType.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource type, can only be these two types: folder, script.",
									},
									"resource_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource ID: directory ID or script ID.",
									},
									"resource_id_for_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Full ID path, used for recursive authentication.",
									},
									"resource_cfs_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CFS path.",
									},
								},
							},
						},
						"authorize_subjects": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Authorization details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subject_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Subject type (user: user, role: role, group: group).",
									},
									"subject_values": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Subject value list.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"privileges": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Permission list.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudWedataCodePermissionsCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_permissions.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewCreateCodePermissionsRequest()
		projectId string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("authorize_permission_objects"); ok {
		for _, item := range v.([]interface{}) {
			authorizePermissionObjectsMap := item.(map[string]interface{})
			exploreAuthorizationObject := wedatav20250806.ExploreAuthorizationObject{}
			if resourceMap, ok := helper.ConvertInterfacesHeadToMap(authorizePermissionObjectsMap["resource"]); ok {
				exploreFileResource := wedatav20250806.ExploreFileResource{}
				if v, ok := resourceMap["resource_type"].(string); ok && v != "" {
					exploreFileResource.ResourceType = helper.String(v)
				}

				if v, ok := resourceMap["resource_id"].(string); ok && v != "" {
					exploreFileResource.ResourceId = helper.String(v)
				}

				if v, ok := resourceMap["resource_id_for_path"].(string); ok && v != "" {
					exploreFileResource.ResourceIdForPath = helper.String(v)
				}

				if v, ok := resourceMap["resource_cfs_path"].(string); ok && v != "" {
					exploreFileResource.ResourceCFSPath = helper.String(v)
				}

				exploreAuthorizationObject.Resource = &exploreFileResource
			}

			if v, ok := authorizePermissionObjectsMap["authorize_subjects"]; ok {
				for _, item := range v.([]interface{}) {
					authorizeSubjectsMap := item.(map[string]interface{})
					exploreAuthorizeSubject := wedatav20250806.ExploreAuthorizeSubject{}
					if v, ok := authorizeSubjectsMap["subject_type"].(string); ok && v != "" {
						exploreAuthorizeSubject.SubjectType = helper.String(v)
					}

					if v, ok := authorizeSubjectsMap["subject_values"]; ok {
						subjectValuesSet := v.(*schema.Set).List()
						for i := range subjectValuesSet {
							subjectValues := subjectValuesSet[i].(string)
							exploreAuthorizeSubject.SubjectValues = append(exploreAuthorizeSubject.SubjectValues, helper.String(subjectValues))
						}
					}

					if v, ok := authorizeSubjectsMap["privileges"]; ok {
						privilegesSet := v.(*schema.Set).List()
						for i := range privilegesSet {
							privileges := privilegesSet[i].(string)
							exploreAuthorizeSubject.Privileges = append(exploreAuthorizeSubject.Privileges, helper.String(privileges))
						}
					}

					exploreAuthorizationObject.AuthorizeSubjects = append(exploreAuthorizationObject.AuthorizeSubjects, &exploreAuthorizeSubject)
				}
			}

			request.AuthorizePermissionObjects = append(request.AuthorizePermissionObjects, &exploreAuthorizationObject)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateCodePermissionsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata code permissions failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata code permissions failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(projectId)
	return resourceTencentCloudWedataCodePermissionsRead(d, meta)
}

func resourceTencentCloudWedataCodePermissionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_permissions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		projectId = d.Id()
	)

	respData, err := service.DescribeWedataCodePermissionsById(ctx, projectId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_code_permissions` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("project_id", projectId)

	rowsList := make([]map[string]interface{}, 0, len(respData))
	for _, rows := range respData {
		rowsMap := map[string]interface{}{}
		if rows.Privileges != nil {
			rowsMap["privileges"] = rows.Privileges
		}

		if rows.RoleType != nil {
			rowsMap["role_type"] = rows.RoleType
		}

		if rows.RoleId != nil {
			rowsMap["role_id"] = rows.RoleId
		}

		resourceMap := map[string]interface{}{}
		if rows.Resource != nil {
			if rows.Resource.ResourceType != nil {
				resourceMap["resource_type"] = rows.Resource.ResourceType
			}

			if rows.Resource.ResourceId != nil {
				resourceMap["resource_id"] = rows.Resource.ResourceId
			}

			if rows.Resource.ResourceIdForPath != nil {
				resourceMap["resource_id_for_path"] = rows.Resource.ResourceIdForPath
			}

			if rows.Resource.ResourceCFSPath != nil {
				resourceMap["resource_cfs_path"] = rows.Resource.ResourceCFSPath
			}

			rowsMap["resource"] = []interface{}{resourceMap}
		}

		if rows.DeleteAble != nil {
			rowsMap["delete_able"] = rows.DeleteAble
		}

		rowsList = append(rowsList, rowsMap)
	}

	_ = d.Set("authorize_permission_objects", rowsList)

	return nil
}

func resourceTencentCloudWedataCodePermissionsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_permissions.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	// var (
	// 	logId     = tccommon.GetLogId(tccommon.ContextNil)
	// 	ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	// 	projectId = d.Id()
	// )

	// if d.HasChange("authorize_permission_objects") {
	// 	oldInterface, newInterface := d.GetChange("authorize_permission_objects")
	// 	olds := oldInterface.(*schema.Set)
	// 	news := newInterface.(*schema.Set)
	// 	remove := olds.Difference(news).List()
	// 	add := news.Difference(olds).List()
	// 	if len(remove) > 0 {

	// 	}

	// 	if len(add) > 0 {

	// 	}
	// }

	return resourceTencentCloudWedataCodePermissionsRead(d, meta)
}

func resourceTencentCloudWedataCodePermissionsDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_permissions.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewDeleteCodePermissionsRequest()
		projectId = d.Id()
	)

	request.ProjectId = &projectId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteCodePermissionsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata code permissions failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
