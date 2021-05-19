import React, { useEffect, useState, useContext } from "react";
import { useParams } from "react-router-dom";
import { Knapp, Fareknapp } from "nav-frontend-knapper";
import { GiTilgang, SlettProdukt } from "./produktTilgangModaler";
import { Normaltekst, Systemtittel } from "nav-frontend-typografi";
import {
  DataProdukt,
  DataProduktResponse,
  DataProduktTilgang,
  hentProdukt,
} from "./produktAPI";
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
  isOwner: boolean;
}

interface ProduktDetaljProps {
  setCrumb: React.Dispatch<React.SetStateAction<string | null>>;
}

const ProduktInfoFaktaboks = ({ produkt, isOwner }: ProduktInfotabellProps) => {
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
      <ProduktTilganger produkt={produkt.data_product} isOwner={isOwner} />
    </div>
  );
};

const ProduktTilganger: React.FC<{
  produkt: DataProdukt | null;
  isOwner: boolean;
}> = ({ produkt, isOwner }) => {
  const userContext = useContext(UserContext);
  const produktTilgang = (tilgang: DataProduktTilgang) => {
    const accessEnd = moment(tilgang.expires).format("LLL");

    return (
      <div className="innslag">
        {tilgang.subject}: til {accessEnd}
      </div>
    );
  };

  const entryShouldBeDisplayed = (
    entry: {
      subject: string;
      expires: Date | null;
    },
    isOwner: boolean
  ): boolean => {
    // Hvis produkteier, vis all tilgang;
    if (isOwner) return true;

    // Ellers, vis kun dine egne tilganger.
    return entry.subject === userContext?.email;
  };

  if (!produkt?.access) return <></>;

  return <div></div>;
};

export const ProduktDetalj = ({
  setCrumb,
}: ProduktDetaljProps): JSX.Element => {
  let { produktID } = useParams<ProduktDetaljParams>();
  const userContext = useContext(UserContext);

  const [produkt, setProdukt] = useState<DataProduktResponse | null>(null);
  const [error, setError] = useState<string | null>();
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [tilgangIsOpen, setTilgangIsOpen] = useState<boolean>(false);
  const [owner, setOwner] = useState<boolean>(false);

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
    hentProdukt(produktID)
      .then((p) => {
        setProdukt(p);
        setError(null);
      })
      .catch((e) => {
        setError(e.toString());
      });
  }, [produktID]);

  useEffect(() => {
    if (produkt != null) {
      setCrumb(produkt?.data_product.name || null);
    }
  }, [produkt, setCrumb]);

  useEffect(() => {
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
        callback={tilgangModalCallback}
        produkt={produkt}
      />

      <div className="produktdetalj">
        <ProduktInfoFaktaboks produkt={produkt} isOwner={owner} />
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
