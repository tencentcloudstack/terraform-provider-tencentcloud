package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudClbServerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbServerAttachmentCreate,
		Read:   resourceTencentCloudClbServerAttachmentRead,
		Update: resourceTencentCloudClbServerAttachmentUpdate,
		Delete: resourceTencentCloudClbServerAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the CLB.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the CLB listener.",
			},
			"rule_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"domain", "url"},
				Description:   "ID of the CLB listener rule. Only supports listeners of `HTTPS` and `HTTP` protocol.",
			},
			"domain": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				RequiredWith:  []string{"url"},
				ConflictsWith: []string{"rule_id"},
				Description:   "Domain of the target forwarding rule. Does not take effect when parameter `rule_id` is provided.",
			},
			"url": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				RequiredWith:  []string{"domain"},
				ConflictsWith: []string{"rule_id"},
				Description:   "URL of the target forwarding rule. Does not take effect when parameter `rule_id` is provided.",
			},
			"protocol_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of protocol within the listener.",
			},
			"targets": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    100,
				Description: "Information of the backends to be attached.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "CVM Instance Id of the backend server, conflict with `eni_ip` but must specify one of them.",
						},
						"eni_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Eni IP address of the backend server, conflict with `instance_id` but must specify one of them.",
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 65535),
							Description:  "Port of the backend server. Valid value ranges: (0~65535).",
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      10,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 100),
							Description:  "Forwarding weight of the backend service. Valid value ranges: (0~100). defaults to `10`.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClbServerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_attachment.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = clb.NewRegisterTargetsRequest()
		locationId string
		domain     string
		url        string
	)

	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}

	clbId := d.Get("clb_id").(string)
	request.LoadBalancerId = helper.String(clbId)
	request.ListenerId = helper.String(listenerId)
	if v, ok := d.GetOk("rule_id"); ok {
		locationId = v.(string)
		checkErr = RuleIdCheck(locationId)
		if checkErr != nil {
			return checkErr
		}

		if locationId != "" {
			request.LocationId = helper.String(locationId)
		}
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("url"); ok {
		request.Url = helper.String(v.(string))
		url = v.(string)
	}

	insList := d.Get("targets").(*schema.Set).List()
	insLen := len(insList)
	for count := 0; count < insLen; count += 20 {
		//this request only support 20 targets at most once time
		request.Targets = make([]*clb.Target, 0, 20)
		for i := 0; i < 20; i++ {
			index := count + i
			if index >= insLen {
				break
			}

			inst := insList[index].(map[string]interface{})
			request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["eni_ip"], inst["port"], inst["weight"]))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			requestId := ""
			ratelimit.Check(request.GetAction())
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().RegisterTargets(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				requestId = *result.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create CLB attachment failed, reason:%+v", logId, err)
			return err
		}
	}

	var id string
	if locationId != "" {
		id = fmt.Sprintf("%s#%v#%v", locationId, d.Get("listener_id"), d.Get("clb_id"))
	} else if domain != "" && url != "" {
		id = fmt.Sprintf("%s,%s#%v#%v", domain, url, d.Get("listener_id"), d.Get("clb_id"))
	} else {
		// only api support for now
		id = fmt.Sprintf("%s#%v#%v", "", d.Get("listener_id"), d.Get("clb_id"))
	}

	d.SetId(id)

	return resourceTencentCloudClbServerAttachmentRead(d, meta)
}

func resourceTencentCloudClbServerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instance   *clb.ListenerBackend
		locationId string
		domain     string
		url        string
	)

	items := strings.Split(d.Id(), "#")
	locationIdOrDomainUrl := items[0]
	listenerId := items[1]
	clbId := items[2]

	if locationIdOrDomainUrl != "" {
		if strings.HasPrefix(locationIdOrDomainUrl, "loc-") && !strings.Contains(locationIdOrDomainUrl, ",") {
			// get locationId
			locationId = locationIdOrDomainUrl
		} else {
			// get domain & url
			domainUrl := strings.Split(locationIdOrDomainUrl, ",")
			domain = domainUrl[0]
			url = domainUrl[1]
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		instance = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read CLB attachment failed, reason:%+v", logId, err)
		return err
	}

	//see if read empty
	if instance == nil ||
		(len(instance.Targets) == 0 && locationId == "" && domain == "" && url == "") ||
		(len(instance.Rules) == 0 && locationId != "") {
		d.SetId("")
		return nil
	}

	_ = d.Set("clb_id", clbId)
	_ = d.Set("listener_id", listenerId)
	_ = d.Set("protocol_type", instance.Protocol)

	if locationId != "" {
		_ = d.Set("rule_id", locationId)
	}

	if domain != "" && url != "" {
		_ = d.Set("domain", domain)
		_ = d.Set("url", url)
	}

	var onlineTargets []*clb.Backend
	if *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTP || *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTPS {
		if len(instance.Rules) > 0 {
			for _, loc := range instance.Rules {
				if (locationId == *loc.LocationId) || (domain == *loc.Domain && url == *loc.Url) {
					onlineTargets = loc.Targets
					break
				}
			}
		}
		// TCP / UDP / TCP_SSL
	} else if instance.Targets != nil {
		onlineTargets = instance.Targets
	}

	targets := make([]interface{}, 0)
	for _, onlineTarget := range onlineTargets {
		if *onlineTarget.Type == CLB_BACKEND_TYPE_CVM {
			target := map[string]interface{}{
				"weight":      int(*onlineTarget.Weight),
				"port":        int(*onlineTarget.Port),
				"instance_id": *onlineTarget.InstanceId,
			}

			targets = append(targets, target)
		} else if *onlineTarget.Type == CLB_BACKEND_TYPE_ENI || *onlineTarget.Type == CLB_BACKEND_TYPE_NAT ||
			*onlineTarget.Type == CLB_BACKEND_TYPE_CCN || *onlineTarget.Type == CLB_BACKEND_TYPE_SRV ||
			*onlineTarget.Type == CLB_BACKEND_TYPE_MS || *onlineTarget.Type == CLB_BACKEND_TYPE_EVM ||
			*onlineTarget.Type == CLB_BACKEND_TYPE_GR || *onlineTarget.Type == CLB_BACKEND_TYPE_IPDC ||
			*onlineTarget.Type == CLB_BACKEND_TYPE_PVGW {
			target := map[string]interface{}{
				"weight": int(*onlineTarget.Weight),
				"port":   int(*onlineTarget.Port),
				"eni_ip": *onlineTarget.PrivateIpAddresses[0],
			}

			targets = append(targets, target)
		}
	}

	_ = d.Set("targets", targets)

	return nil
}

func resourceTencentCloudClbServerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_attachment.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	if !d.HasChange("targets") {
		return nil
	}

	o, n := d.GetChange("targets")
	os := o.(*schema.Set)
	ns := n.(*schema.Set)

	// Set.Difference uses the default hash that includes every inner field
	// (instance_id, eni_ip, port, weight). So a pure weight change for the
	// same (instance_id|eni_ip, port) target lands as one entry in `add` AND
	// one entry in `remove`. Partition the raw add/remove sets into three
	// buckets: weightOnly (Modify), pureRemove (Deregister), pureAdd (Register).
	rawAdd := ns.Difference(os).List()
	rawRemove := os.Difference(ns).List()

	weightOnly, pureAdd, pureRemove := partitionTargetChanges(rawAdd, rawRemove)

	// Order: Deregister → Modify → Register.
	// Deregister first frees backend slots; Modify only touches surviving
	// identities; Register re-fills with whatever remains. No identity is
	// ever interpreted twice within a single apply.
	if err := batchProcessTargets(pureRemove, 20, func(chunk []interface{}) error {
		return resourceTencentCloudClbServerAttachmentRemove(d, meta, chunk)
	}); err != nil {
		return err
	}
	if err := batchProcessTargets(weightOnly, 20, func(chunk []interface{}) error {
		return resourceTencentCloudClbServerAttachmentModifyWeight(d, meta, chunk)
	}); err != nil {
		return err
	}
	if err := batchProcessTargets(pureAdd, 20, func(chunk []interface{}) error {
		return resourceTencentCloudClbServerAttachmentAdd(d, meta, chunk)
	}); err != nil {
		return err
	}

	return resourceTencentCloudClbServerAttachmentRead(d, meta)
}

func resourceTencentCloudClbServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_attachment.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		locationId string
		domain     string
		url        string
	)

	attachmentId := d.Id()
	items := strings.Split(attachmentId, "#")
	if len(items) < 3 {
		return fmt.Errorf("[CHECK][CLB attachment][Delete] check: id %s of resource.tencentcloud_clb_attachment is not match loc-xxx#lbl-xxx#lb-xxx", attachmentId)
	}

	locationIdOrDomainUrl := items[0]
	listenerId := items[1]
	clbId := items[2]

	if locationIdOrDomainUrl != "" {
		if strings.HasPrefix(locationIdOrDomainUrl, "loc-") && !strings.Contains(locationIdOrDomainUrl, ",") {
			// get locationId
			locationId = locationIdOrDomainUrl
		} else {
			// get domain & url
			domainUrl := strings.Split(locationIdOrDomainUrl, ",")
			domain = domainUrl[0]
			url = domainUrl[1]
		}
	}

	request := clb.NewDeregisterTargetsRequest()
	request.ListenerId = &listenerId
	request.LoadBalancerId = helper.String(clbId)
	if locationId != "" {
		request.LocationId = helper.String(locationId)
	}

	if domain != "" && url != "" {
		request.Domain = helper.String(domain)
		request.Url = helper.String(url)
	}

	//insList := d.Get("targets").(*schema.Set).List()

	// filter target group which cvm not existed
	insList := getRemoveCandidates(ctx, clbService, clbId, listenerId, locationId, d.Get("targets").(*schema.Set).List())
	insLen := len(insList)
	for count := 0; count < insLen; count += 20 {
		//this request only support 20 targets at most once time
		request.Targets = make([]*clb.Target, 0, 20)
		for i := 0; i < 20; i++ {
			index := count + i
			if index >= insLen {
				break
			}

			inst := insList[index].(map[string]interface{})
			request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["eni_ip"], inst["port"], inst["weight"]))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			requestId := ""
			ratelimit.Check(request.GetAction())
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().DeregisterTargets(request)
			if e != nil {
				ee, ok := e.(*sdkErrors.TencentCloudSDKError)
				if ok && ee.GetCode() == "InvalidParameter" {
					return nil
				}

				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				requestId = *result.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create CLB attachment failed, reason:%+v", logId, err)
			return err
		}
	}

	return nil
}

func resourceTencentCloudClbServerAttachmentRemove(d *schema.ResourceData, meta interface{}, remove []interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_attachment.remove")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		locationId string
		domain     string
		url        string
	)

	attachmentId := d.Id()
	items := strings.Split(attachmentId, "#")
	if len(items) < 3 {
		return fmt.Errorf("[CHECK][CLB attachment][Remove] check: id %s of resource.tencentcloud_clb_attachment is not match loc-xxx#lbl-xxx#lb-xxx", attachmentId)
	}

	locationIdOrDomainUrl := items[0]
	listenerId := items[1]
	clbId := items[2]
	if locationIdOrDomainUrl != "" {
		if strings.HasPrefix(locationIdOrDomainUrl, "loc-") && !strings.Contains(locationIdOrDomainUrl, ",") {
			// get locationId
			locationId = locationIdOrDomainUrl
		} else {
			// get domain & url
			domainUrl := strings.Split(locationIdOrDomainUrl, ",")
			domain = domainUrl[0]
			url = domainUrl[1]
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		removeCandidates := getRemoveCandidates(ctx, clbService, clbId, listenerId, locationId, remove)
		if len(removeCandidates) == 0 {
			return nil
		}

		e := clbService.DeleteAttachmentById(ctx, clbId, listenerId, locationId, removeCandidates, domain, url)
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s reason[%+v]", logId, err)
		return err
	}

	return nil
}

func resourceTencentCloudClbServerAttachmentAdd(d *schema.ResourceData, meta interface{}, add []interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_attachment.add")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = clb.NewRegisterTargetsRequest()
		locationId string
	)

	listenerId := d.Get("listener_id").(string)
	clbId := d.Get("clb_id").(string)
	request.LoadBalancerId = helper.String(clbId)
	request.ListenerId = helper.String(listenerId)
	if v, ok := d.GetOk("rule_id"); ok {
		locationId = v.(string)
		if locationId != "" {
			request.LocationId = helper.String(locationId)
		}
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url"); ok {
		request.Url = helper.String(v.(string))
	}

	for _, v := range add {
		inst := v.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["eni_ip"], inst["port"], inst["weight"]))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		requestId := ""
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().RegisterTargets(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			requestId = *response.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s add CLB attachment failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

// Destroy CVM instance will dispatch async task to deregister target group from cloudApi backend. Duplicate deregister target groups here will cause Error response.
// If remove diffs created, check existing cvm instance first, filter target groups which already deregister
func getRemoveCandidates(ctx context.Context, clbService ClbService, clbId string, listenerId string, locationId string, remove []interface{}) []interface{} {
	removeCandidates := make([]interface{}, 0)
	listenerBackend, err := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
	if err != nil {
		return removeCandidates
	}

	existTargetGroups := make([]*clb.Backend, 0)
	existTargetGroups = append(existTargetGroups, listenerBackend.Targets...)
	if len(listenerBackend.Rules) > 0 {
		existTargetGroups = append(existTargetGroups, listenerBackend.Rules[0].Targets...)
	}

	for _, item := range remove {
		target := item.(map[string]interface{})
		if targetGroupContainsInstance(existTargetGroups, target["instance_id"]) || targetGroupContainsEni(existTargetGroups, target["eni_ip"]) {
			removeCandidates = append(removeCandidates, target)
		}
	}

	return removeCandidates
}

func targetGroupContainsInstance(targets []*clb.Backend, instanceId interface{}) (contains bool) {
	contains = false
	id, ok := instanceId.(string)
	if !ok || id == "" {
		return
	}

	for _, target := range targets {
		if target.InstanceId == nil {
			continue
		}

		if id == *target.InstanceId {
			log.Printf("[WARN] Instance %s applied.", id)
			return true
		}
	}

	log.Printf("[WARN] Instance %s not exist, skip deregister.", id)

	return
}

// targetIdentityKey returns the bucket key used to detect "same target with a
// different weight". Identity is (instance_id|eni_ip, port). `weight` is the
// only mutable, non-identity inner field.
func targetIdentityKey(m map[string]interface{}) string {
	port, _ := m["port"].(int)
	if id, ok := m["instance_id"].(string); ok && id != "" {
		return fmt.Sprintf("inst:%s:%d", id, port)
	}
	if ip, ok := m["eni_ip"].(string); ok && ip != "" {
		return fmt.Sprintf("eni:%s:%d", ip, port)
	}
	// neither key set — treat as unique-per-pointer to avoid false collisions.
	return fmt.Sprintf("anon:%p:%d", m, port)
}

// partitionTargetChanges splits the raw set-diff (add ∪ remove) into three
// disjoint buckets:
//   - weightOnly: same (instance_id|eni_ip, port) appears in both add and remove,
//     differing only by `weight`. Routed to ModifyTargetWeight.
//   - pureAdd:    identity present only in `add`. Routed to RegisterTargets.
//   - pureRemove: identity present only in `remove`. Routed to DeregisterTargets.
//
// `weightOnly` always carries the NEW (post-change) inner map so its weight is
// the value to push to the API.
func partitionTargetChanges(rawAdd, rawRemove []interface{}) (weightOnly, pureAdd, pureRemove []interface{}) {
	addByKey := make(map[string]map[string]interface{}, len(rawAdd))
	for _, item := range rawAdd {
		m := item.(map[string]interface{})
		addByKey[targetIdentityKey(m)] = m
	}
	removeByKey := make(map[string]map[string]interface{}, len(rawRemove))
	for _, item := range rawRemove {
		m := item.(map[string]interface{})
		removeByKey[targetIdentityKey(m)] = m
	}

	for k, addItem := range addByKey {
		if _, hit := removeByKey[k]; hit {
			weightOnly = append(weightOnly, addItem)
			continue
		}
		pureAdd = append(pureAdd, addItem)
	}
	for k, removeItem := range removeByKey {
		if _, hit := addByKey[k]; hit {
			continue
		}
		pureRemove = append(pureRemove, removeItem)
	}
	return
}

// batchProcessTargets walks `items` in chunks of `chunkSize` and calls fn per
// chunk. Returns the first error from fn.
func batchProcessTargets(items []interface{}, chunkSize int, fn func(chunk []interface{}) error) error {
	total := len(items)
	if total == 0 {
		return nil
	}
	for start := 0; start < total; start += chunkSize {
		end := start + chunkSize
		if end > total {
			end = total
		}
		if err := fn(items[start:end]); err != nil {
			return err
		}
	}
	return nil
}

// resourceTencentCloudClbServerAttachmentModifyWeight issues
// ModifyTargetWeight for each (instance_id|eni_ip, port) whose weight changed
// in place. The API accepts a list of Targets where each Target carries its
// own per-target Weight (request-level Weight is omitted, since per-target
// Weight has higher priority anyway).
func resourceTencentCloudClbServerAttachmentModifyWeight(d *schema.ResourceData, meta interface{}, modify []interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_attachment.modifyWeight")()

	if len(modify) == 0 {
		return nil
	}

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = clb.NewModifyTargetWeightRequest()
		locationId string
	)

	listenerId := d.Get("listener_id").(string)
	clbId := d.Get("clb_id").(string)
	request.LoadBalancerId = helper.String(clbId)
	request.ListenerId = helper.String(listenerId)
	if v, ok := d.GetOk("rule_id"); ok {
		locationId = v.(string)
		if locationId != "" {
			request.LocationId = helper.String(locationId)
		}
	}
	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}
	if v, ok := d.GetOk("url"); ok {
		request.Url = helper.String(v.(string))
	}

	for _, v := range modify {
		inst := v.(map[string]interface{})
		request.Targets = append(request.Targets,
			clbNewTarget(inst["instance_id"], inst["eni_ip"], inst["port"], inst["weight"]))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyTargetWeight(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		requestId := *response.Response.RequestId
		if retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient()); retryErr != nil {
			return resource.NonRetryableError(errors.WithStack(retryErr))
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify CLB target weight failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func targetGroupContainsEni(targets []*clb.Backend, eniIp interface{}) (contains bool) {
	contains = false
	ip, ok := eniIp.(string)
	if !ok || ip == "" {
		return
	}

	for _, target := range targets {
		if len(target.PrivateIpAddresses) == 0 || target.PrivateIpAddresses[0] == nil {
			continue
		}

		if ip == *target.PrivateIpAddresses[0] {
			log.Printf("[WARN] IP %s applied.", ip)
			return true
		}
	}

	log.Printf("[WARN] IP %s not exist, skip deregister.", ip)

	return
}
