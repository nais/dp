import React, { useEffect, useState, useContext } from "react";
import { useParams, useHistory } from "react-router-dom";
import { Knapp, Fareknapp } from "nav-frontend-knapper";
import Modal from "nav-frontend-modal";
import { Normaltekst, Systemtittel } from "nav-frontend-typografi";
import { DataProduktResponse, DataProduktTilgang } from "./produktAPI";
import NavFrontendSpinner from "nav-frontend-spinner";
import "./produktDetalj.less";
import "moment/locale/nb";
import moment from "moment";
import { UserContext } from "./userContext";

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
      <ProduktTilganger tilganger={produkt.data_product.access} />
    </div>
  );
};

interface SlettProduktProps {
  produktID: string;
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

interface GiTilgangProps {
  produkt: DataProduktResponse;
  tilgangIsOpen: boolean;
  setTilgangIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

const GiTilgang = ({
  produkt,
  tilgangIsOpen,
  setTilgangIsOpen,
}: GiTilgangProps): JSX.Element => {
  return (
    <Modal
      isOpen={tilgangIsOpen}
      onRequestClose={() => setTilgangIsOpen(false)}
      closeButton={true}
      contentLabel="Gi tilgang"
    >
      <div className="slette-bekreftelse">
        <Systemtittel>Gi tilgang</Systemtittel>
      </div>
    </Modal>
  );
};

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

const ProduktTilganger: React.FC<{
  tilganger: DataProduktTilgang[] | null;
}> = ({ tilganger }) => {
  const produktTilgang = (tilgang: DataProduktTilgang) => {
    const accessStart = moment(tilgang.start).format("LLL");
    const accessEnd = moment(tilgang.end).format("LLL");

    return (
      <div className="innslag">
        {tilgang.subject}: {accessStart}&mdash;{accessEnd}
      </div>
    );
  };

  if (tilganger == null) return <></>;

  return <div>{tilganger.map((a) => produktTilgang(a))}</div>;
};

export const ProduktDetalj = ({
  setCrumb,
}: ProduktDetaljProps): JSX.Element => {
  let { produktID } = useParams<ProduktDetaljParams>();

  const [produkt, setProdukt] = useState<DataProduktResponse | null>(null);
  const [error, setError] = useState<string | null>();
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [tilgangIsOpen, setTilgangIsOpen] = useState<boolean>(false);
  const [owner, setOwner] = useState<boolean>(false);
  const userContext = useContext(UserContext);

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

  useEffect(() => {
    console.log(userContext);
    if (produkt && userContext) {
      setOwner(userContext.teams.includes(produkt.data_product.owner));
    }
  }, [produkt, userContext]);

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

      <GiTilgang
        tilgangIsOpen={tilgangIsOpen}
        setTilgangIsOpen={setTilgangIsOpen}
        produkt={produkt}
      />

      <div className="produktdetalj">
        <ProduktInfoFaktaboks produkt={produkt} />
        <div className="knapperad">
          {owner ? (
            <Fareknapp onClick={() => setIsOpen(true)}>Slett</Fareknapp>
          ) : userContext ? (
            <Knapp onClick={() => setTilgangIsOpen(true)}>Få tilgang</Knapp>
          ) : null}
        </div>
      </div>
    </div>
  );
};

export default ProduktDetalj;
