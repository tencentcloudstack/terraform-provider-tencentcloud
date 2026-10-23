package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudFwmCfwManaged() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmCfwManagedCreate,
		Read:   resourceTencentCloudFwmCfwManagedRead,
		Update: resourceTencentCloudFwmCfwManagedUpdate,
		Delete: resourceTencentCloudFwmCfwManagedDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"member_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target member id of the organization.",
			},
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "App id of the target member, used to locate the member by AppId filter.",
			},
			"cfw_managed": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Cloud firewall managed switch, 1: enable managed, 0: disable managed.",
			},
		},
	}
}

func resourceTencentCloudFwmCfwManagedCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_managed.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	memberId := d.Get("member_id").(string)
	appId := d.Get("app_id").(string)

	d.SetId(appId + tccommon.FILED_SP + memberId)

	return resourceTencentCloudFwmCfwManagedUpdate(d, meta)
}

func resourceTencentCloudFwmCfwManagedRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_managed.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	appId := items[0]
	memberId := items[1]

	member, err := service.DescribeFwmCfwManagedByAppId(ctx, appId, memberId)
	if err != nil {
		return err
	}

	if member == nil {
		log.Printf("[WARN]%s resource `tencentcloud_fwm_cfw_managed` id=%s not found, maybe the member has left the organization.", logId, d.Id())
		d.SetId("")
		return nil
	}

	if member.MemberId != nil {
		_ = d.Set("member_id", member.MemberId)
	}

	if member.AppId != nil {
		_ = d.Set("app_id", member.AppId)
	}

	if member.CfwManaged != nil {
		_ = d.Set("cfw_managed", member.CfwManaged)
	}

	return nil
}

func resourceTencentCloudFwmCfwManagedUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_managed.update")()
	defer tccommon.InconsistentCheck(d, meta)

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	memberId := items[1]

	cfwManaged := int64(d.Get("cfw_managed").(int))

	current, err := service.DescribeFwmCfwManagedByAppId(ctx, items[0], memberId)
	if err != nil {
		return err
	}

	skip := false
	if current != nil && current.CfwManaged != nil && *current.CfwManaged == cfwManaged {
		skip = true
	}

	if !skip {
		err = service.ModifyFwmCfwManaged(ctx, memberId, cfwManaged)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudFwmCfwManagedRead(d, meta)
}

func resourceTencentCloudFwmCfwManagedDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_managed.delete")()
	defer tccommon.InconsistentCheck(d, meta)

	return nil
}
