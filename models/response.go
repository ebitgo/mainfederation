package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// TransactionDef transaction defined
type TransactionDef struct {
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	MemoType string `json:"memo_type"`
	Memo     string `json:"memo"`
}

// ResponseBaseMsg 返回正确消息定义
type ResponseBaseMsg struct {
	Address   string `json:"stellar_address"`
	Accountid string `json:"account_id"`
	MemoType  string `json:"memo_type"`
	Memo      string `json:"memo"`
}

// ResponseMsg 返回消息定义
type ResponseMsg struct {
	Status    int
	Msg       *ResponseBaseMsg
	ErrMsg    string `json:"detail"`
	QParam    string
	TypeParam string
}

// DecodeParam 获取输入参数
func (ths *ResponseMsg) DecodeParam(vals url.Values) *ResponseMsg {
	ths.Status = 200
	ths.QParam = vals.Get("q")
	err := ths.checkValue(ths.QParam, "q")
	if err != nil {
		ths.Status = 501
		ths.ErrMsg = err.Error()
		return ths
	}
	ths.TypeParam = vals.Get("type")
	err = ths.checkValue(ths.TypeParam, "type")
	if err != nil {
		ths.Status = 501
		ths.ErrMsg = err.Error()
	}
	return ths
}

// Execute 执行
func (ths *ResponseMsg) Execute() *ResponseMsg {
	ths.TypeParam = strings.ToLower(ths.TypeParam)
	switch ths.TypeParam {
	case "name":
		ths.nameRequest()
	case "id":
		ths.idRequest()
	case "txid":
		ths.txidRequest()
	default:
		ths.Status = 501
		ths.ErrMsg = fmt.Sprintf("Undefined type [" + ths.TypeParam + "] string!")
	}
	return ths

}

func (ths *ResponseMsg) checkValue(val, flag string) (err error) {
	if len(val) == 0 {
		err = fmt.Errorf("Input parameter [%s] is not invalid or nil.", flag)
	}
	return
}

func (ths *ResponseMsg) nameRequest() {
	name := ths.getNickName(ths.QParam)
	if strings.HasSuffix(name, "*wechat") {
		ths.getFromWechat(name)
	} else {
		ths.getFromLedgercn(name)
	}
}

func (ths *ResponseMsg) getNickName(param string) string {
	name := strings.ToLower(param)
	suffixIndex := strings.LastIndex(name, "*ebitgo.com")
	if suffixIndex == -1 || suffixIndex == 0 {
		ths.Status = 501
		ths.ErrMsg = fmt.Sprintf("type=name has error [ The value of 'q' format is error ]")
		return ""
	}
	return strings.Trim(string([]byte(name)[0:suffixIndex]), " ")
}

func (ths *ResponseMsg) getFromLedgercn(name string) {
	walInfo, err := DatabaseInstance.GetWalletInfo("", name)
	ths.Status = 404
	if err == nil {
		if walInfo != nil {
			for _, wi := range walInfo {
				if strings.Compare(wi.NickName, name) == 0 {
					ths.Msg = new(ResponseBaseMsg)
					ths.Status = 200
					ths.Msg.Accountid = wi.PublicAddr
					ths.Msg.Address = name + "*ebitgo.com"
					return
				}
			}
		}
		ths.ErrMsg = fmt.Sprintf("nickname = %s is not exist", name)
	} else {
		ths.ErrMsg = err.Error()
	}
}

func (ths *ResponseMsg) getFromWechat(name string) {
	unique := strings.ToUpper(name)
	suffixIndex := strings.LastIndex(unique, "*WECHAT")
	if suffixIndex == -1 || suffixIndex == 0 {
		ths.Status = 501
		ths.ErrMsg = fmt.Sprintf("type=name has error [ The value of 'q' format is error ]")
	}
	unique = strings.Trim(string([]byte(unique)[0:suffixIndex]), " ")

	userInfo, err := DatabaseInstance.GetUserInfo("", unique)
	if err == nil {
		ths.Msg = new(ResponseBaseMsg)
		ths.Status = 200
		ths.Msg.Accountid = userInfo.PublicAddr
		ths.Msg.Address = name + "*ebitgo.com"
	} else {
		ths.Status = 404
		ths.ErrMsg = err.Error()
	}
}

func (ths *ResponseMsg) idRequest() {
	ths.getIDFromLedger(ths.QParam)
	if ths.Status != 200 {
		ths.getIDFromWechat(ths.QParam)
	}
}

func (ths *ResponseMsg) getIDFromLedger(id string) {
	walInfo, err := DatabaseInstance.GetWalletInfo(id, "")
	ths.Status = 404
	if err == nil {
		if walInfo != nil {
			for _, wi := range walInfo {
				if strings.Compare(wi.PublicAddr, id) == 0 {
					ths.Msg = new(ResponseBaseMsg)
					ths.Status = 200
					ths.Msg.Accountid = id
					ths.Msg.Address = wi.NickName + "*ebitgo.com"
					return
				}
			}
		}
		ths.ErrMsg = fmt.Sprintf("id = %s is not exist", id)
	} else {
		ths.ErrMsg = err.Error()
	}

}

func (ths *ResponseMsg) getIDFromWechat(id string) {
	userInfo, err := DatabaseInstance.GetUserInfo(id, "")
	if err == nil {
		ths.Msg = new(ResponseBaseMsg)
		ths.Status = 200
		ths.Msg.Accountid = id
		ths.Msg.Address = userInfo.UniqueId + "*wechat*ebitgo.com"
	} else {
		ths.Status = 404
		ths.ErrMsg = err.Error()
	}

}

func (ths *ResponseMsg) txidRequest() {
	txid := strings.ToLower(ths.QParam)
	quaryAddr := "https://horizon.stellar.org/transactions/" + txid
	resp, err := http.Get(quaryAddr)
	if err != nil {
		ths.Status = 404
		ths.ErrMsg = fmt.Sprintf("stellar server has error : %v", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ths.Status = 404
		ths.ErrMsg = fmt.Sprintf("http get body has error : %v", err)
		return
	}
	tran := &TransactionDef{}
	json.Unmarshal(body, tran)
	if tran.Status == 0 {
		ths.Status = 200
		ths.Msg = new(ResponseBaseMsg)
		ths.Msg.Memo = tran.Memo
		ths.Msg.MemoType = tran.MemoType
	} else {
		ths.Status = tran.Status
		ths.ErrMsg = tran.Detail
	}
}
