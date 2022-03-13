#include "go-cpsv.h"

SaCkptHandleT ckptHandle;
SaCkptCheckpointHandleT checkpointHandle;
SaCkptCallbacksT callbk;
SaVersionT version;
SaNameT ckptName;
SaAisErrorT rc;
SaCkptCheckpointCreationAttributesT ckptCreateAttr;
SaCkptCheckpointOpenFlagsT ckptOpenFlags;
SaCkptSectionCreationAttributesT sectionCreationAttributes;
SaUint32T erroneousVectorIndex;
const void *initialData = "Default data in the section";
SaTimeT timeout = 1000000000;

Status cpsv_ckpt_init(){
	memset(&ckptName, 0, 255);
	ckptName.length = strlen(DEMO_CKPT_NAME);
	memcpy(ckptName.value, DEMO_CKPT_NAME, strlen(DEMO_CKPT_NAME));

	callbk.saCkptCheckpointOpenCallback = AppCkptOpenCallback;
	callbk.saCkptCheckpointSynchronizeCallback = AppCkptSyncCallback;
	version.releaseCode = 'B';
	version.majorVersion = 2;
	version.minorVersion = 2;
	printf("Initialising With Checkpoint Service....\n");
	rc = saCkptInitialize(&ckptHandle, &callbk, &version);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
	} else {
		printf("Failed \n");
		return -1;
	}
	ckptCreateAttr.creationFlags =
	    SA_CKPT_CHECKPOINT_COLLOCATED | SA_CKPT_WR_ACTIVE_REPLICA;
	ckptCreateAttr.checkpointSize = 1024;
	ckptCreateAttr.retentionDuration = 100000;
	ckptCreateAttr.maxSections = 2;
	ckptCreateAttr.maxSectionSize = 700;
	ckptCreateAttr.maxSectionIdSize = 4;

	ckptOpenFlags = SA_CKPT_CHECKPOINT_CREATE | SA_CKPT_CHECKPOINT_READ |
			SA_CKPT_CHECKPOINT_WRITE;
	printf("Opening Collocated Checkpoint = %s with create flags....\n",
	       ckptName.value);
	rc = saCkptCheckpointOpen(ckptHandle, &ckptName, &ckptCreateAttr,
				  ckptOpenFlags, timeout, &checkpointHandle);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
		return 0;
	} else {
		printf("Failed \n");
		return -1;
	}
}

Status cpsv_ckpt_destroy(){
	printf("Ckpt Closed ....\t");
	rc = saCkptCheckpointClose(checkpointHandle);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
	} else {
		printf("Failed \n");
		return -1;
	}

	printf("Ckpt Finalize being called ....\t");
	rc = saCkptFinalize(ckptHandle);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
		return 0;
	} else {
		printf("Failed \n");
		return -1;
	}
}

Status cpsv_sync_read(unsigned char* buffer, SaOffsetT offset, int dataSize){
	SaCkptIOVectorElementT readVector;
	readVector.sectionId.id = (unsigned char *)"11";
	readVector.sectionId.idLen = 2;
	readVector.dataBuffer = buffer;
	readVector.dataSize = dataSize;
	readVector.dataOffset = offset;

	rc = saCkptCheckpointRead(checkpointHandle, &readVector, 1,
					&erroneousVectorIndex);
	printf("Checkpoint Data Read = \"%s\"\n",
		    (char *)readVector.dataBuffer);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
	} else {
		printf("Failed \n");
		return -1;
	}
	printf("Synchronizing My Checkpoint being called ....\n");
	rc = saCkptCheckpointSynchronize(checkpointHandle, timeout);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
		return 0;
	} else {
		printf("Failed \n");
		return -1;
	}
}

Status cpsv_sync_write(char* data, SaOffsetT offset, int dataSize){
	SaCkptIOVectorElementT writeVector;
	printf("Setting the Active Replica for my checkpoint ....\t");
	rc = saCkptActiveReplicaSet(checkpointHandle);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
	} else {
		printf("Failed \n");
		return -1;
	}

	sectionCreationAttributes.sectionId =
		(SaCkptSectionIdT *)malloc(sizeof(SaCkptSectionIdT));
	sectionCreationAttributes.sectionId->id = (unsigned char *)"11";
	sectionCreationAttributes.sectionId->idLen = 2;
	/* 
	 * Cpsv expects `expirationTime` as  absolute time
	 * check  section 3.4.3.2 SaCkptSectionCreationAttributesT of
	 * CKPT Specification for more details  
	 */
	sectionCreationAttributes.expirationTime =
		(SA_TIME_ONE_HOUR +
		(time((time_t *)0) * 1000000000)); /* One Hour */

	printf("Created Section ....\t");
	rc = saCkptSectionCreate(checkpointHandle,
				&sectionCreationAttributes,
				initialData, 28);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
	} else {
		printf("Failed \n");
		return -1;
	}

	writeVector.sectionId.id = (unsigned char *)"11";
	writeVector.sectionId.idLen = 2;
	writeVector.dataBuffer = data;
	writeVector.dataSize = dataSize;
	writeVector.dataOffset = offset;
	writeVector.readSize = 0;

	printf("Writing to Checkpoint %s ....\n", DEMO_CKPT_NAME);
	printf("Section-Id = %s ....\n", writeVector.sectionId.id);
	printf("CheckpointData being written = \"%s\"\n",
		    (char *)writeVector.dataBuffer);
	printf("DataOffset = %llu ....\n", writeVector.dataOffset);
	rc = saCkptCheckpointWrite(checkpointHandle, &writeVector, 1,
					&erroneousVectorIndex);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
	} else {
		printf("Failed \n");
		return -1;
	}
	printf("Synchronizing My Checkpoint being called ....\n");
	rc = saCkptCheckpointSynchronize(checkpointHandle, timeout);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
		return 0;
	} else {
		printf("Failed \n");
		return -1;
	}
}

void AppCkptOpenCallback(SaInvocationT invocation,
			 SaCkptCheckpointHandleT checkpointHandle,
			 SaAisErrorT error)
{
	if (error != SA_AIS_OK) {
		printf("Checkpoint Open Async callback unsuccessful\n");
		return;
	} else {
		printf(
		    "Checkpoint Open Async callback success and ckpt_hdl %llu \n",
		    checkpointHandle);
		return;
	}
}
void AppCkptSyncCallback(SaInvocationT invocation, SaAisErrorT error)
{
	if (error != SA_AIS_OK) {
		printf("Checkpoint Sync Callback unsuccessful\n");
		return;
	} else {
		printf("Checkpoint Sync Callback success\n");
		return;
	}
}