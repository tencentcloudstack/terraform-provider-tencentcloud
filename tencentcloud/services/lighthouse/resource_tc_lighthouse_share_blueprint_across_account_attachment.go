package lighthouse

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const (
	shareBlueprintBatchSize  = 10
	cancelBlueprintBatchSize = 10
)

func ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentCreate,
		Read:   resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentRead,
		Update: resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentUpdate,
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
				Type:        schema.TypeSet,
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

	blueprintId := d.Get("blueprint_id").(string)
	accountIds := helper.InterfacesStrings(d.Get("account_ids").(*schema.Set).List())

	// Batch share accounts (API limit: 10 per request)
	for i := 0; i < len(accountIds); i += shareBlueprintBatchSize {
		end := i + shareBlueprintBatchSize
		if end > len(accountIds) {
			end = len(accountIds)
		}

		request := lighthouse.NewShareBlueprintAcrossAccountsRequest()
		request.BlueprintId = &blueprintId
		batch := accountIds[i:end]
		request.AccountIds = make([]*string, 0, len(batch))
		for _, v := range batch {
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
			log.Printf("[CRITAL]%s create tencentcloud_lighthouse_share_blueprint_across_account_attachment failed, reason:%+v", logId, err)
			return err
		}
	}

	// Set resource ID as blueprint_id
	d.SetId(blueprintId)

	return resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentRead(d, meta)
}

func resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_share_blueprint_across_account_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	blueprintId := d.Id()

	_ = ctx

	request := lighthouse.NewDescribeBlueprintsShareAcrossAccountInfosRequest()
	request.BlueprintIds = []*string{&blueprintId}
	limit := int64(100)
	request.Limit = &limit
	var offset int64 = 0

	// Collect actual account IDs from the response with pagination support
	var actualAccountIds []string
	for {
		request.Offset = &offset

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
			log.Printf("[CRITAL]%s read tencentcloud_lighthouse_share_blueprint_across_account_attachment failed, reason:%+v", logId, err)
			return err
		}

		if response == nil || response.Response == nil {
			log.Printf("[WARN]%s resource `tencentcloud_lighthouse_share_blueprint_across_account_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
			d.SetId("")
			return nil
		}

		if response.Response.BlueprintShareAcrossAccountInfoSet != nil {
			for _, info := range response.Response.BlueprintShareAcrossAccountInfoSet {
				if info.AccountId != nil {
					actualAccountIds = append(actualAccountIds, *info.AccountId)
				}
			}
		}

		// Check if there are more pages to fetch
		if response.Response.TotalCount != nil && offset+int64(len(actualAccountIds)) < *response.Response.TotalCount {
			offset += limit
		} else {
			break
		}
	}

	_ = d.Set("blueprint_id", blueprintId)
	_ = d.Set("account_ids", actualAccountIds)

	return nil
}

func resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_share_blueprint_across_account_attachment.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	blueprintId := d.Id()

	if d.HasChange("account_ids") {
		oldAccountIds, newAccountIds := d.GetChange("account_ids")
		oldSet := helper.InterfacesStrings(oldAccountIds.(*schema.Set).List())
		newSet := helper.InterfacesStrings(newAccountIds.(*schema.Set).List())

		// Calculate accounts to remove (in old but not in new)
		toRemove := diffStringSlices(oldSet, newSet)

		// Calculate accounts to add (in new but not in old)
		toAdd := diffStringSlices(newSet, oldSet)

		// Batch cancel sharing for removed accounts (API limit: 10 per request)
		if len(toRemove) > 0 {
			for i := 0; i < len(toRemove); i += cancelBlueprintBatchSize {
				end := i + cancelBlueprintBatchSize
				if end > len(toRemove) {
					end = len(toRemove)
				}

				removeRequest := lighthouse.NewCancelShareBlueprintAcrossAccountsRequest()
				removeRequest.BlueprintId = &blueprintId
				batch := toRemove[i:end]
				removeRequest.AccountIds = make([]*string, 0, len(batch))
				for _, v := range batch {
					accountId := v
					removeRequest.AccountIds = append(removeRequest.AccountIds, &accountId)
				}

				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().CancelShareBlueprintAcrossAccounts(removeRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, removeRequest.GetAction(), removeRequest.ToJsonString(), result.ToJsonString())
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s update tencentcloud_lighthouse_share_blueprint_across_account_attachment (remove accounts) failed, reason:%+v", logId, err)
					return err
				}
			}
		}

		// Batch share for added accounts (API limit: 10 per request)
		if len(toAdd) > 0 {
			for i := 0; i < len(toAdd); i += shareBlueprintBatchSize {
				end := i + shareBlueprintBatchSize
				if end > len(toAdd) {
					end = len(toAdd)
				}

				addRequest := lighthouse.NewShareBlueprintAcrossAccountsRequest()
				addRequest.BlueprintId = &blueprintId
				batch := toAdd[i:end]
				addRequest.AccountIds = make([]*string, 0, len(batch))
				for _, v := range batch {
					accountId := v
					addRequest.AccountIds = append(addRequest.AccountIds, &accountId)
				}

				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ShareBlueprintAcrossAccounts(addRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, addRequest.GetAction(), addRequest.ToJsonString(), result.ToJsonString())
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s update tencentcloud_lighthouse_share_blueprint_across_account_attachment (add accounts) failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	}

	return resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentRead(d, meta)
}

// diffStringSlices returns elements in 'a' that are not in 'b'
func diffStringSlices(a []string, b []string) []string {
	bMap := make(map[string]bool)
	for _, v := range b {
		bMap[v] = true
	}

	var diff []string
	for _, v := range a {
		if !bMap[v] {
			diff = append(diff, v)
		}
	}
	return diff
}

func resourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_share_blueprint_across_account_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	blueprintId := d.Id()
	accountIds := helper.InterfacesStrings(d.Get("account_ids").(*schema.Set).List())

	// Batch cancel sharing (API limit: 10 per request)
	for i := 0; i < len(accountIds); i += cancelBlueprintBatchSize {
		end := i + cancelBlueprintBatchSize
		if end > len(accountIds) {
			end = len(accountIds)
		}

		request := lighthouse.NewCancelShareBlueprintAcrossAccountsRequest()
		request.BlueprintId = &blueprintId
		batch := accountIds[i:end]
		request.AccountIds = make([]*string, 0, len(batch))
		for _, v := range batch {
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
			log.Printf("[CRITAL]%s delete tencentcloud_lighthouse_share_blueprint_across_account_attachment failed, reason:%+v", logId, err)
			return err
		}
	}

	return nil
}
