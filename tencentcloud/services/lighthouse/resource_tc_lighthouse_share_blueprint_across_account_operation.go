package lighthouse

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudLighthouseShareBlueprintAcrossAccountOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseShareBlueprintAcrossAccountOperationCreate,
		Read:   resourceTencentCloudLighthouseShareBlueprintAcrossAccountOperationRead,
		Delete: resourceTencentCloudLighthouseShareBlueprintAcrossAccountOperationDelete,
		Schema: map[string]*schema.Schema{
			"blueprint_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Blueprint ID. Only custom images in NORMAL state can be shared.",
			},

			"account_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 10,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of account IDs receiving the shared image. Max 10 account IDs. Must be main accounts.",
			},
		},
	}
}

func resourceTencentCloudLighthouseShareBlueprintAcrossAccountOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_share_blueprint_across_account.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = lighthouse.NewShareBlueprintAcrossAccountsRequest()
	)

	if v, ok := d.GetOk("blueprint_id"); ok {
		request.BlueprintId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account_ids"); ok {
		accountIdsList := v.([]interface{})
		request.AccountIds = make([]*string, len(accountIdsList))
		for i, accountId := range accountIdsList {
			request.AccountIds[i] = helper.String(accountId.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ShareBlueprintAcrossAccounts(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse shareBlueprintAcrossAccount failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.BuildToken())

	return resourceTencentCloudLighthouseShareBlueprintAcrossAccountOperationRead(d, meta)
}

func resourceTencentCloudLighthouseShareBlueprintAcrossAccountOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_share_blueprint_across_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseShareBlueprintAcrossAccountOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_share_blueprint_across_account.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
