package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	cfg * TomlConfig
	once sync.Once
)

var indexPage = "https://www.feastogether.com.tw/booking/"
var loginUrl = "https://www.feastogether.com.tw/memberAPI/login"
var otGetPossibleUrl = "https://www.feastogether.com.tw/orderAPI/otGetPossible"
var bookingCheckUrl = "https://www.feastogether.com.tw/booking-check"
var orderSetUrl = "https://www.feastogether.com.tw/orderAPI/orderSet"

func otGetPossible(httpClient *http.Client, cookie string, csrfToken string)  {

	cityName := Config().City
	cityCode := "1"
	mealName := Config().MealTime
	mealCode := "lunch"

	if cityName == "台北市" {
		cityCode = "1"
	} else if cityName == "新北市" {
		cityCode = "2"
	} else if cityName == "桃園市" {
		cityCode = "3"
	} else if cityName == "台中市" {
		cityCode = "4"
	} else if cityName == "台南市" {
		cityCode = "5"
	} else if cityName == "高雄市" {
		cityCode = "6"
	} else if cityName == "新竹市" {
		cityCode = "8"
	}


	if mealName == "午餐" {
		mealCode = "lunch"
	} else if mealName == "下午餐" {
		mealCode = "afternoon-tea"
	} else if mealName == "晚餐" {
		mealCode = "dinner"
	}

	form := url.Values{}
	form.Add("bu_code", "res" + Config().Res)
	form.Add("city", cityCode)
	form.Add("people", Config().People)
	form.Add("date", Config().BookingDate)
	form.Add("meal_time", mealCode)

	req, err := http.NewRequest("POST", otGetPossibleUrl, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("origin", "https://www.feastogether.com.tw")
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("cookie", cookie)
	req.Header.Set("referer", indexPage)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.104 Safari/537.36")
	req.Header.Set("x-csrf-token", csrfToken)
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	possibleResp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer possibleResp.Body.Close()
	possibleResultStr , err := ioutil.ReadAll(possibleResp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

//	 fmt.Printf("%s\n", possibleResultStr)

	if (string(possibleResultStr)) == "\"\""  {
		fmt.Println("訂位日期尚未開始，無法繼續")
		fmt.Println("訂位失敗")
		os.Exit(0)
	}

	type seatInfo map[string][]PossibleSeat
	var info seatInfo
	json.Unmarshal(possibleResultStr, &info)

	/*fmt.Println("info=======>")
	fmt.Println(info)*/

	for _, s := range info["res" + Config().Res] {
		for _,c := range s.Content {

			if c.Store == Config().Store {
				fmt.Println("店名：", c.Store, " , 剩餘座位：", c.Data.Seat)

				people,_ := strconv.Atoi(Config().People)
				if c.Data.Seat < people {
					fmt.Println("剩餘座位不足，訂位失敗")
					os.Exit(0)
				}
				break
			}
		}
	}

}

func bookingCheck(httpClient *http.Client, cookie string, csrfToken string)  {
	fmt.Print("確認訂位資訊...")

	form := url.Values{}
	form.Add("_token", csrfToken)
	form.Add("townShip", "台北市")
	form.Add("peopleNums", "2")
	form.Add("date", "2021-11-30")
	form.Add("mealTime", "晚餐")
	form.Add("time", "<i class=\"fa fa-angle-right\" aria-hidden=\"true\" data-header-member-menu=\"\"></i>")
	form.Add("store", "京站店")
	form.Add("timeMobile", "")
	form.Add("storeMobile", "京站店")

	req, err := http.NewRequest("POST", bookingCheckUrl, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("origin", "https://www.feastogether.com.tw")
	//req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9\n")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("cookie", cookie)
	req.Header.Set("referer", indexPage)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.104 Safari/537.36")
	req.Header.Set("upgrade-insecure-requests", "1")
	/*req.Header.Set("x-csrf-token", csrfToken)
	req.Header.Set("x-requested-with", "XMLHttpRequest")*/

	bookingResp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer bookingResp.Body.Close()

	fmt.Println("完成")
	/*resultStr , err := ioutil.ReadAll(bookingResp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	 fmt.Printf("%s\n", resultStr)*/
}

func orderSet(httpClient *http.Client, cookie string, csrfToken string)  {
	form := url.Values{}
	form.Add("eat_vegetable", Config().Vegetable + " 位")
	form.Add("child_chair", Config().ChildChair + " 張")
	form.Add("special_needs", "")
	form.Add("townShip", Config().City)
	form.Add("peopleNums",  Config().People)
	form.Add("date",  Config().BookingDate)
	form.Add("mealTime",  Config().MealTime)
	form.Add("store",  Config().Store)
	form.Add("order_time",  Config().Time)

	req, err := http.NewRequest("POST", orderSetUrl, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("origin", "https://www.feastogether.com.tw")
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("cookie", cookie)
	req.Header.Set("referer", bookingCheckUrl)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.104 Safari/537.36")
	req.Header.Set("x-csrf-token", csrfToken)
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	bookingResp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer bookingResp.Body.Close()
	bookingResultStr , err := ioutil.ReadAll(bookingResp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%s\n", bookingResultStr)

	var bookingResult BookingResult
	json.Unmarshal(bookingResultStr, &bookingResult)

	if len(bookingResult.OrderNo) > 0 &&  bookingResult.State == "finish" {
		fmt.Println("恭喜訂位成功!!!!!!!!!!!!")
		fmt.Println("訂位編號：", bookingResult.OrderNo)
	} else {
		fmt.Println("很可惜，訂位失敗，請再接再勵!")
	}
}


func Config() *TomlConfig {
	once.Do(func() {
		filePath, err := filepath.Abs("config.toml")
		if err != nil {
			fmt.Println("找不到設定檔，無法繼續")
			os.Exit(0)
		}
		//fmt.Printf("parse toml file once. filePath: %s\n", filePath)
		if _ , err := toml.DecodeFile(filePath, &cfg); err != nil {
			fmt.Println("解析設定檔發生錯誤，無法繼續")
			os.Exit(0)
		}
	})
	return cfg
}

func main() {

	resName := ""

	if Config().Res == "1" {
		resName = "響食天堂"
	} else if Config().Res == "2" {
		resName = "饗饗"
	} else if Config().Res == "10" {
		resName = "旭集"
	}

	fmt.Println("=====設定參數=====")
	fmt.Println("會員帳號:", Config().Account)
	fmt.Println("會員密碼:", Config().Password)
	fmt.Println("訂位餐廳:", resName)
	fmt.Println("訂位區域:", Config().City)
	fmt.Println("訂位人數:", Config().People)
	fmt.Println("訂位日期:", Config().BookingDate)
	fmt.Println("訂位餐次:", Config().MealTime)
	fmt.Println("訂位店別:", Config().Store)
	fmt.Println("進場時段:", Config().Time)
	fmt.Println("素食人數:", Config().Vegetable)
	fmt.Println("兒童座椅:", Config().ChildChair)


	fmt.Println("=====開始訂位=====")

	csrfToken := ""

	resp, err := http.Get(indexPage + Config().Res)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		op, _ := s.Attr("name")
		con, _ := s.Attr("content")
		if op == "csrf-token" {
			// fmt.Println("csrf-token=", con)
			csrfToken = con
		}

	})

	var cookieStr string

	for name, values := range resp.Header {
		for _, value := range values {

			if strings.Contains(name,"Set-Cookie") {
				// res = append(res, fmt.Sprintf("%s: %s", name, value))
				cookieStr = cookieStr + value
			}
		}
	}

	// fmt.Println("CookieStr=====>", cookieStr)

	xsrfTokenFilter, _ := regexp.Compile("^(XSRF-TOKEN=.*?;)")
	laravelSessionFilter, _ := regexp.Compile("securelaravel_session=.+?;")

	xsrfToker := strings.ReplaceAll(xsrfTokenFilter.FindString(cookieStr),"XSRF-TOKEN=","")
	session := strings.ReplaceAll(laravelSessionFilter.FindString(cookieStr),"securelaravel_session=","laravel_session=")
	cookie := session + xsrfToker //+ " _ga=GA1.3.501333287.1637306467; _gid=GA1.3.370617231.1637306467; _gat_gtag_UA_109250939_1=1;"


	/*fmt.Println(xsrfToker)
	fmt.Println(session)*/



	fmt.Print("登入會員中...")
	client := &http.Client{}


	var jsonStr = []byte(`{
  "login_type": 1,
  "account": "0966961001",
  "password": "prwr29hx",
  "channel": 2
}`)

	//這邊可以任意變換 http method  GET、POST、PUT、DELETE
	req, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("origin", "https://www.feastogether.com.tw")
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", cookie )
	req.Header.Set("referer", indexPage)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.104 Safari/537.36")
	req.Header.Set("x-csrf-token", csrfToken)
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	loginResp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer loginResp.Body.Close()
	loginResultStr , err := ioutil.ReadAll(loginResp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	//fmt.Printf("%s\n", loginResultStr)

	var loginResult LoginResult
	json.Unmarshal(loginResultStr, &loginResult)

	fmt.Println("已登入")
	/*userAccessToken := loginResult.Results.UserAccessToken

	fmt.Println(loginResult.Rcrm.Rm)
	fmt.Println("AccessToken:", userAccessToken)*/


	otGetPossible(client, cookie, csrfToken)
	bookingCheck(client, cookie, csrfToken)
	/*orderSet(client, cookie, csrfToken)*/
}
