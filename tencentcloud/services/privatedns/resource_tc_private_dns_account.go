package privatedns

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudPrivateDnsAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPrivateDnsAccountCreate,
		Read:   resourceTencentCloudPrivateDnsAccountRead,
		Delete: resourceTencentCloudPrivateDnsAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_uin": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Uin of the associated account.",
			},
			"account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the associated account.",
			},
			"nickname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Nickname of the associated account.",
			},
		},
	}
}

func resourceTencentCloudPrivateDnsAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_account.create")()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		uin   string
	)

	if v, ok := d.GetOk("account_uin"); ok {
		uin = v.(string)
	}

	service := PrivateDnsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := service.CreatePrivateDnsAccount(ctx, uin)
	if err != nil {
		log.Printf("[CRITAL]%s create private dns account failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(uin)

	return resourceTencentCloudPrivateDnsAccountRead(d, meta)
}

func resourceTencentCloudPrivateDnsAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_account.read")()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		uin   = d.Id()
	)

	service := PrivateDnsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	account, err := service.DescribePrivateDnsAccountByUin(ctx, uin)
	if err != nil {
		log.Printf("[CRITAL]%s read private dns account failed, reason:%+v", logId, err)
		return err
	}

	if account == nil {
		d.SetId("")
		log.Printf("[WARN]%s private dns account [%s] not found, removing from state", logId, uin)
		return nil
	}

	_ = d.Set("account_uin", account.Uin)
	_ = d.Set("account", account.Account)
	_ = d.Set("nickname", account.Nickname)

	return nil
}

func resourceTencentCloudPrivateDnsAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_account.delete")()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		uin   = d.Id()
	)

	service := PrivateDnsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := service.DeletePrivateDnsAccount(ctx, uin)
	if err != nil {
		log.Printf("[CRITAL]%s delete private dns account failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
