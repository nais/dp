import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import { Knapp, Fareknapp } from "nav-frontend-knapper";
import ModalWrapper from "nav-frontend-modal";
import { Ingress, Normaltekst, Systemtittel } from "nav-frontend-typografi";
import { DataProdukt } from "./produktAPI";
import NavFrontendSpinner from "nav-frontend-spinner";
import { Col, Container, Row } from "react-bootstrap";
import "./produktDetalj.less";

interface ProduktDetaljProps {
  produktID: string;
}

interface ProduktInfotabellProps {
  produkt: DataProdukt;
}

const ProduktInfoFaktaboks = ({ produkt }: ProduktInfotabellProps) => {
  return (
    <div className={"produkt-detaljer"}>
      <Systemtittel>{produkt.data_product?.name}</Systemtittel>

      <div className={"produkt-detaljer-faktaboks"}>
        <Ingress>
          Produkteier: {produkt.data_product?.owner || "uvisst"}
        </Ingress>
        <Normaltekst>
          URI: <code>{produkt.data_product?.uri || "uvisst"}</code>
        </Normaltekst>
        <Normaltekst>
          {produkt.data_product?.description || "uvisst"}
        </Normaltekst>
        <ul>
          <li>ID: {produkt.id}</li>
          <li>Opprettet: {produkt.created}</li>
          <li>Oppdatert: {produkt.updated}</li>
        </ul>
      </div>
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
  let history = useHistory();

  const deleteProduct = (id: string) => {
    fetch(`http://localhost:8080/dataproducts/${id}`, {
      method: "delete",
    })
      .then((res) => {
        if (res.status % 100 !== 2) {
          res.text().then((t) => {
            setError(t);
          });
        } else {
          history.push("/");
        }
      })
      .catch((e) => {
        setError(`Nettverksfeil: ${e}`);
        console.log(e);
      });
    console.log("delete this:", { id });
  };
  /*
<ModalWrapper
          isOpen={isOpen}
          onRequestClose={() => toggleOpen(false)}
          closeButton={true}
          contentLabel="Min modalrute"
      >.

        <div style={{ padding: "2rem 2.5rem" }}>Innhold her</div>

 */

  return (
    <div
      style={{ border: "1px solid black", display: isOpen ? "block" : "none" }}
    >
      <Systemtittel>Er du sikker?</Systemtittel>
      {error ? <p>{error}</p> : null}
      <Fareknapp onClick={() => deleteProduct(produktID)}>
        {error ? "Prøv igjen" : "Ja"}
      </Fareknapp>
    </div>
  );
};

export const ProduktTilganger = ({
  produkt,
}: ProduktInfotabellProps): JSX.Element => {
  return (
    <div className={"produkt-detaljer-tilganger"}>
      <Ingress>Tilganger</Ingress>
      <ul>
        {produkt.data_product.access.map((a) => (
          <li>{}</li>
        ))}
      </ul>
    </div>
  );
};

export const ProduktDetalj = (): JSX.Element => {
  let { produktID } = useParams<ProduktDetaljProps>();

  const [produkt, setProdukt] = useState<DataProdukt | undefined>(undefined);
  const [error, setError] = useState<string | null>();
  const [isOpen, setToggleOpen] = useState<boolean>(false);
  const toggleOpen = (input: boolean) => {
    setToggleOpen(input);
  };

  useEffect(() => {
    fetch(`http://localhost:8080/dataproducts/${produktID}`).then((res) => {
      if (!res.ok) {
        res.text().then((text) => setError(`HTTP ${res.status}: ${text}`));
      } else {
        res.json().then((json) => {
          setError(null);
          setProdukt(json);
        });
      }
    });
  }, [produktID]);

  if (error) return <div>{error}</div>;

  if (typeof produkt == "undefined")
    return (
      <div style={{ textAlign: "center" }}>
        <NavFrontendSpinner />
      </div>
    );

  return (
    <div>
      <Container fluid>
        <Row>
          <Col sm={3}>
            <div className="produktdetalj-knapper">
              <Knapp>Gi tilgang</Knapp>
              <Knapp>Fjern tilgang</Knapp>
              <Fareknapp onClick={() => toggleOpen(true)}>Slett</Fareknapp>
            </div>
          </Col>
          <Col sm={9}>
            {produkt ? <ProduktInfoFaktaboks produkt={produkt} /> : null}
            {produkt ? <ProduktTilganger produkt={produkt} /> : null}
            {produkt?.id ? (
              <SlettProdukt isOpen={isOpen} produktID={produkt.id} />
            ) : null}
          </Col>
        </Row>
      </Container>
    </div>
  );
};

export default ProduktDetalj;
