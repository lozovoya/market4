package repository

type Requests []struct {
	Request string `yaml:"request"`
}

type TestData struct {
	Conf struct {
		Setup struct {
			Requests Requests
		}
		Teardown struct {
			Requests Requests
		}
	}
}
