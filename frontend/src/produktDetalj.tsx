import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import { Knapp, Fareknapp } from "nav-frontend-knapper";
import Modal from "nav-frontend-modal";
import { Normaltekst, Systemtittel } from "nav-frontend-typografi";
import { DataProduktResponse, DataProduktTilgang } from "./produktAPI";
import NavFrontendSpinner from "nav-frontend-spinner";
import "./produktDetalj.less";
import "moment/locale/nb";
import moment from "moment";

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
  moment.locale("nb");

  return (
    <div className={"faktaboks"}>
      <Systemtittel className={"produktnavn"}>
        {produkt.data_product?.name}
      </Systemtittel>

      <Normaltekst className="beskrivelse">
        {produkt.data_product?.description || "Ingen beskrivelse"}
      </Normaltekst>
      <Normaltekst>
        Produkteier: {produkt.data_product?.owner || "uvisst"}
      </Normaltekst>

      <Normaltekst>
        Opprettet {moment(produkt.created).format("LLL")}
        {produkt.created !== produkt.updated
          ? ` (Oppdatert: ${moment(produkt.updated).fromNow()})`
          : ""}
      </Normaltekst>

      {produkt.data_product.access.map((a) => (
        <ProduktTilgang tilgang={a} />
      ))}
    </div>
  );
};

interface SlettProduktProps {
  produktID: string;
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

const SlettProdukt = ({
  produktID,
  isOpen,
  setIsOpen,
}: SlettProduktProps): JSX.Element => {
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
    <Modal
      isOpen={isOpen}
      onRequestClose={() => setIsOpen(false)}
      closeButton={true}
      contentLabel="Min modalrute"
    >
      <div className="slette-bekreftelse">
        <Systemtittel>Er du sikker?</Systemtittel>
        {error ? <p>{error}</p> : null}
        <Fareknapp onClick={() => deleteProduct(produktID)}>
          {error ? "Prøv igjen" : "Ja"}
        </Fareknapp>
      </div>
    </Modal>
  );
};

const ProduktTilgang: React.FC<{ tilgang: DataProduktTilgang }> = ({
  tilgang,
}) => {
  const accessStart = moment(tilgang.start).format("LLL");
  const accessEnd = moment(tilgang.end).format("LLL");

  return (
    <div className="innslag">
      {tilgang.subject}: {accessStart}&mdash;{accessEnd}
    </div>
  );
};

export const ProduktDetalj = ({
  setCrumb,
}: ProduktDetaljProps): JSX.Element => {
  let { produktID } = useParams<ProduktDetaljParams>();

  const [produkt, setProdukt] = useState<DataProduktResponse | null>(null);
  const [error, setError] = useState<string | null>();
  const [isOpen, setIsOpen] = useState<boolean>(false);

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
  }, [produkt, setCrumb]);

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
      <SlettProdukt
        isOpen={isOpen}
        setIsOpen={setIsOpen}
        produktID={produkt.id}
      />

      <div className="produktdetalj">
        <ProduktInfoFaktaboks produkt={produkt} />
        <div className="knapperad">
          <Knapp>Gi tilgang</Knapp>
          <Knapp>Fjern tilgang</Knapp>
          <Fareknapp onClick={() => setIsOpen(true)}>Slett</Fareknapp>
        </div>
      </div>
    </div>
  );
};

export default ProduktDetalj;
