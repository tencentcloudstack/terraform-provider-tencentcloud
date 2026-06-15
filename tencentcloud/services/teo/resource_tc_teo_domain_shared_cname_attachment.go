package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoDomainSharedCnameAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDomainSharedCnameAttachmentCreate,
		Read:   resourceTencentCloudTeoDomainSharedCnameAttachmentRead,
		Update: resourceTencentCloudTeoDomainSharedCnameAttachmentUpdate,
		Delete: resourceTencentCloudTeoDomainSharedCnameAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zone ID that the acceleration domain belongs to.",
			},
			"shared_cname": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The shared CNAME to bind to.",
			},
			"domain_names": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The acceleration domain names to bind.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// bindSharedCNAMEInBatches calls BindSharedCNAMEWithContext in batches of up to 20 domain names per request.
func bindSharedCNAMEInBatches(ctx context.Context, logId string, meta interface{}, zoneId, sharedCname, bindType string, domainNames []string) error {
	if len(domainNames) == 0 {
		request := teov20220901.NewBindSharedCNAMERequest()
		request.ZoneId = helper.String(zoneId)
		request.BindType = helper.String(bindType)
		bindMap := teov20220901.BindSharedCNAMEMap{
			SharedCNAME: helper.String(sharedCname),
		}
		request.BindSharedCNAMEMaps = []*teov20220901.BindSharedCNAMEMap{&bindMap}
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().BindSharedCNAMEWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})
		if reqErr != nil {
			return reqErr
		}
	}
	const batchSize = 20
	for i := 0; i < len(domainNames); i += batchSize {
		end := i + batchSize
		if end > len(domainNames) {
			end = len(domainNames)
		}
		batch := domainNames[i:end]

		request := teov20220901.NewBindSharedCNAMERequest()
		request.ZoneId = helper.String(zoneId)
		request.BindType = helper.String(bindType)
		bindMap := teov20220901.BindSharedCNAMEMap{
			SharedCNAME: helper.String(sharedCname),
		}
		for _, dn := range batch {
			bindMap.DomainNames = append(bindMap.DomainNames, helper.String(dn))
		}
		request.BindSharedCNAMEMaps = []*teov20220901.BindSharedCNAMEMap{&bindMap}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().BindSharedCNAMEWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})
		if reqErr != nil {
			return reqErr
		}
	}
	return nil
}

func resourceTencentCloudTeoDomainSharedCnameAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_domain_shared_cname_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		zoneId      string
		sharedCname string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("shared_cname"); ok {
		sharedCname = v.(string)
	}

	var domainNames []string
	if v, ok := d.GetOk("domain_names"); ok {
		for _, domainName := range v.(*schema.Set).List() {
			domainNames = append(domainNames, domainName.(string))
		}
	}

	d.SetId(strings.Join([]string{zoneId, sharedCname}, tccommon.FILED_SP))

	if reqErr := bindSharedCNAMEInBatches(ctx, logId, meta, zoneId, sharedCname, "bind", domainNames); reqErr != nil {
		log.Printf("[CRITAL]%s create teo_domain_shared_cname_attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTeoDomainSharedCnameAttachmentRead(d, meta)
}

func resourceTencentCloudTeoDomainSharedCnameAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_domain_shared_cname_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDescribeSharedCNAMERequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	sharedCname := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("shared-cname"),
			Values: []*string{helper.String(sharedCname)},
		},
	}
	request.Limit = helper.Int64(200)

	var response *teov20220901.DescribeSharedCNAMEResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeSharedCNAMEWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read teo_domain_shared_cname_attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil || len(response.Response.SharedCNAMEInfo) == 0 {
		log.Printf("[WARN]%s resource `teo_domain_shared_cname_attachment` [%s] not found, removing from state", logId, d.Id())
		d.SetId("")
		return nil
	}

	// Find the matching shared CNAME info and collect all bound domain names
	var boundDomainNames []string
	var found bool
	for _, info := range response.Response.SharedCNAMEInfo {
		if info.SharedCNAME != nil && *info.SharedCNAME == sharedCname {
			for _, ref := range info.AccelerationDomains {
				if ref.Instance != nil {
					boundDomainNames = append(boundDomainNames, *ref.Instance)
				}
			}
			found = true
			break
		}
	}

	if !found {
		log.Printf("[WARN]%s resource `teo_domain_shared_cname_attachment` [%s] shared cname not found, removing from state", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("shared_cname", sharedCname)
	_ = d.Set("domain_names", boundDomainNames)

	return nil
}

func resourceTencentCloudTeoDomainSharedCnameAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_domain_shared_cname_attachment.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	sharedCname := idSplit[1]

	if d.HasChange("domain_names") {
		oldVal, newVal := d.GetChange("domain_names")

		oldSet := oldVal.(*schema.Set)
		newSet := newVal.(*schema.Set)

		// Compute domains to unbind (in old but not in new)
		var domainsToUnbind []string
		for _, dn := range oldSet.Difference(newSet).List() {
			domainsToUnbind = append(domainsToUnbind, dn.(string))
		}

		// Compute domains to bind (in new but not in old)
		var domainsToBind []string
		for _, dn := range newSet.Difference(oldSet).List() {
			domainsToBind = append(domainsToBind, dn.(string))
		}

		// Unbind removed domains
		if len(domainsToUnbind) > 0 {
			if reqErr := bindSharedCNAMEInBatches(ctx, logId, meta, zoneId, sharedCname, "unbind", domainsToUnbind); reqErr != nil {
				log.Printf("[CRITAL]%s unbind domains during update teo_domain_shared_cname_attachment failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}

		// Bind new domains
		if len(domainsToBind) > 0 {
			if reqErr := bindSharedCNAMEInBatches(ctx, logId, meta, zoneId, sharedCname, "bind", domainsToBind); reqErr != nil {
				log.Printf("[CRITAL]%s bind domains during update teo_domain_shared_cname_attachment failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	}

	return resourceTencentCloudTeoDomainSharedCnameAttachmentRead(d, meta)
}

func resourceTencentCloudTeoDomainSharedCnameAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_domain_shared_cname_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	sharedCname := idSplit[1]

	var domainNames []string
	if v, ok := d.GetOk("domain_names"); ok {
		for _, domainName := range v.(*schema.Set).List() {
			domainNames = append(domainNames, domainName.(string))
		}
	}

	if len(domainNames) == 0 {
		return nil
	}

	if reqErr := bindSharedCNAMEInBatches(ctx, logId, meta, zoneId, sharedCname, "unbind", domainNames); reqErr != nil {
		log.Printf("[CRITAL]%s delete teo_domain_shared_cname_attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
