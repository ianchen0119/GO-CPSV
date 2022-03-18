/*      -*- OpenSAF  -*-
 *
 * (C) Copyright 2008 The OpenSAF Foundation
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
 * or FITNESS FOR A PARTICULAR PURPOSE. This file and program are licensed
 * under the GNU Lesser General Public License Version 2.1, February 1999.
 * The complete license can be accessed from the following location:
 * http://opensource.org/licenses/lgpl-license.php
 * See the Copying file included with the OpenSAF distribution for full
 * licensing terms.
 *
 * Original Author(s): Emerson Network Power
 */

/*****************************************************************************
..............................................................................
MODULE NAME: cpsv.c (derived from cpsv_test_app.c)

  .............................................................................
  DESCRIPTION:

  This program interacts with OpenSAF CKPT and GO-CPSV.
******************************************************************************/

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

Status cpsv_ckpt_init(char* newName){
	SaAisErrorT rc;
	memset(&ckptName, 0, 255);
	ckptName.length = strlen(newName);
	memcpy(ckptName.value, newName, strlen(newName));

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
	ckptCreateAttr.creationFlags = SA_CKPT_WR_ALL_REPLICAS;
	ckptCreateAttr.checkpointSize = 1024;
	ckptCreateAttr.retentionDuration = 100000;
	ckptCreateAttr.maxSections = 20;
	ckptCreateAttr.maxSectionSize = 1000;
	ckptCreateAttr.maxSectionIdSize = 20;

	ckptOpenFlags = SA_CKPT_CHECKPOINT_CREATE | SA_CKPT_CHECKPOINT_READ |
			SA_CKPT_CHECKPOINT_WRITE;
	printf("Opening Non-Collocated Checkpoint = %s with create flags....\n",
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

unsigned char* cpsv_sync_read(char* sectionId, SaOffsetT offset, int dataSize, unsigned char isFixed, int* dataSizePtr){
	SaAisErrorT rc;
	SaCkptIOVectorElementT readVector;
	int* sizeBuf;
	unsigned char* read_buff;
	if (isFixed != 1) {
		/* 
		 *  Store non-fixed length data
		 *	for example: json, protobuf...
		 */
		sizeBuf = (int*) malloc(sizeof(int));
		*sizeBuf = 4;
		readVector.sectionId.id = (unsigned char *)sectionId;
		readVector.sectionId.idLen = 2;
		readVector.dataBuffer = (unsigned char *)sizeBuf;
		readVector.dataSize = 4;
		readVector.readSize = 4;
		readVector.dataOffset = 0;
		rc = saCkptCheckpointRead(checkpointHandle, &readVector, 1,
					&erroneousVectorIndex);
		printf("Section Data Size = \"%d\"\n",
				*(int*) readVector.dataBuffer);
		if (rc == SA_AIS_OK) {
			printf("PASSED \n");
		} else {
			printf("Failed \n");
			free(sizeBuf);
			return (void*)0;
		}
	}
	if (isFixed != 1) {
		int nonFixedDataSize = *sizeBuf;
		printf("nonFixedDataSize: %d\n", *sizeBuf);
		*dataSizePtr = nonFixedDataSize;
		read_buff = calloc(nonFixedDataSize, sizeof(unsigned char));
		readVector.dataSize = nonFixedDataSize;
		readVector.readSize = nonFixedDataSize;
		readVector.dataOffset = 4;
		free(sizeBuf);
	} else {
		read_buff = calloc(dataSize, sizeof(unsigned char));
		readVector.dataSize = dataSize;
		readVector.readSize = dataSize;
		readVector.dataOffset = offset;
	}
	
	readVector.sectionId.id = (unsigned char *)sectionId;
	readVector.sectionId.idLen = 2;
	readVector.dataBuffer = read_buff;


	printf("Section-Id = %s ....\n", readVector.sectionId.id);
	rc = saCkptCheckpointRead(checkpointHandle, &readVector, 1,
					&erroneousVectorIndex);
	printf("Checkpoint Data Read = \"%s\"\n",
		    (char*) readVector.dataBuffer);
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

Status cpsv_sync_write(char* sectionId, unsigned char* data, SaOffsetT offset, int dataSize, unsigned char isFixed){
	SaAisErrorT rc;
	SaCkptIOVectorElementT writeVector;

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
	} else if (rc == SA_AIS_ERR_EXIST) {
		printf("Section is already exist.\n");
	} else {
		goto err;
	}

	if (isFixed != 1) {
		int* sizeBuf = (int*) malloc(sizeof(int));
		*sizeBuf = dataSize;
		writeVector.sectionId.id = (unsigned char *)sectionId;
		writeVector.sectionId.idLen = strlen(sectionId);
		writeVector.dataBuffer = (unsigned char *)sizeBuf;
		writeVector.dataSize = 4;
		writeVector.dataOffset = 0;
		writeVector.readSize = 0;
		rc = saCkptCheckpointWrite(checkpointHandle, &writeVector, 1,
						&erroneousVectorIndex);
		if (rc != SA_AIS_OK) {
			goto err;
		}
		free(sizeBuf);
	} 

	if (isFixed != 1) {
		writeVector.dataOffset = 4;
	} else {
		writeVector.dataOffset = offset;
	}
	
	writeVector.sectionId.id = (unsigned char *)sectionId;
	writeVector.sectionId.idLen = strlen(sectionId);
	writeVector.dataBuffer = data;
	writeVector.dataSize = dataSize;
	writeVector.readSize = 0;

	printf("Writing to Checkpoint %s ....\n", ckptName.value);
	printf("Section-Id = %s ....\n", writeVector.sectionId.id);
	printf("CheckpointData being written = \"%s\"\n",
		    (char*) writeVector.dataBuffer);
	printf("DataOffset = %llu ....\n", writeVector.dataOffset);
	rc = saCkptCheckpointWrite(checkpointHandle, &writeVector, 1,
					&erroneousVectorIndex);
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
