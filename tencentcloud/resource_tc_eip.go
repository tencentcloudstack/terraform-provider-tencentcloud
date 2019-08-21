package tencentcloud

import (
	"errors"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	cvm "github.com/zqfan/tencentcloud-sdk-go/services/cvm/v20170312"
)

const (
	tencentCloudApiEipStatusBinding      = "BINDING"
	tencentCloudApiEipStatusBind         = "BIND"
	tencentCloudApiEipStatusUnbind       = "UNBIND"
	tencentCloudApiEipStatusCreateFailed = "CREATE_FAILED"
	tencentCloudApiEipStatusBindEni      = "BIND_ENI"
	// tencentCloudApiEipStatusCreating     = "CREATING"
	// tencentCloudApiEipStatusOfflining    = "OFFLINING"
	// tencentCloudApiEipStatusUnbinding    = "UNBINDING"
)

var (
	errCreateEIPFailed   = errors.New("create eip failed")
	errEIPStillBinding   = errors.New("eip still binding")
	errEIPStillUnbinding = errors.New("eip still unbinding")
	errEIPStillCreating  = errors.New("eip still creating")
	errEIPStillDeleting  = errors.New("eip still deleting")
	errEIPNotUnbind      = errors.New("eip should be unbind")
	errEIPInvalidName    = errors.New("eip name is invalid")
	// errEIPNotBind        = errors.New("eip should be bind")
)

func resourceTencentCloudEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipCreate,
		Read:   resourceTencentCloudEipRead,
		Update: resourceTencentCloudEipUpdate,
		Delete: resourceTencentCloudEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(1, 20),
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.create")()

	cvmConn := meta.(*TencentCloudClient).cvmConn

	req := cvm.NewAllocateAddressesRequest()
	req.AddressCount = common.IntPtr(1)
	resp, err := cvmConn.AllocateAddresses(req)
	if err != nil {
		return err
	}
	eipIds := resp.Response.AddressSet
	if len(eipIds) == 0 {
		return errCreateEIPFailed
	}
	eipId := eipIds[0]
	err = waitForEipAvailable(cvmConn, *eipId)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		err = setEipName(cvmConn, *eipId, name)
		if err != nil {
			return err
		}
	}

	d.SetId(*eipId)
	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.read")()

	cvmConn := meta.(*TencentCloudClient).cvmConn
	eipId := d.Id()

	eip, _, err := findEipById(cvmConn, eipId)
	if err != nil {
		if err == errEIPNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("public_ip", *eip.AddressIp)
	d.Set("status", *eip.AddressStatus)
	if eip.AddressName != nil {
		d.Set("name", *eip.AddressName)
	}
	return nil
}

func resourceTencentCloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.update")()

	if d.HasChange("name") {
		eipId := d.Id()
		cvmConn := meta.(*TencentCloudClient).cvmConn
		v, ok := d.GetOk("name")
		if !ok {
			return errEIPInvalidName
		}
		newName := v.(string)
		err := setEipName(cvmConn, eipId, newName)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.delete")()

	cvmConn := meta.(*TencentCloudClient).cvmConn
	eipId := d.Id()

	// NOTE wait until eip is unbind
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		eip, _, err := findEipById(cvmConn, eipId)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		status := *eip.AddressStatus
		if status == tencentCloudApiEipStatusUnbind {
			req := cvm.NewReleaseAddressesRequest()
			req.AddressIds = []*string{
				common.StringPtr(eipId),
			}
			_, err = cvmConn.ReleaseAddresses(req)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		}

		return resource.RetryableError(errEIPStillDeleting)
	})
}
