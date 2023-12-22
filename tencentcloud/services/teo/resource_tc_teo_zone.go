package teo

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoZone() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoZoneRead,
		Create: resourceTencentCloudTeoZoneCreate,
		Update: resourceTencentCloudTeoZoneUpdate,
		Delete: resourceTencentCloudTeoZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site name. When accessing CNAME/NS, please pass the second-level domain (example.com) as the site name; when accessing without a domain name, please leave this value empty.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site access type. The value of this parameter is as follows, and the default is `partial` if not filled in. Valid values: `partial`: CNAME access; `full`: NS access; `noDomainAccess`: No domain access.",
			},

			"alias_zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alias site identifier. Limit the input to a combination of numbers, English, - and _, within 20 characters. For details, refer to the alias site identifier. If there is no such usage scenario, leave this field empty.",
			},

			"area": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "When the `type` value is `partial` or `full`, the acceleration region of the L7 domain name. The following are the values of this parameter, and the default value is `overseas` if not filled in. When the `type` value is `noDomainAccess`, please leave this value empty. Valid values: `global`: Global availability zone; `mainland`: Chinese mainland availability zone; `overseas`: Global availability zone (excluding Chinese mainland).",
			},

			"plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target Plan ID to be bound. When you have an existing Plan in your account, you can fill in this parameter to directly bind the site to the Plan. If you do not have a Plan that can be bound at the moment, please go to the console to purchase a Plan to complete the site creation.",
			},

			"paused": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the site is disabled.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site status. Valid values: `active`: NS is switched; `pending`: NS is not switched; `moved`: NS is moved; `deactivated`: this site is blocked.",
			},

			"ownership_verification": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Ownership verification information. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_verification": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CNAME access, using DNS to resolve the information required for authentication. For details, please refer to [Site/Domain Name Ownership Verification ](https://cloud.tencent.com/document/product/1552/70789#7af6ecf8-afca-4e35-8811-b5797ed1bde5). Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subdomain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Host record.",
									},
									"record_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Record type.",
									},
									"record_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Record the value.",
									},
								},
							},
						},
					},
				},
			},

			"name_servers": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "NS list allocated by Tencent Cloud.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTeoZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_zone.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request  = teo.NewCreateZoneRequest()
		response *teo.CreateZoneResponse
		zoneId   string
	)

	if v, ok := d.GetOk("zone_name"); ok {
		request.ZoneName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("alias_zone_name"); ok {
		request.AliasZoneName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plan_id"); ok {
		request.PlanId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateZone(request)
		if e != nil {
			if tccommon.IsExpectError(e, []string{"ResourceInUse", "ResourceInUse.Others"}) {
				return resource.NonRetryableError(e)
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo zone failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(zoneId)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoZone(ctx, zoneId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if *instance.Status == "pending" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("zone status is %v, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}

	if v, _ := d.GetOkExists("paused"); v != nil {
		if v.(bool) {
			err := service.ModifyZoneStatus(ctx, zoneId, v.(bool), "create")
			if err != nil {
				return err
			}
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		resourceName := fmt.Sprintf("qcs::teo::uin/:zone/%s", zoneId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudTeoZoneRead(d, meta)
}

func resourceTencentCloudTeoZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_zone.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	zoneId := d.Id()

	zone, err := service.DescribeTeoZone(ctx, zoneId)

	if err != nil {
		return err
	}

	if zone == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `zone` %s does not exist.\n", logId, d.Id())
		return nil
	}

	if zone.ZoneName != nil {
		_ = d.Set("zone_name", zone.ZoneName)
	}

	if zone.Type != nil {
		_ = d.Set("type", zone.Type)
	}

	if zone.AliasZoneName != nil {
		_ = d.Set("alias_zone_name", zone.AliasZoneName)
	}

	if zone.Area != nil {
		_ = d.Set("area", zone.Area)
	}

	if zone.Resources != nil && len(zone.Resources) > 0 {
		if zone.Resources[0].PlanId != nil {
			_ = d.Set("plan_id", zone.Resources[0].PlanId)
		}
	}

	if zone.Paused != nil {
		_ = d.Set("paused", zone.Paused)
	}

	if zone.OwnershipVerification != nil {
		ownershipVerificationMap := map[string]interface{}{}
		if zone.OwnershipVerification.DnsVerification != nil {
			dnsVerificationMap := map[string]interface{}{}

			if zone.OwnershipVerification.DnsVerification.Subdomain != nil {
				dnsVerificationMap["subdomain"] = zone.OwnershipVerification.DnsVerification.Subdomain
			}

			if zone.OwnershipVerification.DnsVerification.RecordType != nil {
				dnsVerificationMap["record_type"] = zone.OwnershipVerification.DnsVerification.RecordType
			}

			if zone.OwnershipVerification.DnsVerification.RecordValue != nil {
				dnsVerificationMap["record_value"] = zone.OwnershipVerification.DnsVerification.RecordValue
			}

			ownershipVerificationMap["dns_verification"] = []interface{}{dnsVerificationMap}

		}
		_ = d.Set("ownership_verification", []interface{}{ownershipVerificationMap})
	}

	if zone.NameServers != nil {
		_ = d.Set("name_servers", zone.NameServers)
	}

	if zone.Status != nil {
		_ = d.Set("status", zone.Status)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "teo", "zone", "", zoneId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTeoZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_zone.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	request := teo.NewModifyZoneRequest()

	zoneId := d.Id()
	request.ZoneId = &zoneId

	if d.HasChange("paused") {
		if v, ok := d.GetOkExists("paused"); ok {
			err := service.ModifyZoneStatus(ctx, zoneId, v.(bool), "update")
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("type") || d.HasChange("alias_zone_name") || d.HasChange("area") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}

		if v, ok := d.GetOk("alias_zone_name"); ok {
			request.AliasZoneName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("area"); ok {
			request.Area = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyZone(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create teo zone failed, reason:%+v", logId, err)
			return err
		}

		service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, errRet := service.DescribeTeoZone(ctx, zoneId)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			if *instance.Status == "pending" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("zone status is %v, retry...", *instance.Status))
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("teo", "zone", "", zoneId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoZoneRead(d, meta)
}

func resourceTencentCloudTeoZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_zone.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	zoneId := d.Id()

	instance, err := service.DescribeTeoZone(ctx, zoneId)
	if err != nil {
		return err
	}

	if !*instance.Paused {
		err := service.ModifyZoneStatus(ctx, zoneId, true, "delete")
		if err != nil {
			return err
		}
	}

	if err = service.DeleteTeoZoneById(ctx, zoneId); err != nil {
		return err
	}

	return nil
}
