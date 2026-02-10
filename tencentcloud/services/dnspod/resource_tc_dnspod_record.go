package dnspod

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDnspodRecord() *schema.Resource {
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
				Description: "Weight information. An integer from 1 to 100. Only enterprise VIP domain names are available, does not pass this parameter, means that the weight information is not set.",
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
				Description: "The monitoring status of the record.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Remark of record.",
			},
			"record_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the record.",
			},
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// weight设置后不能单独关闭
			if d.HasChange("weight") {
				old, new := d.GetChange("weight")
				if old.(int) != 0 && new.(int) == 0 {
					return fmt.Errorf("field `weight` cannot be unset once specified")
				}
			}
			return nil
		},
	}
}

func resourceTencentCloudDnspodRecordCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_record.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		recordId uint64
	)
	request := dnspod.NewCreateRecordRequest()
	requestRemark := dnspod.NewModifyRecordRemarkRequest()

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

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().CreateRecord(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		recordId = *response.Response.RecordId

		d.SetId(domain + tccommon.FILED_SP + strconv.FormatUint(recordId, 10))
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create DnsPod record failed, reason:%s\n", logId, err.Error())
		return err
	}

	if v, ok := d.GetOk("remark"); ok {
		requestRemark.Domain = helper.String(domain)
		requestRemark.RecordId = helper.Uint64(recordId)
		requestRemark.Remark = helper.String(v.(string))
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyRecordRemark(requestRemark)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create DnsPod record, modify remark failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudDnspodRecordRead(d, meta)
}

func resourceTencentCloudDnspodRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_record.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	id := d.Id()
	items := strings.Split(id, tccommon.FILED_SP)
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

	var recordInfo *dnspod.RecordInfo

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().DescribeRecord(request)
		if e != nil {
			e, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && e.GetCode() == "InvalidParameter.RecordIdInvalid" {
				// cannot find record id
				return nil
			} else {
				return tccommon.RetryError(e)
			}
		}
		recordInfo = response.Response.RecordInfo
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read DnsPod record failed, reason:%s\n", logId, err.Error())
		return err
	}

	if recordInfo == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("sub_domain", recordInfo.SubDomain)
	_ = d.Set("mx", recordInfo.MX)
	_ = d.Set("ttl", recordInfo.TTL)
	_ = d.Set("monitor_status", recordInfo.MonitorStatus)
	if recordInfo.Weight != nil {
		_ = d.Set("weight", recordInfo.Weight)
	}
	_ = d.Set("domain", items[0])
	_ = d.Set("record_line", recordInfo.RecordLine)
	_ = d.Set("record_type", recordInfo.RecordType)
	if v, ok := d.GetOk("value"); ok {
		value := v.(string)
		if strings.HasSuffix(value, ".") {
			_ = d.Set("value", recordInfo.Value)
		} else {
			_ = d.Set("value", strings.TrimSuffix(*recordInfo.Value, "."))
		}
	} else {
		_ = d.Set("value", recordInfo.Value)
	}
	_ = d.Set("remark", recordInfo.Remark)
	if *recordInfo.Enabled == uint64(0) {
		_ = d.Set("status", "DISABLE")
	} else {
		_ = d.Set("status", "ENABLE")
	}
	_ = d.Set("record_id", items[1])
	return nil
}

func resourceTencentCloudDnspodRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_record.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	id := d.Id()
	items := strings.Split(id, tccommon.FILED_SP)
	if len(items) < 2 {
		return nil
	}
	domain := items[0]
	recordId, err := strconv.Atoi(items[1])
	if err != nil {
		return err
	}
	request := dnspod.NewModifyRecordRequest()
	requestRemark := dnspod.NewModifyRecordRemarkRequest()
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

	if v, ok := d.GetOk("status"); ok {
		status := v.(string)
		request.Status = &status
	}
	if v, ok := d.GetOk("mx"); ok {
		request.MX = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("ttl"); ok {
		ttl := v.(int)
		request.TTL = helper.IntUint64(ttl)
	}
	if v, ok := d.GetOk("weight"); ok {
		weight := v.(int)
		request.Weight = helper.IntUint64(weight)
	}
	d.Partial(true)
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyRecord(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if d.HasChange("remark") {
		remark := d.Get("remark").(string)
		requestRemark.Domain = helper.String(domain)
		requestRemark.Remark = helper.String(remark)
		requestRemark.RecordId = helper.IntUint64(recordId)
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyRecordRemark(requestRemark)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s mdofify DnsPod record remark failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	d.Partial(false)
	return resourceTencentCloudDnspodRecordRead(d, meta)
}

func resourceTencentCloudDnspodRecordDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_record.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	id := d.Id()
	items := strings.Split(id, tccommon.FILED_SP)
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

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().DeleteRecord(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete DnsPod record failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
