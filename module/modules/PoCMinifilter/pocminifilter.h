/*++

Module Name:

    fsmdriver.h

Abstract:

    Header file which contains the structures, type definitions,
    constants, global variables and function prototypes that are
    only visible within the kernel.

Environment:

    Kernel mode

--*/
#ifndef __POCDRIVER_H__
#define __POCDRIVER_H__

#include <fltkernel.h>


#pragma prefast(disable:__WARNING_ENCODE_MEMBER_FUNCTION_POINTER, "Not valid for kernel mode drivers")

typedef struct _POC_GLOBAL_DATA {
    //
    //  The global FLT_FILTER pointer. Many API needs this, such as 
    //  FltAllocateContext(...)
    //

    PFLT_FILTER Filter;

    //
    //  Server-side communicate ports.
    //

    //PFLT_PORT ScanServerPort;
    //PFLT_PORT AbortServerPort;
    //PFLT_PORT QueryServerPort;

    //
    //  The scan client ports.
    //  These ports are assigned at AvConnectNotifyCallback and cleaned at AvDisconnectNotifyCallback
    //
    //  ScanClientPort is the connection port regarding the scan message.
    //  AbortClientPort is the connection port regarding the abort message.
    //  QueryClient is the connection port regarding the query command.
    //

    //PFLT_PORT ScanClientPort;
    //PFLT_PORT AbortClientPort;
    //PFLT_PORT QueryClientPort;

    //
    //  Scan context list head. 
    //  At AvMessageNotifyCallback, when user passes ScanCtxId, we 
    //  have to check the validity of the id by checking this list.
    //

    //LIST_ENTRY ScanCtxListHead;

    //
    //  The lock that synchronizes the accesses of the scan context list above.
    //

    ERESOURCE PoCGlobalLock;

    //
    //  A flag that indicating that the filter is being unloaded.
    //    

    BOOLEAN  Unloading;

} POC_GLOBAL_DATA, * PPOC_GLOBAL_DATA;

POC_GLOBAL_DATA Globals;


FORCEINLINE
VOID
_Acquires_lock_(_Global_critical_region_)
PoCAcquireResourceExclusive(
    _Inout_ _Acquires_exclusive_lock_(*Resource) PERESOURCE Resource
)
{
    FLT_ASSERT(KeGetCurrentIrql() <= APC_LEVEL);
    FLT_ASSERT(ExIsResourceAcquiredExclusiveLite(Resource) ||
        !ExIsResourceAcquiredSharedLite(Resource));

    KeEnterCriticalRegion();
    (VOID)ExAcquireResourceExclusiveLite(Resource, TRUE);
}

FORCEINLINE
VOID
_Acquires_lock_(_Global_critical_region_)
PoCAcquireResourceShared(
    _Inout_ _Acquires_shared_lock_(*Resource) PERESOURCE Resource
)
{
    FLT_ASSERT(KeGetCurrentIrql() <= APC_LEVEL);

    KeEnterCriticalRegion();
    (VOID)ExAcquireResourceSharedLite(Resource, TRUE);
}

FORCEINLINE
VOID
_Releases_lock_(_Global_critical_region_)
_Requires_lock_held_(_Global_critical_region_)
PoCReleaseResource(
    _Inout_ _Requires_lock_held_(*Resource) _Releases_lock_(*Resource) PERESOURCE Resource
)
{
    FLT_ASSERT(KeGetCurrentIrql() <= APC_LEVEL);
    FLT_ASSERT(ExIsResourceAcquiredExclusiveLite(Resource) ||
        ExIsResourceAcquiredSharedLite(Resource));

    ExReleaseResourceLite(Resource);
    KeLeaveCriticalRegion();
}
#endif

