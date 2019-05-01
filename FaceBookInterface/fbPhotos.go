package fb

import (
	"fmt"
	"strings"

	"github.com/A-M-Simmons/FacebookScrape/datastructs"
	"github.com/PuerkitoBio/goquery"
)

func testIfCommentThread(st *goquery.Selection) (isCommentThread bool, commentThreadSingle bool, commentThreadHref string) {
	isCommentThread = false
	commentThreadSingle = false
	commentThreadHref = ""
	// Iterate over every 'a' tag to see if its a link for a Comment thread
	st.Find("a").EachWithBreak(func(_ int, stt *goquery.Selection) bool {
		href, _ := stt.Attr("href") //TODO: Existence handle
		msg, _ := stt.Html()        //TODO: Error handle
		if (strings.Contains(href, "/Comment/replies/") || strings.Contains(href, "/comment/replies/")) && msg != "Reply" && (strings.Contains(msg, "reply") || strings.Contains(msg, "replies")) {
			if strings.Contains(msg, "reply") {
				commentThreadSingle = true
			}
			commentThreadHref = href
			isCommentThread = true
			return false
		}
		return true
	})
	return
}

func (s *Session) getComments(p *datastructs.Photo) {
	var commentThreads []datastructs.Comment

	ctSel := s.browser.Find("div#MPhotoContent")
	ctSel = ctSel.Find("div")
	ctSel = ctSel.Next() // Skip over Uploader details
	ctSel = ctSel.Find("div")
	ctSel = ctSel.Find("div")
	ctSel = ctSel.Find("div")
	ctSel = ctSel.Next()
	ctSel = ctSel.Next()
	ctSel = ctSel.Find("div")

	// Iterate over every Comment of Comment thread
	ctSel = ctSel.Find("div").Has("h3").EachWithBreak(func(_ int, st *goquery.Selection) bool {

		isCommentThread, commentThreadSingle, commentThreadHref := testIfCommentThread(st)
		// Get comments from Comment thread page
		if isCommentThread {
			if commentThreadSingle {
				CommentThread := s.getSingleReplyCommentThread(commentThreadHref)
				commentThreads = append(commentThreads, CommentThread)
			} else {
				CommentThread := s.getCommentThread(commentThreadHref)
				commentThreads = append(commentThreads, CommentThread)
			}
		} else {
			CommentThread := getCommentDetails(st, true)
			commentThreads = append(commentThreads, CommentThread)
		}
		return true
	})

	// Add Comments to Photo Structure
	p.Comments = commentThreads
}

func (s *Session) getSingleReplyCommentThread(url string) datastructs.Comment {
	var CommentThread datastructs.Comment

	url = fmt.Sprintf("https://m.facebook.com%s", url)
	s.Open(url) //TODO: Error handle
	ctSel := s.browser.Find("div#root")
	ctSel = ctSel.Find("div")
	ctSel = ctSel.Children().Has("h3")

	// Get Root message
	com := getCommentDetails(ctSel, false)
	CommentThread.Append(com)

	// Get Single response message
	ctSel = ctSel.Next()
	com = getCommentDetails(ctSel.Find("div"), false)
	CommentThread.Append(com)

	return CommentThread
}

func (s *Session) getCommentThread(url string) datastructs.Comment {
	var CommentThread datastructs.Comment
	url = fmt.Sprintf("https://m.facebook.com%s", url)
	s.Open(url) //TODO: Error handle
	ctSel := s.browser.Find("div#root")
	ctSel = ctSel.Find("div")
	ctSel = ctSel.Children()
	ctSel = ctSel.Next()

	// Get root message
	CommentThread = getCommentDetails(ctSel, false)

	// Get child messages
	ctSel = ctSel.Next()
	ctSel = ctSel.Children().Has("h3")
	for true {
		html, _ := ctSel.Html()
		if html == "" {
			break
		}
		com := getCommentDetails(ctSel, false)
		CommentThread.Append(com)
		ctSel = ctSel.Next()
	}
	return CommentThread
}

func getCommentDetails(st *goquery.Selection, skipDiv bool) datastructs.Comment {
	var comment datastructs.Comment
	if skipDiv == false {
		st = st.Find("div")
	}

	// Get Real Name
	commenterRealName, _ := st.Find("a").Html() //TODO: Error handling

	// Get UserName
	commenterUserName, _ := st.Find("a").Attr("href")
	commenterUserName = strings.Split(strings.Split(commenterUserName, "?")[0], "/")[1]

	// Get Message
	message, _ := st.Find("div").Html() //TODO: Error handling

	// Get Time Stamp
	timeStamp, _ := st.Find("div").Find("abbr").Html() //TODO: Error handling

	// If no errors, save data
	comment.UserID.SetRealName(commenterRealName)
	comment.UserID.SetUsername(commenterUserName)
	comment.SetComment(message) // TODO: Error handle
	comment.SetTimestamp(timeStamp)

	return comment
}

func (s *Session) getPhotoStats(p *datastructs.Photo) {
	// Get URL for full sized image
	var url string
	urlSelector := s.browser.Find("div#MPhotoContent")
	urlSelector = urlSelector.Find("div")
	urlSelector = urlSelector.Find("div")
	urlSelector = urlSelector.Next()
	urlSelector = urlSelector.Find("span")
	urlSelector = urlSelector.Find("div")
	urlSelector = urlSelector.Find("div").Next()
	urlSelector = urlSelector.Find("a")
	for true {
		html, _ := urlSelector.Html()
		//fmt.Println(html)
		if html == "" {
			break
		} else if html == "View full size" || html == "View Full Size" {
			url, _ = urlSelector.Attr("href") // TODO: Add error handling
			break
		}
		urlSelector = urlSelector.Next()
	}
	p.URL = url

	// Get Username and Realname
	nameSelector := s.browser.Find("div#MPhotoContent")
	nameSelector = nameSelector.Find("div")
	nameSelector = nameSelector.Find("div")

	// UserName
	nameID, _ := nameSelector.Find("a").Attr("href") // TODO: Add Error Handling
	p.UserID.UserName = strings.Split(strings.Split(nameID, "?")[0], "/")[1]

	// Real name
	realNameSelector := nameSelector.Find("strong")
	realName, _ := realNameSelector.Html() // TODO: Add Error Handling
	p.UserID.RealName = realName           // TODO: Add constructor function rather than direct element access

	// Timestamp
	timeSelector := s.browser.Find("div#MPhotoContent")
	timeSelector = timeSelector.Find("div")
	timeSelector = timeSelector.Find("div")
	timeSelector = timeSelector.Next()
	timeSelector = timeSelector.Find("div")
	timeSelector = timeSelector.Find("div")
	timeSelector = timeSelector.Find("abbr")
	time, _ := timeSelector.Html() // TODO: Add error handling
	p.DateTime = time
}

// getPhotoPages ...
func (s *Session) getPhotoPage(url string) (datastructs.Photo, error) {
	var p datastructs.Photo
	url = fmt.Sprintf("https://m.facebook.com%s", url)
	err := s.Open(url)

	if err != nil {
		return p, err
	}
	p.URL = url
	s.getPhotoStats(&p)
	s.getComments(&p) // Fix: reference to this function
	return p, nil
}

// GetPhotos ...
func (s *Session) GetPhotos(id string, t string) error {
	var url string
	switch t {
	case "PhotosBy":
		url = fmt.Sprintf("https://m.facebook.com/search/%s/photos-by", id)
	case "PhotosLiked":
		url = fmt.Sprintf("https://m.facebook.com/search/%s/photos-liked", id)
	case "PhotosOfTagged":
		url = fmt.Sprintf("https://m.facebook.com/search/%s/photos-of", id)
	case "PhotosCommented":
		url = fmt.Sprintf("https://m.facebook.com/search/%s/photos-commented", id)
	case "PhotosInteracted":
		url = fmt.Sprintf("https://m.facebook.com/search/%s/photos-interacted", id)
	case "PhotosInterested":
		url = fmt.Sprintf("https://m.facebook.com/search/%s/photos-interested", id)
	case "PhotosRecommended":
		url = fmt.Sprintf("https://m.facebook.com/search/%s/photos-recommended-for", id)
	default:
		fmt.Println(fmt.Sprintf("Case %s not understood"))
		err := fmt.Errorf("%s Photo type not found", t)
		return err
	}
	return s.getPhotosInternal(url)
}

// getPhotosInternal ...
func (s *Session) getPhotosInternal(url string) error {
	var photos []datastructs.Photo
	err := s.Open(url)
	if err != nil {
		return err
	}

	// Iterate over all the pages of pictures
	var results []string
	for true {
		st := s.browser.Find("div#BrowseResultsContainer")
		st = st.Find("div")
		st = st.Find("div")
		st.Find("a").Each(func(_ int, st *goquery.Selection) {
			row, e := st.Attr("href")
			if e == true {
				results = append(results, row)
			}
		})

		// Click "See More Results" at bottom of page
		st = s.browser.Find("div#see_more_pager") // Find div with link
		st.Find("a").AddClass("NextPageLink")     // Add name to link tag so we can click it
		err = s.browser.Click("a.NextPageLink")   // Click link
		if err == nil {
		} else if err.Error() == "Element not found matching expr 'a.NextPageLink'." {
			break
		} else if err != nil {
			panic(err)
		}
	}
	fmt.Println(fmt.Sprintf("Found %d images", len(results)))

	for _, pURL := range results {
		p, _ := s.getPhotoPage(pURL)
		photos = append(photos, p)
	}

	s.DB.SetPhotos(photos)
	return nil
}
