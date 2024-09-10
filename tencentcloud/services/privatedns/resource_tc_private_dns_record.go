package privatedns

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
)

func ResourceTencentCloudPrivateDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDPrivateDnsRecordCreate,
		Read:   resourceTencentCloudDPrivateDnsRecordRead,
		Update: resourceTencentCloudDPrivateDnsRecordUpdate,
		Delete: resourceTencentCloudDPrivateDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Private domain ID.",
			},
			"record_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Record type. Valid values: `A`, `AAAA`, `CNAME`, `MX`, `TXT`, `PTR`.",
			},
			"sub_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subdomain, such as `www`, `m`, and `@`.",
			},
			"record_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Record value, such as IP: 192.168.10.2, CNAME: cname.qcloud.com, and MX: mail.qcloud.com.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Record weight. Value range: 1~100.",
			},
			"mx": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "MX priority, which is required when the record type is MX." +
					" Valid values: 5, 10, 15, 20, 30, 40, 50.",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Record cache time. The smaller the value, the faster the record will take effect. Value range: 1~86400s.",
			},
		},
	}
}

func resourceTencentCloudDPrivateDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_record.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = privatedns.NewCreatePrivateZoneRecordRequest()
		response = privatedns.NewCreatePrivateZoneRecordResponse()
		zoneId   string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("record_type"); ok {
		request.RecordType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sub_domain"); ok {
		request.SubDomain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_value"); ok {
		request.RecordValue = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("weight"); ok {
		request.Weight = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("mx"); ok {
		request.MX = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("ttl"); ok {
		request.TTL = helper.Int64(int64(v.(int)))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().CreatePrivateZoneRecord(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("create PrivateDns record failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create PrivateDns record failed, reason:%s\n", logId, err.Error())
		return err
	}

	recordId := *response.Response.RecordId
	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudDPrivateDnsRecordRead(d, meta)
}

func resourceTencentCloudDPrivateDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_zone.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PrivateDnsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("record id strategy is can't read, id is borken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	records, err := service.DescribePrivateDnsRecordByFilter(ctx, zoneId, nil)
	if err != nil {
		return err
	}

	if len(records) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `PrivateDnsRecord` [%s] not found, please check if it has been deleted.\n", logId, recordId)
		return nil
	}

	var record *privatedns.PrivateZoneRecord
	for _, item := range records {
		if *item.RecordId == recordId {
			record = item
		}
	}

	if record == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PrivateDnsRecord` [%s] not found, please check if it has been deleted.\n", logId, recordId)
		return nil
	}

	_ = d.Set("zone_id", record.ZoneId)
	_ = d.Set("record_type", record.RecordType)
	_ = d.Set("sub_domain", record.SubDomain)
	_ = d.Set("record_value", record.RecordValue)
	_ = d.Set("weight", record.Weight)
	_ = d.Set("mx", record.MX)
	_ = d.Set("ttl", record.TTL)

	return nil
}

func resourceTencentCloudDPrivateDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_record.update")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = privatedns.NewModifyPrivateZoneRecordRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("record id strategy is can't read, id is borken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.RecordId = helper.String(recordId)
	if v, ok := d.GetOk("record_type"); ok {
		request.RecordType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sub_domain"); ok {
		request.SubDomain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_value"); ok {
		request.RecordValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("weight"); ok {
		request.Weight = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("mx"); ok {
		request.MX = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("ttl"); ok {
		request.TTL = helper.Int64(int64(v.(int)))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().ModifyPrivateZoneRecord(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify privateDns record info failed, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudDPrivateDnsRecordRead(d, meta)
}

func resourceTencentCloudDPrivateDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_record.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = privatedns.NewDeletePrivateZoneRecordRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("record id strategy is can't read, id is borken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.RecordId = helper.String(recordId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().DeletePrivateZoneRecord(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete privateDns record failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
