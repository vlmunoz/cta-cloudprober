package main

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DriveState struct {
	bun.BaseModel `bun:"table:drive_state"`

	DriveName                string
	Host                     string
	LogicalLibrary           string
	SessionID                int64 `bun:session_id`
	BytesTransferedInSession int64
	FilesTransferedInSession int64
	SessionStartTime         int64
	SessionElapsedTime       int64
	MountStartTime           int64
	TransferStartTime        int64
	UnloadStartTime          int64
	UnmountStartTime         int64
	DrainingStartTime        int64
	DownOrUpStartTime        int64
	ProbeStartTime           int64
	CleanupStartTime         int64
	StartStartTime           int64
	ShutdownTime             int64
	MountType                string
	DriveStatus              string
	DesiredUp                string
	DesiredForceDown         string
	ReasonUpDown             string
	CurrentVID               string `bun:"current_vid"`
	CTAVersion               string `bun:"cta_version"`
	CurrentPriority          string
	CurrentActivity          string
	CurrentTapePool          string
	NextMountType            string
	NextVID                  string `bun:"next_vid"`
	NextPriority             string
	NextActivity             string
	NextTapePool             string
	DevFileName              string
	RawLibrarySlot           string
	CurrentVO                string `bun:"current_vo"`
	NextVO                   string `bun:"next_vo"`
	UserComment              string
	CreationLogUserName      string
	CreationLogHostName      string
	CreationLogTime          int64
	LastUpdateUserName       string
	LastUpdateHostName       string
	LastUpdateTime           int64
	DiskSystemName           string
	ReservedBytes            int64
	ReservationSessionID     int64 `bun:resrvation_session_id`
}

func snake_case(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		// A is 65, a is 97
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		// v is capital letter here
		// irregard first letter
		// add underscore if last letter is capital letter
		// add underscore when previous letter is lowercase
		// add underscore when next letter is lowercase
		if (i != 0 || i == l-1) && (          // head and tail
		(i > 0 && rune(camel[i-1]) >= 'a') || // pre
			(i < l-1 && rune(camel[i+1]) >= 'a')) { //next
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}

func main() {
	dsn := "postgres://cta:test4CTA!@ifdb07.fnal.gov:5438/cta_dev?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	ctx := context.Background()

	/* db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	)) */

	drive_states := make([]DriveState, 0)
	db.NewSelect().Model(&drive_states).Column("*").Scan(ctx)

	fields := reflect.VisibleFields(reflect.TypeOf(struct{ DriveState }{}))
	for _, drive_state := range drive_states {
		r := reflect.ValueOf(drive_state)
		drive_name := reflect.Indirect(r).FieldByName("DriveName")
		session_id := reflect.Indirect(r).FieldByName("SessionID")
		drive_status := reflect.Indirect(r).FieldByName("DriveStatus")
		current_vid := reflect.Indirect(r).FieldByName("CurrentVID")
		current_tape_pool := reflect.Indirect(r).FieldByName("CurrentTapePool")
		current_activity := reflect.Indirect(r).FieldByName("CurrentActivity")
		for _, field := range fields {
			if field.Name == "BaseModel" {
				continue
			}
			f := reflect.Indirect(r).FieldByName(field.Name)

			// Set column name to snake case by default
			column_name := snake_case(field.Name)

			// Set column name to tag value if overridden with tag
			tag := field.Tag.Get("bun")
			if tag != "" {
				column_name = tag
			}

			// Give 0 value to bytes transferred while working on first file
			/* if field.Name == "BytesTransferredInSession" && session_id.IsValid() && !session_id.IsZero() {
				if f.IsValid() && f.IsZero() {
					x := int64(0)
					f.SetInt(x)
				}
			} */
			// Already has 0 value I just wasn't outputting!

			if f.IsValid() {
				fmt.Printf("%s{drive_name=%s,session_id=%d,drive_status=%s,current_vid=%s,current_tape_pool=%s,current_activity=%s} %d\n",
					column_name, drive_name, session_id, drive_status, current_vid, current_tape_pool, current_activity, f)
			}
		}
	}
}
