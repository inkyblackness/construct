package chunks

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/serial"
)

// AddGameState adds the chunk for the game state
func AddGameState(consumer chunk.Consumer) {
	store := serial.NewByteStore()
	coder := serial.NewPositioningEncoder(store)

	coder.CodeBytes(make([]byte, data.GameStateSize))

	blocks := [][]byte{store.Data()}
	consumer.Consume(res.ResourceID(0x0FA1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Map, blocks))
}
