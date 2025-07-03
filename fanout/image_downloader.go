package fanout

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// Task is to write a fanout image downloader that takes a list of urls from the
// generator and downloads the images concurrently using a fanout pattern.
// Assuming that we have a concurrency limit of 100 , Use a sempahore to limit the
// concurrency of the image downloader to 100. Use a worker pool to download the images.

const (
	totalImages = 5000
)

type APIResult struct {
	AlbumID      int    `json:"albumId"`
	ID           int    `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type Result struct {
	APIResult
	Error error
}

type Work struct {
	ID  int
	wg  *sync.WaitGroup
	In  <-chan string
	out chan<- Result
}

type downloader struct {
	sem              chan struct{}
	concurrencyLimit int
	total            atomic.Int32
}

func NewDownloader(concurrencyLimit int) *downloader {
	return &downloader{
		sem:              make(chan struct{}, concurrencyLimit),
		concurrencyLimit: concurrencyLimit,
	}
}

func (d *downloader) GenerateImageURL(ctx context.Context,
	total int) <-chan string {
	outStream := make(chan string)
	go func() {
		defer close(outStream)
		for i := range total {
			select {
			case <-ctx.Done():
				return
			case outStream <- fmt.Sprintf("https://jsonplaceholder.typicode.com/photos/%d", i+1):
			}
		}
	}()
	return outStream
}

func (d *downloader) work(ctx context.Context, url string) Result {
	var (
		result = APIResult{}
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logrus.Error("error in creating new request ", err)
		return Result{
			Error: err,
		}
	}

	logrus.Info("Fetching url ", url)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error("error in doing request ", err)
		return Result{Error: err}
	}

	logrus.Info("Received status ", res.StatusCode)
	switch res.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(res.Body)
		if err != nil {
			logrus.Error("error in reading body ", err)
			return Result{Error: err}
		}
		err = json.Unmarshal(body, &result)
		if err != nil {
			logrus.Error("error in reading result ", err)
			return Result{Error: err}
		}
		return Result{APIResult: result}
	default:
		err = fmt.Errorf("status Code : %d", res.StatusCode)
		logrus.Error(err)
		return Result{Error: err}
	}
}

func (d *downloader) Worker(ctx context.Context, w *Work) {
	defer w.wg.Done()
	logrus.Infof("Worker %d reporting sir 🫡", w.ID)
	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-w.In:
			if !ok {
				return
			}
			d.sem <- struct{}{}
			logrus.Infof("Worker %d working ⌛ ", w.ID)
			w.out <- d.work(ctx, url)
			<-d.sem
			logrus.Infof("Worked %d released ", w.ID)
		}
	}
}

func (d *downloader) Driver(ctx context.Context) {
	var (
		wg sync.WaitGroup
	)

	out := make(chan Result)
	in := d.GenerateImageURL(ctx, totalImages)

	for i := range d.concurrencyLimit {
		select {
		case <-ctx.Done():
			return
		default:
			wg.Add(1)
			go d.Worker(ctx, &Work{
				ID:  i,
				wg:  &wg,
				In:  in,
				out: out,
			})
		}
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	go func() {
		ticker := time.NewTicker(time.Second * 1)
		f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				logrus.Info("#############################")
				_, err = f.Write(fmt.Appendf(nil, "Number of goroutines: %d Downloaded Count: %d \n", runtime.NumGoroutine(), d.total.Load()))
				if err != nil {
					logrus.Error("error in writing to file ", err)
				}
				logrus.Info("#############################")

			}
		}
	}()

	for rslt := range out {
		if rslt.Error != nil {
			logrus.Error("error in getting result ", rslt.Error)
			continue
		}
		logrus.Info("Received rslt is ", rslt.APIResult.ThumbnailURL)
		d.total.Add(1)
	}
	logrus.Fatal("Total images downloaded: ", d.total.Load())
}
