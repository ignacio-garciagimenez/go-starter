package domain

type ValueObject interface {
	EqualsTo(ValueObject) bool
}

type DomainEvent interface{}

type Repository[K comparable, E Entity[K]] interface {
	FindByID(K) (E, error)
	Save(E) error
}

type Entity[K comparable] interface {
	GetID() K
	addDomainEvent(DomainEvent)
	GetDomainEvents() []DomainEvent
	ClearDomainEvents()
	EqualsTo(Entity[K]) bool
}

type baseEntity[K comparable] struct {
	id           K
	domainEvents []DomainEvent
}

func (e baseEntity[K]) GetID() K {
	return e.id
}

func (e *baseEntity[K]) addDomainEvent(event DomainEvent) {
	e.domainEvents = append(e.domainEvents, event)
}

func (e baseEntity[K]) GetDomainEvents() []DomainEvent {
	var output []DomainEvent
	for _, domainEvent := range e.domainEvents {
		output = append(output, domainEvent)
	}
	return output
}

func (e *baseEntity[K]) ClearDomainEvents() {
	e.domainEvents = []DomainEvent{}
}
