package state

import (
	"fmt"
	"maps"
	"time"

	"github.com/ethanjantz/oobd/pkg/scan"
)

type Notification struct {
	When      time.Time
	Scheduled scan.Result
	Deleted   scan.Result
}

type UserState struct {
	LatestScan    scan.Result
	Scheduled     map[time.Time]scan.Result
	Notifications map[time.Time]Notification
}

type State struct {
	Recursers map[uint32]UserState
}

func Empty() *State {
	return &State{Recursers: map[uint32]UserState{}}
}

var delGrace = time.Hour * 24 * 30

func (s *State) UpdateFor(rcid uint32, latest scan.Result, now time.Time) (*Notification, error) {
	user, ok := s.Recursers[rcid]
	if !ok {
		user = UserState{scan.Result{}, map[time.Time]scan.Result{}, map[time.Time]Notification{}}
	}

	if _, ok := user.Notifications[now]; ok {
		return nil, fmt.Errorf("Already notified user %d at %s. Bug?", rcid, now)
	}

	added := scan.Result{}
	for p, m := range latest {
		if _, ok := user.LatestScan[p]; !ok {
			added[p] = m
		}
	}

	var notif *Notification
	if len(added) > 0 {
		when := now.Add(delGrace)
		notif = &Notification{when, added, nil}
		user.Scheduled[when] = added
	}

	del := scan.Result{}
	newSched := map[time.Time]scan.Result{}
	for t, s := range user.Scheduled {
		pruned := scan.Result{}
		for f, meta := range s {
			if _, ok := latest[f]; ok {
				pruned[f] = meta
			}
		}
		if len(pruned) == 0 {
			continue
		}

		if t.Before(now) {
			maps.Copy(del, pruned)
		} else {
			newSched[t] = pruned
		}
	}
	user.Scheduled = newSched

	if len(del) > 0 {
		if notif == nil {
			notif = &Notification{}
		}
		notif.Deleted = del

	}

	user.LatestScan = latest
	s.Recursers[rcid] = user
	if notif != nil {
		user.Notifications[now] = *notif
	}

	return notif, nil
}
