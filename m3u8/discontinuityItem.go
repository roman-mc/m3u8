package m3u8

// DiscontinuityItem represents a EXT-X-DISCONTINUITY tag to indicate a
// discontinuity between the SegmentItems that proceed and follow it.
type DiscontinuityItem struct{}

// NewDiscontinuityItem returns a *DiscontinuityItem
func NewDiscontinuityItem() (*DiscontinuityItem, error) {
	return &DiscontinuityItem{}, nil
}

func (di *DiscontinuityItem) String() string {
	return DiscontinuityItemTag
}
