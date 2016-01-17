package main

import (
	"encoding/json"
	"fmt"
	goquery "github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func checkError(err error) {
	if err != nil {
		panic(err)
		os.Exit(1)
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkError(err)
	return dir
}

type Currency struct {
	Name   string
	Source string
}

type Email struct {
	Addr string
}

type Config struct {
	Currencys    []Currency
	Emails       []Email
	SMTPServer   string
	SMTPMail     string
	SMTPPassword string
}

func main() {
	mailBody := time.Now().Format("060102")
	var cfg Config
	var err error
	//fmt.Println(getCurrentDirectory())
	cfgFile, err := os.Open(getCurrentDirectory() + `/config.json`)
	checkError(err)
	defer cfgFile.Close()
	cfgJSON, _ := ioutil.ReadAll(cfgFile)
	json.Unmarshal(cfgJSON, &cfg)

	fmt.Println("\n ** Currencys:")
	for _, currency := range cfg.Currencys {
		fmt.Println(currency.Name + " | " + currency.Source)
		docEachCurrency, err := goquery.NewDocument(currency.Source)
		checkError(err)
		selEachCurrency := docEachCurrency.Find(`div.rate`).First()
		//fmt.Println(selEachCurrency.Text())
		eachCurrency := regexp.MustCompile(`\s`).
			ReplaceAllString(selEachCurrency.Text(), ``)
		if matched, _ := regexp.MatchString(`\d+\.\d+\/\d+\.\d+.*`, eachCurrency); matched {
			fmt.Println(`    ** match a format`)
			fmt.Println(eachCurrency)
			eachCurrency = regexp.MustCompile(`.*\/(\d+\.\d+).*`).
				ReplaceAllString(eachCurrency, `$1`)
			fmt.Println(eachCurrency)
		} else {
			fmt.Println(`    ** not match the format`)
			fmt.Println(eachCurrency)
			eachCurrency = regexp.MustCompile(`.*(\d+\.\d+).*`).
				ReplaceAllString(eachCurrency, `$1`)
			fmt.Println(eachCurrency)
			f, err := strconv.ParseFloat(eachCurrency, 32)
			checkError(err)
			eachCurrency = strconv.FormatFloat(1/f, 'f', 4, 32)
		}
		eachCurrency = currency.Name + ": " + eachCurrency
		mailBody += "\n" + eachCurrency
	}
	fmt.Println("\n ** mailBody:\n" + mailBody)

	fmt.Println("\n ** Emails:")
	for _, email := range cfg.Emails {
		fmt.Println(email.Addr)
		fmt.Println(" *** send email")
		for i := 0; i < 10; i++ {
			err = sendToMail(
				cfg.SMTPMail,     /*fromMail*/
				"银联汇率",           /*fromName*/
				cfg.SMTPPassword, /*password*/
				cfg.SMTPServer,   /*smtpServer*/
				email.Addr,
				"",     /*toName*/
				"常见币种", /*subject*/
				mailBody)
			if err == nil {
				fmt.Println("Send mail success!")
				break
			} else {
				fmt.Printf("Send mail fail! Retry %d\n", i)
			}
		}
		checkError(err)
	}
}
