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
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoDnsRecordV9() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordV9Create,
		Read:   resourceTencentCloudTeoDnsRecordV9Read,
		Update: resourceTencentCloudTeoDnsRecordV9Update,
		Delete: resourceTencentCloudTeoDnsRecordV9Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name. For Chinese, Korean, or Japanese domain names, convert to punycode before input.",
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.",
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "MX", "CNAME", "TXT", "NS", "CAA", "SRV"}, false),
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. Fill in content corresponding to the Type value. For Chinese, Korean, or Japanese domain names, convert to punycode before input.",
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS record resolution line. Default is Default, meaning default resolution line and effective for all regions. Only applicable when Type is A, AAAA, or CNAME.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					recordType := d.Get("type").(string)
					if recordType != "A" && recordType != "AAAA" && recordType != "CNAME" {
						return true
					}
					return false
				},
			},

			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Cache time, range 60-86400 seconds. Smaller value means faster propagation of changes. Default is 300.",
				ValidateFunc: validation.IntBetween(60, 86400),
			},

			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "DNS record weight, range -1 to 100. -1 means no weight set, 0 means no resolution. Only applicable when Type is A, AAAA, or CNAME.",
				ValidateFunc: validation.IntBetween(-1, 100),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					recordType := d.Get("type").(string)
					if recordType != "A" && recordType != "AAAA" && recordType != "CNAME" {
						return true
					}
					return false
				},
			},

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "MX record priority, range 0-50. Lower value means higher priority. Only applicable when Type is MX.",
				ValidateFunc: validation.IntBetween(0, 50),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					recordType := d.Get("type").(string)
					if recordType != "MX" {
						return true
					}
					return false
				},
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record resolution status. Valid values: enable (active), disable (stopped).",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the DNS record.",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification time of the DNS record.",
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecordV9Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v9.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	zoneId := d.Get("zone_id").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)
	content := d.Get("content").(string)

	request := teo.NewCreateDnsRecordRequest()
	request.ZoneId = helper.String(zoneId)
	request.Name = helper.String(name)
	request.Type = helper.String(recordType)
	request.Content = helper.String(content)

	if v, ok := d.GetOk("location"); ok {
		recordType := d.Get("type").(string)
		if recordType == "A" || recordType == "AAAA" || recordType == "CNAME" {
			request.Location = helper.String(v.(string))
		}
	}

	if v, ok := d.GetOk("ttl"); ok {
		request.TTL = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("weight"); ok {
		recordType := d.Get("type").(string)
		if recordType == "A" || recordType == "AAAA" || recordType == "CNAME" {
			request.Weight = helper.IntInt64(v.(int))
		}
	}

	if v, ok := d.GetOk("priority"); ok {
		recordType := d.Get("type").(string)
		if recordType == "MX" {
			request.Priority = helper.IntInt64(v.(int))
		}
	}

	var response *teo.CreateDnsRecordResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateDnsRecordWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo DNS record failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s create teo DNS record failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.RecordId == nil || *response.Response.RecordId == "" {
		return fmt.Errorf("create teo DNS record failed, record_id is empty")
	}

	recordId := *response.Response.RecordId
	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	// Wait for the record to be created
	err = resourceTencentCloudTeoDnsRecordV9Read(d, meta)
	if err != nil {
		return fmt.Errorf("wait for teo DNS record to be created failed: %s", err)
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordV9Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v9.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("invalid resource ID format: %s, expected format: zone_id#record_id", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	// Query DNS records with record ID filter
	request := teo.NewDescribeDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.Limit = helper.IntInt64(1)
	request.Filters = []*teo.AdvancedFilter{
		{
			Name:   helper.String("id"),
			Values: []*string{helper.String(recordId)},
		},
	}

	var response *teo.DescribeDnsRecordsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo DNS records failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s read teo DNS record failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.TotalCount == nil || *response.Response.TotalCount == 0 {
		d.SetId("")
		log.Printf("[WARN]%s teo DNS record not found, zone_id=%s, record_id=%s", logId, zoneId, recordId)
		return nil
	}

	if len(response.Response.DnsRecords) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s teo DNS record list is empty", logId)
		return nil
	}

	dnsRecord := response.Response.DnsRecords[0]

	if dnsRecord.ZoneId != nil {
		_ = d.Set("zone_id", *dnsRecord.ZoneId)
	}

	if dnsRecord.Name != nil {
		_ = d.Set("name", *dnsRecord.Name)
	}

	if dnsRecord.Type != nil {
		_ = d.Set("type", *dnsRecord.Type)
	}

	if dnsRecord.Content != nil {
		_ = d.Set("content", *dnsRecord.Content)
	}

	if dnsRecord.Location != nil {
		_ = d.Set("location", *dnsRecord.Location)
	}

	if dnsRecord.TTL != nil {
		_ = d.Set("ttl", *dnsRecord.TTL)
	}

	if dnsRecord.Weight != nil {
		_ = d.Set("weight", *dnsRecord.Weight)
	}

	if dnsRecord.Priority != nil {
		_ = d.Set("priority", *dnsRecord.Priority)
	}

	if dnsRecord.Status != nil {
		_ = d.Set("status", *dnsRecord.Status)
	}

	if dnsRecord.CreatedOn != nil {
		_ = d.Set("created_on", *dnsRecord.CreatedOn)
	}

	if dnsRecord.ModifiedOn != nil {
		_ = d.Set("modified_on", *dnsRecord.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordV9Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v9.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("invalid resource ID format: %s, expected format: zone_id#record_id", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	// First read the current record to get all fields
	dnsRecord, err := readTeoDnsRecord(ctx, meta, zoneId, recordId, logId)
	if err != nil {
		return err
	}
	if dnsRecord == nil {
		return fmt.Errorf("teo DNS record not found: %s", d.Id())
	}

	// Update only the fields that have changed

	if d.HasChange("name") {
		dnsRecord.Name = helper.String(d.Get("name").(string))
	}

	if d.HasChange("type") {
		dnsRecord.Type = helper.String(d.Get("type").(string))
	}

	if d.HasChange("content") {
		dnsRecord.Content = helper.String(d.Get("content").(string))
	}

	if d.HasChange("location") {
		recordType := d.Get("type").(string)
		if recordType == "A" || recordType == "AAAA" || recordType == "CNAME" {
			if v, ok := d.GetOk("location"); ok {
				dnsRecord.Location = helper.String(v.(string))
			} else {
				dnsRecord.Location = nil
			}
		}
	}

	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			dnsRecord.TTL = helper.IntInt64(v.(int))
		} else {
			dnsRecord.TTL = nil
		}
	}

	if d.HasChange("weight") {
		recordType := d.Get("type").(string)
		if recordType == "A" || recordType == "AAAA" || recordType == "CNAME" {
			if v, ok := d.GetOk("weight"); ok {
				dnsRecord.Weight = helper.IntInt64(v.(int))
			} else {
				dnsRecord.Weight = nil
			}
		}
	}

	if d.HasChange("priority") {
		recordType := d.Get("type").(string)
		if recordType == "MX" {
			if v, ok := d.GetOk("priority"); ok {
				dnsRecord.Priority = helper.IntInt64(v.(int))
			} else {
				dnsRecord.Priority = nil
			}
		}
	}

	request := teo.NewModifyDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.DnsRecords = []*teo.DnsRecord{dnsRecord}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify teo DNS record failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s update teo DNS record failed, reason:%+v", logId, err)
		return err
	}

	// Wait for the record to be updated
	err = resourceTencentCloudTeoDnsRecordV9Read(d, meta)
	if err != nil {
		return fmt.Errorf("wait for teo DNS record to be updated failed: %s", err)
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordV9Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v9.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("invalid resource ID format: %s, expected format: zone_id#record_id", d.Id())
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
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo DNS record failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s delete teo DNS record failed, reason:%+v", logId, err)
		return err
	}

	// Wait for the record to be deleted
	err = resourceTencentCloudTeoDnsRecordV9Read(d, meta)
	if err == nil && d.Id() != "" {
		return fmt.Errorf("wait for teo DNS record to be deleted failed")
	}

	return nil
}

// readTeoDnsRecord reads a DNS record by zone ID and record ID
func readTeoDnsRecord(ctx context.Context, meta interface{}, zoneId, recordId, logId string) (*teo.DnsRecord, error) {
	request := teo.NewDescribeDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.Limit = helper.IntInt64(1)
	request.Filters = []*teo.AdvancedFilter{
		{
			Name:   helper.String("id"),
			Values: []*string{helper.String(recordId)},
		},
	}

	var response *teo.DescribeDnsRecordsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo DNS records failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		return nil, err
	}

	if response.Response.TotalCount == nil || *response.Response.TotalCount == 0 {
		return nil, nil
	}

	if len(response.Response.DnsRecords) == 0 {
		return nil, nil
	}

	return response.Response.DnsRecords[0], nil
}
