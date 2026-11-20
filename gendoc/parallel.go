package main

// parallel.go — concurrency plumbing for the doc generation pipeline.
//
// The generator renders ~2300 independent .md files (one per resource /
// data source / framework reference). Rendering each file touches only
// its own inputs and its own output path, so the per-document jobs are
// run through a bounded worker pool (see runDocFuncs).
//
// Failure semantics: the sequential generator printed a [FAIL!] message
// and called os.Exit(1) at the first problem. Under concurrency the
// abort is signalled with panic(docFailure{}) (see failf / fail) and
// recovered in runDocFuncs, which then exits 1 immediately — the
// observable behaviour (message + exit status + timing of the abort)
// stays identical to the sequential version.

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

// docWorkersEnv allows overriding the worker count for local
// experimentation, e.g. GENDOC_WORKERS=1 reproduces the fully
// sequential rendering order.
const docWorkersEnv = "GENDOC_WORKERS"

// docFailure is the panic payload marking a controlled doc-generation
// abort (see failf). Any other panic is a real bug and propagates.
type docFailure struct{}

// failf prints a [FAIL!] message (same format as the sequential
// generator) and aborts the current doc job.
func failf(format string, v ...interface{}) {
	message("[FAIL!]"+format, v...)
	panic(docFailure{})
}

// messageMu serialises message() output so concurrent workers cannot
// interleave partial log lines.
var messageMu sync.Mutex

// docWorkerCount returns the worker pool size: GENDOC_WORKERS when set
// to a positive integer, otherwise runtime.NumCPU().
func docWorkerCount() int {
	if v := os.Getenv(docWorkersEnv); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	if n := runtime.NumCPU(); n > 0 {
		return n
	}
	return 1
}

// runDocFuncs executes the doc jobs through a bounded worker pool. Jobs
// must be independent (distinct output files, no shared mutable state).
// A job aborts by panicking with docFailure (failf). On the first
// failure no further jobs are dispatched; in-flight jobs are allowed to
// finish their current file (so no .md is left truncated), and the
// process then exits 1 — preserving the sequential generator's
// "any single failure aborts the run" contract.
func runDocFuncs(jobs []func()) {
	workers := docWorkerCount()
	if workers <= 1 || len(jobs) <= 1 {
		for _, job := range jobs {
			runDocJob(job)
		}
		return
	}

	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		failed bool
	)
	sem := make(chan struct{}, workers)
	for _, job := range jobs {
		mu.Lock()
		stop := failed
		mu.Unlock()
		if stop {
			break // a previous job failed: dispatch nothing further
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(job func()) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := runDocJobErr(job); err != nil {
				mu.Lock()
				failed = true
				mu.Unlock()
			}
		}(job)
	}
	wg.Wait()
	if failed {
		os.Exit(1)
	}
}

// runDocJob runs one job on the sequential path, converting a
// docFailure panic into an immediate exit(1).
func runDocJob(job func()) {
	if err := runDocJobErr(job); err != nil {
		os.Exit(1)
	}
}

// runDocJobErr runs one job, converting a docFailure panic into a
// returned error so the pool can stop dispatching and exit cleanly.
func runDocJobErr(job func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(docFailure); ok {
				err = fmt.Errorf("doc job failed")
				return
			}
			panic(r)
		}
	}()
	job()
	return nil
}
