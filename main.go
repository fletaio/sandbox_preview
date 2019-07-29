package main // import "github.com/fletaio/sandbox_preview"

import (
	"encoding/hex"
	"log"
	"strconv"
	"time"

	"github.com/fletaio/fleta/common"
	"github.com/fletaio/fleta/common/amount"
	"github.com/fletaio/fleta/common/key"
	"github.com/fletaio/fleta/core/chain"
	"github.com/fletaio/fleta/core/types"
	"github.com/fletaio/fleta/pof"
	"github.com/fletaio/fleta/process/admin"
	"github.com/fletaio/fleta/process/formulator"
	"github.com/fletaio/fleta/process/vault"
	"github.com/fletaio/fleta/service/p2p"
)

func main() {
	if true { // sandbox.master
		k, _ := key.NewMemoryKeyFromString("0000000000000000000000000000000000000000000000000000000000000100")
		KeyHash := common.NewPublicHash(k.PublicKey())
		log.Println("sandbox.master", KeyHash.String())
	}
	if true { // hyper
		k, _ := key.NewMemoryKeyFromString("0000000000000000000000000000000000000000000000000000000000000010")
		KeyHash := common.NewPublicHash(k.PublicKey())
		g, _ := key.NewMemoryKeyFromString("0000000000000000000000000000000000000000000000000000000000000011")
		GenHash := common.NewPublicHash(g.PublicKey())
		log.Println("hyper", KeyHash.String(), GenHash.String())
	}
	if true { //sigma
		k, _ := key.NewMemoryKeyFromString("0000000000000000000000000000000000000000000000000000000000000020")
		KeyHash := common.NewPublicHash(k.PublicKey())
		g, _ := key.NewMemoryKeyFromString("0000000000000000000000000000000000000000000000000000000000000021")
		GenHash := common.NewPublicHash(g.PublicKey())
		log.Println("sigma", KeyHash.String(), GenHash.String())
	}
	if true { //alpha
		k, _ := key.NewMemoryKeyFromString("0000000000000000000000000000000000000000000000000000000000000030")
		KeyHash := common.NewPublicHash(k.PublicKey())
		g, _ := key.NewMemoryKeyFromString("0000000000000000000000000000000000000000000000000000000000000031")
		GenHash := common.NewPublicHash(g.PublicKey())
		log.Println("alpha", KeyHash.String(), GenHash.String())
	}
	if true { //single
		k, _ := key.NewMemoryKeyFromString("0000000000000000000000000000000000000000000000000000000000000200")
		KeyHash := common.NewPublicHash(k.PublicKey())
		log.Println("single", KeyHash.String())
	}

	obstrs := []string{
		"0000000000000000000000000000000000000000000000000000000000000001",
		"0000000000000000000000000000000000000000000000000000000000000002",
		"0000000000000000000000000000000000000000000000000000000000000003",
		"0000000000000000000000000000000000000000000000000000000000000004",
		"0000000000000000000000000000000000000000000000000000000000000005",
	}
	obkeys := make([]key.Key, 0, len(obstrs))
	NetAddressMap := map[common.PublicHash]string{}
	FrNetAddressMap := map[common.PublicHash]string{}
	ObserverKeys := make([]common.PublicHash, 0, len(obstrs))
	for i, v := range obstrs {
		if bs, err := hex.DecodeString(v); err != nil {
			panic(err)
		} else if Key, err := key.NewMemoryKeyFromBytes(bs); err != nil {
			panic(err)
		} else {
			obkeys = append(obkeys, Key)
			pubhash := common.NewPublicHash(Key.PublicKey())
			ObserverKeys = append(ObserverKeys, pubhash)
			NetAddressMap[pubhash] = ":400" + strconv.Itoa(i)
			FrNetAddressMap[pubhash] = "ws://localhost:500" + strconv.Itoa(i)
		}
	}

	MaxBlocksPerFormulator := uint32(10)
	ChainID := uint8(0xFF)

	for i, obkey := range obkeys {
		st, err := chain.NewStore("./_test/odata_"+strconv.Itoa(i), ChainID, "FLETA Sandbox Preview", 0x0001, true)
		if err != nil {
			panic(err)
		}
		defer st.Close()

		cs := pof.NewConsensus(MaxBlocksPerFormulator, ObserverKeys)
		app := NewSandboxApp()
		cn := chain.NewChain(cs, app, st)
		cn.MustAddProcess(admin.NewAdmin(1))
		vp := vault.NewVault(2)
		cn.MustAddProcess(vp)
		fp := formulator.NewFormulator(3)
		cn.MustAddProcess(fp)
		if err := cn.Init(); err != nil {
			panic(err)
		}

		ob := pof.NewObserverNode(obkey, NetAddressMap, cs)
		if err := ob.Init(); err != nil {
			panic(err)
		}

		go ob.Run(":400"+strconv.Itoa(i), ":500"+strconv.Itoa(i))
	}

	frstrs := []string{
		"0000000000000000000000000000000000000000000000000000000000000011",
	}
	frkeys := make([]key.Key, 0, len(frstrs))
	for _, v := range frstrs {
		if bs, err := hex.DecodeString(v); err != nil {
			panic(err)
		} else if Key, err := key.NewMemoryKeyFromBytes(bs); err != nil {
			panic(err)
		} else {
			frkeys = append(frkeys, Key)
		}
	}

	ndstrs := []string{
		"0000000000000000000000000000000000000000000000000000000000000999",
	}
	NdNetAddressMap := map[common.PublicHash]string{}
	ndkeys := make([]key.Key, 0, len(ndstrs))
	for i, v := range ndstrs {
		if bs, err := hex.DecodeString(v); err != nil {
			panic(err)
		} else if Key, err := key.NewMemoryKeyFromBytes(bs); err != nil {
			panic(err)
		} else {
			ndkeys = append(ndkeys, Key)
			pubhash := common.NewPublicHash(Key.PublicKey())
			NdNetAddressMap[pubhash] = ":601" + strconv.Itoa(i)
		}
	}

	for i, frkey := range frkeys {
		st, err := chain.NewStore("./_test/fdata_"+strconv.Itoa(i), ChainID, "FLETA Sandbox Preview", 0x0001, true)
		if err != nil {
			panic(err)
		}
		defer st.Close()

		cs := pof.NewConsensus(MaxBlocksPerFormulator, ObserverKeys)
		app := NewSandboxApp()
		cn := chain.NewChain(cs, app, st)
		cn.MustAddProcess(admin.NewAdmin(1))
		vp := vault.NewVault(2)
		cn.MustAddProcess(vp)
		fp := formulator.NewFormulator(3)
		cn.MustAddProcess(fp)
		if err := cn.Init(); err != nil {
			panic(err)
		}

		fr := pof.NewFormulatorNode(&pof.FormulatorConfig{
			Formulator:              common.MustParseAddress("385ujsGNZt"),
			MaxTransactionsPerBlock: 10000,
		}, frkey, frkey, FrNetAddressMap, NdNetAddressMap, cs, "./_test/fdata_"+strconv.Itoa(i)+"/peer")
		if err := fr.Init(); err != nil {
			panic(err)
		}

		go fr.Run(":600" + strconv.Itoa(i))
	}

	for i, ndkey := range ndkeys {
		st, err := chain.NewStore("./_test/ndata_"+strconv.Itoa(i), ChainID, "FLETA Sandbox Preview", 0x0001, true)
		if err != nil {
			panic(err)
		}
		defer st.Close()

		cs := pof.NewConsensus(MaxBlocksPerFormulator, ObserverKeys)
		app := NewSandboxApp()
		cn := chain.NewChain(cs, app, st)
		cn.MustAddProcess(admin.NewAdmin(1))
		vp := vault.NewVault(2)
		cn.MustAddProcess(vp)
		fp := formulator.NewFormulator(3)
		cn.MustAddProcess(fp)
		if err := cn.Init(); err != nil {
			panic(err)
		}

		nd := p2p.NewNode(ndkey, NdNetAddressMap, cn, "./_test/ndata_"+strconv.Itoa(i)+"/peer")
		if err := nd.Init(); err != nil {
			panic(err)
		}

		go func() {
			time.Sleep(60 * time.Second)
			nd.Run(":601" + strconv.Itoa(i))
		}()
	}

	select {}
}

// SandboxApp is app
type SandboxApp struct {
	*types.ApplicationBase
	pm      types.ProcessManager
	cn      types.Provider
	addrMap map[string]common.Address
}

// NewSandboxApp returns a SandboxApp
func NewSandboxApp() *SandboxApp {
	return &SandboxApp{
		addrMap: map[string]common.Address{
			"fleta.formulator": common.MustParseAddress("3CUsUpv9v"),
		},
	}
}

// Name returns the name of the application
func (app *SandboxApp) Name() string {
	return "SandboxApp"
}

// Version returns the version of the application
func (app *SandboxApp) Version() string {
	return "v1.0.0"
}

// Init initializes the consensus
func (app *SandboxApp) Init(reg *types.Register, pm types.ProcessManager, cn types.Provider) error {
	app.pm = pm
	app.cn = cn
	return nil
}

// InitGenesis initializes genesis data
func (app *SandboxApp) InitGenesis(ctw *types.ContextWrapper) error {
	rewardPolicy := &formulator.RewardPolicy{
		RewardPerBlock:        amount.NewCoinAmount(0, 951293759512937600), // 0.03%
		PayRewardEveryBlocks:  172800,                                      // 1 day
		AlphaEfficiency1000:   1000,                                        // 100%
		SigmaEfficiency1000:   1150,                                        // 115%
		OmegaEfficiency1000:   1300,                                        // 130%
		HyperEfficiency1000:   1300,                                        // 130%
		StakingEfficiency1000: 700,                                         // 70%
	}
	alphaPolicy := &formulator.AlphaPolicy{
		AlphaCreationLimitHeight:  5184000,                         // 30 days
		AlphaCreationAmount:       amount.NewCoinAmount(200000, 0), // 200,000 FLETA
		AlphaUnlockRequiredBlocks: 2592000,                         // 15 days
	}
	sigmaPolicy := &formulator.SigmaPolicy{
		SigmaRequiredAlphaBlocks:  5184000, // 30 days
		SigmaRequiredAlphaCount:   4,       // 4 Alpha (800,000 FLETA)
		SigmaUnlockRequiredBlocks: 2592000, // 15 days
	}
	omegaPolicy := &formulator.OmegaPolicy{
		OmegaRequiredSigmaBlocks:  5184000, // 30 days
		OmegaRequiredSigmaCount:   2,       // 2 Sigma (1,600,000 FLETA)
		OmegaUnlockRequiredBlocks: 2592000, // 15 days
	}
	hyperPolicy := &formulator.HyperPolicy{
		HyperCreationAmount:         amount.NewCoinAmount(5000000, 0), // 5,000,000 FLETA
		HyperUnlockRequiredBlocks:   2592000,                          // 15 days
		StakingUnlockRequiredBlocks: 2592000,                          // 15 days
	}

	if p, err := app.pm.ProcessByName("fleta.formulator"); err != nil {
		return err
	} else if fp, is := p.(*formulator.Formulator); !is {
		return types.ErrNotExistProcess
	} else {
		if err := fp.InitPolicy(ctw,
			rewardPolicy,
			alphaPolicy,
			sigmaPolicy,
			omegaPolicy,
			hyperPolicy,
		); err != nil {
			return err
		}
	}
	if p, err := app.pm.ProcessByName("fleta.admin"); err != nil {
		return err
	} else if ap, is := p.(*admin.Admin); !is {
		return types.ErrNotExistProcess
	} else {
		if err := ap.InitAdmin(ctw, app.addrMap); err != nil {
			return err
		}
	}
	if p, err := app.pm.ProcessByName("fleta.vault"); err != nil {
		return err
	} else if sp, is := p.(*vault.Vault); !is {
		return types.ErrNotExistProcess
	} else {
		addSingleAccount(sp, ctw, common.MustParsePublicHash("fYVjaY2XAoqjYQfnqx1sjxChoPLWZDxdtFxb3e5V8b"), common.MustParseAddress("3CUsUpv9v"), "fleta.formulator", amount.NewCoinAmount(100000000, 0))

		addHyperFormulator(sp, ctw, hyperPolicy, 0, common.MustParsePublicHash("4rukqNN7Hwd7gengrsQo5kZhGHj1GvfuqbcuKD4ohwP"), common.MustParsePublicHash("t9bJ8qsLxd3pnGuANVfDR5QB6ryQaRJRZS1QZ9oGTo"), common.MustParseAddress("385ujsGNZt"), "HashTower")

		addSigmaFormulator(sp, ctw, sigmaPolicy, alphaPolicy, common.MustParsePublicHash("3QMGcxQpEk8ZzYTdxGquER18pqNwkaigJuJ1jcEM4HY"), common.MustParsePublicHash("4K2z2xMqRVKfEC3BEB2XRjAeMwEde2aGEDBoiw3iGQF"), common.MustParseAddress("5CyLcFhpyN"), "node1")

		addAlphaFormulator(sp, ctw, alphaPolicy, 5184000, common.MustParsePublicHash("2Db6wGB3b3Qnx2nosLbsu7RkbtVjBoTGRdbjVf2yddU"), common.MustParsePublicHash("qbjKBGJmv5ENkwm8NdRmPk1rChrHAjh22rBGnFkh4E"), common.MustParseAddress("3aTgBdBUME"), "1G")

		addSingleAccount(sp, ctw, common.MustParsePublicHash("WWPsV3TCJcTew43vaFa98zhnWZ8oGXHNo5MCHzyBaC"), common.MustParseAddress("2nJauYqHDB"), "sandbox.account", amount.MustParseAmount("10000"))
	}
	if p, err := app.pm.ProcessByName("fleta.formulator"); err != nil {
		return err
	} else if fp, is := p.(*formulator.Formulator); !is {
		return types.ErrNotExistProcess
	} else {
		addStaking(fp, ctw, common.MustParseAddress("385ujsGNZt"), common.MustParseAddress("2nJauYqHDB"), amount.MustParseAmount("10000"))

		HyperAddresses := []common.Address{
			common.MustParseAddress("385ujsGNZt"),
		}
		if err := fp.InitStakingMap(ctw, HyperAddresses); err != nil {
			return err
		}
	}
	return nil
}

// OnLoadChain called when the chain loaded
func (app *SandboxApp) OnLoadChain(loader types.LoaderWrapper) error {
	return nil
}

func addSingleAccount(sp *vault.Vault, ctw *types.ContextWrapper, KeyHash common.PublicHash, addr common.Address, name string, am *amount.Amount) {
	acc := &vault.SingleAccount{
		Address_: addr,
		Name_:    name,
		KeyHash:  KeyHash,
	}
	if err := ctw.CreateAccount(acc); err != nil {
		panic(err)
	}
	if !am.IsZero() {
		if err := sp.AddBalance(ctw, acc.Address(), am); err != nil {
			panic(err)
		}
	}
}

func addAlphaFormulator(sp *vault.Vault, ctw *types.ContextWrapper, alphaPolicy *formulator.AlphaPolicy, PreHeight uint32, KeyHash common.PublicHash, GenHash common.PublicHash, addr common.Address, name string) {
	acc := &formulator.FormulatorAccount{
		Address_:       addr,
		Name_:          name,
		FormulatorType: formulator.AlphaFormulatorType,
		KeyHash:        KeyHash,
		GenHash:        GenHash,
		Amount:         alphaPolicy.AlphaCreationAmount,
		PreHeight:      PreHeight,
		UpdatedHeight:  0,
	}
	if err := ctw.CreateAccount(acc); err != nil {
		panic(err)
	}
}

func addSigmaFormulator(sp *vault.Vault, ctw *types.ContextWrapper, sigmaPolicy *formulator.SigmaPolicy, alphaPolicy *formulator.AlphaPolicy, KeyHash common.PublicHash, GenHash common.PublicHash, addr common.Address, name string) {
	acc := &formulator.FormulatorAccount{
		Address_:       addr,
		Name_:          name,
		FormulatorType: formulator.SigmaFormulatorType,
		KeyHash:        KeyHash,
		GenHash:        GenHash,
		Amount:         alphaPolicy.AlphaCreationAmount.MulC(int64(sigmaPolicy.SigmaRequiredAlphaCount)),
		PreHeight:      0,
		UpdatedHeight:  0,
	}
	if err := ctw.CreateAccount(acc); err != nil {
		panic(err)
	}
}

func addHyperFormulator(sp *vault.Vault, ctw *types.ContextWrapper, hyperPolicy *formulator.HyperPolicy, Commission1000 uint32, KeyHash common.PublicHash, GenHash common.PublicHash, addr common.Address, name string) {
	acc := &formulator.FormulatorAccount{
		Address_:       addr,
		Name_:          name,
		FormulatorType: formulator.HyperFormulatorType,
		KeyHash:        KeyHash,
		GenHash:        GenHash,
		Amount:         hyperPolicy.HyperCreationAmount,
		PreHeight:      0,
		UpdatedHeight:  0,
		StakingAmount:  amount.NewCoinAmount(0, 0),
		Policy: &formulator.ValidatorPolicy{
			CommissionRatio1000: Commission1000,
			MinimumStaking:      amount.NewCoinAmount(100, 0),
			PayOutInterval:      1,
		},
	}
	if err := ctw.CreateAccount(acc); err != nil {
		panic(err)
	}
}

func addStaking(fp *formulator.Formulator, ctw *types.ContextWrapper, HyperAddress common.Address, StakingAddress common.Address, am *amount.Amount) {
	if has, err := ctw.HasAccount(StakingAddress); err != nil {
		panic(err)
	} else if !has {
		panic(types.ErrNotExistAccount)
	}
	if acc, err := ctw.Account(HyperAddress); err != nil {
		panic(err)
	} else if frAcc, is := acc.(*formulator.FormulatorAccount); !is {
		panic(formulator.ErrInvalidFormulatorAddress)
	} else if frAcc.FormulatorType != formulator.HyperFormulatorType {
		panic(formulator.ErrNotHyperFormulator)
	} else {
		frAcc.StakingAmount = frAcc.StakingAmount.Add(am)
	}
	fp.AddStakingAmount(ctw, HyperAddress, StakingAddress, am)
}
