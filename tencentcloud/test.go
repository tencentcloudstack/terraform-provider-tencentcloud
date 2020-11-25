package tencentcloud

/*
import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)


import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type Client struct {
	common.Client
}

func AkSkCheck(secretId, secretKey string) (errRet error, errCode string) {

	credential := common.NewCredential(secretId, secretKey)
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.ReqMethod = "POST"
	clientProfile.HttpProfile.ReqTimeout = 300
	clientProfile.Language = "en-US"

	client := &Client{}
	client.Init("").
		WithCredential(credential).
		WithProfile(clientProfile)

	request := &tchttp.BaseRequest{}

	request.Init().WithApiInfo("sts", "2018-08-13", "GetCallerIdentity")
	request.SetDomain("sts.pre.tencentcloudapi.com")

	var response tchttp.Response
	err := client.Send(request, response)

	if err != nil {
		errRet = fmt.Errorf("error, %s", err.Error())
		return
	}
	b, _ := json.Marshal(response)
	fmt.Printf(string(b))
	return

}

func TestAkSkCheck(t *testing.T) {
	var SecretId = "AKIDEb3W0UxQysw3cMA1x6PxmHxiXPHFFGTi"
	var SecretKey = "Lg5W2Z8sYOIIJNTZgsBzx39LjAgFqWF2"
	e, c := AkSkCheck(SecretId, SecretKey)

	fmt.Println(e, c)
}
*/
