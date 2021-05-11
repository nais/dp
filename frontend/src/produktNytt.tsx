import React, { useEffect, useState } from "react";
import { SkjemaGruppe, Input } from "nav-frontend-skjema";
import { Hovedknapp } from "nav-frontend-knapper";
import { Select } from "nav-frontend-skjema";
import { DataProduktSchema, DataLager } from "./produktAPI";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import { date } from "zod";
import { useHistory } from "react-router-dom";

interface RessursVelgerProps {
  ressurs: DataLager | null;
  setter: React.Dispatch<React.SetStateAction<DataLager | null>>;
}

const RessursVelger = ({ ressurs, setter }: RessursVelgerProps) => {
  if (ressurs === null) {
    return null;
  }

  switch (ressurs.type) {
    case "bucket":
      return (
        <div>
          <Input
            label="project_id"
            value={ressurs.project_id || ""}
            onChange={(e) => setter({ ...ressurs, project_id: e.target.value })}
          />
          <Input
            label="bucket_id"
            value={ressurs.bucket_id || ""}
            onChange={(e) => setter({ ...ressurs, bucket_id: e.target.value })}
          />
        </div>
      );
    case "bigquery":
      return (
        <div>
          <Input
            label="Project ID"
            value={ressurs.project_id || ""}
            onChange={(e) => setter({ ...ressurs, project_id: e.target.value })}
          />
          <Input
            label="Resource ID"
            value={ressurs.resource_id || ""}
            onChange={(e) =>
              setter({ ...ressurs, resource_id: e.target.value })
            }
          />
          <Input
            label="Dataset ID"
            value={ressurs.dataset_id || ""}
            onChange={(e) => setter({ ...ressurs, dataset_id: e.target.value })}
          />
        </div>
      );
  }
};
export const ProduktNytt = (): JSX.Element => {
  const [navn, setNavn] = useState<string>("");
  const [beskrivelse, setBeskrivelse] = useState<string>("");
  const [eier, setEier] = useState<string>("");
  const [datastore, setDatastore] = useState<DataLager | null>(null);
  const history = useHistory();

  const createProduct = async () => {
    const nyttProdukt = DataProduktSchema.parse({
      name: navn,
      description: beskrivelse,
      datastore: [datastore],
      owner: eier,
      access: [],
    });

    const res = await fetch("http://localhost:8080/api/v1/dataproducts", {
      method: "POST",
      body: JSON.stringify(nyttProdukt),
    });
    const newID = await res.text();
    history.push(`/produkt/${newID}`);
  };

  const validForm = () => {
    return true;
  };

  return (
    <div>
      <SkjemaGruppe>
        <Input label="Navn" onChange={(e) => setNavn(e.target.value)} />
        <Input
          label="Beskrivelse"
          onChange={(e) => setBeskrivelse(e.target.value)}
        />
        <Input label="Eier (team)" onChange={(e) => setEier(e.target.value)} />
        <p>Ressurs</p>
        <Select
          label="Type?"
          onChange={(e) => {
            if (e.target.value != "")
              setDatastore({ type: e.target.value } as DataLager);
            else setDatastore(null);
          }}
        >
          <option value="">Velg type</option>
          <option value="bigquery">BigQuery</option>
          <option value="bucket">Bucket</option>
        </Select>
        <RessursVelger ressurs={datastore} setter={setDatastore} />
      </SkjemaGruppe>
      <Hovedknapp
        disabled={!validForm()}
        onClick={() => {
          createProduct();
        }}
      >
        Submit
      </Hovedknapp>
    </div>
  );
};
/*

	Subject string    `firestore:"subject" json:"subject,omitempty" validate:"required"`
	Start   time.Time `firestore:"start" json:"start,omitempty" validate:"required"`
	End     time.Time `firestore:"end" json:"end,omitempty" validate:"required"`

	Name        string         `firestore:"name" json:"name,omitempty" validate:"required"`
	Description string         `firestore:"description" json:"description,omitempty" validate:"required"`
	Resource    Resource       `firestore:"resource" json:"resource,omitempty" validate:"required"`
	Owner       string         `firestore:"owner" json:"owner,omitempty" validate:"required"`
	Access      []*AccessEntry `firestore:"access" json:"access,omitempty" validate:"required,dive"`
 */
export default ProduktNytt;
