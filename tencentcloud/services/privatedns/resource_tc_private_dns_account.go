package privatedns

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatednsv20201028 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
			"account": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Private DNS account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uin": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Root account UIN.",
						},
						"account": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Root account name.",
						},
						"nickname": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Account name.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPrivateDnsAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_account.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = privatednsv20201028.NewCreatePrivateDNSAccountRequest()
		uin     string
	)

	if accountMap, ok := helper.InterfacesHeadMap(d, "account"); ok {
		privateDNSAccount := privatednsv20201028.PrivateDNSAccount{}
		if v, ok := accountMap["uin"]; ok {
			privateDNSAccount.Uin = helper.String(v.(string))
			uin = v.(string)
		}

		if v, ok := accountMap["account"]; ok {
			privateDNSAccount.Account = helper.String(v.(string))
		}

		if v, ok := accountMap["nickname"]; ok {
			privateDNSAccount.Nickname = helper.String(v.(string))
		}

		request.Account = &privateDNSAccount
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsV20201028Client().CreatePrivateDNSAccountWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		_ = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create private dns account failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(uin)

	return resourceTencentCloudPrivateDnsAccountRead(d, meta)
}

func resourceTencentCloudPrivateDnsAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = PrivatednsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		uin     = d.Id()
	)

	respData, err := service.DescribePrivateDnsAccountById(ctx, uin)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `private_dns_account` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	tmpList := make([]map[string]interface{}, 0)
	tmpMap := map[string]interface{}{}
	if respData.Uin != nil {
		tmpMap["uin"] = respData.Uin
	}

	if respData.Account != nil {
		tmpMap["account"] = respData.Account
	}

	if respData.Nickname != nil {
		tmpMap["nickname"] = respData.Nickname
	}

	tmpList = append(tmpList, tmpMap)
	_ = d.Set("account", tmpList)

	return nil
}

func resourceTencentCloudPrivateDnsAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_account.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = privatednsv20201028.NewDeletePrivateDNSAccountRequest()
		uin     = d.Id()
	)

	privateDNSAccount := privatednsv20201028.PrivateDNSAccount{}
	privateDNSAccount.Uin = helper.String(uin)
	request.Account = &privateDNSAccount

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsV20201028Client().DeletePrivateDNSAccountWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		_ = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete private dns account failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
