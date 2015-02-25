package sources

import (
	"strconv"
	"time"

	"github.com/tyba/opensource-search/sources/social/http"

	"github.com/PuerkitoBio/goquery"
)

const TwitterBaseURL = "https://twitter.com/%s"

type Twitter struct {
	client *http.Client
}

func NewTwitter(client *http.Client) *Twitter {
	return &Twitter{client}
}

func (g *Twitter) GetProfileByURL(url string) (*TwitterProfile, error) {
	req, err := http.NewRequest(url)
	if err != nil {
		return nil, err
	}

	doc, _, err := g.client.DoHTML(req)
	if err != nil {
		return nil, err
	}

	return NewTwitterProfile(url, doc), nil
}

type TwitterProfile struct {
	Created   time.Time
	Url       string
	Handle    string
	FullName  string
	Location  string
	Bio       string
	Web       string
	Tweets    int
	Followers int
	Following int
	Favorites int
}

func NewTwitterProfile(url string, doc *goquery.Document) *TwitterProfile {
	g := &TwitterProfile{Url: url, Created: time.Now()}
	g.fillBasicInfo(doc)
	g.fillStats(doc)

	return g
}

func (g *TwitterProfile) fillBasicInfo(doc *goquery.Document) {
	g.Handle = doc.Find(".ProfileHeaderCard-screenname span").Text()
	g.FullName = doc.Find(".ProfileHeaderCard-name a").Text()
	g.Location = doc.Find(".ProfileHeaderCard-locationText").Text()
	g.Bio = doc.Find(".ProfileHeaderCard-bio").Text()
	g.Web, _ = doc.Find(".ProfileHeaderCard-url a").Attr("title")
}

func (g *TwitterProfile) fillStats(doc *goquery.Document) {
	tweets := doc.Find("[data-nav='tweets'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(tweets); err == nil {
		g.Tweets = value
	}

	following := doc.Find("[data-nav='following'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(following); err == nil {
		g.Following = value
	}

	followers := doc.Find("[data-nav='followers'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(followers); err == nil {
		g.Followers = value
	}

	favorites := doc.Find("[data-nav='favorites'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(favorites); err == nil {
		g.Favorites = value
	}
}
