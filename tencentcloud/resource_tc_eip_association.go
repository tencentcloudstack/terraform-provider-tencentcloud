package tencentcloud

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/athom/goset"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	cvm "github.com/zqfan/tencentcloud-sdk-go/services/cvm/v20170312"
)

func resourceTencentCloudEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipAssociationCreate,
		Read:   resourceTencentCloudEipAssociationRead,
		Delete: resourceTencentCloudEipAssociationDelete,

		Schema: map[string]*schema.Schema{
			"eip_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
			},

			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"network_interface_id",
					"private_ip",
				},
				ValidateFunc: validateNotEmpty,
			},

			"network_interface_id": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateNotEmpty,
				ConflictsWith: []string{
					"instance_id",
				},
			},

			"private_ip": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateNotEmpty,
				ConflictsWith: []string{
					"instance_id",
				},
			},
		},
	}
}

func resourceTencentCloudEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	cvmConn := meta.(*TencentCloudClient).cvmConn

	v := d.Get("eip_id")
	eipId := v.(string)

	// make sure EIP is in unbind status for better user experience
	eip, _, err := findEipById(cvmConn, eipId)
	if err != nil {
		return err
	}
	if *eip.AddressStatus != tencentCloudApiEipStatusUnbind {
		return errEIPNotUnbind
	}

	// associate instance
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId := v.(string)
		req := cvm.NewAssociateAddressRequest()
		req.AddressId = common.StringPtr(eipId)
		req.InstanceId = common.StringPtr(instanceId)
		_, err := cvmConn.AssociateAddress(req)
		if err != nil {
			return err
		}

		associationId := fmt.Sprintf("%v::%v", eipId, instanceId)
		d.SetId(associationId)
		return resourceTencentCloudEipAssociationRead(d, meta)
	}

	// associate network interface id
	v, ok := d.GetOk("network_interface_id")
	if !ok {
		err = errors.New("network_interface_id is expected to be specified while no instance_id provided")
		return err
	}
	networkInterfaceId := v.(string)

	v, ok = d.GetOk("private_ip")
	if !ok {
		err = errors.New("private_ip is expected to be specified while network_interface_id provided")
		return err
	}
	privateIp := v.(string)
	req := cvm.NewAssociateAddressRequest()
	req.AddressId = common.StringPtr(eipId)
	req.NetworkInterfaceId = common.StringPtr(networkInterfaceId)
	req.PrivateIpAddress = common.StringPtr(privateIp)
	_, err = cvmConn.AssociateAddress(req)
	if err != nil {
		return err
	}

	associationId := fmt.Sprintf("%v::%v::%v", eipId, networkInterfaceId, privateIp)
	d.SetId(associationId)
	return resourceTencentCloudEipAssociationRead(d, meta)
}

func resourceTencentCloudEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	associationId := d.Id()
	association, err := parseAssociationId(associationId)
	if err != nil {
		return err
	}
	d.Set("eip_id", association.eipId)
	// associate with instance
	if len(association.instanceId) > 0 {
		d.Set("instance_id", association.instanceId)
		return nil
	}

	d.Set("network_interface_id", association.networkInterfaceId)
	d.Set("private_ip", association.privateIp)
	return nil
}

func resourceTencentCloudEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	cvmConn := meta.(*TencentCloudClient).cvmConn
	associationId := d.Id()
	association, err := parseAssociationId(associationId)
	if err != nil {
		return err
	}
	eipId := association.eipId

	// NOTE wait until eip is bind
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		eip, _, err := findEipById(cvmConn, eipId)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		status := *eip.AddressStatus
		if goset.IsIncluded([]string{
			tencentCloudApiEipStatusBind,
			tencentCloudApiEipStatusBindEni,
		}, status) {
			req := cvm.NewDisassociateAddressRequest()
			req.AddressId = common.StringPtr(eipId)
			_, err = cvmConn.DisassociateAddress(req)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		}
		return resource.RetryableError(errEIPStillUnbinding)
	})
}

type association struct {
	eipId              string
	instanceId         string
	networkInterfaceId string
	privateIp          string
}

// association id is in a format like: eip-m5vh60me::ins-ojhtwo3k
func parseAssociationId(associationId string) (r association, err error) {
	ids := strings.Split(associationId, "::")
	if len(ids) < 2 || len(ids) > 3 {
		err = fmt.Errorf("Invalid association ID: %v", associationId)
		return
	}
	r.eipId = ids[0]

	// associate with instance
	if len(ids) == 2 {
		r.instanceId = ids[1]
		return
	}

	// associate with network interface
	r.networkInterfaceId = ids[1]
	r.privateIp = ids[2]
	return
}
