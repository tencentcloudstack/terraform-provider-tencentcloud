package privatedns

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
)

func ResourceTencentCloudPrivateDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDPrivateDnsZoneCreate,
		Read:   resourceTencentCloudDPrivateDnsZoneRead,
		Update: resourceTencentCloudDPrivateDnsZoneUpdate,
		Delete: resourceTencentCloudDPrivateDnsZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name, which must be in the format of standard TLD.",
			},
			"tag_set": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Description:   "Tags the private domain when it is created.",
				Deprecated:    "It has been deprecated from version 1.72.4. Use `tags` instead.",
				ConflictsWith: []string{"tags"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key of Tag.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of Tag.",
						},
					},
				},
			},
			"tags": {
				Type:          schema.TypeMap,
				Optional:      true,
				Description:   "Tags of the private dns zone.",
				ConflictsWith: []string{"tag_set"},
			},
			"vpc_set": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Associates the private domain to a VPC when it is created.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC REGION.",
						},
					},
				},
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remarks.",
			},
			"dns_forward_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      DNS_FORWARD_STATUS_DISABLED,
				ValidateFunc: tccommon.ValidateAllowedStringValue(PRIVATE_DNS_FORWARD_STATUS),
				Description:  "Whether to enable subdomain recursive DNS. Valid values: ENABLED, DISABLED. Default value: DISABLED.",
			},
			"account_vpc_set": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of authorized accounts' VPCs to associate with the private domain.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "UIN of the VPC account.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region.",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC NAME.",
						},
					},
				},
			},
			"cname_speedup_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CNAME_SPEEDUP_STATUS_ENABLED,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CNAME_SPEEDUP_STATUS),
				Description:  "CNAME acceleration: ENABLED, DISABLED, Default value is ENABLED.",
			},
		},
	}
}

func resourceTencentCloudDPrivateDnsZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_zone.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request  = privatedns.NewCreatePrivateZoneRequest()
		response = privatedns.NewCreatePrivateZoneResponse()
	)

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tag_set"); ok {
		tagSet := make([]*privatedns.TagInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			tagInfo := privatedns.TagInfo{
				TagKey:   helper.String(m["tag_key"].(string)),
				TagValue: helper.String(m["tag_value"].(string)),
			}

			tagSet = append(tagSet, &tagInfo)
		}

		request.TagSet = tagSet
	}

	if v, ok := d.GetOk("vpc_set"); ok {
		vpcSet := make([]*privatedns.VpcInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			if item != nil {
				m := item.(map[string]interface{})
				vpcInfo := privatedns.VpcInfo{
					UniqVpcId: helper.String(m["uniq_vpc_id"].(string)),
					Region:    helper.String(m["region"].(string)),
				}

				vpcSet = append(vpcSet, &vpcInfo)
			}
		}

		request.VpcSet = vpcSet
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dns_forward_status"); ok {
		request.DnsForwardStatus = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account_vpc_set"); ok {
		accountVpcSet := make([]*privatedns.AccountVpcInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			accountVpcInfo := privatedns.AccountVpcInfo{
				Uin:       helper.String(m["uin"].(string)),
				UniqVpcId: helper.String(m["uniq_vpc_id"].(string)),
				Region:    helper.String(m["region"].(string)),
				VpcName:   helper.String(m["vpc_name"].(string)),
			}

			accountVpcSet = append(accountVpcSet, &accountVpcInfo)
		}

		request.AccountVpcSet = accountVpcSet
	}

	if v, ok := d.GetOk("cname_speedup_status"); ok {
		request.CnameSpeedupStatus = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().CreatePrivateZone(request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create PrivateDns zone failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create PrivateDns zone failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.ZoneId == nil {
		return fmt.Errorf("ZoneId is nil.")
	}

	id := *response.Response.ZoneId
	d.SetId(id)

	// tags
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)
	region := client.Region
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := tccommon.BuildTagResourceName("privatedns", "zone", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudDPrivateDnsZoneRead(d, meta)
}

func resourceTencentCloudDPrivateDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_zone.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request  = privatedns.NewDescribePrivateZoneRequest()
		response *privatedns.DescribePrivateZoneResponse
		zoneId   = d.Id()
	)

	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = zoneId

	request.ZoneId = &zoneId
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient(iacExtInfo).DescribePrivateZone(request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe PrivateDns zone failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read DnsPod Domain failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.PrivateZone == nil {
		log.Printf("[WARN]%s resource `tencentcloud_private_dns_zone` [%s] not found, please check if it has been deleted.\n", logId, zoneId)
		d.SetId("")
		return nil
	}

	info := response.Response.PrivateZone
	if info.Domain != nil {
		_ = d.Set("domain", info.Domain)
	}

	if info.Tags != nil {
		tagSets := make([]map[string]interface{}, 0, len(info.Tags))
		for _, item := range info.Tags {
			tagSets = append(tagSets, map[string]interface{}{
				"tag_key":   item.TagKey,
				"tag_value": item.TagValue,
			})
		}

		_ = d.Set("tag_set", tagSets)
	}

	if info.VpcSet != nil {
		vpcSet := make([]map[string]interface{}, 0, len(info.VpcSet))
		for _, item := range info.VpcSet {
			vpcSet = append(vpcSet, map[string]interface{}{
				"uniq_vpc_id": item.UniqVpcId,
				"region":      item.Region,
			})
		}

		_ = d.Set("vpc_set", vpcSet)
	}

	if info.Remark != nil {
		_ = d.Set("remark", info.Remark)
	}

	if info.DnsForwardStatus != nil {
		_ = d.Set("dns_forward_status", info.DnsForwardStatus)
	}

	if info.CnameSpeedupStatus != nil {
		_ = d.Set("cname_speedup_status", info.CnameSpeedupStatus)
	}

	if info.AccountVpcSet != nil {
		accountVpcSet := make([]map[string]interface{}, 0, len(info.AccountVpcSet))
		for _, item := range info.AccountVpcSet {
			accountVpcSet = append(accountVpcSet, map[string]interface{}{
				"uin":         item.Uin,
				"uniq_vpc_id": item.UniqVpcId,
				"region":      item.Region,
			})
		}

		_ = d.Set("account_vpc_set", accountVpcSet)
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)
	region := client.Region
	tags, err := tagService.DescribeResourceTags(ctx, "privatedns", "zone", region, zoneId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudDPrivateDnsZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_zone.update")()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		zoneId = d.Id()
	)

	if d.HasChange("tag_set") {
		return fmt.Errorf("tag_set do not support change, please use tags instead.")
	}

	if d.HasChange("remark") || d.HasChange("dns_forward_status") || d.HasChange("cname_speedup_status") {
		request := privatedns.NewModifyPrivateZoneRequest()
		request.ZoneId = &zoneId
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		if v, ok := d.GetOk("dns_forward_status"); ok {
			request.DnsForwardStatus = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cname_speedup_status"); ok {
			request.CnameSpeedupStatus = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().ModifyPrivateZone(request)
			if e != nil {
				return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify privateDns zone info failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("vpc_set") || d.HasChange("account_vpc_set") {
		request := privatedns.NewModifyPrivateZoneVpcRequest()
		request.ZoneId = &zoneId
		if v, ok := d.GetOk("vpc_set"); ok {
			var vpcSets = make([]*privatedns.VpcInfo, 0)
			items := v.([]interface{})
			for _, item := range items {
				if item != nil {
					value := item.(map[string]interface{})
					if value["uniq_vpc_id"].(string) != "" && value["region"].(string) != "" {
						vpcInfo := &privatedns.VpcInfo{
							UniqVpcId: helper.String(value["uniq_vpc_id"].(string)),
							Region:    helper.String(value["region"].(string)),
						}

						vpcSets = append(vpcSets, vpcInfo)
					}
				}
			}

			request.VpcSet = vpcSets
		}

		if v, ok := d.GetOk("account_vpc_set"); ok {
			var accVpcSets = make([]*privatedns.AccountVpcInfo, 0)
			items := v.([]interface{})
			for _, item := range items {
				value := item.(map[string]interface{})
				accVpcInfo := &privatedns.AccountVpcInfo{
					UniqVpcId: helper.String(value["uniq_vpc_id"].(string)),
					Region:    helper.String(value["region"].(string)),
					Uin:       helper.String(value["uin"].(string)),
					VpcName:   helper.String(value["vpc_name"].(string)),
				}

				accVpcSets = append(accVpcSets, accVpcInfo)
			}

			request.AccountVpcSet = accVpcSets
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().ModifyPrivateZoneVpc(request)
			if e != nil {
				return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify privateDns zone vpc failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("tags") {
		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(client)
		region := client.Region
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("privatedns", "zone", region, zoneId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudDPrivateDnsZoneRead(d, meta)
}

func resourceTencentCloudDPrivateDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_zone.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = privatedns.NewDeletePrivateZoneRequest()
		zoneId  = d.Id()
	)

	request.ZoneId = &zoneId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().DeletePrivateZone(request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete privateDns zone failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
