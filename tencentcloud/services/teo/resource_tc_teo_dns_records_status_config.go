package teo

import (
	"context"
	"fmt"
	"log"

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

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	d.SetId(zoneId)
	return resourceTencentCloudTeoDnsRecordsStatusUpdate(d, meta)
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
			request.RecordsToEnable = helper.Strings(helper.InterfacesStrings(v.([]interface{})))
		}

		if v, ok := d.GetOk("records_to_disable"); ok {
			request.RecordsToDisable = helper.Strings(helper.InterfacesStrings(v.([]interface{})))
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
	}

	return resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)
}

func resourceTencentCloudTeoDnsRecordsStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_records_status.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
