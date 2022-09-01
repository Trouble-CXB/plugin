package executor

import (
	"github.com/33cn/chain33/util"
	zt "github.com/33cn/plugin/plugin/dapp/zksync/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHistoryProof(t *testing.T) {
	dir, statedb, localdb := util.CreateTestDB()

	defer util.CloseTestDB(dir, statedb)

	zkproof1 := &zt.ZkCommitProof{
		BlockStart:  1,
		BlockEnd:    5,
		IndexEnd:    1,
		OldTreeRoot: "18617692155653794411600229951838919630651308402001068372178576330275141191583",
		NewTreeRoot: "17588927116408652037942430551926564485721758445011528475354328354279637224164",
		PublicInput: "00000002291a2f56d530fa7c128420b452c611ea76f24c2b25474b1cd00324791411f6751ee507d5baef0ac3086086182b9bd3865479bc5cb4b6662a8af7124f362c01fa",
		Proof:       "ef09981e9646dec88300c34c0da6f19a658d8009e24c7e452736311e8d7e06b995db7b3ff473fdef3a054f3455f5a3683f75a4979d2e2619b9a75baab5e869210a60bb0dd624aedb9d53acde714876dab9543407615a186dd3c4b4370735a646ada69f70cfe7b42bfaa02b14a18f0b6382e388be92d3c4df7fb35354dce29648",
		PubDatas: []string{
			"1329227996713370902439899860492615682",
			"40011206537402433073299143204095554199",
			"186177691408961248415817519130861079171",
			"76596692807341630114957937081755163895",
			"269754080702394078355149086586457554944",
			"0",
			"0",
			"0",
		},
		OnChainProofId:    1,
		ProofId:           1,
		CommitBlockHeight: 0,
	}

	zkproof2 := &zt.ZkCommitProof{
		BlockStart:  6,
		BlockEnd:    6,
		IndexEnd:    1,
		OldTreeRoot: "17588927116408652037942430551926564485721758445011528475354328354279637224164",
		NewTreeRoot: "11145935540026218277834844528205882443571486137141000350047738346658056635507",
		PublicInput: "00000002074fa7ba53e4a88f0c015ccfe7e1a813c89513b0da4f592b89c44729ccd333e12ec319b6d9878e9724d0ca7be95816a1c8dc33b19c12ff18010e3f624779f452",
		Proof:       "d772c954e47e7b950a3ea715565103bf1d28932c99df3961dab6ceef030df1db9cbfd5dfe9061e12940e3b96c7cdc91779ff1185bf17ccfec33fcedee6bd9cf70bd8daec68e883c492cec3581c42abf88e2e7a5971f88f21c5ee5d2ddc09fa85c53bd8665bc848b70b6317bcaaf3d98d5e8d6ae45d42cde294ffc0f08e1ed6e0",
		PubDatas: []string{
			"1329227997022855912333302523255324674",
			"40011206537402433073256028580150952411",
			"280707582035715202603196119834593787890",
			"16967169861597165590835161204918053792",
			"252010078148043611589173402425784008704",
			"0",
			"0",
			"0",
		},
		OnChainProofId:    2,
		ProofId:           2,
		CommitBlockHeight: 0,
	}

	zkproof3 := &zt.ZkCommitProof{
		BlockStart:  7,
		BlockEnd:    7,
		IndexEnd:    1,
		OldTreeRoot: "11145935540026218277834844528205882443571486137141000350047738346658056635507",
		NewTreeRoot: "10629830243661241796890761239146633299552964472177879280335063167396679434939",
		PublicInput: "000000020baf5d8249abb690bb7e8371efe6f3b27330a6999f1456e688d515cae5448b860ebd9bd157052bea8199c85127c56d6cd2508a4957475f50c5ab6b4a8169e4b3",
		Proof:       "e2d3945114bcf1c389b0d4370f09d00966458df49c0a2c695afc8dfe2ca2f858aca512764bdc41b044083c9c806d0ed5092da6a07b56645e198950a01f92f205025e33c12bc75f09bdf47172f71e26ec4f2919d1bade61d35b329ee2be02f729a1477e51782176a9bd2ced57edf461ae93e2b6613c7d89f8792b9187ec4882a9",
		PubDatas: []string{
			"7975367975638082090280654717933543597",
			"27130104817958728946896606866427107541",
			"20405913717877370810253725600767691923",
			"148307611380547906348497407918804602520",
			"18777818215363764900717223316757151744",
			"0",
			"0",
			"0",
		},
		OnChainProofId:    3,
		ProofId:           3,
		CommitBlockHeight: 0,
	}

	zkproof4 := &zt.ZkCommitProof{
		BlockStart:  7,
		BlockEnd:    7,
		IndexEnd:    1,
		OldTreeRoot: "10629830243661241796890761239146633299552964472177879280335063167396679434939",
		NewTreeRoot: "2467545001057789562952276125964088483262320452506060494115669489570341992450",
		PublicInput: "000000020dcbd612dc4f15e5c83816bf1e187fa5893a73d590ed77b6e738f4d13a9c3f3b1a7870daf4a9d8bb8bb53b2479e4b46915d81994615317f19aba4b8340081adb",
		Proof:       "e860e55fba4082c278652acd5e62e0c5c3c2221730e3409c14bcd763d832d3dbe302525252ee759cccf0403461fab7824d295e47b16beb0d97163d77f319bcf82bd489b4b5365858a34e9aa51fc55817b3c7cf88d3b681a0c2b7baee129c663ccca924020230e149717dec8fce6e48a10e7220317a9e8326882051c013597260",
		PubDatas: []string{
			"7975367975947529765054940149485951149",
			"133016442201275975300769839516156716285",
			"215347686733310189032849224345406234330",
			"68783560211637872458272545526244982136",
			"336789183525285842789036535418904903680",
			"0",
			"0",
			"0",
		},
		OnChainProofId:    4,
		ProofId:           4,
		CommitBlockHeight: 0,
	}

	zkproof5 := &zt.ZkCommitProof{
		BlockStart:  9,
		BlockEnd:    10,
		IndexEnd:    1,
		OldTreeRoot: "2467545001057789562952276125964088483262320452506060494115669489570341992450",
		NewTreeRoot: "9302649264771194656034123462933659970532359367746520077277645946816426484078",
		PublicInput: "0000000214df6855718cc9e5eb26a9d09dcaf327dbdc164b2a4408e955d71615f269b7bd199fc0bd9938aa1135640c7f88846be2049e6ad06113edda5248698e8b83680a",
		Proof:       "94805905d96f75e265ff622b829cc9bfb27f7a59a6cf59a35605fe721fca4157e68f7ae6eef8d388f8de1c38b4c242d37de740050b2eaca7b25b32fda0e9d4992eed273dae4c40faaa20e183496eb70ba8bc88bbac150973418d7feb9bd9ccdcefdc2147443979ca011d99d94fd95b1bffb8471755f633ba15ccfeb94e3a2f8e",
		PubDatas: []string{
			"10633823967517267022732009539289219072",
			"1329227995885256715931947727497228544",
			"0",
			"0",
			"14621507953943559611836459002311475200",
			"0",
			"0",
			"0",
		},
		OnChainProofId:    0,
		ProofId:           5,
		CommitBlockHeight: 0,
	}

	zkproof6 := &zt.ZkCommitProof{
		BlockStart:  11,
		BlockEnd:    11,
		IndexEnd:    1,
		OldTreeRoot: "9302649264771194656034123462933659970532359367746520077277645946816426484078",
		NewTreeRoot: "19719700540305761665169694243927237153902441322376380373346058669141823295465",
		PublicInput: "000000020193a40f390596bbe74f6a1f7996f138d3b3469f2ee091dd46e47204e96299f8304683e8a01a155c8f31e5587fc58633754dd01a14fd93a7a4de451d761cf2c4",
		Proof:       "dec706dd791c9d15f664d5def9eec8320844c531a225086518fb5836a4d2e01bc9855f68cb1f8f40145c3c86afadc5f22b4f9163e1e7602f141e078b29fe908d286648faf583611fec1d13a770d3763cb13803fc67c87e8895e3c7cb26780ec5eb573a86db72805181c84e9ce901a99b22904d0ad0277313df12e336b4c82ef4",
		PubDatas: []string{
			"1329227996713370902439899860492615682",
			"40011206537402433073299143204095554199",
			"186177691408961248415817519130861079171",
			"76596692807341630114957937081755163895",
			"269754080702394078355149086586457554944",
			"0",
			"0",
			"0",
		},
		OnChainProofId:    1,
		ProofId:           6,
		CommitBlockHeight: 0,
	}

	proofTable := NewCommitProofTable(localdb)
	err := proofTable.Add(zkproof1)
	assert.Equal(t, nil, err)
	err = proofTable.Add(zkproof2)
	assert.Equal(t, nil, err)
	err = proofTable.Add(zkproof3)
	assert.Equal(t, nil, err)
	err = proofTable.Add(zkproof4)
	assert.Equal(t, nil, err)
	err = proofTable.Add(zkproof5)
	assert.Equal(t, nil, err)
	err = proofTable.Add(zkproof6)
	assert.Equal(t, nil, err)
	kvs, err := proofTable.Save()
	for _, kv := range kvs {
		err = localdb.Set(kv.GetKey(), kv.GetValue())
		assert.Equal(t, nil, err)
	}
	rootHash := "19719700540305761665169694243927237153902441322376380373346058669141823295465"
	//ethFeeAddr := "832367164346888E248bd58b9A5f480299F1e88d"
	//chain33FeeAddr := "2c4a5c378be2424fa7585320630eceba764833f1ec1ffb2fafc1af97f27baf5a"
	req := &zt.ZkReqExistenceProof{
		AccountId:3,
		TokenId:1,
		RootHash:rootHash,
		ChainTitleId:1,
	}
	proof, err := getAccountProofInHistory(localdb, req)
	assert.Equal(t, nil, err)
	t.Log(proof)
}
