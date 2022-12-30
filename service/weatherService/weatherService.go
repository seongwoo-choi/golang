package weatherservice

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
	"io"
	"net/http"
	"seongwoo/go/fiber/config"
	"strconv"
)

type GetWeatherDto struct {
	Response struct {
		Header Header `json:"header"`
		Body   struct {
			DataType string `json:"dataType"`
			Items    struct {
				Item []Forecast `json:"item"`
			} `json:"items"`
			PageNo     int `json:"pageNo"`
			NumOfRows  int `json:"numOfRows"`
			TotalCount int `json:"totalCount"`
		} `json:"body"`
	} `json:"response"`
}

type Header struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMag"`
}

type Forecast struct {
	BaseDate  string `json:"baseDate"`
	BaseTime  string `json:"baseTime"`
	Category  string `json:"category"`
	FcstDate  string `json:"fcstDate"`
	FcstTime  string `json:"fcstTime"`
	FcstValue string `json:"fcstValue"`
	Nx        int    `json:"nx"`
	Ny        int    `json:"ny"`
}

type Weather struct {
	Time string `json:"time"`
	PTY  string `json:"pty"`
	SKY  string `json:"sky"`
	TMP  string `json:"tmp"`
	PCP  string `json:"pcp"`
	SNO  string `json:"sno"`
}

type WeatherVO struct {
	Alert    []Alert `json:"alert"`
	AvgTmpAM int     `json:"avgTmpAM"`
}

type Alert struct {
	Time string `json:"time"`
	PCP  string `json:"pcp"`
	SNO  string `json:"sno"`
	PTY  string `json:"pty"`
}

func getWeatherData(fullUrl string) (GetWeatherDto, error) {
	var getWeatherDto GetWeatherDto
	response, err := http.Get(fullUrl)
	if err != nil {
		return getWeatherDto, err
	}
	defer response.Body.Close()

	// Read the response body
	// body is []byte type
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return getWeatherDto, err
	}
	/** JSON Bytes => Struct 형변환 과정 **/

	_ = json.Unmarshal(body, &getWeatherDto)

	return getWeatherDto, nil
}

func getWeathers(data []Forecast) []Weather {
	var weathers []Weather
	weather := Weather{Time: "null", TMP: "null", SKY: "null", PCP: "null", SNO: "null", PTY: "null"}
	for _, v := range data {
		weather.Time = v.FcstTime
		// 없음(0), 비(1), 비/눈(2), 눈(3), 빗방울(5), 빗방울 눈날림(6), 눈날림(7)
		if v.Category == "PTY" {
			switch v.FcstValue {
			case "0":
				weather.PTY = "비 눈 내리지 않음"
				fallthrough
			case "1":
				weather.PTY = "비"
				fallthrough
			case "2":
				weather.PTY = "비/눈"
				fallthrough
			case "3":
				weather.PTY = "눈"
				fallthrough
			case "5":
				weather.PTY = "빗방울"
				fallthrough
			case "6":
				weather.PTY = "빗방울 및 눈날림"
				fallthrough
			case "7":
				weather.PTY = "눈날림"
			}
		}
		if v.Category == "TMP" {
			weather.TMP = v.FcstValue
		}
		if v.Category == "SKY" {
			switch v.FcstValue {
			case "1":
				weather.SKY = "맑음"
				fallthrough
			case "3":
				weather.SKY = "구름많음"
				fallthrough
			case "4":
				weather.SKY = "흐림"
			}
		}
		if v.Category == "PCP" {
			switch v.FcstValue {
			case "강수없음":
				weather.PCP = v.FcstValue
				fallthrough
			default:
				weather.PCP = v.FcstValue
			}
		}
		if v.Category == "SNO" {
			switch v.FcstValue {
			case "적설없음":
				weather.SNO = v.FcstValue
				fallthrough
			default:
				weather.SNO = v.FcstValue
			}
			weathers = append(weathers, weather)
		}
	}
	return weathers
}

func sendMail(text string) error {
	email := config.Config("GMAIL_ID")
	email_app_pwd := config.Config("GMAIL_APP_PWD")
	msg := gomail.NewMessage()
	msg.SetHeader("From", email)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Today Weather Alarm")
	msg.SetBody("text/html", text)

	n := gomail.NewDialer("smtp.gmail.com", 587, email, email_app_pwd)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		return err
	} else {
		return nil
	}
}

func GetWeather(c *fiber.Ctx) error {
	serviceKey := config.Config("API_SERVICE_KEY")
	// Make the GET request to the API
	url := "https://apis.data.go.kr/1360000/VilageFcstInfoService_2.0/getVilageFcst?"
	numOfRows := c.Query("numOfRows")
	base_date := c.Query("base_date")
	base_time := c.Query("base_time")
	dataType := "JSON"
	// 우리집 위경도를 기상청 격자정보로 변환했습니다.
	// https://fronteer.kr/service/kmaxy
	nx := "58"
	ny := "128"
	fullUrl := fmt.Sprintf("%vServiceKey=%v&numOfRows=%v&dataType=%v&base_date=%v&base_time=%v&nx=%v&ny=%v", url, serviceKey, numOfRows, dataType, base_date, base_time, nx, ny)
	fmt.Println(fullUrl)
	now := c.Query("base_date")

	getWeatherDto, err := getWeatherData(fullUrl)
	if err != nil {
		return c.JSON(err)
	}
	data := getWeatherDto.Response.Body.Items.Item

	/**
	PTY(강수 형태): 없음=0, 비(1), 비/눈(2), 눈(3), 빗방울(5), 빗방울 눈날림(6), 눈날림(7)
	SKY(하늘 상태): 맑음=1, 구름많음(3), 흐림(4)
	TMP(1시간 기온): 온도
	PCP(1시간 강수량): 강수없음 or int(mm)
	SNO(1시간 적설량): 적설없음 or int(cm)
	REH(습도): 29 ==> 습도 29%
	TMN(일 최저 기온)
	TMX(일 최고 기온)
	POP(강수확률)
	**/
	weathers := getWeathers(data)

	var weatherVo WeatherVO
	first, _ := strconv.Atoi(weathers[0].TMP)
	weatherVo.AvgTmpAM = first
	for _, w := range weathers {
		// PCP, SNO '강수없음', '적설없음' 이 아닐 경우
		if w.PCP != "강수없음" || w.SNO != "적설없음" {
			weatherVo.Alert = append(weatherVo.Alert, Alert{Time: w.Time, PCP: w.PCP, SNO: w.SNO, PTY: w.PTY})
		}
	}

	var text string
	var we string
	if len(weatherVo.Alert) == 0 {
		text = fmt.Sprintf("오늘은 %v 입니다. 현재 시간은 %v시 입니다. <br><br>현재 온도 평균은 %v도 입니다.<br><br> 오늘은 비 혹은 눈이 내리지 않는 날입니다.<br><br>", now, base_time, weatherVo.AvgTmpAM)
	} else {
		for _, weather := range weatherVo.Alert {
			text = fmt.Sprintf("오늘은 %v 입니다. 현재 시간은 %v시 입니다. <br><br>현재 온도 평균은 %v 입니다.<br><br>", now, base_time, weatherVo.AvgTmpAM)
			we += fmt.Sprintf("%v 시에 %v이(가) 있을 예정이며, 강수량은 %v 이고, 적설량은 %v 입니다.<br><br>", weather.Time, weather.PTY, weather.PCP, weather.SNO)
		}
	}

	text = text + we

	err = sendMail(text)
	if err != nil {
		return c.Status(500).JSON("email 을 보내는 중 오류가 발생하였습니다.")
	}

	return c.Status(200).JSON(weatherVo)
}
