package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"strings"

	"github.com/julienschmidt/httprouter"

	"net/url"
)

// URL 是一個建構式，裡面含兩個成員
type URL struct {
	Source      string `bson:"source"`
	Destination string `bson:"destination"`
}

// Request 主要用來取URL的請求
type Request struct {
	Host string
}

func getwholeDomainURL(wholeDomainURL string) string {
	var resultURL string
	if strings.Contains(wholeDomainURL, "http://localhost") {
		resultURL = "http://localhost"
	} else {
		if strings.HasPrefix(wholeDomainURL, "https") {
			resultURL = fmt.Sprintf("%shttps://", resultURL)
		} else {
			resultURL = fmt.Sprintf("%shttp://", resultURL)
		}

		u, err := url.Parse(wholeDomainURL)
		if err != nil {
			log.Fatal(err)
		}
		parts := strings.Split(u.Hostname(), ".")
		domain := parts[len(parts)-2] + "." + parts[len(parts)-1]

		resultURL = fmt.Sprintf("%s%s", resultURL, domain)
	}
	return resultURL
}

func checkDB() string {
	sess, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sess.Close()
	err = sess.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// MongoDB server 正常運行中
	return "Healthy"
}

func getCollections() *mgo.Collection {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("admin")
	db.Login("adminrj", "123456")
	collections := db.C("urls")
	return collections
}

// Urls 被用來呈現URL列表，以及初始輸入畫面
func Urls(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 解析 HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<div class='container'><h4>歡迎!</h2>\n")
	t, err := template.ParseFiles("urls.html")
	if err != nil {
		fmt.Println(err)
	}

	collections := getCollections()
	result := []URL{}
	allUrls := collections.Find(nil).All(&result)
	if allUrls != nil {
		log.Fatal(allUrls)
	}

	var resultHTML string
	resultHTMLSource := make([]string, len(result))
	resultHTMLDestination := make([]string, len(result))
	// result 格式: [{SOURCE SOURCE_VALUE DESTINATION DESTINATION_VALUE} {SOURCE SOURCE_VALUE DESTINATION DESTINATION_VALUE}]
	tmpIndex := 0

	for _, listUrls := range result {
		// listUrls 格式: {SOURCE SOURCE VALUE DESTINATION DESTINATION_VALUE}
		resultSecond, _ := json.Marshal(listUrls)
		jsonData := []byte(resultSecond)
		var v interface{}
		json.Unmarshal(jsonData, &v)
		data := v.(map[string]interface{})
		// data 格式: map[Destination:DESTINATION_VALUE Source:SOURCE_VALUE]
		for index, value := range data {
			// index 格式: Destination 或是 Source
			// value 格式: DESTINATION_VALUE 或是 SOURCE_VALUE
			if index == "Source" {
				resultHTMLSource[tmpIndex] = fmt.Sprintf("%s", value)
			}
			if index == "Destination" {
				resultHTMLDestination[tmpIndex] = fmt.Sprintf("%s", value)
			}
		}
		tmpIndex++

	}

	i := 0
	for i < len(result) {
		resultHTML = fmt.Sprintf("%s<tr><td>%s</td><td>%s</td></tr>", resultHTML, resultHTMLSource[i], resultHTMLDestination[i])
		i++
	}

	items := struct {
		Result string
	}{
		Result: resultHTML,
	}
	t.Execute(w, items)
	fmt.Fprint(w, "</div>\n")

}

// InsertURL 單純被用來插入資料到Mongo資料庫
func InsertURL(text string, text02 string) string {
	collections := getCollections()
	err := collections.Insert(&URL{Source: text, Destination: text02})
	if err != nil {
		panic(err)
		// 這邊就會被中斷了，所以不用自定義return
	}
	fmt.Println(err)
	return "Successful"
}

// AddURL 用來做完整的插入資料庫的流程(含顯示結果通知的訊息)
func AddURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 這判斷是我為了保險而做的，其實在最下方main()的router地方，就會判斷method是否如期望的method
	if r.Method == "POST" {

		destination := r.FormValue("destination")
		// 如果目的地欄位沒值或是格式不是網址，便會轉址回去原本的列表頁
		if destination == "" || len(destination) == 0 || !isValidURL(destination) {
			http.Redirect(w, r, r.Header.Get("Referer"), 301)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<link rel='stylesheet' href='https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css'>")

		var randString string
		randString = RandStringURL(6, r.FormValue("custom_alias")) //length參考 https://tinyurl.com/

		fmt.Fprint(w, "<div class='container'><h4>結果</h4>")
		// 這是沒有找到相符的URL，這時是我們預期的，然後才要新增URL
		InsertURL(randString, destination)
		fmt.Fprint(w, "<label>插入資料成功</label><br>")
		var resultInsertString string
		resultInsertString = fmt.Sprintf("</label><a href='%s'>Source</a> => <a href='%s'>Destination</a><br>", randString, destination)
		fmt.Fprint(w, resultInsertString)

		resultDomainURL := getwholeDomainURL(r.Header.Get("Referer"))
		resultURL := fmt.Sprintf("%s/urls", resultDomainURL)

		var resultButtonString string
		resultButtonString = fmt.Sprintf("<button class='btn btn-success' onclick='javascript:location.href=\"%s\"'>回到列表頁</button><br>", resultURL)
		fmt.Fprint(w, resultButtonString)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString 單純只是被用來產生隨機字串
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// RandStringURL 產生沒被存在資料庫內的URL
func RandStringURL(n int, alias string) string {
	var randomMiddleValue string
	randomMiddleValue = GenerateRandomString(n)

	var newAlias string
	if alias == "" {
		// 表示這個alias是我方自動算出
		if !strings.Contains(alias, "auto-") {
			newAlias = fmt.Sprintf("auto-%s", GenerateRandomString(n))
		} else {
			newAlias = fmt.Sprintf("%s", GenerateRandomString(n))
		}
	} else {
		// 表示這個alias是使用者自定義
		if !strings.Contains(alias, "custom-") {
			newAlias = fmt.Sprintf("custom-%s", alias)
		} else {
			newAlias = fmt.Sprintf("%s", alias)
		}
	}

	finalURL := fmt.Sprintf("http://localhost/redirect/%s/%s", randomMiddleValue, newAlias)
	fmt.Println(finalURL)

	collections := getCollections()
	result := URL{}
	errFind := collections.Find(bson.M{"source": finalURL}).One(&result)

	if errFind == nil {
		// 代表已經found，不是我預期的，所以再次跑產生URL func
		finalURL = RandStringURL(n, newAlias)
	}

	return finalURL
}

// RedirectURL 用來實現此專案要的轉址功能，就是讀取到資料庫的「來源」要轉址到「目的地」。
func RedirectURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	wholeURL := fmt.Sprintf("http://%s%s", r.Host, r.URL.Path)

	redirectResult := GetRedirectURL(wholeURL)
	if redirectResult != "Failed" && strings.Contains(redirectResult, "http") {
		http.Redirect(w, r, redirectResult, 301)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<div class='container'><h4>結果：轉址失敗(可能「來源/Source」網址不存在資料庫，或是「目的地/Destination」網址格式錯誤)</h4>")
	fmt.Fprintf(w, "您發出的請求中的整個網址: %s <br>\n", wholeURL)
	fmt.Fprintf(w, "您發出的請求中的第一個網址片段: %s <br>\n", ps.ByName("random_string"))
	fmt.Fprintf(w, "您發出的請求中的第二個網址片段: %s <br>\n", ps.ByName("alias"))
	fmt.Fprintf(w, "您發出的請求的「結果」為(若Failed即從資料庫中找不到對應的網址): %s <br>\n", redirectResult)
	fmt.Fprint(w, "<link rel='stylesheet' href='https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css'>")
	fmt.Fprintf(w, "<button class='btn btn-success' onclick='javascript:location.href=\"http://%s%s/urls\"'>回到列表頁</button><br>", r.Host, r.URL.Path)
	fmt.Fprint(w, "</div>")

}
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// GetRedirectURL 只是被用來「取得」資料庫裡要被用來轉址的「網址」
func GetRedirectURL(source string) string {
	collections := getCollections()
	result := URL{}
	errFind := collections.Find(bson.M{"source": source}).One(&result)
	if errFind == nil {
		// 表示這個網址來源有跟資料庫的source欄位match到，所以要回傳回去給"轉址"使用
		return result.Destination
	}
	return "Failed"
}

func main() {

	// Router (RESTful)
	router := httprouter.New()
	router.GET("/urls", Urls)
	router.POST("/url", AddURL)
	router.GET("/redirect/:random_string/:alias", RedirectURL)
	http.ListenAndServe(":80", router)

}
