package common

func View[T any](table string, in <-chan Mixed) <-chan Notification[T] {
	want := ToKey(table)
	out := make(chan Notification[T], 64)
	go func() {
		defer close(out)
		for m := range in {
			if m.Table != want { continue }
			var n Notification[T]
			n.Table, n.Op = m.Table, m.Op
			if v := coercePtr[T](m.New); v != nil { n.New = v }
			if v := coercePtr[T](m.Old); v != nil { n.Old = v }
			out <- n
		}
	}()
	return out
}

func coercePtr[T any](x any) *T {
	if x == nil { return nil }
	if p, ok := x.(*T); ok { return p }
	if v, ok := x.(T); ok { return &v }
	return nil
}
