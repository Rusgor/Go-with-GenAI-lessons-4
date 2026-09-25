package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunPoolAI_TimeoutReturnsDeadlineExceeded(t *testing.T) {
	withTimeout(t, 2*time.Second, func() {
		jobs := make(chan Job, 1)
		jobs <- Job{
			ID: "cooperative-timeout",
			Fetch: func(ctx context.Context) (int, error) {
				<-ctx.Done()
				return 0, ctx.Err()
			},
		}
		close(jobs)

		results := RunPool(jobs, 3, 25*time.Millisecond)
		var got []Result
		for result := range results {
			got = append(got, result)
		}

		if len(got) != 1 {
			t.Fatalf("received %d results, want 1", len(got))
		}
		if got[0].JobID != "cooperative-timeout" {
			t.Errorf("JobID = %q, want cooperative-timeout", got[0].JobID)
		}
		if !errors.Is(got[0].Err, context.DeadlineExceeded) {
			t.Errorf("Err = %v, want context.DeadlineExceeded", got[0].Err)
		}
	})
}
