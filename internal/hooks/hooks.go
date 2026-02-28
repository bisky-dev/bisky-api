package hooks

import "context"

type Event string

const (
	EventShowCreatePre     Event = "show.create.pre"
	EventShowCreatePost    Event = "show.create.post"
	EventShowUpdatePre     Event = "show.update.pre"
	EventShowUpdatePost    Event = "show.update.post"
	EventShowDeletePre     Event = "show.delete.pre"
	EventShowDeletePost    Event = "show.delete.post"
	EventEpisodeCreatePre  Event = "episode.create.pre"
	EventEpisodeCreatePost Event = "episode.create.post"
	EventEpisodeUpdatePre  Event = "episode.update.pre"
	EventEpisodeUpdatePost Event = "episode.update.post"
	EventEpisodeDeletePre  Event = "episode.delete.pre"
	EventEpisodeDeletePost Event = "episode.delete.post"
)

type Dispatcher interface {
	DispatchPre(ctx context.Context, event Event, payload any) error
	DispatchPost(ctx context.Context, event Event, payload any)
}

type NoopDispatcher struct{}

func (NoopDispatcher) DispatchPre(context.Context, Event, any) error {
	return nil
}

func (NoopDispatcher) DispatchPost(context.Context, Event, any) {}
