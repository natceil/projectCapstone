package apiSplunk

type Row struct{
	Preview bool `json:"preview"`
	Offset int `json:"offset"`
	Result map[string]Value `json:"result"`
	LastRow bool `json:"lastrow"`
}

type Value struct{

}