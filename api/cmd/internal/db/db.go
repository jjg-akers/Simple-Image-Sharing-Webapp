package db

type Inserter interface {
	Insert()
}

type Retriever interface {
	Retrieve()
}

type InserterRetriver interface {
	Inserter
	Retriever
}

type SQLParams struct {
	Params []interface{}
}

func NewSQLParams(opt ParamOption) (*SQLParams, error) {
	p := &SQLParams{}
	err := opt(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

type ParamOption func(s *SQLParams) error

func StringParam(strs []string) ParamOption {
	return func(s *SQLParams) error {
		p := make([]interface{}, len(strs))
		for i, v := range strs {
			p[i] = v
		}
		s.Params = p

		return nil
	}
}
