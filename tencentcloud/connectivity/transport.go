package connectivity

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const ReqClient = "Terraform-v1.24.1"

type LogRoundTripper struct {
}

func (me *LogRoundTripper) RoundTrip(request *http.Request) (response *http.Response, errRet error) {

	var inBytes, outBytes []byte

	var start = time.Now().UnixNano()

	defer func() { me.log(inBytes, outBytes, errRet, start) }()

	bodyReader, errRet := request.GetBody()
	if errRet != nil {
		return
	}

	request.Header.Set("X-TC-RequestClient", ReqClient)

	inBytes, errRet = ioutil.ReadAll(bodyReader)
	if errRet != nil {
		return
	}
	appendMessage := []byte(fmt.Sprintf(
		",(host %+v,action:%+v,region:%+v)",
		request.Header["Host"],
		request.Header["X-TC-Action"],
		request.Header["X-TC-Region"],
	))

	inBytes = append(inBytes, appendMessage...)

	response, errRet = http.DefaultTransport.RoundTrip(request)
	if errRet != nil {
		return
	}
	outBytes, errRet = ioutil.ReadAll(response.Body)
	if errRet != nil {
		return
	}
	response.Body = ioutil.NopCloser(bytes.NewBuffer(outBytes))
	return
}

func (me *LogRoundTripper) log(in []byte, out []byte, err error, start int64) {
	var buf bytes.Buffer
	buf.WriteString("######")
	tag := "[DEBUG]"
	if err != nil {
		tag = "[CRITAL]"
	}
	buf.WriteString(tag)
	if len(in) > 0 {
		buf.WriteString("tencentcloud-sdk-go request:")
		buf.Write(in)
	}
	if len(out) > 0 {
		buf.WriteString("; response:")
		out := bytes.Replace(out,
			[]byte("\n"),
			[]byte(""),
			-1)
		out = bytes.Replace(out,
			[]byte(" "),
			[]byte(""),
			-1)
		buf.Write(out)
	}

	if err != nil {
		buf.WriteString("; error:")
		buf.WriteString(err.Error())
	}

	costFormat := fmt.Sprintf(",cost %.3f seconds", float32(time.Now().UnixNano()-start)/float32(time.Second))
	buf.WriteString(costFormat)

	log.Println(buf.String())
}
