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

func ResourceTencentCloudTeoDnsRecord10() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecord10Create,
		Read:   resourceTencentCloudTeoDnsRecord10Read,
		Update: resourceTencentCloudTeoDnsRecord10Update,
		Delete: resourceTencentCloudTeoDnsRecord10Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "DNS record type. Valid values are:\n" +
					"	- A: points domain name to an external IPv4 address, such as 8.8.8.8;\n" +
					"	- AAAA: points domain name to an external IPv6 address;\n" +
					"	- MX: used for email servers. When there are multiple MX records, lower priority value means higher priority;\n" +
					"	- CNAME: points domain name to another domain name, which then resolves to the final IP address;\n" +
					"	- TXT: identifies and describes domain name, commonly used for domain verification and SPF records (anti-spam);\n" +
					"	- NS: if you need to delegate subdomain to another DNS service provider for resolution, you need to add an NS record. Root domain cannot add NS records;\n" +
					"	- CAA: specifies CA that can issue certificates for this site;\n" +
					"	- SRV: identifies a server using a service, commonly used in Microsoft's directory management.\n" +
					"Different record types, such as SRV and CAA records, have different requirements for host record names and record value formats. For detailed descriptions and format examples of each record type, please refer to: [Introduction to DNS record types](https://cloud.tencent.com/document/product/1552/90453#2f681022-91ab-4a9e-ac3d-0a6c454d954e).",
				ValidateFunc: validateDnsRecordType,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. Fill in corresponding content according to the type value. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Default:     "Default",
				Description: "DNS record resolution route. If not specified, default is Default, which means the default resolution route and takes effect in all regions.\n\n- Resolution route configuration is only applicable when type (DNS record type) is A, AAAA, or CNAME.\n- Resolution route configuration is only applicable to standard version and enterprise edition packages. For valid values, please refer to: [Resolution routes and corresponding code enumeration](https://cloud.tencent.com/document/product/1552/112542).",
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Default:      300,
				Description:  "Cache time. Users can specify a value range of 60-86400. Smaller value means faster modification records will take effect in all regions. Default value: 300. Unit: seconds.",
				ValidateFunc: validateTTL,
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Default:      -1,
				Description:  "DNS record weight. Users can specify a value range of -1 to 100. A value of 0 means no resolution. If not specified, default is -1, which means no weight is set. Weight configuration is only applicable when type (DNS record type) is A, AAAA, or CNAME. Note: For the same subdomain, different DNS records with the same resolution route should either all have weights set or none have weights set.",
				ValidateFunc: validateWeight,
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Default:      0,
				Description:  "MX record priority, which takes effect only when type (DNS record type) is MX. Smaller value means higher priority. Users can specify a value range of 0-50. Default value is 0 if not specified.",
				ValidateFunc: validatePriority,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "DNS record resolution status. Valid values:\n" +
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
				Description: "Modification time.",
			},
		},
	}
}

func validateDnsRecordType(v interface{}, k string) (ws []string, es []error) {
	validTypes := []string{"A", "AAAA", "MX", "CNAME", "TXT", "NS", "CAA", "SRV"}
	value := v.(string)
	valid := false
	for _, t := range validTypes {
		if value == t {
			valid = true
			break
		}
	}
	if !valid {
		es = append(es, fmt.Errorf("%q must be one of: %v, got %q", k, validTypes, value))
	}
	return
}

func validateTTL(v interface{}, k string) (ws []string, es []error) {
	value := v.(int)
	if value < 60 || value > 86400 {
		es = append(es, fmt.Errorf("%q must be in the range 60-86400, got %d", k, value))
	}
	return
}

func validateWeight(v interface{}, k string) (ws []string, es []error) {
	value := v.(int)
	if value < -1 || value > 100 {
		es = append(es, fmt.Errorf("%q must be in the range -1-100, got %d", k, value))
	}
	return
}

func validatePriority(v interface{}, k string) (ws []string, es []error) {
	value := v.(int)
	if value < 0 || value > 50 {
		es = append(es, fmt.Errorf("%q must be in the range 0-50, got %d", k, value))
	}
	return
}

func parseDnsRecordId(id string) (zoneId, recordId string, err error) {
	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		err = fmt.Errorf("invalid ID format, must be zone_id%srecord_id", tccommon.FILED_SP)
		return
	}
	zoneId = idSplit[0]
	recordId = idSplit[1]
	return
}

func buildDnsRecordId(zoneId, recordId string) string {
	return strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP)
}

func resourceTencentCloudTeoDnsRecord10Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_10.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId   string
		recordId string
	)

	request := teov20220901.NewCreateDnsRecordRequest()

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
		recordId = *result.Response.RecordId
		return nil
	})
	if err != nil {
		log.Printf("[CRITICAL]%s create teo dns record failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(buildDnsRecordId(zoneId, recordId))

	return resourceTencentCloudTeoDnsRecord10Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord10Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_10.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	zoneId, recordId, err := parseDnsRecordId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.CallDescribeDnsRecords(ctx, meta, zoneId, recordId)
	if err != nil {
		return err
	}
	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_dns_record_10` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.ZoneId != nil {
		_ = d.Set("zone_id", respData.ZoneId)
	}

	if respData.RecordId != nil {
		_ = d.Set("record_id", respData.RecordId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.Content != nil {
		_ = d.Set("content", respData.Content)
	}

	if respData.Location != nil {
		_ = d.Set("location", respData.Location)
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

func resourceTencentCloudTeoDnsRecord10Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_10.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId, recordId, err := parseDnsRecordId(d.Id())
	if err != nil {
		return err
	}

	mutableArgs := []string{"name", "type", "content", "location", "ttl", "weight", "priority"}
	needChange := false
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
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyDnsRecordsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITICAL]%s update teo dns record failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoDnsRecord10Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord10Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_10.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId, recordId, err := parseDnsRecordId(d.Id())
	if err != nil {
		return err
	}

	request := teov20220901.NewDeleteDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.RecordIds = helper.Strings([]string{recordId})

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITICAL]%s delete teo dns record failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
