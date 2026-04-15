package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoDnsRecordV13() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordV13Create,
		Read:   resourceTencentCloudTeoDnsRecordV13Read,
		Update: resourceTencentCloudTeoDnsRecordV13Update,
		Delete: resourceTencentCloudTeoDnsRecordV13Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				Description: "DNS record name. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.",
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CNAME", "MX", "TXT", "NS", "CAA", "SRV"}, false),
				Description: "DNS record type. valid values are:\n" +
					"	- A: points the domain name to an external ipv4 address, such as 8.8.8.8;\n" +
					"	- AAAA: points the domain name to an external ipv6 address;\n" +
					"	- MX: used for email servers. when there are multiple mx records, the lower the priority value, the higher the priority;\n" +
					"	- CNAME: points the domain name to another domain name, which then resolves to the final ip address;\n" +
					"	- TXT: identifies and describes the domain name, commonly used for domain verification and spf records (anti-spam);\n" +
					"	- NS: if you need to delegate the subdomain to another dns service provider for resolution, you need to add an ns record. the root domain cannot add ns records;\n" +
					"	- CAA: specifies the ca that can issue certificates for this site;\n" +
					"	- SRV: identifies a server using a service, commonly used in microsoft's directory management.\n" +
					"Different record types, such as SRV and CAA records, have different requirements for host record names and record value formats. for detailed descriptions and format examples of each record type, please refer to: [introduction to dns record types](https://intl.cloud.tencent.com/document/product/1552/90453?from_cn_redirect=1#2f681022-91ab-4a9e-ac3d-0a6c454d954e).",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. fill in the corresponding content according to the type value. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.",
			},

			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "DNS record resolution route. if not specified, the default is DEFAULT, which means the default resolution route and is effective in all regions.\n\n- resolution route configuration is only applicable when type (dns record type) is A, AAAA, or CNAME.\n- resolution route configuration is only applicable to standard version and enterprise edition packages. for valid values, please refer to: [resolution routes and corresponding code enumeration](https://intl.cloud.tencent.com/document/product/1552/112542?from_cn_redirect=1).",
			},

			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(60, 86400),
				Description:  "Cache time. users can specify a value range of 60-86400. the smaller the value, the faster the modification records will take effect in all regions. default value: 300. unit: seconds.",
			},

			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(-1, 100),
				Description:  "DNS record weight. users can specify a value range of -1 to 100. a value of 0 means no resolution. if not specified, the default is -1, which means no weight is set. weight configuration is only applicable when type (dns record type) is A, AAAA, or CNAME. note: for the same subdomain, different dns records with the same resolution route should either all have weights set or none have weights set.",
			},

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 50),
				Description:  "MX record priority, which takes effect only when type (dns record type) is MX. the smaller the value, the higher the priority. users can specify a value range of 0-50. the default value is 0 if not specified.",
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "DNS record resolution status, the following values:\n" +
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
	}
}

func resourceTencentCloudTeoDnsRecordV13Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v13.create")()
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
		log.Printf("[CRITAL]%s create teo dns record v13 failed, reason:%+v", logId, err)
		return err
	}

	recordId = *response.Response.RecordId

	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	// Wait for DNS record to be created and enabled
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		record, e := resourceTencentCloudTeoDnsRecordV13DescribeRecordById(ctx, meta, zoneId, recordId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if record == nil {
			return resource.NonRetryableError(fmt.Errorf("DNS record not found"))
		}

		if record.Status == nil {
			return resource.RetryableError(fmt.Errorf("DNS record status is nil"))
		}

		if *record.Status != "enable" {
			return resource.RetryableError(fmt.Errorf("DNS record status is %s, waiting for enable", *record.Status))
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to wait for DNS record to be enabled: %v", err)
	}

	return resourceTencentCloudTeoDnsRecordV13Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV13Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v13.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	record, err := resourceTencentCloudTeoDnsRecordV13DescribeRecordById(ctx, meta, zoneId, recordId)
	if err != nil {
		return err
	}

	if record == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("record_id", recordId)
	if record.Name != nil {
		_ = d.Set("name", *record.Name)
	}
	if record.Type != nil {
		_ = d.Set("type", *record.Type)
	}
	if record.Content != nil {
		_ = d.Set("content", *record.Content)
	}
	if record.Location != nil {
		_ = d.Set("location", *record.Location)
	}
	if record.TTL != nil {
		_ = d.Set("ttl", int(*record.TTL))
	}
	if record.Weight != nil {
		_ = d.Set("weight", int(*record.Weight))
	}
	if record.Priority != nil {
		_ = d.Set("priority", int(*record.Priority))
	}
	if record.Status != nil {
		_ = d.Set("status", *record.Status)
	}
	if record.CreatedOn != nil {
		_ = d.Set("created_on", *record.CreatedOn)
	}
	if record.ModifiedOn != nil {
		_ = d.Set("modified_on", *record.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordV13Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v13.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	if d.HasChange("name") || d.HasChange("type") || d.HasChange("content") || d.HasChange("location") || d.HasChange("ttl") || d.HasChange("weight") || d.HasChange("priority") {
		request := teov20220901.NewModifyDnsRecordsRequest()

		request.ZoneId = helper.String(zoneId)

		dnsRecord := &teov20220901.DnsRecord{
			RecordId: helper.String(recordId),
		}

		if v, ok := d.GetOk("name"); ok {
			dnsRecord.Name = helper.String(v.(string))
		}
		if v, ok := d.GetOk("type"); ok {
			dnsRecord.Type = helper.String(v.(string))
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
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyDnsRecordsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify teo dns record v13 failed, reason:%+v", logId, err)
			return err
		}

		// Wait for DNS record to be updated
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			record, e := resourceTencentCloudTeoDnsRecordV13DescribeRecordById(ctx, meta, zoneId, recordId)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if record == nil {
				return resource.NonRetryableError(fmt.Errorf("DNS record not found"))
			}

			// Check if the record has been updated with the new values
			if d.HasChange("name") && record.Name != nil && *record.Name != d.Get("name") {
				return resource.RetryableError(fmt.Errorf("DNS record name not updated yet"))
			}
			if d.HasChange("content") && record.Content != nil && *record.Content != d.Get("content") {
				return resource.RetryableError(fmt.Errorf("DNS record content not updated yet"))
			}
			if d.HasChange("ttl") && record.TTL != nil && int(*record.TTL) != d.Get("ttl") {
				return resource.RetryableError(fmt.Errorf("DNS record TTL not updated yet"))
			}
			if d.HasChange("weight") && record.Weight != nil && int(*record.Weight) != d.Get("weight") {
				return resource.RetryableError(fmt.Errorf("DNS record weight not updated yet"))
			}
			if d.HasChange("priority") && record.Priority != nil && int(*record.Priority) != d.Get("priority") {
				return resource.RetryableError(fmt.Errorf("DNS record priority not updated yet"))
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to wait for DNS record to be updated: %v", err)
		}
	}

	return resourceTencentCloudTeoDnsRecordV13Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV13Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v13.delete")()
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
		request = teov20220901.NewDeleteDnsRecordsRequest()
	)

	request.ZoneId = helper.String(zoneId)
	request.RecordIds = []*string{helper.String(recordId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo dns record v13 failed, reason:%+v", logId, err)
		return err
	}

	// Wait for DNS record to be deleted
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		record, e := resourceTencentCloudTeoDnsRecordV13DescribeRecordById(ctx, meta, zoneId, recordId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if record == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("DNS record still exists, waiting for deletion"))
	})

	if err != nil {
		return fmt.Errorf("failed to wait for DNS record to be deleted: %v", err)
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordV13DescribeRecordById(ctx context.Context, meta interface{}, zoneId, recordId string) (*teov20220901.DnsRecord, error) {
	logId := tccommon.GetLogId(ctx)

	var (
		request  = teov20220901.NewDescribeDnsRecordsRequest()
		response = teov20220901.NewDescribeDnsRecordsResponse()
	)

	request.ZoneId = helper.String(zoneId)
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("id"),
			Values: []*string{helper.String(recordId)},
		},
	}
	request.Limit = helper.IntInt64(1)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s describe teo dns record v13 failed, reason:%+v", logId, err)
		return nil, err
	}

	if response == nil || response.Response == nil || response.Response.DnsRecords == nil || len(response.Response.DnsRecords) == 0 {
		return nil, nil
	}

	return response.Response.DnsRecords[0], nil
}
