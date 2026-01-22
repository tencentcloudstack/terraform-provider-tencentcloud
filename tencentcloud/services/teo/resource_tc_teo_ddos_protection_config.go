package teo

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoDdosProtectionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDdosProtectionConfigCreate,
		Read:   resourceTencentCloudTeoDdosProtectionConfigRead,
		Update: resourceTencentCloudTeoDdosProtectionConfigUpdate,
		Delete: resourceTencentCloudTeoDdosProtectionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},

			"ddos_protection": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Specifies the exclusive Anti-DDoS configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protection_option": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the protection scope of standalone DDoS. valid values:.\n<li>protect_all_domains: specifies exclusive Anti-DDoS protection for all domain names in the site. newly added domain names automatically enable exclusive Anti-DDoS protection. when this parameter is specified, DomainDDoSProtections will not be processed.</li>.\n<li>protect_specified_domains: only applicable to specified domains. specific scope can be set via DomainDDoSProtection parameter.</li>.",
						},
						"domain_ddos_protections": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "Anti-DDoS configuration of the domain. specifies the exclusive ddos protection settings for the domain in request parameters.\n<li>When ProtectionOption remains protect_specified_domains, the domain names not filled in keep their exclusive Anti-DDoS protection configuration unchanged, while explicitly specified domain names are updated according to the input parameters.</li>.\n<li>When ProtectionOption switches from protect_all_domains to protect_specified_domains: if DomainDDoSProtections is empty, disable exclusive DDoS protection for all domains under the site; if DomainDDoSProtections is not empty, disable or maintain exclusive DDoS protection for the domain names specified in the parameter, and disable exclusive DDoS protection for other unlisted domain names.</li>.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Domain name.",
									},
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Standalone DDoS switch of the domain. valid values:.\n<li>on: enabled;</li>.\n<li>off: closed.</li>.",
									},
								},
							},
						},
						"shared_cname_ddos_protections": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the exclusive DDoS protection configuration of a shared CNAME. used as an output parameter.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Standalone DDoS switch of the domain. valid values:.\n<li>on: enabled;</li>.\n<li>off: closed.</li>.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoDdosProtectionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ddos_protection_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	d.SetId(zoneId)
	return resourceTencentCloudTeoDdosProtectionConfigUpdate(d, meta)
}

func resourceTencentCloudTeoDdosProtectionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ddos_protection_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  = d.Id()
	)

	respData, err := service.DescribeTeoDdosProtectionConfigById(ctx, zoneId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_ddos_protection_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	dMap := make(map[string]interface{}, 0)
	if respData.ProtectionOption != nil {
		dMap["protection_option"] = respData.ProtectionOption
	}

	if respData.DomainDDoSProtections != nil {
		domainDDoSProtectionsList := make([]map[string]interface{}, 0, len(respData.DomainDDoSProtections))
		for _, domainDDoSProtections := range respData.DomainDDoSProtections {
			domainDDoSProtectionsMap := map[string]interface{}{}
			if domainDDoSProtections.Domain != nil {
				domainDDoSProtectionsMap["domain"] = domainDDoSProtections.Domain
			}

			if domainDDoSProtections.Switch != nil {
				domainDDoSProtectionsMap["switch"] = domainDDoSProtections.Switch
			}

			domainDDoSProtectionsList = append(domainDDoSProtectionsList, domainDDoSProtectionsMap)
		}

		dMap["domain_ddos_protections"] = domainDDoSProtectionsList
	}

	if respData.SharedCNAMEDDoSProtections != nil {
		sharedCNAMEDDoSProtectionsList := make([]map[string]interface{}, 0, len(respData.SharedCNAMEDDoSProtections))
		for _, sharedCNAMEDDoSProtections := range respData.SharedCNAMEDDoSProtections {
			sharedCNAMEDDoSProtectionsMap := map[string]interface{}{}
			if sharedCNAMEDDoSProtections.Domain != nil {
				sharedCNAMEDDoSProtectionsMap["domain"] = sharedCNAMEDDoSProtections.Domain
			}

			if sharedCNAMEDDoSProtections.Switch != nil {
				sharedCNAMEDDoSProtectionsMap["switch"] = sharedCNAMEDDoSProtections.Switch
			}

			sharedCNAMEDDoSProtectionsList = append(sharedCNAMEDDoSProtectionsList, sharedCNAMEDDoSProtectionsMap)
		}

		dMap["shared_cname_ddos_protections"] = sharedCNAMEDDoSProtectionsList
	}

	_ = d.Set("ddos_protection", []interface{}{dMap})
	return nil
}

func resourceTencentCloudTeoDdosProtectionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ddos_protection_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewModifyDDoSProtectionRequest()
		zoneId  = d.Id()
	)

	if dDoSProtectionMap, ok := helper.InterfacesHeadMap(d, "ddos_protection"); ok {
		dDoSProtection := teov20220901.DDoSProtection{}
		if v, ok := dDoSProtectionMap["protection_option"].(string); ok && v != "" {
			dDoSProtection.ProtectionOption = helper.String(v)
		}

		if v, ok := dDoSProtectionMap["domain_ddos_protections"]; ok {
			for _, item := range v.(*schema.Set).List() {
				domainDDoSProtectionsMap := item.(map[string]interface{})
				domainDDoSProtection := teov20220901.DomainDDoSProtection{}
				if v, ok := domainDDoSProtectionsMap["domain"].(string); ok && v != "" {
					domainDDoSProtection.Domain = helper.String(v)
				}

				if v, ok := domainDDoSProtectionsMap["switch"].(string); ok && v != "" {
					domainDDoSProtection.Switch = helper.String(v)
				}

				dDoSProtection.DomainDDoSProtections = append(dDoSProtection.DomainDDoSProtections, &domainDDoSProtection)
			}
		}

		request.DDoSProtection = &dDoSProtection
	}

	request.ZoneId = &zoneId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyDDoSProtectionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update teo ddos protection config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTeoDdosProtectionConfigRead(d, meta)
}

func resourceTencentCloudTeoDdosProtectionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ddos_protection_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
