/*
Provides a resource to create a mysql proxy

Example Usage

```hcl
resource "tencentcloud_mysql_proxy" "proxy" {
  instance_id = ""
  uniq_vpc_id = ""
  uniq_subnet_id = ""
  proxy_node_custom {
	node_count =
	cpu =
	mem =
	region = ""
	zone = ""
  }
  security_group =
  desc = ""
  connection_pool_limit =
}
```

Import

mysql proxy can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_proxy.proxy proxy_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance Id.",
			},

			"uniq_vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "private network ID.",
			},

			"uniq_subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "private subnet ID.",
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
							Description: "number of CPU cores.",
						},
						"mem": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "memory size.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "region.",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "availability zone.",
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
				Description: "security group.",
			},

			"desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "describe.",
			},

			"connection_pool_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Connection Pool Threshold.",
			},
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
			return resource.RetryableError(fmt.Errorf("%s create mysql rollback status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s create mysql rollback status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql rollback fail, reason:%s\n ", logId, err.Error())
		return err
	}

	proxy, err := service.DescribeMysqlProxyById(ctx, instanceId)
	if err != nil {
		return err
	}

	d.SetId(instanceId + FILED_SP + *proxy.ProxyGroupId)

	return resourceTencentCloudMysqlProxyRead(d, meta)
}

func resourceTencentCloudMysqlProxyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_proxy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	// proxyGroupId := idSplit[1]

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

	// proxyAddress := proxy.ProxyAddress[0]
	// if proxyAddress != nil {
	// 	if proxyAddress.UniqVpcId != nil {
	// 		_ = d.Set("uniq_vpc_id", proxyAddress.UniqVpcId)
	// 	}

	// 	if proxyAddress.UniqSubnetId != nil {
	// 		_ = d.Set("uniq_subnet_id", proxyAddress.UniqSubnetId)
	// 	}

	// 	if proxyAddress.Desc != nil {
	// 		_ = d.Set("desc", proxyAddress.Desc)
	// 	}
	// }

	if proxy.ProxyNode != nil {
		proxyNodeCustomList := []interface{}{}
		for _, proxyNodeCustom := range proxy.ProxyNode {
			proxyNodeCustomMap := map[string]interface{}{}

			proxyNodeCustomMap["node_count"] = len(proxy.ProxyNode)

			if proxyNodeCustom.Cpu != nil {
				proxyNodeCustomMap["cpu"] = proxyNodeCustom.Cpu
			}

			if proxyNodeCustom.Mem != nil {
				proxyNodeCustomMap["mem"] = proxyNodeCustom.Mem
			}

			if proxyNodeCustom.Region != nil {
				proxyNodeCustomMap["region"] = proxyNodeCustom.Region
			}

			if proxyNodeCustom.Zone != nil {
				proxyNodeCustomMap["zone"] = proxyNodeCustom.Zone
			}

			proxyNodeCustomList = append(proxyNodeCustomList, proxyNodeCustomMap)
		}

		_ = d.Set("proxy_node_custom", proxyNodeCustomList)

	}

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

	var asyncRequestId string
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	proxyGroupId := idSplit[1]

	immutableArgs := []string{"instance_id", "uniq_vpc_id", "uniq_subnet_id", "security_group"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("proxy_node_custom") {
		request := mysql.NewAdjustCdbProxyRequest()
		response := mysql.NewAdjustCdbProxyResponse()
		request.InstanceId = &instanceId
		request.ProxyGroupId = &proxyGroupId
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
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().AdjustCdbProxy(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			response = result
			asyncRequestId = *response.Response.AsyncRequestId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mysql proxy failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("connection_pool_limit") {
		request := mysql.NewModifyCdbProxyParamRequest()
		request.InstanceId = &instanceId
		request.ProxyGroupId = &proxyGroupId
		if v, ok := d.GetOk("open_connection_pool"); ok {
			request.ConnectionPoolLimit = helper.IntUint64(v.(int))
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
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

	if asyncRequestId != "" {
		service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("%s create mysql rollback status is %s", instanceId, taskStatus))
			}
			err = fmt.Errorf("%s create mysql rollback status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s create mysql rollback fail, reason:%s\n ", logId, err.Error())
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

	return nil
}
