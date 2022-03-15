#include "go-cpsv.h"

SaCkptHandleT ckptHandle;
SaCkptCheckpointHandleT checkpointHandle;
SaCkptCallbacksT callbk;
SaVersionT version;
SaNameT ckptName;
SaCkptCheckpointCreationAttributesT ckptCreateAttr;
SaCkptCheckpointOpenFlagsT ckptOpenFlags;
SaCkptSectionCreationAttributesT sectionCreationAttributes;
SaUint32T erroneousVectorIndex;
const void *initialData = "Default data in the section";
SaTimeT timeout = 1000000000;

Status cpsv_ckpt_init(){
	SaAisErrorT rc;
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
	// ckptCreateAttr.creationFlags =
	//     SA_CKPT_CHECKPOINT_COLLOCATED | SA_CKPT_WR_ACTIVE_REPLICA;
	ckptCreateAttr.creationFlags = SA_CKPT_WR_ALL_REPLICAS;
	ckptCreateAttr.checkpointSize = 1024;
	ckptCreateAttr.retentionDuration = 100000;
	ckptCreateAttr.maxSections = 20;
	ckptCreateAttr.maxSectionSize = 1000;
	ckptCreateAttr.maxSectionIdSize = 20;

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
	SaAisErrorT rc;
	printf("Unlink My Checkpoint ....\t");
	rc = saCkptCheckpointUnlink(ckptHandle, &ckptName);
	if (rc == SA_AIS_OK){
		printf("PASSED \n");
	} else {
		printf("Failed \n");
	}		
	printf("Ckpt Closed ....\t");
	rc = saCkptCheckpointClose(checkpointHandle);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
	} else {
		printf("Failed \n");
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

unsigned char* cpsv_sync_read(char* sectionId, SaOffsetT offset, int dataSize){
	SaAisErrorT rc;
	SaCkptIOVectorElementT readVector;
	unsigned char* read_buff = calloc(dataSize, sizeof(unsigned char));
	readVector.sectionId.id = (unsigned char *)sectionId;
	readVector.sectionId.idLen = 2;
	readVector.dataBuffer = read_buff;
	readVector.dataSize = dataSize;
	readVector.dataOffset = offset;

	printf("Section-Id = %s ....\n", readVector.sectionId.id);
	rc = saCkptCheckpointRead(checkpointHandle, &readVector, 1,
					&erroneousVectorIndex);
	printf("Checkpoint Data Read = \"%d\"\n",
		    *(int*) readVector.dataBuffer);
	// if (rc == SA_AIS_OK) {
	// 	printf("PASSED \n");
	// } else {
	// 	printf("Failed \n");
	// 	goto err;
	// }
	// printf("Synchronizing My Checkpoint being called ....\n");
	// rc = saCkptCheckpointSynchronize(checkpointHandle, timeout);
	if (rc == SA_AIS_OK) {
		printf("PASSED \n");
		return read_buff;
	} else {
		goto err;
	}
	err:
		printf("Failed \n");
		free(read_buff);
		return (void*)0;
}

Status cpsv_sync_write(char* sectionId, unsigned char* data, SaOffsetT offset, int dataSize){
	SaAisErrorT rc;
	SaCkptIOVectorElementT writeVector;
	// printf("Setting the Active Replica for my checkpoint ....\t");
	// rc = saCkptActiveReplicaSet(checkpointHandle);
	// if (rc == SA_AIS_OK) {
	// 	printf("PASSED \n");
	// } else {
	// 	printf("Failed \n");
	// 	return -1;
	// }

	sectionCreationAttributes.sectionId =
		(SaCkptSectionIdT *)malloc(sizeof(SaCkptSectionIdT));
	sectionCreationAttributes.sectionId->id = (unsigned char *)sectionId;
	sectionCreationAttributes.sectionId->idLen = strlen(sectionId);
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
		goto err;
	}

	writeVector.sectionId.id = (unsigned char *)sectionId;
	writeVector.sectionId.idLen = strlen(sectionId);
	writeVector.dataBuffer = data;
	writeVector.dataSize = dataSize;
	writeVector.dataOffset = offset;
	writeVector.readSize = 0;

	printf("Writing to Checkpoint %s ....\n", DEMO_CKPT_NAME);
	printf("Section-Id = %s ....\n", writeVector.sectionId.id);
	printf("CheckpointData being written = \"%d\"\n",
		    *(int*) writeVector.dataBuffer);
	printf("DataOffset = %llu ....\n", writeVector.dataOffset);
	rc = saCkptCheckpointWrite(checkpointHandle, &writeVector, 1,
					&erroneousVectorIndex);
	// if (rc == SA_AIS_OK) {
	// 	printf("PASSED \n");
	// } else {
	// 	goto err;
	// }
	// printf("Synchronizing My Checkpoint being called ....\n");
	// rc = saCkptCheckpointSynchronize(checkpointHandle, timeout);
	if (rc == SA_AIS_OK) {
		free(sectionCreationAttributes.sectionId);
		printf("PASSED \n");
		return 0;
	} else {
		goto err;
	}
err:
	free(sectionCreationAttributes.sectionId);
	printf("Failed \n");
	return -1;
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