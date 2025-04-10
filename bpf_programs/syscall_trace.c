#include <uapi/linux/ptrace.h>
#include <linux/sched.h>

struct event_t {
    u32 pid;
    u32 tid;
    u64 timestamp;
    char comm[16];
    u32 syscall_nr;
};

BPF_PERF_OUTPUT(events);
BPF_HASH(start, u32, u64);

int syscall_entry(struct pt_regs *ctx) {
    u64 ts = bpf_ktime_get_ns();
    u32 tid = bpf_get_current_pid_tgid();
    
    start.update(&tid, &ts);
    return 0;
}

int syscall_return(struct pt_regs *ctx) {
    u32 tid = bpf_get_current_pid_tgid();
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    
    u64 *tsp = start.lookup(&tid);
    if (!tsp)
        return 0;
        
    struct event_t event = {};
    event.pid = pid;
    event.tid = tid;
    event.timestamp = bpf_ktime_get_ns() - *tsp;
    event.syscall_nr = ctx->orig_ax;
    
    bpf_get_current_comm(&event.comm, sizeof(event.comm));
    events.perf_submit(ctx, &event, sizeof(event));
    
    start.delete(&tid);
    return 0;
}
