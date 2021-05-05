import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  Ingress,
  Normaltekst,
  Sidetittel,
  Systemtittel,
  Undertittel,
} from "nav-frontend-typografi";
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
    <div className={"produkt-detaljer-faktaboks"}>
      <Systemtittel>{produkt.data_product?.name}</Systemtittel>
      <Ingress>Produkteier: {produkt.data_product?.owner || "uvisst"}</Ingress>
      <Normaltekst>
        URI: <code>{produkt.data_product?.uri || "uvisst"}</code>
      </Normaltekst>
      <Normaltekst>{produkt.data_product?.description || "uvisst"}</Normaltekst>
      <ul>
        <li>ID: {produkt.id}</li>
        <li>Opprettet: {produkt.created}</li>
        <li>Oppdatert: {produkt.updated}</li>
      </ul>
    </div>
  );
};

export const ProduktDetalj = (): JSX.Element => {
  let { produktID } = useParams<ProduktDetaljProps>();

  const [produkt, setProdukt] = useState<DataProdukt | undefined>(undefined);
  const [error, setError] = useState<string | null>();

  useEffect(() => {
    fetch(`http://localhost:8080/dataproducts/${produktID}`)
      .then((res) => res.json())
      .then((json) => {
        setError(null);
        setProdukt(json);
      });
  }, []);

  if (typeof produkt == "undefined")
    return (
      <div style={{ textAlign: "center" }}>
        <NavFrontendSpinner />
      </div>
    );
  if (error) return <div>{error}</div>;

  return (
    <div>
      <Sidetittel>Dataprodukt</Sidetittel>

      <Container fluid>
        <Row>
          <Col sm={3}>
            <div>Actions go here</div>
          </Col>
          <Col sm={9}>
            {produkt ? <ProduktInfoFaktaboks produkt={produkt} /> : null}
          </Col>
        </Row>
      </Container>
    </div>
  );
};
export default ProduktDetalj;
