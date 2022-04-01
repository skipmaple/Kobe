package github

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-github/v43/github"
	"github.com/valyala/fastjson"
	"golang.org/x/oauth2"
)

func login(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}

func makeGetupMsg(city string) (body string, isGetupEarly bool) {
	weatherMsg := ""
	if len(city) > 0 {
		weatherMsg = currentWeather(city)
	}

	cstZone := time.FixedZone("GMT", 8*3600) // 东八
	now := time.Now().In(cstZone)
	// 3点到18点起床都有效
	isGetupEarly = 3 < now.Hour() && now.Hour() < 24
	body = fmt.Sprintf("今天的起床时间是--%s\n", now.Format(time.Kitchen))

	if len(weatherMsg) > 0 {
		body = fmt.Sprintf("%s\n %s\n", body, weatherMsg)
	}

	return body, isGetupEarly
}

func GetUp(privateToken, owner, repo, city string) {
	u := login(privateToken)

	ctx := context.Background()

	// 检查有没有issue, 没有就创建
	issueListByOps := &github.IssueListByRepoOptions{}
	issues, _, _ := u.Issues.ListByRepo(ctx, owner, repo, issueListByOps)
	issue := &github.Issue{}
	if len(issues) == 0 {
		issueRequest := &github.IssueRequest{Title: github.String("起床时间记录")}
		newIssue, rsp, err := u.Issues.Create(ctx, owner, repo, issueRequest)
		if err != nil {
			log.Fatalf("Create Issue Error.\nResponse: %v\n Error: %v\n", rsp, err)
		} else {
			issue = newIssue
		}
	} else {
		issue = issues[0]
	}

	//isTodayHaveRecord := isTodayHaveGetup(u, owner, repo, issue)
	//if isTodayHaveRecord {
	//	log.Println("今天已经有起床记录了.")
	//	return
	//}

	msg, isGetupEarly := makeGetupMsg(city)

	if !isGetupEarly {
		log.Println("当前打卡时间 不在有效时间范围内.")
		return
	}

	// 创建issue评论
	issueComment := &github.IssueComment{
		Body: github.String(msg),
	}
	_, _, err := u.Issues.CreateComment(ctx, owner, repo, issue.GetNumber(), issueComment)
	if err != nil {
		log.Fatalf("Create Issue Comment Error: %v\n", err)
	}

	log.Println(msg)

}

func isTodayHaveGetup(u *github.Client, owner, repo string, issue *github.Issue) bool {
	ctx := context.Background()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	year, month, day := time.Now().In(loc).Date()
	since := time.Date(year, month, day, 0, 0, 0, 0, loc)
	issueListCommentOpts := &github.IssueListCommentsOptions{Since: &since, Sort: github.String("created"), Direction: github.String("desc")}
	issueComments, _, err := u.Issues.ListComments(ctx, owner, repo, issue.GetNumber(), issueListCommentOpts)
	if err != nil {
		log.Fatal(err)
	}

	if len(issueComments) == 0 {
		return false
	}

	latestNoteCreatedAt := issueComments[0].GetCreatedAt()
	isToday := latestNoteCreatedAt.Year() == time.Now().Year() && latestNoteCreatedAt.Month() == time.Now().Month() && latestNoteCreatedAt.Day() == time.Now().Day()
	if isToday {
		return true
	}

	return false
}

/*{
"location": {
		"name": "Beijing",
		"region": "Beijing",
		"country": "China",
		"lat": 39.93,
		"lon": 116.39,
		"tz_id": "Asia/Shanghai",
		"localtime_epoch": 1644818515,
		"localtime": "2022-02-14 14:01"
	},
"current": {
		"last_updated_epoch": 1644817500,
		"last_updated": "2022-02-14 13:45",
		"temp_c": -3.0,
		"temp_f": 26.6,
		"is_day": 1,
		"condition": {
				"text": "晴天",
				"icon": "//cdn.weatherapi.com/weather/64x64/day/113.png",
				"code": 1000
			},
		"wind_mph": 4.3,
		"wind_kph": 6.8,
		"wind_degree": 360,
		"wind_dir": "N",
		"pressure_mb": 1030.0,
		"pressure_in": 30.42,
		"precip_mm": 0.0,
		"precip_in": 0.0,
		"humidity": 36,
		"cloud": 0,
		"feelslike_c": -8.4,
		"feelslike_f": 17.0,
		"vis_km": 10.0,
		"vis_miles": 6.0,
		"uv": 2.0,
		"gust_mph": 11.6,
		"gust_kph": 18.7
	}
}*/

func currentWeather(city string) (res string) {
	params := url.Values{}
	params.Add("key", "bdafd6aabd07444b93952349221402")
	params.Add("q", city)
	params.Add("api", "no")
	params.Add("lang", "zh")

	resq := "http://api.weatherapi.com/v1/current.json?" + params.Encode()
	resp, err := http.Get(resq)
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	// Log the request body
	bodyString := string(body)
	log.Print(bodyString)

	var p fastjson.Parser
	weather, err := p.Parse(bodyString)
	if err != nil {
		log.Fatal(err)
	}

	res = fmt.Sprintf("现在的天气是 %.1f °C, %s", weather.GetFloat64("current", "temp_c"), string(weather.GetStringBytes("current", "condition", "text")))
	return res
}
