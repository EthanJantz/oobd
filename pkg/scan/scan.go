package scan

import (
	"maps"
	"time"

	"github.com/dundee/gdu/v5/pkg/analyze"
	"github.com/dundee/gdu/v5/pkg/device"
	"github.com/dundee/gdu/v5/pkg/fs"
)

type Result map[string]Meta

type Meta struct {
	Mtime time.Time
	Size  int64
}

func Dir(p string, o Opts) (Result, error) {
	// parallel for now, just for quicker development
	// should be Seq for less resource churn
	//s := analyze.CreateSeqAnalyzer()
	s := analyze.CreateAnalyzer()
	ms, err := device.Getter.GetMounts()
	if err != nil {
		return nil, err
	}
	r := s.AnalyzeDir(p, ignore(device.GetNestedMountpointsPaths(p, ms)), true)
	<-s.GetDone()
	r.UpdateStats(fs.HardLinkedItems{})
	return wannaDel(o.now.Add(-o.minAge), r, o), nil
}

func wannaDel(t time.Time, r fs.Item, o Opts) Result {
	type rec func(rec, fs.Item) Result
	walk := func(walk rec, f fs.Item) Result {
		p, mt, s := f.GetPath(), f.GetMtime(), f.GetSize()
		if mt.Before(t) && s >= o.minSize && p != r.GetPath() {
			return Result{p: Meta{mt, s}}
		}
		ret := Result{}
		for _, c := range f.GetFiles() {
			maps.Copy(ret, walk(walk, c))
		}
		return ret
	}

	return walk(walk, r)
}

func ignore(ps []string) func(string, string) bool {
	ign := map[string]struct{}{}
	for _, p := range ps {
		ign[p] = struct{}{}
	}
	return func(_, p string) bool { _, ok := ign[p]; return ok }
}
