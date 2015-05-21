package chunks

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/serial"
)

// AddLevel adds one level to the consumer
func AddLevel(consumer chunk.Consumer, levelID int) {
	levelBaseID := res.ResourceID(4000 + 100*levelID)

	AddStaticChunk(consumer, levelBaseID+2, []byte{0x0B, 0x00, 0x00, 0x00})
	AddStaticChunk(consumer, levelBaseID+3, []byte{0x1B, 0x00, 0x00, 0x00})

	AddBasicLevelInformation(consumer, levelBaseID)
	AddMap(consumer, levelBaseID)
	AddStaticChunk(consumer, levelBaseID+6, make([]byte, 8))
	AddLevelTextures(consumer, levelBaseID)
	AddMasterObjectTables(consumer, levelBaseID)
	AddLevelObjects(consumer, levelBaseID)

	AddStaticChunk(consumer, levelBaseID+40, []byte{0x0D, 0x00, 0x00, 0x00})
	AddStaticChunk(consumer, levelBaseID+41, []byte{0x00})
	AddStaticChunk(consumer, levelBaseID+42, make([]byte, 0x1C))

	AddSurveillanceChunk(consumer, levelBaseID)
	AddStaticChunk(consumer, levelBaseID+45, make([]byte, 0x5E))
	AddMapNotes(consumer, levelBaseID)

	AddStaticChunk(consumer, levelBaseID+48, make([]byte, 0x30))
	AddStaticChunk(consumer, levelBaseID+49, make([]byte, 0x01C0))
	AddStaticChunk(consumer, levelBaseID+50, make([]byte, 2))

	AddLoopConfiguration(consumer, levelBaseID)

	// CD-Release only content
	AddStaticChunk(consumer, levelBaseID+52, make([]byte, 2))
	AddStaticChunk(consumer, levelBaseID+53, make([]byte, 0x40))
}

func addData(consumer chunk.Consumer, chunkID res.ResourceID, data interface{}) {
	addTypedData(consumer, chunkID, chunk.BasicChunkType, data)
}

func addTypedData(consumer chunk.Consumer, chunkID res.ResourceID, typeID chunk.TypeID, data interface{}) {
	store := serial.NewByteStore()
	coder := serial.NewPositioningEncoder(store)

	serial.MapData(data, coder)

	blocks := [][]byte{store.Data()}
	consumer.Consume(chunkID, chunk.NewBlockHolder(typeID, res.Map, blocks))
}

// AddBasicLevelInformation adds the basic level info block
func AddBasicLevelInformation(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	info := data.DefaultLevelInformation()

	addTypedData(consumer, levelBaseID+4, chunk.BasicChunkType.WithCompression(), info)
}

// AddMap adds a map
func AddMap(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	tileFactory := func() interface{} {
		entry := data.DefaultTileMapEntry()

		entry.Type = data.Open

		return entry
	}

	table := data.NewTable(64*64, tileFactory)
	for index, entry := range table.Entries {
		if (index < 64) || ((index % 64) == 0) || ((index % 64) == 63) || (index > (64 * 63)) {
			tile := entry.(*data.TileMapEntry)
			tile.Type = data.Solid
		}
	}
	addTypedData(consumer, levelBaseID+5, chunk.BasicChunkType.WithCompression(), table)
}

// AddLevelTextures adds level texture information
func AddLevelTextures(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	data := make([]byte, 54*2)
	data[0] = 0x01
	AddStaticChunk(consumer, levelBaseID+7, data)
}

// AddMasterObjectTables adds main object tables
func AddMasterObjectTables(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	{
		masterCount := 872
		masters := make([]*data.LevelObjectEntry, masterCount)
		for index := range masters {
			master := data.DefaultLevelObjectEntry()

			masters[index] = master
			master.Next = uint16((index + 1) % masterCount)
			master.Previous = uint16((masterCount + index - 1) % masterCount)
		}

		masterTable := &data.Table{Entries: make([]interface{}, masterCount)}
		for i := range masters {
			masterTable.Entries[i] = masters[i]
		}
		addTypedData(consumer, levelBaseID+8, chunk.BasicChunkType.WithCompression(), masterTable)
	}
	{
		refCount := 1600
		references := make([]*data.LevelObjectCrossReference, refCount)
		for index := range references {
			ref := data.DefaultLevelObjectCrossReference()
			references[index] = ref
			ref.NextObjectIndex = uint16((index + 1) % refCount)
		}

		refTable := &data.Table{Entries: make([]interface{}, refCount)}
		for i := range references {
			refTable.Entries[i] = references[i]
		}
		addTypedData(consumer, levelBaseID+9, chunk.BasicChunkType.WithCompression(), refTable)
	}
}

// AddLevelObjects adds level object tables
func AddLevelObjects(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	addLevelObjectTables(consumer, levelBaseID, 0, data.LevelWeaponEntrySize, 16)
	addLevelObjectTables(consumer, levelBaseID, 1, 6, 32)
	addLevelObjectTables(consumer, levelBaseID, 2, 0x28, 32)
	addLevelObjectTables(consumer, levelBaseID, 3, 12, 32)
	addLevelObjectTables(consumer, levelBaseID, 4, 6, 32)
	addLevelObjectTables(consumer, levelBaseID, 5, 7, 8)
	addLevelObjectTables(consumer, levelBaseID, 6, 9, 16)
	addLevelObjectTables(consumer, levelBaseID, 7, data.LevelSceneryEntrySize, 176)
	addLevelObjectTables(consumer, levelBaseID, 8, data.LevelItemEntrySize, 128)
	addLevelObjectTables(consumer, levelBaseID, 9, 0x1E, 64)
	addLevelObjectTables(consumer, levelBaseID, 10, 14, 64)
	addLevelObjectTables(consumer, levelBaseID, 11, 10, 32)
	addLevelObjectTables(consumer, levelBaseID, 12, 0x1C, 160)
	addLevelObjectTables(consumer, levelBaseID, 13, 21, 64)
	addLevelObjectTables(consumer, levelBaseID, 14, 0x2E, 64)
}

type tempStruct struct {
	data.LevelObjectPrefix
	Extra []byte
}

func addLevelObjectTables(consumer chunk.Consumer, levelBaseID res.ResourceID, classID int, entrySize int, entryCount int) {
	table := data.Table{Entries: make([]interface{}, entryCount)}

	for i := range table.Entries {
		table.Entries[i] = &tempStruct{
			LevelObjectPrefix: data.LevelObjectPrefix{
				Next:                  uint16((i + 1) % entryCount),
				Previous:              uint16((entryCount + i - 1) % entryCount),
				LevelObjectTableIndex: 0},
			Extra: make([]byte, entrySize-data.LevelObjectPrefixSize)}
	}
	addData(consumer, levelBaseID+10+res.ResourceID(classID), table)
	//AddStaticChunk(consumer, levelBaseID+10+res.ResourceID(classID), make([]byte, entrySize*entryCount))
	AddStaticChunk(consumer, levelBaseID+25+res.ResourceID(classID), make([]byte, entrySize))
}

// AddSurveillanceChunk adds a chunk for surveillance information
func AddSurveillanceChunk(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	AddStaticChunk(consumer, levelBaseID+43, make([]byte, 8*2))
	AddStaticChunk(consumer, levelBaseID+44, make([]byte, 8*2))
}

// AddMapNotes prepares empty map notes chunks
func AddMapNotes(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	AddStaticChunk(consumer, levelBaseID+46, make([]byte, 0x0800))
	AddStaticChunk(consumer, levelBaseID+47, make([]byte, 4))
}

// AddLoopConfiguration adds an empty loop configuration chunk
func AddLoopConfiguration(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	AddStaticChunk(consumer, levelBaseID+51, make([]byte, 0x03C0))
}
