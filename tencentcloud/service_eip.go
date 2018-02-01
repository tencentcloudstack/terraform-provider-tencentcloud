package tencentcloud

import (
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	cvm "github.com/zqfan/tencentcloud-sdk-go/services/cvm/v20170312"
)

func findEipById(cvmConn *cvm.Client, eipId string) (eip *cvm.Address, retryable bool, err error) {
	req := cvm.NewDescribeAddressesRequest()
	req.Filters = []*cvm.Filter{
		&cvm.Filter{
			Name:   common.StringPtr("address-id"),
			Values: []*string{common.StringPtr(eipId)},
		},
	}
	req.Limit = common.IntPtr(1)
	resp, err := cvmConn.DescribeAddresses(req)
	if err != nil {
		retryable = false
		return
	}
	if *resp.Response.TotalCount == 0 {
		err = errEIPNotFound
		retryable = false
		return
	}
	eips := resp.Response.AddressSet
	if len(eips) != 1 {
		err = errEIPNotFound
		retryable = false
		return
	}
	eip = eips[0]
	return
}

func setEipName(cvmConn *cvm.Client, eipId string, newName string) (err error) {
	req := cvm.NewModifyAddressAttributeRequest()
	req.AddressId = common.StringPtr(eipId)
	req.AddressName = common.StringPtr(newName)
	_, err = cvmConn.ModifyAddressAttribute(req)
	if err != nil {
		return err
	}
	return nil
}

func waitForEipAvailable(cvmConn *cvm.Client, eipId string) (err error) {
	resource.Retry(3*time.Minute, func() *resource.RetryError {
		eip, retryable, err := findEipById(cvmConn, eipId)
		if err != nil {
			if retryable {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		status := *eip.AddressStatus
		if status == tencentCloudApiEipStatusCreateFailed {
			return resource.NonRetryableError(errCreateEIPFailed)
		}

		if *eip.AddressStatus == tencentCloudApiEipStatusUnbind {
			return nil
		}

		return resource.RetryableError(errEIPStillCreating)
	})

	return nil
}
