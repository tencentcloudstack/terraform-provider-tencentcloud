package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoDnsRecordV10() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordV10Create,
		Read:   resourceTencentCloudTeoDnsRecordV10Read,
		Update: resourceTencentCloudTeoDnsRecordV10Update,
		Delete: resourceTencentCloudTeoDnsRecordV10Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone id.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name. if domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "DNS record type. valid values are:\n" +
					"	- A: points domain name to an external ipv4 address, such as 8.8.8.8;\n" +
					"	- AAAA: points domain name to an external ipv6 address;\n" +
					"	- MX: used for email servers. when there are multiple mx records, lower priority value, higher priority;\n" +
					"	- CNAME: points domain name to another domain name, which then resolves to final ip address;\n" +
					"	- TXT: identifies and describes of domain name, commonly used for domain verification and spf records (anti-spam);\n" +
					"	- NS: if you need to delegate subdomain to another dns service provider for resolution, you need to add an ns record. root domain cannot add ns records;\n" +
					"	- CAA: specifies ca that can issue certificates for this site;\n" +
					"	- SRV: identifies a server using a service, commonly used in microsoft's directory management.\n" +
					"Different record types, such as SRV and CAA records, have different requirements for host record names and record value formats. for detailed descriptions and format examples of each record type, please refer to: [introduction to dns record types](https://intl.cloud.tencent.com/document/product/1552/90453?from_cn_redirect=1#2f681022-91ab-4a9e-ac3d-0a6c454d954e).",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. fill in corresponding content according to type value. if domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "DNS record resolution route. if not specified, default is DEFAULT, which means default resolution route and is effective in all regions.\n\n- resolution route configuration is only applicable when type (dns record type) is A, AAAA, or CNAME.\n- resolution route configuration is only applicable to standard version and enterprise edition packages. for valid values, please refer to: [resolution routes and corresponding code enumeration](https://intl.cloud.tencent.com/document/product/1552/112542?from_cn_redirect=1).",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time. users can specify a value range of 60-86400. smaller value, faster modification records will take effect in all regions. default value: 300. unit: seconds.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. users can specify a value range of -1 to 100. a value of 0 means no resolution. if not specified, default is -1, which means no weight is set. weight configuration is only applicable when type (dns record type) is A, AAAA, or CNAME. note: for the same subdomain, different dns records with the same resolution route should either all have weights set or none have weights set.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority, which takes effect only when type (dns record type) is MX. smaller value, higher priority. users can specify a value range of 0-50. default value is 0 if not specified.",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "DNS record resolution status, following values:\n" +
					"	- enable: has taken effect;\n" +
					"	- disable: has been disabled.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modify time.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceTencentCloudTeoDnsRecordV10Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v10.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId   string
		recordId string
	)
	var (
		request  = teov20220901.NewCreateDnsRecordRequest()
		response = teov20220901.NewCreateDnsRecordResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateDnsRecordWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo dns record v10 failed, reason:%+v", logId, err)
		return err
	}

	recordId = *response.Response.RecordId

	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	// Wait for DNS record to become active
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		respData, e := service.DescribeTeoDnsRecordById(ctx, zoneId, recordId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if respData == nil {
			return resource.RetryableError(fmt.Errorf("DNS record %s not found", recordId))
		}
		if respData.Status == nil {
			return resource.RetryableError(fmt.Errorf("DNS record %s status not available", recordId))
		}
		if *respData.Status != "enable" {
			return resource.RetryableError(fmt.Errorf("DNS record %s is not active yet, status: %s", recordId, *respData.Status))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s wait for teo dns record v10 to become active failed, reason:%+v", logId, err)
		return fmt.Errorf("waiting for DNS record to become active failed: %v", err)
	}

	return resourceTencentCloudTeoDnsRecordV10Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV10Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v10.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	respData, err := service.DescribeTeoDnsRecordById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}
	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_dns_record_v10` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.ZoneId != nil {
		_ = d.Set("zone_id", respData.ZoneId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.Location != nil {
		_ = d.Set("location", respData.Location)
	}

	if respData.Content != nil {
		_ = d.Set("content", respData.Content)
	}

	if respData.TTL != nil {
		_ = d.Set("ttl", respData.TTL)
	}

	if respData.Weight != nil {
		_ = d.Set("weight", respData.Weight)
	}

	if respData.Priority != nil {
		_ = d.Set("priority", respData.Priority)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.CreatedOn != nil {
		_ = d.Set("created_on", respData.CreatedOn)
	}

	if respData.ModifiedOn != nil {
		_ = d.Set("modified_on", respData.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordV10Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v10.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	needChange := false
	mutableArgs := []string{"content", "location", "ttl", "weight", "priority"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifyDnsRecordsRequest()
		request.ZoneId = helper.String(zoneId)
		dnsRecord := &teov20220901.DnsRecord{
			RecordId: helper.String(recordId),
		}

		if v, ok := d.GetOk("content"); ok {
			dnsRecord.Content = helper.String(v.(string))
		}

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
		request.DnsRecords = []*teov20220901.DnsRecord{dnsRecord}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyDnsRecordsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo dns record v10 failed, reason:%+v", logId, err)
			return err
		}

		// Wait for DNS record update to become active
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			respData, e := service.DescribeTeoDnsRecordById(ctx, zoneId, recordId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			if respData == nil {
				return resource.RetryableError(fmt.Errorf("DNS record %s not found", recordId))
			}
			if respData.Status == nil {
				return resource.RetryableError(fmt.Errorf("DNS record %s status not available", recordId))
			}
			if *respData.Status != "enable" {
				return resource.RetryableError(fmt.Errorf("DNS record %s is not active yet, status: %s", recordId, *respData.Status))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s wait for teo dns record v10 update to become active failed, reason:%+v", logId, err)
			return fmt.Errorf("waiting for DNS record update to become active failed: %v", err)
		}
	}

	return resourceTencentCloudTeoDnsRecordV10Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV10Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v10.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	var (
		request  = teov20220901.NewDeleteDnsRecordsRequest()
		response = teov20220901.NewDeleteDnsRecordsResponse()
	)

	request.ZoneId = helper.String(zoneId)
	request.RecordIds = helper.Strings([]string{recordId})

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo dns record v10 failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
