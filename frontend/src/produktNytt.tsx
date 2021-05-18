import React, { useContext, useEffect, useState } from "react";
import { SkjemaGruppe, Input } from "nav-frontend-skjema";
import { Hovedknapp } from "nav-frontend-knapper";
import { Select } from "nav-frontend-skjema";
import { DataProduktSchema, DataLager, opprettProdukt } from "./produktAPI";
import "react-datepicker/dist/react-datepicker.css";
import { useHistory } from "react-router-dom";
import { UserContext } from "./userContext";

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
  const user = useContext(UserContext);

  const [navn, setNavn] = useState<string>("");
  const [beskrivelse, setBeskrivelse] = useState<string>("");
  const [eier, setEier] = useState<string>("");
  const [datastore, setDatastore] = useState<DataLager | null>(null);
  const history = useHistory();

  useEffect(() => setEier(user?.teams?.[0] || eier), [user]);

  const createProduct = async (): Promise<void> => {
    try {
      const nyttProdukt = DataProduktSchema.parse({
        name: navn,
        description: beskrivelse,
        datastore: [datastore],
        owner: eier,
        access: [],
      });
      const newID = await opprettProdukt(nyttProdukt);
      history.push(`/produkt/${newID}`);
    } catch (e) {
      console.log(e.toString());
    }
  };

  const validForm = () => {
    return true;
  };

  return (
    <div style={{ margin: "1em 1em 0 1em" }}>
      <SkjemaGruppe>
        <Input label="Navn" onChange={(e) => setNavn(e.target.value)} />
        <Input
          label="Beskrivelse"
          onChange={(e) => setBeskrivelse(e.target.value)}
        />
        <Select
          label="Eier (team)"
          onChange={(e) => setEier(e.target.value)}
          children={
            user?.teams
              ? user.teams.map((t) => (
                  <option key={t} value={t}>
                    {t}
                  </option>
                ))
              : null
          }
        />
        <Select
          label="Ressurstype"
          onChange={(e) => {
            if (e.target.value !== "")
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
        style={{ display: "block", marginLeft: "auto" }}
        disabled={!validForm()}
        onClick={async () => {
          await createProduct();
        }}
      >
        Submit
      </Hovedknapp>
    </div>
  );
};

export default ProduktNytt;
