package tencentcloud

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	lb "github.com/zqfan/tencentcloud-sdk-go/services/lb/unversioned"
	"strings"
	"sync"
	"time"
)

var lbActionMu = &sync.Mutex{}

func resourceTencentCloudAlbServerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAlbServerAttachmentCreate,
		Read:   resourceTencentCloudAlbServerAttachmentRead,
		Delete: resourceTencentCloudAlbServerAttachmentDelete,
		Update: resourceTencentCloudAlbServerAttachmentUpdate,

		Schema: map[string]*schema.Schema{
			"loadbalancer_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"listener_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"protocol_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"backends": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(0, 65535),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 100),
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAlbServerAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).lbConn
	lbActionMu.Lock()
	defer lbActionMu.Unlock()

	req := lb.NewRegisterInstancesWithForwardLBSeventhListenerRequest()
	req.LoadBalancerId = common.StringPtr(d.Get("loadbalancer_id").(string))
	req.ListenerId = common.StringPtr(d.Get("listener_id").(string))

	if loc, ok := d.GetOk("location_id"); ok {
		req.LocationIds = []*string{}
		l := loc.(string)
		req.LocationIds = append(req.LocationIds, &l)
	}
	for _, inst_ := range d.Get("backends").(*schema.Set).List() {
		inst := inst_.(map[string]interface{})
		req.Backends = append(req.Backends, lbNewBackend(inst["instance_id"], inst["port"], inst["weight"]))
	}
	resp, err := client.RegisterInstancesWithForwardLBSeventhListener(req)
	if err != nil {
		return err
	}

	err = lbRequestStatusCheck(m, resp.RequestId)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%v:%v:%v", d.Get("loadbalancer_id"), d.Get("listener_id"), d.Get("location_id"))
	d.SetId(id)

	return resourceTencentCloudAlbServerAttachmentRead(d, m)
}

func resourceTencentCloudAlbServerAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).lbConn
	lbActionMu.Lock()
	defer lbActionMu.Unlock()

	req := lb.NewDeregisterInstancesFromForwardLBRequest()
	req.LoadBalancerId = common.StringPtr(d.Get("loadbalancer_id").(string))
	req.ListenerId = common.StringPtr(d.Get("listener_id").(string))

	if loc, ok := d.GetOk("location_id"); ok {
		req.LocationIds = []*string{}
		l := loc.(string)
		req.LocationIds = append(req.LocationIds, &l)
	}
	for _, inst_ := range d.Get("backends").(*schema.Set).List() {
		inst := inst_.(map[string]interface{})
		req.Backends = append(req.Backends, lbNewBackend(inst["instance_id"], inst["port"], inst["weight"]))
	}
	resp, err := client.DeregisterInstancesFromForwardLB(req)
	if err != nil {
		return err
	}

	err = lbRequestStatusCheck(m, resp.RequestId)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceTencentCloudAlbServerAttachementRemove(d *schema.ResourceData, m interface{}, remove []interface{}) error {
	client := m.(*TencentCloudClient).lbConn

	req := lb.NewDeregisterInstancesFromForwardLBRequest()
	req.LoadBalancerId = common.StringPtr(d.Get("loadbalancer_id").(string))
	req.ListenerId = common.StringPtr(d.Get("listener_id").(string))

	if loc, ok := d.GetOk("location_id"); ok {
		req.LocationIds = []*string{}
		l := loc.(string)
		req.LocationIds = append(req.LocationIds, &l)
	}
	for _, inst_ := range remove {
		inst := inst_.(map[string]interface{})
		req.Backends = append(req.Backends, lbNewBackend(inst["instance_id"], inst["port"], inst["weight"]))
	}
	resp, err := client.DeregisterInstancesFromForwardLB(req)
	if err != nil {
		// 9003 error code backend server with port not exist, not ncesssary to remove
		if strings.HasPrefix(*resp.Message, "(9003)") {
			return nil
		} else {
			return err
		}
	} else {
		err = lbRequestStatusCheck(m, resp.RequestId)
		if err != nil {
			return err
		}
	}
	return nil
}
func resourceTencentCloudAlbServerAttachementAdd(d *schema.ResourceData, m interface{}, add []interface{}) error {
	client := m.(*TencentCloudClient).lbConn

	req := lb.NewRegisterInstancesWithForwardLBSeventhListenerRequest()
	req.LoadBalancerId = common.StringPtr(d.Get("loadbalancer_id").(string))
	req.ListenerId = common.StringPtr(d.Get("listener_id").(string))

	if loc, ok := d.GetOk("location_id"); ok {
		req.LocationIds = []*string{}
		l := loc.(string)
		req.LocationIds = append(req.LocationIds, &l)
	}
	for _, inst_ := range add {
		inst := inst_.(map[string]interface{})
		req.Backends = append(req.Backends, lbNewBackend(inst["instance_id"], inst["port"], inst["weight"]))
	}
	resp, err := client.RegisterInstancesWithForwardLBSeventhListener(req)
	if err != nil {
		return err
	}

	err = lbRequestStatusCheck(m, resp.RequestId)
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudAlbServerAttachmentUpdate(d *schema.ResourceData, m interface{}) error {
	lbActionMu.Lock()
	defer lbActionMu.Unlock()

	if d.HasChange("backends") {
		o, n := d.GetChange("backends")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		log.Println("os", os.List())
		log.Println("ns", ns.List())
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()
		if len(remove) > 0 {
			err := resourceTencentCloudAlbServerAttachementRemove(d, m, remove)
			if err != nil {
				return err
			}
		}
		if len(add) > 0 {
			err := resourceTencentCloudAlbServerAttachementAdd(d, m, add)
			if err != nil {
				return err
			}
		}
		return resourceTencentCloudAlbServerAttachmentRead(d, m)
	}
	return nil
}

func resourceTencentCloudAlbServerAttachmentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).lbConn

	req := lb.NewDescribeForwardLBBackendsRequest()
	req.LoadBalancerId = common.StringPtr(d.Get("loadbalancer_id").(string))
	req.ListenerIds = common.StringPtrs([]string{d.Get("listener_id").(string)})
	resp, err := client.DescribeForwardLBBackends(req)
	if err != nil {
		return err
	}
	if len(resp.Data) == 0 {
		d.SetId("")
		return errors.New("resource not exist")
	}
	d.Set("protocol_type", *resp.Data[0].ProtocolType)
	var dataSet []map[string]interface{}
	locid := d.Get("location_id").(string)
	makemap := func(instanceId string, port, weight int) map[string]interface{} {
		m := make(map[string]interface{})
		m["instance_id"] = instanceId
		m["port"] = port
		m["weight"] = weight
		return m
	}
	if *resp.Data[0].Protocol == lb.LBListenerProtocolHTTP || *resp.Data[0].Protocol == lb.LBListenerProtocolHTTPS {
		if len(locid) > 0 {
			for _, loc := range resp.Data[0].Rules {
				if *loc.LocationId == locid {
					for _, v := range loc.Backends {
						m := makemap(*v.UnInstanceId, *v.Port, *v.Weight)
						dataSet = append(dataSet, m)
					}
				}
			}
		} else {
			if len(resp.Data[0].Rules) > 0 {
				for _, v := range resp.Data[0].Rules[0].Backends {
					m := makemap(*v.UnInstanceId, *v.Port, *v.Weight)
					dataSet = append(dataSet, m)
				}
			}
		}
	} else if *resp.Data[0].Protocol == lb.LBListenerProtocolTCP || *resp.Data[0].Protocol == lb.LBListenerProtocolUDP {
		for _, v := range resp.Data[0].Backends {
			m := makemap(*v.UnInstanceId, *v.Port, *v.Weight)
			dataSet = append(dataSet, m)
		}
	}
	d.Set("backends", dataSet)
	return nil
}

func lbRequestStatusCheck(m interface{}, requestId *int) error {
	client := m.(*TencentCloudClient).lbConn

	req := lb.NewDescribeLoadBalancersTaskResultRequest()
	req.RequestId = requestId

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.DescribeLoadBalancersTaskResult(req)

		if err != nil {
			return resource.RetryableError(err)
		}
		switch *resp.Data.Status {
		case lb.LBTaskSuccess:
			return nil
		case lb.LBTaskFail:
			return resource.NonRetryableError(errors.New("action fail"))
		case lb.LBTaskDoing:
			// lb action is still doing, wait and retry
			return resource.RetryableError(errors.New("retry"))
		}

		return nil
	})
}

func lbNewBackend(instanceId, port, weight interface{}) *lb.Backend {
	id := instanceId.(string)
	p, _ := port.(int)
	bk := lb.Backend{
		InstanceId: &id,
		Port:       &p,
	}
	if w, ok := weight.(int); ok {
		bk.Weight = &w
	}
	return &bk
}
