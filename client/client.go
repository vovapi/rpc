package client

import (
	// Vendor
	"github.com/pkg/errors"

	// RPC
	"github.com/asuleymanov/rpc"
	"github.com/asuleymanov/rpc/transactions"
	"github.com/asuleymanov/rpc/transports/websocket"
	"github.com/asuleymanov/rpc/types"
)

const fdt = `"20060102t150405"`

var Key_List = make(map[string]Keys)

type Keys struct {
	PKey string
	AKey string
	OKey string
	MKey string
}

type Client struct {
	Rpc   *rpc.Client
	Chain *transactions.Chain
}

type BResp struct {
	ID       string
	BlockNum uint32
	TrxNum   uint32
	Expired  bool
}

func initclient(url []string) *rpc.Client {
	// Инициализация Websocket
	t, err := websocket.NewTransport(url)
	if err != nil {
		panic(errors.Wrapf(err, "Error Websocket: "))
	}

	// Инициализация RPC клиента
	client, err := rpc.NewClient(t)
	if err != nil {
		panic(errors.Wrapf(err, "Error RPC: "))
	}
	//defer client.Close()
	return client
}

func initChainId(str string) *transactions.Chain {
	var ChainId transactions.Chain
	// Определяем ChainId
	switch str {
	case "steem":
		ChainId = *transactions.SteemChain
	}
	return &ChainId
}

func NewApi(url []string, chain string) *Client {
	return &Client{
		Rpc:   initclient(url),
		Chain: initChainId(chain),
	}
}

func (api *Client) Send_Trx(username string, strx types.Operation) (*BResp, error) {
	// Получение необходимых параметров
	props, err := api.Rpc.Database.GetDynamicGlobalProperties()
	if err != nil {
		return nil, errors.Wrapf(err, "Error get DynamicGlobalProperties: ")
	}

	// Создание транзакции
	refBlockPrefix, err := transactions.RefBlockPrefix(props.HeadBlockID)
	if err != nil {
		return nil, err
	}
	tx := transactions.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(props.HeadBlockNumber),
		RefBlockPrefix: refBlockPrefix,
	})

	// Добавление операций в транзакцию
	tx.PushOperation(strx)

	// Получаем необходимый для подписи ключ
	privKeys := api.Signing_Keys(username, strx)

	// Подписываем транзакцию
	if err := tx.Sign(privKeys, api.Chain); err != nil {
		return nil, errors.Wrapf(err, "Error Sign: ")
	}

	// Отправка транзакции
	resp, err := api.Rpc.NetworkBroadcast.BroadcastTransactionSynchronous(tx.Transaction)

	if err != nil {
		return nil, errors.Wrapf(err, "Error BroadcastTransactionSynchronous: ")
	} else {
		var bresp BResp

		bresp.ID = resp.ID
		bresp.BlockNum = resp.BlockNum
		bresp.TrxNum = resp.TrxNum
		bresp.Expired = resp.Expired

		return &bresp, nil
	}
}

func (api *Client) Send_Arr_Trx(username string, strx []types.Operation) (*BResp, error) {
	// Получение необходимых параметров
	props, err := api.Rpc.Database.GetDynamicGlobalProperties()
	if err != nil {
		return nil, errors.Wrapf(err, "Error get DynamicGlobalProperties: ")
	}

	// Создание транзакции
	refBlockPrefix, err := transactions.RefBlockPrefix(props.HeadBlockID)
	if err != nil {
		return nil, err
	}
	tx := transactions.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(props.HeadBlockNumber),
		RefBlockPrefix: refBlockPrefix,
	})

	// Добавление операций в транзакцию
	for _, val := range strx {
		tx.PushOperation(val)
	}

	// Получаем необходимый для подписи ключ
	privKeys := api.Signing_Keys(username, strx[0])

	// Подписываем транзакцию
	if err := tx.Sign(privKeys, api.Chain); err != nil {
		return nil, errors.Wrapf(err, "Error Sign: ")
	}

	// Отправка транзакции
	resp, err := api.Rpc.NetworkBroadcast.BroadcastTransactionSynchronous(tx.Transaction)

	if err != nil {
		return nil, errors.Wrapf(err, "Error BroadcastTransactionSynchronous: ")
	} else {
		var bresp BResp

		bresp.ID = resp.ID
		bresp.BlockNum = resp.BlockNum
		bresp.TrxNum = resp.TrxNum
		bresp.Expired = resp.Expired

		return &bresp, nil
	}
}

func (api *Client) Verify_Trx(username string, strx types.Operation) (bool, error) {
	// Получение необходимых параметров
	props, err := api.Rpc.Database.GetDynamicGlobalProperties()
	if err != nil {
		return false, errors.Wrapf(err, "Error get DynamicGlobalProperties: ")
	}

	// Создание транзакции
	refBlockPrefix, err := transactions.RefBlockPrefix(props.HeadBlockID)
	if err != nil {
		return false, err
	}
	tx := transactions.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(props.HeadBlockNumber),
		RefBlockPrefix: refBlockPrefix,
	})

	// Добавление операций в транзакцию
	tx.PushOperation(strx)

	// Получаем необходимый для подписи ключ
	privKeys := api.Signing_Keys(username, strx)

	// Подписываем транзакцию
	if err := tx.Sign(privKeys, api.Chain); err != nil {
		return false, errors.Wrapf(err, "Error Sign: ")
	}

	// Отправка транзакции
	resp, err := api.Rpc.Database.GetVerifyAuthoruty(tx.Transaction)

	if err != nil {
		return false, errors.Wrapf(err, "Error BroadcastTransactionSynchronous: ")
	} else {
		return resp, nil
	}
}
