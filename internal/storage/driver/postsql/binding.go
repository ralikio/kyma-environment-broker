package postsql

import (
	"fmt"

	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/dberr"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/dbmodel"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/postsql"
	log "github.com/sirupsen/logrus"
)

type Binding struct {
	postsql.Factory
	cipher     Cipher
}

func NewBinding(sess postsql.Factory, cipher Cipher) *Binding {
	return &Binding{
		Factory:    sess,
		cipher:     cipher,
	}
}

// TODO: Wrap retries in single method WithRetries
func (s *Binding) Get(bindingId string) (*internal.Binding, error) {
	sess := s.NewReadSession()
	bindingDTO := dbmodel.BindingDTO{}
	bindingDTO, lastErr := sess.GetBindingByID(bindingId)
	if lastErr != nil {
		if dberr.IsNotFound(lastErr) {
			return nil, dberr.NotFound("Binding with id %s not exist", bindingId)
		}
		log.Errorf("while getting instanceDTO by ID %s: %v", bindingId, lastErr)
		return nil, lastErr
	}
	binding, err := s.toBinding(bindingDTO)
	if err != nil {
		return nil, err
	}

	return &binding, nil
}


// keep
func (s *Binding) Insert(binding *internal.Binding) error {
	_, err := s.Get(binding.ID)
	if err == nil {
		return dberr.AlreadyExists("instance with id %s already exist", binding.ID)
	}

	dto, err := s.toBindingDTO(binding)
	if err != nil {
		return err
	}

	sess := s.NewWriteSession()
	err = sess.InsertBinding(dto)
	if err != nil {
		return fmt.Errorf("while saving binding with ID %s: %w", binding.ID, err)	
	}

	return nil
}

func (s *Binding) Update(binding *internal.Binding) (*internal.Binding, error) {
	sess := s.NewWriteSession()
	dto, err := s.toBindingDTO(binding)
	if err != nil {
		return nil, err
	}
	var lastErr dberr.Error
		lastErr = sess.UpdateBinding(dto)

		switch {
		case dberr.IsNotFound(lastErr):
			_, lastErr = s.NewReadSession().GetBindingByID(binding.ID)
			if dberr.IsNotFound(lastErr) {
				return nil, dberr.NotFound("Binding with id %s not exist", binding.ID)
			}
			if lastErr != nil {
				return nil, fmt.Errorf("while getting Operation: %w", lastErr)
			}

			return nil, lastErr
		case lastErr != nil:
			return nil, fmt.Errorf("while updating instance ID %s: %w", binding.ID, lastErr)
		}
	if err != nil {
		return nil, lastErr
	}
	binding.Version = binding.Version + 1
	return binding, nil
}



// leave
func (s *Binding) Delete(instanceID string) error {
	sess := s.NewWriteSession()
	return sess.DeleteBinding(instanceID)
}

func (s *Binding) List(filter dbmodel.BindingFilter) ([]internal.Binding, int, int, error) {
	dtos, count, totalCount, err := s.NewReadSession().ListBindings(filter)
	if err != nil {
		return []internal.Binding{}, 0, 0, err
	}
	var bindings []internal.Binding
	for _, dto := range dtos {
		instance, err := s.toBinding(dto)
		if err != nil {
			return []internal.Binding{}, 0, 0, err
		}

		bindings = append(bindings, instance)
	}
	return bindings, count, totalCount, err
}



func (s *Binding) toBindingDTO(binding *internal.Binding) (dbmodel.BindingDTO, error) {
	encrypted, err := s.cipher.Encrypt([]byte(binding.Kubeconfig))
	if err != nil {
		return dbmodel.BindingDTO{}, fmt.Errorf("while encrypting kubeconfig: %w", err)
	}
	
	return dbmodel.BindingDTO{
		Kubeconfig: 				string(encrypted),
		ID: 					   binding.ID,
		RuntimeID:                   binding.RuntimeID,
		CreatedAt:                   binding.CreatedAt,
		UpdatedAt:                   binding.UpdatedAt,
		DeletedAt:                   binding.DeletedAt,
		ExpiredAt:                   binding.ExpiredAt,
		Version:                     binding.Version,
	}, nil
}

func (s *Binding) toBinding(dto dbmodel.BindingDTO) (internal.Binding, error) {
	decrypted, err := s.cipher.Decrypt([]byte(dto.Kubeconfig))
	if err != nil {
		return internal.Binding{}, fmt.Errorf("while decrypting kubeconfig: %w", err)	
	}

	return internal.Binding{
		Kubeconfig: string(decrypted),
		ID: 					   dto.ID,
		RuntimeID:                   dto.RuntimeID,
		CreatedAt:                   dto.CreatedAt,
		UpdatedAt:                   dto.UpdatedAt,
		DeletedAt:                   dto.DeletedAt,
		ExpiredAt:                   dto.ExpiredAt,
		Version:                     dto.Version,
	}, nil
}