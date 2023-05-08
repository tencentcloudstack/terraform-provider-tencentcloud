/*
Provide a resource to create a DnsPod record.

Example Usage

```hcl
resource "tencentcloud_dnspod_record" "demo" {
  domain="mikatong.com"
  record_type="A"
  record_line="默认"
  value="1.2.3.9"
  sub_domain="demo"
}
```

Import

DnsPod Domain record can be imported using the Domain#RecordId, e.g.

```
$ terraform import tencentcloud_dnspod_record.demo arunma.com#1194109872
```
*/
package tencentcloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDnspodRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodRecordCreate,
		Read:   resourceTencentCloudDnspodRecordRead,
		Update: resourceTencentCloudDnspodRecordUpdate,
		Delete: resourceTencentCloudDnspodRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Domain.",
			},
			"record_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The record type.",
			},
			"record_line": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The record line.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The record value.",
			},
			"sub_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "@",
				Description: "The host records, default value is `@`.",
			},
			"mx": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "MX priority, valid when the record type is MX, range 1-20. Note: must set when record type equal MX.",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     600,
				Description: "TTL, the range is 1-604800, and the minimum value of different levels of domain names is different. Default is 600.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Weight information. An integer from 0 to 100. Only enterprise VIP domain names are available, 0 means off, does not pass this parameter, means that the weight information is not set. Default is 0.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ENABLE",
				Description: "Records the initial state, with values ranging from ENABLE and DISABLE. The default is ENABLE, and if DISABLE is passed in, resolution will not take effect and the limits of load balancing will not be verified.",
			},
			"monitor_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The D monitoring status of the record.",
			},
		},
	}
}

func resourceTencentCloudDnspodRecordCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_record.create")()
	logId := getLogId(contextNil)
	request := dnspod.NewCreateRecordRequest()

	domain := d.Get("domain").(string)
	recordType := d.Get("record_type").(string)
	recordLine := d.Get("record_line").(string)
	value := d.Get("value").(string)
	subDomain := d.Get("sub_domain").(string)
	ttl := d.Get("ttl").(int)
	status := d.Get("status").(string)
	request.Domain = &domain
	request.RecordType = &recordType
	request.RecordLine = &recordLine
	request.Value = &value
	request.SubDomain = &subDomain
	if v, ok := d.GetOk("mx"); ok {
		request.MX = helper.IntUint64(v.(int))
	}
	request.TTL = helper.IntUint64(ttl)
	if v, ok := d.GetOk("weight"); ok {
		request.Weight = helper.IntUint64(v.(int))
	}
	request.Status = &status

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().CreateRecord(request)
		if e != nil {
			return retryError(e)
		}
		recordId := *response.Response.RecordId

		d.SetId(domain + FILED_SP + fmt.Sprint(recordId))

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create DnsPod record failed, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudDnspodRecordRead(d, meta)
}

func resourceTencentCloudDnspodRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_record.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	id := d.Id()
	items := strings.Split(id, FILED_SP)
	if len(items) < 2 {
		return nil
	}
	request := dnspod.NewDescribeRecordRequest()
	request.Domain = helper.String(items[0])
	recordId, err := strconv.Atoi(items[1])
	if err != nil {
		return err
	}
	request.RecordId = helper.IntUint64(recordId)

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().DescribeRecord(request)
		if e != nil {
			return retryError(e)
		}

		recordInfo := response.Response.RecordInfo

		_ = d.Set("sub_domain", recordInfo.SubDomain)
		_ = d.Set("mx", recordInfo.MX)
		_ = d.Set("ttl", recordInfo.TTL)
		_ = d.Set("monitor_status", recordInfo.MonitorStatus)
		_ = d.Set("weight", recordInfo.Weight)
		_ = d.Set("domain", items[0])
		_ = d.Set("record_line", recordInfo.RecordLine)
		_ = d.Set("record_type", recordInfo.RecordType)
		_ = d.Set("value", recordInfo.Value)
		if *recordInfo.Enabled == uint64(0) {
			_ = d.Set("status", "DISABLE")
		} else {
			_ = d.Set("status", "ENABLE")
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read DnsPod record failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}

func resourceTencentCloudDnspodRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_record.update")()

	id := d.Id()
	items := strings.Split(id, FILED_SP)
	if len(items) < 2 {
		return nil
	}
	domain := items[0]
	recordId, err := strconv.Atoi(items[1])
	if err != nil {
		return err
	}
	request := dnspod.NewModifyRecordRequest()
	request.Domain = &domain
	request.RecordId = helper.IntUint64(recordId)
	recordType := d.Get("record_type").(string)
	recordLine := d.Get("record_line").(string)
	value := d.Get("value").(string)
	subDomain := d.Get("sub_domain").(string)
	request.RecordType = &recordType
	request.RecordLine = &recordLine
	request.Value = &value
	request.SubDomain = &subDomain

	if d.HasChange("status") {
		status := d.Get("status").(string)
		request.Status = &status
	}
	if d.HasChange("mx") {
		if v, ok := d.GetOk("mx"); ok {
			request.MX = helper.IntUint64(v.(int))
		}
	}
	if d.HasChange("ttl") {
		ttl := d.Get("ttl").(int)
		request.TTL = helper.IntUint64(ttl)
	}
	if d.HasChange("weight") {
		weight := d.Get("weight").(int)
		request.TTL = helper.IntUint64(weight)
	}
	d.Partial(true)
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().ModifyRecord(request)
		if e != nil {
			return retryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceTencentCloudDnspodRecordRead(d, meta)
}

func resourceTencentCloudDnspodRecordDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_record.delete")()

	logId := getLogId(contextNil)
	id := d.Id()
	items := strings.Split(id, FILED_SP)
	if len(items) < 2 {
		return nil
	}
	domain := items[0]
	recordId, err := strconv.Atoi(items[1])
	if err != nil {
		return err
	}
	request := dnspod.NewDeleteRecordRequest()
	request.Domain = helper.String(domain)
	request.RecordId = helper.IntUint64(recordId)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().DeleteRecord(request)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete DnsPod record failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
