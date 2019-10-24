package tencentcloud

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

var (
	errEIPNotFound = errors.New("eip not found")
)

func dataSourceTencentCloudEip() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.20.0. Please use 'tencentcloud_eips' instead.",
		Read:               dataSourceTencentCloudEipRead,

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

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	filter := make(map[string][]string)
	filters, ok := d.GetOk("filter")
	if ok {
		for _, v := range filters.(*schema.Set).List() {
			vv := v.(map[string]interface{})
			name := vv["name"].(string)
			filter[name] = []string{}
			for _, vvv := range vv["values"].([]interface{}) {
				filter[name] = append(filter[name], vvv.(string))
			}
		}
	}

	var eips []*vpc.Address
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eips, errRet = vpcService.DescribeEipByFilter(ctx, filter)
		if errRet != nil {
			return retryError(errRet, "InternalError")
		}
		return nil
	})
	if err != nil {
		return err
	}

	includeArrears := false
	if v, ok := d.GetOk("include_arrears"); ok {
		includeArrears = v.(bool)
	}
	includeBlocked := false
	if v, ok := d.GetOk("include_blocked"); ok {
		includeBlocked = v.(bool)
	}

	if len(eips) == 0 {
		return errEIPNotFound
	}

	var filteredEips []*vpc.Address
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
