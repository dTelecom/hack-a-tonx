package ton

import (
	"github.com/xssnick/tonutils-go/tvm/cell"
	"time"
)

func extractCreateCallMessage(msg []byte) (callId uint64, validUntil time.Time) {
	slice := cell.BeginCell().MustStoreSlice(msg, 8*uint(len(msg))).EndCell().BeginParse()
	callId = slice.MustLoadUInt(64)
	validUntil = time.Unix(int64(slice.MustLoadUInt(32)), 0)
	return
}

func extractEndCallMessage(msg []byte) (callId uint64, validUntil time.Time, spentMinutes uint32) {
	slice := cell.BeginCell().MustStoreSlice(msg, 8*uint(len(msg))).EndCell().BeginParse()
	callId = slice.MustLoadUInt(64)
	validUntil = time.Unix(int64(slice.MustLoadUInt(32)), 0)
	spentMinutes = uint32(slice.MustLoadUInt(32))
	return
}
