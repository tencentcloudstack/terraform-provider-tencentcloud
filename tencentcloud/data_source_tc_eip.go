package tencentcloud

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	cvm "github.com/zqfan/tencentcloud-sdk-go/services/cvm/v20170312"
)

var (
	errEIPNotFound = errors.New("eip not found")
)

func dataSourceTencentCloudEip() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEipRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceTencentCloudFiltersSchema(),

			"include_arrears": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"include_blocked": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"id": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceTencentCloudEipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eip.read")()

	cvmConn := meta.(*TencentCloudClient).cvmConn

	req := cvm.NewDescribeAddressesRequest()
	req.Filters = []*cvm.Filter{}

	filters, filtersOk := d.GetOk("filter")
	if filtersOk {
		filterList := filters.(*schema.Set)
		req.Filters = buildFiltersParamForSDK(filterList)
	}
	req.Limit = common.IntPtr(100)
	resp, err := cvmConn.DescribeAddresses(req)
	if err != nil {
		return err
	}
	if *resp.Response.TotalCount == 0 {
		return errEIPNotFound
	}

	eips := resp.Response.AddressSet
	if len(eips) == 0 {
		return errEIPNotFound
	}

	includeArrears := false
	if v, ok := d.GetOk("include_arrears"); ok {
		includeArrears = v.(bool)
	}
	includeBlocked := false
	if v, ok := d.GetOk("include_blocked"); ok {
		includeBlocked = v.(bool)
	}

	var filteredEips []*cvm.Address
	for _, eip := range eips {
		if *eip.IsArrears && !includeArrears {
			continue
		}
		if *eip.IsBlocked && !includeBlocked {
			continue
		}
		filteredEips = append(filteredEips, eip)
	}
	if len(filteredEips) == 0 {
		return errEIPNotFound
	}

	eip := filteredEips[0]

	d.SetId(*eip.AddressId)
	d.Set("public_ip", *eip.AddressIp)
	d.Set("status", *eip.AddressStatus)
	return nil
}
