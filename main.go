package main

import (
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/mmcdole/gofeed"
)

// loadArticles 抓取博客 RSS 的最新文章。
// 抓取失败时返回空列表而不是 panic，避免博客临时不可用时 README 自动更新流程中断。
func loadArticles() []*gofeed.Item {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://gs2ygc.online/rss2.xml")
	if err != nil {
		log.Printf("warn: parse rss err: %v", err)
		return nil
	}
	if len(feed.Items) > 5 {
		return feed.Items[:5]
	}
	return feed.Items
}

func main() {
	b, err := ioutil.ReadFile("./readme.tpl")
	if err != nil {
		log.Panicf("read tpl err: %v", err)
	}
	tpl, err := template.New("").Parse(string(b))
	if err != nil {
		log.Panicf("tpl init err: %v", err)
	}
	f, err := os.OpenFile("./README.md", os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Panicf("open readme err: %v", err)
	}

	err = tpl.Execute(f, map[string]interface{}{
		"Articles": loadArticles(),
	})
	if err != nil {
		log.Panicf("tpl exec err: %v", err)
	}
}
