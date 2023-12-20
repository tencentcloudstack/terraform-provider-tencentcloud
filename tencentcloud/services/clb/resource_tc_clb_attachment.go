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
		Delete: resourceTencentCloudClbServerAttachmentDelete,
		Update: resourceTencentCloudClbServerAttachmentUpdate,
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
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of the CLB listener rule. Only supports listeners of `HTTPS` and `HTTP` protocol.",
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
							Description: "CVM Instance Id of the backend server, conflict with `eni_ip` but must specify one of them.",
						},
						"eni_ip": {
							Type:        schema.TypeString,
							Optional:    true,
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

	logId := tccommon.GetLogId(tccommon.ContextNil)
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	clbId := d.Get("clb_id").(string)
	locationId := ""
	request := clb.NewRegisterTargetsRequest()
	request.LoadBalancerId = helper.String(clbId)
	request.ListenerId = helper.String(listenerId)
	if v, ok := d.GetOk("rule_id"); ok {
		locationId = v.(string)
		checkErr := RuleIdCheck(locationId)
		if checkErr != nil {
			return checkErr
		}
		if locationId != "" {
			request.LocationId = helper.String(locationId)
		}
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
	id := fmt.Sprintf("%s#%v#%v", locationId, d.Get("listener_id"), d.Get("clb_id"))
	d.SetId(id)

	return resourceTencentCloudClbServerAttachmentRead(d, meta)
}

func resourceTencentCloudClbServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_attachment.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	clbService := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	attachmentId := d.Id()

	items := strings.Split(attachmentId, "#")
	if len(items) < 3 {
		return fmt.Errorf("[CHECK][CLB attachment][Delete] check: id %s of resource.tencentcloud_clb_attachment is not match loc-xxx#lbl-xxx#lb-xxx", attachmentId)
	}

	locationId := items[0]
	listenerId := items[1]
	clbId := items[2]

	request := clb.NewDeregisterTargetsRequest()
	request.ListenerId = &listenerId
	request.LoadBalancerId = helper.String(clbId)
	if locationId != "" {
		request.LocationId = helper.String(locationId)
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

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	attachmentId := d.Id()
	items := strings.Split(attachmentId, "#")
	if len(items) < 3 {
		return fmt.Errorf("[CHECK][CLB attachment][Remove] check: id %s of resource.tencentcloud_clb_attachment is not match loc-xxx#lbl-xxx#lb-xxx", attachmentId)
	}
	locationId := items[0]
	listenerId := items[1]
	clbId := items[2]

	request := clb.NewDeregisterTargetsRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(clbId)
	if locationId != "" {
		request.LocationId = helper.String(locationId)
	}

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		removeCandidates := getRemoveCandidates(ctx, clbService, clbId, listenerId, locationId, remove)
		if len(removeCandidates) == 0 {
			return nil
		}
		e := clbService.DeleteAttachmentById(ctx, clbId, listenerId, locationId, removeCandidates)
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
	logId := tccommon.GetLogId(tccommon.ContextNil)

	listenerId := d.Get("listener_id").(string)
	clbId := d.Get("clb_id").(string)
	locationId := ""
	request := clb.NewRegisterTargetsRequest()
	request.LoadBalancerId = helper.String(clbId)
	request.ListenerId = helper.String(listenerId)

	if v, ok := d.GetOk("rule_id"); ok {
		locationId = v.(string)
		if locationId != "" {
			request.LocationId = helper.String(locationId)
		}
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

func resourceTencentCloudClbServerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_attachment.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	if d.HasChange("targets") {
		o, n := d.GetChange("targets")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()
		addLen := len(add)
		removeLen := len(remove)
		if removeLen > 0 {
			for count := 0; count < removeLen; count += 20 {
				removeList := make([]interface{}, 0, 20)
				for i := 0; i < 20; i++ {
					index := count + i
					if index >= removeLen {
						break
					}
					removeList = append(removeList, remove[index])
				}
				err := resourceTencentCloudClbServerAttachmentRemove(d, meta, removeList)
				if err != nil {
					return err
				}
			}
		}
		if addLen > 0 {
			for count := 0; count < addLen; count += 20 {
				addList := make([]interface{}, 0, 20)
				for i := 0; i < 20; i++ {
					index := count + i
					if index >= addLen {
						break
					}
					addList = append(addList, add[index])
				}
				err := resourceTencentCloudClbServerAttachmentAdd(d, meta, addList)
				if err != nil {
					return err
				}
			}
		}
		return resourceTencentCloudClbServerAttachmentRead(d, meta)
	}

	return nil
}

func resourceTencentCloudClbServerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), "#")
	locationId := items[0]
	listenerId := items[1]
	clbId := items[2]

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var instance *clb.ListenerBackend
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

	if instance == nil || (len(instance.Targets) == 0 && locationId == "") || (len(instance.Rules) == 0 && locationId != "") {
		d.SetId("")
		return nil
	}

	_ = d.Set("clb_id", clbId)
	_ = d.Set("listener_id", listenerId)
	_ = d.Set("protocol_type", instance.Protocol)

	var onlineTargets []*clb.Backend
	if *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTP || *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTPS {
		_ = d.Set("rule_id", locationId)
		if len(instance.Rules) > 0 {
			for _, loc := range instance.Rules {
				if locationId == "" || locationId == *loc.LocationId {
					onlineTargets = loc.Targets
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
		} else if *onlineTarget.Type == CLB_BACKEND_TYPE_ENI {
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
