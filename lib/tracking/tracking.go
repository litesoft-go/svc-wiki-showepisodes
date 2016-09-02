package tracking

import (
	lTracker "lib-builtin/lib/logrhythmictracker"
	dispatch "lib-builtin/lib/serveradaptor"
	json "lib-builtin/lib/jsonBuilder"
	"lib-builtin/lib/channels"
	"lib-builtin/lib/ints"
	"lib-builtin/lib/rw"
	"lib-builtin/lib/tracker"
)

type nonTemporalState struct {
	rw.Mutex // Protect:
	mShowsMonitored int
}

func (this *nonTemporalState) setShowsMonitored(pShowsMonitored int) {
	defer this.WLock().Unlock()

	this.mShowsMonitored = pShowsMonitored
}

func (this *nonTemporalState) addShowMonitored() {
	defer this.WLock().Unlock()

	this.mShowsMonitored++
}

func (this *nonTemporalState) AddToJsonObject(pObject *json.Builder) {
	defer this.RLock().Unlock()

	if this.mShowsMonitored != 0 {
		pObject.AddAttributeInt("ShowsMonitored", this.mShowsMonitored)
	}
}

type trackedPlus struct {
	mSuccessfulRequests int
	mShowsAdded         int
}

func newPlusTracked() tracker.Tracked {
	return &trackedPlus{}
}

func (copy trackedPlus) Copy() tracker.Tracked {
	return &copy
}

func (this *trackedPlus) IsEmpty() bool {
	return ints.AllZero(
		this.mSuccessfulRequests,
		this.mShowsAdded)
}

func (this *trackedPlus) AddIn(pTracked tracker.Tracked) {
	them := pTracked.(*trackedPlus)
	this.mSuccessfulRequests += them.mSuccessfulRequests
	this.mShowsAdded += them.mShowsAdded
}

func (this *trackedPlus) AddToJsonObject(pObject *json.Builder, pExcept []int) {
	pObject.
	AddAttributeInt("Successful", this.mSuccessfulRequests, pExcept...).
		AddAttributeInt("ShowsAdded", this.mShowsAdded, pExcept...)
}

type Tracker struct {
	tracker.InternalTracker
	nonTemporalState
}

//noinspection GoUnusedFunction
func NewTracker(pIssueLogger dispatch.DispatchIssues, pShutDownChannelAccessor channels.ShutDownChannelAccessor) (rTracker *Tracker) {
	rTracker = &Tracker{}
	rTracker.nonTemporalState.Mutex = rw.NewMutex()
	rTracker.Init(newPlusTracked, &rTracker.nonTemporalState, pIssueLogger, pShutDownChannelAccessor, lTracker.Sec, lTracker.TenHour)
	return
}

func (this *Tracker) SetShowsMonitored(pShowsMonitored int) {
	this.nonTemporalState.setShowsMonitored(pShowsMonitored)
}

func (this *Tracker) AddShowsMonitored() {
	this.AddPlus(&trackedPlus{mSuccessfulRequests:1, mShowsAdded:1})
	this.nonTemporalState.addShowMonitored()
}
