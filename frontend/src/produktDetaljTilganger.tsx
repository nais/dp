import React, { useContext, useEffect, useState } from "react";
import {
  DataProdukt,
  DataProduktResponse,
  DataProduktTilgangListe,
  DataProduktTilgangResponse,
  hentTilganger,
  isOwner,
} from "./produktAPI";
import { UserContext } from "./userContext";
import moment from "moment";
import { Systemtittel } from "nav-frontend-typografi";
import { Child } from "@navikt/ds-icons";

export const ProduktTilganger: React.FC<{
  produkt: DataProduktResponse | null;
  tilganger: DataProduktTilgangListe | null;
}> = ({ produkt, tilganger }) => {
  const userContext = useContext(UserContext);

  const produktTilgang = (tilgang: DataProduktTilgangResponse) => {
    const accessEnd = moment(tilgang.expires).format("LLL");

    return (
      <div className="innslag">
        {tilgang.subject}: til {tilgang.expires}
      </div>
    );
  };

  const entryShouldBeDisplayed = (subject: string | undefined): boolean => {
    if (!produkt?.data_product || !userContext?.teams) return false;
    // Hvis produkteier, vis all tilgang;
    return isOwner(produkt?.data_product, userContext?.teams);
    // Ellers, vis kun dine egne tilganger.
    return subject === userContext?.email;
  };

  if (!tilganger) return <></>;

  const tilgangsLinje = (tilgang: DataProduktTilgangResponse) => {
    return (
      <ul>
        <li>{tilgang.author}</li>
        <li>{tilgang.action}</li>
        <li>{tilgang.subject}</li>
        <li>{tilgang.expires}</li>
      </ul>
    );
  };

  const synligeTilganger = tilganger
    .filter((tilgang) => tilgang.action !== "verify")
    .filter((tilgang) => entryShouldBeDisplayed(tilgang.subject));

  if (!synligeTilganger.length)
    return <p>Ingen relevante tilganger definert</p>;

  return (
    <div className={"datalagerBoks"}>
      <Systemtittel>Tilganger</Systemtittel>
      <div className={"datalagerentry"}>
        <Child style={{ height: "100%", width: "auto" }} />
        {synligeTilganger.map((tilgang) => tilgangsLinje(tilgang))}
      </div>
    </div>
  );
};
