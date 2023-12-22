package dnspod

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDnspodModifyDomainOwnerOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodModifyDomainOwnerOperationCreate,
		Read:   resourceTencentCloudDnspodModifyDomainOwnerOperationRead,
		Delete: resourceTencentCloudDnspodModifyDomainOwnerOperationDelete,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"account": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The account to which the domain needs to be transferred, supporting Uin or email format.",
			},

			"domain_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.",
			},
		},
	}
}

func resourceTencentCloudDnspodModifyDomainOwnerOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_modify_domain_owner_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = dnspod.NewModifyDomainOwnerRequest()
		domain  string
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account"); ok {
		request.Account = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		request.DomainId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyDomainOwner(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dnspod modify_domain_owner failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(domain)

	return resourceTencentCloudDnspodModifyDomainOwnerOperationRead(d, meta)
}

func resourceTencentCloudDnspodModifyDomainOwnerOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_modify_domain_owner_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDnspodModifyDomainOwnerOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_modify_domain_owner_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
