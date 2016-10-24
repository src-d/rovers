package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	apiUrl = "https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&start=%d"
	q      = "(%22generated+by+cgit%22+OR+%22powered+by+cgit%22)"

	domainGlobal      = "global"
	reasonInvalid     = "invalid"
	domainUsageLimits = "usageLimits"
	reasonKeyInvalid  = "keyInvalid"
)

var errPageNotFound error = errors.New("Page not found")
var errInvalidKey error = errors.New("Invalid key")

type googleCseApi struct {
	cachedPages map[int]*result
	key         string
	cx          string
	client      *http.Client
}

func newGoogleCseApi(key string, cx string) *googleCseApi {
	return &googleCseApi{
		cachedPages: make(map[int]*result),
		key:         key,
		cx:          cx,
		client:      &http.Client{},
	}
}

func (gca *googleCseApi) GetPage(index int) (*result, error) {
	result, ok := gca.cachedPages[index]
	if ok {
		return result, nil
	}

	result, err := gca.requestData(index)
	if err != nil {
		return nil, err
	}
	gca.cachedPages[index] = result

	return result, nil
}

func (gca *googleCseApi) PageExists(index int) (bool, error) {
	_, err := gca.GetPage(index)
	switch err {
	case errPageNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

func (gca *googleCseApi) Reset() {
	gca.cachedPages = make(map[int]*result)
}

func (gca *googleCseApi) requestData(index int) (*result, error) {
	searchUrl := fmt.Sprintf(apiUrl, gca.key, gca.cx, q, index)
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := gca.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		record, err := gca.result(resp.Body)
		if err != nil {
			return nil, err
		}
		return record, nil
	case http.StatusBadRequest:
		return nil, gca.handleBadRequestError(resp.Body)
	default:
		return nil, fmt.Errorf("Unhandled CSE API error: %s", resp.Status)
	}
}

func (gca *googleCseApi) handleBadRequestError(body io.ReadCloser) error {
	br, err := gca.badRequestResult(body)
	if err != nil {
		return err
	}

	if len(br.Error.Errors) > 0 {
		e := br.Error.Errors[0]
		if e.Domain == domainGlobal && e.Reason == reasonInvalid {
			// This case is for ANY bad param, page index included.
			// We must assume that all the data provided by the user (cx) are ok
			return errPageNotFound
		}

		if e.Domain == domainUsageLimits && e.Reason == reasonKeyInvalid {
			return errInvalidKey
		}
	}

	return fmt.Errorf("Bad request CSE API error: %s", br.Error.Message)
}

func (gca *googleCseApi) badRequestResult(body io.ReadCloser) (*badRequest, error) {
	var record badRequest
	if err := json.NewDecoder(body).Decode(&record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (gca *googleCseApi) result(body io.ReadCloser) (*result, error) {
	var record result
	if err := json.NewDecoder(body).Decode(&record); err != nil {
		return nil, err
	}

	return &record, nil
}
