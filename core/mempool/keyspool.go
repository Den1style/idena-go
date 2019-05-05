package mempool

import (
	"errors"
	"github.com/deckarep/golang-set"
	"idena-go/blockchain/types"
	"idena-go/blockchain/validation"
	"idena-go/common"
	"idena-go/common/eventbus"
	"idena-go/core/appstate"
	"idena-go/events"
	"idena-go/log"
	"sync"
)

type KeysPool struct {
	appState  *appstate.AppState
	flipKeys  map[common.Address]*types.FlipKey
	knownKeys mapset.Set
	mutex     sync.Mutex
	bus       eventbus.Bus
	head      *types.Header
	log       log.Logger
}

func NewKeysPool(appState *appstate.AppState, bus eventbus.Bus) *KeysPool {
	return &KeysPool{
		appState:  appState,
		bus:       bus,
		knownKeys: mapset.NewSet(),
		log:       log.New(),
		flipKeys:  make(map[common.Address]*types.FlipKey),
	}
}

func (p *KeysPool) Initialize(head *types.Header) {
	p.head = head
}

func (p *KeysPool) Add(key *types.FlipKey) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	hash := key.Hash()

	if p.knownKeys.Contains(hash) {
		return errors.New("flipKey with same hash already exists")
	}

	sender, _ := types.SenderFlipKey(key)

	if _, ok := p.flipKeys[sender]; ok {
		return errors.New("sender has already published his key")
	}

	appState := p.appState.ForCheck(p.head.Height())

	if err := validation.ValidateFlipKey(appState, key); err != nil {
		log.Warn("FlipKey is not valid", "hash", hash.Hex(), "err", err)
		return err
	}

	p.knownKeys.Add(hash)
	p.flipKeys[sender] = key

	p.bus.Publish(&events.NewFlipKeyEvent{
		Key: key,
	})

	return nil
}

func (p *KeysPool) GetFlipKeys() []*types.FlipKey {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var list []*types.FlipKey

	for _, tx := range p.flipKeys {
		list = append(list, tx)
	}
	return list
}

func (p *KeysPool) GetFlipKey(address common.Address) *types.FlipKey {
	return p.flipKeys[address]
}

func (p *KeysPool) Clear() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.knownKeys = mapset.NewSet()
	p.flipKeys = make(map[common.Address]*types.FlipKey)
}