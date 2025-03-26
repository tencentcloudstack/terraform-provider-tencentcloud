package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafAntiFake() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafAntiFakeCreate,
		Read:   resourceTencentCloudWafAntiFakeRead,
		Update: resourceTencentCloudWafAntiFakeUpdate,
		Delete: resourceTencentCloudWafAntiFakeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Rule Name.",
			},
			"uri": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Uri.",
			},
			"status": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      ANTI_FAKE_URL_STATUS_1,
				ValidateFunc: tccommon.ValidateAllowedIntValue(ANTI_FAKE_URL_STATUS),
				Description:  "Status. 0: Turn off rules and log switches, 1: Turn on the rule switch and Turn off the log switch; 2: Turn off the rule switch and turn on the log switch;3: Turn on the log switch.",
			},
			"rule_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Rule ID.",
			},
			"protocol": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Protocol.",
			},
		},
	}
}

func resourceTencentCloudWafAntiFakeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_anti_fake.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = waf.NewAddAntiFakeUrlRequest()
		response = waf.NewAddAntiFakeUrlResponse()
		id       string
		domain   string
	)

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uri"); ok {
		request.Uri = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().AddAntiFakeUrl(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf antiFake failed, reason:%+v", logId, err)
		return err
	}

	id = *response.Response.Id
	d.SetId(strings.Join([]string{id, domain}, tccommon.FILED_SP))

	// set status
	if v, ok := d.GetOkExists("status"); ok {
		status := v.(int)
		if status != API_SAFE_STATUS_1 {
			modifyAntiFakeUrlStatusRequest := waf.NewModifyAntiFakeUrlStatusRequest()
			idUInt, _ := strconv.ParseUint(id, 10, 64)
			modifyAntiFakeUrlStatusRequest.Ids = common.Uint64Ptrs([]uint64{idUInt})
			modifyAntiFakeUrlStatusRequest.Domain = &domain
			modifyAntiFakeUrlStatusRequest.Status = helper.IntUint64(v.(int))
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyAntiFakeUrlStatus(modifyAntiFakeUrlStatusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyAntiFakeUrlStatusRequest.GetAction(), modifyAntiFakeUrlStatusRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update waf antiFake status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudWafAntiFakeRead(d, meta)
}

func resourceTencentCloudWafAntiFakeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_anti_fake.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	id := idSplit[0]
	domain := idSplit[1]

	antiFake, err := service.DescribeWafAntiFakeById(ctx, id, domain)
	if err != nil {
		return err
	}

	if antiFake == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafAntiFake` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if antiFake.Domain != nil {
		_ = d.Set("domain", antiFake.Domain)
	}

	if antiFake.Name != nil {
		_ = d.Set("name", antiFake.Name)
	}

	if antiFake.Uri != nil {
		_ = d.Set("uri", antiFake.Uri)
	}

	if antiFake.Status != nil {
		_ = d.Set("status", antiFake.Status)
	}

	if antiFake.Id != nil {
		_ = d.Set("rule_id", antiFake.Id)
	}

	if antiFake.Protocol != nil {
		_ = d.Set("protocol", antiFake.Protocol)
	}

	return nil
}

func resourceTencentCloudWafAntiFakeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_anti_fake.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                          = tccommon.GetLogId(tccommon.ContextNil)
		modifyAntiFakeUrlRequest       = waf.NewModifyAntiFakeUrlRequest()
		modifyAntiFakeUrlStatusRequest = waf.NewModifyAntiFakeUrlStatusRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	id := idSplit[0]
	domain := idSplit[1]

	immutableArgs := []string{"domain"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idInt, _ := strconv.ParseInt(id, 10, 64)
	modifyAntiFakeUrlRequest.Id = &idInt
	modifyAntiFakeUrlRequest.Domain = &domain

	if v, ok := d.GetOk("name"); ok {
		modifyAntiFakeUrlRequest.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uri"); ok {
		modifyAntiFakeUrlRequest.Uri = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyAntiFakeUrl(modifyAntiFakeUrlRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyAntiFakeUrlRequest.GetAction(), modifyAntiFakeUrlRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update waf antiFake failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			modifyAntiFakeUrlStatusRequest.Status = helper.IntUint64(v.(int))
		}

		idUInt, _ := strconv.ParseUint(id, 10, 64)
		modifyAntiFakeUrlStatusRequest.Ids = common.Uint64Ptrs([]uint64{idUInt})
		modifyAntiFakeUrlStatusRequest.Domain = &domain
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyAntiFakeUrlStatus(modifyAntiFakeUrlStatusRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyAntiFakeUrlStatusRequest.GetAction(), modifyAntiFakeUrlStatusRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update waf antiFake status failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudWafAntiFakeRead(d, meta)
}

func resourceTencentCloudWafAntiFakeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_anti_fake.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	id := idSplit[0]
	domain := idSplit[1]

	if err := service.DeleteWafAntiFakeById(ctx, id, domain); err != nil {
		return err
	}

	return nil
}
