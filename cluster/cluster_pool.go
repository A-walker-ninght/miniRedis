package cluster

import (
	"context"
	"errors"
	client2 "github.com/A-walker-ninght/miniRedis/resp/client"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type connectionFactory struct {
	Peer string
}

func (c connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	client, err := client2.MakeClient(c.Peer)
	if err != nil {
		return nil, err
	}
	client.Start()
	return pool.NewPooledObject(client), nil
}

func (c connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	client, ok := object.Object.(*client2.Client)
	if !ok {
		return errors.New("type mismatch")
	}
	client.Close()
	return nil
}

func (c connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (p connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (p connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
