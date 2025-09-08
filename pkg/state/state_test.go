package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethanjantz/oobd/pkg/scan"
)

func TestUpdate(t *testing.T) {
	now := time.Now()
	s := Empty()

	notif, err := s.UpdateFor(1000, scan.Result{}, now)
	require.NoError(t, err)
	assert.Nil(t, notif)
	assert.Equal(t, map[uint32]UserState{
		1000: UserState{
			LatestScan:    scan.Result{},
			Notifications: map[time.Time]Notification{},
			Scheduled:     map[time.Time]scan.Result{},
		},
	}, s.Recursers)

	f := scan.Meta{now, 42}
	fs := scan.Result{"foo": f}
	notif, err = s.UpdateFor(1000, fs, now)
	require.NoError(t, err)
	require.Equal(t, &Notification{
		When:      now.Add(delGrace),
		Scheduled: fs,
		Deleted:   nil,
	}, notif)
	assert.Equal(t, map[uint32]UserState{
		1000: UserState{
			LatestScan: fs,
			Notifications: map[time.Time]Notification{
				now: *notif,
			},
			Scheduled: map[time.Time]scan.Result{
				now.Add(delGrace): fs,
			},
		},
	}, s.Recursers)

	notif2, err := s.UpdateFor(1000, scan.Result{}, now.Add(time.Second))
	require.NoError(t, err)
	assert.Nil(t, notif2)
	assert.Equal(t, map[uint32]UserState{
		1000: UserState{
			LatestScan: scan.Result{},
			Notifications: map[time.Time]Notification{
				now: *notif,
			},
			Scheduled: map[time.Time]scan.Result{},
		},
	}, s.Recursers)

	_, err = s.UpdateFor(1000, fs, now)
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "Already notified user 1000")

	notif2, err = s.UpdateFor(1000, fs, now.Add(time.Second))
	require.NoError(t, err)
	///require.NotNil(t, err)
	assert.Equal(t, map[uint32]UserState{
		1000: UserState{
			LatestScan: fs,
			Notifications: map[time.Time]Notification{
				now: *notif,
				now.Add(time.Second): *notif2,
			},
			Scheduled: map[time.Time]scan.Result{
				now.Add(delGrace+time.Second): fs,
			},
		},
	}, s.Recursers)

	fs2 := scan.Result{
		"foo": f,
		"bar": f,
		"baz": f,
	}
	now3 := now.Add(31*24*time.Hour)
	notif3, err := s.UpdateFor(1000, fs2, now3)
	require.NoError(t, err)
	require.Equal(t, &Notification{
		When:      now3.Add(delGrace),
		Scheduled: scan.Result{
			"bar": f,
			"baz": f,
		},
		Deleted:   fs,
	}, notif3)
	assert.Equal(t, map[uint32]UserState{
		1000: UserState{
			LatestScan: fs2,
			Notifications: map[time.Time]Notification{
				now: *notif,
				now.Add(time.Second): *notif2,
				now3: *notif3,
			},
			Scheduled: map[time.Time]scan.Result{
				now3.Add(delGrace): scan.Result{
					"bar": f,
					"baz": f,
				},
			},
		},
	}, s.Recursers)
}
