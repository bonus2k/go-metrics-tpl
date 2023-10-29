package utils

import (
	"fmt"
	"os"
	"time"
)

func RetryAfterError(f func() error) (err error) {
	attempts := 3
	sleep := 1 * time.Second
	for i := 1; ; i++ {
		now := time.Now()
		err = f()
		if err == nil {
			return
		}

		if i > attempts {
			break
		}

		time.Sleep(sleep)
		fmt.Fprintf(os.Stdout, "[%s]Attempt %d, retrying after error: %v\n", time.Since(now).Round(time.Millisecond), i, err)
		sleep = sleep + 2*time.Second

	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
