/*
 * Copyright (c) 2019-2020 Datacequia LLC. All rights reserved.
 *
 * This program is licensed to you under the Apache License Version 2.0,
 * and you may not use this file except in compliance with the Apache License Version 2.0.
 * You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the Apache License Version 2.0 is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
 */

package cmd

import (
	"github.com/datacequia/go-dogg3rz/resource"
	//	"github.com/datacequia/go-dogg3rz/util"
)

type dgrzCreateType struct {
	Class    dgrzCreateTypeClass    `command:"class"  description:"create an instance of an RDF Schema Class (rdfs:Class)"`
	Datatype dgrzCreateTypeDatatype `command:"datatype" description:"create an instance of an RDF Schema Datatype (rdfs:Datatype)"`
	Property dgrzCreateTypeProperty `command:"property" description:"create an instance of an RDF Property (rdf:Property)"`
}

type dgrzTypeOptions struct {
	Comment string `long:"comment" description:"RDF Schema label" required:"false"`
	Label   string `long:"label" description:"RDF Schema comment" required:"false"`
}

type dgrzTypeContext struct {
	DatasetPath string `positional-arg-name:"DATASET_PATH" description:"repository path to dataset" required:"yes"`
	ID          string `positional-arg-name:"ID" description:"[relative] IRI or term" required:"yes" `
}

type dgrzClassOptions struct {
	dgrzTypeOptions
	SubclassOf string `long:"subclass-of" description:"parent class" required:"false"`
}

type dgrzPropertyOptions struct {
	dgrzTypeOptions
	SubpropertyOf string `long:"subproperty-of" description:"parent class" required:"false"`
	Domain        string `long:"domain" description:"parent class" required:"true"`
	Range         string `long:"range" description:"parent class" required:"true"`
}

type dgrzCreateTypeClass struct {
	Positional dgrzTypeContext `positional-args:"yes"`
	Options    dgrzClassOptions
}

type dgrzCreateTypeDatatype struct {
	Positional dgrzTypeContext `positional-args:"yes"`
	Options    dgrzClassOptions
}

type dgrzCreateTypeProperty struct {
	Positional dgrzTypeContext `positional-args:"yes"`
	Options    dgrzPropertyOptions
}

///////////////////////////////////////////////////////////
// CREATE TYPE  FUNCTIONS
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
// CREATE TYPE CLASS SUBCOMMAND FUNCTIONS
///////////////////////////////////////////////////////////
func (o *dgrzCreateTypeClass) Execute(args []string) error {

	ctxt := getCmdContext()

	repo := resource.GetRepositoryResource(ctxt)

	if err := repo.CreateTypeClass(ctxt,
		createCmd.Repository,
		o.Positional.DatasetPath,
		o.Positional.ID, o.Options.SubclassOf,
		o.Options.Label, o.Options.Comment); err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////
// CREATE TYPE DATATYPE SUBCOMMAND FUNCTIONS
///////////////////////////////////////////////////////////
func (o *dgrzCreateTypeDatatype) Execute(args []string) error {

	ctxt := getCmdContext()

	repo := resource.GetRepositoryResource(ctxt)

	if err := repo.CreateTypeDatatype(ctxt,
		createCmd.Repository, o.Positional.DatasetPath,
		o.Positional.ID, o.Options.SubclassOf,
		o.Options.Label, o.Options.Comment); err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////
// CREATE TYPE PROPERTY SUBCOMMAND FUNCTIONS
///////////////////////////////////////////////////////////
func (o *dgrzCreateTypeProperty) Execute(args []string) error {

	ctxt := getCmdContext()
	// TODO: ALLOW FOR MULTIPLE RANGES (AND DOMAINS) DURING CREATION
	repo := resource.GetRepositoryResource(ctxt)

	if err := repo.CreateTypeProperty(ctxt,
		createCmd.Repository, o.Positional.DatasetPath,
		o.Positional.ID, o.Options.SubpropertyOf,
		o.Options.Domain, o.Options.Range,
		o.Options.Label, o.Options.Comment); err != nil {

		return err
	}

	return nil
}
