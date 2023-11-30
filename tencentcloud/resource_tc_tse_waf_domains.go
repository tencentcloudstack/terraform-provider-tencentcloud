package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTseWafDomains() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseWafDomainsCreate,
		Read:   resourceTencentCloudTseWafDomainsRead,
		Delete: resourceTencentCloudTseWafDomainsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The waf protected domain name.",
			},
		},
	}
}

func resourceTencentCloudTseWafDomainsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_waf_domains.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tse.NewCreateWafDomainsRequest()
		gatewayId string
		domain    string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domains = append(request.Domains, helper.String(v.(string)))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CreateWafDomains(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse wafDomains failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{gatewayId, domain}, FILED_SP))

	return resourceTencentCloudTseWafDomainsRead(d, meta)
}

func resourceTencentCloudTseWafDomainsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_waf_domains.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	domain := idSplit[1]

	wafDomains, err := service.DescribeTseWafDomainsById(ctx, gatewayId)
	if err != nil {
		return err
	}

	if wafDomains == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseWafDomains` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)

	if wafDomains.Domains != nil {
		for _, v := range wafDomains.Domains {
			if *v == domain {
				_ = d.Set("domain", domain)
				break
			}
		}
	}

	return nil
}

func resourceTencentCloudTseWafDomainsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_waf_domains.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	domain := idSplit[1]

	if err := service.DeleteTseWafDomainsById(ctx, gatewayId, domain); err != nil {
		return err
	}

	return nil
}
