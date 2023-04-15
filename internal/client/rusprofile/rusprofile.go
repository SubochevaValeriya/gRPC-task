package rusprofile

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	gRPC_task "rusprofile/proto"
	"strconv"
	"time"
)

type Client interface {
	GetProfile(INN int64) (gRPC_task.CompanyInfo, error)
}

type client struct {
	httpClient http.Client
	addr       string
}

func NewClient(addr string, timeout time.Duration) Client {
	client := client{
		httpClient: http.Client{Timeout: timeout},
		addr:       addr,
	}
	return client
}

var notINNError = errors.New("incorrect format for INN")
var companyNotFoundError = errors.New("company not found")

func (c client) GetProfile(INN int64) (gRPC_task.CompanyInfo, error) {
	if !(INN >= 1000000000 && INN <= 9999999999) {
		return gRPC_task.CompanyInfo{}, notINNError
	}

	res, err := c.httpClient.Get(fmt.Sprintf("https://%v/search?query=%v", c.addr, INN))
	if err != nil {
		return gRPC_task.CompanyInfo{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return gRPC_task.CompanyInfo{}, errors.New(res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return gRPC_task.CompanyInfo{}, err
	}

	if doc.Find(".search-result__container").Find("p").Text() == "Попробуйте изменить поисковой запрос" {
		return gRPC_task.CompanyInfo{}, companyNotFoundError
	}

	info := gRPC_task.CompanyInfo{}
	info.INN = INN
	doc.Find(".copy_target").Each(func(i int, s *goquery.Selection) {
		if s.AttrOr("id", "") == "clip_kpp" {
			info.KPP, _ = strconv.ParseInt(s.Text(), 10, 64)
		}
	})

	info.Name = doc.Find(".company-name").Text()
	doc.Find(".company-info__text").Each(func(i int, s *goquery.Selection) {
		if s.Find("a").AttrOr("data-goal-param", "") == "interactions, person_ul" {
			info.Director = s.Text()
		}
	})

	return info, nil
}
