package lighthouse

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentCreate,
		Read:   resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentRead,
		Delete: resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"blueprint_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Lighthouse blueprint ID.",
			},

			"account_ids": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of target TencentCloud account IDs to share the blueprint with.",
			},
		},
	}
}

func resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_share_blueprint_across_account_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = lighthouse.NewShareBlueprintAcrossAccountsRequest()
	)

	blueprintId := d.Get("blueprint_id").(string)
	accountIds := helper.InterfacesStrings(d.Get("account_ids").([]interface{}))

	request.BlueprintId = &blueprintId
	request.AccountIds = make([]*string, 0, len(accountIds))
	for _, v := range accountIds {
		accountId := v
		request.AccountIds = append(request.AccountIds, &accountId)
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
		log.Printf("[CRITAL]%s create lighthouse shareBlueprintAcrossAccountAttachment failed, reason:%+v", logId, err)
		return err
	}

	// Build composite ID: blueprint_id#sorted_account_ids
	sort.Strings(accountIds)
	d.SetId(blueprintId + tccommon.FILED_SP + strings.Join(accountIds, tccommon.FILED_SP))

	return resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentRead(d, meta)
}

func resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_share_blueprint_across_account_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	blueprintId := idSplit[0]
	expectedAccountIds := idSplit[1:]

	// Sort expected account IDs
	sort.Strings(expectedAccountIds)

	_ = ctx

	request := lighthouse.NewDescribeBlueprintsShareAcrossAccountInfosRequest()
	request.BlueprintIds = []*string{&blueprintId}
	limit := int64(100)
	request.Limit = &limit

	var response *lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().DescribeBlueprintsShareAcrossAccountInfos(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read lighthouse shareBlueprintAcrossAccountAttachment failed, reason:%+v", logId, err)
		return err
	}

	// Check if response is valid
	if response == nil || response.Response == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseShareBlueprintAcrossAccountAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	// Collect actual account IDs from the response
	var actualAccountIds []string
	if response.Response.BlueprintShareAcrossAccountInfoSet != nil {
		for _, info := range response.Response.BlueprintShareAcrossAccountInfoSet {
			if info.AccountId != nil {
				actualAccountIds = append(actualAccountIds, *info.AccountId)
			}
		}
	}

	// Sort actual account IDs
	sort.Strings(actualAccountIds)

	// Check if the expected account IDs match the actual account IDs
	if len(expectedAccountIds) != len(actualAccountIds) {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseShareBlueprintAcrossAccountAttachment` [%s] not found, account count mismatch (expected=%d, actual=%d).\n", logId, d.Id(), len(expectedAccountIds), len(actualAccountIds))
		return nil
	}

	for i := range expectedAccountIds {
		if expectedAccountIds[i] != actualAccountIds[i] {
			d.SetId("")
			log.Printf("[WARN]%s resource `LighthouseShareBlueprintAcrossAccountAttachment` [%s] not found, account mismatch at index %d (expected=%s, actual=%s).\n", logId, d.Id(), i, expectedAccountIds[i], actualAccountIds[i])
			return nil
		}
	}

	_ = d.Set("blueprint_id", blueprintId)
	_ = d.Set("account_ids", actualAccountIds)

	return nil
}

func resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_share_blueprint_across_account_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	blueprintId := idSplit[0]
	accountIds := idSplit[1:]

	request := lighthouse.NewCancelShareBlueprintAcrossAccountsRequest()
	request.BlueprintId = &blueprintId
	request.AccountIds = make([]*string, 0, len(accountIds))
	for _, v := range accountIds {
		accountId := v
		request.AccountIds = append(request.AccountIds, &accountId)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().CancelShareBlueprintAcrossAccounts(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete lighthouse shareBlueprintAcrossAccountAttachment failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
