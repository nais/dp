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
  ressurs: DataLager;
  setter: React.Dispatch<React.SetStateAction<DataLager>>;
}
const RessursVelger = ({ ressurs, setter }: RessursVelgerProps) => {
  switch (ressurs.type) {
    case "bucket":
      return (
        <div>
          <Input
            label="project_id"
            value={ressurs.project_id}
            onChange={(e) => setter({ ...ressurs, project_id: e.target.value })}
          />
          <Input
            label="bucket_id"
            value={ressurs.bucket_id}
            onChange={(e) => setter({ ...ressurs, bucket_id: e.target.value })}
          />
        </div>
      );
    case "bigquery-table":
      return (
        <div>
          <Input
            label="query"
            value={bigqueryTable}
            onChange={(e) => setBigqueryTable(e.target.value)}
          />
        </div>
      );
    case "bucket":
      return (
        <div>
          <Input
            label="bucket"
            value={bucket}
            onChange={(e) => setBucket(e.target.value)}
          />
        </div>
      );
    case "":
      return null;
  }
};
export const ProduktNytt = (): JSX.Element => {
  const [ressursType, setRessursType] = useState<string>("");
  const [navn, setNavn] = useState<string>("");
  const [beskrivelse, setBeskrivelse] = useState<string>("");
  const [eier, setEier] = useState<string>("");
  const [bigqueryView, setBigqueryView] = useState<object>({});
  const [bigqueryTable, setBigqueryTable] = useState<object>({});
  const [bucket, setBucket] = useState<string>("");
  const history = useHistory();

  useEffect(() => {
    setBigqueryTable("");
    setBigqueryView("");
    setBucket("");
  }, [ressursType]);

  const createProduct = async () => {
    const nyttProdukt = DataProduktSchema.parse({
      name: navn,
      description: beskrivelse,
      datastore: {
        type: ressursType,
        project_id: "placeholder",
        dataset_id: "placeholder",
      },
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
    if (
      navn?.length > 2 &&
      eier?.length > 2 &&
      (bigqueryView?.length > 2 ||
        bigqueryTable?.length > 2 ||
        bucket?.length > 2)
    )
      return false;
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
        <Select label="Type?" onChange={(e) => setRessursType(e.target.value)}>
          <option value="">Velg type</option>
          <option value="bigquery-view">BigQuery view</option>
          <option value="bigquery-table">BigQuery table</option>
          <option value="bucket">Bucket</option>
        </Select>
        {ressursVelger(ressursType)}
      </SkjemaGruppe>
      <Hovedknapp
        disabled={validForm()}
        onClick={() => {
          createProduct();
          //console.log(navn, beskrivelse, bucket);
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
