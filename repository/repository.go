package repository

type ReceiptRepository interface {
	StoreReceipt(uuid string, receipt Receipt)
	GetPoints(uuid string) (int, bool)
}

type receiptMap map[string]Receipt

func NewReceiptRepository() ReceiptRepository {
	rc := receiptMap{}
	return rc
}

func (receiptMap receiptMap) StoreReceipt(uuid string, receipt Receipt) {
	receiptMap[uuid] = receipt
}

func (receiptMap receiptMap) GetPoints(uuid string) (int, bool) {
	receipt, ok := receiptMap[uuid]
	return receipt.Points, ok
}
