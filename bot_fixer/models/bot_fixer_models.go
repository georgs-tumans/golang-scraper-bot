package bot_fixer_models

type TelegramBotAPIResponse struct {
	Ok          bool   `json:"ok"`
	Result      bool   `json:"result"`
	Description string `json:"description"`
}
