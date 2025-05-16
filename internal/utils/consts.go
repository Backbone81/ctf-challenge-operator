package utils

import "time"

// DurationEpsilon is the difference between two times which we still see as equal. This is useful when tests want
// to check if the requeue is at the correct point in time, while still allowing for a bit of deviation because of
// host load.
const DurationEpsilon = 3 * time.Second
