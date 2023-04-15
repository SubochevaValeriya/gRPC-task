package internal

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	gRPC_task "rusprofile/proto"
	"strconv"
)

type CompanyInfo struct {
	INN      int    `json:"INN"`
	KPP      int    `json:"KPP"`
	Name     string `json:"name"`
	Director string `json:"director"`
}

var notINNError = errors.New("incorrect format for INN")
var companyNotFoundError = errors.New("company not found")

func INNValidation(INN int64) (bool, error) {
	matched, err := regexp.MatchString(`^(\d{10})$`, fmt.Sprintf("%v", INN))
	if err != nil {
		return false, err
	}
	if !matched || len(fmt.Sprintf("%v", INN)) != 10 {
		return false, notINNError
	}

	return true, nil
}

func RusProfileParse(INN int64) (gRPC_task.CompanyInfo, error) {
	info := gRPC_task.CompanyInfo{}

	res, err := http.Get(fmt.Sprintf("https://www.rusprofile.ru/search?query=%v", INN))
	if err != nil {
		return info, companyNotFoundError
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return info, errors.New(res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return info, err
	}

	if doc.Find(".search-result__container").Find("p").Text() == "Попробуйте изменить поисковой запрос" {
		return info, companyNotFoundError
	}

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
