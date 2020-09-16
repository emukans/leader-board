package model

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
