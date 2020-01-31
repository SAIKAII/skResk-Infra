package httpclient

import (
	"fmt"
	"github.com/SAIKAII/skReskInfra/lb"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/go-eureka-client/eureka"
	"github.com/tietang/props/ini"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHttpClient_Do(t *testing.T) {
	// 创建一个eureka clilent
	conf := ini.NewIniFileConfigSource("ec_test.ini")
	client := eureka.NewClient(conf)
	client.Start()
	client.Applications, _ = client.GetApplications()
	apps := &lb.Apps{Client: client}
	c := NewHttpClient(apps, &Option{Timeout: defaultHttpTimeout})
	Convey("测试HTTP客户端", t, func() {
		req, err := c.NewRequest(
			http.MethodGet,
			"http://skResk/",
			nil, nil)
		So(err, ShouldBeNil)
		So(req, ShouldNotBeNil)
		res, err := c.Do(req)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		defer res.Body.Close()
		d, err := ioutil.ReadAll(res.Body)
		So(err, ShouldBeNil)
		So(d, ShouldNotBeNil)
		fmt.Println(string(d))
	})
}
