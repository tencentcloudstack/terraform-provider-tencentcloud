package ses

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSesBlackListDelete() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesBlackListDeleteCreate,
		Read:   resourceTencentCloudSesBlackListDeleteRead,
		Delete: resourceTencentCloudSesBlackListDeleteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email_address": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Email addresses to be unblocklisted.",
			},
		},
	}
}

func resourceTencentCloudSesBlackListDeleteCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_black_list_delete.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = ses.NewDeleteBlackListRequest()
		emailAddress string
	)
	if v, ok := d.GetOk("email_address"); ok {
		emailAddress = v.(string)
		request.EmailAddressList = append(request.EmailAddressList, helper.String(v.(string)))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().DeleteBlackList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ses BlackList failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(emailAddress)

	return resourceTencentCloudSesBlackListDeleteRead(d, meta)
}

func resourceTencentCloudSesBlackListDeleteRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_black_list_delete.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSesBlackListDeleteDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_black_list_delete.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
