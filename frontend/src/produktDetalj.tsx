import React, { useEffect, useState, useContext } from "react";
import { useParams } from "react-router-dom";
import { GiTilgang, SlettProdukt } from "./produktTilgangModaler";
import { Feilmelding } from "nav-frontend-typografi";
import {
  DataProduktResponse,
  DataProduktTilgangListe,
  hentProdukt,
  hentTilganger,
} from "./produktAPI";
import NavFrontendSpinner from "nav-frontend-spinner";
import "./produktDetalj.less";
import { UserContext } from "./userContext";
import { ProduktInfoFaktaboks } from "./produktDetaljFaktaboks";
import { ProduktKnapperad } from "./produktDetaljKnapperad";

interface ProduktDetaljParams {
  produktID: string;
}

export const ProduktDetalj: React.FC<{
  setCrumb: React.Dispatch<React.SetStateAction<string | null>>;
}> = ({ setCrumb }): JSX.Element => {
  let { produktID } = useParams<ProduktDetaljParams>();
  const [tilgangIsOpen, setTilgangIsOpen] = useState<boolean>(false);
  const [isOpen, setIsOpen] = useState<boolean>(false);

  const [produkt, setProdukt] = useState<DataProduktResponse | null>(null);
  const [error, setError] = useState<string | null>();
  const [tilganger, setTilganger] = useState<DataProduktTilgangListe | null>(
    null
  );
  const [tilgangerError, setTilgangerError] = useState<string | null>();

  useEffect(() => {
    if (!produkt) return;
    hentTilganger(produkt.id)
      .then((p) => {
        setTilganger(p);
        setTilgangerError(null);
      })
      .catch((e) => {
        setTilgangerError(e.toString());
      });
  }, [produkt]);

  useEffect(() => {
    hentProdukt(produktID)
      .then((p) => {
        setProdukt(p);
        setError(null);
      })
      .catch((e) => {
        setError(e.toString());
      });
  }, [produktID]);

  const tilgangModalCallback = (hasChanged: boolean) => {
    setTilgangIsOpen(false);
    if (hasChanged) {
      // TODO: Use hooks more elegantly
      hentProdukt(produktID)
        .then((p) => {
          setProdukt(p);
          setError(null);
        })
        .catch((e) => {
          setError(e.toString());
        });
    }
  };

  useEffect(() => {
    if (produkt != null) {
      setCrumb(produkt?.data_product.name || null);
    }
  }, [produkt, setCrumb]);

  if (error)
    return (
      <Feilmelding>
        <code>{error}</code>
      </Feilmelding>
    );

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
        <SlettProdukt
          isOpen={isOpen}
          setIsOpen={setIsOpen}
          produktID={produkt.id}
        />

        <GiTilgang
          tilgangIsOpen={tilgangIsOpen}
          callback={tilgangModalCallback}
          produkt={produkt}
        />
        <ProduktInfoFaktaboks tilganger={tilganger} produkt={produkt} />
        <ProduktKnapperad
          produkt={produkt}
          tilganger={tilganger}
          openSlett={() => setIsOpen(true)}
          openTilgang={() => setTilgangIsOpen(true)}
        />
      </div>
    </div>
  );
};

export default ProduktDetalj;
