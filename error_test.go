package sqlite3

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/loveuer/go-sqlite3/internal/util"
)

func Test_assertErr(t *testing.T) {
	err := util.AssertErr()
	if s := err.Error(); !strings.HasPrefix(s, "sqlite3: assertion failed") || !strings.HasSuffix(s, "error_test.go:13)") {
		t.Errorf("got %q", s)
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	var ecode ErrorCode
	var xcode xErrorCode
	err := &Error{code: 0x8080}
	if !errors.As(err, &err) {
		t.Fatal("want true")
	}
	if ecode := err.Code(); ecode != 0x80 {
		t.Errorf("got %#x, want 0x80", uint8(ecode))
	}
	if ok := errors.As(err, &ecode); !ok || ecode != ErrorCode(0x80) {
		t.Errorf("got %#x, want 0x80", uint8(ecode))
	}
	if !errors.Is(err, ErrorCode(0x80)) {
		t.Errorf("want true")
	}
	if xcode := err.ExtendedCode(); xcode != 0x8080 {
		t.Errorf("got %#x, want 0x8080", uint16(xcode))
	}
	if ok := errors.As(err, &xcode); !ok || xcode != xErrorCode(0x8080) {
		t.Errorf("got %#x, want 0x8080", uint16(xcode))
	}
	if !errors.Is(err, xErrorCode(0x8080)) {
		t.Errorf("want true")
	}
	if s := err.Error(); s != "sqlite3: 32896" {
		t.Errorf("got %q", s)
	}
	if ok := errors.As(err.ExtendedCode(), &ecode); !ok || ecode != ErrorCode(0x80) {
		t.Errorf("got %#x, want 0x80", uint8(ecode))
	}
	if !errors.Is(err.ExtendedCode(), ErrorCode(0x80)) {
		t.Errorf("want true")
	}
}

func TestError_Temporary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code uint64
		want bool
	}{
		{"ERROR", uint64(ERROR), false},
		{"BUSY", uint64(BUSY), true},
		{"BUSY_RECOVERY", uint64(BUSY_RECOVERY), true},
		{"BUSY_SNAPSHOT", uint64(BUSY_SNAPSHOT), true},
		{"BUSY_TIMEOUT", uint64(BUSY_TIMEOUT), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			{
				err := &Error{code: tt.code}
				if got := err.Temporary(); got != tt.want {
					t.Errorf("Error.Temporary(%d) = %v, want %v", tt.code, got, tt.want)
				}
			}
			{
				err := ErrorCode(tt.code)
				if got := err.Temporary(); got != tt.want {
					t.Errorf("ErrorCode.Temporary(%d) = %v, want %v", tt.code, got, tt.want)
				}
			}
			{
				err := ExtendedErrorCode(tt.code)
				if got := err.Temporary(); got != tt.want {
					t.Errorf("ExtendedErrorCode.Temporary(%d) = %v, want %v", tt.code, got, tt.want)
				}
			}
		})
	}
}

func TestError_Timeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code uint64
		want bool
	}{
		{"ERROR", uint64(ERROR), false},
		{"BUSY", uint64(BUSY), false},
		{"BUSY_RECOVERY", uint64(BUSY_RECOVERY), false},
		{"BUSY_SNAPSHOT", uint64(BUSY_SNAPSHOT), false},
		{"BUSY_TIMEOUT", uint64(BUSY_TIMEOUT), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			{
				err := &Error{code: tt.code}
				if got := err.Timeout(); got != tt.want {
					t.Errorf("Error.Timeout(%d) = %v, want %v", tt.code, got, tt.want)
				}
			}
			{
				err := ExtendedErrorCode(tt.code)
				if got := err.Timeout(); got != tt.want {
					t.Errorf("Error.Timeout(%d) = %v, want %v", tt.code, got, tt.want)
				}
			}
		})
	}
}

func Test_ErrorCode_Error(t *testing.T) {
	t.Parallel()

	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Test all error codes.
	for i := 0; i == int(ErrorCode(i)); i++ {
		want := "sqlite3: "
		r := db.call("sqlite3_errstr", uint64(i))
		want += util.ReadString(db.mod, uint32(r), _MAX_NAME)

		got := ErrorCode(i).Error()
		if got != want {
			t.Fatalf("got %q, want %q, with %d", got, want, i)
		}
	}
}

func Test_ExtendedErrorCode_Error(t *testing.T) {
	t.Parallel()

	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Test all extended error codes.
	for i := 0; i == int(ExtendedErrorCode(i)); i++ {
		want := "sqlite3: "
		r := db.call("sqlite3_errstr", uint64(i))
		want += util.ReadString(db.mod, uint32(r), _MAX_NAME)

		got := ExtendedErrorCode(i).Error()
		if got != want {
			t.Fatalf("got %q, want %q, with %d", got, want, i)
		}
	}
}

func Test_errorCode(t *testing.T) {
	tests := []struct {
		arg      error
		wantMsg  string
		wantCode uint32
	}{
		{nil, "", _OK},
		{ERROR, "", util.ERROR},
		{IOERR, "", util.IOERR},
		{IOERR_READ, "", util.IOERR_READ},
		{&Error{code: util.ERROR}, "", util.ERROR},
		{fmt.Errorf("%w", ERROR), ERROR.Error(), util.ERROR},
		{fmt.Errorf("%w", IOERR), IOERR.Error(), util.IOERR},
		{fmt.Errorf("%w", IOERR_READ), IOERR_READ.Error(), util.IOERR_READ},
		{fmt.Errorf("error"), "error", util.ERROR},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotMsg, gotCode := errorCode(tt.arg, ERROR)
			if gotMsg != tt.wantMsg {
				t.Errorf("errorCode() gotMsg = %q, want %q", gotMsg, tt.wantMsg)
			}
			if gotCode != uint32(tt.wantCode) {
				t.Errorf("errorCode() gotCode = %d, want %d", gotCode, tt.wantCode)
			}
		})
	}
}
