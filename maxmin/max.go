package maxmin

import "golang.org/x/exp/constraints"

func Max[V constraints.Ordered](in ...V) (out V) {
	if len(in) == 0 {
		return out
	}

	for _, v := range in {
		if v >= out {
			out = v
		}
	}
	return out
}

func Min[V constraints.Ordered](in ...V) (out V) {
	if len(in) == 0 {
		return out
	}

	out = in[0]

	for _, v := range in {
		if v <= out {
			out = v
		}
	}
	return out
}
