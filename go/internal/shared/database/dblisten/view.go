package dblisten

func View[T any](table string, in <-chan Mixed) <-chan Notification[T] {
	want := toKey(table)
	out := make(chan Notification[T], 64)
	go func() {
		defer close(out)
		for m := range in {
			if m.Table != want {
				continue
			}
			var n Notification[T]
			n.Table, n.Op = m.Table, m.Op
			if m.New != nil {
				if v, ok := m.New.(*T); ok {
					n.New = v
				} else if v, ok := m.New.(T); ok {
					n.New = &v
				}
			}
			if m.Old != nil {
				if v, ok := m.Old.(*T); ok {
					n.Old = v
				} else if v, ok := m.Old.(T); ok {
					n.Old = &v
				}
			}
			out <- n
		}
	}()
	return out
}
