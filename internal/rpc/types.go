package rpc

type Block struct {
	Hash              string   `json:"hash"`
	Confirmations     int64    `json:"confirmations"`
	Height            int64    `json:"height"`
	Version           int64    `json:"version"`
	VersionHex        string   `json:"versionHex"`
	MerkleRoot        string   `json:"merkleroot"`
	Time              int64    `json:"time"`
	MedianTime        int64    `json:"mediantime"`
	Nonce             uint64   `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	ChainWork         string   `json:"chainwork"`
	PreviousBlockHash string   `json:"previousblockhash"`
	NextBlockHash     string   `json:"nextblockhash"`
	Tx                []string `json:"tx"`
}

type RawTransaction struct {
	TxID     string `json:"txid"`
	Hash     string `json:"hash"`
	Version  int32  `json:"version"`
	LockTime uint32 `json:"locktime"`

	Vin []Vin `json:"vin"`

	Vout []any `json:"vout"`
}

type Vin struct {
	Coinbase string `json:"coinbase"`
}
