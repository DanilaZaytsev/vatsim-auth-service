package model

type UserInfoResponse struct {
	Data struct {
		CID      string `json:"cid"`
		Personal struct {
			FirstName string `json:"name_first"`
			LastName  string `json:"name_last"`
			FullName  string `json:"name_full"`
			Email     string `json:"email"`
			Country   struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"country"` // <-- вот здесь у тебя сейчас ошибка, если стоит string
		} `json:"personal"`
		VATSIM struct {
			Rating struct {
				ID    int    `json:"id"`
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"rating"`
			PilotRating struct {
				ID    int    `json:"id"`
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"pilotrating"`
			Division struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"division"`
			Region struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"region"`
			Subdivision struct {
				ID   *string `json:"id"`   // иногда null
				Name *string `json:"name"` // иногда null
			} `json:"subdivision"`
		} `json:"vatsim"`
		OAuth struct {
			TokenValid string `json:"token_valid"`
		} `json:"oauth"`
	} `json:"data"`
}
