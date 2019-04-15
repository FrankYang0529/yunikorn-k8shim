/*
Copyright 2019 The Unity Scheduler Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
    "context"
    "sync"
)

type DoWorkPieceFunc func(piece int)

// Copy from k8s.
func ParallelizeUntil(ctx context.Context, workers, pieces int, doWorkPiece DoWorkPieceFunc) {
    var stop <-chan struct{}
    if ctx != nil {
        stop = ctx.Done()
    }

    toProcess := make(chan int, pieces)
    for i := 0; i < pieces; i++ {
        toProcess <- i
    }
    close(toProcess)

    if pieces < workers {
        workers = pieces
    }

    wg := sync.WaitGroup{}
    wg.Add(workers)
    for i := 0; i < workers; i++ {
        go func() {
            defer wg.Done()
            for piece := range toProcess {
                select {
                case <-stop:
                    return
                default:
                    doWorkPiece(piece)
                }
            }
        }()
    }
    wg.Wait()
}