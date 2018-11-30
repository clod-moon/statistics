package main

import (
	"statistics/myelastic"
	"flag"
	"sync"
	"github.com/widuu/goini"
	"statistics/card"
	"statistics/syncdata"
	"fmt"
	"time"
	"strings"
	"statistics/calendar"
)

var (
	//配置文件的信息
	conf *map[string]map[string]string

	//配置文件的默认保存路径和名称
	configFile = flag.String("config", "./config/config.ini", "配置文件所在路径")

	esServerInfo *myelastic.EsServerInfo

	wg     sync.WaitGroup //定义一个同步等待的组
	wchild sync.WaitGroup //定义一个同步等待的组

	analyzer []card.Analyzer

	syncTime *syncdata.SyncTime

	currDate string

	analyzerlog *card.Analyzerlog
)

func initGlobalVar() error {

	cycle, err := syncdata.Mstrtof((*conf)["sync_time_interval"]["interval"])
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	syncTime = syncdata.NewSyncTime(
		(*conf)["sync_time_interval"]["start"],
		(*conf)["sync_time_interval"]["end"],
		cycle)

	esServerInfo = myelastic.NewEsServerInfo(
		(*conf)["esinfo"]["addr"],
		(*conf)["esinfo"]["port"],
		(*conf)["esinfo"]["index"],
		(*conf)["esinfo"]["type"],
		&wg,
	)

	for key, value := range (*conf)["analyzer"] {
		scale, err := syncdata.Mstrtof(value)
		if err != nil {
			return err
		}
		analyzer = append(analyzer, *card.NewAnalyzer(key, scale))
	}

	analyzerlog =card.NewAnalyzerlog(
		fmt.Sprintf("./log/.%s_analyzerlog",calendar.GetYestoday(calendar.GetDays())))

	analyzerlog.Analyzer = &analyzer

	esServerInfo.Init()

	currDate = time.Now().Format("20060102")
	return nil
}

func readConfig(path string) *map[string]map[string]string {
	conini := goini.SetConfig(path)
	tmpConf := conini.ReadList()
	return &tmpConf
}

func syncData() {
	icafeCards := card.NewIcafeInfo((*conf)["icafe_account"]["username"],
		(*conf)["icafe_account"]["password"])
	esCards := card.NewEsCardsInfo(&analyzer)
	syncdata := syncdata.NewSyncData(
		syncTime, icafeCards, esServerInfo, esCards,analyzerlog, &wg)
	syncdata.FuncSyncData()
}

func main() {

	conf = readConfig(*configFile)

	err := initGlobalVar()
	if err != nil {
		return
	}

	wchild.Add(1)
	go func() {
		for{
			if strings.Compare(syncTime.Next,calendar.GetCurrMinuteTime()) >= 0{
				break
			}
			syncTime.CalculateNextSyncTime()
			time.Sleep(time.Second)
		}

		for {
			newDate := time.Now().Format("20060102")
			if newDate != currDate {
				currDate = newDate
				analyzerlog.Path = fmt.Sprintf("./.%s_analyzerlog",calendar.GetYestoday(calendar.GetDays()))
			}

			if calendar.TimeEqual(syncTime.Next) &&
				time.Now().Weekday() != time.Saturday&&
				time.Now().Weekday() != time.Sunday{
				syncTime.CalculateNextSyncTime()
				syncData()
			}
			time.Sleep(time.Second * 60)
		}
		wchild.Done()
	}()

	wchild.Wait()
	return
}
