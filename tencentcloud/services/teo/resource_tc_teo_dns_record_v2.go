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

func ResourceTencentCloudTeoDnsRecordV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordV2Create,
		Read:   resourceTencentCloudTeoDnsRecordV2Read,
		Update: resourceTencentCloudTeoDnsRecordV2Update,
		Delete: resourceTencentCloudTeoDnsRecordV2Delete,
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
				Description: "Site ID.",
			},
			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.",
				ValidateFunc: validateDnsRecordType,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS record resolution line. Default: Default. Only applicable when Type is A, AAAA, or CNAME.",
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Cache time in seconds. Range: 60-86400. Default: 300.",
				ValidateFunc: validateDnsRecordTTL,
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "DNS record weight. Range: -1~100. Default: -1. Only applicable when Type is A, AAAA, or CNAME.",
				ValidateFunc: validateDnsRecordWeight,
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "MX record priority. Range: 0~50. Default: 0. Only applicable when Type is MX.",
				ValidateFunc: validateDnsRecordPriority,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record status. Valid values: enable, disable.",
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

func resourceTencentCloudTeoDnsRecordV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

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
		request.Location = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ttl"); ok {
		request.TTL = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("weight"); ok {
		request.Weight = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("priority"); ok {
		request.Priority = helper.Int64(int64(v.(int)))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateDnsRecordWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo dns record failed, Response is nil."))
		}

		recordId := result.Response.RecordId
		if recordId == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo dns record failed, record ID is nil."))
		}

		d.SetId(fmt.Sprintf("%s#%s", zoneId, *recordId))
		d.Set("record_id", *recordId)

		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s create teo dns record failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoDnsRecordV2Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId, recordId, err := resourceTencentCloudTeoDnsRecordV2ParseId(d.Id())
	if err != nil {
		return err
	}

	request := teo.NewDescribeDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.Limit = helper.Int64(20)

	// Filter by record ID
	request.Filters = []*teo.AdvancedFilter{
		{
			Name:   helper.String("id"),
			Values: []*string{&recordId},
		},
	}

	var response *teo.DescribeDnsRecordsResponse
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo dns records failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s describe teo dns records failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.DnsRecords) == 0 {
		log.Printf("[DEBUG]%s teo dns record not found, id:%s", logId, d.Id())
		d.SetId("")
		return nil
	}

	dnsRecord := response.Response.DnsRecords[0]
	_ = d.Set("zone_id", zoneId)
	_ = d.Set("record_id", recordId)
	_ = d.Set("name", dnsRecord.Name)
	_ = d.Set("type", dnsRecord.Type)
	_ = d.Set("content", dnsRecord.Content)
	_ = d.Set("location", dnsRecord.Location)
	_ = d.Set("ttl", dnsRecord.TTL)
	_ = d.Set("weight", dnsRecord.Weight)
	_ = d.Set("priority", dnsRecord.Priority)
	_ = d.Set("status", dnsRecord.Status)
	_ = d.Set("created_on", dnsRecord.CreatedOn)
	_ = d.Set("modified_on", dnsRecord.ModifiedOn)

	return nil
}

func resourceTencentCloudTeoDnsRecordV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId, recordId, err := resourceTencentCloudTeoDnsRecordV2ParseId(d.Id())
	if err != nil {
		return err
	}

	// Read current record data first
	currentRecord, err := readDnsRecord(meta, zoneId, recordId)
	if err != nil {
		return err
	}

	if currentRecord == nil {
		return fmt.Errorf("DNS record not found, zone_id: %s, record_id: %s", zoneId, recordId)
	}

	// Check if any field has changed
	if !d.HasChange("name") && !d.HasChange("type") && !d.HasChange("content") &&
		!d.HasChange("location") && !d.HasChange("ttl") &&
		!d.HasChange("weight") && !d.HasChange("priority") {
		log.Printf("[DEBUG]%s no changes detected for teo dns record", logId)
		return resourceTencentCloudTeoDnsRecordV2Read(d, meta)
	}

	// Build update request with current values and updates
	dnsRecord := teo.DnsRecord{
		RecordId: helper.String(recordId),
	}

	// Use updated values or current values
	if d.HasChange("name") {
		dnsRecord.Name = helper.String(d.Get("name").(string))
	} else {
		dnsRecord.Name = currentRecord.Name
	}

	if d.HasChange("type") {
		dnsRecord.Type = helper.String(d.Get("type").(string))
	} else {
		dnsRecord.Type = currentRecord.Type
	}

	if d.HasChange("content") {
		dnsRecord.Content = helper.String(d.Get("content").(string))
	} else {
		dnsRecord.Content = currentRecord.Content
	}

	if d.HasChange("location") {
		if v, ok := d.GetOk("location"); ok {
			dnsRecord.Location = helper.String(v.(string))
		}
	} else {
		dnsRecord.Location = currentRecord.Location
	}

	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			dnsRecord.TTL = helper.Int64(int64(v.(int)))
		}
	} else {
		dnsRecord.TTL = currentRecord.TTL
	}

	if d.HasChange("weight") {
		if v, ok := d.GetOk("weight"); ok {
			dnsRecord.Weight = helper.Int64(int64(v.(int)))
		}
	} else {
		dnsRecord.Weight = currentRecord.Weight
	}

	if d.HasChange("priority") {
		if v, ok := d.GetOk("priority"); ok {
			dnsRecord.Priority = helper.Int64(int64(v.(int)))
		}
	} else {
		dnsRecord.Priority = currentRecord.Priority
	}

	request := teo.NewModifyDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.DnsRecords = []*teo.DnsRecord{&dnsRecord}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify teo dns records failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s modify teo dns records failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoDnsRecordV2Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId, recordId, err := resourceTencentCloudTeoDnsRecordV2ParseId(d.Id())
	if err != nil {
		return err
	}

	request := teo.NewDeleteDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.RecordIds = []*string{&recordId}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DeleteDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo dns records failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s delete teo dns records failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

// resourceTencentCloudTeoDnsRecordV2ParseId parses the composite ID (zone_id#record_id)
func resourceTencentCloudTeoDnsRecordV2ParseId(id string) (zoneId, recordId string, err error) {
	parts := strings.Split(id, "#")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource ID format: %s, expected format: zone_id#record_id", id)
	}
	return parts[0], parts[1], nil
}

// readDnsRecord reads a single DNS record from the cloud API
func readDnsRecord(meta interface{}, zoneId, recordId string) (*teo.DnsRecord, error) {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := teo.NewDescribeDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.Limit = helper.Int64(20)

	// Filter by record ID
	request.Filters = []*teo.AdvancedFilter{
		{
			Name:   helper.String("id"),
			Values: []*string{&recordId},
		},
	}

	var response *teo.DescribeDnsRecordsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo dns records failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(response.Response.DnsRecords) == 0 {
		return nil, nil
	}

	return response.Response.DnsRecords[0], nil
}

// Validation functions

func validateDnsRecordType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validTypes := map[string]bool{
		"A":     true,
		"AAAA":  true,
		"MX":    true,
		"CNAME": true,
		"TXT":   true,
		"NS":    true,
		"CAA":   true,
		"SRV":   true,
	}
	if !validTypes[value] {
		errors = append(errors, fmt.Errorf("%s must be one of: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV", k))
	}
	return
}

func validateDnsRecordTTL(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 60 || value > 86400 {
		errors = append(errors, fmt.Errorf("%s must be between 60 and 86400", k))
	}
	return
}

func validateDnsRecordWeight(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < -1 || value > 100 {
		errors = append(errors, fmt.Errorf("%s must be between -1 and 100", k))
	}
	return
}

func validateDnsRecordPriority(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 50 {
		errors = append(errors, fmt.Errorf("%s must be between 0 and 50", k))
	}
	return
}
