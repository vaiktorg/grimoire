package grimoire

type StateConfig struct {
	allowable map[Trigger]State
	onEntry   func(State) error
	onExit    func(State) error
	onFail    func(error)
	onSuccess func()
	fallback  State
}

func (c *StateConfig) OnEntry(onEntry func(State) error) *StateConfig {
	c.onEntry = onEntry
	return c
}

func (c *StateConfig) OnExit(onExit func(State) error) *StateConfig {
	c.onExit = onExit
	return c
}
func (c *StateConfig) FallbackTo(fallbackState State) *StateConfig {
	c.fallback = fallbackState
	return c
}
func (c *StateConfig) Allow(trigger Trigger, state State) *StateConfig {
	if _, ok := c.allowable[trigger]; !ok {
		c.allowable[trigger] = state
	}
	return c
}
