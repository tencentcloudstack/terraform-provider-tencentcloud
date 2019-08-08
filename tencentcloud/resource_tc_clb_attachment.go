/*
Provides a resource to create a CLB attachment.

Example Usage

```hcl
resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9#lb-k2zjp9lv"
  rule_id     = "loc-4xxr2cy7#lbl-hh141sn9#lb-k2zjp9lv"

  targets {
    instance_id = "ins-1flbqyp8"
    port        = 80
    weight      = 10
  }
}
```

Import

CLB attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_attachment.foo loc-4xxr2cy7#lbl-hh141sn9#lb-7a0t6zqb
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbServerAttachment() *schema.Resource {
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
				Description: "ID of the clb.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: " ID of the clb listener.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of the clb listener rule. Only supports listeners of 'HTTPS'/'HTTP' protocol.",
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
				Description: " Information of the backends to be attached.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Id of the backend server.",
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 65535),
							Description:  "Port of the backend server.",
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 100),
							Description:  "Forwarding weight of the backend service, the range of [0, 100], defaults to 10.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClbServerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_attachment.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(nil)

	items := strings.Split(d.Get("listener_id").(string), "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_clb_listener is wrong ?%d %s", len(items), d.Get("listener_id").(string))
	}

	listenerId := items[0]
	clbId := items[1]
	locationId := ""
	request := clb.NewRegisterTargetsRequest()
	request.LoadBalancerId = stringToPointer(clbId)
	request.ListenerId = stringToPointer(listenerId)
	if v, ok := d.GetOk("rule_id"); ok {
		items := strings.Split(v.(string), "#")
		locationId = items[0]
		if locationId != "" {
			request.LocationId = stringToPointer(locationId)
		}
	}

	for _, inst_ := range d.Get("targets").(*schema.Set).List() {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["port"], inst["weight"]))
	}

	requestId := ""
	response, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().RegisterTargets(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		requestId = *response.Response.RequestId
		retryErr := retrySet(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
		if retryErr != nil {
			return retryErr
		}
	}

	id := fmt.Sprintf("%s#%v", locationId, d.Get("listener_id"))
	d.SetId(id)

	return resourceTencentCloudClbServerAttachmentRead(d, meta)
}

func resourceTencentCloudClbServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_attachment.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	attachmentId := d.Id()
	items := strings.Split(attachmentId, "#")
	if len(items) < 3 {
		return fmt.Errorf("id of resource.tencentcloud_clb_listener is wrong")
	}
	locationId := items[0]
	listenerId := items[1]
	clbId := items[2]
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := clbService.DeleteAttachmentById(ctx, clbId, listenerId, locationId, d.Get("targets").(*schema.Set).List())

	if err != nil {
		log.Printf("[CRITAL]%s reason[%s]\n", logId, err.Error())
		return err
	}

	return nil
}

func resourceTencentCloudClbServerAttachementRemove(d *schema.ResourceData, meta interface{}, remove []interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_attachment.remove")()

	logId := getLogId(nil)
	attachmentId := d.Id()
	items := strings.Split(attachmentId, "#")
	if len(items) < 3 {
		return fmt.Errorf("id of resource.tencentcloud_clb_attachment is wrong")
	}
	locationId := items[0]
	listenerId := items[1]
	clbId := items[2]
	request := clb.NewDeregisterTargetsRequest()
	request.ListenerId = stringToPointer(listenerId)
	request.LoadBalancerId = stringToPointer(clbId)
	if locationId != "" {
		request.LocationId = stringToPointer(locationId)
	}
	for _, inst_ := range remove {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["port"], inst["weight"]))
	}

	requestId := ""
	response, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().DeregisterTargets(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		requestId = *response.Response.RequestId
		retryErr := retrySet(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
		if retryErr != nil {
			return retryErr
		}
	}
	return nil
}

func resourceTencentCloudClbServerAttachementAdd(d *schema.ResourceData, meta interface{}, add []interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_attachment.add")()
	logId := getLogId(nil)
	items := strings.Split(d.Get("listener_id").(string), "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_clb_attachment is wrong")
	}

	listenerId := items[0]
	clbId := items[1]
	locationId := ""
	request := clb.NewRegisterTargetsRequest()
	request.LoadBalancerId = stringToPointer(clbId)
	request.ListenerId = stringToPointer(listenerId)

	if v, ok := d.GetOk("rule_id"); ok {
		items := strings.Split(v.(string), "#")
		locationId = items[0]
		if locationId != "" {
			request.LocationId = stringToPointer(locationId)
		}
	}

	for _, inst_ := range add {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["port"], inst["weight"]))
	}

	requestId := ""
	response, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().RegisterTargets(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		requestId = *response.Response.RequestId
		retryErr := retrySet(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
		if retryErr != nil {
			return retryErr
		}
	}
	return nil
}

func resourceTencentCloudClbServerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_attachment.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	if d.HasChange("targets") {
		o, n := d.GetChange("targets")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()
		if len(remove) > 0 {
			err := resourceTencentCloudClbServerAttachementRemove(d, meta, remove)
			if err != nil {
				return err
			}
		}
		if len(add) > 0 {
			err := resourceTencentCloudClbServerAttachementAdd(d, meta, add)
			if err != nil {
				return err
			}
		}
		return resourceTencentCloudClbServerAttachmentRead(d, meta)
	}

	return nil
}

func resourceTencentCloudClbServerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_attachment.read")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), "#")
	locationId := items[0]
	listenerId := items[1]
	clbId := items[2]

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instance, err := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
	if err != nil {
		return err
	}

	d.Set("clb_id", clbId)
	d.Set("listener_id", listenerId+"#"+clbId)
	d.Set("protocol_type", instance.Protocol)

	if *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTP || *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTPS {
		d.Set("rule_id", locationId+"#"+listenerId+"#"+clbId)
		if len(instance.Rules) > 0 {
			for _, loc := range instance.Rules {
				if locationId == "" || locationId == *loc.LocationId {
					d.Set("targets", flattenBackendList(loc.Targets))
				}
			}
		}
	} else if *instance.Protocol == CLB_LISTENER_PROTOCOL_TCP || *instance.Protocol == CLB_LISTENER_PROTOCOL_UDP {
		d.Set("targets", flattenBackendList(instance.Targets))
	}

	return nil
}
