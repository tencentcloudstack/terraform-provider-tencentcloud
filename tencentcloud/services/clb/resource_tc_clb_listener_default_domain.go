package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbListenerDefaultDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbListenerDefaultDomainCreate,
		Read:   resourceTencentCloudClbListenerDefaultDomainRead,
		Update: resourceTencentCloudClbListenerDefaultDomainUpdate,
		Delete: resourceTencentCloudClbListenerDefaultDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of CLB instance.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of CLB listener.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name of the listener rule. Single domain rules are passed to `domain`, and multi domain rules are passed to `domains`.",
			},

			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of this CLB listener rule.",
			},
		},
	}
}

func resourceTencentCloudClbListenerDefaultDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_default_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = clb.NewModifyDomainAttributesRequest()
		response   *clb.ModifyDomainAttributesResponse
		clbId      string
		listenerId string
	)

	if v, ok := d.GetOk("clb_id"); ok {
		clbId = v.(string)
		request.LoadBalancerId = helper.String(clbId)
	}

	if v, ok := d.GetOk("listener_id"); ok {
		listenerId = v.(string)
		request.ListenerId = helper.String(listenerId)
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	request.DefaultServer = helper.Bool(true)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := client.ModifyDomainAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("modify domain failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls topic failed, reason:%+v", logId, err)
		return err
	}

	taskId := *response.Response.RequestId

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	_ = waitTaskReady(ctx, client, taskId)

	d.SetId(clbId + tccommon.FILED_SP + listenerId)
	return resourceTencentCloudClbListenerDefaultDomainRead(d, meta)
}

func resourceTencentCloudClbListenerDefaultDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_default_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Id()

	items := strings.Split(resourceId, tccommon.FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("id is broken,%s", resourceId)
	}
	clbId := items[0]
	listenerId := items[1]

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	filter := map[string]string{"listener_id": listenerId, "clb_id": clbId}
	var instances []*clb.RuleOutput
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeRulesByFilter(ctx, filter)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instances = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listener rule failed, reason:%+v", logId, err)
		return err
	}

	if len(instances) == 0 {
		d.SetId("")
		return nil
	}

	var (
		domain string
		ruleId string
	)

	for _, rule := range instances {
		if *rule.DefaultServer {
			domain = *rule.Domain
			ruleId = *rule.LocationId
			break
		}
	}

	_ = d.Set("clb_id", clbId)
	_ = d.Set("listener_id", listenerId)
	_ = d.Set("domain", domain)
	_ = d.Set("rule_id", ruleId)

	return nil
}

func resourceTencentCloudClbListenerDefaultDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_default_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	if d.HasChange("domain") {

		var (
			logId      = tccommon.GetLogId(tccommon.ContextNil)
			request    = clb.NewModifyDomainAttributesRequest()
			response   *clb.ModifyDomainAttributesResponse
			clbId      string
			listenerId string
		)

		if v, ok := d.GetOk("clb_id"); ok {
			clbId = v.(string)
			request.LoadBalancerId = helper.String(clbId)
		}

		if v, ok := d.GetOk("listener_id"); ok {
			listenerId = v.(string)
			request.ListenerId = helper.String(listenerId)
		}

		if v, ok := d.GetOk("domain"); ok {
			request.Domain = helper.String(v.(string))
		}

		request.DefaultServer = helper.Bool(true)

		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient()

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := client.ModifyDomainAttributes(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil {
				e = fmt.Errorf("modify domain failed")
				return resource.NonRetryableError(e)
			}

			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create cls topic failed, reason:%+v", logId, err)
			return err
		}

		taskId := *response.Response.RequestId

		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		_ = waitTaskReady(ctx, client, taskId)

	}

	return resourceTencentCloudClbListenerDefaultDomainRead(d, meta)
}

func resourceTencentCloudClbListenerDefaultDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_domain_default.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
