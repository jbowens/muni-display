package predictions

import "time"

type updatePredicate func(now time.Time, lastUpdated time.Time, m *Module) *bool

func composite(preds ...updatePredicate) updatePredicate {
	return func(now time.Time, lastUpdated time.Time, m *Module) *bool {
		for _, p := range preds {
			if res := p(now, lastUpdated, m); res != nil {
				return res
			}
		}
		return nil
	}
}

func interval(duration time.Duration) updatePredicate {
	return func(now time.Time, lastUpdated time.Time, m *Module) *bool {
		res := now.After(lastUpdated.Add(duration))
		return &res
	}
}

func onWeekends(inner updatePredicate) updatePredicate {
	return func(now time.Time, lastUpdated time.Time, m *Module) *bool {
		weekday := now.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			return nil
		}

		return inner(now, lastUpdated, m)
	}
}

func inMornings(inner updatePredicate) updatePredicate {
	return func(now time.Time, lastUpdated time.Time, m *Module) *bool {
		h := now.Hour()
		if h >= 8 && h < 11 {
			return inner(now, lastUpdated, m)
		}

		return nil
	}
}

func atNight(inner updatePredicate) updatePredicate {
	return func(now time.Time, lastUpdated time.Time, m *Module) *bool {
		h := now.Hour()
		if h == 23 || h < 7 {
			return inner(now, lastUpdated, m)
		}
		return nil
	}
}
