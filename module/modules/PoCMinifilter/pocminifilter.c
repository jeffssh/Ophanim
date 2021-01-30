//#include <ntddk.h>
//#include <wdf.h>
#include "pocminifilter.h"


DRIVER_INITIALIZE DriverEntry;

NTSTATUS
PoCUnload(
    _Unreferenced_parameter_ FLT_FILTER_UNLOAD_FLAGS Flags
);

FLT_PREOP_CALLBACK_STATUS
PreFileCreate(
    _Inout_ PFLT_CALLBACK_DATA Data,
    _In_ PCFLT_RELATED_OBJECTS FltObjects,
    _Flt_CompletionContext_Outptr_ PVOID* CompletionContext
);

// https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/fltkernel/ns-fltkernel-_flt_parameters
CONST FLT_OPERATION_REGISTRATION Callbacks[] = {
    { IRP_MJ_CREATE,
      0,
      PreFileCreate,
      NULL},
      //PostFileCreate },

    { IRP_MJ_OPERATION_END }
};

CONST FLT_REGISTRATION FilterRegistration = {

    sizeof(FLT_REGISTRATION),         //  Size
    FLT_REGISTRATION_VERSION,           //  Version
    0,                                  //  Flags - TODO FLTFL_REGISTRATION_SUPPORT_NPFS_MSFS

    NULL,                //  Context - could be null?
    Callbacks,                          //  Operation callbacks

    PoCUnload,                           //  MiniFilterUnload

    NULL,                     //  InstanceSetup - TODO... I think we want attach to all here.
    NULL,                     // InstanceQueryTeardownCallback - https://docs.microsoft.com/en-us/windows-hardware/drivers/ifs/loading-and-unloading
                              // I don't think this matters for our purpose, since we're just reporting on file opens. 
                              // All pre/post ops continue execution, and we're not opening files in those so unmounting a volume shouldn't matter.
    NULL,                     //  InstanceTeardownStart
    NULL,                     //  InstanceTeardownComplete

    NULL,                               //  GenerateFileName
    NULL,                               //  NormalizeNameComponentCallback
    NULL,                               //  NormalizeContextCleanupCallback
    NULL,          //  TransactionNotificationCallback
    NULL,                               //  NormalizeNameComponentExCallback
    NULL            //  SectionNotificationCallback
};

NTSTATUS
DriverEntry(
    _In_ PDRIVER_OBJECT     DriverObject,
    _In_ PUNICODE_STRING    RegistryPath
)
{

    //UNREFERENCED_PARAMETER(DriverObject);
    //UNREFERENCED_PARAMETER(RegistryPath);

    // NTSTATUS variable to record success or failure
    NTSTATUS status = STATUS_SUCCESS;
    KdPrintEx((DPFLTR_IHVDRIVER_ID, DPFLTR_INFO_LEVEL, "[PoC]: DriverEntry: Entered\n"));
    KdPrintEx((DPFLTR_IHVDRIVER_ID, DPFLTR_INFO_LEVEL, "[PoC]: DriverEntry: RegistryPath: %wZ\n", RegistryPath));
    

    RtlZeroMemory(&Globals, sizeof(Globals));
    ExInitializeResourceLite(&Globals.PoCGlobalLock);

    try {

        status = FltRegisterFilter(DriverObject, &FilterRegistration, &Globals.Filter);
        if (!NT_SUCCESS(status)) {
            KdPrintEx((DPFLTR_IHVDRIVER_ID, DPFLTR_INFO_LEVEL, "FltRegisterFilter FAILED. status = 0x%x\n", status));
            leave;
        }

        // todo set up grpc stuff

        status = FltStartFiltering(Globals.Filter);
        if (!NT_SUCCESS(status)) {
            KdPrintEx((DPFLTR_IHVDRIVER_ID, DPFLTR_INFO_LEVEL, "FltStartFiltering FAILED. status = 0x%x\n", status));
            leave;
        }
    }
    finally {
        if (!NT_SUCCESS(status)) {

            if (NULL != Globals.Filter) {

                FltUnregisterFilter(Globals.Filter);
                Globals.Filter = NULL;
            }

            ExDeleteResourceLite(&Globals.PoCGlobalLock);
        }
    }


    KdPrintEx((DPFLTR_IHVDRIVER_ID, DPFLTR_INFO_LEVEL, "[PoC]: DriverEntry: Loaded!\n"));
    return status;

}


NTSTATUS
PoCUnload(
    _Unreferenced_parameter_ FLT_FILTER_UNLOAD_FLAGS Flags
)
{
    PAGED_CODE();

    UNREFERENCED_PARAMETER(Flags);

    //
    //  Traverse the scan context list, and cancel the scan if it exists.
    //

    PoCAcquireResourceExclusive(&Globals.PoCGlobalLock);
    Globals.Unloading = TRUE;
    PoCReleaseResource(&Globals.PoCGlobalLock);

    // TODO close grpc connection
    /*
    FltCloseCommunicationPort( Globals.ScanServerPort );
    Globals.ScanServerPort = NULL;
    FltCloseCommunicationPort( Globals.AbortServerPort );
    Globals.AbortServerPort = NULL;
    FltCloseCommunicationPort( Globals.QueryServerPort );
    Globals.QueryServerPort = NULL;

    */

    FltUnregisterFilter(Globals.Filter);  // This will typically trigger instance tear down.
    Globals.Filter = NULL;

    ExDeleteResourceLite(&Globals.PoCGlobalLock);

    KdPrintEx((DPFLTR_IHVDRIVER_ID, DPFLTR_INFO_LEVEL, "[PoC]: PoCUnload: Completed!\n"));

    return STATUS_SUCCESS;
}


FLT_PREOP_CALLBACK_STATUS
PreFileCreate(
    _Inout_ PFLT_CALLBACK_DATA Data,
    _In_ PCFLT_RELATED_OBJECTS FltObjects,
    _Flt_CompletionContext_Outptr_ PVOID* CompletionContext
)
{
    UNREFERENCED_PARAMETER(FltObjects);
    UNREFERENCED_PARAMETER(Data);
    UNREFERENCED_PARAMETER(CompletionContext);

    //char* message; todo

    PAGED_CODE();
    //FltObjects->FileObject->FileName.Buffer
    KdPrintEx((DPFLTR_IHVDRIVER_ID, DPFLTR_INFO_LEVEL, "[PoC]: PreFileCreate: %wZ\n", FltObjects->FileObject->FileName));
    /*
    if (FlagOn(Data->Iopb->Parameters.Create.Options, FILE_DIRECTORY_FILE)) {

    }
    */
    return FLT_PREOP_SUCCESS_NO_CALLBACK;
}
