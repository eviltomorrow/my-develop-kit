package gitlib

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	upgradeTimeout = 2 * time.Minute
	goroutines     = 3
	retryCount     = 1
)

// SetUpgradeTimeout per action timeout set
func SetUpgradeTimeout(timeout time.Duration) {
	upgradeTimeout = timeout
}

// SetGoroutines set goroutines
func SetGoroutines(num int) {
	goroutines = num
}

// SetRetryCount set retry count
func SetRetryCount(count int) {
	retryCount = count
}

// UpgradeGitLibLatest upgrade gitlib to latest
func UpgradeGitLibLatest(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("Not a folder")
	}

	var chpath = FindGitDir(path)
	var data = newCache(chpath)

	var waiter sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		waiter.Add(1)
		go func(c *cache) {
		loop:
			for {
				select {
				case current, ok := <-data.path:
					if !ok {
						break loop
					}

					var record = &record{
						Status:       failure,
						Count:        0,
						GitLocalPath: current,
						Timstamp:     time.Now(),
					}
					result, err := GitRemoteWithV(current, upgradeTimeout)
					if err != nil {
						record.Err = fmt.Errorf("Get remote path failure, nest error: %v", err)
						break
					} else {
						record.GitRemotePath = result
					}
					_, err = GitFetch(current, upgradeTimeout)
					if err != nil {
						record.Err = fmt.Errorf("Fetch remote lib failure, nest error: %v", err)
						break
					}
					result, err = GitCurrentBranch(current, upgradeTimeout)
					if err != nil {
						record.Err = fmt.Errorf("Get current branch failure, nest error: %v", err)
						break
					} else {
						record.Branch = result
					}
					_, err = GitPull(current, upgradeTimeout)
					if err != nil {
						record.Err = fmt.Errorf("Pull remote lib failure, nest error: %v", err)
						break
					}
					record.Status = success
					data.records <- record
				}
			}
			waiter.Done()
		}(data)
	}

	var signal = make(chan struct{})
	go func() {
		for record := range data.records {
			fmt.Println(record.String())
		}
		signal <- struct{}{}
	}()
	waiter.Wait()
	close(data.records)
	<-signal
	return nil
}
