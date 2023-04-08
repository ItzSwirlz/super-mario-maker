package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/PretendoNetwork/super-mario-maker-secure/utility"
)

func FollowingsLatestCourseSearchObject(err error, client *nex.Client, callID uint32, dataStoreSearchParam *datastore.DataStoreSearchParam, extraData []string) {
	pRankingResults := make([]*datastore_super_mario_maker.DataStoreCustomRankingResult, 0)

	for _, pid := range dataStoreSearchParam.OwnerIds {
		courseMetadatas := database.GetCourseMetadatasByPID(pid)

		for _, courseMetadata := range courseMetadatas {
			pRankingResults = append(pRankingResults, utility.CourseMetadataToDataStoreCustomRankingResult(courseMetadata))
		}
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(pRankingResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodFollowingsLatestCourseSearchObject, rmcResponseBody)

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
}
