package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWafCcAutoStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafCcAutoStatusCreate,
		Read:   resourceTencentCloudWafCcAutoStatusRead,
		Delete: resourceTencentCloudWafCcAutoStatusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},
			"edition": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(EDITION_TYPE),
				Description:  "Waf edition. clb-waf means clb-waf, sparta-waf means saas-waf.",
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "cc auto status, 1 means open, 0 means close.",
			},
		},
	}
}

func resourceTencentCloudWafCcAutoStatusCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_cc_auto_status.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = waf.NewUpsertCCAutoStatusRequest()
		domain  string
		edition string
	)

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("edition"); ok {
		request.Edition = helper.String(v.(string))
		edition = v.(string)
	}

	request.Value = helper.IntInt64(1)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().UpsertCCAutoStatus(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf CcAutoStatus failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{domain, edition}, FILED_SP))

	return resourceTencentCloudWafCcAutoStatusRead(d, meta)
}

func resourceTencentCloudWafCcAutoStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_cc_auto_status.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	domain := idSplit[0]
	edition := idSplit[1]

	CcAutoStatus, err := service.DescribeWafCcAutoStatusById(ctx, domain)
	if err != nil {
		return err
	}

	if CcAutoStatus == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafCcAutoStatus` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("edition", edition)

	if CcAutoStatus.AutoCCSwitch != nil {
		_ = d.Set("status", CcAutoStatus.AutoCCSwitch)
	}

	return nil
}

func resourceTencentCloudWafCcAutoStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_cc_auto_status.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	domain := idSplit[0]
	edition := idSplit[1]

	if err := service.DeleteWafCcAutoStatusById(ctx, domain, edition); err != nil {
		return err
	}

	return nil
}
