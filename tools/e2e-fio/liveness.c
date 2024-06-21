// This file implements functions to monitor liveness of fio target files or block devices
//
// For each monitored target, a thread (reader thread) is created which periodically
// (read interval) reads the same block from a target and records the time on success.
//
// liveness of a device is determined by comparing the delta between current time and the time of
// the last successful read on the target and a timeout value (seconds).
//
// To work correctly, the timeout value should be at least 3 times larger than the read
// intervals.
//
#define _GNU_SOURCE

#include <assert.h>
#include <errno.h>
#include <fcntl.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <sys/types.h>
#include <time.h>
#include <unistd.h>
#include <pthread.h>
#include "liveness.h"

#define BLOCKSIZE (4096)
#define ALIGNMENT BLOCKSIZE

// read buffer for all threads.
// for liveness checks the contents of read operations are irrelevant,
// only success or failure of the read ops matters.
static char read_buffer[BLOCKSIZE] __attribute__((__aligned__(ALIGNMENT)));

// reader thread struct (this is the only shared data that reader threads access)
typedef struct {
    // time when last read was completed
    // produced by reader thread
    // consumed by monitoring code
    volatile time_t tv_sec;
    // stop flag for reader thread termination
    // produced by controlling/monitoring code.
    // consumed by reader thread
    // used to signal reader thread that it should complete.
    volatile bool stop;
    // produced reader thread
    // consumed by monitoring/controlling code,
    // used to avoid blocking pthread_join call.
    volatile bool done;
    // produced reader thread
    // consumed by monitoring code.
    volatile bool read_ok;
    // file descriptor to read
    // produced by add_target function
    // consumed by reader thread
    int fd;
    // read op interval
    // produced by add_target function
    // consumed by reader thread
    unsigned read_interval;
} reader_context;

// monitoring struct (note: not accessed by reader thread)
typedef struct s_reader_monitor {
    // reader struct
    reader_context reader;
    // reader thread sampling interval, this should be larger than
    // the largest reader thread interval by at least a factor of 3
    time_t timeout;
    // reader thread id
    pthread_t thread_id;
    // linked list pointer
    struct s_reader_monitor* next;
} reader_monitor_context;

// liveness struct
struct struc_liveness_context {
    // linked list of readers
    reader_monitor_context* rd_ctxts;
};

// reader thread function prototype
static void* target_reader(void* vargp);

// add a target to be checked for liveness
bool add_target(liveness_context lc, const char* targetpath, unsigned read_interval, time_t timeout) {
    if (lc == NULL) {
        return false;
    }

    int fd = open(targetpath, O_RDONLY | O_DIRECT);
    if (fd >= 0) {
        reader_monitor_context *rc = calloc(1, sizeof(*rc));
        if (rc != NULL) {
            rc->reader.fd = fd;
            rc->reader.read_interval = read_interval;
            rc->reader.stop = false;
            rc->timeout = timeout;
            rc->next = lc->rd_ctxts;
            lc->rd_ctxts = rc;
            printf("%p: liveness check for %s at intervals of %u seconds, timeout %ld seconds\n", rc, targetpath, read_interval, timeout);
            return true;
        } else {
            perror("failed to allocate memory");
            close(fd);
        }
    } else {
        perror("failed to open target");
    }
    return false;
}

// stop liveness checking - signal reader threads to stop.
void stop_liveness_checks(liveness_context lc) {
    if (lc == NULL) {
        return;
    }
    // signal all the reader threads to stop
    // when all the reader threads have stopped,
    // then the monitor thread will stop.
    reader_monitor_context* rc;
    for(rc = lc->rd_ctxts; rc != NULL; rc = rc->next) {
        rc->reader.stop = true;
    }
}

// start the liveness checks - start reader threads
bool start_liveness_checks(liveness_context lc) {
    bool ok = false;
    if (lc != NULL) {
        reader_monitor_context* rc = lc->rd_ctxts;
        while(rc != NULL) {
            if (rc->thread_id == 0) {
                pthread_create(&rc->thread_id, NULL, target_reader, rc);
                rc = rc->next;
            }
            ok = ok && (rc->thread_id != 0);
        }
    }
    return ok;
}

// free reader_monitor_context linked list, recursive.
// if a reader thread is not "done" then associated
// memory is not freed, and all ancestors in the
// linked list are not freed either.
// This is fine as when the container terminates the
// memory will be returned to the system.
static bool free_reader_list(reader_monitor_context * rc) {
    bool freed = true;
    if (rc != NULL) {
        if (free_reader_list(rc->next)) {
            if (rc->reader.done) {
                printf("%p: freed\n", rc);
                free(rc);
            } else {
                // the reader thread is still running
                // elect to keep the memory allocated
                freed = false;
            }
        }
    }
    return freed;
}

// destroy liveness context
// - close filedescriptors
// - free reader memory if reader thread has completed.
// - free liveness context if all readers have completed
// returns NULL on success otherwise the pointer to the
// liveness context.
// Note: reader threads may be "stuck" so we do not do pthread_join,
// but merely set the stop flag.
// If all readers were not determined to have stopped (done flag not set)
// then the context is not destroyed.
// For the use case where this program is run within a container,
// this does not constitute a memory leak - the memory will be
// returned to the system when the container terminates.
// Crucially the thread will not "scribble" on memory
// which has been freed - if it ever revives.
liveness_context destroy_liveness_context(liveness_context lc) {
    if (lc == NULL) {
        return lc;
    }
    for(reader_monitor_context* rc = lc->rd_ctxts; rc != NULL; rc = rc->next) {
        if (rc->reader.fd >= 0) {
            close(rc->reader.fd);
            rc->reader.fd = -1;
        }
    }
    for(reader_monitor_context* rc = lc->rd_ctxts; rc != NULL; rc = rc->next) {
        if (!rc->reader.done) {
            stop_liveness_checks(lc);
        }
        sleep(rc->reader.read_interval * 2);
    }
    if (free_reader_list(lc->rd_ctxts)) {
        free(lc);
        lc = NULL;
    }
    return lc;
}

// Create and initialise a liveness context.
liveness_context make_liveness_context(void) {
    liveness_context lc = calloc(sizeof(struct struc_liveness_context), 1);
    return lc;
}

// return liveness calculated based on current time.
// Note: if any reader thread has been stopped (by setting the stop flag),
// this function will return a false "not live" result.
// returns true if nothing to check
bool liveness_check(liveness_context lc) {
    bool live = true;
    if (lc != NULL) {
        struct timespec ts;
        if (clock_gettime(CLOCK_MONOTONIC, &ts) != 0) {
            perror("clock_gettime (liveness_check)");
            return false;
        }
        for(reader_monitor_context* rc = lc->rd_ctxts; rc != NULL; rc = rc->next) {
            live = live && ((ts.tv_sec - rc->reader.tv_sec) <= rc->timeout);
        }
    }
    return live;
}

// dump liveness state.
void dump_liveness_context(liveness_context lc) {
    if (lc == NULL) {
        return;
    }
    struct timespec ts;
    if (clock_gettime(CLOCK_MONOTONIC, &ts) != 0) {
        perror("clock_gettime (dump_liveness_context)");
        return;
    }
    for(reader_monitor_context* rc = lc->rd_ctxts; rc != NULL; rc = rc->next) {
        bool read_time_delta = (ts.tv_sec - rc->reader.tv_sec) <= rc->timeout;
        printf("%p: read_time_delta:%d, reader:tv_sec=%ld read_ok=%d stop=%d done=%d,fd=%d\n",
                rc, read_time_delta,
                rc->reader.tv_sec, rc->reader.read_ok, rc->reader.stop, rc->reader.done, rc->reader.fd);
    }
}

// reader thread
static void *target_reader(void *vargp)
{
    reader_context *rc = vargp;
    const int block_no = 1;
    int err = 0;
    bool ok = true;
    struct timespec ts;
    err = clock_gettime(CLOCK_MONOTONIC, &ts);
    if (err != 0) {
        perror("clock_gettime");
        rc->done = true;
        rc->tv_sec = -1;
        return NULL;
    }
    rc->tv_sec = ts.tv_sec;
    sleep(1);

    for(; !rc->stop ;) {
        err = lseek(rc->fd, (int64_t)(block_no * BLOCKSIZE), 0);
        if (err >= 0) {
            int64_t num_bytes;
            num_bytes = read(rc->fd, read_buffer, BLOCKSIZE);
            if (num_bytes == BLOCKSIZE) {
                err = clock_gettime(CLOCK_MONOTONIC, &ts);
                if (err != 0) {
                    rc->done = 0;
                    rc->tv_sec = -1;
                    perror("clock_gettime (L)");
                } else {
                    rc->tv_sec = ts.tv_sec;
                    ok = true;
                }
            } else {
                if (ok) {
                    // emit error message on transitions from success to failure
                    perror("read failure");
                    fprintf(stderr, "%p: could not read at block %d, wanted %d, got %ld\n",
                            rc, block_no, BLOCKSIZE, num_bytes);
                }
                ok = false;
            }
        } else {
            if (ok) {
                // emit error message on transitions from success to failure
                perror("seek failure");
                fprintf(stderr, "%p: could not seek at block %d\n",
                        rc, block_no);
            }
            ok = false;
        }
        rc->read_ok = ok;
        sleep(rc->read_interval);
    }
    rc->done = true;
    printf("%p: end of liveness check\n", rc);
    return NULL;
}
