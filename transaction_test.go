package ergo

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTxId(t *testing.T) {
	txIdStr := "93d344aa527e18e5a221db060ea1a868f46b61e4537e6e5f69ecc40334c15e38"

	testTxId, _ := NewTxId(txIdStr)
	resTxIdStr, _ := testTxId.String()

	assert.Equal(t, txIdStr, resTxIdStr)
}

func TestTxBuilder_Build(t *testing.T) {
	recipient, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")
	boxJson := `{
          "boxId": "e56847ed19b3dc6b72828fcfb992fdf7310828cf291221269b7ffc72fd66706e",
          "value": 67500000000,
          "ergoTree": "100204a00b08cd021dde34603426402615658f1d970cfa7c7bd92ac81a8b16eeebff264d59ce4604ea02d192a39a8cc7a70173007301",
          "assets": [],
          "creationHeight": 284761,
          "additionalRegisters": {},
          "transactionId": "9148408c04c2e38a6402a7950d6157730fa7d49e9ab3b9cadec481d7769918e9",
          "index": 1
        }`
	unspentBox, _ := NewBoxFromJson(boxJson)
	unspentBoxes := NewBoxes()
	unspentBoxes.Add(unspentBox)

	testContract, _ := NewContractPayToAddress(recipient)

	outBoxValue := SafeUserMinBoxValue()
	outbox, _ := NewBoxCandidateBuilder(outBoxValue, testContract, 0).Build()

	txOutputs := NewBoxCandidates()
	txOutputs.Add(outbox)

	fee := SuggestedTxFee()
	changeAddr, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")

	testDataInputs := NewDataInputs()
	testBoxSelector := NewSimpleBoxSelector()

	targetBalance, _ := SumOfBoxValues(outBoxValue, fee)

	testBoxSelection, _ := testBoxSelector.Select(unspentBoxes, targetBalance, NewTokens())

	testTxBuilder := NewTxBuilder(testBoxSelection, txOutputs, 0, fee, changeAddr)
	testTxBuilder.SetDataInputs(testDataInputs)

	tx, txErr := testTxBuilder.Build()
	assert.NoError(t, txErr)
	txJson, txJsonErr := tx.JsonEIP12()
	assert.NoError(t, txJsonErr)
	assert.Equal(t, `{"inputs":[{"boxId":"e56847ed19b3dc6b72828fcfb992fdf7310828cf291221269b7ffc72fd66706e","extension":{}}],"data_inputs":[],"outputs":[{"value":"1000000","ergoTree":"0008cd02229ac0a22560d7bdfa4eb1de64e688390e85339c08aaf018b22d5ce93593192f","assets":[],"additionalRegisters":{},"creationHeight":0},{"value":"67497900000","ergoTree":"0008cd02229ac0a22560d7bdfa4eb1de64e688390e85339c08aaf018b22d5ce93593192f","assets":[],"additionalRegisters":{},"creationHeight":0},{"value":"1100000","ergoTree":"1005040004000e36100204a00b08cd0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798ea02d192a39a8cc7a701730073011001020402d19683030193a38cc7b2a57300000193c2b2a57301007473027303830108cdeeac93b1a57304","assets":[],"additionalRegisters":{},"creationHeight":0}]}`, txJson)
}

func testBlockHeadersFromJson() BlockHeaders {
	blockHeaderJsn := `       {
        "extensionId": "d16f25b14457186df4c5f6355579cc769261ce1aebc8209949ca6feadbac5a3f",
        "difficulty": "626412390187008",
        "votes": "040000",
        "timestamp": 1618929697400,
        "size": 221,
        "stateRoot": "8ad868627ea4f7de6e2a2fe3f98fafe57f914e0f2ef3331c006def36c697f92713",
        "height": 471746,
        "nBits": 117586360,
        "version": 2,
        "id": "4caa17e62fe66ba7bd69597afdc996ae35b1ff12e0ba90c22ff288a4de10e91b",
        "adProofsRoot": "d882aaf42e0a95eb95fcce5c3705adf758e591532f733efe790ac3c404730c39",
        "transactionsRoot": "63eaa9aff76a1de3d71c81e4b2d92e8d97ae572a8e9ab9e66599ed0912dd2f8b",
        "extensionHash": "3f91f3c680beb26615fdec251aee3f81aaf5a02740806c167c0f3c929471df44",
        "powSolutions": {
          "pk": "02b3a06d6eaa8671431ba1db4dd427a77f75a5c2acbd71bfb725d38adc2b55f669",
          "w": "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798",
          "n": "5939ecfee6b0d7f4",
          "d": 0
        },
        "adProofsId": "86eaa41f328bee598e33e52c9e515952ad3b7874102f762847f17318a776a7ae",
        "transactionsId": "ac80245714f25aa2fafe5494ad02a26d46e7955b8f5709f3659f1b9440797b3e",
        "parentId": "6481752bace5fa5acba5d5ef7124d48826664742d46c974c98a2d60ace229a34"
        }`
	testBlockHeaders := NewBlockHeaders()
	for i := 0; i < 10; i++ {
		testBlockHeader, _ := NewBlockHeader(blockHeaderJsn)
		testBlockHeaders.Add(testBlockHeader)
	}
	return testBlockHeaders
}

func TestWallet_SignTransaction(t *testing.T) {
	sk := NewSecretKey()
	inputContract, _ := NewContractPayToAddress(sk.Address())
	testTxIdStr := "93d344aa527e18e5a221db060ea1a868f46b61e4537e6e5f69ecc40334c15e38"
	testTxId, _ := NewTxId(testTxIdStr)
	inputBoxVal, _ := NewBoxValue(1000000000)
	inputBox, _ := NewBox(inputBoxVal, 0, inputContract, testTxId, 0, NewTokens())

	recipient, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")
	unspentBoxes := NewBoxes()
	unspentBoxes.Add(inputBox)
	testContract, _ := NewContractPayToAddress(recipient)
	outBoxValue := SafeUserMinBoxValue()
	outbox, _ := NewBoxCandidateBuilder(outBoxValue, testContract, 0).Build()
	txOutputs := NewBoxCandidates()
	txOutputs.Add(outbox)
	fee := SuggestedTxFee()
	changeAddr, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")
	testDataInputs := NewDataInputs()
	testBoxSelector := NewSimpleBoxSelector()

	targetBalance, _ := SumOfBoxValues(outBoxValue, fee)
	testBoxSelection, _ := testBoxSelector.Select(unspentBoxes, targetBalance, NewTokens())

	testTxBuilder := NewTxBuilder(testBoxSelection, txOutputs, 0, fee, changeAddr)
	testTxBuilder.SetDataInputs(testDataInputs)

	tx, txErr := testTxBuilder.Build()
	assert.NoError(t, txErr)
	_, txJsonErr := tx.JsonEIP12()
	assert.NoError(t, txJsonErr)

	txDataInputs := NewBoxes()

	testBlockHeaders := testBlockHeadersFromJson()
	testBlockHeader, _ := testBlockHeaders.Get(0)
	testPreHeader := NewPreHeader(testBlockHeader)

	ctx, _ := NewStateContext(testPreHeader, testBlockHeaders)
	testSecretKeys := NewSecretKeys()
	testSecretKeys.Add(sk)
	testWallet := NewWalletFromSecretKeys(testSecretKeys)

	signedTx, signingErr := testWallet.SignTransaction(ctx, tx, unspentBoxes, txDataInputs)
	assert.NoError(t, signingErr)
	_, signedTxJsonErr := signedTx.JsonEIP12()
	assert.NoError(t, signedTxJsonErr)
}

func TestMintToken(t *testing.T) {
	recipient, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")
	boxJson := `{
          "boxId": "e56847ed19b3dc6b72828fcfb992fdf7310828cf291221269b7ffc72fd66706e",
          "value": 67500000000,
          "ergoTree": "100204a00b08cd021dde34603426402615658f1d970cfa7c7bd92ac81a8b16eeebff264d59ce4604ea02d192a39a8cc7a70173007301",
          "assets": [],
          "creationHeight": 284761,
          "additionalRegisters": {},
          "transactionId": "9148408c04c2e38a6402a7950d6157730fa7d49e9ab3b9cadec481d7769918e9",
          "index": 1
        }`
	unspentBox, _ := NewBoxFromJson(boxJson)
	unspentBoxes := NewBoxes()
	unspentBoxes.Add(unspentBox)

	testContract, _ := NewContractPayToAddress(recipient)
	outBoxValue := SafeUserMinBoxValue()
	fee := SuggestedTxFee()

	testBoxSelector := NewSimpleBoxSelector()

	targetBalance, _ := SumOfBoxValues(outBoxValue, fee)
	testBoxSelection, _ := testBoxSelector.Select(unspentBoxes, targetBalance, NewTokens())

	// Mint token
	mintBox, _ := testBoxSelection.Boxes().Get(0)
	testTokenId := NewTokenIdFromBoxId(mintBox.BoxId())
	testTokenAmount, _ := NewTokenAmount(1)
	testToken := NewToken(testTokenId, testTokenAmount)
	testBoxBuilder := NewBoxCandidateBuilder(outBoxValue, testContract, 0)
	testBoxBuilder.MintToken(testToken, "TKN", "token desc", 2)
	outbox, _ := testBoxBuilder.Build()

	txOutputs := NewBoxCandidates()
	txOutputs.Add(outbox)
	changeAddr, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")
	testDataInputs := NewDataInputs()

	testTxBuilder := NewTxBuilder(testBoxSelection, txOutputs, 0, fee, changeAddr)
	testTxBuilder.SetDataInputs(testDataInputs)
	tx, txErr := testTxBuilder.Build()
	assert.NoError(t, txErr)
	_, txJsonErr := tx.JsonEIP12()
	assert.NoError(t, txJsonErr)
}

func TestBurnToken(t *testing.T) {
	recipient, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")
	boxJson := `{
                  "boxId": "0cf7b9e71961cc473242de389c8e594a4e5d630ddd2e4e590083fb0afb386341",
                  "value": 11491500000,
                  "ergoTree": "100f040005c801056404000e2019719268d230fd9093e4db0e2e42a07883ffe976e77c7419efc1bb218a05d4ba04000500043c040204c096b10204020101040205c096b1020400d805d601b2a5730000d602e4c6a70405d6039c9d720273017302d604b5db6501fed9010463ededed93e4c67204050ec5a7938cb2db6308720473030001730492e4c672040605997202720390e4c6720406059a72027203d605b17204ea02d1edededededed93cbc27201e4c6a7060e917205730593db63087201db6308a793e4c6720104059db072047306d9010641639a8c720601e4c68c72060206057e72050593e4c6720105049ae4c6a70504730792c1720199c1a77e9c9a720573087309058cb072048602730a730bd901063c400163d802d6088c720601d6098c72080186029a7209730ceded8c72080293c2b2a5720900d0cde4c68c720602040792c1b2a5720900730d02b2ad7204d9010663cde4c672060407730e00",
                  "assets": [
                    {
                      "tokenId": "19475d9a78377ff0f36e9826cec439727bea522f6ffa3bda32e20d2f8b3103ac",
                      "amount": 1
                    }
                  ],
                  "creationHeight": 348198,
                  "additionalRegisters": {
                    "R4": "059acd9109",
                    "R5": "04f2c02a",
                    "R6": "0e20277c78751ff6f68d4dcd082eeea9506324911a875b6b9cd4d177d4fcab061327"
                  },
                  "transactionId": "5ed0e572a8c097b053965519a696f413f7be02754345e8ed650377e29a6dedb3",
                  "index": 0
                }`
	unspentBox, _ := NewBoxFromJson(boxJson)
	unspentBoxes := NewBoxes()
	unspentBoxes.Add(unspentBox)

	testTokenId, _ := NewTokenId("19475d9a78377ff0f36e9826cec439727bea522f6ffa3bda32e20d2f8b3103ac")
	testTokenAmount, _ := NewTokenAmount(1)
	testToken := NewToken(testTokenId, testTokenAmount)

	testBoxSelector := NewSimpleBoxSelector()

	testTokens := NewTokens()
	testTokens.Add(testToken)

	outBoxValue := SafeUserMinBoxValue()
	fee := SuggestedTxFee()

	targetBalance, _ := SumOfBoxValues(outBoxValue, fee)
	testBoxSelection, _ := testBoxSelector.Select(unspentBoxes, targetBalance, testTokens)
	testContract, _ := NewContractPayToAddress(recipient)
	testBoxBuilder := NewBoxCandidateBuilder(outBoxValue, testContract, 0)
	outbox, _ := testBoxBuilder.Build()
	txOutputs := NewBoxCandidates()
	txOutputs.Add(outbox)
	changeAddr, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")
	testDataInputs := NewDataInputs()

	testTxBuilder := NewTxBuilder(testBoxSelection, txOutputs, 0, fee, changeAddr)
	testTxBuilder.SetDataInputs(testDataInputs)
	_, buildErr := testTxBuilder.Build()
	assert.Error(t, buildErr)
	testTxBuilder.SetTokenBurnPermit(testTokens)
	_, txErr := testTxBuilder.Build()
	assert.NoError(t, txErr)
}

func TestMultiSigTx(t *testing.T) {
	aliceByteSecret, _ := hex.DecodeString("e726ad60a073a49f7851f4d11a83de6c9c7f99e17314fcce560f00a51a8a3d18")
	aliceSecret, _ := NewSecretKeyFromBytes(aliceByteSecret)
	bobByteSecret, _ := hex.DecodeString("9e6616b4e44818d21b8dfdd5ea87eb822480e7856ab910d00f5834dc64db79b3")
	bobSecret, _ := NewSecretKeyFromBytes(bobByteSecret)
	alicePkBytes, _ := hex.DecodeString("cd03c8e1527efae4be9868cea6767157fcccac66489842738efed0a302e4f81710d0")

	// Pay 2 Script address of a multi_sig contract with contract { alicePK && bobPK }
	multiSigAddress, _ := NewAddress("JryiCXrc7x5D8AhS9DYX1TDzW5C5mT6QyTMQaptF76EQkM15cetxtYKq3u6LymLZLVCyjtgbTKFcfuuX9LLi49Ec5m2p6cwsg5NyEsCQ7na83yEPN")
	inputContract, _ := NewContractPayToAddress(multiSigAddress)
	testTxId, _ := NewTxId("0000000000000000000000000000000000000000000000000000000000000000")
	testInputBoxValue, _ := NewBoxValue(1000000000)
	testInputBox, _ := NewBox(testInputBoxValue, 0, inputContract, testTxId, 0, NewTokens())

	// create a transaction that spends the "simulated" box
	recipient, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")
	unspentBoxes := NewBoxes()
	unspentBoxes.Add(testInputBox)
	testContract, _ := NewContractPayToAddress(recipient)
	outBoxValue := SafeUserMinBoxValue()
	fee := SuggestedTxFee()
	outbox, _ := NewBoxCandidateBuilder(outBoxValue, testContract, 0).Build()
	txOutputs := NewBoxCandidates()
	txOutputs.Add(outbox)
	changeAddr, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")
	testDataInputs := NewDataInputs()
	testBoxSelector := NewSimpleBoxSelector()

	targetBalance, _ := SumOfBoxValues(outBoxValue, fee)
	testBoxSelection, _ := testBoxSelector.Select(unspentBoxes, targetBalance, NewTokens())
	testTxBuilder := NewTxBuilder(testBoxSelection, txOutputs, 0, fee, changeAddr)
	testTxBuilder.SetDataInputs(testDataInputs)

	tx, buildErr := testTxBuilder.Build()
	assert.NoError(t, buildErr)

	txDataInputs := NewBoxes()

	testBlockHeaders := testBlockHeadersFromJson()
	testBlockHeader, _ := testBlockHeaders.Get(0)
	testPreHeader := NewPreHeader(testBlockHeader)
	ctx, _ := NewStateContext(testPreHeader, testBlockHeaders)

	sksAlice := NewSecretKeys()
	sksAlice.Add(aliceSecret)
	walletAlice := NewWalletFromSecretKeys(sksAlice)

	sksBob := NewSecretKeys()
	sksBob.Add(bobSecret)
	walletBob := NewWalletFromSecretKeys(sksBob)

	bobGeneratedCommitments, _ := walletBob.GenerateCommitments(ctx, tx, unspentBoxes, txDataInputs)
	bobHints := bobGeneratedCommitments.AllHintsForInput(0)
	bobKnown, _ := bobHints.Get(0)
	bobOwn, _ := bobHints.Get(1)

	testHintsBag := NewHintsBag()
	testHintsBag.Add(bobKnown)

	aliceTxHintsBag := NewTransactionHintsBag()
	aliceTxHintsBag.AddHintsForInput(0, testHintsBag)

	partialSigned, _ := walletAlice.SignTransactionMulti(ctx, tx, unspentBoxes, txDataInputs, aliceTxHintsBag)

	realPropositions := NewPropositions()
	simulatedPropositions := NewPropositions()
	_ = realPropositions.Add(alicePkBytes)

	bobExtractedHints, _ := ExtractHintsFromSignedTransaction(partialSigned, ctx, unspentBoxes, txDataInputs, realPropositions, simulatedPropositions)
	bobHintsBag := bobExtractedHints.AllHintsForInput(0)
	bobHintsBag.Add(bobOwn)

	bobTxHintsBag := NewTransactionHintsBag()
	bobTxHintsBag.AddHintsForInput(0, bobHintsBag)

	_, signErr := walletBob.SignTransactionMulti(ctx, tx, unspentBoxes, txDataInputs, bobTxHintsBag)
	assert.NoError(t, signErr)
}
