/*++

Copyright (c) 2011  Microsoft Corporation

Module Name:

    fsmflib.h

Abstract:

    This header file defines the common data structure used by kernel and user.

Environment:

    User mode
    Kernel mode

--*/

#ifndef __FSMFLIB_H__
#define __FSMFLIB_H__

#if defined(_MSC_VER)
#if (_MSC_VER >= 1200)
#pragma warning(push)
#pragma warning(disable:4201) // nonstandard extension used : nameless struct/union
#endif
#endif

//
//  Name of AV filter server ports
//

#define AV_SCAN_PORT_NAME                    L"\\MicrosoftAvSampleFilterScanPort"
#define AV_ABORT_PORT_NAME                   L"\\MicrosoftAvSampleFilterAbortPort"
#define AV_QUERY_PORT_NAME                   L"\\MicrosoftAvSampleFilterQueryPort"


//
//  Definition of invalide section handle for data scan
//

#define AV_INVALID_SECTION_HANDLE   ((HANDLE)((LONG_PTR)(-1)))


//
//  Command type enumeration, please see COMMAND_MESSAGE below
//

typedef enum _FileSystemMinifilter_COMMAND {

    AvIsFileModified,
    AvCmdCreateSectionForDataScan,
    AvCmdCloseSectionForDataScan

} FileSystemMinifilter_COMMAND;

//
//  Message type enumeration, please see AV_SCANNER_NOTIFICATION below
//

typedef enum _FileSystemMinifilter_MESSAGE {

    AvMsgStartScanning,
    AvMsgAbortScanning,
    AvMsgFilterUnloading

} FileSystemMinifilter_MESSAGE;

typedef enum _FileSystemMinifilter_REASON {
    FileSystemMinifilterOnOpen,
    FileSystemMinifilterOnCleanup

} FileSystemMinifilter_REASON;

typedef enum _FileSystemMinifilter_RESULT {

    FileSystemMinifilterResultUndetermined,
    FileSystemMinifilterResultInfected,
    FileSystemMinifilterResultClean

} FileSystemMinifilter_RESULT;

//
//  Defines the commands between the user program and the filter
//  Command: User -> Kernel
//

typedef struct _COMMAND_MESSAGE {

    //
    //  Command type
    //

    FileSystemMinifilter_COMMAND      Command;
    
    //
    //  Scan identifier.
    //  This argument will be checked in message notificaiton callback.
    //
    
    LONGLONG  ScanId;
    
    //
    //  Scan thread id. This id will be used in cancel message passing.
    //  So that we will know which scan thread to cancel.
    //
    
    ULONG  ScanThreadId;
    
    union {
    
        //
        //  When user program is connecting for query (AvConnectForQuery)
        //  it has to pass the file handle to query the status of the file.
        //  Valid when Command == AvIsFileModified
        //
        
        HANDLE FileHandle;
        
        //
        //  The result result.
        //  Valid when Command == AvCmdCloseSectionForDataScan
        //
        FileSystemMinifilter_RESULT ScanResult;
    };
    
} COMMAND_MESSAGE, *PCOMMAND_MESSAGE;

//
//  Message: Kernel -> User Message
//

typedef struct _SCANNER_NOTIFICATION {

    //
    //  Message type
    //
    
    FileSystemMinifilter_MESSAGE Message;

    //
    //  Reason
    //

    FileSystemMinifilter_REASON  Reason;
    
    //
    //  Scan identifier.
    //  This argument will be checked in message notificaiton callback.
    //
    
    LONGLONG  ScanId;
    
    //
    //  Scan thread id. This id will be used in cancel message passing.
    //  So that we will know which scan thread to cancel.
    //
    
    ULONG  ScanThreadId;
    
} AV_SCANNER_NOTIFICATION, *PAV_SCANNER_NOTIFICATION;

//
//  Connection type enumeration. It would be mainly used in connection context.
//

typedef enum _FileSystemMinifilter_CONNECTION_TYPE {

    AvConnectForScan = 1,
    AvConnectForAbort,
    AvConnectForQuery

} FileSystemMinifilter_CONNECTION_TYPE, *PFileSystemMinifilter_CONNECTION_TYPE;

//
//  Connection context. It will be passed through FilterConnectCommunicationPort(...)
//

typedef struct _AV_CONNECTION_CONTEXT {

    FileSystemMinifilter_CONNECTION_TYPE   Type;

} AV_CONNECTION_CONTEXT, *PAV_CONNECTION_CONTEXT;

//
//  The following string is actully "message to be found"
//

#define AV_DEFAULT_SEARCH_PATTERN       "7?));=?z.5z8?z<5/4>"
#define AV_DEFAULT_SEARCH_PATTERN_SIZE  sizeof(AV_DEFAULT_SEARCH_PATTERN)
#define AV_DEFAULT_PATTERN_XOR_KEY      90

#if defined(_MSC_VER)
#if (_MSC_VER >= 1200)
#pragma warning(pop)
#endif
#endif

#endif

