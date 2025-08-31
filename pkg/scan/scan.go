package scan

import (
	"fmt"
	"github.com/dundee/gdu/v5/pkg/analyze"
	"github.com/dundee/gdu/v5/pkg/fs"
	"time"
)

func Dir(p string, o Opts) {
	// parallel for now, just for quicker development
	// should be Seq for less resource churn
	//s := analyze.CreateSeqAnalyzer()
	s := analyze.CreateAnalyzer()
	r := s.AnalyzeDir(p, ignore, true)
	<-s.GetDone()
	r.UpdateStats(fs.HardLinkedItems{})
	walk(o.now.Add(-o.minAge), r, o, 2)
}

func walk(t time.Time, f fs.Item, o Opts, n uint) {
	s := f.GetSize()
	if s > o.minSize && f.GetMtime().Before(t) {
		fmt.Printf("%s %d %d\n", f.GetPath(), f.GetUsage(), s)
	}
	if n > 0 {
		for _, c := range f.GetFiles() {
			walk(t, c, o, n-1)
		}
	}
}

func ignore(name, p string) bool { return false }
