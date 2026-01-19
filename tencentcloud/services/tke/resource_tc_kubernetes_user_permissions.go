package tke

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const maxPermissionsPerRequest = 100

func ResourceTencentCloudKubernetesUserPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesUserPermissionsCreate,
		Read:   resourceTencentCloudKubernetesUserPermissionsRead,
		Update: resourceTencentCloudKubernetesUserPermissionsUpdate,
		Delete: resourceTencentCloudKubernetesUserPermissionsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_uin": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier of the user to be authorized (supports sub-account UIN and role UIN).",
			},

			"permissions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Complete list of permissions that the user should ultimately have. Uses declarative semantics, the passed list represents all permissions the user should ultimately have, the system will automatically calculate differences and perform necessary create/delete operations. When empty or not provided, all permissions for this user will be cleared. Maximum support for 100 permission items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cluster ID.",
						},
						"role_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Role name. Predefined roles include: tke:admin (cluster administrator), tke:ops (operations personnel), tke:dev (developer), tke:ro (read-only user), tke:ns:dev (namespace developer), tke:ns:ro (namespace read-only user), others are user-defined roles.",
						},
						"role_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Authorization type. Enum values: cluster (cluster-level permissions, corresponding to ClusterRoleBinding), namespace (namespace-level permissions, corresponding to RoleBinding).",
						},
						"is_custom": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether it is a custom role, default false.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Namespace. Required when RoleType is namespace.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudKubernetesUserPermissionsCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_user_permissions.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		targetUin string
	)

	if v, ok := d.GetOk("target_uin"); ok {
		targetUin = v.(string)
	}

	var allPermissions []*tkev20180525.PermissionItem
	if v, ok := d.GetOk("permissions"); ok {
		for _, item := range v.(*schema.Set).List() {
			permissionsMap := item.(map[string]interface{})
			permissionItem := tkev20180525.PermissionItem{}
			if v, ok := permissionsMap["cluster_id"].(string); ok && v != "" {
				permissionItem.ClusterId = helper.String(v)
			}

			if v, ok := permissionsMap["role_name"].(string); ok && v != "" {
				permissionItem.RoleName = helper.String(v)
			}

			if v, ok := permissionsMap["role_type"].(string); ok && v != "" {
				permissionItem.RoleType = helper.String(v)
			}

			if v, ok := permissionsMap["is_custom"].(bool); ok {
				permissionItem.IsCustom = helper.Bool(v)
			}

			if v, ok := permissionsMap["namespace"].(string); ok && v != "" {
				permissionItem.Namespace = helper.String(v)
			}

			allPermissions = append(allPermissions, &permissionItem)
		}
	}

	// max permissions is 100 for once request
	totalPermissions := len(allPermissions)
	for i := 0; i < totalPermissions; i += maxPermissionsPerRequest {
		end := i + maxPermissionsPerRequest
		if end > totalPermissions {
			end = totalPermissions
		}

		batchPermissions := allPermissions[i:end]
		request := tkev20180525.NewGrantUserPermissionsRequest()
		request.TargetUin = &targetUin
		request.Permissions = batchPermissions
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().GrantUserPermissionsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success for batch %d-%d, request body [%s], response body [%s]\n", logId, request.GetAction(), i+1, end, request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s create kubernetes user permissions failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	d.SetId(targetUin)
	return resourceTencentCloudKubernetesUserPermissionsRead(d, meta)
}

func resourceTencentCloudKubernetesUserPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_user_permissions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		targetUin = d.Id()
	)

	respData, err := service.DescribeKubernetesUserPermissionsById(ctx, targetUin)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_user_permissions` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.TargetUin != nil {
		_ = d.Set("target_uin", respData.TargetUin)
	}

	if respData.Permissions != nil {
		permissionsList := make([]map[string]interface{}, 0, len(respData.Permissions))
		for _, permissions := range respData.Permissions {
			permissionsMap := map[string]interface{}{}
			if permissions.ClusterId != nil {
				permissionsMap["cluster_id"] = permissions.ClusterId
			}

			if permissions.RoleName != nil {
				permissionsMap["role_name"] = permissions.RoleName
			}

			if permissions.RoleType != nil {
				permissionsMap["role_type"] = permissions.RoleType
			}

			if permissions.IsCustom != nil {
				permissionsMap["is_custom"] = permissions.IsCustom
			}

			if permissions.Namespace != nil {
				permissionsMap["namespace"] = permissions.Namespace
			}

			permissionsList = append(permissionsList, permissionsMap)
		}

		_ = d.Set("permissions", permissionsList)
	}

	return nil
}

func resourceTencentCloudKubernetesUserPermissionsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_user_permissions.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudKubernetesUserPermissionsRead(d, meta)
}

func resourceTencentCloudKubernetesUserPermissionsDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_user_permissions.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		targetUin = d.Id()
	)

	// get all permissions
	respData, err := service.DescribeKubernetesUserPermissionsById(ctx, targetUin)
	if err != nil {
		return err
	}

	if respData == nil || respData.TargetUin == nil || respData.Permissions == nil || len(respData.Permissions) == 0 {
		return nil
	}

	// max permissions is 100 for once request
	totalPermissions := len(respData.Permissions)
	for i := 0; i < totalPermissions; i += maxPermissionsPerRequest {
		end := i + maxPermissionsPerRequest
		if end > totalPermissions {
			end = totalPermissions
		}

		batchPermissions := respData.Permissions[i:end]
		request := tkev20180525.NewDeleteUserPermissionsRequest()
		request.TargetUin = &targetUin
		request.Permissions = batchPermissions
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DeleteUserPermissionsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success for batch %d-%d, request body [%s], response body [%s]\n", logId, request.GetAction(), i+1, end, request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s delete kubernetes user permissions failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return nil
}
