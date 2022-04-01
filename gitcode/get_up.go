package gitcode

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/valyala/fastjson"
)

func login(token string) *gitlab.Client {
	u, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitcode.net"))
	if err != nil {
		log.Fatal(err)
	}
	return u
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

func GetUp(privateToken string, projectId, issueId int, city string) {

	u := login(privateToken)

	// POST /projects/:id/issues/:issue_iid/notes

	// curl --request POST --header "PRIVATE-TOKEN: <your_access_token>" "https://gitlab.example.com/api/v4/projects/5/issues/11/notes?body=note"

	//isTodayHaveRecord := isTodayHaveGetup(u, projectId, issueId)
	//if isTodayHaveRecord {
	//	return
	//}

	msg, isGetupEarly := makeGetupMsg(city)

	if !isGetupEarly {
		return
	}

	opt := &gitlab.CreateIssueNoteOptions{
		Body: gitlab.String(msg),
	}

	note, _, err := u.Notes.CreateIssueNote(projectId, issueId, opt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(note)

}

func isTodayHaveGetup(u *gitlab.Client, projectId, issueId int) bool {
	opt := &gitlab.ListIssueNotesOptions{
		ListOptions: gitlab.ListOptions{},
		OrderBy:     gitlab.String("created_at"),
		Sort:        gitlab.String("asc"),
	}
	issueNotes, _, err := u.Notes.ListIssueNotes(projectId, issueId, opt)
	if err != nil {
		log.Fatal(err)
	}

	if len(issueNotes) == 0 {
		return false
	}

	lastNoteCreatedAt := issueNotes[len(issueNotes)-1].CreatedAt
	isToday := lastNoteCreatedAt.Year() == time.Now().Year() && lastNoteCreatedAt.Month() == time.Now().Month() && lastNoteCreatedAt.Day() == time.Now().Day()
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
	defer resp.Body.Close()
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
