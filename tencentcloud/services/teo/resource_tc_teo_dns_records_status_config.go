package teo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoDnsRecordsStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordsStatusCreate,
		Read:   resourceTencentCloudTeoDnsRecordsStatusRead,
		Update: resourceTencentCloudTeoDnsRecordsStatusUpdate,
		Delete: resourceTencentCloudTeoDnsRecordsStatusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions, each filter element contains `name`, `values` and `fuzzy`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Filter values of the field.",
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable fuzzy query.",
						},
					},
				},
			},

			"sort_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort by field.",
			},

			"sort_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort order, valid values: `asc`, `desc`.",
			},

			"match": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Match mode, valid values: `all`, `any`.",
			},

			"records_to_enable": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "DNS record ID list to be enabled, only manages a single resource, pass in a single record ID.",
			},

			"records_to_disable": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "DNS record ID list to be disabled, only manages a single resource, pass in a single record ID.",
			},

			"dns_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DNS record list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Site ID.",
						},
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record type.",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record resolution line.",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record content.",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cache time, unit: seconds.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DNS record weight.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MX record priority.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record resolution status, valid values: `enable`, `disable`.",
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
				},
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecordsStatusCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_records_status.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		zoneId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	log.Printf("[CRUD]%s create teo_dns_records_status, zoneId=%s", logId, zoneId)

	recordsToEnable := d.Get("records_to_enable").([]interface{})
	recordsToDisable := d.Get("records_to_disable").([]interface{})

	if len(recordsToEnable) == 0 && len(recordsToDisable) == 0 {
		d.SetId(zoneId)
		return resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)
	}

	request := teov20220901.NewModifyDnsRecordsStatusRequest()
	request.ZoneId = helper.String(zoneId)

	if len(recordsToEnable) > 0 {
		request.RecordsToEnable = helper.Strings(interfaceToStringSlice(recordsToEnable))
	}

	if len(recordsToDisable) > 0 {
		request.RecordsToDisable = helper.Strings(interfaceToStringSlice(recordsToDisable))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyDnsRecordsStatusWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo_dns_records_status failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo_dns_records_status failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(zoneId)

	return resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)
}

func resourceTencentCloudTeoDnsRecordsStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_records_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		zoneId     = d.Id()
		request    = teov20220901.NewDescribeDnsRecordsRequest()
		dnsRecords []*teov20220901.DnsRecord
	)

	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("filters"); ok {
		filters := make([]*teov20220901.AdvancedFilter, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			filterMap := item.(map[string]interface{})
			filter := teov20220901.AdvancedFilter{}
			if v, ok := filterMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}
			if v, ok := filterMap["values"].([]interface{}); ok && len(v) > 0 {
				filter.Values = helper.Strings(interfaceToStringSlice(v))
			}
			if v, ok := filterMap["fuzzy"].(bool); ok {
				filter.Fuzzy = helper.Bool(v)
			}
			filters = append(filters, &filter)
		}
		request.Filters = filters
	}

	if v, ok := d.GetOk("sort_by"); ok {
		request.SortBy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_order"); ok {
		request.SortOrder = helper.String(v.(string))
	}

	if v, ok := d.GetOk("match"); ok {
		request.Match = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Read teo_dns_records_status failed, Response is nil."))
		}

		dnsRecords = result.Response.DnsRecords
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read teo_dns_records_status failed, reason:%+v", logId, err)
		return err
	}

	if len(dnsRecords) == 0 {
		log.Printf("[CRUD] teo_dns_records_status id=%s", d.Id())
		d.SetId("")
		return nil
	}

	dnsRecordsList := make([]map[string]interface{}, 0, len(dnsRecords))
	for _, record := range dnsRecords {
		recordMap := map[string]interface{}{}
		if record.ZoneId != nil {
			recordMap["zone_id"] = record.ZoneId
		}
		if record.RecordId != nil {
			recordMap["record_id"] = record.RecordId
		}
		if record.Name != nil {
			recordMap["name"] = record.Name
		}
		if record.Type != nil {
			recordMap["type"] = record.Type
		}
		if record.Location != nil {
			recordMap["location"] = record.Location
		}
		if record.Content != nil {
			recordMap["content"] = record.Content
		}
		if record.TTL != nil {
			recordMap["ttl"] = record.TTL
		}
		if record.Weight != nil {
			recordMap["weight"] = record.Weight
		}
		if record.Priority != nil {
			recordMap["priority"] = record.Priority
		}
		if record.Status != nil {
			recordMap["status"] = record.Status
		}
		if record.CreatedOn != nil {
			recordMap["created_on"] = record.CreatedOn
		}
		if record.ModifiedOn != nil {
			recordMap["modified_on"] = record.ModifiedOn
		}
		dnsRecordsList = append(dnsRecordsList, recordMap)
	}

	if len(dnsRecordsList) > 0 {
		_ = d.Set("dns_records", dnsRecordsList)
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordsStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_records_status.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		zoneId = d.Id()
	)

	log.Printf("[CRUD]%s update teo_dns_records_status, zoneId=%s", logId, zoneId)

	needChange := false
	mutableArgs := []string{"records_to_enable", "records_to_disable"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifyDnsRecordsStatusRequest()
		request.ZoneId = helper.String(zoneId)

		if v, ok := d.GetOk("records_to_enable"); ok {
			request.RecordsToEnable = helper.Strings(interfaceToStringSlice(v.([]interface{})))
		}

		if v, ok := d.GetOk("records_to_disable"); ok {
			request.RecordsToDisable = helper.Strings(interfaceToStringSlice(v.([]interface{})))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyDnsRecordsStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Update teo_dns_records_status failed, Response is nil."))
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo_dns_records_status failed, reason:%+v", logId, err)
			return err
		}

		// Because ModifyDnsRecordsStatus is an asynchronous API, poll DescribeDnsRecords
		// until the target record's status reaches the expected value or timeout.
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			readErr := resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)
			if readErr != nil {
				return resource.NonRetryableError(readErr)
			}

			dnsRecords := d.Get("dns_records").([]interface{})
			if len(dnsRecords) == 0 {
				return resource.RetryableError(fmt.Errorf("teo_dns_records_status: dns_records is empty, waiting for status to take effect."))
			}

			if checkDnsRecordsStatus(dnsRecords, d) {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("teo_dns_records_status: record status has not reached the expected value, waiting for async operation to complete."))
		})

		if err != nil {
			log.Printf("[CRITAL]%s poll teo_dns_records_status status failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)
}

func resourceTencentCloudTeoDnsRecordsStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_records_status.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

// checkDnsRecordsStatus checks whether the target record's status has reached the expected value.
func checkDnsRecordsStatus(dnsRecords []interface{}, d *schema.ResourceData) bool {
	recordsToEnable := d.Get("records_to_enable").([]interface{})
	recordsToDisable := d.Get("records_to_disable").([]interface{})

	enableMap := make(map[string]bool)
	for _, v := range recordsToEnable {
		enableMap[v.(string)] = true
	}
	disableMap := make(map[string]bool)
	for _, v := range recordsToDisable {
		disableMap[v.(string)] = true
	}

	for _, record := range dnsRecords {
		recordMap := record.(map[string]interface{})
		recordId, ok := recordMap["record_id"].(string)
		if !ok || recordId == "" {
			continue
		}
		status, ok := recordMap["status"].(string)
		if !ok {
			continue
		}
		if enableMap[recordId] && status != "enable" {
			return false
		}
		if disableMap[recordId] && status != "disable" {
			return false
		}
	}

	return true
}

// interfaceToStringSlice converts []interface{} to []string.
func interfaceToStringSlice(list []interface{}) []string {
	result := make([]string, 0, len(list))
	for _, v := range list {
		result = append(result, v.(string))
	}
	return result
}
