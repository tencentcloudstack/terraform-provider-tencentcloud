/*
Provide a resource to create a SnatIp of CLB instance.

~> **NOTE:** Target CLB instance must enable `snat_pro` before creating snat ips.
~> **NOTE:** Dynamic allocate IP doesn't support for now.

Example Usage

```hcl
resource "tencentcloud_clb_instance" "snat_test" {
  network_type = "OPEN"
  clb_name     = "tf-clb-snat-test"
}

resource "tencentcloud_clb_snat_ip" "foo" {
  clb_id = tencentcloud_clb_instance.snat_test.id
  ips {
  	subnet_id = "subnet-12345678"
    ip = "172.16.0.1"
  }
  ips {
  	subnet_id = "subnet-12345678"
    ip = "172.16.0.2"
  }
}

```


Import

ClbSnatIp instance can be imported by clb instance id, e.g.
```
$ terraform import tencentcloud_clb_snat_ip.test clb_id
```

*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbSnatIp() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudClbSnatIpRead,
		Create: resourceTencentCloudClbSnatIpCreate,
		Update: resourceTencentCloudClbSnatIpUpdate,
		Delete: resourceTencentCloudClbSnatIpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CLB instance ID.",
			},
			"ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Snat IP address config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Snat IP.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet ID.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClbSnatIpRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_snat_ip.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := ClbService{client}

	clbId := d.Id()

	var (
		instance *clb.LoadBalancer
		err      error
	)

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, err = service.DescribeLoadBalancerById(ctx, clbId)
		if err != nil {
			return retryError(err)
		}
		return nil
	})

	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(clbId)

	if ipLen := len(instance.SnatIps); ipLen > 0 {
		snatIps := make([]interface{}, 0)
		for i := range instance.SnatIps {
			item := instance.SnatIps[i]
			snatIps = append(snatIps, map[string]interface{}{
				"subnet_id": *item.SubnetId,
				"ip":        *item.Ip,
			})
		}
		_ = d.Set("ips", snatIps)
	}

	_ = d.Set("clb_id", clbId)

	return nil
}

func resourceTencentCloudClbSnatIpCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_snat_ip.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := ClbService{client}

	clbId := d.Get("clb_id").(string)
	ips := d.Get("ips").(*schema.Set).List()
	snatIps := make([]*clb.SnatIp, 0, len(ips))
	for i := range ips {
		item := ips[i].(map[string]interface{})
		subnetId := item["subnet_id"].(string)
		ip := item["ip"].(string)
		snatIp := &clb.SnatIp{
			SubnetId: &subnetId,
			Ip:       &ip,
		}
		snatIps = append(snatIps, snatIp)
	}

	var taskId string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		reqId, err := service.CreateLoadBalancerSnatIps(ctx, clbId, snatIps)
		if err != nil {
			return retryError(err, clb.FAILEDOPERATION)
		}
		taskId = reqId
		return nil
	})

	if err != nil {
		return err
	}

	if err := waitForTaskFinish(taskId, client.UseClbClient()); err != nil {
		return err
	}

	d.SetId(clbId)

	return resourceTencentCloudClbSnatIpRead(d, meta)
}

func resourceTencentCloudClbSnatIpUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_snat_ip.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	clbId := d.Id()
	client := meta.(*TencentCloudClient).apiV3Conn
	service := ClbService{client}

	//instanceSnatIps := make([]interface{}, 0)
	//instance, err := service.DescribeLoadBalancerById(ctx, clbId)
	//if err != nil {
	//	return err
	//}

	//if len(instance.SnatIps) > 0 {
	//	for i := range instance.SnatIps {
	//		subnet := *instance.SnatIps[i].SubnetId
	//		ip := *instance.SnatIps[i].Ip
	//		instanceSnatIps = append(instanceSnatIps, map[string]interface{}{
	//			"ip":        ip,
	//			"subnet_id": subnet,
	//		})
	//	}
	//}

	if d.HasChange("ips") {
		o, n := d.GetChange("ips")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()

		for i := range add {
			var (
				addIps []*clb.SnatIp
				taskId string
				err    error
			)
			item := add[i].(map[string]interface{})
			ip := item["ip"].(string)
			subnet := item["subnet_id"].(string)
			addIps = append(addIps, &clb.SnatIp{
				SubnetId: &subnet,
				Ip:       &ip,
			})
			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				taskId, err = service.CreateLoadBalancerSnatIps(ctx, clbId, addIps)
				if err != nil {
					return retryError(err, clb.FAILEDOPERATION)
				}
				return nil
			})
			if err != nil {
				return err
			}

			err = waitForTaskFinish(taskId, client.UseClbClient())
			if err != nil {
				return err
			}
		}

		if len(remove) > 0 {
			var (
				removeIps []*string
				taskId    string
				err       error
			)

			for i := range remove {
				item := remove[i].(map[string]interface{})
				ip := item["ip"].(string)
				removeIps = append(removeIps, &ip)
			}

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				taskId, err = service.DeleteLoadBalancerSnatIps(ctx, clbId, removeIps)
				if err != nil {
					return retryError(err, clb.FAILEDOPERATION)
				}
				return nil
			})
			if err != nil {
				return err
			}

			err = waitForTaskFinish(taskId, client.UseClbClient())
			if err != nil {
				return err
			}
		}

	}

	return resourceTencentCloudClbSnatIpRead(d, meta)
}

func resourceTencentCloudClbSnatIpDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_snat_ip.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := ClbService{client}

	clbId := d.Id()

	instanceSnatIps := make([]*string, 0)
	instance, err := service.DescribeLoadBalancerById(ctx, clbId)
	if err != nil {
		return err
	}

	if len(instance.SnatIps) > 0 {
		for i := range instance.SnatIps {
			ip := instance.SnatIps[i].Ip
			instanceSnatIps = append(instanceSnatIps, ip)
		}
	}

	var taskId string
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		reqId, err := service.DeleteLoadBalancerSnatIps(ctx, clbId, instanceSnatIps)
		if err != nil {
			return retryError(err, clb.FAILEDOPERATION)
		}
		taskId = reqId
		return nil
	})

	if err != nil {
		return err
	}

	if err := waitForTaskFinish(taskId, client.UseClbClient()); err != nil {
		return err
	}

	d.SetId("")

	return err
}
