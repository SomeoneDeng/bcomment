package main

import (
	"breplies/dao"
	"breplies/data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
func getReplyData(pn int, ps int, ttype int, oid int, sort int, replyChan chan data.Reply) {
	reqString := fmt.Sprintf("%s?pn=%d&type=%d&oid=%d&sort=%d&ps=%d", REPLY_BASE_URL, pn, ttype, oid, sort, ps)
	//println("getting replies from ---> ", reqString)
	resp, e := http.Get(reqString)
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
	}
	for _, v := range respBody.Data.Replies {
		replyChan <- v
	}
	if len(respBody.Data.Replies) > 0 {
		go getReplyData(pn+1, ps, ttype, oid, sort, replyChan)
	}
}

// 从视频获取评论
func getArchiveReplies(scraper *Scraper) {
	for archive := range scraper.archiveChan {
		pn := 1
		// goroutine here to speed up
		getReplyData(pn, 20, 1, archive.Aid, 0, scraper.replyChan)
	}
}

// extract data from response
func getArchiveData(ps int, rid int, pn int) data.ArchiveResponse {
	reqString := fmt.Sprintf("%s?ps=%d&rid=%d&pn=%d", ARCHIVE_BASE_URL, ps, rid, pn)
	resp, e := http.Get(reqString)
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
	}
	return respBody
}

// 获取某个分类下的视频
func getArchives(ps int, rid int, pn int, archiveChan chan data.Archive) {
	respBody := getArchiveData(ps, rid, pn)
	// get rest pages
	for respBody.Code != -404 {
		for _, v := range respBody.Data.Archives {
			println("[ " + v.Title + "------" + strconv.Itoa(v.Aid) + " ]")
			// send to channel
			archiveChan <- v
		}
		pn += 1
		respBody = getArchiveData(ps, rid, pn)
	}
}

// 爬虫
type Scraper struct {
	Rids        [] int
	archiveChan chan data.Archive
	replyChan   chan data.Reply
}

func (scraper *Scraper) Start() {
	for _, rid := range scraper.Rids {
		go getArchives(20, rid, 1, scraper.archiveChan)
	}
	go scraper.ReplyConsume()
	getArchiveReplies(scraper)
}

func main() {
	//getReplies(2, 1, 53278740, 1)

	//getArchives(20, 17, 2)

	scraper := &Scraper{
		Rids:        []int{17, 189, 15}, // 分类
		archiveChan: make(chan data.Archive),
		replyChan:   make(chan data.Reply),
	}
	scraper.Start()
}
