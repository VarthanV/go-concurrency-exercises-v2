Here's a professional and comprehensive `README.md` file for your codebase:

---

# 🕸️ Naive Limited URL Fetcher

A **Go-based concurrent URL fetcher** with a *naive rate-limiting mechanism*. It allows a configurable number of concurrent tasks (fetches) and resets the concurrency limit at a fixed interval.

## 🧠 Key Features

* 🔁 Periodic reset of concurrency limits
* 📦 Supports adding dynamic tasks (URLs)
* ⚙️ Customizable concurrency and reset interval
* 🧵 Uses goroutines and channels to manage concurrency
* 🔒 Simple semaphore-based throttling logic

---

## 📦 Use Case

If you're fetching thousands of URLs and want to limit the number of concurrent HTTP requests (e.g., to avoid rate limiting by remote servers), this tool provides a minimal and customizable approach.

---

## 🚀 Getting Started

### 🧱 Prerequisites

* Go 1.18+
* Internet connection (to test URL fetching)

### 📥 Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/naive-limited-runner.git
cd naive-limited-runner
```

Build and run:

```bash
go run main.go
```

---

## 🛠️ How It Works

1. **`naiveLimitedRunner`**:

   * Holds a channel-based semaphore (`sem`) to manage concurrency.
   * Has a `resetLimit` goroutine that periodically refills the semaphore.
   * Runs `doAction` workers that consume tasks and fetch URLs.

2. **Concurrency Control**:

   * Only `concurrencyLimit` tasks can run in parallel.
   * The limit resets every `resetInterval`.

3. **Task Management**:

   * A task is just a struct with a `URL`.
   * All tasks are pushed into a generator channel, and workers consume from it.

---

## 📄 Example

Here’s how to use it in `main`:

```go
func main() {
	ctx := context.Background()

	runner := NewNaive(5, 10*time.Second) // 5 concurrent fetches, reset every 10s

	tasks := []Task{
		{URL: "https://example.com"},
		{URL: "https://golang.org"},
		// Add more URLs...
	}

	runner.AddTasks(tasks)
	if err := runner.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
```

---

## 📘 API Overview

### `NewNaive(limit int, resetInterval time.Duration) *naiveLimitedRunner`

Creates a new runner with the specified concurrency limit and reset interval.

### `AddTasks(tasks []Task)`

Adds a list of tasks (URLs) to be fetched.

### `Start(ctx context.Context) error`

Begins processing tasks concurrently while respecting rate limits.

---

## ⚠️ Notes

* The implementation is **naive**: it does not use a token bucket or leaky bucket algorithm.
* Currently fetches only URLs via HTTP `GET`. You can modify `fetch` to handle other HTTP methods or payloads.

---

## 🧪 Example Output

```
Spawning workers 10000
Allowed concurrency limit 5
Fetching url https://example.com
Status code 200
...
resetting the concurrency limit to 5
```
---
