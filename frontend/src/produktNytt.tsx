import React, { useEffect, useState } from "react";
import { SkjemaGruppe, Input } from "nav-frontend-skjema";
import { Hovedknapp } from "nav-frontend-knapper";
import { Select } from "nav-frontend-skjema";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import { date } from "zod";

export const ProduktNytt = (): JSX.Element => {
  const [ressursType, setRessursType] = useState<string>("");
  const [navn, setNavn] = useState<string>("");
  const [beskrivelse, setBeskrivelse] = useState<string>("");
  const [eier, setEier] = useState<string>("");
  const [bigqueryView, setBigqueryView] = useState<string>("");
  const [bigqueryTable, setBigqueryTable] = useState<string>("");
  const [bucket, setBucket] = useState<string>("");

  useEffect(() => {
    setBigqueryTable("");
    setBigqueryView("");
    setBucket("");
  }, [ressursType]);

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

  const ressursVelger = (ressursType: string) => {
    switch (ressursType) {
      case "bigquery-view":
        return (
          <div>
            <Input
              label="view"
              value={bigqueryView}
              onChange={(e) => setBigqueryView(e.target.value)}
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
  return (
    <div>
      <SkjemaGruppe>
        <Input label="navn" onChange={(e) => setNavn(e.target.value)} />
        <Input
          label="beskrivelse"
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
          console.log(navn, beskrivelse, bucket);
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
