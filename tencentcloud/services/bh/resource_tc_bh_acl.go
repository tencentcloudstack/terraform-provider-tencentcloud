package bh

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhAclCreate,
		Read:   resourceTencentCloudBhAclRead,
		Update: resourceTencentCloudBhAclUpdate,
		Delete: resourceTencentCloudBhAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access permission name, maximum 32 characters, cannot contain whitespace characters.",
			},

			"allow_disk_redirect": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable disk mapping.",
			},

			"allow_any_account": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to allow any account to log in.",
			},

			"allow_clip_file_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable clipboard file upload.",
			},

			"allow_clip_file_down": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable clipboard file download.",
			},

			"allow_clip_text_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable clipboard text (including images) upload.",
			},

			"allow_clip_text_down": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable clipboard text (including images) download.",
			},

			"allow_file_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable SFTP file upload.",
			},

			"max_file_up_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "File transfer upload size limit (reserved parameter, not currently used).",
			},

			"allow_file_down": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable SFTP file download.",
			},

			"max_file_down_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "File transfer download size limit (reserved parameter, not currently used).",
			},

			"allow_disk_file_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable RDP disk mapping file upload.",
			},

			"allow_disk_file_down": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable RDP disk mapping file download.",
			},

			"allow_shell_file_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable rz sz file upload.",
			},

			"allow_shell_file_down": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable rz sz file download.",
			},

			"allow_file_del": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable SFTP file deletion.",
			},

			"allow_access_credential": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to allow the use of access credentials. Default is allowed.",
			},

			"allow_keyboard_logger": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to allow keyboard logging.",
			},

			"max_access_credential_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum validity period of access credentials (in seconds). Must be a multiple of 86400 when access credentials are enabled.",
			},

			"user_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associated user ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"user_group_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associated user group ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"device_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associated asset ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"app_asset_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associated application asset ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"device_group_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associated asset group ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"account_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associated account set.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"cmd_template_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associated high-risk command template ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"ac_template_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associated high-risk DB template ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"validate_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Access permission effective time in ISO8601 format, e.g.: `2021-09-22T00:00:00+00:00`. If not set, the permission is permanently valid.",
			},

			"validate_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Access permission expiration time in ISO8601 format, e.g.: `2021-09-23T00:00:00+00:00`. If not set, the permission is permanently valid.",
			},

			"department_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Department ID to which the access permission belongs, e.g.: `1.2.3`.",
			},

			// computed
			"acl_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Access permission ID.",
			},
		},
	}
}

func resourceTencentCloudBhAclCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_acl.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = bhv20230418.NewCreateAclRequest()
		response = bhv20230418.NewCreateAclResponse()
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("allow_disk_redirect"); ok {
		request.AllowDiskRedirect = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_any_account"); ok {
		request.AllowAnyAccount = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_clip_file_up"); ok {
		request.AllowClipFileUp = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_clip_file_down"); ok {
		request.AllowClipFileDown = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_clip_text_up"); ok {
		request.AllowClipTextUp = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_clip_text_down"); ok {
		request.AllowClipTextDown = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_file_up"); ok {
		request.AllowFileUp = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("max_file_up_size"); ok {
		request.MaxFileUpSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("allow_file_down"); ok {
		request.AllowFileDown = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("max_file_down_size"); ok {
		request.MaxFileDownSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("allow_disk_file_up"); ok {
		request.AllowDiskFileUp = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_disk_file_down"); ok {
		request.AllowDiskFileDown = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_shell_file_up"); ok {
		request.AllowShellFileUp = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_shell_file_down"); ok {
		request.AllowShellFileDown = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_file_del"); ok {
		request.AllowFileDel = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_access_credential"); ok {
		request.AllowAccessCredential = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("allow_keyboard_logger"); ok {
		request.AllowKeyboardLogger = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("max_access_credential_duration"); ok {
		request.MaxAccessCredentialDuration = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("user_id_set"); ok {
		for _, id := range v.(*schema.Set).List() {
			request.UserIdSet = append(request.UserIdSet, helper.IntUint64(id.(int)))
		}
	}

	if v, ok := d.GetOk("user_group_id_set"); ok {
		for _, id := range v.(*schema.Set).List() {
			request.UserGroupIdSet = append(request.UserGroupIdSet, helper.IntUint64(id.(int)))
		}
	}

	if v, ok := d.GetOk("device_id_set"); ok {
		for _, id := range v.(*schema.Set).List() {
			request.DeviceIdSet = append(request.DeviceIdSet, helper.IntUint64(id.(int)))
		}
	}

	if v, ok := d.GetOk("app_asset_id_set"); ok {
		for _, id := range v.(*schema.Set).List() {
			request.AppAssetIdSet = append(request.AppAssetIdSet, helper.IntUint64(id.(int)))
		}
	}

	if v, ok := d.GetOk("device_group_id_set"); ok {
		for _, id := range v.(*schema.Set).List() {
			request.DeviceGroupIdSet = append(request.DeviceGroupIdSet, helper.IntUint64(id.(int)))
		}
	}

	if v, ok := d.GetOk("account_set"); ok {
		for _, acc := range v.(*schema.Set).List() {
			request.AccountSet = append(request.AccountSet, helper.String(acc.(string)))
		}
	}

	if v, ok := d.GetOk("cmd_template_id_set"); ok {
		for _, id := range v.(*schema.Set).List() {
			request.CmdTemplateIdSet = append(request.CmdTemplateIdSet, helper.IntUint64(id.(int)))
		}
	}

	if v, ok := d.GetOk("ac_template_id_set"); ok {
		for _, id := range v.(*schema.Set).List() {
			request.ACTemplateIdSet = append(request.ACTemplateIdSet, helper.String(id.(string)))
		}
	}

	if v, ok := d.GetOk("validate_from"); ok {
		request.ValidateFrom = helper.String(v.(string))
	}

	if v, ok := d.GetOk("validate_to"); ok {
		request.ValidateTo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		request.DepartmentId = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().CreateAclWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bh acl failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bh acl failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Id == nil {
		return fmt.Errorf("acl Id is nil.")
	}

	d.SetId(helper.UInt64ToStr(*response.Response.Id))
	return resourceTencentCloudBhAclRead(d, meta)
}

func resourceTencentCloudBhAclRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_acl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	aclId := d.Id()

	respData, err := service.DescribeBhAclById(ctx, aclId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_acl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	departments, err := service.DescribeBhDepartments(ctx)
	if err != nil {
		return err
	}

	if respData.Id != nil {
		_ = d.Set("acl_id", respData.Id)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.AllowDiskRedirect != nil {
		_ = d.Set("allow_disk_redirect", respData.AllowDiskRedirect)
	}

	if respData.AllowAnyAccount != nil {
		_ = d.Set("allow_any_account", respData.AllowAnyAccount)
	}

	if respData.AllowClipFileUp != nil {
		_ = d.Set("allow_clip_file_up", respData.AllowClipFileUp)
	}

	if respData.AllowClipFileDown != nil {
		_ = d.Set("allow_clip_file_down", respData.AllowClipFileDown)
	}

	if respData.AllowClipTextUp != nil {
		_ = d.Set("allow_clip_text_up", respData.AllowClipTextUp)
	}

	if respData.AllowClipTextDown != nil {
		_ = d.Set("allow_clip_text_down", respData.AllowClipTextDown)
	}

	if respData.AllowFileUp != nil {
		_ = d.Set("allow_file_up", respData.AllowFileUp)
	}

	if respData.MaxFileUpSize != nil {
		_ = d.Set("max_file_up_size", respData.MaxFileUpSize)
	}

	if respData.AllowFileDown != nil {
		_ = d.Set("allow_file_down", respData.AllowFileDown)
	}

	if respData.MaxFileDownSize != nil {
		_ = d.Set("max_file_down_size", respData.MaxFileDownSize)
	}

	if respData.AllowDiskFileUp != nil {
		_ = d.Set("allow_disk_file_up", respData.AllowDiskFileUp)
	}

	if respData.AllowDiskFileDown != nil {
		_ = d.Set("allow_disk_file_down", respData.AllowDiskFileDown)
	}

	if respData.AllowShellFileUp != nil {
		_ = d.Set("allow_shell_file_up", respData.AllowShellFileUp)
	}

	if respData.AllowShellFileDown != nil {
		_ = d.Set("allow_shell_file_down", respData.AllowShellFileDown)
	}

	if respData.AllowFileDel != nil {
		_ = d.Set("allow_file_del", respData.AllowFileDel)
	}

	if respData.AllowAccessCredential != nil {
		_ = d.Set("allow_access_credential", respData.AllowAccessCredential)
	}

	if respData.AllowKeyboardLogger != nil {
		_ = d.Set("allow_keyboard_logger", respData.AllowKeyboardLogger)
	}

	if respData.ValidateFrom != nil {
		_ = d.Set("validate_from", respData.ValidateFrom)
	}

	if respData.ValidateTo != nil {
		_ = d.Set("validate_to", respData.ValidateTo)
	}

	if departments != nil && departments.Enabled != nil && *departments.Enabled {
		if respData.Department != nil && respData.Department.Id != nil {
			_ = d.Set("department_id", respData.Department.Id)
		}
	}

	if respData.UserSet != nil && len(respData.UserSet) > 0 {
		userIdSet := make([]interface{}, 0, len(respData.UserSet))
		for _, user := range respData.UserSet {
			if user.Id != nil {
				userIdSet = append(userIdSet, int(*user.Id))
			}
		}
		_ = d.Set("user_id_set", userIdSet)
	}

	if respData.UserGroupSet != nil && len(respData.UserGroupSet) > 0 {
		userGroupIdSet := make([]interface{}, 0, len(respData.UserGroupSet))
		for _, group := range respData.UserGroupSet {
			if group.Id != nil {
				userGroupIdSet = append(userGroupIdSet, int(*group.Id))
			}
		}
		_ = d.Set("user_group_id_set", userGroupIdSet)
	}

	if respData.DeviceSet != nil && len(respData.DeviceSet) > 0 {
		deviceIdSet := make([]interface{}, 0, len(respData.DeviceSet))
		for _, device := range respData.DeviceSet {
			if device.Id != nil {
				deviceIdSet = append(deviceIdSet, int(*device.Id))
			}
		}
		_ = d.Set("device_id_set", deviceIdSet)
	}

	if respData.DeviceGroupSet != nil && len(respData.DeviceGroupSet) > 0 {
		deviceGroupIdSet := make([]interface{}, 0, len(respData.DeviceGroupSet))
		for _, group := range respData.DeviceGroupSet {
			if group.Id != nil {
				deviceGroupIdSet = append(deviceGroupIdSet, int(*group.Id))
			}
		}
		_ = d.Set("device_group_id_set", deviceGroupIdSet)
	}

	if respData.AppAssetSet != nil && len(respData.AppAssetSet) > 0 {
		appAssetIdSet := make([]interface{}, 0, len(respData.AppAssetSet))
		for _, asset := range respData.AppAssetSet {
			if asset.Id != nil {
				appAssetIdSet = append(appAssetIdSet, int(*asset.Id))
			}
		}
		_ = d.Set("app_asset_id_set", appAssetIdSet)
	}

	if respData.AccountSet != nil && len(respData.AccountSet) > 0 {
		accountSet := make([]interface{}, 0, len(respData.AccountSet))
		for _, acc := range respData.AccountSet {
			if acc != nil {
				accountSet = append(accountSet, *acc)
			}
		}
		_ = d.Set("account_set", accountSet)
	}

	if respData.CmdTemplateSet != nil && len(respData.CmdTemplateSet) > 0 {
		cmdTemplateIdSet := make([]interface{}, 0, len(respData.CmdTemplateSet))
		for _, tmpl := range respData.CmdTemplateSet {
			if tmpl.Id != nil {
				cmdTemplateIdSet = append(cmdTemplateIdSet, int(*tmpl.Id))
			}
		}
		_ = d.Set("cmd_template_id_set", cmdTemplateIdSet)
	}

	if respData.ACTemplateSet != nil && len(respData.ACTemplateSet) > 0 {
		acTemplateIdSet := make([]interface{}, 0, len(respData.ACTemplateSet))
		for _, tmpl := range respData.ACTemplateSet {
			if tmpl.TemplateId != nil {
				acTemplateIdSet = append(acTemplateIdSet, *tmpl.TemplateId)
			}
		}
		_ = d.Set("ac_template_id_set", acTemplateIdSet)
	}

	return nil
}

func resourceTencentCloudBhAclUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_acl.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	needChange := false
	mutableArgs := []string{
		"name", "allow_disk_redirect", "allow_any_account",
		"allow_clip_file_up", "allow_clip_file_down", "allow_clip_text_up", "allow_clip_text_down",
		"allow_file_up", "max_file_up_size", "allow_file_down", "max_file_down_size",
		"allow_disk_file_up", "allow_disk_file_down", "allow_shell_file_up", "allow_shell_file_down",
		"allow_file_del", "allow_access_credential", "allow_keyboard_logger", "max_access_credential_duration",
		"user_id_set", "user_group_id_set", "device_id_set", "app_asset_id_set", "device_group_id_set",
		"account_set", "cmd_template_id_set", "ac_template_id_set",
		"validate_from", "validate_to", "department_id",
	}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := bhv20230418.NewModifyAclRequest()
		request.Id = helper.StrToUint64Point(d.Id())

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("allow_disk_redirect"); ok {
			request.AllowDiskRedirect = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_any_account"); ok {
			request.AllowAnyAccount = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_clip_file_up"); ok {
			request.AllowClipFileUp = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_clip_file_down"); ok {
			request.AllowClipFileDown = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_clip_text_up"); ok {
			request.AllowClipTextUp = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_clip_text_down"); ok {
			request.AllowClipTextDown = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_file_up"); ok {
			request.AllowFileUp = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("max_file_up_size"); ok {
			request.MaxFileUpSize = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("allow_file_down"); ok {
			request.AllowFileDown = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("max_file_down_size"); ok {
			request.MaxFileDownSize = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("allow_disk_file_up"); ok {
			request.AllowDiskFileUp = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_disk_file_down"); ok {
			request.AllowDiskFileDown = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_shell_file_up"); ok {
			request.AllowShellFileUp = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_shell_file_down"); ok {
			request.AllowShellFileDown = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_file_del"); ok {
			request.AllowFileDel = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_access_credential"); ok {
			request.AllowAccessCredential = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("allow_keyboard_logger"); ok {
			request.AllowKeyboardLogger = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("max_access_credential_duration"); ok {
			request.MaxAccessCredentialDuration = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("user_id_set"); ok {
			for _, id := range v.(*schema.Set).List() {
				request.UserIdSet = append(request.UserIdSet, helper.IntUint64(id.(int)))
			}
		}

		if v, ok := d.GetOk("user_group_id_set"); ok {
			for _, id := range v.(*schema.Set).List() {
				request.UserGroupIdSet = append(request.UserGroupIdSet, helper.IntUint64(id.(int)))
			}
		}

		if v, ok := d.GetOk("device_id_set"); ok {
			for _, id := range v.(*schema.Set).List() {
				request.DeviceIdSet = append(request.DeviceIdSet, helper.IntUint64(id.(int)))
			}
		}

		if v, ok := d.GetOk("app_asset_id_set"); ok {
			for _, id := range v.(*schema.Set).List() {
				request.AppAssetIdSet = append(request.AppAssetIdSet, helper.IntUint64(id.(int)))
			}
		}

		if v, ok := d.GetOk("device_group_id_set"); ok {
			for _, id := range v.(*schema.Set).List() {
				request.DeviceGroupIdSet = append(request.DeviceGroupIdSet, helper.IntUint64(id.(int)))
			}
		}

		if v, ok := d.GetOk("account_set"); ok {
			for _, acc := range v.(*schema.Set).List() {
				request.AccountSet = append(request.AccountSet, helper.String(acc.(string)))
			}
		}

		if v, ok := d.GetOk("cmd_template_id_set"); ok {
			for _, id := range v.(*schema.Set).List() {
				request.CmdTemplateIdSet = append(request.CmdTemplateIdSet, helper.IntUint64(id.(int)))
			}
		}

		if v, ok := d.GetOk("ac_template_id_set"); ok {
			for _, id := range v.(*schema.Set).List() {
				request.ACTemplateIdSet = append(request.ACTemplateIdSet, helper.String(id.(string)))
			}
		}

		if v, ok := d.GetOk("validate_from"); ok {
			request.ValidateFrom = helper.String(v.(string))
		}

		if v, ok := d.GetOk("validate_to"); ok {
			request.ValidateTo = helper.String(v.(string))
		}

		if v, ok := d.GetOk("department_id"); ok {
			request.DepartmentId = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyAclWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bh acl failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBhAclRead(d, meta)
}

func resourceTencentCloudBhAclDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_acl.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewDeleteAclsRequest()
	)

	request.IdSet = []*uint64{helper.StrToUint64Point(d.Id())}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().DeleteAclsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bh acl failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
