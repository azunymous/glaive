package asagi

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"glaive/board"
	"path"
	"sort"
	"time"
)

type Loader struct {
	db *sqlx.DB
}

func NewLoader(dsn string) (*Loader, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Loader{
		db: db,
	}, nil
}

const GetPosts = `SELECT c.* FROM c_threads INNER JOIN c ON c.thread_num = c_threads.thread_num WHERE time_op <= ? AND time_last >= ?  ORDER BY doc_id;`
const GetThread = `SELECT * FROM c WHERE thread_num = ? ORDER BY doc_id`

func (l *Loader) GetPostsByThread(threadNo int) ([]Post, error) {
	ctx := context.Background()

	var posts []Post
	// The timestamp is seconds since 1st Jan 1970 NYC rather than UTC
	err := l.db.SelectContext(ctx, &posts, GetThread, threadNo)
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
	err := l.db.SelectContext(ctx, &posts, GetPosts, est.Unix(), est.Unix())
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func PostToThread(posts []Post) *board.Thread {
	threads := PostToThreads(posts, DoNotDiscard)
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

func PostToThreads(posts []Post, discardRule func(board.Post) bool) []*board.Thread {
	var threads []*board.Thread
	threadMap := make(map[uint64]*board.Thread)
	for _, post := range posts {
		if post.Op == 1 {
			p, _ := PostToPost(post)
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
				p, transformations := PostToPost(post)
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

func PostToPost(post Post) (board.Post, []board.Transform) {
	p := board.Post{
		No:              uint64(post.Num),
		Timestamp:       board.TS(NYCToUTC(post.Timestamp)),
		Name:            post.Name.String,
		Email:           post.Email.String,
		Comment:         post.Comment.String,
		CommentSegments: nil,
		// TODO this is only the thumbnail, the real file is post.MediaOrig
		Image:          toPath(post.PreviewOrig.String, thumb),
		ThumbnailImage: toPath(post.PreviewOrig.String, thumb),
		Filename:       post.MediaFilename.String,
		Meta:           "",
		QuotedBy:       nil,
	}
	return p.Update()
}

const thumb = "thumb"

func toPath(s string, sizeName string) string {
	if s == "" {
		return ""
	}

	if len(s) < 6 {
		return s
	}

	return path.Join("c", sizeName, s[:4], s[4:6], s)

}

func NYCToUTC(timestamp int) time.Time {
	return time.Unix(int64(timestamp), 0).Add(4 * time.Hour)
}

func UTCToNYC(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).Add(-4 * time.Hour)
}
