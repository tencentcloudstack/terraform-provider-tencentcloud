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
			Read:   schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the site related with the DNS record.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "DNS record name. If it is a Chinese, Korean, or Japanese domain name, it needs to be converted to punycode before input.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "DNS record content. Fill in the corresponding content according to the Type value. If it is a Chinese, Korean, or Japanese domain name, it needs to be converted to punycode before input.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "DNS record resolution line. Default is Default, which means default resolution line and represents all regions. Valid only when Type is A, AAAA, or CNAME.",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time. Valid range: 60-86400, unit: seconds. Default is 300.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. Valid range: -1-100. Default is -1, which means no weight is set. When set to 0, it means no resolution. Valid only when Type is A, AAAA, or CNAME.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority. Valid only when Type is MX. Valid range: 0-50, smaller value, higher priority. Default is 0.",
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
				Description: "Creation time of the DNS record.",
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecord10Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_10.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)
	content := d.Get("content").(string)

	// Check for idempotency: query existing records to avoid creating duplicates
	existingRecordId, err := resourceTencentCloudTeoDnsRecord10FindExistingRecord(meta, logId, zoneId, name, recordType, content)
	if err != nil {
		return fmt.Errorf("error checking for existing DNS record: %s", err)
	}

	if existingRecordId != "" {
		log.Printf("[INFO] DNS record already exists with RecordId: %s", existingRecordId)
		d.SetId(fmt.Sprintf("%s#%s", zoneId, existingRecordId))
		return resourceTencentCloudTeoDnsRecord10Read(d, meta)
	}

	// Build create request
	request := teo.NewCreateDnsRecordRequest()
	request.ZoneId = &zoneId
	request.Name = &name
	request.Type = &recordType
	request.Content = &content

	if v, ok := d.GetOk("location"); ok {
		request.Location = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ttl"); ok {
		request.TTL = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("weight"); ok {
		request.Weight = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	// Call API to create DNS record
	teoService := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()
	response, err := teoService.CreateDnsRecordWithContext(context.TODO(), request)
	if err != nil {
		log.Printf("[CRITAL]%s create DNS record failed, request=%s, response=%s, reason=%s", logId, request.ToJsonString(), response.ToJsonString(), err.Error())
		return err
	}

	recordId := *response.Response.RecordId
	d.SetId(fmt.Sprintf("%s#%s", zoneId, recordId))

	// Wait for DNS record to be available
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := resourceTencentCloudTeoDnsRecord10FindRecordById(meta, logId, zoneId, recordId)
		if err != nil {
			if isRecordNotFoundError(err) {
				return resource.RetryableError(fmt.Errorf("waiting for DNS record to be created: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for DNS record creation: %s", err)
	}

	return resourceTencentCloudTeoDnsRecord10Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord10Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_10.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	id := d.Id()
	zoneId, recordId, err := resourceTencentCloudTeoDnsRecord10ParseId(id)
	if err != nil {
		return fmt.Errorf("error parsing DNS record ID: %s", err)
	}

	record, err := resourceTencentCloudTeoDnsRecord10FindRecordById(meta, logId, zoneId, recordId)
	if err != nil {
		if isRecordNotFoundError(err) {
			log.Printf("[WARN] DNS record %s does not exist, removing from state", id)
			d.SetId("")
			return nil
		}
		return err
	}

	// Set resource attributes
	_ = d.Set("zone_id", zoneId)
	_ = d.Set("name", record.Name)
	_ = d.Set("type", record.Type)
	_ = d.Set("content", record.Content)
	_ = d.Set("location", record.Location)
	_ = d.Set("ttl", record.TTL)
	_ = d.Set("weight", record.Weight)
	_ = d.Set("priority", record.Priority)
	_ = d.Set("record_id", record.RecordId)
	_ = d.Set("status", record.Status)
	_ = d.Set("created_on", record.CreatedOn)

	return nil
}

func resourceTencentCloudTeoDnsRecord10Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_10.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	id := d.Id()
	zoneId, recordId, err := resourceTencentCloudTeoDnsRecord10ParseId(id)
	if err != nil {
		return fmt.Errorf("error parsing DNS record ID: %s", err)
	}

	if !d.HasChange("name") && !d.HasChange("type") && !d.HasChange("content") && !d.HasChange("location") && !d.HasChange("ttl") && !d.HasChange("weight") && !d.HasChange("priority") {
		log.Printf("[INFO] No changes detected for DNS record %s", id)
		return nil
	}

	// Build update request
	dnsRecord := &teo.DnsRecord{
		RecordId: &recordId,
	}

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
		dnsRecord.Location = helper.String(d.Get("location").(string))
	}

	if d.HasChange("ttl") {
		dnsRecord.TTL = helper.IntInt64(d.Get("ttl").(int))
	}

	if d.HasChange("weight") {
		dnsRecord.Weight = helper.IntInt64(d.Get("weight").(int))
	}

	if d.HasChange("priority") {
		dnsRecord.Priority = helper.IntInt64(d.Get("priority").(int))
	}

	request := teo.NewModifyDnsRecordsRequest()
	request.ZoneId = &zoneId
	request.DnsRecords = []*teo.DnsRecord{dnsRecord}

	// Call API to update DNS record
	teoService := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()
	response, err := teoService.ModifyDnsRecordsWithContext(context.TODO(), request)
	if err != nil {
		log.Printf("[CRITAL]%s update DNS record failed, request=%s, response=%s, reason=%s", logId, request.ToJsonString(), response.ToJsonString(), err.Error())
		return err
	}

	// Wait for DNS record update to be effective
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		currentRecord, err := resourceTencentCloudTeoDnsRecord10FindRecordById(meta, logId, zoneId, recordId)
		if err != nil {
			if isRecordNotFoundError(err) {
				return resource.NonRetryableError(fmt.Errorf("DNS record %s not found", id))
			}
			return resource.RetryableError(fmt.Errorf("waiting for DNS record update: %s", err))
		}

		// Check if the values have been updated
		if d.HasChange("name") && *currentRecord.Name != d.Get("name").(string) {
			return resource.RetryableError(fmt.Errorf("waiting for name to be updated"))
		}
		if d.HasChange("type") && *currentRecord.Type != d.Get("type").(string) {
			return resource.RetryableError(fmt.Errorf("waiting for type to be updated"))
		}
		if d.HasChange("content") && *currentRecord.Content != d.Get("content").(string) {
			return resource.RetryableError(fmt.Errorf("waiting for content to be updated"))
		}
		if d.HasChange("location") {
			currentLocation := "Default"
			if currentRecord.Location != nil {
				currentLocation = *currentRecord.Location
			}
			expectedLocation := d.Get("location").(string)
			if expectedLocation != "" && currentLocation != expectedLocation {
				return resource.RetryableError(fmt.Errorf("waiting for location to be updated"))
			}
		}
		if d.HasChange("ttl") && currentRecord.TTL != nil && *currentRecord.TTL != int64(d.Get("ttl").(int)) {
			return resource.RetryableError(fmt.Errorf("waiting for TTL to be updated"))
		}
		if d.HasChange("weight") && currentRecord.Weight != nil && *currentRecord.Weight != int64(d.Get("weight").(int)) {
			return resource.RetryableError(fmt.Errorf("waiting for weight to be updated"))
		}
		if d.HasChange("priority") && currentRecord.Priority != nil && *currentRecord.Priority != int64(d.Get("priority").(int)) {
			return resource.RetryableError(fmt.Errorf("waiting for priority to be updated"))
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for DNS record update: %s", err)
	}

	return resourceTencentCloudTeoDnsRecord10Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord10Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_10.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	id := d.Id()
	zoneId, recordId, err := resourceTencentCloudTeoDnsRecord10ParseId(id)
	if err != nil {
		return fmt.Errorf("error parsing DNS record ID: %s", err)
	}

	// Build delete request
	request := teo.NewDeleteDnsRecordsRequest()
	request.ZoneId = &zoneId
	request.RecordIds = []*string{&recordId}

	// Call API to delete DNS record
	teoService := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()
	response, err := teoService.DeleteDnsRecordsWithContext(context.TODO(), request)
	if err != nil {
		log.Printf("[CRITAL]%s delete DNS record failed, request=%s, response=%s, reason=%s", logId, request.ToJsonString(), response.ToJsonString(), err.Error())
		return err
	}

	d.SetId("")
	return nil
}

// Helper functions

func resourceTencentCloudTeoDnsRecord10ParseId(id string) (zoneId, recordId string, err error) {
	parts := strings.Split(id, "#")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid DNS record ID format: %s, expected format: zoneId#recordId", id)
	}
	return parts[0], parts[1], nil
}

func resourceTencentCloudTeoDnsRecord10FindExistingRecord(meta interface{}, logId, zoneId, name, recordType, content string) (string, error) {
	teoService := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()

	request := teo.NewDescribeDnsRecordsRequest()
	request.ZoneId = &zoneId
	request.Limit = helper.IntInt64(1000)

	// Add filters to find existing record with same name, type, and content
	filters := []*teo.AdvancedFilter{
		{Name: helper.String("name"), Values: []*string{&name}},
		{Name: helper.String("type"), Values: []*string{&recordType}},
		{Name: helper.String("content"), Values: []*string{&content}},
	}
	request.Filters = filters

	response, err := teoService.DescribeDnsRecordsWithContext(context.TODO(), request)
	if err != nil {
		return "", err
	}

	if response.Response != nil && len(response.Response.DnsRecords) > 0 {
		return *response.Response.DnsRecords[0].RecordId, nil
	}

	return "", nil
}

func resourceTencentCloudTeoDnsRecord10FindRecordById(meta interface{}, logId, zoneId, recordId string) (*teo.DnsRecord, error) {
	teoService := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()

	request := teo.NewDescribeDnsRecordsRequest()
	request.ZoneId = &zoneId
	request.Limit = helper.IntInt64(1000)

	// Add filter to find record by ID
	filters := []*teo.AdvancedFilter{
		{Name: helper.String("id"), Values: []*string{&recordId}},
	}
	request.Filters = filters

	response, err := teoService.DescribeDnsRecordsWithContext(context.TODO(), request)
	if err != nil {
		return nil, err
	}

	if response.Response == nil || len(response.Response.DnsRecords) == 0 {
		return nil, fmt.Errorf("DNS record not found")
	}

	return response.Response.DnsRecords[0], nil
}

func isRecordNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "does not exist") ||
		strings.Contains(err.Error(), "ResourceNotFound")
}
