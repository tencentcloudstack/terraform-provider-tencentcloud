/*
Provides a resource to create a mysql proxy

Example Usage

```hcl
resource "tencentcloud_mysql_proxy" "proxy" {
  instance_id    = "cdb-e4z5r8p9"
  uniq_vpc_id    = "vpc-4owdpnwr"
  uniq_subnet_id = "subnet-ahv6swf2"
  proxy_node_custom {
    node_count = 1
    cpu        = 2
    mem        = 4000
    region     = "ap-guangzhou"
    zone       = "ap-guangzhou-3"
  }
  security_group        = ["sg-edmur627"]
  desc                  = "desc"
  connection_pool_limit = 2
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlProxyCreate,
		Read:   resourceTencentCloudMysqlProxyRead,
		Update: resourceTencentCloudMysqlProxyUpdate,
		Delete: resourceTencentCloudMysqlProxyDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"uniq_vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Vpc id.",
			},

			"uniq_subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subnet id.",
			},

			"proxy_node_custom": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Node specification configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of nodes.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of CPU cores.",
						},
						"mem": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Memory size.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region.",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Zone.",
						},
					},
				},
			},

			"security_group": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group.",
			},

			"desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Describe.",
			},

			"connection_pool_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Connection Pool Threshold.",
			},

			// Computed

		},
	}
}

func resourceTencentCloudMysqlProxyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_proxy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = mysql.NewCreateCdbProxyRequest()
		response   = mysql.NewCreateCdbProxyResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_subnet_id"); ok {
		request.UniqSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_node_custom"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			proxyNodeCustom := mysql.ProxyNodeCustom{}
			if v, ok := dMap["node_count"]; ok {
				proxyNodeCustom.NodeCount = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["cpu"]; ok {
				proxyNodeCustom.Cpu = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["mem"]; ok {
				proxyNodeCustom.Mem = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["region"]; ok {
				proxyNodeCustom.Region = helper.String(v.(string))
			}
			if v, ok := dMap["zone"]; ok {
				proxyNodeCustom.Zone = helper.String(v.(string))
			}
			request.ProxyNodeCustom = append(request.ProxyNodeCustom, &proxyNodeCustom)
		}
	}

	if v, ok := d.GetOk("security_group"); ok {
		securityGroupSet := v.(*schema.Set).List()
		for i := range securityGroupSet {
			securityGroup := securityGroupSet[i].(string)
			request.SecurityGroup = append(request.SecurityGroup, &securityGroup)
		}
	}

	if v, ok := d.GetOk("desc"); ok {
		request.Desc = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("connection_pool_limit"); ok {
		request.ConnectionPoolLimit = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().CreateCdbProxy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql proxy failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s create mysql proxy status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s create mysql proxy status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql proxy fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlProxyRead(d, meta)
}

func resourceTencentCloudMysqlProxyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_proxy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()
	proxy, err := service.DescribeMysqlProxyById(ctx, instanceId)
	if err != nil {
		return err
	}

	if proxy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlProxy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	// if proxy.ProxyAddress.UniqVpcId != nil {
	// 	_ = d.Set("uniq_vpc_id", proxy.UniqVpcId)
	// }

	// if proxy.UniqSubnetId != nil {
	// 	_ = d.Set("uniq_subnet_id", proxy.UniqSubnetId)
	// }

	// if proxy.ProxyNodeCustom != nil {
	// 	proxyNodeCustomList := []interface{}{}
	// 	for _, proxyNodeCustom := range proxy.ProxyNodeCustom {
	// 		proxyNodeCustomMap := map[string]interface{}{}

	// 		if proxy.ProxyNodeCustom.NodeCount != nil {
	// 			proxyNodeCustomMap["node_count"] = proxy.ProxyNodeCustom.NodeCount
	// 		}

	// 		if proxy.ProxyNodeCustom.Cpu != nil {
	// 			proxyNodeCustomMap["cpu"] = proxy.ProxyNodeCustom.Cpu
	// 		}

	// 		if proxy.ProxyNodeCustom.Mem != nil {
	// 			proxyNodeCustomMap["mem"] = proxy.ProxyNodeCustom.Mem
	// 		}

	// 		if proxy.ProxyNodeCustom.Region != nil {
	// 			proxyNodeCustomMap["region"] = proxy.ProxyNodeCustom.Region
	// 		}

	// 		if proxy.ProxyNodeCustom.Zone != nil {
	// 			proxyNodeCustomMap["zone"] = proxy.ProxyNodeCustom.Zone
	// 		}

	// 		proxyNodeCustomList = append(proxyNodeCustomList, proxyNodeCustomMap)
	// 	}

	// 	_ = d.Set("proxy_node_custom", proxyNodeCustomList)

	// }

	// if proxy.SecurityGroup != nil {
	// 	_ = d.Set("security_group", proxy.SecurityGroup)
	// }

	// if proxy.Desc != nil {
	// 	_ = d.Set("desc", proxy.Desc)
	// }

	if proxy.ConnectionPoolLimit != nil {
		_ = d.Set("connection_pool_limit", proxy.ConnectionPoolLimit)
	}

	return nil
}

func resourceTencentCloudMysqlProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_proxy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	proxy, err := service.DescribeMysqlProxyById(ctx, instanceId)
	if err != nil {
		return err
	}

	if proxy == nil || proxy.ProxyGroupId == nil {
		return fmt.Errorf("Instance `%s` proxy does not exist", instanceId)
	}

	immutableArgs := []string{"instance_id", "uniq_vpc_id", "uniq_subnet_id", "security_group", "desc"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("proxy_node_custom") {
		request := mysql.NewAdjustCdbProxyRequest()
		response := mysql.NewAdjustCdbProxyResponse()

		request.InstanceId = &instanceId
		request.ProxyGroupId = proxy.ProxyGroupId

		if v, ok := d.GetOk("proxy_node_custom"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				proxyNodeCustom := mysql.ProxyNodeCustom{}
				if v, ok := dMap["node_count"]; ok {
					proxyNodeCustom.NodeCount = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["cpu"]; ok {
					proxyNodeCustom.Cpu = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["mem"]; ok {
					proxyNodeCustom.Mem = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["region"]; ok {
					proxyNodeCustom.Region = helper.String(v.(string))
				}
				if v, ok := dMap["zone"]; ok {
					proxyNodeCustom.Zone = helper.String(v.(string))
				}
				request.ProxyNodeCustom = append(request.ProxyNodeCustom, &proxyNodeCustom)
			}
		}

		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().AdjustCdbProxy(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mysql proxy failed, reason:%+v", logId, err)
			return err
		}

		asyncRequestId := *response.Response.AsyncRequestId
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("%s update mysql proxy status is %s", instanceId, taskStatus))
			}
			err = fmt.Errorf("%s update mysql proxy status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s update mysql proxy fail, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	if d.HasChange("connection_pool_limit") {
		connectionPoolLimit := d.Get("connection_pool_limit")
		request := mysql.NewModifyCdbProxyParamRequest()
		request.InstanceId = &instanceId
		request.ProxyGroupId = proxy.ProxyGroupId
		request.ConnectionPoolLimit = helper.IntUint64(connectionPoolLimit.(int))

		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ModifyCdbProxyParam(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mysql proxy failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMysqlProxyRead(d, meta)
}

func resourceTencentCloudMysqlProxyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_proxy.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	if err := service.DeleteMysqlProxyById(ctx, instanceId); err != nil {
		return err
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		proxy, err := service.DescribeMysqlProxyById(ctx, instanceId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if proxy == nil {
			return nil
		}
		if proxy != nil {
			return resource.RetryableError(fmt.Errorf("%s delete mysql proxy status is %s", instanceId, *proxy.Status))
		}
		err = fmt.Errorf("%s delete mysql proxy status is %s,we won't wait for it finish", instanceId, *proxy.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete mysql proxy fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
