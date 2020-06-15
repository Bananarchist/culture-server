package graph

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Bananarchist/culture-server/graph/model"
	"github.com/Bananarchist/culture-server/graph/sqlutil"
)

// Resolver for sharing the db luv
type Resolver struct {
	Conn *sql.DB
}

const taxonomyFieldString = "taxonomy_id, genus, species, subspecies"

func taxonomyFields(t *model.Taxonomy) []interface{} {
	var Subspecies sql.NullString
	return []interface{}{&t.TaxonomyID, &t.Genus, &t.Species, Subspecies}
}

func taxonomyFromFields(t []interface{}) model.Taxonomy {
	return model.Taxonomy{
		TaxonomyID: t[0].(string),
		Genus:      t[1].(string),
		Species:    t[2].(string),
		Subspecies: t[3].(*string)}
}

//GetTaxonomyByID returns a *model.Taxonomy if said id exists
func (r *Resolver) GetTaxonomyByID(id string) (t model.Taxonomy, err error) {
	query := r.GetObjectByIDQuery("taxonomy", "taxonomy_id", taxonomyFieldString)
	row := r.Conn.QueryRow(query, id)
	var Subspecies sql.NullString
	err = row.Scan(&t.TaxonomyID, &t.Genus, &t.Species, &Subspecies)
	if Subspecies.Valid {
		t.Subspecies = &Subspecies.String
	}
	return
}

//GetObjectByIDQuery takes a table, idfieldname and a list of fields selected and returns a select statement string
func (r *Resolver) GetObjectByIDQuery(table, idFieldName, fields string) string {
	return fmt.Sprintf("select %s from %s where %s=$1", fields, table, idFieldName)
}

//QueryReturning ads
func (r *Resolver) QueryReturning(query, returning string) string {
	return query + " returning " + returning
}

//QueryReturningColumns sad
func (r *Resolver) QueryReturningColumns(query string, columns []string) string {
	return r.QueryReturning(query, strings.Join(columns, ", "))
}

//CreateObjectInsertQuery s da--
func (r *Resolver) CreateObjectInsertQuery(table string, columns []string, objectCount int) string {
	valString := sqlutil.GetRowsStringForColumnsList(1, 1, columns)
	colString := strings.Join(columns, ", ")
	return fmt.Sprintf("insert into %s (%s) values %s", table, colString, valString)
}

//HandleInsertErrors panics for norow, but otherwise just passes error/nil on
func (r *Resolver) HandleInsertErrors(err error) error {
	switch err {
	case sql.ErrNoRows:
		panic(err)
	case nil:
		return nil
	default:
		return err
	}
}

//InsertTaxonomy returns a *model.Taxonomy if it can create one
func (r *Resolver) InsertTaxonomy(input model.NewTaxonomy) (t model.Taxonomy, err error) {
	query := r.QueryReturning(r.CreateObjectInsertQuery("taxonomy", []string{"genus", "species", "subspecies"}, 1), taxonomyFieldString)
	var Subspecies sql.NullString
	err = r.Conn.QueryRow(query, input.Genus, input.Species, input.Subspecies).Scan(&t.TaxonomyID, &t.Species, &t.Genus, &Subspecies)
	if Subspecies.Valid {
		t.Subspecies = &Subspecies.String
	}
	err = r.HandleInsertErrors(err)
	return
}

//UpdateTaxonomyFields returns a *model.Taxonomy if said id exists
func (r *Resolver) UpdateTaxonomyFields(id string, genus, species, subspecies *string) (err error) {
	return r.SetObjectFields(id, "taxonomy_id", "taxonomy", map[interface{}]interface{}{"genus": genus, "species": species, "subspecies": subspecies})
}

//SetObjectFields afdsds
func (r *Resolver) SetObjectFields(id, idFieldName, table string, val2testmap map[interface{}]interface{}) (err error) {
	argList := make([]interface{}, len(val2testmap))
	idx := 0
	for _, t := range val2testmap {
		argList[idx] = t
		idx++
	}
	if len(argList) < 1 {
		panic("NO ARGS FOUND")
	}
	args := sqlutil.GetArgsArr(argList...)
	args = append(args, id)
	columns := sqlutil.GetColumnsArr(val2testmap)
	query := fmt.Sprintf("update %s set %s where %s=$%d", table, sqlutil.SQLUpdateArgsString(columns), idFieldName, len(args))
	_, err = r.Conn.Exec(query, args...)
	return
}
