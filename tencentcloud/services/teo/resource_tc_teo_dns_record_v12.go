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

func ResourceTencentCloudTeoDnsRecordV12() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordV12Create,
		Read:   resourceTencentCloudTeoDnsRecordV12Read,
		Update: resourceTencentCloudTeoDnsRecordV12Update,
		Delete: resourceTencentCloudTeoDnsRecordV12Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the site.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content.",
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS record resolution route. Default is Default.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record cache time. Value range: 60~86400, unit: seconds. Default: 300.",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. Value range: -1~100, -1 means no weight, 0 means no resolution. Default: -1.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority. Value range: 0~50. Default: 0. Only valid when Type is MX.",
			},

			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record ID.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record resolution status. Valid values: enable, disable.",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record creation time.",
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecordV12Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.teo_dns_record.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(context.Background())
	ctx := context.WithValue(context.Background(), tccommon.LogIdKey, logId)

	zoneId := d.Get("zone_id").(string)
	name := d.Get("name").(string)
	dnsType := d.Get("type").(string)
	content := d.Get("content").(string)

	request := teo.NewCreateDnsRecordRequest()
	request.ZoneId = helper.String(zoneId)
	request.Name = helper.String(name)
	request.Type = helper.String(dnsType)
	request.Content = helper.String(content)

	if v, ok := d.GetOk("location"); ok {
		request.Location = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("ttl"); ok {
		request.TTL = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("weight"); ok {
		request.Weight = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	var response *teo.CreateDnsRecordResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateDnsRecordWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo dns record failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil {
		return fmt.Errorf("create teo dns record failed, response is nil")
	}

	recordId := helper.PString(response.Response.RecordId)
	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoDnsRecordV12Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV12Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.teo_dns_record.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(context.Background())
	ctx := context.WithValue(context.Background(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("invalid resource id format, expected: zone_id#record_id")
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	record, err := service.DescribeTeoDnsRecordById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}

	if record == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", helper.PString(record.ZoneId))
	_ = d.Set("name", helper.PString(record.Name))
	_ = d.Set("type", helper.PString(record.Type))
	_ = d.Set("content", helper.PString(record.Content))
	_ = d.Set("location", helper.PString(record.Location))
	_ = d.Set("ttl", helper.PInt64(record.TTL))
	_ = d.Set("weight", helper.PInt64(record.Weight))
	_ = d.Set("priority", helper.PInt64(record.Priority))
	_ = d.Set("record_id", helper.PString(record.RecordId))
	_ = d.Set("status", helper.PString(record.Status))
	_ = d.Set("created_on", helper.PString(record.CreatedOn))

	return nil
}

func resourceTencentCloudTeoDnsRecordV12Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.teo_dns_record.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(context.Background())
	ctx := context.WithValue(context.Background(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("invalid resource id format, expected: zone_id#record_id")
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	if d.HasChange("name") || d.HasChange("type") || d.HasChange("content") ||
		d.HasChange("location") || d.HasChange("ttl") || d.HasChange("weight") || d.HasChange("priority") {

		dnsRecord := &teo.DnsRecord{}
		dnsRecord.ZoneId = helper.String(zoneId)
		dnsRecord.RecordId = helper.String(recordId)
		dnsRecord.Name = helper.String(d.Get("name").(string))
		dnsRecord.Type = helper.String(d.Get("type").(string))
		dnsRecord.Content = helper.String(d.Get("content").(string))

		if v, ok := d.GetOk("location"); ok {
			dnsRecord.Location = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("ttl"); ok {
			dnsRecord.TTL = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("weight"); ok {
			dnsRecord.Weight = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("priority"); ok {
			dnsRecord.Priority = helper.IntInt64(v.(int))
		}

		request := teo.NewModifyDnsRecordsRequest()
		request.ZoneId = helper.String(zoneId)
		request.DnsRecords = []*teo.DnsRecord{dnsRecord}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyDnsRecordsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo dns record failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoDnsRecordV12Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV12Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.teo_dns_record.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(context.Background())
	ctx := context.WithValue(context.Background(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("invalid resource id format, expected: zone_id#record_id")
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	request := teo.NewDeleteDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.RecordIds = []*string{helper.String(recordId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DeleteDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo dns record failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
