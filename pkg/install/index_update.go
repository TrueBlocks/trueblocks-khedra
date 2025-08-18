package install

// UpdateIndexStrategy updates the draft's index acquisition strategy and detail
// level and persists the rough estimates directly into the draft metadata. It
// does not save the draft to disk; caller should invoke SaveDraftAtomic.
func UpdateIndexStrategy(d *Draft, strategy, detail string) {
	if d == nil {
		return
	}
	d.Config.General.Strategy = strategy
	d.Config.General.Detail = detail
	disk, hrs := EstimateIndex(strategy, detail)
	d.Meta.EstDiskGB = disk
	d.Meta.EstHours = hrs
}
