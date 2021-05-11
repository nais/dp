import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import { Knapp, Fareknapp } from "nav-frontend-knapper";
import ModalWrapper from "nav-frontend-modal";
import {
  Ingress,
  Normaltekst,
  Sidetittel,
  Systemtittel,
  Undertittel,
} from "nav-frontend-typografi";
import { DataLager, DataProduktResponse } from "./produktAPI";
import NavFrontendSpinner from "nav-frontend-spinner";
import "./produktDetalj.less";

interface ProduktDetaljParams {
  produktID: string;
}

interface ProduktInfotabellProps {
  produkt: DataProduktResponse;
}

interface ProduktDetaljProps {
  setCrumb: React.Dispatch<React.SetStateAction<string | null>>;
}

const ProduktInfoFaktaboks = ({ produkt }: ProduktInfotabellProps) => {
  return (
    <div className={"faktaboks"}>
      <Systemtittel className={"produktnavn"}>
        {produkt.data_product?.name}
      </Systemtittel>
      <Undertittel>Produkteier:</Undertittel>{" "}
      <Normaltekst className="faktatekst">
        {produkt.data_product?.owner || "uvisst"}
      </Normaltekst>
      <Undertittel>Beskrivelse:</Undertittel>{" "}
      <Normaltekst className="faktatekst">
        {produkt.data_product?.description || "uvisst"}
      </Normaltekst>
      <Undertittel>Detaljer:</Undertittel>
      <table>
        <tr>
          <th>ID:</th>
          <td>{produkt.id}</td>
        </tr>
        <tr>
          <th>Opprettet:</th>
          <td>{produkt.created}</td>
        </tr>
        <tr>
          <th>Oppdatert:</th>
          <td>{produkt.updated}</td>
        </tr>
      </table>
    </div>
  );
};

interface SlettProduktProps {
  produktID: string;
  isOpen: boolean;
}

const SlettProdukt = ({
  produktID,
  isOpen,
}: SlettProduktProps): JSX.Element => {
  const toggleOpen = (f: boolean) => {};
  const [error, setError] = useState<string | null>(null);
  const history = useHistory();

  const deleteProduct = async (id: string) => {
    try {
      const res = await fetch(
        `http://localhost:8080/api/v1/dataproducts/${id}`,
        {
          method: "delete",
        }
      );

      if (res.status !== 204) {
        setError(`Feil: ${await res.text()}`);
      } else {
        history.push("/");
      }
    } catch (e) {
      setError(`Nettverksfeil: ${e}`);
      console.log(e);
    }
    console.log("delete this:", { id });
  };

  return (
    <ModalWrapper
      isOpen={isOpen}
      onRequestClose={() => toggleOpen(false)}
      closeButton={true}
      contentLabel="Min modalrute"
    >
      <div>
        <Systemtittel>Er du sikker?</Systemtittel>
        {error ? <p>{error}</p> : null}
        <Fareknapp onClick={() => deleteProduct(produktID)}>
          {error ? "Prøv igjen" : "Ja"}
        </Fareknapp>
      </div>
    </ModalWrapper>
  );
};
export const ProduktTilganger = ({
  produkt,
}: ProduktInfotabellProps): JSX.Element => {
  return (
    <div className={"produkt-detaljer-tilganger"}>
      <Ingress>Tilganger</Ingress>
      <ul>
        {produkt.data_product.access
          ? produkt.data_product.access.map((a) => <li>{}</li>)
          : null}
      </ul>
    </div>
  );
};

export const ProduktDetalj = ({
  setCrumb,
}: ProduktDetaljProps): JSX.Element => {
  let { produktID } = useParams<ProduktDetaljParams>();

  const [produkt, setProdukt] = useState<DataProduktResponse | null>(null);
  const [error, setError] = useState<string | null>();
  const [isOpen, setToggleOpen] = useState<boolean>(false);
  const toggleOpen = (input: boolean) => {
    setToggleOpen(input);
  };

  useEffect(() => {
    fetch(`http://localhost:8080/api/v1/dataproducts/${produktID}`).then(
      (res) => {
        if (!res.ok) {
          res.text().then((text) => setError(`HTTP ${res.status}: ${text}`));
        } else {
          res.json().then((json) => {
            setError(null);
            setProdukt(json);
          });
        }
      }
    );
  }, [produktID]);

  useEffect(() => {
    if (produkt != null) {
      setCrumb(produkt?.data_product.name || null);
    }
  }, [produkt]);

  if (error) return <div>{error}</div>;

  if (typeof produkt == "undefined")
    return (
      <div style={{ textAlign: "center" }}>
        <NavFrontendSpinner />
      </div>
    );

  if (produkt == null) return <></>;

  return (
    <div>
      <div className="produktdetalj">
        <div className="bokser">
          <ProduktInfoFaktaboks produkt={produkt} />
          <ProduktTilganger produkt={produkt} />
          <SlettProdukt isOpen={isOpen} produktID={produkt.id} />
        </div>
        <div className="knapperad">
          <Knapp>Gi tilgang</Knapp>
          <Knapp>Fjern tilgang</Knapp>
          <Fareknapp onClick={() => toggleOpen(true)}>Slett</Fareknapp>
        </div>
      </div>
    </div>
  );
};

export default ProduktDetalj;
