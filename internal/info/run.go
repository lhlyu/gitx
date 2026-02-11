package info

func Run() error {
	service := NewService()

	data, err := service.Collect()
	if err != nil {
		return err
	}

	Print(data)
	return nil
}
