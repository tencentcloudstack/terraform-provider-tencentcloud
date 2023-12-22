package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudGaapSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapSecurityPolicyCreate,
		Read:   resourceTencentCloudGaapSecurityPolicyRead,
		Update: resourceTencentCloudGaapSecurityPolicyUpdate,
		Delete: resourceTencentCloudGaapSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"proxy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the GAAP proxy.",
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"ACCEPT", "DROP"}),
				ForceNew:     true,
				Description:  "Default policy. Valid value: `ACCEPT` and `DROP`.",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether policy is enable, default value is `true`.",
			},
		},
	}
}

func resourceTencentCloudGaapSecurityPolicyCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_security_policy.create")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	proxyId := d.Get("proxy_id").(string)
	action := d.Get("action").(string)
	enable := d.Get("enable").(bool)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id, err := service.CreateSecurityPolicy(ctx, proxyId, action)
	if err != nil {
		return err
	}

	d.SetId(id)

	if enable {
		if err := service.EnableSecurityPolicy(ctx, proxyId, id); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapSecurityPolicyRead(d, m)
}

func resourceTencentCloudGaapSecurityPolicyRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_security_policy.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	proxyId, status, action, exist, err := service.DescribeSecurityPolicy(ctx, id)
	if err != nil {
		return err
	}

	if !exist {
		d.SetId("")
		return nil
	}

	_ = d.Set("proxy_id", proxyId)
	_ = d.Set("action", action)
	_ = d.Set("enable", status == GAAP_SECURITY_POLICY_BOUND)

	return nil
}

func resourceTencentCloudGaapSecurityPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_security_policy.update")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	proxyId := d.Get("proxy_id").(string)
	enable := d.Get("enable").(bool)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if enable {
		if err := service.EnableSecurityPolicy(ctx, proxyId, id); err != nil {
			return err
		}
	} else {
		if err := service.DisableSecurityPolicy(ctx, proxyId, id); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapSecurityPolicyRead(d, m)
}

func resourceTencentCloudGaapSecurityPolicyDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_security_policy.delete")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	enable := d.Get("enable").(bool)
	proxyId := d.Get("proxy_id").(string)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if enable {
		if err := service.DisableSecurityPolicy(ctx, proxyId, id); err != nil {
			return err
		}
	}

	return service.DeleteSecurityPolicy(ctx, id)
}
