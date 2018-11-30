package syncdata

import (
	"statistics/card"
	"strconv"
	"strings"
	"sync"
	"statistics/myelastic"
	"math"
	"fmt"
)

func Mstrtof(num string) (float64, error) {

	arr := strings.Split(num, ".")
	integerPart, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0.0, err
	}
	decimalPart := 0
	if len(arr) > 1 {
		decimalPart, err = strconv.Atoi(arr[1])
		if err != nil {
			return 0.0, err
		}
	}

	return float64(integerPart) + float64(decimalPart)/10, nil
}

type SyncTime struct {
	Begin    string
	End      string
	Cycle    float64
	Next     string
	IsUpdate bool
}

func NewSyncTime(begin, end string, cycle float64) *SyncTime {
	return &SyncTime{
		Begin:    begin,
		End:      end,
		Cycle:    cycle,
		Next:     begin,
		IsUpdate: true,
	}
}

func (this *SyncTime) CalculateNextSyncTime(){
	if this.Next > this.End {
		this.Next = this.Begin
		return
	}
	tmpNext := strings.Split(this.Next, ":")
	hour, err := strconv.Atoi(tmpNext[0])
	if err != nil {
		return
	}
	minute, err := strconv.Atoi(tmpNext[1])
	if err != nil {
		return
	}

	h := int(math.Floor(this.Cycle))
	minute += int((this.Cycle - float64(h)) * 60)
	hour += int(math.Floor(this.Cycle))
	hour += minute / 60
	minute = minute % 60
	if minute <10{
		this.Next = fmt.Sprintf("%d:0%d", hour, minute)
	}else{
		this.Next = fmt.Sprintf("%d:%d", hour, minute)
	}
}

type SyncData struct {
	SyncTime    *SyncTime
	IcafeCards  *card.IcafeInfo
	EsCards     *card.EsCardsInfo
	EsServer    *myelastic.EsServerInfo
	Analyzerlog *card.Analyzerlog
	Wg          *sync.WaitGroup
}

func NewSyncData(syncTime *SyncTime,
	icafeCards *card.IcafeInfo,
	esServer *myelastic.EsServerInfo,
	esCards *card.EsCardsInfo,
	analyzerlog	*card.Analyzerlog,
	wg *sync.WaitGroup) *SyncData {
	return &SyncData{SyncTime: syncTime,
		IcafeCards: icafeCards,
		EsServer: esServer,
		EsCards: esCards,
		Analyzerlog:analyzerlog,
		Wg: wg,
	}
}

func (this *SyncData) FuncSyncData() {
	this.IcafeCards.GetAllCase()
	this.Analyzerlog.GetClassified()
	//fmt.Println(this.Analyzerlog)
	for _, cards := range this.IcafeCards.Cards.Cards {
		escard := card.NewEsCard()
		cards.ToEsCard(escard)
		if escard.Module == "PNC" &&
			len(escard.CaseType[0]) == 0 &&
			len(escard.Analyzer) == 0 {
				_,ok := this.Analyzerlog.Classified[escard.CaseId]
				if ok{
					//fmt.Println("ok")
					continue
				}else{
					//fmt.Println("not ok:",escard.CaseId)
					this.EsCards.Cards = append(this.EsCards.Cards, escard)
					continue
				}
		}
		this.Wg.Add(1)
		if this.EsServer.QueryById(escard) {
			go this.EsServer.Update(escard)
		} else {
			go this.EsServer.Insert(escard)
		}
	}

	if len(this.EsCards.Cards) == 0{
		return
	}
	this.EsCards.Allocation(this.Analyzerlog)
	for _, c := range this.EsCards.Cards {
		this.Wg.Add(1)
		if this.EsServer.QueryById(c) {
			go this.EsServer.Update(c)
		} else {
			go this.EsServer.Insert(c)
		}
	}
	this.Wg.Wait()
}
