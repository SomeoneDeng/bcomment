package main

import (
	"breplies/dao"
	"breplies/data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	REPLY_BASE_URL   = "https://api.bilibili.com/x/v2/reply"
	ARCHIVE_BASE_URL = "https://api.bilibili.com/x/web-interface/dynamic/region"
)

var daoHelper *dao.DaoHelper

func init() {
	daoHelper = new(dao.DaoHelper)
	daoHelper.New()
}

func (scraper *Scraper) ReplyConsume() {
	for reply := range scraper.replyChan {
		println(reply.Content.Message)
		daoHelper.SaveComment(reply)
		for _, v := range reply.Replies {
			daoHelper.SaveComment(v)
		}
	}
}

// 获取某个视频下的回复
func getReplyData(pn int, ps int, ttype int, oid int, sort int, scraper *Scraper) {
	reqString := fmt.Sprintf("%s?pn=%d&type=%d&oid=%d&sort=%d&ps=%d", REPLY_BASE_URL, pn, ttype, oid, sort, ps)
	a := rand.Intn(len(scraper.proxys))
	proxy := scraper.proxys[a]
	println(reqString)
	println("proxy: ", proxy, " index: ", a)
	proxyUrl, _ := url.Parse(proxy)
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	resp, e := client.Get(reqString)
	if e != nil {
		println("获取reply网络错误", e.Error())
	}
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		println(e.Error())
	}
	var respBody data.ReplyResponse
	e = json.Unmarshal(body, &respBody)
	if e != nil {
		println(e.Error())
		println(string(body))
	}
	for _, v := range respBody.Data.Replies {
		scraper.replyChan <- v
	}
	if len(respBody.Data.Replies) > 0 {
		time.Sleep(time.Millisecond * 500)
		go getReplyData(pn+1, ps, ttype, oid, sort, scraper)
	}
}

// 从视频获取评论
func getArchiveReplies(scraper *Scraper) {
	for archive := range scraper.archiveChan {
		pn := 1
		// goroutine here to speed up
		getReplyData(pn, 20, 1, archive.Aid, 0, scraper)
	}
}

// extract data from response
func getArchiveData(ps int, rid int, pn int, scraper *Scraper) data.ArchiveResponse {
	reqString := fmt.Sprintf("%s?ps=%d&rid=%d&pn=%d", ARCHIVE_BASE_URL, ps, rid, pn)
	a := rand.Intn(len(scraper.proxys))
	proxy := scraper.proxys[a]
	println(reqString)
	println("proxy: ", proxy, " index: ", a)
	proxyUrl, _ := url.Parse(proxy)
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	resp, e := client.Get(reqString)
	if e != nil {
		println("获取archive网络错误", e.Error())
	}
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		println(e.Error())
	}
	var respBody data.ArchiveResponse
	e = json.Unmarshal(body, &respBody)
	if e != nil {
		println(e.Error())
		println(string(body))
	}
	return respBody
}

// 获取某个分类下的视频
func getArchives(ps int, rid int, pn int, scraper *Scraper) {
	respBody := getArchiveData(ps, rid, pn, scraper)
	// get rest pages
	for respBody.Code != -404 {
		for _, v := range respBody.Data.Archives {
			println("[ " + v.Title + "------" + strconv.Itoa(v.Aid) + " ]")
			time.Sleep(time.Second)
			// send to channel
			scraper.archiveChan <- v
		}
		pn += 1
		respBody = getArchiveData(ps, rid, pn, scraper)
	}
}

// 爬虫
type Scraper struct {
	Rids        [] int
	archiveChan chan data.Archive
	replyChan   chan data.Reply
	proxys      []string
}

func (scraper *Scraper) Start() {
	for _, rid := range scraper.Rids {
		go getArchives(20, rid, 1, scraper)
	}
	go scraper.ReplyConsume()
	getArchiveReplies(scraper)
}

func main() {
	scraper := &Scraper{
		Rids:        []int{17}, // 分类Id
		archiveChan: make(chan data.Archive),
		replyChan:   make(chan data.Reply),
		proxys: []string{
			"http://127.0.0.1:4397",
		},
	}
	scraper.Start()
}
