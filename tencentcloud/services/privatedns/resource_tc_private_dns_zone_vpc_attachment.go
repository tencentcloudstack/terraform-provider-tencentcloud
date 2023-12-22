package privatedns

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPrivateDnsZoneVpcAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPrivateDnsZoneVpcAttachmentCreate,
		Read:   resourceTencentCloudPrivateDnsZoneVpcAttachmentRead,
		Delete: resourceTencentCloudPrivateDnsZoneVpcAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "PrivateZone ID.",
			},
			"vpc_set": {
				Optional:     true,
				ForceNew:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"account_vpc_set"},
				Type:         schema.TypeList,
				Description:  "New add vpc info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Uniq Vpc Id.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Vpc region.",
						},
					},
				},
			},
			"account_vpc_set": {
				Optional:     true,
				ForceNew:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"vpc_set"},
				Type:         schema.TypeList,
				Description:  "New add account vpc info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Uniq Vpc Id.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Vpc region.",
						},
						"uin": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Vpc owner uin. To grant role authorization to this account.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPrivateDnsZoneVpcAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_zone_vpc_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		request      = privatedns.NewAddSpecifyPrivateZoneVpcRequest()
		asyncRequest = privatedns.NewQueryAsyncBindVpcStatusRequest()
		zoneId       string
		uniqVpcId    string
		uniqId       string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("vpc_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			vpcInfo := new(privatedns.VpcInfo)
			if v, ok := dMap["uniq_vpc_id"]; ok {
				vpcInfo.UniqVpcId = helper.String(v.(string))
				uniqVpcId = v.(string)
			}

			if v, ok := dMap["region"]; ok {
				vpcInfo.Region = helper.String(v.(string))
			}

			request.VpcSet = append(request.VpcSet, vpcInfo)
		}
	}

	if v, ok := d.GetOk("account_vpc_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			accountVpcInfo := new(privatedns.AccountVpcInfo)
			if v, ok := dMap["uniq_vpc_id"]; ok {
				accountVpcInfo.UniqVpcId = helper.String(v.(string))
				uniqVpcId = v.(string)
			}

			if v, ok := dMap["region"]; ok {
				accountVpcInfo.Region = helper.String(v.(string))
			}

			if v, ok := dMap["uin"]; ok {
				accountVpcInfo.Uin = helper.String(v.(string))
			}

			accountVpcInfo.VpcName = helper.String("")

			request.AccountVpcSet = append(request.AccountVpcSet, accountVpcInfo)
		}
	}

	request.Sync = helper.Bool(false)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().AddSpecifyPrivateZoneVpc(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response.UniqId == nil {
			e = fmt.Errorf("PrivateDns ZoneVpcAttachment not exists")
			return resource.NonRetryableError(e)
		}

		uniqId = *result.Response.UniqId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create PrivateDns ZoneVpcAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{zoneId, uniqVpcId}, tccommon.FILED_SP))

	// wait
	asyncRequest.UniqId = &uniqId
	err = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().QueryAsyncBindVpcStatus(asyncRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, asyncRequest.GetAction(), asyncRequest.ToJsonString(), asyncRequest.ToJsonString())
		}

		if *result.Response.Status == "success" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("query async bind vpc status is %s.", *result.Response.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s query async bind vpc status failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPrivateDnsZoneVpcAttachmentRead(d, meta)
}

func resourceTencentCloudPrivateDnsZoneVpcAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_zone_vpc_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PrivateDnsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	zoneId := idSplit[0]
	uniqVpcId := idSplit[1]

	ZoneVpcAttachment, err := service.DescribePrivateDnsZoneVpcAttachmentById(ctx, zoneId)
	if err != nil {
		return err
	}

	if ZoneVpcAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PrivateDnsZoneVpcAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ZoneVpcAttachment.ZoneId != nil {
		_ = d.Set("zone_id", ZoneVpcAttachment.ZoneId)
	}

	if ZoneVpcAttachment.VpcSet != nil {
		vpcSetList := []interface{}{}
		for _, vpcSet := range ZoneVpcAttachment.VpcSet {
			vpcSetMap := map[string]interface{}{}

			if *vpcSet.UniqVpcId == uniqVpcId {
				vpcSetMap["uniq_vpc_id"] = *vpcSet.UniqVpcId
				vpcSetMap["region"] = *vpcSet.Region
				vpcSetList = append(vpcSetList, vpcSetMap)
				break
			}
		}

		_ = d.Set("vpc_set", vpcSetList)
	}

	if ZoneVpcAttachment.AccountVpcSet != nil {
		accountVpcSetList := []interface{}{}
		for _, accountVpcSet := range ZoneVpcAttachment.AccountVpcSet {
			accountVpcSetMap := map[string]interface{}{}

			if *accountVpcSet.UniqVpcId == uniqVpcId {
				accountVpcSetMap["uniq_vpc_id"] = *accountVpcSet.UniqVpcId
				accountVpcSetMap["region"] = *accountVpcSet.Region
				accountVpcSetMap["uin"] = *accountVpcSet.Uin
				accountVpcSetList = append(accountVpcSetList, accountVpcSetMap)
				break
			}
		}

		_ = d.Set("account_vpc_set", accountVpcSetList)
	}

	return nil
}

func resourceTencentCloudPrivateDnsZoneVpcAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_zone_vpc_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PrivateDnsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		region  string
		uin     string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	zoneId := idSplit[0]
	uniqVpcId := idSplit[1]

	// get vpc detail
	ZoneVpcAttachment, err := service.DescribePrivateDnsZoneVpcAttachmentById(ctx, zoneId)
	if err != nil {
		return err
	}

	if ZoneVpcAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PrivateDnsZoneVpcAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ZoneVpcAttachment.VpcSet != nil {
		for _, vpcSet := range ZoneVpcAttachment.VpcSet {
			if *vpcSet.UniqVpcId == uniqVpcId {
				region = *vpcSet.Region
				break
			}
		}
	}

	if ZoneVpcAttachment.AccountVpcSet != nil {
		for _, accountVpcSet := range ZoneVpcAttachment.AccountVpcSet {
			if *accountVpcSet.UniqVpcId == uniqVpcId {
				region = *accountVpcSet.Region
				uin = *accountVpcSet.Uin
				break
			}
		}
	}

	if err = service.DeletePrivateDnsZoneVpcAttachmentById(ctx, zoneId, uniqVpcId, region, uin); err != nil {
		return err
	}

	return nil
}
