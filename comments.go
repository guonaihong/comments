package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/guonaihong/gout"
	"os"
	"regexp"
)

//http://fanyi.youdao.com/translate?&doctype=json&type=AUTO&i=计算
// {"type":"ZH_CN2EN","errorCode":0,"elapsedTime":0,"translateResult":[[{"src":"计算","tgt":"To calculate"}]]}
type youdao struct {
	Doctype string `query:"doctype"`
	Type    string `query:"auto"`
	I       string `query:"i"`
}

type line struct {
	Src string `json:"src"`
	Tgt string `json:"tgt"`
}

type result struct {
	ErrorCode       int      `json:"errorcode"`
	TranslateResult [][]line `json:"translateResult"`
}

func (r *result) getSentence() string {
	if len(r.TranslateResult) == 0 {
		return ""
	}

	if len(r.TranslateResult[0]) == 0 {
		return ""
	}

	return r.TranslateResult[0][0].Tgt
}

var defYoudao = youdao{Doctype: "json", Type: "auto"}

func getEnglish(s string) string {
	if len(s) == 0 {
		return s
	}

	d := defYoudao
	d.I = s

	r := result{}
	err := gout.
		GET("http://fanyi.youdao.com/translate").
		Debug(true).
		SetQuery(d).
		BindJSON(&r).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	rv := r.getSentence()
	if len(rv) == 0 {
		return s
	}

	return rv
}

func translate(inFile, outFile string, debug bool) {
	inFd, err := os.Open(inFile)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer inFd.Close()

	outFd, err := os.OpenFile(outFile, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0644)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer outFd.Close()

	br := bufio.NewReader(inFd)

	re := regexp.MustCompile("[\u4e00-\u9fa5]*")

	for {

		l, e := br.ReadBytes('\n')
		if len(l) == 0 && e != nil {
			break
		}

		rv := re.ReplaceAllStringFunc(string(l), getEnglish)

		outFd.WriteString(rv)
	}
}

func main() {
	in := flag.String("in", "", "(must)input file")
	out := flag.String("out", "", "(must)output file")
	debug := flag.Bool("debug", false, "debug mode")

	flag.Parse()

	if len(*in) == 0 || len(*out) == 0 {
		flag.Usage()
		return
	}

	translate(*in, *out, *debug)
}
