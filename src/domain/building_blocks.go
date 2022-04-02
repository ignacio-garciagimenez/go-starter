package domain

type ValueObject interface{}
type DomainEvent interface{}

type Repository[K comparable, E Entity[K]] interface {
	FindByID(key K) (E, error)
	Save(entity E) error
}

type Entity[K comparable] interface {
	GetID() K
	addDomainEvent(DomainEvent)
	GetDomainEvents() []DomainEvent
	ClearDomainEvents()
}

type baseEntity[K comparable] struct {
	id           K
	domainEvents []DomainEvent
}

func (e baseEntity[K]) GetID() K {
	return e.id
}

func (e baseEntity[K]) EqualsTo(other Entity[K]) bool {
	return e.id == other.GetID()
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
