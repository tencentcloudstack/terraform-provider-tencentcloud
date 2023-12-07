package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcSubnetCreate,
		Read:   resourceTencentCloudVpcSubnetRead,
		Update: resourceTencentCloudVpcSubnetUpdate,
		Delete: resourceTencentCloudVpcSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC to be associated.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The availability zone within which the subnet should be created.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "The name of subnet to be created.",
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
				Description:  "A network address block of the subnet.",
			},
			"is_multicast": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether multicast is enabled. The default value is 'true'.",
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of a routing table to which the subnet should be associated.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the subnet.",
			},

			// Computed values
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether it is the default VPC for this region.",
			},
			"available_ip_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of available IPs.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of subnet resource.",
			},
		},
	}
}

func resourceTencentCloudVpcSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_subnet.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		vpcId            string
		availabilityZone string
		name             string
		cidrBlock        string
		isMulticast      bool
		routeTableId     string
		tags             map[string]string
	)
	if temp, ok := d.GetOk("vpc_id"); ok {
		vpcId = temp.(string)
		if len(vpcId) < 1 {
			return fmt.Errorf("vpc_id should be not empty string")
		}
	}
	if temp, ok := d.GetOk("availability_zone"); ok {
		availabilityZone = temp.(string)
		if len(availabilityZone) < 1 {
			return fmt.Errorf("availability_zone should be not empty string")
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}
	if temp, ok := d.GetOk("cidr_block"); ok {
		cidrBlock = temp.(string)
	}

	isMulticast = d.Get("is_multicast").(bool)

	if temp, ok := d.GetOk("route_table_id"); ok {
		routeTableId = temp.(string)
		if len(routeTableId) < 1 {
			return fmt.Errorf("route_table_id should be not empty string")
		}
	}

	if routeTableId != "" {
		_, has, err := vpcService.IsRouteTableInVpc(ctx, routeTableId, vpcId)
		if err != nil {
			return err
		}
		if has != 1 {
			err = fmt.Errorf("error,route_table [%s]  not found in vpc [%s]", routeTableId, vpcId)
			log.Printf("[CRITAL]%s %s", logId, err.Error())
			return err
		}
	}

	if temp := helper.GetTags(d, "tags"); len(temp) > 0 {
		tags = temp
	}

	subnetId, err := vpcService.CreateSubnet(ctx, vpcId, name, cidrBlock, availabilityZone, tags)
	if err != nil {
		return err
	}
	d.SetId(subnetId)

	err = vpcService.ModifySubnetAttribute(ctx, subnetId, name, isMulticast)
	if err != nil {
		return err
	}

	if routeTableId != "" {
		err = vpcService.ReplaceRouteTableAssociation(ctx, subnetId, routeTableId)
		if err != nil {
			return err
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}

		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:subnet/%s", region, subnetId)

		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpcSubnetRead(d, meta)
}

func resourceTencentCloudVpcSubnetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_subnet.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region
	var (
		info VpcSubnetBasicInfo
		has  int
		e    error
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, e = vpcService.DescribeSubnet(ctx, id, nil, "", "")
		if e != nil {
			return retryError(e)
		}

		// deleted
		if has == 0 {
			d.SetId("")
			return nil
		}

		if has != 1 {
			errRet := fmt.Errorf("one subnet_id read get %d subnet info", has)
			log.Printf("[CRITAL]%s %s", logId, errRet.Error())
			return resource.NonRetryableError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "subnet", region, id)
	if err != nil {
		return err
	}

	_ = d.Set("vpc_id", info.vpcId)
	_ = d.Set("availability_zone", info.zone)
	_ = d.Set("name", info.name)
	_ = d.Set("cidr_block", info.cidr)
	_ = d.Set("is_multicast", info.isMulticast)
	_ = d.Set("route_table_id", info.routeTableId)
	_ = d.Set("is_default", info.isDefault)
	_ = d.Set("available_ip_count", info.availableIpCount)
	_ = d.Set("create_time", info.createTime)
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpcSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_subnet.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name        string
		isMulticast bool
	)
	old, now := d.GetChange("name")
	if d.HasChange("name") {
		name = now.(string)
	} else {
		name = old.(string)
	}

	old, now = d.GetChange("is_multicast")
	if d.HasChange("is_multicast") {
		isMulticast = now.(bool)
	} else {
		isMulticast = old.(bool)
	}

	d.Partial(true)

	if err := service.ModifySubnetAttribute(ctx, id, name, isMulticast); err != nil {
		return err
	}

	if d.HasChange("route_table_id") {
		routeTableId := d.Get("route_table_id").(string)
		if len(routeTableId) < 1 {
			return fmt.Errorf("route_table_id should be not empty string")
		}

		_, has, err := service.IsRouteTableInVpc(ctx, routeTableId, d.Get("vpc_id").(string))
		if err != nil {
			return err
		}
		if has != 1 {
			err = fmt.Errorf("error,route_table [%s]  not found in vpc [%s]", routeTableId, d.Get("vpc_id").(string))
			log.Printf("[CRITAL]%s %s", logId, err.Error())
			return err
		}

		if err := service.ReplaceRouteTableAssociation(ctx, id, routeTableId); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}

		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:subnet/%s", region, id)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudVpcSubnetRead(d, meta)
}

func resourceTencentCloudVpcSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_subnet.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteSubnet(ctx, d.Id()); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}
