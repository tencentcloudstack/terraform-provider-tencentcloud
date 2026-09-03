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

func ResourceTencentCloudTeoDnsRecordsStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordsStatusCreate,
		Read:   resourceTencentCloudTeoDnsRecordsStatusRead,
		Update: resourceTencentCloudTeoDnsRecordsStatusUpdate,
		Delete: resourceTencentCloudTeoDnsRecordsStatusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"records_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNS record ID, combined with `zone_id` as the unique ID of the resource.",
			},

			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record status. Valid values: `enable` (enabled), `disable` (disabled).",
				ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
					value := v.(string)
					if value != "enable" && value != "disable" {
						errs = append(errs, fmt.Errorf("%q must be either `enable` or `disable`, got: %s", k, value))
					}
					return
				},
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecordsStatusCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_records_status.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		zoneId    string
		recordsId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("records_id"); ok {
		recordsId = v.(string)
	}

	d.SetId(zoneId + tccommon.FILED_SP + recordsId)
	return resourceTencentCloudTeoDnsRecordsStatusUpdate(d, meta)
}

func resourceTencentCloudTeoDnsRecordsStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_records_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	recordsId := idSplit[1]

	request := teov20220901.NewDescribeDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.Limit = helper.Int64(1000)
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("id"),
			Values: []*string{helper.String(recordsId)},
		},
	}

	var dnsRecord *teov20220901.DnsRecord
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

		dnsRecord = result.Response.DnsRecords[0]
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read teo_dns_records_status failed, reason:%+v", logId, err)
		return err
	}

	if dnsRecord != nil {
		if dnsRecord.Status != nil {
			_ = d.Set("status", *dnsRecord.Status)
		}
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordsStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_records_status.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	recordsId := idSplit[1]

	log.Printf("[CRUD]%s update teo_dns_records_status, zoneId=%s, recordsId=%s", logId, zoneId, recordsId)

	if d.HasChange("status") {
		status := d.Get("status").(string)

		request := teov20220901.NewModifyDnsRecordsStatusRequest()
		request.ZoneId = helper.String(zoneId)

		if status == "enable" {
			request.RecordsToEnable = []*string{helper.String(recordsId)}
		} else if status == "disable" {
			request.RecordsToDisable = []*string{helper.String(recordsId)}
		} else {
			return fmt.Errorf("invalid status %s, must be `enable` or `disable`", status)
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

		// ModifyDnsRecordsStatus is an async interface, poll DescribeDnsRecords until status reaches the target value.
		describeRequest := teov20220901.NewDescribeDnsRecordsRequest()
		describeRequest.ZoneId = helper.String(zoneId)
		describeRequest.Limit = helper.Int64(1000)
		describeRequest.Filters = []*teov20220901.AdvancedFilter{
			{
				Name:   helper.String("id"),
				Values: []*string{helper.String(recordsId)},
			},
		}

		if _, waitErr := (&resource.StateChangeConf{
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
			Pending:    []string{"pending", ""},
			Target:     []string{status},
			Timeout:    tccommon.ReadRetryTimeout,
			Refresh: func() (interface{}, string, error) {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeDnsRecordsWithContext(ctx, describeRequest)
				if e != nil {
					return nil, "", e
				}
				if result == nil || result.Response == nil || len(result.Response.DnsRecords) == 0 {
					return nil, "pending", nil
				}
				recordStatus := ""
				if result.Response.DnsRecords[0].Status != nil {
					recordStatus = *result.Response.DnsRecords[0].Status
				}
				return result, recordStatus, nil
			},
		}).WaitForStateContext(ctx); waitErr != nil {
			log.Printf("[CRITAL]%s wait teo_dns_records_status status to %s failed, reason:%+v", logId, status, waitErr)
			return fmt.Errorf("waiting for teo_dns_records_status (%s) to become %s: %s", d.Id(), status, waitErr)
		}
	}

	return resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)
}

func resourceTencentCloudTeoDnsRecordsStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_records_status.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
