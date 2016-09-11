package triton

type StepCreateKey struct {
	KeyName     string
	KeyMaterial string
}

func (s *StepCreateKey) Run(state TritonStateBag) StepAction {
	client := state.Client().Keys()

	key, err := client.CreateKey(&CreateKeyInput{
		Name: s.KeyName,
		Key:  s.KeyMaterial,
	})
	if err != nil {
		state.AppendError(err)
		return Halt
	}

	state.Put("key", key)
	return Continue
}

func (s *StepCreateKey) Cleanup(state TritonStateBag) {
	client := state.Client().Keys()

	if err := client.DeleteKey(&DeleteKeyInput{
		KeyName: s.KeyName,
	}); err != nil {
		if IsResourceNotFound(err) {
			return
		}
		state.AppendError(err)
	}
}

type StepGetKey struct {
	StateBagKey     string
	KeyName         string
	ExpectedErrorKey string
}

func (s *StepGetKey) Run(state TritonStateBag) StepAction {
	client := state.Client().Keys()

	key, err := client.GetKey(&GetKeyInput{
		KeyName: s.KeyName,
	})
	if err != nil {
		if s.ExpectedErrorKey != "" {
			state.Put(s.ExpectedErrorKey, err)
			return Continue
		}

		state.AppendError(err)
		return Halt
	}

	state.Put(s.StateBagKey, key)
	return Continue
}

func (s *StepGetKey) Cleanup(state TritonStateBag) {
	return
}

type StepDeleteKey struct {
	KeyName string
}

func (s *StepDeleteKey) Run(state TritonStateBag) StepAction {
	client := state.Client().Keys()

	err := client.DeleteKey(&DeleteKeyInput{
		KeyName: s.KeyName,
	})
	if err != nil {
		if IsResourceNotFound(err) {
			return Continue
		}
		state.AppendError(err)
		return Halt
	}

	return Continue
}

func (s *StepDeleteKey) Cleanup(state TritonStateBag) {
	return
}
