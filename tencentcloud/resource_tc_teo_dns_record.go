/*
Provides a resource to create a teo dns_record

~> **NOTE:** This resource has been deprecated in Terraform TencentCloud Provider Version 1.79.19.

Example Usage

```hcl
resource "tencentcloud_teo_dns_record" "dns_record" {
  zone_id   = "zone-297z8rf93cfw"
  type      = "A"
  name      = "www.toutiao2.com"
  content   = "150.109.8.2"
  mode      = "proxied"
  ttl       = "1"
  priority  = 1
}

```
Import

teo dns_record can be imported using the zone_id#dns_record_id, e.g.
```
$ terraform import tencentcloud_teo_dns_record.dns_record zone-297z8rf93cfw#record-297z9ei9b9oc
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoDnsRecord() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoDnsRecordRead,
		Create: resourceTencentCloudTeoDnsRecordCreate,
		Update: resourceTencentCloudTeoDnsRecordUpdate,
		Delete: resourceTencentCloudTeoDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"dns_record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record ID.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record Type. Valid values: `A`, `AAAA`, `CNAME`, `MX`, `TXT`, `NS`, `CAA`, `SRV`.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record Name.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record Content.",
			},

			"mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Proxy mode. Valid values:- `dns_only`: only DNS resolution of the subdomain is enabled.- `proxied`: subdomain is proxied and accelerated.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Time to live of the DNS record cache in seconds.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Priority of the record. Valid value range: 1-50, the smaller value, the higher priority.",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date.",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification date.",
			},

			"locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the DNS record is locked.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Resolution status. Valid values: `active`, `pending`.",
			},

			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME address. Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"domain_status": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Whether this domain enable load balancing, security, or l4 proxy capability. Valid values: `lb`, `security`, `l4`.",
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_record.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = teo.NewCreateDnsRecordRequest()
		response    *teo.CreateDnsRecordResponse
		zoneId      string
		dnsRecordId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mode"); ok {
		request.Mode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ttl"); ok {
		request.TTL = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := service.CheckZoneComplete(ctx, zoneId)
	if err != nil {
		log.Printf("[CRITAL]%s create teo dnsRecord failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateDnsRecord(request)
		if e != nil {
			return retryError(e, "OperationDenied", "UnauthorizedOperation")
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo dnsRecord failed, reason:%+v", logId, err)
		return err
	}

	dnsRecordId = *response.Response.DnsRecordId

	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoDnsRecord(ctx, zoneId, dnsRecordId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.Status == "pending" {
			return resource.RetryableError(fmt.Errorf("dnsRecord status is %v, retry...", *instance.Status))
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(zoneId + FILED_SP + dnsRecordId)
	return resourceTencentCloudTeoDnsRecordRead(d, meta)
}

func resourceTencentCloudTeoDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_record.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	dnsRecordId := idSplit[1]

	dnsRecord, err := service.DescribeTeoDnsRecord(ctx, zoneId, dnsRecordId)
	if err != nil {
		return err
	}

	if dnsRecord == nil {
		d.SetId("")
		return fmt.Errorf("resource `dnsRecord` %s does not exist", dnsRecordId)
	}

	if dnsRecord.ZoneId != nil {
		_ = d.Set("zone_id", dnsRecord.ZoneId)
	}

	if dnsRecord.DnsRecordId != nil {
		_ = d.Set("dns_record_id", dnsRecord.DnsRecordId)
	}

	if dnsRecord.DnsRecordType != nil {
		_ = d.Set("type", dnsRecord.DnsRecordType)
	}

	if dnsRecord.DnsRecordName != nil {
		_ = d.Set("name", dnsRecord.DnsRecordName)
	}

	if dnsRecord.Content != nil {
		_ = d.Set("content", dnsRecord.Content)
	}

	if dnsRecord.Mode != nil {
		_ = d.Set("mode", dnsRecord.Mode)
	}

	if dnsRecord.TTL != nil {
		_ = d.Set("ttl", dnsRecord.TTL)
	}

	if dnsRecord.Priority != nil {
		_ = d.Set("priority", dnsRecord.Priority)
	}

	if dnsRecord.CreatedOn != nil {
		_ = d.Set("created_on", dnsRecord.CreatedOn)
	}

	if dnsRecord.ModifiedOn != nil {
		_ = d.Set("modified_on", dnsRecord.ModifiedOn)
	}

	if dnsRecord.Locked != nil {
		_ = d.Set("locked", dnsRecord.Locked)
	}

	if dnsRecord.Status != nil {
		_ = d.Set("status", dnsRecord.Status)
	}

	if dnsRecord.Cname != nil {
		_ = d.Set("cname", dnsRecord.Cname)
	}

	if dnsRecord.DomainStatus != nil {
		_ = d.Set("domain_status", dnsRecord.DomainStatus)
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_record.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	request := teo.NewModifyDnsRecordRequest()
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	dnsRecordId := idSplit[1]

	request.ZoneId = &zoneId
	request.DnsRecordId = &dnsRecordId

	if d.HasChange("zone_id") {
		return fmt.Errorf("`zone_id` do not support change now.")
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.DnsRecordType = helper.String(v.(string))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.DnsRecordName = helper.String(v.(string))
		}
	}

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))
		}
	}

	if d.HasChange("mode") {
		if v, ok := d.GetOk("mode"); ok {
			request.Mode = helper.String(v.(string))
		}
	}

	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			request.TTL = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("priority") {
		if v, ok := d.GetOk("priority"); ok {
			request.Priority = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDnsRecord(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[UPDATE]%s update teo dnsRecord failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoDnsRecordRead(d, meta)
}

func resourceTencentCloudTeoDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_record.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	dnsRecordId := idSplit[1]

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if e := service.DeleteTeoDnsRecordById(ctx, zoneId, dnsRecordId); e != nil {
			return retryError(e, "OperationDenied", InternalError)
		}
		return nil
	})

	if err != nil {
		log.Printf("[DELETE]%s delete teo dnsRecord failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
