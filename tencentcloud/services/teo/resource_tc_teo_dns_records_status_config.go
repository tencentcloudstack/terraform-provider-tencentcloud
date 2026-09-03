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

	// Because ModifyDnsRecordsStatus is an asynchronous API, poll DescribeDnsRecords
	// until the target record's status reaches the expected value or timeout.
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		readErr := resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)
		if readErr != nil {
			return resource.NonRetryableError(readErr)
		}

		if checkDnsRecordsStatus(d, meta) {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("teo_dns_records_status: record status has not reached the expected value, waiting for async operation to complete."))
	})

	if err != nil {
		log.Printf("[CRITAL]%s poll teo_dns_records_status status failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordsStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_records_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		zoneId  = d.Id()
		request = teov20220901.NewDescribeDnsRecordsRequest()
	)

	request.ZoneId = helper.String(zoneId)

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

		if len(result.Response.DnsRecords) == 0 {
			log.Printf("[CRUD] teo_dns_records_status id=%s", d.Id())
			d.SetId("")
			return nil
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read teo_dns_records_status failed, reason:%+v", logId, err)
		return err
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

			if checkDnsRecordsStatus(d, meta) {
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

// checkDnsRecordsStatus queries DescribeDnsRecords and checks whether the target record's
// status has reached the expected value (enable/disable).
func checkDnsRecordsStatus(d *schema.ResourceData, meta interface{}) bool {
	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		zoneId  = d.Id()
		request = teov20220901.NewDescribeDnsRecordsRequest()
	)

	request.ZoneId = helper.String(zoneId)

	recordsToEnable := d.Get("records_to_enable").([]interface{})
	recordsToDisable := d.Get("records_to_disable").([]interface{})

	if len(recordsToEnable) == 0 && len(recordsToDisable) == 0 {
		return true
	}

	enableMap := make(map[string]bool)
	for _, v := range recordsToEnable {
		enableMap[v.(string)] = true
	}
	disableMap := make(map[string]bool)
	for _, v := range recordsToDisable {
		disableMap[v.(string)] = true
	}

	result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeDnsRecordsWithContext(ctx, request)
	if e != nil {
		log.Printf("[CRITAL]%s poll teo_dns_records_status status failed, reason:%+v", logId, e)
		return false
	}

	if result == nil || result.Response == nil || len(result.Response.DnsRecords) == 0 {
		return false
	}

	for _, record := range result.Response.DnsRecords {
		if record == nil || record.RecordId == nil {
			continue
		}
		recordId := *record.RecordId
		if record.Status == nil {
			return false
		}
		status := *record.Status
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
