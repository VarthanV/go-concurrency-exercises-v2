package runner

/*
	Task is to create a controlled runner that runs concurrently tasks with
	a specific time period , When the task exceeds the time period it should be vleaned
*/

type runner struct {
	// channel to signal the completion of tasks
	sem chan struct{}
	// represents current active users
	users map[string]struct{}
	// limits the number of tasks that can run concurrently
	concurrencyLimit int
}

func New(concurrencyLimit int) *runner {
	return &runner{
		sem:              make(chan struct{}, 100), // Set a concurrency limit
		users:            make(map[string]struct{}),
		concurrencyLimit: concurrencyLimit,
	}
}
func (r *runner) AddUser(user string) {
	r.users[user] = struct{}{}
}
