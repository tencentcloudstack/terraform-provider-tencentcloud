package tencentcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	lb "github.com/zqfan/tencentcloud-sdk-go/services/lb/unversioned"
)

func resourceTencentCloudLB() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLBCreate,
		Read:   resourceTencentCloudLBRead,
		Update: resourceTencentCloudLBUpdate,
		Delete: resourceTencentCloudLBDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"forward": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudLBCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).lbConn
	req := lb.NewCreateLoadBalancerRequest()
	req.LoadBalancerType = common.IntPtr(d.Get("type").(int))
	if t, ok := d.GetOk("forward"); ok {
		req.Forward = common.IntPtr(t.(int))
	}
	if n, ok := d.GetOk("name"); ok {
		req.LoadBalancerName = common.StringPtr(n.(string))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		req.VpcId = common.StringPtr(v.(string))
	}
	if p, ok := d.GetOk("project_id"); ok {
		req.ProjectId = common.IntPtr(p.(int))
	}
	resp, err := client.CreateLoadBalancer(req)
	if err != nil {
		return err
	}
	dealId := *resp.DealIds[0]
	lbid := (*resp.UnLoadBalancerIds)[dealId][0]
	err = waitForLBReady(client, lbid)
	if err != nil {
		return err
	}
	d.SetId(*lbid)
	return resourceTencentCloudLBRead(d, meta)
}

func resourceTencentCloudLBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).lbConn
	lbid := d.Id()
	req := lb.NewDescribeLoadBalancersRequest()
	req.LoadBalancerIds = []*string{&lbid}
	resp, err := client.DescribeLoadBalancers(req)
	if err != nil {
		return err
	}
	if *resp.TotalCount == 0 {
		d.SetId("")
		return nil
	}
	d.Set("type", *resp.LoadBalancerSet[0].LoadBalancerType)
	d.Set("forward", *resp.LoadBalancerSet[0].Forward)
	d.Set("name", *resp.LoadBalancerSet[0].LoadBalancerName)
	d.Set("vpc_id", *resp.LoadBalancerSet[0].VpcId)
	d.Set("project_id", *resp.LoadBalancerSet[0].ProjectId)
	d.Set("status", *resp.LoadBalancerSet[0].Status)
	return nil
}

func resourceTencentCloudLBUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).lbConn
	lbid := d.Id()
	if !d.HasChange("name") {
		return nil
	}
	v, _ := d.GetOk("name")
	if f := d.Get("forward").(int); f == lb.LBForwardTypeApplication {
		req := lb.NewModifyForwardLBNameRequest()
		req.LoadBalancerId = common.StringPtr(lbid)
		req.LoadBalancerName = common.StringPtr(v.(string))
		_, err := client.ModifyForwardLBName(req)
		if err != nil {
			return err
		}
	} else {
		req := lb.NewModifyLoadBalancerAttributesRequest()
		req.LoadBalancerId = common.StringPtr(lbid)
		req.LoadBalancerName = common.StringPtr(v.(string))
		resp, err := client.ModifyLoadBalancerAttributes(req)
		if err != nil {
			return err
		}
		if err := waitForLBTaskFinish(client, resp.RequestId); err != nil {
			return err
		}
	}
	return resourceTencentCloudLBRead(d, meta)
}

func resourceTencentCloudLBDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).lbConn
	lbid := d.Id()
	req := lb.NewDeleteLoadBalancersRequest()
	req.LoadBalancerIds = []*string{&lbid}
	resp, err := client.DeleteLoadBalancers(req)
	if err != nil {
		return err
	}
	taskid := resp.RequestId
	if err := waitForLBTaskFinish(client, taskid); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
