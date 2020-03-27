/*
 *  Dogg3rz is a decentralized metadata version control system
 *  Copyright (C) 2019 D. Andrew Padilla dba Datacequia
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package config

import (
	"io/ioutil"
	"os"

	"github.com/datacequia/go-dogg3rz/resource/config"
)

func SetConfigDefault(c config.Dogg3rzConfig) error {

	var dgrzConfS string
	var err error
	dgrzConfS, err = config.GenerateDefault(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath(), []byte(dgrzConfS), os.FileMode(0660))

	return err

}
