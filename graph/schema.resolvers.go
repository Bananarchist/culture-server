package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Bananarchist/culture-server/graph/generated"
	"github.com/Bananarchist/culture-server/graph/model"
)

func (r *mutationResolver) CreateTaxonomy(ctx context.Context, input model.NewTaxonomy) (*model.Taxonomy, error) {
	t, err := r.InsertTaxonomy(input)
	return &t, err
}

func (r *mutationResolver) UpdateTaxonomy(ctx context.Context, id string, genus *string, species *string, subspecies *string) (*model.Taxonomy, error) {
	r.UpdateTaxonomyFields(id, genus, species, subspecies)
	t, err := r.GetTaxonomyByID(id)
	return &t, err
}

func (r *mutationResolver) CreateSpecimen(ctx context.Context, input model.NewSpecimen) (*model.Specimen, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateSpecimen(ctx context.Context, id string, taxonomyID *string, parentSpecimenID *string, nickname *string) (*model.Specimen, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFormula(ctx context.Context, input model.NewFormula) (*model.Formula, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFormula(ctx context.Context, id string, description *string, nickname *string) (*model.Formula, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateSubstrate(ctx context.Context, input model.NewSubstrate) (*model.Substrate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateSubstrate(ctx context.Context, id string, formulaID *string, quantity *float64, unit *model.VolumetricUnit) (*model.Substrate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateJar(ctx context.Context, input model.NewJar) (*model.Jar, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateJar(ctx context.Context, id string, jarType *model.ContainerType, description *string, volume *float64, unit *model.VolumetricUnit) (*model.Jar, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateCulture(ctx context.Context, input model.NewCulture) (*model.Culture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateCulture(ctx context.Context, id string, jarID *string, specimenID *string, substrateID *string, cultured *string) (*model.Culture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CloneCulture(ctx context.Context, id string, jarID *string, specimenID *string, substrateID *string, cultured *string) (*model.Culture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Taxonomy(ctx context.Context, taxonomyID string) (*model.Taxonomy, error) {
	t, err := r.GetTaxonomyByID(taxonomyID)
	return &t, err
}

func (r *queryResolver) Taxonomies(ctx context.Context) ([]*model.Taxonomy, error) {
	rows, err := r.Conn.Query("select taxonomy_id, species, genus, subspecies from taxonomy")
	if err != nil {
		panic("Could not query taxonomy (species)")
		//should actually be [return nil, err] ?
	}
	defer rows.Close()
	taxas := make([]*model.Taxonomy, 0, 10)
	for rows.Next() {
		tax := model.Taxonomy{}
		err = rows.Scan(&tax.TaxonomyID, &tax.Species, &tax.Genus, &tax.Subspecies)
		if err != nil {
			panic("Some error in scanning")
		}
		images, _ := r.TaxonomyImages(ctx, tax.TaxonomyID)
		tax.Images = images
		taxas = append(taxas, &tax)
	}
	return taxas, nil
}

func (r *queryResolver) Specimen(ctx context.Context, specimenID string) (*model.Specimen, error) {
	row := r.Conn.QueryRow("select specimen_id, taxonomy_id, parent_specimen_id, nickname from specimen where specimen_id=$1", specimenID)
	s := model.Specimen{}
	var taxonomyID string
	var parentSpecimenID sql.NullString
	switch err := row.Scan(&s.SpecimenID, &taxonomyID, &parentSpecimenID, &s.Nickname); err {
	case sql.ErrNoRows:
		return nil, err
	case nil:
		break
	default:
		panic(err)
	}
	species, _ := r.Taxonomy(ctx, taxonomyID)
	s.Taxonomy = species
	//panicIfNotNil(err)
	if parentSpecimenID.Valid {
		parent, _ := r.Specimen(ctx, parentSpecimenID.String)
		s.ParentSpecimen = parent
	}
	return &s, nil
}

func (r *queryResolver) Specimens(ctx context.Context) ([]*model.Specimen, error) {
	rows, err := r.Conn.Query("select specimen_id, taxonomy_id, parent_specimen_id, nickname from specimen")
	if err != nil {
		panic("Could not query specimen")
	}

	defer rows.Close()
	specimen := make([]*model.Specimen, 0, 10)
	for rows.Next() {
		speciman := model.Specimen{}
		var taxonomyID string
		var parentSpecimenID sql.NullString
		err = rows.Scan(&speciman.SpecimenID, &taxonomyID, &parentSpecimenID, &speciman.Nickname)

		if err != nil {
			panic("Some error in scanning specimen")
		}
		species, _ := r.Taxonomy(ctx, taxonomyID)
		speciman.Taxonomy = species
		//panicIfNotNil(err)
		if parentSpecimenID.Valid {
			parent, _ := r.Specimen(ctx, parentSpecimenID.String)
			speciman.ParentSpecimen = parent
		}
		specimen = append(specimen, &speciman)
	}
	return specimen, nil
}

func (r *queryResolver) Images(ctx context.Context) ([]*model.TaxonomyImage, error) {
	rows, err := r.Conn.Query("select taxonomy_image_id, taxonomy_id, filename from taxonomy_image")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	taxas := make([]*model.TaxonomyImage, 0, 10)
	for rows.Next() {
		img := model.TaxonomyImage{}
		err = rows.Scan(&img.TaxonomyImageID, &img.TaxonomyID, &img.Filename)
		if err != nil {
			panic("Some error in scanning images")
		}
		taxas = append(taxas, &img)
	}
	return taxas, nil
}

func (r *queryResolver) TaxonomyImages(ctx context.Context, taxonomyID string) ([]*model.TaxonomyImage, error) {
	rows, err := r.Conn.Query("select taxonomy_image_id, taxonomy_id, filename from taxonomy_image where taxonomy_id=$1", taxonomyID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	taxas := make([]*model.TaxonomyImage, 0, 10)
	for rows.Next() {
		img := model.TaxonomyImage{}
		err = rows.Scan(&img.TaxonomyImageID, &img.TaxonomyID, &img.Filename)
		if err != nil {
			panic("Some error in scanning images")
		}
		taxas = append(taxas, &img)
	}
	return taxas, nil
}

func (r *queryResolver) Jar(ctx context.Context, jarID string) (*model.Jar, error) {
	row := r.Conn.QueryRow(`
		select j.jar_id, j.description, j.volume, u.name as unit, t.name as jar_type 
		from jar j 
		left join jtype t on t.jtype_id=j.jtype_id 
		left join vunit u on u.vunit_id=j.vunit_id
		where j.jar_id=$1`, jarID)
	j := model.Jar{}
	switch err := row.Scan(&j.JarID, &j.Description, &j.Volume, &j.Unit, &j.JarType); err {
	case sql.ErrNoRows:
		//panic("No row returned for id")
		return nil, err
	case nil:
		break
	default:
		panic(err)
	}
	return &j, nil
}

func (r *queryResolver) Jars(ctx context.Context) ([]*model.Jar, error) {
	rows, err := r.Conn.Query(`
		select j.jar_id, j.description, j.volume, u.name as unit, t.name as jar_type 
		from jar j 
		left join jtype t on t.jtype_id=j.jtype_id 
		left join vunit u on u.vunit_id=j.vunit_id`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	jars := make([]*model.Jar, 0, 10)
	for rows.Next() {
		j := model.Jar{}
		err = rows.Scan(&j.JarID, &j.Description, &j.Volume, &j.Unit, &j.JarType)
		if err != nil {
			return nil, err
		}
		jars = append(jars, &j)
	}
	return jars, nil
}

func (r *queryResolver) Formula(ctx context.Context, formulaID string) (*model.Formula, error) {
	row := r.Conn.QueryRow(`
		select formula_id, description, nickname 
		from formula f
		where formula_id=$1`, formulaID)
	f := model.Formula{}
	switch err := row.Scan(&f.FormulaID, &f.Description, &f.Nickname); err {
	case sql.ErrNoRows:
		//panic("No row returned for id")
		return nil, err
	case nil:
		break
	default:
		panic(err)
	}
	return &f, nil
}

func (r *queryResolver) Formulae(ctx context.Context) ([]*model.Formula, error) {
	rows, err := r.Conn.Query(`
		select j.jar_id, j.description, j.volume, u.name as unit, t.name as jar_type 
		from jar j 
		left join jtype t on t.jtype_id=j.jtype_id 
		left join vunit u on u.vunit_id=j.vunit_id`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	formulae := make([]*model.Formula, 0, 10)
	for rows.Next() {
		f := model.Formula{}
		err = rows.Scan(&f.FormulaID, &f.Description, &f.Nickname)
		if err != nil {
			return nil, err
		}
		formulae = append(formulae, &f)
	}
	return formulae, nil
}

func (r *queryResolver) Substrate(ctx context.Context, substrateID string) (*model.Substrate, error) {
	row := r.Conn.QueryRow(`
		select s.substrate_id, s.formula_id, s.quantity, u.name as unit
		from substrate s 
		left join vunit u on u.vunit_id=s.vunit_id
		where s.substrate_id=$1`, substrateID)
	s := model.Substrate{}
	var formulaID string
	switch err := row.Scan(&s.SubstrateID, &formulaID, &s.Quantity, &s.Unit); err {
	case sql.ErrNoRows:
		//panic("No row returned for id")
		return nil, err
	case nil:
		break
	default:
		panic(err)
	}
	formula, _ := r.Formula(ctx, formulaID)
	s.Formula = formula
	return &s, nil
}

func (r *queryResolver) Substrates(ctx context.Context) ([]*model.Substrate, error) {
	rows, err := r.Conn.Query(`
	select s.substrate_id, s.formula_id, s.quantity, u.name as unit
	from substrate s 
	left join vunit u on u.vunit_id=s.vunit_id`)
	if err != nil {
		panic("Could not query specimen")
	}
	defer rows.Close()
	substrates := make([]*model.Substrate, 0, 10)
	for rows.Next() {
		s := model.Substrate{}
		var formulaID string
		err = rows.Scan(&s.SubstrateID, &formulaID, &s.Quantity, &s.Unit)
		if err != nil {
			panic("Some error in scanning specimen")
		}
		formula, _ := r.Formula(ctx, formulaID)
		s.Formula = formula
		substrates = append(substrates, &s)
	}
	return substrates, nil
}

func (r *queryResolver) Culture(ctx context.Context, cultureID string) (*model.Culture, error) {
	row := r.Conn.QueryRow(`
		select 
			c.cultured, c.culture_id,
			c.jar_id, c.specimen_id,
			c.substrate_id
		from culture c
		where c.culture_id=$1`, cultureID)
	c := model.Culture{}
	var specimenID, jarID, substrateID string
	switch err := row.Scan(&c.Cultured, &c.CultureID, &jarID, &specimenID, &substrateID); err {
	case sql.ErrNoRows:
		return nil, err
	case nil:
		break
	default:
		return nil, err
	}
	specimen, _ := r.Specimen(ctx, specimenID)
	c.Specimen = specimen
	jar, _ := r.Jar(ctx, jarID)
	c.Jar = jar
	substrate, _ := r.Substrate(ctx, substrateID)
	c.Substrate = substrate
	return &c, nil
}

func (r *queryResolver) Cultures(ctx context.Context) ([]*model.Culture, error) {
	rows, err := r.Conn.Query(`
	select 
			c.cultured, c.culture_id,
			c.jar_id, c.specimen_id,
			c.substrate_id
		from culture c`)
	if err != nil {
		panic("Could not query cultures")
	}
	defer rows.Close()
	cultures := make([]*model.Culture, 0, 10)
	for rows.Next() {
		c := model.Culture{}
		var specimenID, jarID, substrateID string
		err := rows.Scan(&c.Cultured, &c.CultureID, &jarID, &specimenID, &substrateID)
		if err != nil {
			panic("Some error in scanning culture")
		}
		specimen, _ := r.Specimen(ctx, specimenID)
		c.Specimen = specimen
		jar, _ := r.Jar(ctx, jarID)
		c.Jar = jar
		substrate, _ := r.Substrate(ctx, substrateID)
		c.Substrate = substrate
		cultures = append(cultures, &c)
	}
	return cultures, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
