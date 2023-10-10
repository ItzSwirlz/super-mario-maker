package nex_datastore_super_mario_maker

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func PrepareAttachFile(err error, client *nex.Client, callID uint32, dataStoreAttachFileParam *datastore_super_mario_maker_types.DataStoreAttachFileParam) uint32 {
	key := fmt.Sprintf("image/%d.jpg", dataStoreAttachFileParam.ReferDataID)
	bucket := os.Getenv("S3_BUCKET_NAME")
	date := strconv.Itoa(int(time.Now().Unix()))
	pid := strconv.Itoa(int(client.PID()))

	data := pid + bucket + key + date

	hmac := hmac.New(sha256.New, []byte{})
	hmac.Write([]byte(data))

	signature := hex.EncodeToString(hmac.Sum(nil))

	fieldBucket := datastore_types.NewDataStoreKeyValue()
	fieldBucket.Key = "bucket"
	fieldBucket.Value = bucket

	fieldKey := datastore_types.NewDataStoreKeyValue()
	fieldKey.Key = "key"
	fieldKey.Value = key

	fieldACL := datastore_types.NewDataStoreKeyValue()
	fieldACL.Key = "acl"
	fieldACL.Value = "public-read"

	fieldContentType := datastore_types.NewDataStoreKeyValue()
	fieldContentType.Key = "content-type"
	fieldContentType.Value = "image/jpeg"

	fieldPID := datastore_types.NewDataStoreKeyValue()
	fieldPID.Key = "pid"
	fieldPID.Value = pid

	fieldDate := datastore_types.NewDataStoreKeyValue()
	fieldDate.Key = "date"
	fieldDate.Value = date

	fieldSignature := datastore_types.NewDataStoreKeyValue()
	fieldSignature.Key = "signature"
	fieldSignature.Value = signature

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	pReqPostInfo := datastore_types.NewDataStoreReqPostInfo()

	pReqPostInfo.DataID = dataStoreAttachFileParam.ReferDataID
	pReqPostInfo.URL = os.Getenv("DATASTORE_UPLOAD_URL")
	pReqPostInfo.RequestHeaders = []*datastore_types.DataStoreKeyValue{}
	pReqPostInfo.FormFields = []*datastore_types.DataStoreKeyValue{fieldBucket, fieldKey, fieldACL, fieldContentType, fieldPID, fieldDate, fieldSignature}
	pReqPostInfo.RootCACert = []byte{}

	rmcResponseStream.WriteStructure(pReqPostInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodPrepareAttachFile, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)

	return 0
}
