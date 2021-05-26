import React, { useContext, useEffect, useState } from "react";
import {
    DataProdukt,
    DataProduktResponse,
    DataProduktTilgangListe,
    DataProduktTilgangResponse, deleteAccess,
    hentTilganger,
    isOwner,
} from "./produktAPI";
import { UserContext } from "./userContext";
import "moment/locale/nb"
import moment from "moment";
import {Normaltekst, Systemtittel, Undertekst, Undertittel} from "nav-frontend-typografi";
import { Child } from "@navikt/ds-icons";
import "./produktDetaljTilganger.less"
import { Xknapp } from "nav-frontend-ikonknapper";

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
    if (isOwner(produkt?.data_product, userContext?.teams)) return true;
    // Ellers, vis kun dine egne tilganger.
    return subject === userContext?.email;
  };

  console.log(tilganger)
  if (!tilganger) return <></>;

  const tilgangsLinje = (tilgang: DataProduktTilgangResponse) => {
      const accessEnd = moment(tilgang.expires).format("LLL");

      return (
            <>
                <Child style={{ height: "100%", width: "auto" }} />
                <div className={"tilgangsTekst"}>
                <Undertittel>{tilgang.subject}</Undertittel>
                    <Undertekst>Innvilget av <em>{tilgang.author}</em> til {accessEnd}</Undertekst>
                </div>
                <Xknapp mini type={'fare'} onClick={async () => {
                    if(produkt?.id && tilgang?.subject)
                        await deleteAccess(produkt.id, tilgang.subject, 'user')
                }}/>
        </>
    );
  };

  const synligeTilganger = tilganger
    .filter((tilgang) => tilgang.action !== "verify")
    .filter((tilgang) => entryShouldBeDisplayed(tilgang.subject));


  let chronologicalEvents = synligeTilganger.sort((x, y) => (new Date(x.time).getTime() - new Date(y.time).getTime()));

  if (!synligeTilganger.length)
    return <p>Ingen relevante tilganger definert</p>;

  return (
      <div className={"tilgangsBoks"}>
          <Systemtittel>Tilganger</Systemtittel>
          <div className={"tilgangerContainer"}>
             {
                 synligeTilganger.map((tilgang, index) => (
                 <div key={index} className={"tilgangEntry"}>
                     {tilgangsLinje(tilgang)}
                 </div>
                 ))
             }
          </div>
      </div>
  );
};
