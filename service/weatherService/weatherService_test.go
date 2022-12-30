package weatherservice

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"seongwoo/go/fiber/config"
	"testing"
)

func TestGetWeatherEight(t *testing.T) {
	assert := assert.New(t)
	// 1. 시나리오 생성
	// 2. 메서드 호출
	// 3. 내가 생성한 값과 메서드에서 호출되서 나온 값이 같은지 확인
	serviceKey := config.Config("API_SERVICE_KEY")
	getWeatherDto := GetWeatherDto{}
	foreCastDto := []Forecast{
		{
			BaseDate:  "20221223",
			BaseTime:  "0800",
			Category:  "TMP",
			FcstDate:  "20221223",
			FcstTime:  "0900",
			FcstValue: "-13",
			Nx:        60,
			Ny:        127,
		},
		{
			BaseDate:  "20221223",
			BaseTime:  "0800",
			Category:  "UUU",
			FcstDate:  "20221223",
			FcstTime:  "0900",
			FcstValue: "4.1",
			Nx:        60,
			Ny:        127,
		},
		{
			BaseDate:  "20221223",
			BaseTime:  "0800",
			Category:  "VVV",
			FcstDate:  "20221223",
			FcstTime:  "0900",
			FcstValue: "-2.4",
			Nx:        60,
			Ny:        127,
		},
		{
			BaseDate:  "20221223",
			BaseTime:  "0800",
			Category:  "VEC",
			FcstDate:  "20221223",
			FcstTime:  "0900",
			FcstValue: "301",
			Nx:        60,
			Ny:        127,
		},
		{
			BaseDate:  "20221223",
			BaseTime:  "0800",
			Category:  "WSD",
			FcstDate:  "20221223",
			FcstTime:  "0900",
			FcstValue: "4.8",
			Nx:        60,
			Ny:        127,
		},
		{
			BaseDate:  "20221223",
			BaseTime:  "0800",
			Category:  "SKY",
			FcstDate:  "20221223",
			FcstTime:  "0900",
			FcstValue: "3",
			Nx:        60,
			Ny:        127,
		},
		{
			BaseDate:  "20221223",
			BaseTime:  "0800",
			Category:  "PTY",
			FcstDate:  "20221223",
			FcstTime:  "0900",
			FcstValue: "0",
			Nx:        60,
			Ny:        127,
		},
		{
			BaseDate:  "20221223",
			BaseTime:  "0800",
			Category:  "POP",
			FcstDate:  "20221223",
			FcstTime:  "0900",
			FcstValue: "20",
			Nx:        60,
			Ny:        127,
		},
		{
			BaseDate:  "20221223",
			BaseTime:  "0800",
			Category:  "WAV",
			FcstDate:  "20221223",
			FcstTime:  "0900",
			FcstValue: "0",
			Nx:        60,
			Ny:        127,
		},
		{
			BaseDate:  "20221223",
			BaseTime:  "0800",
			Category:  "PCP",
			FcstDate:  "20221223",
			FcstTime:  "0900",
			FcstValue: "강수없음",
			Nx:        60,
			Ny:        127,
		},
	}
	fullUrl := fmt.Sprintf("https://apis.data.go.kr/1360000/VilageFcstInfoService_2.0/getVilageFcst?ServiceKey=%v&numOfRows=10&base_date=20221223&base_time=0800&nx=60&ny=127&dataType=JSON", serviceKey)
	res, err := http.Get(fullUrl)
	assert.Equal("200 OK", res.Status)
	fmt.Println(err)
	if err != nil {
		assert.Error(err, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	assert.Equal("200 OK", res.Status)
	if err != nil {
		assert.Error(err, err)
	}

	_ = json.Unmarshal(body, &getWeatherDto)
	item := getWeatherDto.Response.Body.Items.Item
	fmt.Println(item)
	fmt.Println(foreCastDto)
	assert.Equal(item, foreCastDto, "foreCastDto 와 item 는 같아야 합니다.")
}
