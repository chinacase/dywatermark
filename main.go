package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	//抖音分享的URL
	regDyURL   = `https://v.douyin.com/+[\w]+?/`
	regItemIds = `/video/([\d]+?)/`
	//url_list":["URL"]
	//https://aweme.snssdk.com/aweme/v1/playwm/?video_id=v0d00f500000c0ppd496q7knshhfoa80&ratio=720p&line=0
	redVideoURL = `"url_list":\["(https://aweme.snssdk.com/[\w/&?.=]+?)?"\]`
)

//匹配抖音分享的URL
func matchURLByExpr(shareStr, expr string) (dyURL string, err error) {
	r, err := regexp.Compile(expr)
	if err != nil {
		return
	}
	dyURL = r.FindString(shareStr)
	if len(dyURL) == 0 {
		err = errors.New("未匹配")
		return
	}
	return
}

func matchAllByExpr(matchStr, expr string) [][]string {
	matchs := make([][]string, 0)
	r, err := regexp.Compile(expr)
	if err != nil {
		return matchs
	}
	matchs = r.FindAllStringSubmatch(matchStr, -1)
	return matchs
}

func getInfo(url string) (resHTML string, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4406.0 Safari/537.36")
	resp, err1 := client.Do(req)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()
	//fmt.Println(doc.Html())
	body, _ := ioutil.ReadAll(resp.Body)
	resHTML = string(body)
	//resHTML = mahonia.NewDecoder("gbk").ConvertString(string(body)) //gbk=>utf8
	return
}

func getNowoterURL(url string) (resHTML string, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
	resp, err1 := client.Do(req)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()
	//fmt.Println(doc.Html())
	resHTML = resp.Request.URL.String()
	//resHTML = mahonia.NewDecoder("gbk").ConvertString(string(body)) //gbk=>utf8
	return
}

//匹配itemids
func getItemIds(url string) (itemIds string, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4406.0 Safari/537.36")
	resp, err1 := client.Do(req)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()
	locltionURL := resp.Request.URL.String()
	itemidss := matchAllByExpr(locltionURL, regItemIds)
	if len(itemidss) == 0 {
		err = errors.New("未匹配到itemids")
		return
	}
	//fmt.Println(itemidss)
	itemIds = itemidss[0][1]
	return
}

func getVideoURL(itemids string) string {
	//https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=6932049937869982991
	url := "https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids="
	url = url + itemids
	res, err := getInfo(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	r := matchAllByExpr(res, redVideoURL)

	if len(r) == 0 {
		return ""
	}
	videoURL := r[0][1]
	videoURL = strings.Replace(videoURL, "playwm", "play", 1)
	return videoURL
}

func main() {
	//抖音分享地址
	shareStr := "不懂的你可以问老师🙂%黑丝 %大长腿  https://v.douyin.com/e15Fdtq/ 鳆制此鏈接，打kaiDouyin搜索，直接观kan视频"

	url, err := matchURLByExpr(shareStr, regDyURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	itemids, err := getItemIds(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	videoURL := getVideoURL(itemids)
	nowaterURL, _ := getNowoterURL(videoURL)
	fmt.Println("抖音无水印地址：", nowaterURL)
}
