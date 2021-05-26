import React, { useEffect, useState, useContext } from "react";
import { useParams } from "react-router-dom";
import { Knapp, Fareknapp } from "nav-frontend-knapper";
import { GiTilgang, SlettProdukt } from "./produktTilgangModaler";
import { Feilmelding, Normaltekst, Systemtittel } from "nav-frontend-typografi";
import {
  DataProduktResponse,
  DataProduktTilgangListe,
  DataProduktTilgangResponse,
  hentProdukt,
  hentTilganger,
} from "./produktAPI";
import NavFrontendSpinner from "nav-frontend-spinner";
import "./produktDetalj.less";
import "moment/locale/nb";
import moment from "moment";
import { UserContext } from "./userContext";
import { DatalagerInfo } from "./produktDatalager";

interface ProduktDetaljParams {
  produktID: string;
}

interface ProduktInfotabellProps {
  produkt: DataProduktResponse;
  tilganger: DataProduktTilgangListe;
  isOwner: boolean;
}

interface ProduktDetaljProps {
  setCrumb: React.Dispatch<React.SetStateAction<string | null>>;
}

const ProduktInfoFaktaboks = ({
  tilganger,
  produkt,
  isOwner,
}: ProduktInfotabellProps) => {
  moment.locale("nb");

  return (
    <div className={"faktaboks"}>
      <Systemtittel className={"produktnavn"}>
        {produkt.data_product?.name}
      </Systemtittel>

      <Normaltekst>
        Produkteier: {produkt.data_product?.team || "uvisst"}
      </Normaltekst>
      <Normaltekst>
        Opprettet {moment(produkt.created).format("LLL")}
        {produkt.created !== produkt.updated
          ? ` (Oppdatert: ${moment(produkt.updated).fromNow()})`
          : ""}
      </Normaltekst>
      <Normaltekst className="beskrivelse">
        {produkt.data_product?.description || "Ingen beskrivelse"}
      </Normaltekst>
      {produkt.data_product?.datastore &&
        produkt.data_product?.datastore.map((ds) => <DatalagerInfo ds={ds} />)}

      <ProduktTilganger tilganger={tilganger} isOwner={isOwner} />
    </div>
  );
};

const ProduktTilganger: React.FC<{
  tilganger: DataProduktTilgangListe | null;
  isOwner: boolean;
}> = ({ tilganger, isOwner }) => {
  const userContext = useContext(UserContext);

  const produktTilgang = (tilgang: DataProduktTilgangResponse) => {
    const accessEnd = moment(tilgang.expires).format("LLL");

    return (
      <div className="innslag">
        {tilgang.subject}: til {tilgang.expires}
      </div>
    );
  };

  const entryShouldBeDisplayed = (
    subject: string | undefined,
    isOwner: boolean
  ): boolean => {
    // Hvis produkteier, vis all tilgang;
    if (isOwner) return true;

    // Ellers, vis kun dine egne tilganger.
    return subject === userContext?.email;
  };

  if (!tilganger) return <></>;

  const tilgangsLinje = (tilgang: DataProduktTilgangResponse) => {
    return (
      <tr>
        <td>{tilgang.author}</td>
        <td>{tilgang.action}</td>
        <td>{tilgang.subject}</td>
        <td>{tilgang.expires}</td>
      </tr>
    );
  };

  const synligeTilganger = tilganger
    .filter((tilgang) => tilgang.action !== "verify")
    .filter((tilgang) => entryShouldBeDisplayed(tilgang.subject, isOwner));

  if (!synligeTilganger.length)
    return <p>Ingen relevante tilganger definert</p>;

  return (
    <table>
      <tr>
        <th>author</th>
        <th>action</th>
        <th>subject</th>
        <th>expires</th>
      </tr>
      {synligeTilganger.map((tilgang) => tilgangsLinje(tilgang))}
    </table>
  );
};

export const ProduktDetalj = ({
  setCrumb,
}: ProduktDetaljProps): JSX.Element => {
  let { produktID } = useParams<ProduktDetaljParams>();
  const userContext = useContext(UserContext);

  const [produkt, setProdukt] = useState<DataProduktResponse | null>(null);
  const [tilganger, setTilganger] = useState<DataProduktTilgangListe | null>(
    null
  );
  const [error, setError] = useState<string | null>();
  const [tilgangerError, setTilgangerError] = useState<string | null>();
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
    hentTilganger(produktID)
      .then((p) => {
        setTilganger(p);
        setTilgangerError(null);
      })
      .catch((e) => {
        setTilgangerError(e.toString());
      });
  }, [produktID]);

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
      setOwner(userContext.teams.includes(produkt.data_product.team));
    }
  }, [produkt, userContext]);

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
        <ProduktInfoFaktaboks
          tilganger={tilganger}
          produkt={produkt}
          isOwner={owner}
        />
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
