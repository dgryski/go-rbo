// Package rbo computes the rank-biased overlap for two sorted result sets
/*

"A Similarity Measure for Indefinite Rankings", William Webber, Alistair
Moffat, and Justin Zobel, ACM Transactions on Information Systems, Volume 28,
Issue 4 (November, 2010), pages 20:1--38.

    http://www.williamwebber.com/research/papers/wmz10_tois.pdf

This package is based on the reference implementation in C available from

    http://www.williamwebber.com/research/

Licensed under the MIT license.
*/
package rbo

import "math"

// State is an in-progress RBO calculation
type State struct {
	depth        int
	overlap      int
	wgt          float64
	rbo          float64
	p            float64
	seen         map[int]struct{}
	shortDepth   int
	shortOverlap int
}

// New returns a new state
func New(p float64) *State {
	return &State{
		p:            p,
		wgt:          (1 - p) / p,
		seen:         make(map[int]struct{}),
		shortDepth:   -1,
		shortOverlap: -1,
	}
}

// Update adds the next item from each list to the current state
func (s *State) Update(e1, e2 int) {
	if s.shortDepth != -1 {
		panic("rbo: Update() called after EndShort()")
	}

	if e1 == e2 {
		s.overlap++
	} else {
		if _, ok := s.seen[e1]; ok {
			delete(s.seen, e1)
			s.overlap++
		} else {
			s.seen[e1] = struct{}{}
		}

		if _, ok := s.seen[e2]; ok {
			delete(s.seen, e2)
			s.overlap++
		} else {
			s.seen[e2] = struct{}{}
		}
	}
	s.depth++
	s.wgt *= s.p
	s.rbo += (float64(s.overlap) / float64(s.depth)) * s.wgt
}

// EndShort indicates the end of the shorter of the two lists has been reached
func (s *State) EndShort() {
	s.shortDepth = s.depth
	s.shortOverlap = s.overlap
}

// UpdateUneven adds the entries from the longer list to the state
func (s *State) UpdateUneven(e int) {
	if s.shortDepth == -1 {
		panic("rbo: UpdateUneven() called without EndShort()")
	}

	if _, ok := s.seen[e]; ok {
		s.overlap++
		delete(s.seen, e)
	}

	s.depth++
	s.wgt *= s.p

	s.rbo += (float64(s.overlap) / float64(s.depth)) * s.wgt

	/* Extrapolation of overlap at end of short list */

	s.rbo += (float64(s.shortOverlap*(s.depth-s.shortDepth)) / float64(s.depth*s.shortDepth)) * s.wgt
}

// CalcExtrapolated returns the calculated extrapolated RBO
func (s *State) CalcExtrapolated() float64 {
	pl := math.Pow(s.p, float64(s.depth))

	//    assert(fabs((s.wgt * s.p) / (1 - s.p) - p_l) < 0.00001);
	if s.shortDepth == -1 {
		s.EndShort()
	}

	return s.rbo + (float64(s.overlap-s.shortOverlap)/float64(s.depth)+(float64(s.shortOverlap)/float64(s.shortDepth)))*pl
}

// Calculate is a helper function which computes the RBO for two integer arrays
func Calculate(s, t []int, p float64) float64 {

	st := New(p)

	if len(t) < len(s) {
		s, t = t, s
	}

	for i, v := range s {
		st.Update(v, t[i])
	}

	st.EndShort()

	if len(t) > len(s) {
		for _, v := range t[len(s):] {
			st.UpdateUneven(v)
		}
	}

	return st.CalcExtrapolated()
}
