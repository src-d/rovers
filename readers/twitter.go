package readers

import (
	"strconv"
	"strings"

	"github.com/src-d/domain/container"
	"github.com/src-d/domain/models/social"
	"github.com/src-d/rovers/client"

	"github.com/PuerkitoBio/goquery"
)

const TwitterBaseURL = "https://twitter.com/%s"

type TwitterReader struct {
	client *client.Client
}

func NewTwitterReader(client *client.Client) *TwitterReader {
	return &TwitterReader{client}
}

func (t *TwitterReader) GetProfileByURL(url string) (*social.TwitterProfile, error) {
	req, err := client.NewRequest(url)
	if err != nil {
		return nil, err
	}

	doc, _, err := t.client.DoHTML(req)
	if err != nil {
		return nil, err
	}

	store := container.GetDomainModelsSocialTwitterProfileStore()
	profile, err := store.New(url)
	if err != nil {
		return nil, err
	}
	t.fillBasicInfo(doc, profile)
	t.fillStats(doc, profile)

	return profile, nil
}

func (t *TwitterReader) fillBasicInfo(doc *goquery.Document, p *social.TwitterProfile) {
	p.Handle = doc.Find(".ProfileHeaderCard-screenname span").Text()
	p.FullName = doc.Find(".ProfileHeaderCard-name a").Text()
	p.Location = strings.Trim(doc.Find(".ProfileHeaderCard-locationText").Text(), "\n\r\t ")
	p.Bio = doc.Find(".ProfileHeaderCard-bio").Text()
	p.Web, _ = doc.Find(".ProfileHeaderCard-url a").Attr("title")
}

func (t *TwitterReader) fillStats(doc *goquery.Document, p *social.TwitterProfile) {
	tweets := doc.Find("[data-nav='tweets'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(tweets); err == nil {
		p.Tweets = value
	}

	following := doc.Find("[data-nav='following'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(following); err == nil {
		p.Following = value
	}

	followers := doc.Find("[data-nav='followers'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(followers); err == nil {
		p.Followers = value
	}

	favorites := doc.Find("[data-nav='favorites'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(favorites); err == nil {
		p.Favorites = value
	}
}
