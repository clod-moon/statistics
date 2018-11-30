package card

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
)

type Empty struct{}

type Analyzerlog struct {
	Path       string
	Analyzer   *[]Analyzer `json:"analyzer"`
	Classified map[string]Empty
}

type tmpAnalyzerlog struct {
	Path       string
	Analyzer   []Analyzer `json:"analyzer"`
	Classified map[string]Empty
}
func newTmpAnalyzerlog(path string) *Analyzerlog{
	return &Analyzerlog{
		Path:path,
		Classified:make(map[string]Empty),
	}
}

func NewAnalyzerlog(path string) *Analyzerlog{
	return &Analyzerlog{
		Path:path,
		Classified:make(map[string]Empty),
	}
}

func (this *Analyzerlog) GetClassified() error {
	b, err := ioutil.ReadFile(this.Path)
	if err != nil {
		fmt.Print(err)
		return err
	}

	return json.Unmarshal(b,this)
}


func (this *Analyzerlog) UpdateClassified() error {
	log,err := json.Marshal(this)
	if err != nil{
		fmt.Println(err)
		return err
	}
	return ioutil.WriteFile(this.Path, log, 0666)
}

func (this *Analyzerlog) Remove(){

}
