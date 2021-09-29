package gorpUtil

import (
	"fmt"
)

const (
	Join_Inner = "inner"
	Join_Left  = "left"
	Join_Right = "right"
)

type Relations []*Relation

func (r Relations) Conditions(s Model, tables, joins *[]string, cs *[]*Condition) error {
	for _, ir := range r {
		if err := ir.Conditions(s, tables, joins, cs); err != nil {
			return err
		}
	}
	return nil
}

type Relation struct {
	model   Model
	relCond *Condition
	edge    string
	join    *string
	rels    []*Relation

	_subJoin []string
}

func (r *Relation) TableName() string {
	tn := r.model.TableName()
	if r._subJoin != nil {
		for _, j := range r._subJoin {
			tn = fmt.Sprintf("%v %v", tn, j)
		}
	}
	return tn
}

func (r *Relation) Join(ds ...string) *Relation {
	direction := Join_Inner
	if len(ds) > 0 {
		direction = ds[0]
	}
	r.join = &direction
	return r
}

func (r *Relation) Conditions(s Model, tables, joins *[]string, cs *[]*Condition) error {
	if r.rels != nil {
		r._subJoin = []string{}
		if err := Relations(r.rels).Conditions(r.model, tables, &r._subJoin, cs); err != nil {
			return err
		}
	}
	f := r.conditions
	t := tables
	if r.join != nil {
		f = r.jConditions
		t = joins
	}
	if err := f(s, t, cs); err != nil {
		return err
	}
	return nil
}

type (
	conHandler func(ft string) error
	relHandler func(sk, fk *Field, ft string) error
	mulHandler func(sk, smk, fmk, fk *Field, mt, ft string) error
)

func (r *Relation) _contition(s Model, pk *Field, cf conHandler, rf relHandler, mf mulHandler) error {
	if r.relCond != nil {
		if err := cf(r.TableName()); err != nil {
			return err
		}
	} else if r.edge != "" {
		if fk, ok := r.model.Relation(r.edge); ok {
			if err := rf(pk, fk, r.TableName()); err != nil {
				return err
			}
		} else if fk, ok := s.Relation(r.edge); ok && pk.Table() == fk.Table() {
			mpk, _ := s.PK()
			if err := rf(fk, mpk, r.TableName()); err != nil {
				return err
			}
		} else if rt, sk, fk, ok := r.model.Multiple(r.edge); ok {
			mpk, _ := s.PK()
			if err := mf(pk, sk, fk, mpk, rt, r.TableName()); err != nil {
				return err
			}
		} else if rt, sk, fk, ok := s.Multiple(r.edge); ok {
			mpk, _ := s.PK()
			if err := mf(mpk, sk, fk, pk, rt, r.TableName()); err != nil {
				return err
			}
		} else {
			return ErrRelNotFound
		}
	} else {
		return ErrRelNotFound
	}
	return nil
}

type (
	joinUtils []*joinUtil
	joinUtil  struct {
		direction string
		tableName string
		c         *Condition
		sub       []*joinUtil
	}
)

func ju(d, tn string, c *Condition, sub ...*joinUtil) *joinUtil {
	return &joinUtil{
		direction: d,
		tableName: tn,
		c:         c,
		sub:       sub,
	}
}

func (j *joinUtil) String() (string, error) {
	inner := ""
	if len(j.sub) > 0 {
		sub, err := joinUtils(j.sub).String()
		if err != nil {
			return "", err
		}
		inner = fmt.Sprintf(" %v", sub)
	}
	cs, err := j.c.Relation().String(0)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v join %v%v on %v", j.direction, j.tableName, inner, cs), nil
}

func (js joinUtils) String() (string, error) {
	s := ""
	for _, j := range js {
		if s != "" {
			s = fmt.Sprintf(" %v", s)
		}
		_s, err := j.String()
		if err != nil {
			return "", err
		}
		s = fmt.Sprintf("%v%v", s, _s)
	}
	return s, nil
}

func (r *Relation) jConditions(s Model, joins *[]string, cs *[]*Condition) error {
	rpk, v := r.model.PK()
	if i, ok := v.(int); ok && i > 0 {
		*cs = append(*cs, rpk.EQ(v))
	}
	return r._contition(s, rpk,
		func(ft string) error {
			s, err := ju(
				*r.join,
				ft,
				r.relCond,
			).String()
			if err != nil {
				return err
			}
			*joins = append(*joins, s)
			return nil
		},
		func(sk, fk *Field, ft string) error {
			s, err := ju(
				*r.join,
				ft,
				sk.EQ(fk),
			).String()
			if err != nil {
				return err
			}
			*joins = append(*joins, s)
			return nil
		},
		func(sk, smk, fmk, fk *Field, mt, ft string) error {
			s, err := ju(
				*r.join,
				mt,
				fmk.EQ(fk),
				ju(
					*r.join,
					ft,
					sk.EQ(smk),
				),
			).String()
			if err != nil {
				return err
			}
			*joins = append(*joins, s)
			return nil
		},
	)
}

func (r *Relation) conditions(s Model, tables *[]string, cs *[]*Condition) error {
	*tables = append(*tables, r.TableName())
	rpk, v := r.model.PK()
	if i, ok := v.(int); ok && i > 0 {
		*cs = append(*cs, rpk.EQ(v))
	}
	return r._contition(s, rpk,
		func(ft string) error {
			*cs = append(*cs, r.relCond.Relation())
			return nil
		},
		func(sk, fk *Field, ft string) error {
			*cs = append(*cs, sk.EQ(fk).Relation())
			return nil
		},
		func(sk, smk, fmk, fk *Field, mt, ft string) error {
			*tables = append(*tables, mt)
			*cs = append(*cs,
				sk.EQ(smk).Relation(),
				fmk.EQ(fk).Relation(),
			)
			return nil
		},
	)
}

func R(m Model, edge string, rs ...*Relation) *Relation {
	return &Relation{
		model: m,
		edge:  edge,
		rels:  rs,
	}
}

func RWith(m Model, c *Condition, rs ...*Relation) *Relation {
	return &Relation{
		model:   m,
		relCond: c,
		rels:    rs,
	}
}
