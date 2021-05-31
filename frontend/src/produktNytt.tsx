import React, { useContext, useEffect, useState } from "react";
import { SkjemaGruppe, Input } from "nav-frontend-skjema";
import { Hovedknapp } from "nav-frontend-knapper";
import { Select } from "nav-frontend-skjema";
import {
  DataProduktSchema,
  DataLager,
  opprettProdukt,
  DataLagerBigquery,
  DataLagerBucket,
  DataLagerBigquerySchema,
  DataLagerBucketSchema,
} from "./produktAPI";
import "react-datepicker/dist/react-datepicker.css";
import { useHistory } from "react-router-dom";
import { UserContext } from "./userContext";
import { Feilmelding } from "nav-frontend-typografi";
import { ZodError } from "zod";

const RessursVelger: React.FC<{
  datastore: DataLager | null;
  setDatastore: React.Dispatch<React.SetStateAction<DataLager | null>>;
  parseErrors: ZodError["formErrors"] | undefined;
}> = ({ datastore, setDatastore, parseErrors }) => {
  const bucketFields = (datastore: DataLagerBucket) => (
    <>
      <Input
        label="project_id"
        feil={parseErrors?.fieldErrors?.project_id}
        value={datastore.project_id || ""}
        onChange={(e) =>
          setDatastore({ ...datastore, project_id: e.target.value })
        }
      />
      <Input
        label="bucket_id"
        feil={parseErrors?.fieldErrors?.bucket_id}
        value={datastore.bucket_id || ""}
        onChange={(e) =>
          setDatastore({ ...datastore, bucket_id: e.target.value })
        }
      />
    </>
  );

  const bigqueryFields = (datastore: DataLagerBigquery) => (
    <>
      <Input
        label="Project ID"
        feil={parseErrors?.fieldErrors?.project_id}
        value={datastore.project_id || ""}
        onChange={(e) =>
          setDatastore({ ...datastore, project_id: e.target.value })
        }
      />
      <Input
        label="Resource ID"
        feil={parseErrors?.fieldErrors?.resource_id}
        value={datastore.resource_id || ""}
        onChange={(e) =>
          setDatastore({ ...datastore, resource_id: e.target.value })
        }
      />
      <Input
        label="Dataset ID"
        feil={parseErrors?.fieldErrors?.dataset_id}
        value={datastore.dataset_id || ""}
        onChange={(e) =>
          setDatastore({ ...datastore, dataset_id: e.target.value })
        }
      />
    </>
  );

  return (
    <>
      <Select
        feil={parseErrors?.fieldErrors?.type}
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
      {datastore?.type === "bigquery" && bigqueryFields(datastore)}
      {datastore?.type === "bucket" && bucketFields(datastore)}
    </>
  );
};

const ProduktSkjema: React.FC<{
    produktID: string | null;
}> = ({ produktID}) => {
    const history = useHistory();
    const user = useContext(UserContext);
    const [navn, setNavn] = useState<string>("");
    const [beskrivelse, setBeskrivelse] = useState<string>("");
    const [eier, setEier] = useState<string>("");
    const [datastore, setDatastore] = useState<DataLager | null>(null);
    const [formErrors, setFormErrors] = useState<ZodError["formErrors"]>();
    const [dsFormErrors, setDsFormErrors] = useState<ZodError["formErrors"]>();
    // Correctly initialize form to default value on page render
    useEffect(() => setEier((e) => user?.teams?.[0] || e), [user]);

    const parseDatastore = () => {
        if (datastore?.type === "bigquery") {
            DataLagerBigquerySchema.parse(datastore);
        } else if (datastore?.type === "bucket") {
            DataLagerBucketSchema.parse(datastore);
        } else {
            setDsFormErrors({
                formErrors: [],
                fieldErrors: { type: ["Required"] },
            });
        }
    };

    const handleSubmit = async (): Promise<void> => {
        // First make sure we have a valid datastore
        try {
            parseDatastore();
        } catch (e) {
            setDsFormErrors(e.flatten());
        }
        try {
            const nyttProdukt = DataProduktSchema.parse({
                name: navn,
                description: beskrivelse,
                datastore: [datastore],
                team: eier,
                access: {},
            });
            const newID = await opprettProdukt(nyttProdukt);
            history.push(`/produkt/${newID}`);
        } catch (e) {
            console.log(e.toString());

            if (e instanceof ZodError) {
                setFormErrors(e.flatten());
            } else {
                setFormErrors({ formErrors: e.toString(), fieldErrors: {} });
            }
        }
    };
    return (
        <div style={{ margin: "1em 1em 0 1em" }}>
            <SkjemaGruppe>
                <Input
                    label="Navn"
                    feil={formErrors?.fieldErrors?.name}
                    onChange={(e) => setNavn(e.target.value)}
                />
                <Input
                    label="Beskrivelse"
                    feil={formErrors?.fieldErrors?.description}
                    onChange={(e) => setBeskrivelse(e.target.value)}
                />
                <Select
                    label="Eier (team)"
                    feil={formErrors?.fieldErrors?.owner}
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

                <RessursVelger
                    datastore={datastore}
                    setDatastore={setDatastore}
                    parseErrors={dsFormErrors}
                />
            </SkjemaGruppe>
            {!!formErrors?.formErrors?.length && (
                <Feilmelding>
                    Feil: <br />
                    <code>{formErrors.formErrors}</code>
                </Feilmelding>
            )}
            <Hovedknapp
                style={{ display: "block", marginLeft: "auto" }}
                onClick={async () => {
                    await handleSubmit();
                }}
            >
                Submit
            </Hovedknapp>
        </div>
    );
}

export const ProduktNytt = (): JSX.Element => {
  return <ProduktSkjema produktID={""}></ProduktSkjema>
};
export const ProduktOppdatering = (): JSX.Element => {
    return <ProduktSkjema produktID={""}></ProduktSkjema>
};

export default ProduktNytt;
