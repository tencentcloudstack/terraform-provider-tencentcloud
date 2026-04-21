package teo

import (
	"fmt"
	"log"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTeoCheckCnameStatusOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoCheckCnameStatusOperationCreate,
		Read:   resourceTencentCloudTeoCheckCnameStatusOperationRead,
		Delete: resourceTencentCloudTeoCheckCnameStatusOperationDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"record_names": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of record names to check CNAME status.",
			},
			"cname_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "CNAME status information for each record name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record name.",
						},
						"cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CNAME address. May be null.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CNAME status. Valid values: active, moved.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoCheckCnameStatusOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_check_cname_status_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	recordNamesRaw := d.Get("record_names").([]interface{})

	if zoneId == "" {
		return fmt.Errorf("zone_id is required")
	}

	if len(recordNamesRaw) == 0 {
		return fmt.Errorf("record_names is required")
	}

	recordNames := make([]*string, 0, len(recordNamesRaw))
	for _, v := range recordNamesRaw {
		recordNames = append(recordNames, helper.String(v.(string)))
	}

	request := teov20220901.NewCheckCnameStatusRequest()
	request.ZoneId = helper.String(zoneId)
	request.RecordNames = recordNames

	var response *teov20220901.CheckCnameStatusResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CheckCnameStatus(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	if response == nil || response.Response == nil {
		return fmt.Errorf("CheckCnameStatus API returned empty response")
	}

	d.SetId(zoneId)

	cnameStatusList := make([]map[string]interface{}, 0)
	for _, status := range response.Response.CnameStatus {
		cnameStatus := map[string]interface{}{}
		if status.RecordName != nil {
			cnameStatus["record_name"] = *status.RecordName
		}
		if status.Cname != nil {
			cnameStatus["cname"] = *status.Cname
		}
		if status.Status != nil {
			cnameStatus["status"] = *status.Status
		}
		cnameStatusList = append(cnameStatusList, cnameStatus)
	}

	if err := d.Set("cname_status", cnameStatusList); err != nil {
		return fmt.Errorf("Set cname_status failed: %s", err)
	}

	return resourceTencentCloudTeoCheckCnameStatusOperationRead(d, meta)
}

func resourceTencentCloudTeoCheckCnameStatusOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_check_cname_status_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoCheckCnameStatusOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_check_cname_status_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
