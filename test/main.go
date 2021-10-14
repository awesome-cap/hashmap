package test

type EventConfig struct {
	Events []Event `json:"events"`
}

type Event struct {
	ID        int64     `json:"id"`
	Predicate Predicate `json:"predicate"`
	Execute   Execute   `json:"execute"`
	Relations []int64   `json:"relations"`
}

//type Execute struct {
//	Args []interface{} `json:"args"`
//	Sql  string        `json:"sql"` // update xxx set xx = ? where uid = #{currentUser}
//}

type Execute struct {
	Units []ActionUnit `json:"units"`
}

type ActionUnit struct {
	Name string `json:"name"` // addBadge
	Desc string `json:"desc"` // 增加徽章
}

type Predicate struct {
	Key        string      `json:"key"`
	Type       string      `json:"type"` // eq gt gte lt lte between and?
	Value      interface{} `json:"value"`
	Predicates []Predicate `json:"predicates"`
}
