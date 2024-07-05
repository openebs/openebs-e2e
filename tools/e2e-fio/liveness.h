#ifndef _e2e_fio_liveness_h
#define _e2e_fio_liveness_h

#include <time.h>

/* Note: none of these functions are thread safe */
typedef struct struc_liveness_context *liveness_context;
// Create and initialise a liveness context.
liveness_context make_liveness_context(void);
// destroy liveness context and associated resources,
// returns NULL on success
liveness_context destroy_liveness_context(liveness_context lc);
// add a target to be checked for liveness
bool add_target(liveness_context lc, const char* devicepath, unsigned read_interval, time_t timeout);
// return liveness of target, true if no liveness checks have been created.
bool liveness_check(liveness_context lc);
// start the liveness checks - start reader threads
bool start_liveness_checks(liveness_context lc);
// stop liveness checking - signal reader threads to stop.
void stop_liveness_checks(liveness_context lc);
// dump liveness state to stdout
void dump_liveness_context(liveness_context lc);

#endif // _e2e_fio_liveness_h

