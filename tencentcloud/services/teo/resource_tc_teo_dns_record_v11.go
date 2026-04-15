package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudTeoDnsRecordV11() *schema.Resource {
	return resourceTencentCloudTeoDnsRecordV11()
}

func resourceTencentCloudTeoDnsRecordV11() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordV11Create,
		Read:   resourceTencentCloudTeoDnsRecordV11Read,
		Update: resourceTencentCloudTeoDnsRecordV11Update,
		Delete: resourceTencentCloudTeoDnsRecordV11Delete,
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
				Description: "ID of the site related with the DNS record.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNS record name. If it is Chinese, Korean, or Japanese domain name, it needs to be converted to punycode before input.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"A", "AAAA", "MX", "CNAME", "TXT", "NS", "CAA", "SRV"}),
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. If it is Chinese, Korean, or Japanese domain name, it needs to be converted to punycode before input.",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time. Value range: 60~86400, default is 300, unit: seconds.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. Value range: -1~100, default is -1 (not set). 0 means the record will not resolve. Only applicable for A, AAAA, and CNAME record types.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority. Value range: 0~50, default is 0. Only applicable for MX record type.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "DNS record resolution line. Default is Default, indicating the default resolution line, representing all regions. Only applicable for A, AAAA, and CNAME record types.",
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
			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time of the DNS record.",
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecordV11Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v11.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		param   = make(map[string]interface{})
		zoneId  string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		param["zone_id"] = v.(string)
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		param["name"] = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		param["type"] = v.(string)
	}

	if v, ok := d.GetOk("content"); ok {
		param["content"] = v.(string)
	}

	if v, ok := d.GetOkExists("ttl"); ok {
		param["ttl"] = v.(int)
	}

	if v, ok := d.GetOkExists("weight"); ok {
		param["weight"] = v.(int)
	}

	if v, ok := d.GetOkExists("priority"); ok {
		param["priority"] = v.(int)
	}

	if v, ok := d.GetOk("location"); ok {
		param["location"] = v.(string)
	}

	recordId, err := service.CreateDnsRecord(ctx, param)
	if err != nil {
		return err
	}

	d.SetId(buildDnsRecordV11Id(zoneId, recordId))

	return resourceTencentCloudTeoDnsRecordV11Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV11Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v11.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	zoneId, recordId, err := parseDnsRecordV11Id(d.Id())
	if err != nil {
		return err
	}

	dnsRecord, err := service.DescribeDnsRecordById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}

	if dnsRecord == nil {
		d.SetId("")
		log.Printf("[WARN]%s dns record not found, zoneId: %s, recordId: %s\n", logId, zoneId, recordId)
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("name", dnsRecord.Name)
	_ = d.Set("type", dnsRecord.Type)
	_ = d.Set("content", dnsRecord.Content)
	_ = d.Set("ttl", dnsRecord.TTL)
	_ = d.Set("weight", dnsRecord.Weight)
	_ = d.Set("priority", dnsRecord.Priority)
	_ = d.Set("location", dnsRecord.Location)
	_ = d.Set("status", dnsRecord.Status)
	_ = d.Set("created_on", dnsRecord.CreatedOn)
	_ = d.Set("modified_on", dnsRecord.ModifiedOn)

	return nil
}

func resourceTencentCloudTeoDnsRecordV11Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v11.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		param   = make(map[string]interface{})
	)

	zoneId, recordId, err := parseDnsRecordV11Id(d.Id())
	if err != nil {
		return err
	}

	param["record_id"] = recordId

	// Check if name or type was changed (which is not allowed)
	if d.HasChange("name") || d.HasChange("type") {
		return fmt.Errorf("name and type fields cannot be modified after creation")
	}

	// Only allow updating specific fields
	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			param["content"] = v.(string)
		}
	}

	if d.HasChange("ttl") {
		if v, ok := d.GetOkExists("ttl"); ok {
			param["ttl"] = v.(int)
		}
	}

	if d.HasChange("weight") {
		if v, ok := d.GetOkExists("weight"); ok {
			param["weight"] = v.(int)
		}
	}

	if d.HasChange("priority") {
		if v, ok := d.GetOkExists("priority"); ok {
			param["priority"] = v.(int)
		}
	}

	if d.HasChange("location") {
		if v, ok := d.GetOk("location"); ok {
			param["location"] = v.(string)
		}
	}

	err = service.ModifyDnsRecord(ctx, zoneId, param)
	if err != nil {
		return err
	}

	return resourceTencentCloudTeoDnsRecordV11Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV11Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v11.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	zoneId, recordId, err := parseDnsRecordV11Id(d.Id())
	if err != nil {
		return err
	}

	err = service.DeleteDnsRecordById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}

	return nil
}

func buildDnsRecordV11Id(zoneId, recordId string) string {
	return fmt.Sprintf("%s#%s", zoneId, recordId)
}

func parseDnsRecordV11Id(id string) (zoneId, recordId string, err error) {
	parts := strings.Split(id, "#")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid dns record id format, expected: zoneId#recordId, got: %s", id)
	}
	return parts[0], parts[1], nil
}
