package tencentcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	lb "github.com/zqfan/tencentcloud-sdk-go/services/lb/unversioned"
)

const (
	lbNetworkTypeOpen        = "OPEN"
	lbNetworkTypeInternal    = "INTERNAL"
	lbForwardTypeClassic     = "CLASSIC"
	lbForwardTypeApplication = "APPLICATION"
)

func resourceTencentCloudLB() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.15.0. Please use 'tencentcloud_clb_instance' instead.",
		Create:             resourceTencentCloudLBCreate,
		Read:               resourceTencentCloudLBRead,
		Update:             resourceTencentCloudLBUpdate,
		Delete:             resourceTencentCloudLBDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != lbNetworkTypeOpen && value != lbNetworkTypeInternal {
						errors = append(errors, fmt.Errorf("Invalid value %s for '%s', choices: OPEN, INTERNAL", value, k))
					}
					return
				},
			},
			"forward": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != lbForwardTypeClassic && value != lbForwardTypeApplication {
						errors = append(errors, fmt.Errorf("Invalid value %s for '%s', choices: CLASSIC, APPLICATION", value, k))
					}
					return
				},
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
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudLBCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).lbConn
	req := lb.NewCreateLoadBalancerRequest()
	if d.Get("type").(string) == lbNetworkTypeOpen {
		req.LoadBalancerType = common.IntPtr(lb.LBNetworkTypePublic)
	} else {
		req.LoadBalancerType = common.IntPtr(lb.LBNetworkTypePrivate)
	}
	if t, ok := d.GetOk("forward"); ok {
		if t.(string) == lbForwardTypeClassic {
			req.Forward = common.IntPtr(lb.LBForwardTypeClassic)
		} else {
			req.Forward = common.IntPtr(lb.LBForwardTypeApplication)
		}
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
	if *resp.LoadBalancerSet[0].LoadBalancerType == lb.LBNetworkTypePublic {
		d.Set("type", lbNetworkTypeOpen)
	} else {
		d.Set("type", lbNetworkTypeInternal)
	}
	if *resp.LoadBalancerSet[0].Forward == lb.LBForwardTypeClassic {
		d.Set("forward", lbForwardTypeClassic)
	} else {
		d.Set("forward", lbForwardTypeApplication)
	}
	d.Set("name", *resp.LoadBalancerSet[0].LoadBalancerName)
	d.Set("vpc_id", *resp.LoadBalancerSet[0].VpcId)
	d.Set("project_id", *resp.LoadBalancerSet[0].ProjectId)
	if *resp.LoadBalancerSet[0].Status == lb.LBStatusReady {
		d.Set("status", "NORMAL")
	} else if *resp.LoadBalancerSet[0].Status == lb.LBStatusCreating {
		d.Set("status", "CREATING")
	} else {
		d.Set("status", "UNKNOWN")
	}
	return nil
}

func resourceTencentCloudLBUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).lbConn
	lbid := d.Id()
	if !d.HasChange("name") {
		return nil
	}
	v, _ := d.GetOk("name")
	if f := d.Get("forward").(string); f == lbForwardTypeApplication {
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
