package handlers

type APIHandler struct {
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
}

func (ah *APIHandler) HandleAPICommand(chatId int64, command string) error {
	// TODO: figure out if there is a more generic/common way to handle commands

	//var err error

	// Some commands may have parameters
	// commandParts := strings.Split(command, " ")
	// mainCommand := commandParts[0]
	// param := ""

	// if len(commandParts) > 1 {
	// 	param = commandParts[1]
	// }

	// log.Printf("[BondsHandler] Handling command: %s; parameters: %s", mainCommand, param)

	// switch mainCommand {

	// case "/bonds_start":
	// 	bh.handleBondsStart(chatId)

	// case "/bonds_stop":
	// 	bh.handleBondsStop(chatId)

	// case "/bonds_status":
	// 	bh.handleBondsStatus(chatId)

	// case "/bonds_set_interval":
	// 	err = bh.handleBondsSetInterval(chatId, param)
	// }

	// return err

	return nil

}
