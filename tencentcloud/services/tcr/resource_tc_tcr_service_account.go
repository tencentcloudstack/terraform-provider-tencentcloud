package tcr

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcrServiceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrServiceAccountCreate,
		Read:   resourceTencentCloudTcrServiceAccountRead,
		Update: resourceTencentCloudTcrServiceAccountUpdate,
		Delete: resourceTencentCloudTcrServiceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service account name.",
			},

			"permissions": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "strategy list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "resource path, currently only supports Namespace. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"actions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Actions, currently support: `tcr:PushRepository`, `tcr:PullRepository`, `tcr:CreateRepository`, `tcr:CreateHelmChart`, `tcr:DescribeHelmCharts`. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service account description.",
			},

			"duration": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "expiration date (unit: day), calculated from the current time, priority is higher than ExpiresAt Service account description.",
			},

			"expires_at": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Service account expiration time (time stamp, unit: milliseconds).",
			},

			"disable": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "whether to disable Service accounts.",
			},

			"password": {
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "Password of the service account.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTcrServiceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_service_account.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = tcr.NewCreateServiceAccountRequest()
		response   = tcr.NewCreateServiceAccountResponse()
		registryId string
		name       string
	)

	if v, ok := d.GetOk("registry_id"); ok {
		request.RegistryId = helper.String(v.(string))
		registryId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
		name = v.(string)
	}

	if v, ok := d.GetOk("permissions"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			permission := tcr.Permission{}
			if v, ok := dMap["resource"]; ok {
				permission.Resource = helper.String(v.(string))
			}

			if v, ok := dMap["actions"]; ok {
				actionsSet := v.(*schema.Set).List()
				for i := range actionsSet {
					if actionsSet[i] != nil {
						actions := actionsSet[i].(string)
						permission.Actions = append(permission.Actions, &actions)
					}
				}
			}

			request.Permissions = append(request.Permissions, &permission)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("duration"); ok {
		request.Duration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("expires_at"); ok {
		request.ExpiresAt = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("disable"); ok {
		request.Disable = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTCRClient().CreateServiceAccount(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tcr ServiceAccount failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tcr ServiceAccount failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Name == nil {
		return fmt.Errorf("Name is nil.")
	}

	if !strings.Contains(*response.Response.Name, name) {
		return fmt.Errorf("The name[%s] return from response is not equal to the name[%s] of tf code.", *response.Response.Name, name)
	}

	d.SetId(strings.Join([]string{registryId, name}, tccommon.FILED_SP))

	var deafultPwd string
	if response.Response.Password != nil {
		deafultPwd = *response.Response.Password
	}

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::tcr:%s:uin/:instance/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	// set custom password OR set default password
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		service := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		password, err := service.ModifyServiceAccountPassword(ctx, registryId, name, v.(string))
		if err != nil {
			return err
		}

		_ = d.Set("password", password)
	} else {
		_ = d.Set("password", deafultPwd)
	}

	return resourceTencentCloudTcrServiceAccountRead(d, meta)
}

func resourceTencentCloudTcrServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_service_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	registryId := idSplit[0]
	name := idSplit[1]

	ServiceAccount, err := service.DescribeTcrServiceAccountById(ctx, registryId, name)
	if err != nil {
		return err
	}

	if ServiceAccount == nil {
		log.Printf("[WARN]%s resource `tencentcloud_tcr_service_account` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("registry_id", registryId)
	_ = d.Set("name", name)

	if ServiceAccount.Permissions != nil {
		permissionsList := []interface{}{}
		for _, permission := range ServiceAccount.Permissions {
			permissionsMap := map[string]interface{}{}
			if permission.Resource != nil {
				permissionsMap["resource"] = permission.Resource
			}

			if len(permission.Actions) > 0 {
				permissionsMap["actions"] = helper.StringsInterfaces(permission.Actions)
			}

			permissionsList = append(permissionsList, permissionsMap)
		}

		_ = d.Set("permissions", permissionsList)
	}

	if ServiceAccount.Description != nil {
		_ = d.Set("description", ServiceAccount.Description)
	}

	if ServiceAccount.ExpiresAt != nil {
		_ = d.Set("expires_at", ServiceAccount.ExpiresAt)
	}

	if ServiceAccount.Disable != nil {
		_ = d.Set("disable", ServiceAccount.Disable)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "tcr", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTcrServiceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_service_account.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	registryId := idSplit[0]
	name := idSplit[1]

	immutableArgs := []string{"registry_id", "name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	needChange := false
	mutableArgs := []string{"permissions", "description", "duration", "expires_at", "disable"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := tcr.NewModifyServiceAccountRequest()
		if v, ok := d.GetOk("permissions"); ok {
			for _, item := range v.([]interface{}) {
				permission := tcr.Permission{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["resource"]; ok {
					permission.Resource = helper.String(v.(string))
				}

				if v, ok := dMap["actions"]; ok {
					actionsSet := v.(*schema.Set).List()
					for i := range actionsSet {
						if actionsSet[i] != nil {
							actions := actionsSet[i].(string)
							permission.Actions = append(permission.Actions, &actions)
						}
					}
				}
				request.Permissions = append(request.Permissions, &permission)
			}
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("duration"); ok {
			request.Duration = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("expires_at"); ok {
			request.ExpiresAt = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("disable"); ok {
			request.Disable = helper.Bool(v.(bool))
		}

		request.RegistryId = &registryId
		request.Name = helper.String(TCR_NAME_PREFIX + name)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTCRClient().ModifyServiceAccount(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update tcr ServiceAccount failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("tcr", "instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	if d.HasChange("password") {
		if v, ok := d.GetOk("password"); ok {
			password, err := service.ModifyServiceAccountPassword(ctx, registryId, name, v.(string))
			if err != nil {
				return err
			}

			_ = d.Set("password", password)
		}
	}

	return resourceTencentCloudTcrServiceAccountRead(d, meta)
}

func resourceTencentCloudTcrServiceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_service_account.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	registryId := idSplit[0]
	name := TCR_NAME_PREFIX + idSplit[1]

	if err := service.DeleteTcrServiceAccountById(ctx, registryId, name); err != nil {
		return err
	}

	return nil
}
