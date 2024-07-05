#define _GNU_SOURCE

#include <stdio.h>
#include <signal.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/stat.h>
#include <sys/vfs.h>
#include <fcntl.h>
#include <errno.h>
#include <time.h>
#include "liveness.h"
#include "e2e_fio_version.h"

const char *workspace_path = "./workspace";
char session_id[128+1];


int parse_cmds(const char* arg);
static liveness_context liveness_ctx;
unsigned post_op_sleep = 0;

// convert whole strings to unsigned integer or fail
bool strtounsigned(const char* chars, unsigned *p_val) {
    char *end;
    *p_val = '\0';
    unsigned long v = strtoul(chars, &end, 10);
    *p_val = (unsigned)v;
    return *end == 0;
}

// convert whole strings to long integer or fail
bool strtolong(const char* chars, long *p_val) {
    char *end;
    *p_val = 0;
    long v = strtol(chars, &end, 10);
    *p_val = v;
    return *end == 0;
}

// convert whole strings to size_t or fail
bool strtosize_t(const char* chars, size_t *p_val) {
    char *end;
    *p_val = 0;
    size_t v = strtoul(chars, &end, 10);
    *p_val = v;
    return *end == 0;
}


#define NOT_A_KEYWORD   -1
#define KW_EOL          -2
/*
 * array of key words and matching enum.
 * position of string in the array should match symbol position in enum
 */
const char* cmd_keywords[] = {
    "--", "---", "&&", ";",
    "sleep", "segfault", "sigterm", "makefile", "exitv", "liveness",
    "zerofill", "sessionId", "postopsleep", "filesize",
};
enum {
    KW_FIO_IMPLICIT,
    KW_TASK_START,
    KW_TASK_WAIT,
    KW_TASK_END,

    KW_SLEEP,
    KW_SEGFAULT,
    KW_SIGTERM,
    KW_MAKEFILE,
    KW_EXITV,
    KW_LIVENESS,
    KW_ZEROFILL,
    KW_SESSION_ID,
    KW_POST_OP_SLEEP,
    KW_FILESIZE,
};

/*
 *  array of file size key words and matching enum
 * position of string in the array should match symbol position in enum
 */
const char* filesize_keywords[] = {
    "availblockspercent", "availblockslessby", "bytes"};
enum {
    BLOCKS_PERCENT,
    BLOCKS_LESSBY,
    SIZE_BYTES,
};

/* struct for linked list of child processes */
typedef struct e2e_process {
    struct e2e_process* next;
    pid_t   pid;
    int     status;
    int     exitv;
    int     termsig;
    int     abnormal_exit;
    int     finished;
    char*   cmd;
    bool    wait;
} e2e_process;

/* struct for linked list of signal generations */
typedef struct e2e_signal_gen {
    struct e2e_signal_gen* next;
    time_t endtick;
    void (*func)(void);
} e2e_signal_gen;

/* struct for linked list of files created */
typedef struct e2e_fio_files {
    struct e2e_fio_files* next;
    char* filename;
} e2e_fio_files;


/* head of the linked list of child processes */
static e2e_process* proc_list = NULL;
static e2e_signal_gen* sig_gen_list = NULL;
static e2e_fio_files* files_list = NULL;

/* Create and run fio in child processes as defined in the list */
int start_proc(e2e_process* proc_ptr ) {
    proc_ptr->pid = fork();

    if ( 0 == proc_ptr->pid ) {
        /* Change working directory to avoid trivial collisions in file
         * space across multiple jobs
         */
        char wkspace[64];
        snprintf(wkspace, 64, "%s/%d", workspace_path, getpid());
        if (0 == mkdir(wkspace,0777) && (0 == chdir(wkspace)) ) {
            execl("/bin/sh", "sh", "-c", proc_ptr->cmd, NULL);
            printf("** execl %s failed %d **\n", proc_ptr->cmd, errno);
        } else {
            printf("** mkdir %s failed **\n", wkspace);
        }
        return errno;
    }

    printf("pid:%d\n%s\n", proc_ptr->pid, proc_ptr->cmd);
    fflush(stdout);
    if (proc_ptr->wait) {
        printf("waiting for pid:%d\n", proc_ptr->pid);
        fflush(stdout);
        if ( 0 > waitpid(proc_ptr->pid, &proc_ptr->status, 0 )) {
            perror("");
            printf("waitpid failed: pid:%d %d\n", proc_ptr->pid, proc_ptr->status);
        }
        fflush(stdout);
        proc_ptr->finished = 1;
        if (WIFEXITED(proc_ptr->status))  {
            proc_ptr->exitv = WEXITSTATUS(proc_ptr->status);
        }
        return proc_ptr->status;
    }
    return 0;
}

/*
 * parse command line arguments,
 *  - create an e2e_process struct,
 *  - append the process struct to the global list of e2e processes
 */
int do_exec(const char **argv, int count, bool wait) {
    e2e_process *proc_ptr = NULL;

    /* Tis' C so we do it the "hard way" */
    char *pinsert;
    size_t buflen = 0;

    puts("EXEC");
    system("date");

    /* 1. work out the size of the buffer required to copy the arguments.*/
    for(const char **argv_scan=argv; *argv_scan != NULL; ++argv_scan) {
        /* +1 for space delimiter */
        buflen += strlen(*argv_scan) + 1;
    }

    if (buflen == 0) {
        puts("ERROR: empty command line");
        return 1;
    }
    ++buflen;

    /* 2. create the process struct and initialise it */
    proc_ptr = calloc(sizeof(*proc_ptr), 1);
    if (proc_ptr == NULL) {
        puts("ERROR: failed to allocate memory for e2e_process");
        return 1;
    }
    proc_ptr->wait = wait;

    /* 3. allocate a 0 intialised buffer for the command line */
    proc_ptr->cmd = calloc(sizeof(unsigned char), buflen);
    if (proc_ptr->cmd == NULL) {
        free(proc_ptr);
        puts("ERROR: failed to allocate memory for command line");
        return 1;
    }

    /* 4. construct the command line, using strcat */
    pinsert = proc_ptr->cmd;
    for(; *argv != NULL; ++argv) {
        strcat(pinsert, *argv);
        pinsert += strlen(pinsert);
        *pinsert = ' ';
        ++pinsert;
    }

    /* 5. append the process to the list */
    {
        e2e_process** insert_proc = &proc_list;
        while (*insert_proc != NULL) {
            insert_proc = &(*insert_proc)->next;
        }
        *insert_proc = proc_ptr;
    }

    /* 6. start the process */
    return start_proc(proc_ptr);
}

/* Kill all processes as defined in the list */
void kill_procs(int signal) {
    for (e2e_process* proc_ptr = proc_list; NULL != proc_ptr; proc_ptr = proc_ptr->next) {
        if (proc_ptr->pid && !proc_ptr->finished) {
            printf("kill process %d\n", proc_ptr->pid);
            kill(proc_ptr->pid, signal);
        }
    }
}


/*
 * argv[0] = path
 * argv[1] = relative path to file (from argv[0])
 * argv[2] = <"availblockspercent" || "availblockslessby" || "bytes">
 * argv[3] = positive integer string; 0-100 for "availblockspercent", unbounded for "availblockslessby" and "bytes"
 */
int do_makefile(const char **argv, int count) {
    const char *argstrs[4];
    int indx;

    puts("MAKE FILE");
    if (count != 4) {
        puts("ERROR: invalid argument count");
        return 1;
    }

    for(indx=0; indx < sizeof(argstrs)/sizeof(argstrs[0]); ++indx, ++argv) {
        if (*argv == NULL || parse_cmds(*argv) != NOT_A_KEYWORD || 0 == strlen(*argv)) {
            printf("ERROR:argument missing %d for makefile command\n", indx);
            return 1;
        }
        argstrs[indx] = *argv;
    }

    {
        const char *fspath = argstrs[0];
        const char *file_rel_path = argstrs[1];
        size_t val;
        struct statfs fs_stat;
        char* fsfile = calloc(1, strlen(fspath) + strlen(file_rel_path) + 2);
        int op = NOT_A_KEYWORD;
        e2e_fio_files* eff = calloc(1, sizeof(*eff));

        if (fsfile == NULL) {
            printf("ERROR: malloc failed for %s/%s \n", fspath, file_rel_path);
            return 1;
        }

        if (eff == NULL) {
            printf("ERROR: malloc failed for %ld \n", sizeof(*eff));
            return 1;
        }

        eff->filename = fsfile;
        eff->next = files_list;
        files_list = eff;

        for (int ix=0; ix < sizeof(filesize_keywords)/sizeof(filesize_keywords[0]); ++ix) {
            if ( 0 == strcmp(argstrs[2], filesize_keywords[ix])) {
                op = ix;
                break;
            }
        }
        if (op == NOT_A_KEYWORD) {
            printf("ERROR: invalid directive %s\n", argstrs[2]);
        }

        if (!strtosize_t(argstrs[3], &val)) {
            printf("ERROR: parsing integer %s\n", argstrs[3]);
            return 1;
        }

        if ((op == BLOCKS_PERCENT) && val > 100) {
            puts("ERROR: percentage spec is > 100");
            return 1;
        }
        sprintf(fsfile, "%s/%s", fspath, file_rel_path);
        printf("unlinking fio file %s\n", fsfile);
        if (unlink(fsfile) != 0) {
            perror("unlink failed");
        }
        if (statfs(fspath, &fs_stat) == 0 ) {
            size_t use_bytes, use_blocks;
            printf("block size: %ld, avail blocks: %ld, avail bytes: %ld, total blocks: %ld, total bytes: %ld\n",
                    fs_stat.f_bsize, fs_stat.f_bavail,
                    fs_stat.f_bsize * fs_stat.f_bavail,
                    fs_stat.f_blocks,
                    fs_stat.f_bsize * fs_stat.f_blocks);

            switch(op) {
                case BLOCKS_PERCENT:
                    use_blocks = ((fs_stat.f_bavail * val)/100);
                    if (val == 100) {
                        /* sometimes fallocate fails when using all blocks,
                         * but works with if at least 1 block is left free.
                         */
                        use_blocks -= 1;
                    }
                    use_bytes = use_blocks * fs_stat.f_bsize;
                break;
                case BLOCKS_LESSBY:
                    if (val >= fs_stat.f_bavail) {
                        puts("ERROR: calculated blocks for file size <= 0");
                        return 1;
                    }
                    use_blocks = fs_stat.f_bavail - val;
                    use_bytes = use_blocks * fs_stat.f_bsize;
                break;
                case SIZE_BYTES:
                    use_blocks = (val + fs_stat.f_bsize - 1) / fs_stat.f_bsize;
                    use_bytes = val;
                break;
            }
            printf("fs: %s, file: %s blocks: %ld bytes: %ld\n", fspath, fsfile, use_blocks, use_bytes);
            printf("creat %s\n", fsfile);

            int fd = creat(fsfile, S_IRWXU| S_IRWXG | S_IRWXO);
            if (fd > -1) {
                printf("fallocate request %ld bytes, %ld MiB, %ld GiB\n",
                        use_bytes,
                        use_bytes / (1024 * 1024),
                        use_bytes / (1024 * 1024 * 1024)
                      );
//                int rv =  posix_fallocate(fd, 0, use_bytes);
                int rv =  fallocate(fd, 0, 0, use_bytes);
                if (rv != 0) {
                    perror("fallocate failed");
                }
                close(fd);
                return rv;
            } else {
                perror("creat failed");
            }
        } else {
            perror("statfs failed");
        }
    }
    return 1;
}

/*
 * argv[0] = filepath / device path
 */
int do_zerofill(const char **argv, int count) {
    int rv = 1;
    static char buffer[4096];
    puts("ZERO FILL");
    fflush(stdout);
    if (count == 1) {
        const char *fdpath = argv[0];
        off_t bytesize;
        int fd = open(fdpath, O_RDONLY);
        if (fd < 0) {
            printf("unable to access files %s\n", fdpath);
            return rv;
        }
        bytesize = lseek(fd, 0, SEEK_END);
        close(fd);
        if (bytesize > 0) {
            /* use unix utilities to zero fill - the behaviour is well understood */
            rv = snprintf(buffer, sizeof(buffer), "dd if=/dev/zero bs=%ld count=1 of=%s oflag=direct",
                    bytesize, fdpath);
            if (rv < sizeof(buffer)) {
                /* using time() to calculate elapsed time is good enough */
                time_t start = time(NULL);
                printf("executing %s\n", buffer);
                fflush(stdout);
                rv = system(buffer);
                system("sync");
                printf("elapsed time %ld seconds\n", time(NULL) - start);
                fflush(stdout);
                fflush(stderr);
            } else {
                puts("ERROR: command buffer is too small");
            }
        } else {
            printf("ERROR:unable to determine size of file %s\n", fdpath);
        }
    } else {
        printf("ERROR: invalid set of arguments\n");
    }
    return rv;
}

int do_sleep(const char **argv, int count) {
    puts("SLEEP");
    if (count == 1 ) {
        unsigned sleep_time;
        if (!strtounsigned(*argv, &sleep_time)) {
            printf("ERROR: invalid sleep time %s\n", *argv);
            return 1;
        }
        printf("sleeping %d seconds\n", sleep_time);
        fflush(stdout);
        sleep(sleep_time);
        return 0;
    } else {
        puts("ERROR: invalid argument count for sleep command");
    }
    return 1;
}

int do_post_op_sleep(const char **argv, int count) {
    puts("POST OP SLEEP");
    if (count == 1 ) {
        unsigned sleep_time;
        if (!strtounsigned(*argv, &sleep_time)) {
            printf("ERROR: invalid sleep time %s\n", *argv);
            return 1;
        }
        post_op_sleep = sleep_time;
        return 0;
    } else {
        puts("ERROR: invalid argument count for post op sleep command");
    }
    return 1;
}

/*
 * argv[0] = path to file/device
 */
int do_filesize(const char **argv, int count) {
    const char *argstrs[1];
    int indx;

    puts("FILESIZE");
    if (count != 1) {
        puts("ERROR: invalid argument count");
        return 1;
    }

    for(indx=0; indx < sizeof(argstrs)/sizeof(argstrs[0]); ++indx, ++argv) {
        if (*argv == NULL || parse_cmds(*argv) != NOT_A_KEYWORD || 0 == strlen(*argv)) {
            printf("ERROR:argument missing %d for filesize command\n", indx);
            return 1;
        }
        argstrs[indx] = *argv;
    }

    {
        const char *filepath = argstrs[0];
        int fd = open(filepath, O_RDONLY | O_NONBLOCK);
        if ( fd >= 0) {
            off_t fdsize = lseek(fd, 0, SEEK_END);
            close(fd);
            printf("\nJSON{\"fio_target_size\": %ld, \"path\": \"%s\"}\n", fdsize, filepath);
            return 0;
        }
    }
    return 1;
}


void generate_segfault(void) {
    kill_procs(SIGKILL);
    sleep(1);
    puts("Segfaulting now!");
    fflush(stdout);
    raise(SIGSEGV);
}

int do_segfault(const char **argv, int count) {
    struct timespec ts;
    unsigned sleep_time;

    puts("SEGFAULT");
    if (count != 1 ) {
        puts("ERROR: invalid argument count for segfault command");
        return 1;
    }
    if (!strtounsigned(*argv, &sleep_time)) {
        printf("ERROR: invalid sleep time %s\n", *argv);
        return 1;
    }
    e2e_signal_gen* sig_gen_ptr = calloc(sizeof(*sig_gen_ptr), 1);
    if (sig_gen_ptr == NULL) {
        perror("calloc failed");
        return 1;
    }

    clock_gettime(CLOCK_MONOTONIC, &ts);
    sig_gen_ptr->next = sig_gen_list;
    sig_gen_list = sig_gen_ptr;
    sig_gen_ptr->endtick = ts.tv_sec + sleep_time;
    sig_gen_ptr->func = generate_segfault;

    printf("Segfaulting after %d seconds\n", sleep_time);
    fflush(stdout);
    return 0;
}

static void send_sigterm(void) {
    puts("sending SIGTERM to all processes now!");
    fflush(stdout);
    for (e2e_process* proc_ptr = proc_list; NULL != proc_ptr; proc_ptr = proc_ptr->next) {
        if (proc_ptr->pid && !proc_ptr->finished) {
            printf("SIGTERM -> %d\n", proc_ptr->pid);
            if (0 == kill(proc_ptr->pid, SIGTERM)) {
                int retry_count = 5;
                // wait for SIGTERM to be effective
                while (0 == waitpid(proc_ptr->pid, &proc_ptr->status, WNOHANG)
                        && retry_count >= 0) {
                    if (WIFSIGNALED(proc_ptr->status)) {
                        proc_ptr->termsig = WTERMSIG(proc_ptr->status);
                        break;
                    }
                    --retry_count;
                    printf("failed to signal process %d successfully, sleep(5) then retry\n", proc_ptr->pid);
                    fflush(stdout);
                    sleep(5);
                    // send sigterm again
                    kill(proc_ptr->pid, SIGTERM);
                }
                // FIX up exit value
                if (WIFEXITED(proc_ptr->status))  {
                    proc_ptr->exitv = WEXITSTATUS(proc_ptr->status);
                    if (proc_ptr->exitv == 128) {
                        proc_ptr->exitv = 0;
                    }
                    proc_ptr->finished = 1;
                }
            } else {
                printf("Fail to signal %d with SIGTERM\n", proc_ptr->pid);
                fflush(stdout);
            }
        }
    }
}

int do_sigterm(const char **argv, int count) {
    unsigned sleep_time;
    struct timespec ts;

    if (count != 1 ) {
        puts("ERROR: invalid argument count for sigterm command");
        return 1;
    }
    if (!strtounsigned(*argv, &sleep_time)) {
        printf("ERROR: invalid sleep time %s for sigterm command\n", *argv);
        fflush(stdout);
        return 1;
    }
    e2e_signal_gen* sig_gen_ptr = calloc(sizeof(*sig_gen_ptr), 1);
    if (sig_gen_ptr == NULL) {
        perror("calloc failed");
        return 1;
    }

    clock_gettime(CLOCK_MONOTONIC, &ts);
    sig_gen_ptr->next = sig_gen_list;
    sig_gen_list = sig_gen_ptr;
    sig_gen_ptr->endtick = ts.tv_sec + sleep_time;
    sig_gen_ptr->func = send_sigterm;

    printf("SIGTERM will be sent to all child processes after %d seconds\n", sleep_time);
    fflush(stdout);
    return 0;
}

int do_liveness(const char **argv, int count) {
    int retv = 1;
    unsigned read_interval;
    time_t timeout;
    puts("LIVENESS");
    if (count == 3) {
       if (strtounsigned(argv[1], &read_interval) && strtolong(argv[2], &timeout)) {
            if (add_target(liveness_ctx, argv[0], read_interval, timeout)) {
                retv = 0;
            }
       } else {
            printf("ERROR: invalid arguments for liveness %d, %s %s %s\n", count, argv[0], argv[1], argv[2]);
       }
    } else {
        puts("ERROR: invalid argument count for liveness");
    }
    return retv;
}

int do_session_id(const char **argv, int count) {
    int retv = 1;

    if (count != 1 ) {
        puts("ERROR: invalid argument count for session_id command");
        return 1;
    }

    if (strlcpy(session_id, argv[0], sizeof(session_id)) > sizeof(session_id)) {
        printf("session_id is too large\n");
    } else {
        printf("session_id set to %s\n", session_id);
        retv = 0;
    }
    fflush(stdout);
    return retv;
}


/* Wait for all processes in the list to complete.*/
int wait_procs() {
    int exitv = 0;
    int pending;
    struct timespec ts;
    bool live_ok;

    do {
        sleep(2);
        pending = 0;
        live_ok = liveness_check(liveness_ctx);
        if (!live_ok) {
            puts("*** liveness check failed ***");
            dump_liveness_context(liveness_ctx);
            fflush(stdout);
            kill_procs(SIGINT);
        }
        clock_gettime(CLOCK_MONOTONIC, &ts);
        for (e2e_signal_gen* sig_gen_ptr = sig_gen_list; NULL != sig_gen_ptr; sig_gen_ptr = sig_gen_ptr->next) {
            if (ts.tv_sec >= sig_gen_ptr->endtick && sig_gen_ptr->func != NULL) {
                sig_gen_ptr->func();
                sig_gen_ptr->func = NULL;
            }
        }
        for (e2e_process* proc_ptr = proc_list; NULL != proc_ptr; proc_ptr = proc_ptr->next) {
            if (proc_ptr->finished) {
                continue;
            }

            if ( 0 == waitpid(proc_ptr->pid, &proc_ptr->status, WNOHANG)) {
                pending += 1;
                continue;
            }

            proc_ptr->finished = 1;
            printf("** %d finished **\n", proc_ptr->pid);
            if (WIFEXITED(proc_ptr->status)) {
                proc_ptr->exitv = WEXITSTATUS(proc_ptr->status);
                if (0 != proc_ptr->exitv) {
                    printf("** exit value = %d for %s **\n", proc_ptr->exitv, proc_ptr->cmd);
                }
            } else if (WIFSIGNALED(proc_ptr->status)) {
                proc_ptr->termsig = WTERMSIG(proc_ptr->status);
                printf("** termsig %d, %s **\n", proc_ptr->termsig, proc_ptr->cmd);
            } else {
                /* Should not reach here */
                printf("** Bug in handling waitpid status **\n");
                proc_ptr->abnormal_exit = 1;
            }
            fflush(stdout);
        }
    } while(pending);

    dump_liveness_context(liveness_ctx);
    if (live_ok) {
        for (e2e_process* proc_ptr = proc_list; NULL != proc_ptr; proc_ptr = proc_ptr->next) {
            if (proc_ptr->exitv) {
                exitv = proc_ptr->exitv;
            } else if (proc_ptr->termsig) {
                exitv = 254;
            } else if (proc_ptr->abnormal_exit) {
                exitv = 255;
            }
        }
    } else {
        fflush(stdout);
        exitv = 253;
    }
    return exitv;
}

/* Print contents of processes in the list. */
void print_procs() {
    for (e2e_process* proc_ptr = proc_list; NULL != proc_ptr; proc_ptr = proc_ptr->next) {
        printf("pid:%d, status=%d, exit=%d, termsig=%d, abnormal_exit=%d finished=%d wait=%d\ncmd=%s\n",
               proc_ptr->pid,
               proc_ptr->status,
               proc_ptr->exitv,
               proc_ptr->termsig,
               proc_ptr->abnormal_exit,
               proc_ptr->finished,
               proc_ptr->wait,
               proc_ptr->cmd);
    }
}


int parse_cmds(const char* arg) {
    if ( arg == NULL )
        return KW_EOL;

    for (int ix=0; ix < sizeof(cmd_keywords)/sizeof(cmd_keywords[0]); ++ix) {
        if ( 0 == strcmp(arg, cmd_keywords[ix])) {
            return ix;
        }
    }
    return NOT_A_KEYWORD;
}

int process_args(int token, const char **args, int count, bool wait, int *p_exitv) {
    puts("------------------");
    args[count] = NULL;
    switch (token) {
        case KW_SLEEP:
            return do_sleep(args, count);
        case KW_SEGFAULT:
            return do_segfault(args, count);
        case KW_SIGTERM:
            return do_sigterm(args, count);
        case KW_LIVENESS:
            return do_liveness(args, count);
        case KW_MAKEFILE:
            return do_makefile(args, count);
        case KW_ZEROFILL:
            return do_zerofill(args, count);
        case KW_SESSION_ID:
            return do_session_id(args, count);
        case KW_POST_OP_SLEEP:
            return do_post_op_sleep(args, count);
        case KW_FILESIZE:
            return do_filesize(args, count);
        case KW_EXITV:
            if (count == 1) {
                *p_exitv = atoi(args[0]);
                printf("exit value set to %d\n", *p_exitv);
                return 0;
            }
            printf("invalid argument count\n");
            *p_exitv = -1;
            return 1;
        default:
            return do_exec(args, count, wait);
    }
}

/*
 * Usage:
 * See README.md
 */
int main(int argc, const char **argv)
{
    int token = NOT_A_KEYWORD;
    int idx;
    bool wait = false;
    int exitv = 0;
    int procs_exitv = 0;
    int parse_return_value = 0;
    const char **args = calloc(sizeof(*argv), argc);
    time_t start = time(NULL);

    system("date");
    printf("e2e_fio: version %s\n", VERSION);
    {
        puts("\ncommand line:");
        for(const char **argv_scan=argv; *argv_scan != NULL; ++argv_scan) {
            printf("%s ", *argv_scan);
        }
    }
    puts("\n");
    if (0 != mkdir(workspace_path, 0777)) {
        if (errno != EEXIST) {
            printf("failed to create workspace directory %s\n", workspace_path);
            return -1;
        }
    }
    liveness_ctx = make_liveness_context();
    if (liveness_ctx == NULL) {
        puts("out of memory");
        return -1;
    }
    fflush(stdout);

    if (args == NULL ) {
        puts("calloc failed");
        return -1;
    }

    ++argv;
    idx = 0;
    wait = false;
    do {
        int t = parse_cmds(*argv);
        if (token == NOT_A_KEYWORD) {
            /* not parsing a command,
             * so expect a token and start parsing command or ignore
             */
            switch(t) {
                /* start parsing */
                case KW_FIO_IMPLICIT: /* "--" */
                    idx=0;
                    args[idx] = "fio";
                    ++idx;
                    token = t;
                    wait = false;
                    break;
                case KW_TASK_START:  /* "---" */
                case KW_SLEEP:
                case KW_POST_OP_SLEEP:
                case KW_SEGFAULT:
                case KW_SIGTERM:
                case KW_LIVENESS:
                case KW_MAKEFILE:
                case KW_ZEROFILL:
                case KW_SESSION_ID:
                case KW_FILESIZE:
                case KW_EXITV:
                    token = t;
                    idx = 0;
                    wait = true;
                    break;
                case KW_EOL:
                    break;
                /* ignore */
                default:
                    printf("ignoring %s\n", *argv);
                    break;
            }
        } else {
            /* parsing a command,
             * so expect
             *  - start of next command
             *  - end of current command
             *  - command argument
             */
            switch(t) {
                /* start of next command  - implicit end of this command */
                case KW_TASK_START: /* "---" */
                case KW_FIO_IMPLICIT: /* "--" */
                    --argv;
                /* end of this command */
                case KW_TASK_WAIT: /* "&&" */
                case KW_TASK_END: /* ";" */
                case KW_EOL:
                    if ( t == KW_TASK_WAIT) { wait = true; }
                    if ( t == KW_TASK_END) { wait = false; }
                    parse_return_value = process_args(token, args, idx, wait, &exitv);
                    fflush(stdout);
                    token = NOT_A_KEYWORD;
                    break;
                /* command argument */
                case NOT_A_KEYWORD:
                    args[idx] = *argv;
                    ++idx;
                    break;
                default:
                    printf("Unexpected token value %d\n", t);
                    return 252;
            }
        }
    } while(*argv != NULL && ++argv && parse_return_value == 0);
    free(args);

    start_liveness_checks(liveness_ctx);
    /* wait for forked processes to complete */
    procs_exitv = wait_procs();
    puts("");
    print_procs();
    stop_liveness_checks(liveness_ctx);

    if (0 == exitv) {
        exitv = procs_exitv;
        if (0 == procs_exitv) {
            exitv = parse_return_value;
        }
    }

    if (strlen(session_id) > 0) {
        char cmd[512];
        puts("------------------");
        fflush(stdout);
        sprintf(cmd, "./archive_wksp.sh %s", session_id);
        system(cmd);
    }
//    puts("executing sync");
//    fflush(stdout);
//    system("sync");
    {
        while (files_list != NULL) {
            e2e_fio_files* eff = files_list;
            printf("unlinking %s\n", eff->filename);
            fflush(stdout);
            if (unlink(eff->filename) !=0) {
                perror("unlink  failed");
            }
            eff = eff->next;
            printf("#### 1 %p %p %s\n", files_list, eff, files_list->filename);
            free(files_list->filename);
            free(files_list);
            files_list = eff;
            printf("#### 1 %p %p\n", eff, files_list);
        }
    }
    printf("Exit value is %d\n", exitv);
    printf("\nJSON{\"exit_value\": %d, \"elapsed_seconds\" : %ld}\n", exitv, time(NULL) - start);
    system("date");
    if (post_op_sleep != 0) {
        printf("post op sleep %d seconds\n", post_op_sleep);
        fflush(stdout);
        sleep(post_op_sleep);
    }
    return exitv;
}
