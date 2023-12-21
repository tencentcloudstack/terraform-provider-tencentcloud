package dnspod

import (
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDnspodModifyRecordGroupOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodModifyRecordGroupOperationCreate,
		Read:   resourceTencentCloudDnspodModifyRecordGroupOperationRead,
		Delete: resourceTencentCloudDnspodModifyRecordGroupOperationDelete,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Record Group ID.",
			},

			"record_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Record ID, multiple IDs are separated by a vertical line |.",
			},

			"domain_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.",
			},
		},
	}
}

func resourceTencentCloudDnspodModifyRecordGroupOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_modify_record_group_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = dnspod.NewModifyRecordToGroupRequest()
		domain   string
		recordId string
	)

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_id"); ok {
		recordId = v.(string)
		request.RecordId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		request.DomainId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("group_id"); ok {
		request.GroupId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyRecordToGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dnspod modify_record_group failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{domain, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudDnspodModifyRecordGroupOperationRead(d, meta)
}

func resourceTencentCloudDnspodModifyRecordGroupOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_modify_record_group_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDnspodModifyRecordGroupOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_modify_record_group_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
