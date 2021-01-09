package asagi

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"glaive/board"
	"path"
	"sort"
	"strings"
	"text/template"
	"time"
)

type Loader struct {
	db           DB
	getPostsSQL  string
	getThreadSQL string
	board        string
}

type DB interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

const GetPostsTemplate = `SELECT {{. }}.* FROM {{. }}_threads INNER JOIN {{. }} ON {{. }}.thread_num = {{. }}_threads.thread_num WHERE time_op <= ? AND time_last >= ?  ORDER BY doc_id;`
const GetThreadTemplate = `SELECT * FROM {{. }} WHERE thread_num = ? ORDER BY doc_id`

func NewLoader(board string, db DB) (*Loader, error) {
	var err error
	var sb strings.Builder
	err = template.Must(template.New("getPostsSQL").Parse(GetPostsTemplate)).Execute(&sb, board)
	if err != nil {
		return nil, err
	}
	getPosts := sb.String()
	sb.Reset()
	err = template.Must(template.New("getThreadSQL").Parse(GetThreadTemplate)).Execute(&sb, board)
	if err != nil {
		return nil, err
	}
	getThread := sb.String()
	return &Loader{
		db:           db,
		getPostsSQL:  getPosts,
		getThreadSQL: getThread,
		board:        board,
	}, nil
}

func NewSqlConn(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}

func (l *Loader) GetPostsByThread(threadNo int) ([]Post, error) {
	ctx := context.Background()

	var posts []Post
	// The timestamp is seconds since 1st Jan 1970 NYC rather than UTC
	err := l.db.SelectContext(ctx, &posts, l.getThreadSQL, threadNo)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (l *Loader) GetPosts(time time.Time) ([]Post, error) {
	ctx := context.Background()
	est := UTCToNYC(time.Unix())
	var posts []Post
	// The timestamp is seconds since 1st Jan 1970 NYC rather than UTC
	err := l.db.SelectContext(ctx, &posts, l.getPostsSQL, est.Unix(), est.Unix())
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (l *Loader) PostToThread(posts []Post) *board.Thread {
	threads := l.PostToThreads(posts, DoNotDiscard)
	if len(threads) > 1 {
		panic("Only expected one thread")
	}
	if len(threads) == 0 {
		panic("No threads!")
	}
	return threads[0]
}

func DiscardIfAfter(time time.Time) func(board.Post) bool {
	return func(post board.Post) bool {
		return post.Timestamp.After(time)
	}
}

func DoNotDiscard(_ board.Post) bool {
	return false
}

func (l *Loader) PostToThreads(posts []Post, discardRule func(board.Post) bool) []*board.Thread {
	var threads []*board.Thread
	threadMap := make(map[uint64]*board.Thread)
	for _, post := range posts {
		if post.Op == 1 {
			p, _ := l.PostToPost(post)
			t := &board.Thread{
				Post:              p,
				Subject:           post.Title.String,
				Replies:           make([]board.Post, 0),
				LastPostTimestamp: p.Timestamp.Time,
			}
			threads = append(threads, t)
			threadMap[t.No] = t
		} else {
			if t, ok := threadMap[uint64(post.ThreadNum)]; ok {
				p, transformations := l.PostToPost(post)
				if discardRule(p) {
					continue
				}
				t.Replies = append(t.Replies, p)
				if p.Meta == "" {
					t.LastPostTimestamp = p.Timestamp.Time
				}
				for _, transformation := range transformations {
					*t = transformation(*t)
				}
			} else {
				panic("Orphan Post!")
			}
		}

	}

	sort.SliceStable(threads, func(i, j int) bool {
		return threads[i].LastPostTimestamp.After(threads[j].LastPostTimestamp)
	})

	return threads
}

const image = "image"
const thumb = "thumb"

func (l *Loader) PostToPost(post Post) (board.Post, []board.Transform) {
	p := board.Post{
		No:              uint64(post.Num),
		Timestamp:       board.TS(NYCToUTC(post.Timestamp)),
		Name:            post.Name.String,
		Email:           post.Email.String,
		Comment:         post.Comment.String,
		CommentSegments: nil,
		// TODO this is only the thumbnail, the real file is post.MediaOrig
		Image:          l.toPath(post.MediaOrig.String, image),
		ThumbnailImage: l.toPath(post.PreviewOrig.String, thumb),
		Filename:       post.MediaFilename.String,
		Meta:           "",
		QuotedBy:       nil,
	}
	return p.Update()
}

func (l *Loader) toPath(s string, sizeName string) string {
	if s == "" {
		return ""
	}

	if len(s) < 6 {
		return s
	}

	return path.Join(l.board, sizeName, s[:4], s[4:6], s)

}

func NYCToUTC(timestamp int) time.Time {
	return time.Unix(int64(timestamp), 0).Add(4 * time.Hour)
}

func UTCToNYC(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).Add(-4 * time.Hour)
}
