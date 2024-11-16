package lock

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math"
	"sync"
	"testing"
	"time"
)

func TestSubLock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "secret_redis",
		DB:       0,
	})

	var key1, key2, key3, keyout int

	defer func() {
		fmt.Printf("dolock: %d, noticeLock: %d, relock:%d, timeouterr: %d", key1, key2, key3, keyout)
	}()

	wg := sync.WaitGroup{}
	defer wg.Wait()

	for i1 := 0; i1 < 100; i1++ {
		wg.Add(1)

		key := fmt.Sprintf("test:%d", i1)

		go func() {
			defer wg.Done()
			lc := NewSubLock(context.Background(), client, key, SetExpire(10*time.Second)).Lock()

			for i := 0; i < 97; i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()

					lc1 := NewSubLock(context.Background(), client, key, SetExpire(10*time.Second)).Lock()
					defer lc1.Release()

					if lc1.Error() != nil {

						if errors.Is(lc1.Error(), FailTimeoutErr) {
							keyout++
						}

						return
					}

					switch lc1.status {
					case DoLockOk:
						key1++
					case NoticeOk:
						key2++
					case ReLockOk:
						key3++
					}

				}()
			}

			lc.Release()

		}()

	}

}

func TestName(t *testing.T) {
	var a = 3

	t.Log(int(math.Ceil(float64(a) / 2)))

}
