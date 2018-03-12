package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	vpc "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
)

var errEipUnassigned = errors.New("assigned_eip_set list need at least one")

func resourceTencentCloudNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudNatGatewayCreate,
		Read:   resourceTencentCloudNatGatewayRead,
		Update: resourceTencentCloudNatGatewayUpdate,
		Delete: resourceTencentCloudNatGatewayDelete,

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
			},
			"max_concurrent": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"bandwidth": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"assigned_eip_set": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 10,
			},
		},
	}
}

func resourceTencentCloudNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {

	args := vpc.NewCreateNatGatewayRequest()

	args.VpcId = common.StringPtr(d.Get("vpc_id").(string))
	args.NatName = common.StringPtr(d.Get("name").(string))
	if v, ok := d.GetOk("max_concurrent"); ok {
		args.MaxConcurrent = common.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("bandwidth"); ok {
		args.Bandwidth = common.IntPtr(v.(int))
	}

	eips := d.Get("assigned_eip_set").(*schema.Set).List()
	if len(eips) > 0 {
		args.AssignedEipSet = common.StringPtrs(expandStringList(eips))
	} else {
		return errEipUnassigned
	}

	client := meta.(*TencentCloudClient)
	conn := client.vpcConn
	response, err := conn.CreateNatGateway(args)
	b, _ := json.Marshal(response)
	log.Printf("[DEBUG] conn.CreateNatGateway response: %s", b)
	if _, ok := err.(*common.APIError); ok {
		return fmt.Errorf("conn.CreateNatGateway error: %v", err)
	}

	//Polling NAT gateway production status
	if _, err := client.PollingVpcBillResult(response.BillId); err != nil {
		return err
	}

	log.Printf("[DEBUG] conn.CreateNatGateway NatGatewayId: %s", *response.NatGatewayId)

	d.SetId(*response.NatGatewayId)
	return nil
}

func resourceTencentCloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*TencentCloudClient).vpcConn

	descReq := vpc.NewDescribeNatGatewayRequest()
	descReq.NatId = common.StringPtr(d.Id())

	descResp, err := conn.DescribeNatGateway(descReq)
	b, _ := json.Marshal(descResp)
	log.Printf("[DEBUG] conn.DescribeNatGateway response: %s", b)
	if _, ok := err.(*common.APIError); ok {
		return fmt.Errorf("conn.DescribeNatGateway error: %v", err)
	}

	nat := descResp.Data[0]

	d.Set("name", *nat.NatName)
	d.Set("max_concurrent", *nat.MaxConcurrent)
	d.Set("bandwidth", *nat.Bandwidth)
	d.Set("assigned_eip_set", nat.EipSet)
	return nil
}

func resourceTencentCloudNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*TencentCloudClient)
	conn := client.vpcConn

	d.Partial(true)
	attributeUpdate := false

	updateReq := vpc.NewModifyNatGatewayRequest()
	updateReq.VpcId = common.StringPtr(d.Get("vpc_id").(string))
	updateReq.NatId = common.StringPtr(d.Id())

	if d.HasChange("name") {
		d.SetPartial("name")
		var name string
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		} else {
			return fmt.Errorf("cann't change name to empty string")
		}
		updateReq.NatName = common.StringPtr(name)

		attributeUpdate = true
	}

	if d.HasChange("bandwidth") {
		d.SetPartial("bandwidth")
		var bandwidth int
		if v, ok := d.GetOk("bandwidth"); ok {
			bandwidth = v.(int)
		} else {
			return fmt.Errorf("cann't change bandwidth to empty string")
		}
		updateReq.Bandwidth = common.IntPtr(bandwidth)

		attributeUpdate = true
	}

	if attributeUpdate {
		updateResp, err := conn.ModifyNatGateway(updateReq)
		b, _ := json.Marshal(updateResp)
		log.Printf("[DEBUG] conn.ModifyNatGateway response: %s", b)
		if _, ok := err.(*common.APIError); ok {
			return fmt.Errorf("conn.ModifyNatGateway error: %v", err)
		}
	}

	if d.HasChange("max_concurrent") {
		d.SetPartial("max_concurrent")
		old_mc, new_mc := d.GetChange("max_concurrent")
		old_max_concurrent := old_mc.(int)
		new_max_concurrent := new_mc.(int)
		if new_max_concurrent <= old_max_concurrent {
			return fmt.Errorf("max_concurrent only supports upgrade")
		}
		upgradeReq := vpc.NewUpgradeNatGatewayRequest()
		upgradeReq.VpcId = updateReq.VpcId
		upgradeReq.NatId = updateReq.NatId
		upgradeReq.MaxConcurrent = common.IntPtr(new_max_concurrent)

		upgradeResp, err := conn.UpgradeNatGateway(upgradeReq)
		b, _ := json.Marshal(upgradeResp)
		log.Printf("[DEBUG] conn.UpgradeNatGateway response: %s", b)
		if _, ok := err.(*common.APIError); ok {
			return fmt.Errorf("conn.UpgradeNatGateway error: %v", err)
		}

		if _, err := client.PollingVpcBillResult(upgradeResp.BillId); err != nil {
			return err
		}
	}

	if d.HasChange("assigned_eip_set") {
		o, n := d.GetChange("assigned_eip_set")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		old_eip_set := os.List()
		new_eip_set := ns.List()

		if len(old_eip_set) > 0 && len(new_eip_set) > 0 {

			// Unassign old EIP
			unassignIps := os.Difference(ns)
			if unassignIps.Len() != 0 {
				unbindReq := vpc.NewEipUnBindNatGatewayRequest()
				unbindReq.VpcId = updateReq.VpcId
				unbindReq.NatId = updateReq.NatId
				unbindReq.AssignedEipSet = common.StringPtrs(expandStringList(unassignIps.List()))
				unbindResp, err := conn.EipUnBindNatGateway(unbindReq)
				b, _ := json.Marshal(unbindResp)
				log.Printf("[DEBUG] conn.EipUnBindNatGateway response: %s", b)
				if _, ok := err.(*common.APIError); ok {
					return fmt.Errorf("conn.EipUnBindNatGateway error: %v", err)
				}

				if _, err := client.PollingVpcTaskResult(unbindResp.TaskId); err != nil {
					return err
				}
			}

			// Assign new EIP
			assignIps := ns.Difference(os)
			if assignIps.Len() != 0 {
				bindReq := vpc.NewEipBindNatGatewayRequest()
				bindReq.VpcId = updateReq.VpcId
				bindReq.NatId = updateReq.NatId
				bindReq.AssignedEipSet = common.StringPtrs(expandStringList(assignIps.List()))
				bindResp, err := conn.EipBindNatGateway(bindReq)
				b, _ := json.Marshal(bindResp)
				log.Printf("[DEBUG] conn.EipBindNatGateway response: %s", b)
				if _, ok := err.(*common.APIError); ok {
					return fmt.Errorf("conn.EipBindNatGateway error: %v", err)
				}

				if _, err := client.PollingVpcTaskResult(bindResp.TaskId); err != nil {
					return err
				}
			}

		} else {
			return errEipUnassigned
		}

		d.SetPartial("assigned_eip_set")
	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*TencentCloudClient)

	deleteReq := vpc.NewDeleteNatGatewayRequest()
	deleteReq.VpcId = common.StringPtr(d.Get("vpc_id").(string))
	deleteReq.NatId = common.StringPtr(d.Id())

	deleteResp, err := client.vpcConn.DeleteNatGateway(deleteReq)
	b, _ := json.Marshal(deleteResp)
	log.Printf("[DEBUG] client.vpcConn.DeleteNatGateway response: %s", b)
	if _, ok := err.(*common.APIError); ok {
		return fmt.Errorf("[ERROR] client.vpcConn.DeleteNatGateway error: %v", err)
	}

	_, err = client.PollingVpcTaskResult(deleteResp.TaskId)
	return err
}
