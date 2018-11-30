package card

import (
	"fmt"
	"statistics/calendar"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
)

type IcafeCardsInfo struct {
	Result      string `json:"result"`
	Code        int    `"json:"code""`
	Message     string `json:"message"`
	Total       int    `json:"total"`
	PageSize    int    `json:"pageSize"`
	CurrentPage int    `json:"currentPage"`
	Cards       []Card `json:"cards"`
}

type ResponsiblePeople struct {
	Email string `json:"email"`
}

type Property struct {
	LocalId      int    `json:"localId"`
	PropertyName string `json:"propertyName"`
	Value        string `json:"value"`
	DisplayValue string `json:"displayValue"`
}

type Card struct {
	CaseId       int                 `json:"sequence"`          //ID
	Status       string              `json:"status"`            //流程状态
	CaseEmail    []ResponsiblePeople `json:"responsiblePeople"` //负责人
	Title        string              `json:"title"`             //标题
	CaseProperty []Property          `json:"properties"`        //case属性
}

func (this *Card) ToEsCard(escard *EsCard) {
	escard.CaseId = fmt.Sprintf("MRDR-%d", this.CaseId)
	escard.Title = this.Title
	escard.ProgressStatus = this.Status
	escard.IcafeURL = fmt.Sprintf("%s%s%s", ICAFEURLHEAD, escard.CaseId, ICAFEURLTRAIL)
	for _, value := range this.CaseEmail {
		escard.CaseEmail = append(escard.CaseEmail, value.Email)
	}
	for _, property := range this.CaseProperty {
		switch  property.LocalId {
		case 122324:
			escard.SubModule = append(escard.SubModule, property.DisplayValue)
			break
		case 122135:
			escard.Module = property.DisplayValue
			break
		case 21666:
			escard.CaseType = append(escard.CaseType, property.DisplayValue)
			break
		case 122137:
			escard.CaseCause = property.DisplayValue
			break
		case 122128:
			escard.CaseTime = property.DisplayValue
			break
		case 123545:
			escard.CarId = property.DisplayValue
			break
		case 122131:
			escard.Taskpurpose = property.DisplayValue
			break
		case 123541:
			escard.IsoVersion = property.DisplayValue
			break
		case 123542:
			escard.MapVersion = property.DisplayValue
			break
		case 123543:
			escard.Weather = property.DisplayValue
			break
		case 122129:
			escard.InterventionSubType = property.DisplayValue
			break
		case 123500:
			escard.CaseDescription = property.DisplayValue
			break
		case 123501:
			escard.MonitorInfo = property.DisplayValue
			break
		case 123919:
			escard.MapRegion = property.DisplayValue
			break
		case 125428:
			escard.CaseDate = property.DisplayValue
			break
		case 126358:
			escard.TaskUrl = property.DisplayValue
			break
		case 126359:
			escard.DvUrl = property.DisplayValue
			escard.DvUrl = strings.Replace(escard.DvUrl, "amp;amp;", "", 5)
			break
		case 126360:
			escard.RecordUrl = property.DisplayValue
			break
		case 125430:
			escard.CaseDetailTime = property.DisplayValue
			break
		case 129222:
			escard.InterventionType = property.DisplayValue
			break
		default:
			break
		}
	}
}

type IcafeInfo struct {
	At_issueCount int
	Account       string
	Passwd        string
	Url           string
	BeginDate     string
	EndDate       string
	TotalPage     int
	CurrPage      int
	Cards         IcafeCardsInfo
}

func NewIcafeInfo(account, passwd string) *IcafeInfo {
	return &IcafeInfo{
		Account:  account,
		Passwd:   passwd,
		CurrPage: 1,
	}
}

func (this *IcafeInfo) getDate() {
	interval := calendar.GetDays()
	this.BeginDate = calendar.GetYestoday(interval)
	this.EndDate = calendar.GetYestoday(1)
}

func (this *IcafeInfo) UpdateIcafeUrl(page int) string {
	this.getDate()
	this.Url = ""
	this.Url += "http://icafe.baidu.com/api/spaces/MRDR/cards/"
	this.Url += fmt.Sprintf("?u=%s&pw=%s", this.Account, this.Passwd)
	this.Url += "&q=%5BissueTypeId%5D%5Bin%5D%5B5011%5D%2B%5B122128%5D%5Bbetween%5D%5B"
	this.Url += this.BeginDate
	this.Url += "%2000%3A00_"
	this.Url += this.EndDate
	this.Url += "%2023%3A59%5D"
	this.Url += fmt.Sprintf("&page=%d", page)
	return this.Url
}

func (this *IcafeInfo) GetIcafeReq() []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", this.Url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return body
}

func (this *IcafeInfo) UnMarshal(body []byte) {
	if this.CurrPage == 1 {
		err := json.Unmarshal(body, &this.Cards)
		if err != nil {
			fmt.Println(err.Error())
		}
		if this.Cards.Code != 200 {
			fmt.Println(this.Cards.Result)
			fmt.Println(this.Cards.Message)
			return
		}
		this.TotalPage = this.Cards.PageSize
	} else {
		var cards IcafeCardsInfo
		err := json.Unmarshal(body, &cards)
		if err != nil {
			fmt.Println(err.Error())
		}
		if cards.Code != 200 {
			fmt.Println(cards.Result)
			fmt.Println(cards.Message)
			return
		}
		this.Cards.CurrentPage = cards.CurrentPage
		this.Cards.Cards = append(this.Cards.Cards, cards.Cards...)
	}
}

func (this *IcafeInfo) Init() {
	this.UpdateIcafeUrl(1)
	this.UnMarshal(this.GetIcafeReq())
}

func (this *IcafeInfo) GetAllCase() *IcafeCardsInfo {
	this.Init()
	for i := 2; i <= this.TotalPage; i++ {
		this.CurrPage = i
		this.UpdateIcafeUrl(i)
		this.UnMarshal(this.GetIcafeReq())
	}
	return &this.Cards
}
