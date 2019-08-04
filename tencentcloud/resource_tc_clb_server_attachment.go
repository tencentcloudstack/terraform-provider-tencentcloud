/*
Provide a resource to create a CLB instance.

Example Usage

```hcl
resource "tencentcloud_clb_server_attachment" "attachment" {
  listener_id   = "lbl-hh141sn9#lb-k2zjp9lv"
  clb_id        = "lb-k2zjp9lv"
  protocol_type = "tcp"
  location_id   = "loc-4xxr2cy7"
  targets = {
    instance_id = "ins-1flbqyp8"
    port        = 50
    weight      = 10
  }
}
```

Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_server_attachment.attachment loc-4xxr2cy7#lbl-hh141sn9#lb-7a0t6zqb
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
				Description: "Id of the cloud load balancer. ",
			},

			"listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Id of the cloud load balance listener. ",
			},
			"protocol_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of protocol within the listener, and available values include 'TCP', 'UDP', 'HTTP', 'HTTPS' and 'TCP_SSL'. ",
			},
			"location_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Id of the cloud load balance listener rule. ",
			},
			"targets": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    100,
				Description: "Backend infos.",
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
							Description:  "Weight of the backend server.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClbServerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	//ctx := context.WithValue(context.TODO(), "logId", logId)

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
	if v, ok := d.GetOk("location_id"); ok {
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

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	defer LogElapsed(logId + "resource.tencentcloud_clb_listener_rule.delete")()

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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	defer LogElapsed(logId + "resource.tencentcloud_clb_listener_rule.delete")()

	attachmentId := d.Id()
	items := strings.Split(attachmentId, "#")
	if len(items) < 3 {
		return fmt.Errorf("id of resource.tencentcloud_clb_server_attachment is wrong")
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
func resourceTencentCloudClbServerAttachementAdd(d *schema.ResourceData, meta interface{}, add []interface{}) error {
	logId := GetLogId(nil)
	items := strings.Split(d.Get("listener_id").(string), "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_clb_server_attachment is wrong")
	}

	listenerId := items[0]
	clbId := items[1]
	locationId := ""
	request := clb.NewRegisterTargetsRequest()
	request.LoadBalancerId = stringToPointer(clbId)
	request.ListenerId = stringToPointer(listenerId)

	if v, ok := d.GetOk("location_id"); ok {
		items := strings.Split(v.(string), "#")
		locationId = items[0]
		request.LocationId = stringToPointer(locationId)
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
	return nil
}

func resourceTencentCloudClbServerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("targets") {
		o, n := d.GetChange("targets")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		log.Println("os", os.List())
		log.Println("ns", ns.List())
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
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_clb_instance.read")()
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
	d.Set("location_id", locationId+"#"+listenerId+"#"+clbId)
	if *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTP || *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTPS {
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
