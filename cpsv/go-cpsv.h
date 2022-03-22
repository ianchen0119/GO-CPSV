#ifndef __GO_CPSV_H__
#define __GO_CPSV_H__
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <saCkpt.h>
#include <time.h>
#include <limits.h>

#define SECTION_SIZE 20000
#define SECTIONS 100000

#define Status int

void AppCkptOpenCallback(SaInvocationT invocation,
			 SaCkptCheckpointHandleT checkpointHandle,
			 SaAisErrorT error);
void AppCkptSyncCallback(SaInvocationT invocation, SaAisErrorT error);
Status cpsv_ckpt_init(char* newName);
Status cpsv_ckpt_destroy();
unsigned char* cpsv_sync_read(char* sectionId, SaOffsetT offset, int dataSize, unsigned char isFixed, int* dataSizePtr);
Status cpsv_sync_write(char* sectionId, unsigned char* data, SaOffsetT offset, int dataSize, unsigned char isFixed);

#endif