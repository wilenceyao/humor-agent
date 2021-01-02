package location

import (
	"errors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strings"
)

type Location struct {
	ExternalIP string
	Country    string
	Province   string
	City       string
}

func GetMyLocation() (*Location, error) {
	res, err := http.Get("https://ip.tool.lu/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	btArr, _ := ioutil.ReadAll(res.Body)
	s := string(btArr)
	log.Info().Msgf("ip location original res: %s", s)
	strArr := strings.Split(s, "\n")
	// 当前IP: xx.xx.xx.xx
	// 归属地: 中国 xxx xxx

	loc := &Location{}
	// handle IP
	loc.ExternalIP = strings.TrimSpace(strings.Split(strArr[0], ":")[1])
	// handle city
	if len(strArr) >= 2 {
		subStrArr := strings.Split(strings.TrimSpace(
			strings.Split(strArr[1], ":")[1]), " ")
		if len(subStrArr) >= 3 {
			loc.Country = subStrArr[0]
			loc.Province = subStrArr[1]
			loc.City = subStrArr[2]
		} else {
			err = errors.New("invalid original location res")
			log.Err(err)
			return nil, err
		}
	}
	return loc, nil
}
