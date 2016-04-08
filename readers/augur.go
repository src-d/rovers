package readers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/src-d/rovers/client"
	"github.com/src-d/rovers/metrics"
	"gop.kg/src-d/domain@v5/container"
	"gop.kg/src-d/domain@v5/models/social"
)

const (
	AugurInsightsURL = "https://api.augur.io/v1/insights"
	AugurKey         = "2bwn2e88g9dbva8pjolgxeid0nz9m4ne"
)

// AugurInsightsAPI works as a "fire and forget" service, it'll run as fast as
// it can for as long as it can until your API token reaches its monthly rate
// limit.
type AugurInsightsAPI struct {
	client       *client.Client
	insightStore *social.AugurInsightStore
	reachedLimit bool
}

func NewAugurInsightsAPI(client *client.Client) *AugurInsightsAPI {
	return &AugurInsightsAPI{
		client:       client,
		insightStore: container.GetDomainModelsSocialAugurInsightStore(),
	}
}

func (a *AugurInsightsAPI) SearchByEmail(email string) (*social.AugurInsight, *http.Response, error) {
	metrics.AugurRequested.Inc()

	q := &url.Values{}
	q.Add("email", email)

	body, res, err := a.doRequest(q)
	if err == ErrRateLimitExceeded {
		a.reachedLimit = true
		metrics.AugurFailed.WithLabelValues("api_rate_limit").Inc()
		return nil, res, err
	}
	if err != nil {
		metrics.AugurFailed.WithLabelValues("request_err").Inc()
		return nil, res, err
	}

	insight, err := a.processResponse(body)
	if err != nil {
		metrics.AugurFailed.WithLabelValues("process_response_err").Inc()
		return nil, res, err
	}
	insight.InputEmail = email
	return insight, res, nil
}

func (a *AugurInsightsAPI) buildURL(q *url.Values) *url.URL {
	q.Add("key", AugurKey)

	u, _ := url.Parse(AugurInsightsURL)
	u.RawQuery = q.Encode()

	return u
}

func (a *AugurInsightsAPI) doRequest(q *url.Values) ([]byte, *http.Response, error) {

	req, err := client.NewRequest(a.buildURL(q).String())
	if err != nil {
		return nil, nil, err
	}

	start := time.Now()
	res, err := a.client.Do(req)
	if err != nil {
		return nil, res, err
	}
	defer res.Body.Close()

	defer func() {
		elapsed := time.Since(start)
		microseconds := float64(elapsed) / float64(time.Microsecond)
		metrics.AugurRequestDur.Observe(microseconds)
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res, err
	}

	switch res.StatusCode {
	case 200, 202:
		return body, res, nil
	case 420, 429:
		return body, res, ErrRateLimitExceeded
	default:
		return body, res, ErrUnexpectedStatusCode
	}
}

func (a *AugurInsightsAPI) processResponse(body []byte) (*social.AugurInsight, error) {
	doc := a.insightStore.New()

	var raw interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	var ri RawInsight
	if err := json.Unmarshal(body, &ri); err != nil {
		return nil, err
	}

	doc.Raw = raw

	doc.Email = getValues(ri.Email)
	doc.GithubURL = getValues(ri.GithubURL)
	doc.LinkedinURL = getValues(ri.LinkedinURL)
	doc.Location = getValues(ri.Location)
	doc.Name = getValues(ri.Name)
	doc.TwitterURL = getValues(ri.TwitterURL)

	doc.HasData = hasData(doc)
	doc.Last = time.Now()
	doc.LastStatus = ri.LastStatus
	doc.Done = true
	doc.TwitterDone = false
	doc.GitHubDone = false

	return doc, nil
}

func getValues(values []RawValue) []string {
	var s []string
	for _, value := range values {
		s = append(s, value.Value)
	}
	return s
}

func hasData(doc *social.AugurInsight) bool {
	if doc.LastStatus == 200 {
		return true
	}

	dataLengths := []int{
		len(doc.GithubURL),
		len(doc.LinkedinURL),
		len(doc.Location),
		len(doc.Name),
		len(doc.TwitterURL),
	}

	for _, length := range dataLengths {
		if length > 0 {
			return true
		}
	}
	return false
}

type RawInsight struct {
	Email       []RawValue `json:"email" bson:"email"`
	GithubURL   []RawValue `json:"github-url" bson:"github_url"`
	LinkedinURL []RawValue `json:"linkedin-url" bson:"linkedin_url"`
	Location    []RawValue `json:"location" bson:"location_formatted"`
	Name        []RawValue `json:"name" bson:"name"`
	TwitterURL  []RawValue `json:"twitter-url" bson:"twitter_url"`
	LastStatus  int        `json:"status" bson:"status"`
}

type RawValue struct {
	Score json.Number `json:"score"`
	Value string      `json:"value"`
}
