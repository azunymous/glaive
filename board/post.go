package board

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Thread struct {
	Post              `json:"post"`
	Subject           string    `json:"subject"`
	Replies           []Post    `json:"replies"`
	LastPostTimestamp time.Time `json:"last_post_timestamp"`
}

// Returns the index and post of the reply with the given post no.
// Returns -1 if the no is not found in the replies.
func (t Thread) getReplyWithPostNo(no uint64) (index int, p Post) {
	for i, p := range t.Replies {
		if p.No == no {
			return i, p
		}
	}
	return -1, Post{}
}

type Post struct {
	No              uint64    `json:"no"`
	Timestamp       Timestamp `json:"timestamp"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Comment         string    `json:"comment"`
	CommentSegments []Segment `json:"comment_segments"`
	Image           string    `json:"image"`
	ThumbnailImage  string    `json:"thumbnail_image"`
	Filename        string    `json:"filename"`
	Meta            string    `json:"meta"`
	QuotedBy        []uint64  `json:"quoted_by"`
}

type Timestamp struct {
	time.Time
}

func TS(t time.Time) Timestamp {
	return Timestamp{t}
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	timestamp := t.Time.UnixNano() / int64(time.Millisecond)
	return []byte(strconv.Itoa(int(timestamp))), nil
}

func (p Post) quotedBy(postQuotingNo uint64) Post {
	p.QuotedBy = append(p.QuotedBy, postQuotingNo)
	return p
}

func (p Post) Update() (Post, []Transform) {
	post := p

	if len(post.Name) < 1 {
		post.Name = "Anonymous"
	}

	if post.Email == "sage" || post.Email == "noko" || post.Email == "nokosage" {
		post.Meta = post.Email
		post.Email = ""
	}

	post, threadTransformations := post.parse()
	return post, threadTransformations
}

// One line with either a format or not
type Segment struct {
	Format  []string `json:"format"`
	Segment string   `json:"segment"`
}

type format struct {
	regex                  *regexp.Regexp
	class                  string
	transformationProvider func(f format, line string, p Post) func(t Thread) Thread
}

type Transform func(t Thread) Thread

func getFormats() []format {
	return []format{
		{
			regex:                  regexp.MustCompile(`^>([^>].*)`),
			class:                  "quote",
			transformationProvider: nil,
		},
		{
			regex:                  regexp.MustCompile(`^>>(\d+)[ \t]*`),
			class:                  "noQuote",
			transformationProvider: addPostReplyToQuotedPost,
		},

		{
			regex:                  regexp.MustCompile(`^[ \t]*Objection![ \t]*`),
			class:                  "objection",
			transformationProvider: nil, // FOR NOW!
		},
	}
}

func (p Post) parse() (Post, []Transform) {
	post := p
	postContent := post.Comment
	var segments []Segment
	var transformations []Transform
	for _, line := range strings.Split(postContent, "\n") {
		addingSegment := Segment{[]string{}, line}
		for _, f := range getFormats() {
			find := f.regex.FindString(line)
			if find != "" {
				addingSegment = Segment{[]string{f.class}, line}
				if f.transformationProvider != nil {
					transformations = append(transformations, f.transformationProvider(f, line, p))
				}
			}
		}
		segments = append(segments, addingSegment)
	}

	post.CommentSegments = segments
	return post, transformations
}

func addPostReplyToQuotedPost(f format, line string, p Post) func(t Thread) Thread {
	return func(t Thread) Thread {
		submatches := f.regex.FindStringSubmatch(line)
		if len(submatches[1]) > 0 {
			quotedPostNo, _ := strconv.ParseUint(submatches[1], 10, 0)
			index, transformedPost := t.getReplyWithPostNo(quotedPostNo)
			if quotedPostNo == t.No {
				t.Post = t.Post.quotedBy(p.No)
			} else if index > -1 {
				t.Replies[index] = transformedPost.quotedBy(p.No)
			}
		}
		return t
	}
}
