/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package node

//import (
//	"flag"
//	"fmt"
//	"git.internal.yunify.com/qxp/flow/internal/convert"
//	"git.internal.yunify.com/qxp/flow/internal/flow"
//	"git.internal.yunify.com/qxp/flow/pkg/config"
//	"git.internal.yunify.com/qxp/flow/pkg/utils"
//	"git.internal.yunify.com/qxp/misc/id2"
//	"io/ioutil"
//	"net/http"
//	"testing"
//)
//
//var (
//	BaseURL    = "http://127.0.0.1"
//	httpClient http.Client
//	Code       = fmt.Sprintf("flow:%s_%s", convert.CallbackOfDelay, id2.GenID())
//)
//
//func TestMain(m *testing.M) {
//	var (
//		configPath = flag.String("config", "../../configs/config.yml", "-config 配置文件地址")
//	)
//	flag.Parse()
//	err := config.Init(*configPath)
//	if err != nil {
//		panic(err)
//	}
//	m.Run()
//}
//
//// TestDelayNode 延时节点测试
//func TestDelayNode(t *testing.T) {
//	path := "/api/v1/flow/triggerFlow"
//	url := fmt.Sprintf("%s%s%s", BaseURL, config.Config.Port, path)
//	reqData := flow.FormMsg{
//		TableID: "",
//		Entity:  "",
//		Magic:   "",
//		Seq:     "",
//		Version: "0.1",
//		Method:  "POST",
//	}
//	buf, err := utils.Struct2Bytes(&reqData)
//	if err != nil {
//		t.Fatal(err)
//	}
//	resp, err := httpClient.Post(url, "application/json", buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if resp.StatusCode >= 400 {
//		t.Fatal(fmt.Sprintf("error response code: %d", resp.StatusCode))
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Printf("Response: %s", string(body))
//}
//
//// TestSaveFlowData 保存流程数据
//func TestSaveFlowData(t *testing.T) {
//
//}
//
//// TestDelayNodeCallback 延时节点dispatcher回调测试
//func TestDelayNodeCallback(t *testing.T) {
//	path := fmt.Sprintf("/api/v1/handout/%s", Code)
//	url := fmt.Sprintf("%s%s%s", BaseURL, config.Config.Port, path)
//
//	buf, err := utils.Struct2Bytes(`{}`)
//	if err != nil {
//		t.Fatal(err)
//	}
//	resp, err := httpClient.Post(url, "application/json", buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if resp.StatusCode >= 400 {
//		t.Fatal(fmt.Sprintf("error response code: %d", resp.StatusCode))
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Printf("Response: %s", string(body))
//}
//
//func TestMainOrder(t *testing.T) {
//	t.Run("TestDelayNode", TestDelayNode)
//}
