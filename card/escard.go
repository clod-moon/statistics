package card

import (
	"fmt"
	"math"
)

const (
	ICAFEURLHEAD  = "http://newicafe.baidu.com/issue/"
	ICAFEURLTRAIL = "/show?from=page"
)

type Analyzer struct {
	Name  string   `json:"name"`
	Scale float64  `json:"scale"`
	Cases []string `json:"cases"`
}


func NewAnalyzer(name string, scale float64) *Analyzer {
	return &Analyzer{Name: name, Scale: scale}
}

type EsCardsInfo struct {
	Cards    []*EsCard
	Analyzer *[]Analyzer
}


func (this *EsCardsInfo) SortByScale() {
	for n := 0; n <= len(*this.Analyzer); n++ {
		for i := 1; i < len(*this.Analyzer)-n; i++ {
			if (*this.Analyzer)[i].Scale < (*this.Analyzer)[i-1].Scale {
				(*this.Analyzer)[i], (*this.Analyzer)[i-1] = (*this.Analyzer)[i-1], (*this.Analyzer)[i]
			}
		}
	}
}

func NewEsCardsInfo(Analyzer *[]Analyzer) *EsCardsInfo {
	return &EsCardsInfo{Analyzer: Analyzer}
}

func (this *EsCardsInfo) Allocation(analyzerlog *Analyzerlog) {
	var totalScale float64
	this.SortByScale()
	for _, Analyzer := range (*this.Analyzer) {
		totalScale += Analyzer.Scale
	}
	//fmt.Println("len(this.Cards)",(this.Cards))
	eachNumber := float64(len(this.Cards)+len(analyzerlog.Classified)) / totalScale
	total := 0
	for i:=0;i<len(*analyzerlog.Analyzer);i++{
	//	fmt.Println("cases:",len((*analyzerlog.Analyzer)[i].Cases))
		analyzerTotal := int(math.Floor((*analyzerlog.Analyzer)[i].Scale * eachNumber))
		classified :=len((*analyzerlog.Analyzer)[i].Cases)
		add := 0
		if analyzerTotal>classified{
			add = analyzerTotal - classified
		}else{
			continue
		}
		tmpTotal := total
		fmt.Println("n:",add)
		total += add
		for i := tmpTotal; i < total; i++ {
			this.Cards[i].Analyzer = (*analyzerlog.Analyzer)[i].Name
			analyzerlog.Classified[this.Cards[i].CaseId]=Empty{}
			(*analyzerlog.Analyzer)[i].Cases = append((*analyzerlog.Analyzer)[i].Cases,this.Cards[i].CaseId)
		}
	}
	remaining := len(this.Cards)-total
	//fmt.Println("remaining:",remaining)
	//fmt.Println("len(*this.Analyzer):",len(*this.Analyzer))
	//fmt.Println("this.Cards:",len(this.Cards))
	//fmt.Println("total:",total)
	for i := len(*this.Analyzer) - remaining; i < len(*this.Analyzer); i++ {
		this.Cards[total].Analyzer = (*this.Analyzer)[i].Name
		analyzerlog.Classified[this.Cards[total].CaseId]=Empty{}
		//fmt.Println(len((*analyzerlog.Analyzer)[i].Cases))
		(*analyzerlog.Analyzer)[i].Cases = append((*analyzerlog.Analyzer)[i].Cases,this.Cards[total].CaseId)
		//fmt.Println(len((*analyzerlog.Analyzer)[i].Cases))
		total++
	}
	analyzerlog.UpdateClassified()
}

type EsCard struct {
	CaseId              string   `json:"case_id"`               //caseID
	ProgressStatus      string   `json:"progress_status"`       //流程状态
	CaseEmail           []string `json:"case_email"`            //负责人
	Title               string   `json:"title"`                 //标题
	IcafeURL            string   `json:"icafe_url"`             //case ifcafe_url
	Analyzer            string   `json:"analyzer"`              //case分析人负责人
	Weather             string   `json:"weather"`               //天气
	Taskpurpose         string   `json:"taskpurpose"`           //路跑or路测
	InterventionType    string   `json:"intervention_type"`     //接管类型
	InterventionSubType string   `json:"intervention_sub_type"` //接管子类型
	CaseDescription     string   `json:"case_desp"`             //case描述
	MonitorInfo         string   `json:"monitor_info"`          //监控异常信息
	TaskUrl             string   `json:"task_url"`              //case task_url
	DvUrl               string   `json:"dv_url"`                //case dv_url
	RecordUrl           string   `json:"record_url"`            //case record_url
	CaseDate            string   `json:"case_date"`             //问题发生日期
	CaseTime            string   `json:"case_time"`             //case发生时间
	CaseDetailTime      string   `json:"case_detail_time"`      //case精确发生时间
	IsoVersion          string   `json:"iso_version"`           //ISO版本
	Module              string   `json:"module"`                //case模块
	SubModule           []string `json:"sub_module"`            //case子模块
	CaseType            []string `json:"case_type"`             //case归类
	CaseCause           string   `json:"case_cause"`            //case原因
	CarId               string   `json:"car_id"`                //车辆ID
	MapVersion          string   `json:"map_version"`           //地图版本号
	MapRegion           string   `json:"map_region"`            //地图区域名称
}

func NewEsCard() *EsCard {
	return &EsCard{}
}

func (this *EsCard) MakceId() string {
	return fmt.Sprintf("%s_%s", this.CaseDate, this.CaseId)
}
