package readers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tyba/srcd-domain/container"
	"github.com/tyba/srcd-domain/models"
	"github.com/tyba/srcd-domain/models/social"
	"github.com/tyba/srcd-rovers/client"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrPartialResponse      = errors.New("received partial data")
)

const (
	augurInsightsURL = "https://api.augur.io/v2/insights"
	augurKey         = "2bwn2e88g9dbva8pjolgxeid0nz9m4ne"
	augurRateLimit   = 1 * time.Second
)

type AugurInsightsAPI struct {
	client       *client.Client
	next         time.Time
	insightStore *social.AugurInsightStore
}

func NewAugurInsightsAPI(client *client.Client) *AugurInsightsAPI {
	return &AugurInsightsAPI{
		client:       client,
		next:         time.Now(),
		insightStore: container.GetDomainModelsSocialAugurInsightStore(),
	}
}

func (a *AugurInsightsAPI) SearchByEmail(email string) (*social.AugurInsight, *http.Response, error) {
	if time.Now().Before(a.next) {
		time.Sleep(time.Now().Sub(a.next))
	}
	a.next = time.Now().Add(augurRateLimit)

	q := &url.Values{}
	q.Add("email", email)

	var (
		i = a.insightStore.New()
		r RawInsight
	)

	res, err := a.doRequest(q, &r)
	if err != nil {
		return nil, res, err
	}

	i.Demographics.Gender = getFirstValue(r.Demographics.Gender)
	i.Demographics.Language = getFirstValue(r.Demographics.Language)
	i.Geographics.Locale = getFirstValue(r.Geographics.Locale)
	i.Geographics.Location = getFirstValue(r.Geographics.Location)
	i.Private.Bio = getFirstValue(r.Private.Bio)
	i.Private.Description = getFirstValue(r.Private.Description)
	i.Private.Email = getFirstValue(r.Private.Email)
	i.Private.Homepage = getFirstValue(r.Private.Homepage)
	i.Private.Name = getFirstValue(r.Private.Name)
	i.Private.Phone = getFirstValue(r.Private.Phone)
	i.Private.Photo = getFirstValue(r.Private.Photo)
	i.Profiles.Handle = getFirstValue(r.Profiles.Handle)
	i.Profiles.Post = getFirstValue(r.Profiles.Post)
	i.Profiles.Service = getFirstValue(r.Profiles.Service)
	i.Profiles.URL = getFirstValue(r.Profiles.URL)
	i.Psychographics.Color = getFirstValue(r.Psychographics.Color)
	i.Psychographics.Keyword = getFirstValue(r.Psychographics.Keyword)
	i.Misc = r.Misc
	i.Status = r.Status

	return i, res, nil
}

func getFirstValue(values []Value) string {
	if len(values) == 0 {
		return ""
	}
	return values[0].Value
}

func (a *AugurInsightsAPI) buildURL(q *url.Values) *url.URL {
	q.Add("key", augurKey)

	url, _ := url.Parse(augurInsightsURL)
	url.RawQuery = q.Encode()

	return url
}

func (a *AugurInsightsAPI) doRequest(q *url.Values, result interface{}) (*http.Response, error) {
	req, err := client.NewRequest(a.buildURL(q).String())
	if err != nil {
		return nil, err
	}

	res, err := a.client.DoJSON(req, result)
	if err != nil {
		return res, err
	}

	switch res.StatusCode {
	case 200:
		return res, nil
	case 202:
		return res, ErrPartialResponse
	default:
		return res, ErrUnexpectedStatusCode
	}
}

type RawInsight struct {
	Demographics struct {
		Gender   []Value `json:"gender"`
		Language []Value `json:"language"`
	} `json:"DEMOGRAPHICS"`
	Geographics struct {
		Locale   []Value `json:"locale"`
		Location []Value `json:"location"`
	} `json:"GEOGRAPHICS"`
	Private struct {
		Bio         []Value `json:"bio"`
		Description []Value `json:"description"`
		Email       []Value `json:"email"`
		Homepage    []Value `json:"homepage"`
		Name        []Value `json:"name"`
		Phone       []Value `json:"phone"`
		Photo       []Value `json:"photo"`
	} `json:"PRIVATE"`
	Profiles struct {
		Handle  []Value `json:"handle"`
		Post    []Value `json:"post"`
		Service []Value `json:"service"`
		URL     []Value `json:"url"`
	} `json:"PROFILES"`
	Psychographics struct {
		Color   []Value `json:"color"`
		Keyword []Value `json:"keyword"`
	} `json:"PSYCHOGRAPHICS"`
	Misc struct {
		BackgroundImage       []Value `json:"background_image"`
		ColorBackground       []Value `json:"color_background"`
		ColorForeground       []Value `json:"color_foreground"`
		FacebookHandle        []Value `json:"facebook_handle"`
		FacebookID            []Value `json:"facebook_id"`
		FacebookURL           []Value `json:"facebook_url"`
		FirstName             []Value `json:"first_name"`
		KloutHandle           []Value `json:"klout_handle"`
		KloutURL              []Value `json:"klout_url"`
		LastName              []Value `json:"last_name"`
		LinkedinHandle        []Value `json:"linkedin_handle"`
		LinkedinURL           []Value `json:"linkedin_url"`
		LocationCity          []Value `json:"location_city"`
		LocationCountry       []Value `json:"location_country"`
		LocationFormatted     []Value `json:"location_formatted"`
		LocationState         []Value `json:"location_state"`
		MentionName           []Value `json:"mention_name"`
		MentionTwitterHandle  []Value `json:"mention_twitter_handle"`
		MentionTwitterID      []Value `json:"mention_twitter_id"`
		PersonUid             string  `json:"person_uid"`
		PostHashtag           []Value `json:"post_hashtag"`
		PostLink              []Value `json:"post_link"`
		PostPhoto             []Value `json:"post_photo"`
		RecentRetweetedCount  []Value `json:"recent_retweeted_count"`
		ReplyToTwitterHandle  []Value `json:"reply_to_twitter_handle"`
		ReplyToTwitterID      []Value `json:"reply_to_twitter_id"`
		RepostedCount         []Value `json:"reposted_count"`
		SingleFollowingCount  []Value `json:"single_following_count"`
		SingleFolowersCount   []Value `json:"single_folowers_count"`
		SinglePostsCount      []Value `json:"single_posts_count"`
		Timestamp             float64 `json:"timestamp"` // 1378273006.96988
		TimeZone              []Value `json:"time_zone"`
		TweetPostSource       []Value `json:"tweet_post_source"`
		TwitterFollowingCount []Value `json:"twitter_following_count"`
		TwitterFolowersCount  []Value `json:"twitter_folowers_count"`
		TwitterHandle         []Value `json:"twitter_handle"`
		TwitterID             []Value `json:"twitter_id"`
		TwitterListedCount    []Value `json:"twitter_listed_count"`
		TwitterPostsCount     []Value `json:"twitter_posts_count"`
		TwitterURL            []Value `json:"twitter_url"`
		UrlDomain             []Value `json:"url_domain"`
	} `json:"MISC"`
	Status int `json:"status"`
}

type Value struct {
	Score json.Number `json:"score"`
	Value string      `json:"value"`
}

type AugurEmailSource interface {
	Next() bool
	Get() (string, error)
}

type AugurPeopleSource struct {
	results *models.PersonResultSet
	emails  []string
}

func NewAugurPeopleSource() *AugurPeopleSource {
	store := container.GetDomainModelsPersonStore()
	q := store.Query()
	return &AugurPeopleSource{
		results: store.MustFind(q),
	}
}

func (s *AugurPeopleSource) Next() bool {
	return s.results.Next()
}

func (s *AugurPeopleSource) Get() (string, error) {
	if len(s.emails) > 0 {
		email := s.emails[0]
		s.emails = s.emails[1:]
		return email, nil
	}
	person, err := s.results.Get()
	if err != nil {
		return "", err
	}
	var emails []string
	for _, email := range person.Email {
		emails = append(emails, email)
	}
	email := emails[0]
	s.emails = emails[1:]
	return email, nil
}

type AugurFileSource struct {
	scanner *bufio.Scanner
}

func NewAugurFileSource(filename string) *AugurFileSource {
	f, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("couldn't open %q - error: %s", filename, err))
	}
	return &AugurFileSource{
		scanner: bufio.NewScanner(f),
	}
}

func (s *AugurFileSource) Next() bool {
	return s.scanner.Scan()
}

func (s *AugurFileSource) Get() (string, error) {
	return s.scanner.Text(), s.scanner.Err()
}
