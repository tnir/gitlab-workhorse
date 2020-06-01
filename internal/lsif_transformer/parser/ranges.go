package parser

import (
	"encoding/json"
	"io"
	"strconv"
)

const Definitions = "definitions"
const References = "references"

type Ranges struct {
	Entries map[Id]*Range
	DefRefs map[Id]*DefRef
	Hovers  *Hovers
}

type RawRange struct {
	Id   Id    `json:"id"`
	Data Range `json:"start"`
}

type Range struct {
	Line      int32 `json:"line"`
	Character int32 `json:"character"`
	RefId     Id
}

type RawDefRef struct {
	Property string `json:"property"`
	RefId    Id     `json:"outV"`
	RangeIds []Id   `json:"inVs"`
	DocId    Id     `json:"document"`
}

type DefRef struct {
	Line  string
	DocId Id
}

type SerializedRange struct {
	StartLine      int32           `json:"start_line"`
	StartChar      int32           `json:"start_char"`
	DefinitionPath string          `json:"definition_path,omitempty"`
	Hover          json.RawMessage `json:"hover"`
}

func NewRanges(tempDir string) (*Ranges, error) {
	hovers, err := NewHovers(tempDir)
	if err != nil {
		return nil, err
	}

	return &Ranges{
		Entries: make(map[Id]*Range),
		DefRefs: make(map[Id]*DefRef),
		Hovers:  hovers,
	}, nil
}

func (r *Ranges) Read(label string, line []byte) error {
	switch label {
	case "range":
		if err := r.addRange(line); err != nil {
			return err
		}
	case "item":
		if err := r.addItem(line); err != nil {
			return err
		}
	default:
		return r.Hovers.Read(label, line)
	}

	return nil
}

func (r *Ranges) Serialize(f io.Writer, rangeIds []Id, docs map[Id]string) error {
	encoder := json.NewEncoder(f)
	n := len(rangeIds)

	if _, err := f.Write([]byte("[")); err != nil {
		return err
	}

	for i, rangeId := range rangeIds {
		entry := r.Entries[rangeId]
		serializedRange := SerializedRange{
			StartLine:      entry.Line,
			StartChar:      entry.Character,
			DefinitionPath: r.definitionPathFor(docs, entry.RefId),
			Hover:          r.Hovers.For(entry.RefId),
		}
		if err := encoder.Encode(serializedRange); err != nil {
			return err
		}
		if i+1 < n {
			if _, err := f.Write([]byte(",")); err != nil {
				return err
			}
		}
	}

	if _, err := f.Write([]byte("]")); err != nil {
		return err
	}

	return nil
}

func (r *Ranges) Close() error {
	return r.Hovers.Close()
}

func (r *Ranges) definitionPathFor(docs map[Id]string, refId Id) string {
	defRef, ok := r.DefRefs[refId]
	if !ok {
		return ""
	}

	defPath := docs[defRef.DocId] + "#L" + defRef.Line

	return defPath
}

func (r *Ranges) addRange(line []byte) error {
	var rg RawRange
	if err := json.Unmarshal(line, &rg); err != nil {
		return err
	}

	r.Entries[rg.Id] = &rg.Data

	return nil
}

func (r *Ranges) addItem(line []byte) error {
	var defRef RawDefRef
	if err := json.Unmarshal(line, &defRef); err != nil {
		return err
	}

	if defRef.Property != Definitions && defRef.Property != References {
		return nil
	}

	for _, rangeId := range defRef.RangeIds {
		if entry, ok := r.Entries[rangeId]; ok {
			entry.RefId = defRef.RefId
		}
	}

	if defRef.Property != Definitions {
		return nil
	}

	defRange := r.Entries[defRef.RangeIds[0]]

	r.DefRefs[defRef.RefId] = &DefRef{
		Line:  strconv.Itoa(int(defRange.Line + 1)),
		DocId: defRef.DocId,
	}

	return nil
}