// https://github.com/savaki/myfitnesspal/blob/master/login.go
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {

	// Create a cookiejar so we can keep the received cookies in subsequent HTTP calls
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	failOnErr(err, "Failed to create cookiejar")

	browser := &http.Client{Jar: jar}

	// Retrieve the login credentials from environment
	email := os.Getenv("PUSHER_EMAIL")
	password := os.Getenv("PUSHER_PASSWORD")

	// Succesful Login results in a redirect to the keys overview page
	overviewPage, err := login(email, password, browser)
	failOnErr(err, "Failed to create cookiejar")

	appIdFilter := ""
	//findId := "110326"

	appIds, err := extractAppIdsFromPage(browser, overviewPage)

	if appIdFilter == "" {
		// If no appId specified , list all keys
		for appId, _ := range appIds {
			getKeyPair(browser, appId)
		}
	} else {
		// Retrieve a specific keypair
		getKeyPair(browser, appIdFilter)
	}

}

func getKeyPair(browser *http.Client, appId string) error {

	var detailUrl = "https://app.pusher.com/apps/" + appId + "/api_access"

	var detailResponse *http.Response
	detailResponse, err := browser.Get(detailUrl)
	failOnErr(err, "error fetching keypair details")

	defer detailResponse.Body.Close()

	detailPage, err := goquery.NewDocumentFromReader(detailResponse.Body)
	failOnErr(err, "error parsing keypair details page")

	detailPage.Find(".name").Each(func(i int, s *goquery.Selection) {
		name := s.Find("h3").Text()
		fmt.Print(name)
	})

	detailPage.Find("code").Each(func(i int, s *goquery.Selection) {
		code := s.Text()

		if strings.HasPrefix(code, "key =") {

			pairRegex, _ := regexp.Compile(`key.=.'(\w+)'\nsecret.=.'(\w+)'`)

			result := pairRegex.FindStringSubmatch(code)
			fmt.Printf("|%s|%s|%s\n", appId, result[1], result[2])
		}
	})

	failOnErr(err, "could not find any valid keypairs for appId: "+string(appId))

	return nil

}

func login(email string, password string, client *http.Client) (*goquery.Document, error) {

	var login_url = "https://app.pusher.com/accounts/sign_in"

	var loginResponse *http.Response
	loginResponse, err := client.Get(login_url)
	defer loginResponse.Body.Close()

	loginPage, err := goquery.NewDocumentFromReader(loginResponse.Body)
	failOnErr(err, "Error parsing login Page")

	// Look for the form information. Especially the authenticity_token as it's used CSRF
	var authenticityToken string
	var utf8 string

	loginPage.Find("input[name=authenticity_token]").First().Each(func(i int, s *goquery.Selection) {
		authenticityToken, _ = s.Attr("value")
	})

	loginPage.Find("input[name=utf8]").First().Each(func(i int, s *goquery.Selection) {
		utf8, _ = s.Attr("value")
	})

	// Prepare for actual login
	var send_login_url = "https://app.pusher.com/accounts/sign_in"

	params := url.Values{}
	params.Set("utf8", utf8)
	params.Set("authenticity_token", authenticityToken)
	params.Set("account[email]", email)
	params.Set("account[password]", password)
	params.Set("remember_me", "1")

	sendLoginResponse, err := client.Post(send_login_url, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	failOnErr(err, "error sending login credentials")

	sendLoginPage, err := goquery.NewDocumentFromReader(sendLoginResponse.Body)

	defer sendLoginResponse.Body.Close()

	return sendLoginPage, err
}

/*
 * From the details page, extra
 */

func extractAppIdsFromPage(browser *http.Client, overviewPage *goquery.Document) (map[string]int, error) {
	var m = make(map[string]int)
	var err error

	overviewPage.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		name := s.Text()

		if strings.HasSuffix(href, "api_access") && !strings.HasPrefix(name, "App Keys") {

			re, err := regexp.Compile(`/apps/(\d+)/api_access`)

			failOnErr(err, "error compiing regex")

			result := re.FindStringSubmatch(href)

			pairId := result[1]
			m[pairId] = 1
		}
	})

	return m, err

}

/*
 * Helper page to exit on error with a nice message
 */
func failOnErr(err error, reason string) {
	if err != nil {
		log.Fatal("Failed: %s, %s\n\n", reason, err)
		os.Exit(-1)
	}

	return
}
