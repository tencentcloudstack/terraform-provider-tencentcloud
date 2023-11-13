/*
Provides a resource to create a cynosdb wan

Example Usage

```hcl
resource "tencentcloud_cynosdb_wan" "wan" {
  instance_grp_id = ""
}
```

Import

cynosdb wan can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_wan.wan wan_id
```
*/
package tencentcloud

import (
"context"
"fmt"
"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
"log"
"time"
)


func resourceTencentCloudCynosdbWan () * schema.Resource {
return & schema.Resource {
Create : resourceTencentCloudCynosdbWanCreate ,
Read : resourceTencentCloudCynosdbWanRead ,
Update : resourceTencentCloudCynosdbWanUpdate ,
Delete : resourceTencentCloudCynosdbWanDelete ,
Importer : & schema.ResourceImporter {
State : schema.ImportStatePassthrough ,
} ,
Schema : map[string] * schema.Schema {
"instance_grp_id": {
  Required: true,
  Type: schema.TypeString,
  Description: "ï¿½®žä¾‹ç»„id.",
},

} ,
}
} 

func resourceTencentCloudCynosdbWanCreate (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_cynosdb_wan.create") ()
defer inconsistentCheck (d,meta) () 

var clusterId string
if v,ok := d . GetOk ("cluster_id");ok {
clusterId = v .(string)
} 

d . SetId (clusterId) 

return resourceTencentCloudCynosdbWanUpdate (d,meta)
} 

func resourceTencentCloudCynosdbWanRead (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_cynosdb_wan.read") ()
defer inconsistentCheck (d,meta) () 

logId := getLogId (contextNil) 

ctx := context.WithValue (context.TODO (),logIdKey,logId) 

service := CynosdbService {client:meta .(* TencentCloudClient) . apiV3Conn} 

wanId := d . Id ()


wan,err := service . DescribeCynosdbWanById (ctx,clusterId)
if err != nil {
return err
} 

if wan == nil {
d . SetId ("")
log.Printf ("[WARN]%s resource `CynosdbWan` [%s] not found, please check if it has been deleted.\n",logId,d . Id ())
return nil
} 

if wan . InstanceGrpId != nil {
_ = d . Set ("instance_grp_id",wan . InstanceGrpId)
} 

return nil
} 

func resourceTencentCloudCynosdbWanUpdate (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_cynosdb_wan.update") ()
defer inconsistentCheck (d,meta) () 

logId := getLogId (contextNil) 

request := cynosdb.NewCloseWanRequest () 



wanId := d . Id ()


request . ClusterId = & clusterId


immutableArgs := [] string {"instance_grp_id"}


for _,v := range immutableArgs {
if d . HasChange (v) {
return fmt.Errorf ("argument `%s` cannot be changed",v)
}
}


err := resource.Retry (writeRetryTimeout,func () * resource.RetryError {
result,e := meta .(* TencentCloudClient) . apiV3Conn . UseCynosdbClient () . CloseWan (request)
if e != nil {
return  retryError (e)
} else {
log.Printf ("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",logId,request . GetAction (),request . ToJsonString (),result . ToJsonString ())
}
return  nil
})
if err != nil {
log.Printf ("[CRITAL]%s update cynosdb wan failed, reason:%+v",logId,err)
return  err
} 

service := CynosdbService {client:meta .(* TencentCloudClient) . apiV3Conn} 

conf := BuildStateChangeConf ([] string {},[] string {"success"},30 * readRetryTimeout,time.Second,service . CynosdbWanStateRefreshFunc (d . Id (),[] string {})) 

if _,e := conf . WaitForState ();e != nil {
return  e
} 

return resourceTencentCloudCynosdbWanRead (d,meta)
} 

func resourceTencentCloudCynosdbWanDelete (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_cynosdb_wan.delete") ()
defer inconsistentCheck (d,meta) () 

return nil
} 
