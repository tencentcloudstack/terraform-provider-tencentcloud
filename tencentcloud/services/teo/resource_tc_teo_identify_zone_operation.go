package teo

import (
	"fmt"
	"log"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTeoIdentifyZoneOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoIdentifyZoneCreate,
		Read:   resourceTencentCloudTeoIdentifyZoneRead,
		Delete: resourceTencentCloudTeoIdentifyZoneDelete,
		Schema: map[string]*schema.Schema{
			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone name.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Subdomain under the zone. Required only when verifying a subdomain.",
			},
			"ascription": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "DNS verification information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subdomain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record host.",
						},
						"record_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record type.",
						},
						"record_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record value.",
						},
					},
				},
			},
			"file_ascription": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "File verification information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identify_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File verification path.",
						},
						"identify_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File verification content.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoIdentifyZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_identify_zone_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneName := d.Get("zone_name").(string)
	if zoneName == "" {
		return fmt.Errorf("zone_name is required")
	}

	var domain string
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	request := teov20220901.NewIdentifyZoneRequest()
	request.ZoneName = helper.String(zoneName)
	if domain != "" {
		request.Domain = helper.String(domain)
	}

	// Get Teo service
	service := NewTeoService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	// Call IdentifyZone API
	log.Printf("[DEBUG]%s api[%s] request body [%s]\n", logId, request.GetAction(), request.ToJsonString())

	ascription, fileAscription, err := service.TeoIdentifyZone(zoneName, domain)
	if err != nil {
		return fmt.Errorf("IdentifyZone failed: %s", err.Error())
	}

	log.Printf("[DEBUG]%s api[%s] success, ascription: %+v, file_ascription: %+v\n",
		logId, request.GetAction(), ascription, fileAscription)

	// Set ID
	d.SetId(zoneName)

	// Set ascription
	if ascription != nil {
		ascriptionMap := []map[string]interface{}{
			{
				"subdomain":    ascription.Subdomain,
				"record_type":  ascription.RecordType,
				"record_value": ascription.RecordValue,
			},
		}
		if err := d.Set("ascription", ascriptionMap); err != nil {
			return err
		}
	}

	// Set file_ascription
	if fileAscription != nil {
		fileAscriptionMap := []map[string]interface{}{
			{
				"identify_path":    fileAscription.IdentifyPath,
				"identify_content": fileAscription.IdentifyContent,
			},
		}
		if err := d.Set("file_ascription", fileAscriptionMap); err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoIdentifyZoneRead(d, meta)
}

func resourceTencentCloudTeoIdentifyZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_identify_zone_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoIdentifyZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_identify_zone_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
