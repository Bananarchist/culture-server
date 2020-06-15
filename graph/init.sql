create table if not exists vunit (
	vunit_id serial not null primary key,
	name varchar(16) not null
);

create table if not exists jtype (
	jtype_id serial not null primary key,
	name varchar(16) not null
);

create table if not exists species (
	species_id serial not null primary key,
	genus varchar(32) not null,
	species varchar(32) not null,
	subspecies varchar(32)
);

create table if not exists specimen (
	specimen_id serial not null primary key,
	species_id int not null references species (species_id),
	parent_specimen_id int references specimen (specimen_id),
	nickname varchar(32)
);

create table if not exists formula (
	formula_id serial not null primary key,
	description text,
	nickname varchar(32) not null
);

create table if not exists substrate (
	substrate_id serial not null primary key,
	formula_id int not null references formula(formula_id),
	quantity float not null,
	vunit_id int not null references vunit (vunit_id)
);

create table if not exists jar (
	jar_id serial not null primary key,
	jtype_id int not null references jtype(jtype_id),
	description text,
	volume float not null,
	vunit_id int not null references vunit(vunit_id)
);

create table if not exists culture (
	culture_id serial not null primary key,
	jar_id int not null references jar (jar_id),
	specimen_id int not null references specimen (specimen_id),
	substrate_id int not null references substrate (substrate_id),
	cultured date not null default current_date
);


create table if not exists culture_ancestry (
	parent_culture_id int not null references culture (culture_id),
	child_culture_id int not null references culture (culture_id)
);

create table if not exists culture_event_type (
	culture_event_type_id serial not null primary key,
	name varchar(32) not null,
	description text
);

create table if not exists culture_event (
	culture_event_id serial not null primary key,
	culture_event_type_id int not null references culture_event_type (culture_event_type_id),
	culture_id int not null references culture (culture_id),
	recorded date not null default current_date,
	data json
);
