package jhinih_embedder

import (
	"context"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"os"
)

func NewArkEmbedder(ctx context.Context) *ark.Embedder {
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("EMBEDDER"),
	})
	if err != nil {
		panic(err)
	}
	return embedder
}
