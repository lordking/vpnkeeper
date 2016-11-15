package vpn

var (
	std = New()
)

func Fetch() ([]Service, error) {
	return std.Fetch()
}

func Select(seq int) error {
	return std.Select(seq)
}

func Stop() error {
	return std.Stop(std.Selected)
}

func RunServ() {
	std.RunServ()
}
