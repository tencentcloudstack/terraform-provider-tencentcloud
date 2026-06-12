package teo

import (
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
			"bind_shared_cname_maps": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The binding relationships between acceleration domains and shared CNAMEs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shared_cname": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The shared CNAME to bind to.",
						},
						"domain_names": {
							Type:        schema.TypeList,
							Required:    true,
							ForceNew:    true,
							Description: "The acceleration domain names to bind, up to 20.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoDomainSharedCnameAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_domain_shared_cname_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = teov20220901.NewBindSharedCNAMERequest()
		zoneId  string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	request.BindType = helper.String("bind")

	var sharedCname string
	var domainNames []string

	if v, ok := d.GetOk("bind_shared_cname_maps"); ok {
		for _, item := range v.([]interface{}) {
			mapItem := item.(map[string]interface{})
			bindMap := teov20220901.BindSharedCNAMEMap{}
			if v, ok := mapItem["shared_cname"].(string); ok && v != "" {
				bindMap.SharedCNAME = helper.String(v)
				sharedCname = v
			}

			if v, ok := mapItem["domain_names"]; ok {
				for _, domainName := range v.([]interface{}) {
					dn := domainName.(string)
					bindMap.DomainNames = append(bindMap.DomainNames, helper.String(dn))
					domainNames = append(domainNames, dn)
				}
			}

			request.BindSharedCNAMEMaps = append(request.BindSharedCNAMEMaps, &bindMap)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().BindSharedCNAME(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo_domain_shared_cname_attachment failed, Response is nil"))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo_domain_shared_cname_attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create teo_domain_shared_cname_attachment, logId: %s, id: %s", logId, logId, d.Id())

	if sharedCname == "" || len(domainNames) == 0 {
		return fmt.Errorf("Create teo_domain_shared_cname_attachment failed, shared_cname or domain_names is empty")
	}

	d.SetId(strings.Join([]string{zoneId, sharedCname, strings.Join(domainNames, ",")}, tccommon.FILED_SP))

	return resourceTencentCloudTeoDomainSharedCnameAttachmentRead(d, meta)
}

func resourceTencentCloudTeoDomainSharedCnameAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_domain_shared_cname_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = teov20220901.NewDescribeSharedCNAMERequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	sharedCname := idSplit[1]
	domainNamesStr := idSplit[2]
	domainNames := strings.Split(domainNamesStr, ",")

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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeSharedCNAME(request)
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

	// Find the matching shared CNAME info
	var found bool
	for _, info := range response.Response.SharedCNAMEInfo {
		if info.SharedCNAME != nil && *info.SharedCNAME == sharedCname {
			// Check if all expected domains are bound
			boundDomains := make(map[string]bool)
			for _, ref := range info.AccelerationDomains {
				if ref.Instance != nil {
					boundDomains[*ref.Instance] = true
				}
			}

			allBound := true
			for _, dn := range domainNames {
				if !boundDomains[dn] {
					allBound = false
					break
				}
			}

			if !allBound {
				log.Printf("[WARN]%s resource `teo_domain_shared_cname_attachment` [%s] domains not all bound, removing from state", logId, d.Id())
				d.SetId("")
				return nil
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

	bindSharedCnameMaps := []map[string]interface{}{
		{
			"shared_cname": sharedCname,
			"domain_names": domainNames,
		},
	}
	_ = d.Set("bind_shared_cname_maps", bindSharedCnameMaps)

	return nil
}

func resourceTencentCloudTeoDomainSharedCnameAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_domain_shared_cname_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = teov20220901.NewBindSharedCNAMERequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	sharedCname := idSplit[1]
	domainNamesStr := idSplit[2]
	domainNames := strings.Split(domainNamesStr, ",")

	request.ZoneId = helper.String(zoneId)
	request.BindType = helper.String("unbind")

	bindMap := teov20220901.BindSharedCNAMEMap{
		SharedCNAME: helper.String(sharedCname),
	}
	for _, dn := range domainNames {
		bindMap.DomainNames = append(bindMap.DomainNames, helper.String(dn))
	}
	request.BindSharedCNAMEMaps = []*teov20220901.BindSharedCNAMEMap{&bindMap}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().BindSharedCNAME(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo_domain_shared_cname_attachment failed, Response is nil"))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo_domain_shared_cname_attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
