/*
 * Copyright (C) distroy
 */

package gocoverage

type FileCoverages struct {
	Filename     string
	Coverages    Coverages
	NonCoverages Coverages
}

func (f *FileCoverages) GetCount() Count {
	return Count{
		Coverages:    f.Coverages.GetCount(),
		NonCoverages: f.NonCoverages.GetCount(),
	}
}

func (f *FileCoverages) Add(c Coverage, filters ...filter) {
	if c.Count > 0 {
		f.Coverages = f.addToCoverages(f.Coverages, c, filters)
		return
	}

	filters = append(filters, func(file string, lineNo int) bool {
		return !f.Coverages.IsIn(lineNo, lineNo)
	})
	f.NonCoverages = f.addToCoverages(f.NonCoverages, c, filters)
}

func (f *FileCoverages) addToCoverages(s Coverages, c Coverage, filters []filter) Coverages {
	// log.Printf(" === before. coverages:%v, coverage:%v", s, c)
	for i := c.BeginLine; i <= c.EndLine; i++ {
		if !doFilters(c.Filename, i, filters) {
			continue
		}

		s = s.Add(Coverage{
			Filename:  c.Filename,
			BeginLine: i,
			EndLine:   i,
			Count:     c.Count,
		})
	}
	// log.Printf(" === after. coverages:%v, coverage:%v", s, c)
	return s
}

type Files map[string]*FileCoverages

func NewFileCoverages(coverages []Coverage, filters ...filter) Files {
	res := make(Files)
	for _, c := range coverages {
		if c.Count > 0 {
			res.Add(c, filters...)
		}
	}

	for _, c := range coverages {
		if c.Count <= 0 {
			res.Add(c, filters...)
		}
	}

	// for _, v := range res {
	// 	log.Printf(" === %v", *v)
	// }
	return res
}

func (f Files) Add(c Coverage, filters ...filter) *FileCoverages {
	v := f[c.Filename]
	if v == nil {
		v = &FileCoverages{
			Filename: c.Filename,
		}
		f[c.Filename] = v
	}
	v.Add(c, filters...)
	return v
}

func (f Files) GetCount() Count {
	count := Count{}
	for _, f := range f {
		res := f.GetCount()
		count.Coverages += res.Coverages
		count.NonCoverages += res.NonCoverages
	}
	return count
}
