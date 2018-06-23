package tox

type callHookMethods struct {
	ConferenceJoin     func(friendNumber uint32, groupNumber uint32, cookie string)
	ConferenceDelete   func(groupNumber uint32)
	ConferenceNew      func(groupNumber uint32)
	ConferenceSetTitle func(groupNumber uint32, title string)
}

// include av group
func (t *Tox) HookConferenceJoin(fn func(friendNumber uint32, groupNumber uint32, cookie string)) {
	t.hooks.ConferenceJoin = fn
}

func (t *Tox) HookConferenceDelete(fn func(groupNumber uint32)) {
	t.hooks.ConferenceDelete = fn
}

func (t *Tox) HookConferenceNew(fn func(groupNumber uint32)) {
	t.hooks.ConferenceDelete = fn
}

func (t *Tox) HookConferenceSetTitle(fn func(groupNumber uint32, title string)) {
	t.hooks.ConferenceSetTitle = fn
}
