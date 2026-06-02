package cls

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsConsole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsConsoleCreate,
		Read:   resourceTencentCloudClsConsoleRead,
		Update: resourceTencentCloudClsConsoleUpdate,
		Delete: resourceTencentCloudClsConsoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_mode": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Access mode list. Valid values: `public` (public network), `internal` (intranet).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"login_mode": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Login mode. Valid values: `0` (account-password authentication), `1` (anonymous login), `2` (third-party authentication login).",
			},

			"domain_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom domain prefix.",
			},

			"accounts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "User account information. Required when `login_mode` is `0` (account-password authentication).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User name.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "User password.",
						},
						"secret_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Tencent Cloud account SecretId.",
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Tencent Cloud account SecretKey.",
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Email address used to send verification codes.",
						},
					},
				},
			},

			"anonymous_login": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Anonymous login account information. Required when `login_mode` is `1` (anonymous login).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Anonymous login account SecretId.",
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Anonymous login account SecretKey.",
						},
					},
				},
			},

			"intranet_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Intranet type. Default is 0.",
			},

			"intranet_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Intranet region.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Intranet VPC ID.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Intranet subnet ID.",
			},

			"auth_roles": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Auth role information. Required when `login_mode` is `2` (third-party authentication login).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Auth role name.",
						},
						"secret_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "SecretId for the auth role permission.",
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "SecretKey for the auth role permission.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Tag bindings. ModifyConsole does not accept tags, so changing this field forces resource replacement.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"hide_params": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Custom hidden parameters.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"access_control_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Access control rules. Required when `login_mode` is `2` (third-party authentication login) and AccessMode contains `internal` with Action ACCEPT rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_blocks": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "CIDR blocks or IPs, supporting IPv4 or IPv6.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule action. Valid values: `ACCEPT`, `DROP`.",
						},
						"access_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access mode for the rule. Valid values: `public`, `internal`.",
						},
					},
				},
			},

			"remarks": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remarks.",
			},

			"menus": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Custom display menus.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// computed
			"console_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DataSight console ID.",
			},

			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public network access domain.",
			},

			"intranet_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intranet access domain.",
			},
		},
	}
}

func resourceTencentCloudClsConsoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_console.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = cls.NewCreateConsoleRequest()
		response  = cls.NewCreateConsoleResponse()
		consoleId string
	)

	if v, ok := d.GetOk("access_mode"); ok {
		accessModeList := v.([]interface{})
		for _, item := range accessModeList {
			request.AccessMode = append(request.AccessMode, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOkExists("login_mode"); ok {
		request.LoginMode = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("domain_prefix"); ok {
		request.DomainPrefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("accounts"); ok {
		accountsList := v.([]interface{})
		for _, item := range accountsList {
			account := item.(map[string]interface{})
			consoleAccount := &cls.ConsoleAccount{}
			if v, ok := account["user_name"].(string); ok && v != "" {
				consoleAccount.UserName = helper.String(v)
			}
			if v, ok := account["password"].(string); ok && v != "" {
				consoleAccount.Password = helper.String(v)
			}
			if v, ok := account["secret_id"].(string); ok && v != "" {
				consoleAccount.SecretId = helper.String(v)
			}
			if v, ok := account["secret_key"].(string); ok && v != "" {
				consoleAccount.SecretKey = helper.String(v)
			}
			if v, ok := account["email"].(string); ok && v != "" {
				consoleAccount.Email = helper.String(v)
			}
			request.Accounts = append(request.Accounts, consoleAccount)
		}
	}

	if v, ok := d.GetOk("anonymous_login"); ok {
		anonymousList := v.([]interface{})
		if len(anonymousList) > 0 && anonymousList[0] != nil {
			anonymous := anonymousList[0].(map[string]interface{})
			anonymousLogin := &cls.AnonymousLoginInfo{}
			if v, ok := anonymous["secret_id"].(string); ok && v != "" {
				anonymousLogin.SecretId = helper.String(v)
			}
			if v, ok := anonymous["secret_key"].(string); ok && v != "" {
				anonymousLogin.SecretKey = helper.String(v)
			}
			request.AnonymousLogin = anonymousLogin
		}
	}

	if v, ok := d.GetOkExists("intranet_type"); ok {
		request.IntranetType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("intranet_region"); ok {
		request.IntranetRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_roles"); ok {
		authRolesList := v.([]interface{})
		for _, item := range authRolesList {
			role := item.(map[string]interface{})
			authRole := &cls.AuthRoleInfo{}
			if v, ok := role["role_name"].(string); ok && v != "" {
				authRole.RoleName = helper.String(v)
			}
			if v, ok := role["secret_id"].(string); ok && v != "" {
				authRole.SecretId = helper.String(v)
			}
			if v, ok := role["secret_key"].(string); ok && v != "" {
				authRole.SecretKey = helper.String(v)
			}
			request.AuthRoles = append(request.AuthRoles, authRole)
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsList := v.([]interface{})
		for _, item := range tagsList {
			tagMap := item.(map[string]interface{})
			tag := &cls.Tag{}
			if v, ok := tagMap["key"].(string); ok && v != "" {
				tag.Key = helper.String(v)
			}
			if v, ok := tagMap["value"].(string); ok && v != "" {
				tag.Value = helper.String(v)
			}
			request.Tags = append(request.Tags, tag)
		}
	}

	if v, ok := d.GetOk("hide_params"); ok {
		hideParamsList := v.([]interface{})
		for _, item := range hideParamsList {
			request.HideParams = append(request.HideParams, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("access_control_rules"); ok {
		rulesList := v.([]interface{})
		for _, item := range rulesList {
			ruleMap := item.(map[string]interface{})
			rule := &cls.AccessControlRule{}
			if v, ok := ruleMap["cidr_blocks"].([]interface{}); ok && len(v) > 0 {
				for _, cidr := range v {
					rule.CidrBlocks = append(rule.CidrBlocks, helper.String(cidr.(string)))
				}
			}
			if v, ok := ruleMap["action"].(string); ok && v != "" {
				rule.Action = helper.String(v)
			}
			if v, ok := ruleMap["access_mode"].(string); ok && v != "" {
				rule.AccessMode = helper.String(v)
			}
			request.AccessControlRules = append(request.AccessControlRules, rule)
		}
	}

	if v, ok := d.GetOk("remarks"); ok {
		request.Remarks = helper.String(v.(string))
	}

	if v, ok := d.GetOk("menus"); ok {
		menusList := v.([]interface{})
		for _, item := range menusList {
			request.Menus = append(request.Menus, helper.String(item.(string)))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateConsoleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls console failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cls console failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.ConsoleId == nil {
		return fmt.Errorf("ConsoleId is nil.")
	}

	consoleId = *response.Response.ConsoleId
	d.SetId(consoleId)
	return resourceTencentCloudClsConsoleRead(d, meta)
}

func resourceTencentCloudClsConsoleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_console.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		consoleId = d.Id()
	)

	respData, err := service.DescribeClsConsoleById(ctx, consoleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cls_console` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.AccessMode != nil {
		_ = d.Set("access_mode", helper.StringsInterfaces(respData.AccessMode))
	}

	if respData.LoginMode != nil {
		_ = d.Set("login_mode", respData.LoginMode)
	}

	if respData.DomainPrefix != nil {
		_ = d.Set("domain_prefix", respData.DomainPrefix)
	}

	// `accounts`, `anonymous_login`, and `auth_roles` are treated as write-only
	// fields: DescribeConsoles intentionally hides their sensitive members
	// (password / secret_id / secret_key) for security. To prevent every Read
	// from clobbering state with empty values (which would yield a perpetual
	// diff on the next plan), Read does NOT touch these three fields at all.
	// State retains whatever was last written by Create / Update, and any user
	// edit to these fields propagates through Update -> ModifyConsole.

	if respData.IntranetType != nil {
		_ = d.Set("intranet_type", respData.IntranetType)
	}

	if respData.IntranetRegion != nil {
		_ = d.Set("intranet_region", respData.IntranetRegion)
	}

	if respData.VpcId != nil {
		_ = d.Set("vpc_id", respData.VpcId)
	}

	if respData.SubnetId != nil {
		_ = d.Set("subnet_id", respData.SubnetId)
	}

	// `auth_roles` is part of the write-only credential set documented above.

	if respData.Tags != nil {
		tagsList := make([]map[string]interface{}, 0, len(respData.Tags))
		for _, tag := range respData.Tags {
			if tag == nil {
				continue
			}
			tagMap := map[string]interface{}{}
			if tag.Key != nil {
				tagMap["key"] = tag.Key
			}
			if tag.Value != nil {
				tagMap["value"] = tag.Value
			}
			tagsList = append(tagsList, tagMap)
		}
		_ = d.Set("tags", tagsList)
	}

	if respData.HideParams != nil {
		_ = d.Set("hide_params", helper.StringsInterfaces(respData.HideParams))
	}

	if respData.AccessControlRules != nil {
		rulesList := make([]map[string]interface{}, 0, len(respData.AccessControlRules))
		for _, rule := range respData.AccessControlRules {
			if rule == nil {
				continue
			}
			ruleMap := map[string]interface{}{}
			if rule.CidrBlocks != nil {
				ruleMap["cidr_blocks"] = helper.StringsInterfaces(rule.CidrBlocks)
			}
			if rule.Action != nil {
				ruleMap["action"] = rule.Action
			}
			if rule.AccessMode != nil {
				ruleMap["access_mode"] = rule.AccessMode
			}
			rulesList = append(rulesList, ruleMap)
		}
		_ = d.Set("access_control_rules", rulesList)
	}

	if respData.Remarks != nil {
		_ = d.Set("remarks", respData.Remarks)
	}

	if respData.Menus != nil {
		_ = d.Set("menus", helper.StringsInterfaces(respData.Menus))
	}

	if respData.ConsoleId != nil {
		_ = d.Set("console_id", respData.ConsoleId)
	}

	if respData.Domain != nil {
		_ = d.Set("domain", respData.Domain)
	}

	if respData.IntranetDomain != nil {
		_ = d.Set("intranet_domain", respData.IntranetDomain)
	}

	return nil
}

func resourceTencentCloudClsConsoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_console.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		consoleId = d.Id()
	)

	needChange := false
	mutableArgs := []string{
		"access_mode", "login_mode", "domain_prefix", "accounts",
		"anonymous_login", "intranet_type", "intranet_region", "vpc_id",
		"subnet_id", "auth_roles", "hide_params", "access_control_rules",
		"remarks", "menus",
	}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cls.NewModifyConsoleRequest()
		request.ConsoleId = helper.String(consoleId)

		if v, ok := d.GetOk("access_mode"); ok {
			accessModeList := v.([]interface{})
			for _, item := range accessModeList {
				request.AccessMode = append(request.AccessMode, helper.String(item.(string)))
			}
		}

		if v, ok := d.GetOkExists("login_mode"); ok {
			request.LoginMode = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("domain_prefix"); ok {
			request.DomainPrefix = helper.String(v.(string))
		}

		if v, ok := d.GetOk("accounts"); ok {
			accountsList := v.([]interface{})
			for _, item := range accountsList {
				account := item.(map[string]interface{})
				consoleAccount := &cls.ConsoleAccount{}
				if v, ok := account["user_name"].(string); ok && v != "" {
					consoleAccount.UserName = helper.String(v)
				}
				if v, ok := account["password"].(string); ok && v != "" {
					consoleAccount.Password = helper.String(v)
				}
				if v, ok := account["secret_id"].(string); ok && v != "" {
					consoleAccount.SecretId = helper.String(v)
				}
				if v, ok := account["secret_key"].(string); ok && v != "" {
					consoleAccount.SecretKey = helper.String(v)
				}
				if v, ok := account["email"].(string); ok && v != "" {
					consoleAccount.Email = helper.String(v)
				}
				request.Accounts = append(request.Accounts, consoleAccount)
			}
		}

		if v, ok := d.GetOk("anonymous_login"); ok {
			anonymousList := v.([]interface{})
			if len(anonymousList) > 0 && anonymousList[0] != nil {
				anonymous := anonymousList[0].(map[string]interface{})
				anonymousLogin := &cls.AnonymousLoginInfo{}
				if v, ok := anonymous["secret_id"].(string); ok && v != "" {
					anonymousLogin.SecretId = helper.String(v)
				}
				if v, ok := anonymous["secret_key"].(string); ok && v != "" {
					anonymousLogin.SecretKey = helper.String(v)
				}
				request.AnonymousLogin = anonymousLogin
			}
		}

		if v, ok := d.GetOkExists("intranet_type"); ok {
			request.IntranetType = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("intranet_region"); ok {
			request.IntranetRegion = helper.String(v.(string))
		}

		if v, ok := d.GetOk("vpc_id"); ok {
			request.VpcId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("subnet_id"); ok {
			request.SubnetId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("auth_roles"); ok {
			authRolesList := v.([]interface{})
			for _, item := range authRolesList {
				role := item.(map[string]interface{})
				authRole := &cls.AuthRoleInfo{}
				if v, ok := role["role_name"].(string); ok && v != "" {
					authRole.RoleName = helper.String(v)
				}
				if v, ok := role["secret_id"].(string); ok && v != "" {
					authRole.SecretId = helper.String(v)
				}
				if v, ok := role["secret_key"].(string); ok && v != "" {
					authRole.SecretKey = helper.String(v)
				}
				request.AuthRoles = append(request.AuthRoles, authRole)
			}
		}

		if v, ok := d.GetOk("hide_params"); ok {
			hideParamsList := v.([]interface{})
			for _, item := range hideParamsList {
				request.HideParams = append(request.HideParams, helper.String(item.(string)))
			}
		}

		if v, ok := d.GetOk("access_control_rules"); ok {
			rulesList := v.([]interface{})
			for _, item := range rulesList {
				ruleMap := item.(map[string]interface{})
				rule := &cls.AccessControlRule{}
				if v, ok := ruleMap["cidr_blocks"].([]interface{}); ok && len(v) > 0 {
					for _, cidr := range v {
						rule.CidrBlocks = append(rule.CidrBlocks, helper.String(cidr.(string)))
					}
				}
				if v, ok := ruleMap["action"].(string); ok && v != "" {
					rule.Action = helper.String(v)
				}
				if v, ok := ruleMap["access_mode"].(string); ok && v != "" {
					rule.AccessMode = helper.String(v)
				}
				request.AccessControlRules = append(request.AccessControlRules, rule)
			}
		}

		if v, ok := d.GetOk("remarks"); ok {
			request.Remarks = helper.String(v.(string))
		}

		if v, ok := d.GetOk("menus"); ok {
			menusList := v.([]interface{})
			for _, item := range menusList {
				request.Menus = append(request.Menus, helper.String(item.(string)))
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyConsoleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify cls console failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update cls console failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudClsConsoleRead(d, meta)
}

func resourceTencentCloudClsConsoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_console.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = cls.NewDeleteConsoleRequest()
		consoleId = d.Id()
	)

	request.ConsoleId = helper.String(consoleId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DeleteConsoleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete cls console failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cls console failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
