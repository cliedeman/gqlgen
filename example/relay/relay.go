// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package relay

import (
	"context"
)

type Stub struct {
	FactionResolver struct {
		ID func(ctx context.Context, obj *Faction) (string, error)
	}
	MutationResolver struct {
		IntroduceShip func(ctx context.Context, input IntroduceShipInput) (*IntroduceShipPayload, error)
	}
	QueryResolver struct {
		Rebels func(ctx context.Context) (*Faction, error)
		Empire func(ctx context.Context) (*Faction, error)
		Node   func(ctx context.Context, id string) (Node, error)
	}
	ShipResolver struct {
		ID func(ctx context.Context, obj *Ship) (string, error)
	}
}

func (r *Stub) Faction() FactionResolver {
	return &stubFaction{r}
}
func (r *Stub) Mutation() MutationResolver {
	return &stubMutation{r}
}
func (r *Stub) Query() QueryResolver {
	return &stubQuery{r}
}
func (r *Stub) Ship() ShipResolver {
	return &stubShip{r}
}

type stubFaction struct{ *Stub }

func (r *stubFaction) ID(ctx context.Context, obj *Faction) (string, error) {
	return r.FactionResolver.ID(ctx, obj)
}

type stubMutation struct{ *Stub }

func (r *stubMutation) IntroduceShip(ctx context.Context, input IntroduceShipInput) (*IntroduceShipPayload, error) {
	return r.MutationResolver.IntroduceShip(ctx, input)
}

type stubQuery struct{ *Stub }

func (r *stubQuery) Rebels(ctx context.Context) (*Faction, error) {
	return r.QueryResolver.Rebels(ctx)
}
func (r *stubQuery) Empire(ctx context.Context) (*Faction, error) {
	return r.QueryResolver.Empire(ctx)
}
func (r *stubQuery) Node(ctx context.Context, id string) (Node, error) {
	return r.QueryResolver.Node(ctx, id)
}

type stubShip struct{ *Stub }

func (r *stubShip) ID(ctx context.Context, obj *Ship) (string, error) {
	return r.ShipResolver.ID(ctx, obj)
}
