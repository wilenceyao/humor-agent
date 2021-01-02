package cmasdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/internal/defination"
	"io/ioutil"
	"net/http"
	"strings"
)

type WeatherLiveData struct {
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
		Code int64            `json:"code"`
		Data *WeatherLiveData `json:"data"`
		Msg  string           `json:"msg"`
		Rc   struct {
			C int64  `json:"c"`
			P string `json:"p"`
		} `json:"rc"`
	} `json:"result"`
}

type Station struct {
	CNAME     string
	V06001    string
	V05001    string
	StationID string
}

type StationResponse struct {
	Code   int      `json:"code"`
	Result *Station `json:"result"`
}

func GetWeatherByCity(city string) (*WeatherLiveData, error) {
	log.Info().Msgf("GetWeatherByCity: %s", city)
	// 深圳市 -> 深圳
	index := strings.Index(city, "市")
	if index > 0 {
		city = city[0:index]
	}
	stat, err := getStationIDByCity(city)
	if err != nil {
		return nil, err
	}
	// http://data.cma.cn/kbweb/home/live
	param := make(map[string]string)
	param["lat"] = stat.V05001
	param["lon"] = stat.V06001
	param["type"] = "1"
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
	resObj := &WeatherLiveResponse{}
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
	return resObj.Result.Data, nil
}

func getStationIDByCity(city string) (*Station, error) {
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
