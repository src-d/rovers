package linkedin

import (
	"github.com/tyba/srcd-rovers/client"

	. "gopkg.in/check.v1"
)

const (
	TybaCompanyId = 924688
	CookieFixture = `bcookie="v=2&21733669-d772-42c2-8d0c-6631f3b494b7"; bscookie="v=1&20150904164640a97c0766-597a-4097-83e6-aa8ddf3ce0f6AQEBd5soPhRmh7HBB1afmID0u3OdZN6_"; visit="v=1&M"; sessionid="eyJkamFuZ29fdGltZXpvbmUiOiJFdXJvcGUvQmVybGluIn0:1ZbNU6:6Qlqhdr9bWBDhQNGA567O14PNrY"; csrftoken=3T0lfgzrejTNI80YLJduCviqdTtqYxLx; __utma=23068709.221348018.1441385201.1442832110.1442832110.1; __utmz=23068709.1442832110.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utmv=23068709.guest; L1c=38153dc8; wutan=5XW4mxShwpmVW0IhvylSv8h5tXLFDtAg9qc5/+mAQoI=; L1e=1e953b94; li_at=AQEDAQD4rTcDfJngAAABT8gc3_wAAAFP-euXSk4ArFmrYrEtDQ5abVZGsQQ9cLzu1Htm0WZ22vJ5bhJRMpd1-o9a6FET44xN5vG90I5Mst_NtpsPnGUyNPcJR-O3sT15_SXimb7ObEwtDpqiISBkAjCz; liap=true; sl="v=1&x_gci"; JSESSIONID="ajax:7765494870178832670"; oz_props_fetch_size1_16297271=6; share_setting=PUBLIC; sdsc=1%3A1SZM1shxDNbLt36wZwCgPgvN58iw%3D; lidc="b=TB71:g=114:u=115:i=1443005316:t=1443074927:s=AQE_9GpXb_1hvkwazhni2mcwgIO3bXA4"; _ga=GA1.2.221348018.1441385201; _gat=1; RT=s=1443005316840&r=https%3A%2F%2Fwww.linkedin.com%2Fcompany%2F924688; _lipt=0_0DWUvqOlLiwrAUp1_qzuTvYKhO30OcJ9TEc3PczGd9dDgjcZ3KAXmhKA8eI6zryVYkmWcw-jFWWKI9Y1axh_16jv7p-SSo-G8o3kNVxDF5uQ_pKPUcwECfQt5cKUp1RtrgPEwPNCNvfZj_EL8mSAkNMgC_n_MQ_djS9R8jd-4DNXjh6uHbexZL3ZMiyEwiWXVvkjSDpXz2ZY9mPD022h5J6eQtIt313hKiRyrB3sMn_T6I0bxXRmS3Ob-Q3TW_JFl3pU9euQDHhmsn3eaOVcHf; lang="v=2&lang=en-us"`
)

func (s *S) TestLinkedIn_GetEmployees(c *C) {
	cli := client.NewClient(true)
	wc := NewLinkedInWebCrawler(cli, CookieFixture)
	_, err := wc.GetEmployees(TybaCompanyId)
	c.Assert(err, IsNil)
	// c.Assert(len(employees), Equals, 46)
	// c.Assert(err, Equals, client.NotFound)
}
