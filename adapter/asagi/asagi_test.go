package asagi

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNYCToUTC(t *testing.T) {
	type args struct {
		timestamp int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "converts NYC Epoch to UTC Epoch",
			args: args{
				timestamp: 0,
			},
			want: time.Unix(0, 0).Add(4 * time.Hour),
		},
		{
			name: "converts 4chan time to UTC",
			args: args{
				timestamp: 1343066185,
			},
			want: time.Unix(1343080585, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NYCToUTC(tt.args.timestamp); got.Unix() != tt.want.Unix() {
				t.Errorf("NYCToUTC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPostsByThread(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	const threadNum = 1

	mock.ExpectQuery("SELECT \\* FROM c WHERE thread_num = \\? ORDER BY doc_id").WithArgs(threadNum).
		WillReturnRows(sqlmock.NewRows([]string{
			"doc_id",
			"media_id",
			"poster_ip",
			"num",
			"subnum",
			"thread_num",
			"op",
			"timestamp",
			"timestamp_expired",
			"preview_orig",
			"preview_w",
			"preview_h",
			"media_filename",
			"media_w",
			"media_h",
			"media_size",
			"media_hash",
			"media_orig",
			"spoiler",
			"deleted",
			"capcode",
			"email",
			"name",
			"trip",
			"title",
			"comment",
			"delpass",
			"sticky",
			"poster_hash",
			"poster_country",
			"exif",
			"locked",
		}).AddRow(1,
			"1",
			"0",
			"1880736",
			"0",
			"1880736",
			"1",
			"1343066185",
			"0",
			"1343080585660s.jpg",
			"187",
			"250",
			"Berserk Poteymo.jpg",
			"600",
			"800",
			"112640",
			"5uFlSeyLaG+tNIiAkaXI9g==",
			"1343080585660.jpg",
			"0",
			"1",
			"N",
			sql.NullString{},
			"Anonymous",
			sql.NullString{},
			"Gar Moe Thread",
			"\"Hey /c/, let's have a Gar-Moe thread.\"",
			sql.NullString{},
			0,
			sql.NullString{},
			sql.NullString{},
			sql.NullString{},
			0).
			AddRow(2, 2, 0, 1880738, 0, 1880736, 0, 1343066321, 0, "1343080721444s.jpg", 89, 125, "Desu Ex 3.jpg", 2208, 3080, 2170552, "6YPXJCRRpTKIsnjF898gAw==", "1343080721444.jpg", 0, 0, "N", sql.NullString{}, "Anonymous", sql.NullString{}, sql.NullString{}, "Desu Ex 3", sql.NullString{}, 0, sql.NullString{}, sql.NullString{}, sql.NullString{}, 0))
	loader, err := NewLoader("c", sqlx.NewDb(db, "sqlmock"))
	assert.NoError(t, err)

	posts, err := loader.GetPostsByThread(threadNum)
	assert.NoError(t, err)

	expectedPosts := []Post{{
		DocID:            1,
		MediaID:          1,
		PosterIP:         0,
		Num:              1880736,
		Subnum:           0,
		ThreadNum:        1880736,
		Op:               1,
		Timestamp:        1343066185,
		TimestampExpired: 0,
		PreviewOrig:      sql.NullString{String: "1343080585660s.jpg", Valid: true},
		PreviewW:         187,
		PreviewH:         250,
		MediaFilename:    sql.NullString{String: "Berserk Poteymo.jpg", Valid: true},
		MediaW:           600,
		MediaH:           800,
		MediaSize:        112640,
		MediaHash:        sql.NullString{String: "5uFlSeyLaG+tNIiAkaXI9g==", Valid: true},
		MediaOrig:        sql.NullString{String: "1343080585660.jpg", Valid: true},
		Spoiler:          0,
		Deleted:          1,
		Capcode:          "N",
		Email:            sql.NullString{String: "", Valid: false},
		Name:             sql.NullString{String: "Anonymous", Valid: true},
		Trip:             sql.NullString{String: "", Valid: false},
		Title:            sql.NullString{String: "Gar Moe Thread", Valid: true},
		Comment:          sql.NullString{String: "\"Hey /c/, let's have a Gar-Moe thread.\"", Valid: true},
		Delpass:          sql.NullString{String: "", Valid: false},
		Sticky:           0,
		PosterHash:       sql.NullString{String: "", Valid: false},
		PosterCountry: sql.NullString{
			String: "",
			Valid:  false,
		}},
		{
			DocID:            2,
			MediaID:          2,
			PosterIP:         0,
			Num:              1880738,
			Subnum:           0,
			ThreadNum:        1880736,
			Op:               0,
			Timestamp:        1343066321,
			TimestampExpired: 0,
			PreviewOrig:      sql.NullString{String: "1343080721444s.jpg", Valid: true},
			PreviewW:         89,
			PreviewH:         125,
			MediaFilename:    sql.NullString{String: "Desu Ex 3.jpg", Valid: true},
			MediaW:           2208,
			MediaH:           3080,
			MediaSize:        2170552,
			MediaHash:        sql.NullString{String: "6YPXJCRRpTKIsnjF898gAw==", Valid: true},
			MediaOrig:        sql.NullString{String: "1343080721444.jpg", Valid: true},
			Spoiler:          0,
			Deleted:          0,
			Capcode:          "N",
			Email:            sql.NullString{String: "", Valid: false},
			Name:             sql.NullString{String: "Anonymous", Valid: true},
			Trip:             sql.NullString{String: "", Valid: false},
			Title:            sql.NullString{String: "", Valid: false},
			Comment:          sql.NullString{String: "Desu Ex 3", Valid: true},
			Delpass:          sql.NullString{String: "", Valid: false},
			Sticky:           0,
			PosterHash:       sql.NullString{String: "", Valid: false},
			PosterCountry:    sql.NullString{String: "", Valid: false},
			Exif:             sql.NullString{String: "", Valid: false},
			Locked:           0,
		},
	}

	assert.Equal(t, expectedPosts, posts)
}

func TestGetPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	const timeInNYCSec = 1299985600
	const timeInUTC = 1300000000

	mock.ExpectQuery("SELECT c\\.\\* FROM c_threads INNER JOIN c ON c\\.thread_num = c_threads.thread_num WHERE time_op <= \\? AND time_last >= \\?  ORDER BY doc_id;").WithArgs(timeInNYCSec, timeInNYCSec).
		WillReturnRows(sqlmock.NewRows([]string{
			"doc_id",
			"media_id",
			"poster_ip",
			"num",
			"subnum",
			"thread_num",
			"op",
			"timestamp",
			"timestamp_expired",
			"preview_orig",
			"preview_w",
			"preview_h",
			"media_filename",
			"media_w",
			"media_h",
			"media_size",
			"media_hash",
			"media_orig",
			"spoiler",
			"deleted",
			"capcode",
			"email",
			"name",
			"trip",
			"title",
			"comment",
			"delpass",
			"sticky",
			"poster_hash",
			"poster_country",
			"exif",
			"locked",
		}).AddRow(1,
			"1",
			"0",
			"1880736",
			"0",
			"1880736",
			"1",
			"1343066185",
			"0",
			"1343080585660s.jpg",
			"187",
			"250",
			"Berserk Poteymo.jpg",
			"600",
			"800",
			"112640",
			"5uFlSeyLaG+tNIiAkaXI9g==",
			"1343080585660.jpg",
			"0",
			"1",
			"N",
			sql.NullString{},
			"Anonymous",
			sql.NullString{},
			"Gar Moe Thread",
			"\"Hey /c/, let's have a Gar-Moe thread.\"",
			sql.NullString{},
			0,
			sql.NullString{},
			sql.NullString{},
			sql.NullString{},
			0).
			AddRow(2, 2, 0, 1880738, 0, 1880736, 0, 1343066321, 0, "1343080721444s.jpg", 89, 125, "Desu Ex 3.jpg", 2208, 3080, 2170552, "6YPXJCRRpTKIsnjF898gAw==", "1343080721444.jpg", 0, 0, "N", sql.NullString{}, "Anonymous", sql.NullString{}, sql.NullString{}, "Desu Ex 3", sql.NullString{}, 0, sql.NullString{}, sql.NullString{}, sql.NullString{}, 0))
	loader, err := NewLoader("c", sqlx.NewDb(db, "sqlmock"))
	assert.NoError(t, err)

	posts, err := loader.GetPosts(time.Unix(timeInUTC, 0))
	assert.NoError(t, err)

	expectedPosts := []Post{{
		DocID:            1,
		MediaID:          1,
		PosterIP:         0,
		Num:              1880736,
		Subnum:           0,
		ThreadNum:        1880736,
		Op:               1,
		Timestamp:        1343066185,
		TimestampExpired: 0,
		PreviewOrig:      sql.NullString{String: "1343080585660s.jpg", Valid: true},
		PreviewW:         187,
		PreviewH:         250,
		MediaFilename:    sql.NullString{String: "Berserk Poteymo.jpg", Valid: true},
		MediaW:           600,
		MediaH:           800,
		MediaSize:        112640,
		MediaHash:        sql.NullString{String: "5uFlSeyLaG+tNIiAkaXI9g==", Valid: true},
		MediaOrig:        sql.NullString{String: "1343080585660.jpg", Valid: true},
		Spoiler:          0,
		Deleted:          1,
		Capcode:          "N",
		Email:            sql.NullString{String: "", Valid: false},
		Name:             sql.NullString{String: "Anonymous", Valid: true},
		Trip:             sql.NullString{String: "", Valid: false},
		Title:            sql.NullString{String: "Gar Moe Thread", Valid: true},
		Comment:          sql.NullString{String: "\"Hey /c/, let's have a Gar-Moe thread.\"", Valid: true},
		Delpass:          sql.NullString{String: "", Valid: false},
		Sticky:           0,
		PosterHash:       sql.NullString{String: "", Valid: false},
		PosterCountry: sql.NullString{
			String: "",
			Valid:  false,
		}},
		{
			DocID:            2,
			MediaID:          2,
			PosterIP:         0,
			Num:              1880738,
			Subnum:           0,
			ThreadNum:        1880736,
			Op:               0,
			Timestamp:        1343066321,
			TimestampExpired: 0,
			PreviewOrig:      sql.NullString{String: "1343080721444s.jpg", Valid: true},
			PreviewW:         89,
			PreviewH:         125,
			MediaFilename:    sql.NullString{String: "Desu Ex 3.jpg", Valid: true},
			MediaW:           2208,
			MediaH:           3080,
			MediaSize:        2170552,
			MediaHash:        sql.NullString{String: "6YPXJCRRpTKIsnjF898gAw==", Valid: true},
			MediaOrig:        sql.NullString{String: "1343080721444.jpg", Valid: true},
			Spoiler:          0,
			Deleted:          0,
			Capcode:          "N",
			Email:            sql.NullString{String: "", Valid: false},
			Name:             sql.NullString{String: "Anonymous", Valid: true},
			Trip:             sql.NullString{String: "", Valid: false},
			Title:            sql.NullString{String: "", Valid: false},
			Comment:          sql.NullString{String: "Desu Ex 3", Valid: true},
			Delpass:          sql.NullString{String: "", Valid: false},
			Sticky:           0,
			PosterHash:       sql.NullString{String: "", Valid: false},
			PosterCountry:    sql.NullString{String: "", Valid: false},
			Exif:             sql.NullString{String: "", Valid: false},
			Locked:           0,
		},
	}

	assert.Equal(t, expectedPosts, posts)
}
