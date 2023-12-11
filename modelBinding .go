func StructToMap(u1 any) (m map[string]interface{}) {
	b, _ := json.Marshal(&u1)
	_ = json.Unmarshal(b, &m)
	return
}
