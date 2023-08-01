package m3u8

import "errors"

var (
	// ErrPlaylistInvalid represents playlist error when playlist does not start with #EXTM3U
	ErrPlaylistInvalid = errors.New("invalid playlist, must start with #EXTM3U")

	// ErrPlaylistInvalidType represents playlist error when it's mixed between master and media playlist
	ErrPlaylistInvalidType = errors.New("invalid playlist, mixed master and media")

	// ErrSegmentItemInvalid represents error when a segment item is invalid
	ErrSegmentItemInvalid = errors.New("invalid segment item")

	// ErrPlaylistItemInvalid represents error when a playlist item is invalid
	ErrPlaylistItemInvalid = errors.New("invalid playlist item")
)
