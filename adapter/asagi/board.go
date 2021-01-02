package asagi

import "database/sql"

type Post struct {
	DocID            int     `db:"doc_id"`
	MediaID          int     `db:"media_id"`
	PosterIP         float64 `db:"poster_ip"`
	Num              int
	Subnum           int
	ThreadNum        int `db:"thread_num"`
	Op               int
	Timestamp        int
	TimestampExpired int            `db:"timestamp_expired"`
	PreviewOrig      sql.NullString `db:"preview_orig"`
	PreviewW         int            `db:"preview_w"`
	PreviewH         int            `db:"preview_h"`
	MediaFilename    sql.NullString `db:"media_filename"`
	MediaW           int            `db:"media_w"`
	MediaH           int            `db:"media_h"`
	MediaSize        int            `db:"media_size"`
	MediaHash        sql.NullString `db:"media_hash"`
	MediaOrig        sql.NullString `db:"media_orig"`
	Spoiler          int
	Deleted          int
	Capcode          string
	Email            sql.NullString
	Name             sql.NullString
	Trip             sql.NullString
	Title            sql.NullString
	Comment          sql.NullString
	Delpass          sql.NullString
	Sticky           int
	PosterHash       sql.NullString `db:"poster_hash"`
	PosterCountry    sql.NullString `db:"poster_country"`
	Exif             sql.NullString
	Locked           int
}

type Thread struct {
	ThreadNum        int
	TimeOp           int
	TimeLast         int
	TimeBump         int
	TimeGhost        sql.NullInt64
	TimeGhostBump    sql.NullInt64
	TimeLastModified int
	Nreplies         int
	Nimages          int
	Sticky           int
	Locked           int
}

type Image struct {
	MediaID      int
	MediaHash    string
	Media        sql.NullString
	PreviewOp    sql.NullString
	PreviewReply sql.NullString
	Total        int
	Banned       int
}
