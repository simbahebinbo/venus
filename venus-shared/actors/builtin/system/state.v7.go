// FETCHED FROM LOTUS: builtin/system/state.go.template

package system

import (
	"fmt"

	actorstypes "github.com/filecoin-project/go-state-types/actors"
	"golang.org/x/xerrors"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/venus/venus-shared/actors"
	"github.com/filecoin-project/venus/venus-shared/actors/adt"

	system7 "github.com/filecoin-project/specs-actors/v7/actors/builtin/system"
)

var _ State = (*state7)(nil)

func load7(store adt.Store, root cid.Cid) (State, error) {
	out := state7{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make7(store adt.Store) (State, error) {
	out := state7{store: store}
	out.State = system7.State{}
	return &out, nil
}

type state7 struct {
	system7.State
	store adt.Store
}

func (s *state7) GetState() interface{} {
	return &s.State
}

func (s *state7) GetBuiltinActors() cid.Cid {

	return cid.Undef

}

func (s *state7) SetBuiltinActors(c cid.Cid) error {

	return xerrors.New("cannot set manifest cid before v8")

}

func (s *state7) ActorKey() string {
	return actors.SystemKey
}

func (s *state7) ActorVersion() actorstypes.Version {
	return actorstypes.Version7
}

func (s *state7) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
