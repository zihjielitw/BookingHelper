package main

type TomlConfig struct {
	Account string
	Password string
	Res string
	City string
	People string
	BookingDate string
	MealTime string
	Store string
	Time string
	Vegetable string
	ChildChair string
}

type LoginResult struct {
	Rcrm struct {
		Rc string `json:"RC"`
		Rm string `json:"RM"`
	} `json:"rcrm"`
	Results struct {
		UserAccessToken string      `json:"user_access_token"`
		MachineID       interface{} `json:"machine_id"`
		StoreName       interface{} `json:"store_name"`
		VipClass        string      `json:"vip_class"`
		VipClassImage   interface{} `json:"vip_class_image"`
		VipType         string      `json:"vip_type"`
		Remarks         interface{} `json:"remarks"`
	} `json:"results"`
}


type PossibleSeat struct {

		Date     string `json:"date"`
		City     int    `json:"city"`
		MealTime string `json:"mealTime"`
		Content  []struct {
			Store string `json:"store"`
			Data  struct {
				Seat         int    `json:"seat"`
				MealTimeWord string `json:"mealTimeWord"`
			} `json:"data"`
			Calendar interface{} `json:"calendar"`
			Type     string `json:"type"`
		} `json:"content"`

}


type BookingResult struct {
	State   string `json:"state"`
	OrderNo string `json:"order_no"`
}