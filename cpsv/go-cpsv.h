#ifndef __GO_CPSV_H__
#define __GO_CPSV_H__
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <saCkpt.h>
#include <time.h>

#define DEMO_CKPT_NAME "safCkpt=DemoCkpt,safApp=safCkptService"
#define Status int

void AppCkptOpenCallback(SaInvocationT invocation,
			 SaCkptCheckpointHandleT checkpointHandle,
			 SaAisErrorT error);
void AppCkptSyncCallback(SaInvocationT invocation, SaAisErrorT error);
Status cpsv_ckpt_init();
Status cpsv_ckpt_destroy();
Status cpsv_sync_read(char* sectionId, unsigned char* buffer, SaOffsetT offset, int dataSize);
Status cpsv_sync_write(char* sectionId, char* data, SaOffsetT offset, int dataSize);

#endif