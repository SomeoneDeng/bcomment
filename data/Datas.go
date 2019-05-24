package data

import "strconv"

// 回复
type Reply struct {
	Rpid      int     `json:"rpid"`
	Oid       int     `json:"oid"`
	Type      int     `json:"type"`
	Mid       int     `json:"mid"`
	Root      int     `json:"root"`
	Parent    int     `json:"parent"`
	Dialog    int     `json:"dialog"`
	Count     int     `json:"count"`
	Rcount    int     `json:"rcount"` // reply count
	State     int     `json:"state"`
	Fansgrade int     `json:"fansgrade"`
	Attr      int     `json:"attr"`
	Ctime     int     `json:"ctime"`
	RpidStr   string  `json:"rpid_str"`
	RootStr   string  `json:"root_str"`
	ParentStr string  `json:"parent_str"`
	Like      int     `json:"like"`
	Action    int     `json:"action"`
	Replies   []Reply `json:"replies"`
	Content   struct {
		Message string `json:"message"`
		Plat    int    `json:"plat"`
		Device  string `json:"device"`
	}
	Member struct{
		Mid string `json:"mid"`
		Uname string `json:"uname"`
	} `json:"member"`
}

// 回复响应
type ReplyResponse struct {
	Data struct {
		Page struct {
			Num    int `json:"num"`
			Size   int `json:"size"`
			Count  int `json:"count"`
			Acount int `json:"acount"`
		} `json:"page"`
		Replies []Reply `json:"replies"`
	} `json:"data"`
}

// 每个分类下的视频信息
type Archive struct {
	Aid       int    `json:"aid"`
	Videos    int    `json:"videos"`
	Tid       int    `json:"tid"`
	Tname     string `json:"tname"`
	Copyright int    `json:"copyright"`
	Pic       string `json:"pic"`
	Title     string `json:"title"`
	Pubdate   int    `json:"pubdate"`
	Ctime     int    `json:"ctime"`
	Desc      string `json:"desc"`
	State     int    `json:"state"`
	Attribute int    `json:"attribute"`
	Duration  int    `json:"duration"`
	Rights    struct {
		Bp            int `json:"bp"`
		Elec          int `json:"elec"`
		Download      int `json:"download"`
		Movie         int `json:"movie"`
		Pay           int `json:"pay"`
		Hd5           int `json:"hd5"`
		NoReprint     int `json:"no_reprint"`
		Autoplay      int `json:"autoplay"`
		UgcPay        int `json:"ugc_pay"`
		IsCooperation int `json:"is_cooperation"`
		UgcPayPreview int `json:"ugc_pay_preview"`
	} `json:"rights"`
	Owner struct {
		Mid  int    `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"owner"`
	Stat struct {
		Aid      int `json:"aid"`
		View     int `json:"view"`
		Danmaku  int `json:"danmaku"`
		Reply    int `json:"reply"`
		Favorite int `json:"favorite"`
		Coin     int `json:"coin"`
		Share    int `json:"share"`
		NowRank  int `json:"now_rank"`
		HisRank  int `json:"his_rank"`
		Like     int `json:"like"`
		Dislike  int `json:"dislike"`
	} `json:"stat"`
	Dynamic   string `json:"dynamic"`
	Cid       int    `json:"cid"`
	Dimension struct {
		Width  int `json:"width"`
		Height int `json:"height"`
		Rotate int `json:"rotate"`
	} `json:"dimension"`
}

// 视频响应
type ArchiveResponse struct {
	Code int `json:"code"`
	Data struct {
		Page struct {
			Num   int `json:"num"`
			Size  int `json:"size"`
			Count int `json:"count"`
		}
		Archives [] Archive `json:"archives"`
	} `json:"data"`
}

func (content *Reply) ContentString() string {
	return "Message: " + content.Content.Message +
		"\nDevice: " +
		content.Content.Device +
		"\nPlat: " + strconv.Itoa(content.Content.Plat) +
		"\nreply parent: " + strconv.Itoa(content.Parent) +
		"\n====================\n"
}

func (data *ReplyResponse) PageString() string {
	return "Num: " + strconv.Itoa(data.Data.Page.Num) +
		"\nAcount: " + strconv.Itoa(data.Data.Page.Acount) +
		"\nCount: " + strconv.Itoa(data.Data.Page.Count) +
		"\nSize: " + strconv.Itoa(data.Data.Page.Size) +
		"\n====================\n"
}
