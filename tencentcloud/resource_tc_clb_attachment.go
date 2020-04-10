/*
Provides a resource to create a CLB attachment.

Example Usage

```hcl
resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  rule_id     = "loc-4xxr2cy7"

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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Description: "Id of the CLB.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Id of the CLB listener.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Id of the CLB listener rule. Only supports listeners of 'HTTPS' and 'HTTP' protocol.",
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
							Default:      10,
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

	logId := getLogId(contextNil)
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

	for _, inst_ := range d.Get("targets").(*schema.Set).List() {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["port"], inst["weight"]))
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		requestId := ""
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().RegisterTargets(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			requestId = *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
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
	id := fmt.Sprintf("%s#%v#%v", locationId, d.Get("listener_id"), d.Get("clb_id"))
	d.SetId(id)

	return resourceTencentCloudClbServerAttachmentRead(d, meta)
}

func resourceTencentCloudClbServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_attachment.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	attachmentId := d.Id()

	items := strings.Split(attachmentId, "#")
	if len(items) < 3 {
		return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB attachment][Delete] check: id %s of resource.tencentcloud_clb_attachment is not match loc-xxx#lbl-xxx#lb-xxx", attachmentId)
	}

	locationId := items[0]
	listenerId := items[1]
	clbId := items[2]

	//check exists
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	request := clb.NewDeregisterTargetsRequest()
	request.ListenerId = &listenerId
	request.LoadBalancerId = helper.String(clbId)
	if locationId != "" {
		request.LocationId = helper.String(locationId)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteAttachmentById(ctx, clbId, listenerId, locationId, d.Get("targets").(*schema.Set).List())
		if e != nil {
			return retryError(e)
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s reason[%+v]", logId, err)
		return err
	}

	return nil
}

func resourceTencentCloudClbServerAttachementRemove(d *schema.ResourceData, meta interface{}, remove []interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_attachment.remove")()

	logId := getLogId(contextNil)

	attachmentId := d.Id()
	items := strings.Split(attachmentId, "#")
	if len(items) < 3 {
		return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB attachment][Remove] check: id %s of resource.tencentcloud_clb_attachment is not match loc-xxx#lbl-xxx#lb-xxx", attachmentId)
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

	for _, inst_ := range remove {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["port"], inst["weight"]))

	}

	if len(request.Targets) > 0 {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			requestId := ""
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().DeregisterTargets(request)
			if e != nil {
				ee, ok := e.(*sdkErrors.TencentCloudSDKError)
				if !ok {
					return retryError(ee)
				}
				if ee.Code == "InvalidParameter" {
					return nil
				} else {
					return retryError(e)
				}

			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId = *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s remove CLB attachment failed, reason:%+v", logId, err)
			return err
		}
	}
	return nil
}

func resourceTencentCloudClbServerAttachementAdd(d *schema.ResourceData, meta interface{}, add []interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_attachment.add")()
	logId := getLogId(contextNil)

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

	for _, inst_ := range add {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["port"], inst["weight"]))

	}

	if len(request.Targets) > 0 {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			requestId := ""
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().RegisterTargets(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId = *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
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

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), "#")
	locationId := items[0]
	listenerId := items[1]
	clbId := items[2]

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *clb.ListenerBackend
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if e != nil {
			return retryError(e)
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

	if *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTP || *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTPS {
		_ = d.Set("rule_id", locationId)
		if len(instance.Rules) > 0 {
			for _, loc := range instance.Rules {
				if locationId == "" || locationId == *loc.LocationId {
					_ = d.Set("targets", flattenBackendList(loc.Targets))
				}
			}
		}
	} else if *instance.Protocol == CLB_LISTENER_PROTOCOL_TCP || *instance.Protocol == CLB_LISTENER_PROTOCOL_UDP {
		_ = d.Set("targets", flattenBackendList(instance.Targets))
	}

	return nil
}
