package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoDnsRecordV11() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordV11Create,
		Read:   resourceTencentCloudTeoDnsRecordV11Read,
		Update: resourceTencentCloudTeoDnsRecordV11Update,
		Delete: resourceTencentCloudTeoDnsRecordV11Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},

			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.",
			},

			"record_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.",
			},

			"record_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. Fill in the corresponding content according to the Type value.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time. The value range is 60~86400, in seconds. The smaller the value, the faster the modification takes effect globally. Default is 300.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority. This parameter is only valid when Type (DNS record type) is MX. The value range is 0~50. The smaller the value, the higher the priority. Default is 0.",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. The value range is -1~100. -1 means no weight is set, and 0 means no resolution. Weight configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME. Default is -1.",
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "DNS record resolution line. Default is Default, indicating the default resolution line, which takes effect in all regions. Resolution line configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME.",
			},

			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record ID.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record resolution status. Valid values: enable (active), disable (inactive).",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time.",
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecordV11Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v11.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId   string
		domain   string
		recordId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	request := teo.NewCreateDnsRecordRequest()
	request.ZoneId = helper.String(zoneId)
	request.Name = helper.String(domain)

	if v, ok := d.GetOk("record_type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_value"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ttl"); ok {
		ttl := int64(v.(int))
		request.TTL = &ttl
	}

	if v, ok := d.GetOk("priority"); ok {
		priority := int64(v.(int))
		request.Priority = &priority
	}

	if v, ok := d.GetOk("weight"); ok {
		weight := int64(v.(int))
		request.Weight = &weight
	}

	if v, ok := d.GetOk("location"); ok {
		request.Location = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateDnsRecordWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.RecordId == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo dns record failed, Response is nil."))
		}

		recordId = *result.Response.RecordId
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s create teo dns record failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoDnsRecordV11Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV11Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v11.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("parse resource ID %s failed, format should be zone_id#record_id", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	request := teo.NewDescribeDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.Filters = []*teo.AdvancedFilter{
		{
			Name:   helper.String("id"),
			Values: []*string{helper.String(recordId)},
		},
	}

	var dnsRecord *teo.DnsRecord

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.DnsRecords == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo dns record failed, Response is nil."))
		}

		if len(result.Response.DnsRecords) == 0 {
			return resource.NonRetryableError(fmt.Errorf("DNS record not found"))
		}

		for _, record := range result.Response.DnsRecords {
			if record.RecordId != nil && *record.RecordId == recordId {
				dnsRecord = record
				break
			}
		}

		if dnsRecord == nil {
			return resource.NonRetryableError(fmt.Errorf("DNS record not found"))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s read teo dns record failed, reason:%+v", logId, err)
		return err
	}

	if dnsRecord == nil {
		log.Printf("[WARN]%s dns record not exist, remove from state", logId)
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("domain", dnsRecord.Name)
	_ = d.Set("record_type", dnsRecord.Type)
	_ = d.Set("record_value", dnsRecord.Content)
	_ = d.Set("ttl", dnsRecord.TTL)
	_ = d.Set("priority", dnsRecord.Priority)
	_ = d.Set("weight", dnsRecord.Weight)
	_ = d.Set("location", dnsRecord.Location)
	_ = d.Set("record_id", dnsRecord.RecordId)
	_ = d.Set("status", dnsRecord.Status)
	_ = d.Set("created_at", dnsRecord.CreatedOn)
	_ = d.Set("updated_at", dnsRecord.ModifiedOn)

	return nil
}

func resourceTencentCloudTeoDnsRecordV11Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v11.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("parse resource ID %s failed, format should be zone_id#record_id", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	dnsRecord := teo.DnsRecord{
		RecordId: helper.String(recordId),
	}

	if d.HasChange("record_value") {
		if v, ok := d.GetOk("record_value"); ok {
			dnsRecord.Content = helper.String(v.(string))
		}
	}

	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			ttl := int64(v.(int))
			dnsRecord.TTL = &ttl
		}
	}

	if d.HasChange("priority") {
		if v, ok := d.GetOk("priority"); ok {
			priority := int64(v.(int))
			dnsRecord.Priority = &priority
		}
	}

	if d.HasChange("weight") {
		if v, ok := d.GetOk("weight"); ok {
			weight := int64(v.(int))
			dnsRecord.Weight = &weight
		}
	}

	if d.HasChange("location") {
		if v, ok := d.GetOk("location"); ok {
			dnsRecord.Location = helper.String(v.(string))
		}
	}

	request := teo.NewModifyDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.DnsRecords = []*teo.DnsRecord{&dnsRecord}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), "")
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s update teo dns record failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoDnsRecordV11Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV11Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v11.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("parse resource ID %s failed, format should be zone_id#record_id", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	request := teo.NewDeleteDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.RecordIds = []*string{helper.String(recordId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DeleteDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), "")
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s delete teo dns record failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
