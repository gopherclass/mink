package mink

func insertionSortSlice(t []Element, f func(a Element, b Element) bool) []Element {
	rt := make([]Element, len(t))
	copy(rt, t)
	for i := 1; i < len(rt); i++ {
		for j := 0; j < i; j++ {
			if f(rt[i], rt[j]) {
				temp := rt[i]
				for k := i; k > j; k-- {
					rt[k] = rt[k-1]
				}
				rt[j] = temp
			}
		}
	}
	return rt
}

func mergeSortSlice(t []Element, f func(a Element, b Element) bool) []Element {
	if len(t) <= 16 {
		return insertionSortSlice(t, f)
	}
	m := len(t) / 2
	l := mergeSortSlice(t[:m], f)
	r := mergeSortSlice(t[m:], f)
	rt := make([]Element, len(t))

	for li, ri, rti := 0, 0, 0; ; {
		if li < len(l) && ri < len(r) {
			switch f(l[li], r[ri]) {
			case true:
				rt[rti] = l[li]
				li++
			case false:
				rt[rti] = r[ri]
				ri++
			}
			rti++
			continue
		}
		if li < len(l) && ri >= len(r) {
			rt[rti] = l[li]
			li++
			rti++
			continue
		}
		if li >= len(l) && ri < len(r) {
			rt[rti] = r[ri]
			ri++
			rti++
			continue
		}
		if li >= len(l) && ri >= len(r) {
			break
		}
	}

	return rt
}

func asyncMergeSortSlice(t []Element, f func(a Element, b Element) bool) <-chan []Element {
	ch := make(chan []Element)
	go func() {
		if len(t) <= 65536 {
			ch <- mergeSortSlice(t, f)
			return
		}
		m := len(t) / 2
		lc := asyncMergeSortSlice(t[:m], f)
		rc := asyncMergeSortSlice(t[m:], f)
		l := <-lc
		r := <-rc
		rt := make([]Element, 0, len(t))

		for li, ri, rti := 0, 0, 0; ; {
			if li < len(l) && ri < len(r) {
				switch f(l[li], r[ri]) {
				case true:
					rt[rti] = l[li]
					li++
				case false:
					rt[rti] = r[ri]
					ri++
				}
				rti++
				continue
			}
			if li < len(l) && ri >= len(r) {
				rt[rti] = l[li]
				li++
				rti++
				continue
			}
			if li >= len(l) && ri < len(r) {
				rt[rti] = r[ri]
				ri++
				rti++
				continue
			}
			if li >= len(l) && ri >= len(r) {
				break
			}
		}

		ch <- rt
		return
	}()
	return ch
}
