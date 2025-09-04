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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 100),
				Description:  "Record weight. Value range: 1~100.",
			},
			"mx": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "MX priority, which is required when the record type is MX. Valid values: 5, 10, 15, 20, 30, 40, 50.",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Record cache time. The smaller the value, the faster the record will take effect. Value range: 1~86400s.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"enabled", "disabled"}),
				Description:  "Record status. Valid values: `enabled`, `disabled`.",
			},
		},
	}
}

func resourceTencentCloudDPrivateDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_record.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = PrivateDnsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
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
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.RecordId == nil {
			return resource.NonRetryableError(fmt.Errorf("Create PrivateDns record failed, Response is nil."))
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

	// wait
	_, err = service.DescribePrivateDnsRecordById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}

	// set record status
	if v, ok := d.GetOk("status"); ok {
		status := v.(string)
		if status == "disabled" {
			request := privatedns.NewModifyRecordsStatusRequest()
			request.ZoneId = &zoneId
			request.RecordIds = []*int64{helper.StrToInt64Point(recordId)}
			request.Status = helper.String(status)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().ModifyRecordsStatus(request)
				if e != nil {
					return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify PrivateDns record status failed, reason:%s\n", logId, err.Error())
				return err
			}
		}
	}

	return resourceTencentCloudDPrivateDnsRecordRead(d, meta)
}

func resourceTencentCloudDPrivateDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_record.read")()
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

	record, err := service.DescribePrivateDnsRecordById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}

	if record == nil {
		log.Printf("[WARN]%s resource `tencentcloud_private_dns_record` [%s] not found, please check if it has been deleted.\n", logId, recordId)
		d.SetId("")
		return nil
	}

	if record.ZoneId != nil {
		_ = d.Set("zone_id", record.ZoneId)
	}

	if record.RecordType != nil {
		_ = d.Set("record_type", record.RecordType)
	}

	if record.SubDomain != nil {
		_ = d.Set("sub_domain", record.SubDomain)
	}

	if record.RecordValue != nil {
		_ = d.Set("record_value", record.RecordValue)
	}

	if record.Weight != nil {
		_ = d.Set("weight", record.Weight)
	}

	if record.MX != nil {
		_ = d.Set("mx", record.MX)
	}

	if record.TTL != nil {
		_ = d.Set("ttl", record.TTL)
	}

	if record.Enabled != nil {
		if *record.Enabled == 1 {
			_ = d.Set("status", "enabled")
		} else {
			_ = d.Set("status", "disabled")
		}
	}

	return nil
}

func resourceTencentCloudDPrivateDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_record.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = privatedns.NewModifyPrivateZoneRecordRequest()
		needChange bool
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("record id strategy is can't read, id is borken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	mutableArgs := []string{"record_type", "sub_domain", "record_value", "weight", "mx", "ttl"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
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

		request.ZoneId = helper.String(zoneId)
		request.RecordId = helper.String(recordId)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().ModifyPrivateZoneRecord(request)
			if e != nil {
				return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify privateDns record info failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			request := privatedns.NewModifyRecordsStatusRequest()
			request.ZoneId = &zoneId
			request.RecordIds = []*int64{helper.StrToInt64Point(recordId)}
			request.Status = helper.String(v.(string))
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivateDnsClient().ModifyRecordsStatus(request)
				if e != nil {
					return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify PrivateDns record status failed, reason:%s\n", logId, err.Error())
				return err
			}
		}
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
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete privateDns record failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
