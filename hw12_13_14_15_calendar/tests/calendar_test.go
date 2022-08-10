package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages-go/v16"
)

type calandarTest struct {
	responseStatusCode int
	responseBody       []byte
}

func (test *calandarTest) iSendRequestTo(ctx context.Context, httpMethod, addr string) error {
	var req *http.Request
	var r *http.Response
	var err error

	switch httpMethod {
	case http.MethodGet, http.MethodDelete:
		req, err = http.NewRequestWithContext(ctx, httpMethod, addr, bytes.NewReader([]byte("")))
		if err != nil {
			return err
		}
		r, err = http.DefaultClient.Do(req)
		defer r.Body.Close()
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return err
	}

	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)
	return err
}

func (test *calandarTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *calandarTest) theResponseCodeShouldNotBe(code int) error {
	if test.responseStatusCode == code {
		return fmt.Errorf("unexpected status code: %d = %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *calandarTest) iSendRequestToWithData(
	ctx context.Context,
	httpMethod, addr, contentType string,
	data *messages.PickleDocString,
) error {
	var req *http.Request
	var r *http.Response
	var err error

	replacer := strings.NewReplacer("\n", "", "\t", "")
	cleanJSON := replacer.Replace(data.Content)

	switch httpMethod {
	case http.MethodPost, http.MethodPut:
		req, err = http.NewRequestWithContext(ctx, httpMethod, addr, bytes.NewReader([]byte(cleanJSON)))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", contentType)
		r, err = http.DefaultClient.Do(req)
		defer r.Body.Close()
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return err
	}

	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)
	return err
}

func (test *calandarTest) theResponseShouldMatchText(data *messages.PickleDocString) error {
	replacer := strings.NewReplacer("\n", "", "\t", "")
	cleanDataContent := replacer.Replace(data.Content)
	if string(test.responseBody) != cleanDataContent {
		return fmt.Errorf("unexpected text: %s != %s", test.responseBody, cleanDataContent)
	}
	return nil
}

func (test *calandarTest) iWaitAndSendRequestTo(ctx context.Context, delay, httpMethod, addr string) error {
	waitInterval, err := time.ParseDuration(delay)
	if err != nil {
		return err
	}

	log.Printf("wait %s before send request", waitInterval)
	time.Sleep(waitInterval)

	return test.iSendRequestTo(ctx, httpMethod, addr)
}

func (test *calandarTest) theResponseShouldNotContainText(data *messages.PickleDocString) error {
	replacer := strings.NewReplacer("\n", "", "\t", "")
	cleanDataContent := replacer.Replace(data.Content)
	if strings.Contains(string(test.responseBody), cleanDataContent) {
		return fmt.Errorf("unexpected text: %s contains %s", test.responseBody, cleanDataContent)
	}
	return nil
}

func InitializeScenario(s *godog.ScenarioContext) {
	test := new(calandarTest)

	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The response code should not be (\d+)$`, test.theResponseCodeShouldNotBe)
	s.Step(`^The response should match text:$`, test.theResponseShouldMatchText)
	s.Step(`^I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^I wait "([^"]*)" and send "([^"]*)" request to "([^"]*)"$`, test.iWaitAndSendRequestTo)
	s.Step(`^The response should not contain text:$`, test.theResponseShouldNotContainText)
}
