package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/internal/defination"
	agentapi "github.com/wilenceyao/humor-api/agent/humor"
	"github.com/wilenceyao/humor-api/common"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// 数据来源：国家气象科学数据中心 http://data.cma.cn/

var DefaultWeatherService *WeatherService = &WeatherService{}

type WeatherService struct {
}

type WeatherLive1Data struct {
	AqiForecast []struct {
		Date        string `json:"date"`
		PublishTime string `json:"publishTime"`
		Value       int64  `json:"value"`
	} `json:"aqiForecast"`
	AqiForecastHourly []struct {
		ForecastTime string `json:"forecastTime"`
		PublishTime  string `json:"publishTime"`
		Value        int64  `json:"value"`
	} `json:"aqiForecastHourly"`
	AqiHistory []struct {
		PubTime int64 `json:"pubTime"`
		Value   int64 `json:"value"`
	} `json:"aqiHistory"`
	AqiPoint []struct {
		CityID    int64   `json:"cityId"`
		Co        string  `json:"co"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		No2       string  `json:"no2"`
		O3        string  `json:"o3"`
		Pm10      string  `json:"pm10"`
		Pm25      string  `json:"pm25"`
		PointName string  `json:"pointName"`
		Primary   string  `json:"primary"`
		PubTime   int64   `json:"pubTime"`
		So2       string  `json:"so2"`
		Value     int64   `json:"value"`
	} `json:"aqiPoint"`
	AqiRank []struct {
		Aqi          int64  `json:"aqi"`
		CityName     string `json:"cityName"`
		Order        int64  `json:"order"`
		ProvinceName string `json:"provinceName"`
		PublishTime  string `json:"publishTime"`
	} `json:"aqiRank"`
	City struct {
		CityID        int64  `json:"cityId"`
		Counname      string `json:"counname"`
		Ianatimezone  string `json:"ianatimezone"`
		Name          string `json:"name"`
		Pname         string `json:"pname"`
		Secondaryname string `json:"secondaryname"`
		Timezone      string `json:"timezone"`
	} `json:"city"`
	Condition struct {
		Condition           string `json:"condition"`
		ConditionID         string `json:"conditionId"`
		Humidity            string `json:"humidity"`
		Icon                string `json:"icon"`
		Precipitation       string `json:"precipitation"`
		Pressure            string `json:"pressure"`
		RainSnowType        string `json:"rain_snow_type"`
		RainfallIntensity   string `json:"rainfall_intensity"`
		RainfallIntensity1h string `json:"rainfall_intensity_1h"`
		RealFeel            string `json:"realFeel"`
		SunRise             string `json:"sunRise"`
		SunSet              string `json:"sunSet"`
		Temp                string `json:"temp"`
		Tips                string `json:"tips"`
		Updatetime          string `json:"updatetime"`
		Uvi                 string `json:"uvi"`
		Vis                 string `json:"vis"`
		WindDegrees         string `json:"windDegrees"`
		WindDir             string `json:"windDir"`
		WindLevel           string `json:"windLevel"`
		WindSpeed           string `json:"windSpeed"`
	} `json:"condition"`
	Forecast []struct {
		ConditionDay     string `json:"conditionDay"`
		ConditionIDDay   string `json:"conditionIdDay"`
		ConditionIDNight string `json:"conditionIdNight"`
		ConditionNight   string `json:"conditionNight"`
		Humidity         string `json:"humidity"`
		Moonphase        string `json:"moonphase"`
		Moonrise         string `json:"moonrise"`
		Moonset          string `json:"moonset"`
		Pop              string `json:"pop"`
		PredictDate      string `json:"predictDate"`
		Qpf              string `json:"qpf"`
		Sunrise          string `json:"sunrise"`
		Sunset           string `json:"sunset"`
		TempDay          string `json:"tempDay"`
		TempNight        string `json:"tempNight"`
		Updatetime       string `json:"updatetime"`
		Uvi              string `json:"uvi"`
		WindDegreesDay   string `json:"windDegreesDay"`
		WindDegreesNight string `json:"windDegreesNight"`
		WindDirDay       string `json:"windDirDay"`
		WindDirNight     string `json:"windDirNight"`
		WindLevelDay     string `json:"windLevelDay"`
		WindLevelNight   string `json:"windLevelNight"`
		WindSpeedDay     string `json:"windSpeedDay"`
		WindSpeedNight   string `json:"windSpeedNight"`
	} `json:"forecast"`
	Hourly []struct {
		Condition   string `json:"condition"`
		ConditionID string `json:"conditionId"`
		Date        string `json:"date"`
		Hour        string `json:"hour"`
		Humidity    string `json:"humidity"`
		IconDay     string `json:"iconDay"`
		IconNight   string `json:"iconNight"`
		Pop         string `json:"pop"`
		Pressure    string `json:"pressure"`
		Qpf         string `json:"qpf"`
		RealFeel    string `json:"realFeel"`
		Snow        string `json:"snow"`
		Temp        string `json:"temp"`
		Updatetime  string `json:"updatetime"`
		Uvi         string `json:"uvi"`
		WindDegrees string `json:"windDegrees"`
		WindDir     string `json:"windDir"`
		WindSpeed   string `json:"windSpeed"`
		Windlevel   string `json:"windlevel"`
	} `json:"hourly"`
	LiveIndexData map[string][]LiveIndex `json:"liveIndex"`
	Sfc           struct {
		Banner   string `json:"banner"`
		NearRain int64  `json:"nearRain"`
		Notice   string `json:"notice"`
		Percent  []struct {
			Dbz     int64   `json:"dbz"`
			Desc    string  `json:"desc"`
			Icon    int64   `json:"icon"`
			Percent float32 `json:"percent"`
		} `json:"percent"`
		Rain         int64 `json:"rain"`
		RainLastTime int64 `json:"rainLastTime"`
		SfCondition  int64 `json:"sfCondition"`
		Timestamp    int64 `json:"timestamp"`
	} `json:"sfc"`
}

type LiveIndex struct {
	Code       int64  `json:"code"`
	Day        string `json:"day"`
	Desc       string `json:"desc"`
	Level      string `json:"level"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Updatetime string `json:"updatetime"`
}
type WeatherLiveResponse struct {
	Code   int64 `json:"code"`
	Result struct {
		Code int64             `json:"code"`
		Data *WeatherLive1Data `json:"data"`
		Msg  string            `json:"msg"`
		Rc   struct {
			C int64  `json:"c"`
			P string `json:"p"`
		} `json:"rc"`
	} `json:"result"`
}

type WeatherStation struct {
	CNAME     string
	V06001    string
	V05001    string
	StationID string
}

type StationResponse struct {
	Code   int             `json:"code"`
	Result *WeatherStation `json:"result"`
}

type WeatherLive0Response struct {
	Code   int64 `json:"code"`
	Result struct {
		Ds            []*WeatherLive0Data `json:"DS"`
		ColCount      string              `json:"colCount"`
		FieldNames    string              `json:"fieldNames"`
		FieldUnits    string              `json:"fieldUnits"`
		RequestTime   string              `json:"requestTime"`
		ResponseTime  string              `json:"responseTime"`
		ReturnCode    int64               `json:"returnCode"`
		ReturnMessage string              `json:"returnMessage"`
		RowCount      string              `json:"rowCount"`
		TakeTime      string              `json:"takeTime"`
	} `json:"result"`
}

type WeatherLive0Data struct {
	Datetime string `json:"DATETIME"`
	// 纬度
	Lat string `json:"LAT"`
	// 经度
	Lon string `json:"LON"`
	// 12小时累计降水(10分钟)
	PRE12H string `json:"PRE_12H"`
	// 1小时累计降水(10分钟)
	PRE1H string `json:"PRE_1H"`
	// 24小时累计降水(10分钟)
	PRE24H string `json:"PRE_24H"`
	// 3小时累计降水(10分钟)
	PRE3H string `json:"PRE_3H"`
	// 6小时累计降水(10分钟)
	PRE6H string `json:"PRE_6H"`
	// 相对湿度
	Rhu string `json:"RHU"`
	// 总云量
	Tcdc string `json:"TCDC"`
	// 温度
	Tem string `json:"TEM"`
	// 能见度
	Vis string `json:"VIS"`
	// 天气现象
	Wea string `json:"WEA"`
	// 风向
	Wind string `json:"WIND"`
	// 风力
	Wins string `json:"WINS"`
}

func (s *WeatherService) GetWeatherByCity(city string) (*WeatherLive0Data, error) {
	log.Info().Msgf("GetWeatherByCity: %s", city)
	// 深圳市 -> 深圳
	index := strings.Index(city, "市")
	if index > 0 {
		city = city[0:index]
	}
	stat, err := s.getStationIDByCity(city)
	if err != nil {
		return nil, err
	}
	// http://data.cma.cn/kbweb/home/live
	param := make(map[string]string)
	param["lat"] = stat.V05001
	param["lon"] = stat.V06001
	param["type"] = "0"
	paramBtArr, _ := json.Marshal(param)
	req, _ := http.NewRequest("POST", "http://data.cma.cn/kbweb/home/live",
		bytes.NewBuffer(paramBtArr))
	req.Header.Set("Origin", "http://data.cma.cn")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", defination.USER_AGENT)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Msgf("GetWeatherByCity err: %+v", err)
		return nil, err
	}
	defer res.Body.Close()
	btArr, _ := ioutil.ReadAll(res.Body)
	resObj := &WeatherLive0Response{}
	err = json.Unmarshal(btArr, resObj)
	if err != nil {
		log.Error().Msgf("unmarshal res err: %+v", err)
		return nil, err
	}
	if resObj.Code != 0 {
		err = fmt.Errorf("GetWeatherByCity failed: %d", resObj.Code)
		log.Err(err)
		return nil, err
	}
	if len(resObj.Result.Ds) == 0 {
		err = fmt.Errorf("GetWeatherByCity data rs is empty")
		log.Err(err)
		return nil, err
	}
	return resObj.Result.Ds[0], nil
}

func (s *WeatherService) getStationIDByCity(city string) (*WeatherStation, error) {
	param := make(map[string]string)
	param["city"] = city
	paramBtArr, _ := json.Marshal(param)
	log.Info().Msgf("getStationIDByCity, original req: %+s", string(paramBtArr))
	req, _ := http.NewRequest("POST", "http://data.cma.cn/kbweb/home/getStationID",
		bytes.NewBuffer(paramBtArr))
	req.Header.Set("Origin", "http://data.cma.cn")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", defination.USER_AGENT)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Msgf("getStationID err: %+v", err)
		return nil, err
	}
	defer res.Body.Close()
	btArr, _ := ioutil.ReadAll(res.Body)
	log.Info().Msgf("getStationIDByCity, original res: %+s", string(btArr))
	resObj := &StationResponse{}
	err = json.Unmarshal(btArr, resObj)
	if err != nil {
		log.Error().Msgf("getStationID unmarshal res err: %+v", err)
		return nil, err
	}
	if resObj.Code != 0 {
		err = fmt.Errorf("getStationID failed: %d", resObj.Code)
		log.Err(err)
		return nil, err
	}
	return resObj.Result, nil
}

func (s *WeatherService) LocalWeather(req *agentapi.WeatherRequest, res *agentapi.WeatherResponse) {
	log.Info().Msg("localWeather start")
	loc, err := DefaultLocationService.GetMyLocation()
	if err != nil {
		log.Error().Msgf("GetMyLocation err: %+v", err)
		res.Response.Code = common.ErrorCode_EXTERNAL_ERROR
		res.Response.Msg = "GetMyLocation err"
		return
	}
	log.Info().Msgf("loc: %+v", loc)
	weather, err := DefaultWeatherService.GetWeatherByCity(loc.City)
	if err != nil {
		log.Error().Msgf("GetWeatherByCity err: %+v", err)
		res.Response.Code = common.ErrorCode_EXTERNAL_ERROR
		res.Response.Msg = "GetWeatherByCity err"
		return
	}
	log.Info().Msgf("weather: %+v", weather)
	windDirection, err := strconv.Atoi(weather.Wind)
	direction := s.getWindDirection(windDirection)
	buf := make([]byte, 0, 16)
	buf = append(buf, loc.City...)
	buf = append(buf, ","...)
	buf = append(buf, fmt.Sprintf("当前气温%s度，", weather.Tem)...)
	buf = append(buf, fmt.Sprintf("相对湿度百分之%s。", weather.Rhu)...)
	buf = append(buf, fmt.Sprintf("风向%s风%s级，", direction, weather.Wins)...)
	buf = append(buf, fmt.Sprintf("能见度%s米，", weather.Vis)...)
	buf = append(buf, fmt.Sprintf("未来12小时预计降水量%s毫米。", weather.PRE12H)...)
	text := string(buf)
	log.Info().Msg(text)
	ttsReq := &agentapi.TtsRequest{
		Text: text,
		Request: &common.BaseRequest{
			RequestID: req.Request.RequestID,
		},
	}
	ttsRes := &agentapi.TtsResponse{
		Response: &common.BaseResponse{},
	}
	go DefaultTtsService.TextToVoice(ttsReq, ttsRes)
	res.Response = ttsRes.Response
}

func (s *WeatherService) getWindDirection(windDirection int) string {
	direction := ""
	if windDirection == 0 {
		direction = "北"
	} else if windDirection > 0 && windDirection < 90 {
		direction = "东北"
	} else if windDirection == 90 {
		direction = "东"
	} else if windDirection > 90 && windDirection < 180 {
		direction = "东南"
	} else if windDirection == 180 {
		direction = "南"
	} else if windDirection > 180 && windDirection < 270 {
		direction = "西南"
	} else if windDirection == 270 {
		direction = "西"
	} else if windDirection > 270 {
		direction = "西北"
	}
	return direction
}
