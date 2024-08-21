package vo

import "encoding/json"

type Address struct {
	street string
	city   string
	state  string
	zip    string
}

func NewAddress(street, city, state, zip string) Address {
	return Address{
		street: street,
		city:   city,
		state:  state,
		zip:    zip,
	}
}

type tmpAddress struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

func NewAddressFromJson(data string) Address {
	var address = tmpAddress{}
	err := json.Unmarshal([]byte(data), &address)
	if err != nil {
		return Address{}
	}
	return Address{
		street: address.Street,
		city:   address.City,
		state:  address.State,
		zip:    address.Zip,
	}
}

func (a Address) GetStreet() string {
	return a.street
}

func (a Address) GetCity() string {
	return a.city
}

func (a Address) GetState() string {
	return a.state
}

func (a Address) GetZip() string {
	return a.zip
}

func (a Address) ToJson() string {
	t := tmpAddress{
		Street: a.street,
		City:   a.city,
		State:  a.state,
		Zip:    a.zip,
	}
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}
